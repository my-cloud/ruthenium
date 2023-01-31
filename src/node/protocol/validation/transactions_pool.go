package validation

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
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
	registry              protocol.Registry
	validatorAddress      string

	validationTimer time.Duration
	watch           clock.Watch

	logger log.Logger
}

func NewTransactionsPool(blockchain protocol.Blockchain, minimalTransactionFee uint64, registry protocol.Registry, validatorAddress string, validationTimer time.Duration, watch clock.Watch, logger log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blockchain = blockchain
	pool.minimalTransactionFee = minimalTransactionFee
	pool.registry = registry
	pool.validatorAddress = validatorAddress
	pool.validationTimer = validationTimer
	pool.watch = watch
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transactionRequest *network.TransactionRequest, neighbors []network.Neighbor) {
	err := pool.addTransaction(transactionRequest)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
		return
	}
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
	totalTransactionsValueBySenderAddress := make(map[string]uint64)
	transactionResponses := pool.transactionResponses
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
			for _, validatedTransaction := range blockResponses[len(blockResponses)-1].Transactions {
				if pool.transactions[i].Equals(validatedTransaction) {
					pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction is already in the blockchain, transaction: %v", transaction))
					rejectedTransactions = append(rejectedTransactions, transaction)
					continue
				}
			}
		}
	}
	for _, transaction := range rejectedTransactions {
		transactionResponses = removeTransaction(transactionResponses, transaction)
	}
	for _, transaction := range transactionResponses {
		fee := transaction.Fee
		totalTransactionsValueBySenderAddress[transaction.SenderAddress] += transaction.Value + fee
		reward += fee
	}
	var newAddresses []string
	for senderAddress, totalTransactionsValue := range totalTransactionsValueBySenderAddress {
		senderTotalAmount := currentBlockchain.CalculateTotalAmount(timestamp, senderAddress)
		if totalTransactionsValue > senderTotalAmount {
			rejectedTransactions = nil
			rand.Seed(pool.watch.Now().UnixNano())
			rand.Shuffle(len(transactionResponses), func(i, j int) {
				transactionResponses[i], transactionResponses[j] = transactionResponses[j], transactionResponses[i]
			})
			for _, transaction := range transactionResponses {
				if transaction.SenderAddress == senderAddress {
					rejectedTransactions = append(rejectedTransactions, transaction)
					fee := transaction.Fee
					totalTransactionsValue -= transaction.Value + fee
					reward -= fee
					if totalTransactionsValue <= senderTotalAmount {
						break
					}
				}
			}
			for _, transaction := range rejectedTransactions {
				transactionResponses = removeTransaction(transactionResponses, transaction)
				pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, total transactions value exceeds its sender wallet amount, transaction: %v", transaction))
			}
		}
		if totalTransactionsValue > 0 {
			newAddresses = append(newAddresses, senderAddress)
		}
	}
	if lastBlockResponse.Timestamp == timestamp {
		pool.logger.Error("unable to create block, a block with the same timestamp is already in the blockchain")
		return
	}
	rewardTransaction := NewRewardTransaction(pool.validatorAddress, timestamp, reward)
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
	fee := *transactionRequest.Fee
	if fee < pool.minimalTransactionFee {
		return fmt.Errorf("the transaction fee is too low, fee: %d, minimal fee: %d", fee, pool.minimalTransactionFee)
	}
	transaction, err := NewTransactionFromRequest(transactionRequest)
	if err != nil {
		return fmt.Errorf("failed to instantiate transaction: %w", err)
	}
	currentBlockchain := pool.blockchain.Copy()
	blocks := currentBlockchain.Blocks()
	if len(blocks) > 1 {
		timestamp := transaction.Timestamp()
		nextBlockTimestamp := blocks[len(blocks)-1].Timestamp + 2*pool.validationTimer.Nanoseconds()
		if nextBlockTimestamp < timestamp {
			return fmt.Errorf("the transaction timestamp is too far in the future: %v, now: %v", time.Unix(0, timestamp), time.Unix(0, nextBlockTimestamp))
		}
		currentBlockTimestamp := blocks[len(blocks)-1].Timestamp
		if timestamp < currentBlockTimestamp {
			return fmt.Errorf("the transaction timestamp is too old: %v, current block timestamp: %v", time.Unix(0, timestamp), time.Unix(0, currentBlockTimestamp))
		}
		for _, validatedTransaction := range blocks[len(blocks)-1].Transactions {
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
	if err = transaction.VerifySignature(); err != nil {
		return errors.New("failed to verify transaction")
	}
	var senderWalletAmount uint64
	if len(blocks) > 0 {
		senderWalletAmount = currentBlockchain.CalculateTotalAmount(blocks[len(blocks)-1].Timestamp, transaction.SenderAddress())
	}
	insufficientBalance := senderWalletAmount < transaction.Value()+transaction.Fee()
	if insufficientBalance {
		return errors.New("not enough balance in the sender wallet")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.transactions = append(pool.transactions, transaction)
	pool.transactionResponses = append(pool.transactionResponses, transaction.GetResponse())
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
