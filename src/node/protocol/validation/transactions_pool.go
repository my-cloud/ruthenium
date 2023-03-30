package validation

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"math/rand"
	"sync"
	"time"
)

type TransactionsPool struct {
	transactions         []*Transaction
	transactionResponses []*network.TransactionResponse
	mutex                sync.RWMutex

	blockchain            protocol.Blockchain
	minimalTransactionFee uint64
	synchronizer          network.Synchronizer
	validatorAddress      string

	validationTimer time.Duration
	watch           clock.Watch

	logger log.Logger
}

func NewTransactionsPool(blockchain protocol.Blockchain, minimalTransactionFee uint64, synchronizer network.Synchronizer, validatorAddress string, validationTimer time.Duration, logger log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blockchain = blockchain
	pool.minimalTransactionFee = minimalTransactionFee
	pool.synchronizer = synchronizer
	pool.validatorAddress = validatorAddress
	pool.validationTimer = validationTimer
	pool.watch = tick.NewWatch()
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transactionRequest *network.TransactionRequest, hostTarget string) {
	err := pool.addTransaction(transactionRequest)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
		return
	}
	if transactionRequest.TransactionBroadcasterTarget != nil {
		pool.synchronizer.Incentive(*transactionRequest.TransactionBroadcasterTarget)
	}
	transactionRequest.TransactionBroadcasterTarget = &hostTarget
	neighbors := pool.synchronizer.Neighbors()
	for _, neighbor := range neighbors {
		go func(neighbor network.Neighbor) {
			_ = neighbor.AddTransaction(*transactionRequest)
		}(neighbor)
	}
}

func (pool *TransactionsPool) Transactions() []*network.TransactionResponse {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()
	return pool.transactionResponses
}

func (pool *TransactionsPool) Validate(timestamp int64) {
	currentBlockchain := pool.blockchain.Copy()
	blockResponses := currentBlockchain.Blocks()
	lastBlockResponse := blockResponses[len(blockResponses)-1]
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	isAlreadySpentByInput := make(map[*network.InputResponse]bool)
	transactionResponses := pool.transactionResponses
	rand.Seed(pool.watch.Now().UnixNano())
	rand.Shuffle(len(transactionResponses), func(i, j int) {
		transactionResponses[i], transactionResponses[j] = transactionResponses[j], transactionResponses[i]
	})
	var reward uint64
	var rejectedTransactions []*network.TransactionResponse
	for i, transaction := range transactionResponses {
		if timestamp+pool.validationTimer.Nanoseconds() < transaction.Timestamp {
			pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction timestamp is too far in the future, transaction: %v", transaction))
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		if len(blockResponses) > 1 {
			if transaction.Timestamp < blockResponses[len(blockResponses)-1].Timestamp {
				pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction timestamp is too old, transaction: %v", transaction))
				rejectedTransactions = append(rejectedTransactions, transaction)
				continue
			}
			var skip bool
			for _, validatedTransaction := range blockResponses[len(blockResponses)-1].Transactions {
				if pool.transactions[i].Equals(validatedTransaction) {
					pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction is already in the blockchain, transaction: %v", transaction))
					rejectedTransactions = append(rejectedTransactions, transaction)
					skip = true
					continue
				}
			}
			if skip {
				continue
			}
			fee := currentBlockchain.CalculateFee(transaction, timestamp)
			for _, input := range transaction.Inputs {
				if _, ok := isAlreadySpentByInput[input]; ok {
					pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, an input has already been spent, transaction: %v", transaction))
					rejectedTransactions = append(rejectedTransactions, transaction)
					skip = true
					continue
				}
				isAlreadySpentByInput[input] = true
			}
			if skip {
				continue
			}
			reward += fee
		}
	}
	for _, transaction := range rejectedTransactions {
		transactionResponses = removeTransaction(transactionResponses, transaction)
	}
	var newAddresses []string
	for _, transaction := range transactionResponses {
		for _, output := range transaction.Outputs {
			newAddresses = append(newAddresses, output.Address)
		}
	}
	if lastBlockResponse.Timestamp == timestamp {
		pool.logger.Error("unable to create block, a block with the same timestamp is already in the blockchain")
		return
	}
	rewardTransaction := NewRewardTransaction(pool.validatorAddress, len(blockResponses), timestamp, reward)
	transactionResponses = append(transactionResponses, rewardTransaction)
	err := pool.blockchain.AddBlock(timestamp, transactionResponses, newAddresses)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block: %w", err).Error())
		return
	}
	pool.clear()
	pool.logger.Debug(fmt.Sprintf("reward: %d", reward))
}

func (pool *TransactionsPool) addTransaction(transactionRequest *network.TransactionRequest) error {
	currentBlockchain := pool.blockchain.Copy()
	blocks := currentBlockchain.Blocks()
	if len(blocks) == 0 {
		return errors.New("the blockchain is empty")
	}
	transaction, err := NewTransactionFromRequest(transactionRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate transaction: %w", err)
	}
	currentBlock := blocks[len(blocks)-1]
	transactionResponse := transaction.GetResponse()
	fee := currentBlockchain.CalculateFee(transactionResponse, currentBlock.Timestamp+pool.validationTimer.Nanoseconds())
	if fee < pool.minimalTransactionFee {
		return fmt.Errorf("the transaction fee is too low, fee: %d, minimal fee: %d", fee, pool.minimalTransactionFee)
	}
	if len(blocks) > 1 {
		timestamp := transaction.Timestamp()
		nextBlockTimestamp := currentBlock.Timestamp + 2*pool.validationTimer.Nanoseconds()
		if nextBlockTimestamp < timestamp {
			return fmt.Errorf("the transaction timestamp is too far in the future: %v, now: %v", time.Unix(0, timestamp), time.Unix(0, nextBlockTimestamp))
		}
		currentBlockTimestamp := currentBlock.Timestamp
		if timestamp < currentBlockTimestamp {
			return fmt.Errorf("the transaction timestamp is too old: %v, current block timestamp: %v", time.Unix(0, timestamp), time.Unix(0, currentBlockTimestamp))
		}
		for _, validatedTransaction := range currentBlock.Transactions {
			if transaction.Equals(validatedTransaction) {
				return errors.New("the transaction is already in the blockchain")
			}
		}
	}
	for _, pendingTransaction := range pool.transactionResponses {
		if transaction.Equals(pendingTransaction) {
			return errors.New("the transaction is already in the transactions pool")
		}
	}
	if err = transaction.VerifySignatures(); err != nil {
		return errors.New("failed to verify transaction")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.transactions = append(pool.transactions, transaction)
	pool.transactionResponses = append(pool.transactionResponses, transactionResponse)
	return nil
}

func (pool *TransactionsPool) clear() {
	pool.transactions = nil
	pool.transactionResponses = nil
}

func removeTransaction(transactions []*network.TransactionResponse, removedTransaction *network.TransactionResponse) []*network.TransactionResponse {
	for i := 0; i < len(transactions); i++ {
		if transactions[i] == removedTransaction {
			transactions = append(transactions[:i], transactions[i+1:]...)
			return transactions
		}
	}
	return transactions
}
