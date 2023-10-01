package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/src/node/config"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"math/rand"
	"sync"
	"time"
)

type TransactionsPool struct {
	transactions []*verification.Transaction
	mutex        sync.RWMutex

	blockchain       protocol.Blockchain
	settings         config.Settings
	synchronizer     network.Synchronizer
	validatorAddress string

	validationTimestamp int64
	watch               clock.Watch

	logger log.Logger
}

func NewTransactionsPool(blockchain protocol.Blockchain, settings config.Settings, synchronizer network.Synchronizer, validatorAddress string, validationTimer time.Duration, logger log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blockchain = blockchain
	pool.settings = settings
	pool.synchronizer = synchronizer
	pool.validatorAddress = validatorAddress
	pool.validationTimestamp = validationTimer.Nanoseconds()
	pool.watch = tick.NewWatch()
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(requestBytes []byte, hostTarget string) {
	var request TransactionRequest
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		pool.logger.Debug(fmt.Errorf("failed to unmarshal transaction: %w", err).Error())
		return
	}
	err := pool.addTransaction(request.Transaction())
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
		return
	}
	pool.synchronizer.Incentive(request.TransactionBroadcasterTarget())
	marshaledTransaction, err := json.Marshal(NewTransactionRequest(request.Transaction(), hostTarget))
	neighbors := pool.synchronizer.Neighbors()
	for _, neighbor := range neighbors {
		go func(neighbor network.Neighbor) {
			_ = neighbor.AddTransaction(marshaledTransaction)
		}(neighbor)
	}
}

func (pool *TransactionsPool) Transactions() []byte {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()
	transactionsBytes, err := json.Marshal(pool.transactions)
	if err != nil {
		pool.logger.Error(err.Error())
		return nil
	}
	return transactionsBytes
}

func (pool *TransactionsPool) Validate(timestamp int64) {
	blockchainCopy := pool.blockchain.Copy()
	lastBlockTimestamp := blockchainCopy.LastBlockTimestamp()
	nextBlockTimestamp := lastBlockTimestamp + pool.validationTimestamp
	var reward uint64
	var newAddresses []string
	var hasIncome bool
	if lastBlockTimestamp == 0 {
		reward = pool.settings.GenesisAmountInParticles
		newAddresses = []string{pool.validatorAddress}
		hasIncome = true
	} else if lastBlockTimestamp == timestamp {
		pool.logger.Error("unable to create block, a block with the same timestamp is already in the blockchain")
		return
	} else if timestamp > nextBlockTimestamp {
		pool.logger.Error("unable to create block, a block is missing in the blockchain")
		return
	}
	err := blockchainCopy.AddBlock(timestamp, nil, nil)
	if err != nil {
		pool.logger.Error("failed to add temporary block")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	isAlreadySpentByOutputIdByTransactionIndex := make(map[string]map[uint16]bool)
	transactions := pool.transactions
	rand.Seed(pool.watch.Now().UnixNano())
	rand.Shuffle(len(transactions), func(i, j int) {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	})
	firstBlockTimestamp := blockchainCopy.FirstBlockTimestamp()
	var rejectedTransactions []*verification.Transaction
	for _, transaction := range transactions {
		if timestamp < transaction.Timestamp() {
			pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction timestamp is too far in the future, transaction: %v", transaction))
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		if transaction.Timestamp() < lastBlockTimestamp {
			pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction timestamp is too old, transaction: %v", transaction))
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		if err = transaction.VerifySignatures(); err != nil {
			pool.logger.Warn(fmt.Errorf("failed to verify transaction: %w", err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		fee, err := transaction.Fee(firstBlockTimestamp, pool.settings, timestamp, pool.validationTimestamp, blockchainCopy.Utxo)
		if err != nil {
			pool.logger.Warn(fmt.Errorf("transaction removed from the transactions pool, transaction: %v\n %w", transaction, err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		var isAlreadySpent bool
		for _, input := range transaction.Inputs() {
			if isAlreadySpentByOutputIndex, ok := isAlreadySpentByOutputIdByTransactionIndex[input.TransactionId()]; ok {
				if _, isAlreadySpent = isAlreadySpentByOutputIndex[input.OutputIndex()]; isAlreadySpent {
					pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, an input has already been spent, transaction: %v", transaction))
					rejectedTransactions = append(rejectedTransactions, transaction)
					break
				}
			} else {
				isAlreadySpentByOutputIdByTransactionIndex[input.TransactionId()] = make(map[uint16]bool)
			}
			isAlreadySpentByOutputIdByTransactionIndex[input.TransactionId()][input.OutputIndex()] = true
		}
		if !isAlreadySpent {
			reward += fee
		}
	}
	for _, transaction := range rejectedTransactions {
		transactions = removeTransaction(transactions, transaction)
	}
	for _, transaction := range transactions {
		for _, output := range transaction.Outputs() {
			if output.HasIncome() {
				newAddresses = append(newAddresses, output.Address())
			}
		}
	}
	rewardTransaction, err := verification.NewRewardTransaction(pool.validatorAddress, hasIncome, timestamp, reward)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block, failed to create reward transaction: %w", err).Error())
		return
	}
	transactions = append(transactions, rewardTransaction)
	secondBlockchainCopy := pool.blockchain.Copy()
	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		pool.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
	}
	err = secondBlockchainCopy.AddBlock(timestamp, transactionsBytes, nil)
	if err != nil {
		pool.logger.Error(fmt.Errorf("next block creation would fail: %w", err).Error())
	}
	err = secondBlockchainCopy.AddBlock(nextBlockTimestamp, nil, nil)
	if err != nil {
		pool.logger.Error(fmt.Errorf("later block creation would fail: %w", err).Error())
		return
	}
	err = pool.blockchain.AddBlock(timestamp, transactionsBytes, newAddresses)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block: %w", err).Error())
		return
	}
	pool.clear()
	if lastBlockTimestamp == 0 {
		pool.logger.Info("first block validation done, the node is now fully operational")
	}
	pool.logger.Debug(fmt.Sprintf("reward: %d", reward))
}

func (pool *TransactionsPool) addTransaction(transaction *verification.Transaction) error {
	blockchainCopy := pool.blockchain.Copy()
	lastBlockTimestamp := blockchainCopy.LastBlockTimestamp()
	if lastBlockTimestamp == 0 {
		return errors.New("the blockchain is empty")
	}
	nextBlockTimestamp := lastBlockTimestamp + pool.validationTimestamp
	err := blockchainCopy.AddBlock(nextBlockTimestamp, nil, nil)
	if err != nil {
		return errors.New("failed to add temporary block")
	}
	timestamp := transaction.Timestamp()
	if nextBlockTimestamp < timestamp {
		return fmt.Errorf("the transaction timestamp is too far in the future: %v, now: %v", time.Unix(0, timestamp), time.Unix(0, nextBlockTimestamp))
	}
	currentBlockTimestamp := lastBlockTimestamp
	if timestamp < currentBlockTimestamp {
		return fmt.Errorf("the transaction timestamp is too old: %v, current block timestamp: %v", time.Unix(0, timestamp), time.Unix(0, currentBlockTimestamp))
	}
	for _, pendingTransaction := range pool.transactions {
		if transaction.Equals(pendingTransaction) {
			return errors.New("the transaction is already in the transactions pool")
		}
	}
	if err = transaction.VerifySignatures(); err != nil {
		return fmt.Errorf("failed to verify transaction: %w", err)
	}
	firstBlockTimestamp := blockchainCopy.FirstBlockTimestamp()
	_, err = transaction.Fee(firstBlockTimestamp, pool.settings, nextBlockTimestamp, pool.validationTimestamp, blockchainCopy.Utxo)
	if err != nil {
		return fmt.Errorf("failed to verify fee: %w", err)
	}
	transactions := []*verification.Transaction{transaction}
	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}
	err = blockchainCopy.AddBlock(nextBlockTimestamp+pool.validationTimestamp, transactionsBytes, nil)
	if err != nil {
		return fmt.Errorf("next block creation would fail: %w", err)
	}
	err = blockchainCopy.AddBlock(nextBlockTimestamp+pool.validationTimestamp, nil, nil)
	if err != nil {
		return fmt.Errorf("later block creation would fail: %w", err)
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.transactions = append(pool.transactions, transaction)
	return nil
}

func (pool *TransactionsPool) clear() {
	pool.transactions = nil
}

func removeTransaction(transactions []*verification.Transaction, removedTransaction *verification.Transaction) []*verification.Transaction {
	for i := 0; i < len(transactions); i++ {
		if transactions[i] == removedTransaction {
			transactions = append(transactions[:i], transactions[i+1:]...)
			return transactions
		}
	}
	return transactions
}
