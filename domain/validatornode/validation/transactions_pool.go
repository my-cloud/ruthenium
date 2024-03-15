package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/domain"
	"github.com/my-cloud/ruthenium/domain/clock"
	"github.com/my-cloud/ruthenium/domain/ledger"
	"github.com/my-cloud/ruthenium/domain/network"
	"github.com/my-cloud/ruthenium/domain/validatornode"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"math/rand"
	"sync"
	"time"
)

type TransactionsPool struct {
	transactions []*ledger.Transaction
	mutex        sync.RWMutex

	blocksManager    domain.BlocksManager
	settings         validatornode.SettingsProvider
	neighborsManager network.NeighborsManager
	validatorAddress string

	watch  domain.TimeProvider
	logger log.Logger
}

func NewTransactionsPool(blocksManager domain.BlocksManager, settings validatornode.SettingsProvider, neighborsManager network.NeighborsManager, validatorAddress string, logger log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blocksManager = blocksManager
	pool.settings = settings
	pool.neighborsManager = neighborsManager
	pool.validatorAddress = validatorAddress
	pool.watch = clock.NewWatch()
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transactionRequestBytes []byte, hostTarget string) {
	var transactionRequest ledger.TransactionRequest
	if err := json.Unmarshal(transactionRequestBytes, &transactionRequest); err != nil {
		pool.logger.Debug(fmt.Errorf("failed to unmarshal transaction: %w", err).Error())
		return
	}
	transaction := transactionRequest.Transaction()
	err := pool.addTransaction(transaction)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
		return
	}
	pool.neighborsManager.Incentive(transactionRequest.TransactionBroadcasterTarget())
	newTransactionRequest := ledger.NewTransactionRequest(transaction, hostTarget)
	marshaledTransactionRequest, err := json.Marshal(newTransactionRequest)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to marshal transaction request: %w", err).Error())
		return
	}
	neighbors := pool.neighborsManager.Neighbors()
	for _, neighbor := range neighbors {
		go func(neighbor network.Neighbor) {
			_ = neighbor.AddTransaction(marshaledTransactionRequest)
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
	blocksManagerCopy := pool.blocksManager.Copy()
	lastBlockTimestamp := blocksManagerCopy.LastBlockTimestamp()
	nextBlockTimestamp := lastBlockTimestamp + pool.settings.ValidationTimestamp()
	var reward uint64
	var newAddresses []string
	var isYielding bool
	if lastBlockTimestamp == 0 {
		reward = pool.settings.GenesisAmount()
		newAddresses = []string{pool.validatorAddress}
		isYielding = true
	} else if lastBlockTimestamp == timestamp {
		pool.logger.Error("unable to create block, a block with the same timestamp is already in the blockchain")
		return
	} else if timestamp > nextBlockTimestamp {
		pool.logger.Error("unable to create block, a block is missing in the blockchain")
		return
	}
	err := blocksManagerCopy.AddBlock(timestamp, nil, nil)
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
	var rejectedTransactions []*ledger.Transaction
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
		if err = transaction.VerifySignatures(blocksManagerCopy.Utxo); err != nil {
			pool.logger.Warn(fmt.Errorf("failed to verify transaction: %w", err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		fee, err := transaction.Fee(pool.settings, timestamp, blocksManagerCopy.Utxo)
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
			if output.IsYielding() {
				newAddresses = append(newAddresses, output.Address())
			}
		}
	}
	rewardTransaction, err := ledger.NewRewardTransaction(pool.validatorAddress, isYielding, timestamp, reward)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block, failed to create reward transaction: %w", err).Error())
		return
	}
	transactions = append(transactions, rewardTransaction)
	secondBlockchainCopy := pool.blocksManager.Copy()
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
	err = pool.blocksManager.AddBlock(timestamp, transactionsBytes, newAddresses)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block: %w", err).Error())
		return
	}
	pool.clear()
	pool.logger.Debug(fmt.Sprintf("reward: %d", reward))
}

func (pool *TransactionsPool) addTransaction(transaction *ledger.Transaction) error {
	blocksManagerCopy := pool.blocksManager.Copy()
	lastBlockTimestamp := blocksManagerCopy.LastBlockTimestamp()
	if lastBlockTimestamp == 0 {
		return errors.New("the blockchain is empty")
	}
	nextBlockTimestamp := lastBlockTimestamp + pool.settings.ValidationTimestamp()
	err := blocksManagerCopy.AddBlock(nextBlockTimestamp, nil, nil)
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
	if err = transaction.VerifySignatures(blocksManagerCopy.Utxo); err != nil {
		return fmt.Errorf("failed to verify transaction: %w", err)
	}
	_, err = transaction.Fee(pool.settings, nextBlockTimestamp, blocksManagerCopy.Utxo)
	if err != nil {
		return fmt.Errorf("failed to verify fee: %w", err)
	}
	transactions := []*ledger.Transaction{transaction}
	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}
	err = blocksManagerCopy.AddBlock(nextBlockTimestamp+pool.settings.ValidationTimestamp(), transactionsBytes, nil)
	if err != nil {
		return fmt.Errorf("next block creation would fail: %w", err)
	}
	err = blocksManagerCopy.AddBlock(nextBlockTimestamp+2*pool.settings.ValidationTimestamp(), nil, nil)
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

func removeTransaction(transactions []*ledger.Transaction, removedTransaction *ledger.Transaction) []*ledger.Transaction {
	for i := 0; i < len(transactions); i++ {
		if transactions[i] == removedTransaction {
			transactions = append(transactions[:i], transactions[i+1:]...)
			return transactions
		}
	}
	return transactions
}
