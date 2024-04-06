package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/presentation"
)

type TransactionsPool struct {
	transactions []*protocol.Transaction
	mutex        sync.RWMutex

	blocksManager    ledger.BlocksManager
	settings         ledger.SettingsProvider
	neighborsManager network.NeighborsManager
	utxosManager     ledger.UtxosManager
	validatorAddress string

	logger log.Logger
}

func NewTransactionsPool(blocksManager ledger.BlocksManager, settings ledger.SettingsProvider, neighborsManager network.NeighborsManager, utxosManager ledger.UtxosManager, validatorAddress string, logger log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blocksManager = blocksManager
	pool.settings = settings
	pool.neighborsManager = neighborsManager
	pool.utxosManager = utxosManager
	pool.validatorAddress = validatorAddress
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transactionRequestBytes []byte, hostTarget string) {
	var transactionRequest protocol.TransactionRequest
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
	newTransactionRequest := protocol.NewTransactionRequest(transaction, hostTarget)
	marshaledTransactionRequest, err := json.Marshal(newTransactionRequest)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to marshal transaction request: %w", err).Error())
		return
	}
	neighbors := pool.neighborsManager.Neighbors()
	for _, neighbor := range neighbors {
		go func(neighbor presentation.NeighborCaller) {
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
	lastBlockTimestamp := pool.blocksManager.LastBlockTimestamp()
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
	lastBlockTransactions := pool.blocksManager.LastBlockTransactions()
	utxosManagerCopy := pool.utxosManager.Copy()
	if err := utxosManagerCopy.UpdateUtxos(lastBlockTransactions, nextBlockTimestamp); err != nil {
		pool.logger.Error(fmt.Errorf("failed to update UTXOs: %w", err).Error())
		return
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	transactions := pool.transactions
	rand.Seed(timestamp)
	rand.Shuffle(len(transactions), func(i, j int) {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	})
	var rejectedTransactions []*protocol.Transaction
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
		if err := transaction.VerifySignatures(); err != nil {
			pool.logger.Warn(fmt.Errorf("transaction removed from the transactions pool, failed to verify signature, transaction: %v\n %w", transaction, err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		// TODO choose transactions bytes or interfaces
		var inputs []ledger.InputInfoProvider
		for _, input := range transaction.Inputs() {
			inputs = append(inputs, input)
		}
		var outputs []ledger.OutputInfoProvider
		for _, output := range transaction.Outputs() {
			outputs = append(outputs, output)
		}
		fee, err := utxosManagerCopy.CalculateFee(inputs, outputs, timestamp)
		if err != nil {
			pool.logger.Warn(fmt.Errorf("transaction removed from the transactions pool, failed to calculate fee, transaction: %v\n %w", transaction, err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		transactionBytes, err := json.Marshal([]*protocol.Transaction{transaction})
		if err != nil {
			pool.logger.Warn(fmt.Errorf("transaction removed from the transactions pool, failed to marshal transaction, transaction: %v\n %w", transaction, err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		if err = utxosManagerCopy.UpdateUtxos(transactionBytes, nextBlockTimestamp); err != nil {
			pool.logger.Warn(fmt.Errorf("transaction removed from the transactions pool, failed to update UTXOs, transaction: %v\n %w", transaction, err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		reward += fee
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
	rewardTransaction, err := protocol.NewRewardTransaction(pool.validatorAddress, isYielding, timestamp, reward)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block, failed to create reward transaction: %w", err).Error())
		return
	}
	transactions = append(transactions, rewardTransaction)
	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		pool.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
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

func (pool *TransactionsPool) addTransaction(transaction *protocol.Transaction) error {
	lastBlockTimestamp := pool.blocksManager.LastBlockTimestamp()
	if lastBlockTimestamp == 0 {
		return errors.New("the blockchain is empty")
	}
	nextBlockTimestamp := lastBlockTimestamp + pool.settings.ValidationTimestamp()
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
	transactionsBytes, err := json.Marshal(pool.transactions)
	if err != nil {
		return fmt.Errorf("failed to marshal transactions: %w", err)
	}
	if err = transaction.VerifySignatures(); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}
	utxoManagerCopy := pool.utxosManager.Copy()
	lastBlockTransactions := pool.blocksManager.LastBlockTransactions()
	if err = utxoManagerCopy.UpdateUtxos(lastBlockTransactions, nextBlockTimestamp); err != nil {
		return fmt.Errorf("failed to update UTXOs: %w", err)
	}
	if err = utxoManagerCopy.UpdateUtxos(transactionsBytes, nextBlockTimestamp); err != nil {
		return fmt.Errorf("failed to update UTXOs: %w", err)
	}
	// TODO choose transactions bytes or interfaces
	var inputs []ledger.InputInfoProvider
	for _, input := range transaction.Inputs() {
		inputs = append(inputs, input)
	}
	var outputs []ledger.OutputInfoProvider
	for _, output := range transaction.Outputs() {
		outputs = append(outputs, output)
	}
	_, err = utxoManagerCopy.CalculateFee(inputs, outputs, nextBlockTimestamp)
	if err != nil {
		return fmt.Errorf("failed to verify fee: %w", err)
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.transactions = append(pool.transactions, transaction)
	return nil
}

func (pool *TransactionsPool) clear() {
	pool.transactions = nil
}

func removeTransaction(transactions []*protocol.Transaction, removedTransaction *protocol.Transaction) []*protocol.Transaction {
	for i := 0; i < len(transactions); i++ {
		if transactions[i] == removedTransaction {
			transactions = append(transactions[:i], transactions[i+1:]...)
			return transactions
		}
	}
	return transactions
}
