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
)

type TransactionsPool struct {
	transactions []*protocol.Transaction
	mutex        sync.RWMutex

	blocksManager    ledger.BlocksManager
	settings         ledger.SettingsProvider
	sendersManager   network.SendersManager
	utxosManager     ledger.UtxosManager
	validatorAddress string

	logger log.Logger
}

func NewTransactionsPool(blocksManager ledger.BlocksManager, settings ledger.SettingsProvider, sendersManager network.SendersManager, utxosManager ledger.UtxosManager, validatorAddress string, logger log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blocksManager = blocksManager
	pool.settings = settings
	pool.sendersManager = sendersManager
	pool.utxosManager = utxosManager
	pool.validatorAddress = validatorAddress
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transaction *protocol.Transaction, broadcasterTarget string, hostTarget string) {
	err := pool.addTransaction(transaction)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
		return
	}
	pool.sendersManager.Incentive(broadcasterTarget)
	newTransactionRequest := protocol.NewTransactionRequest(transaction, hostTarget)
	marshaledTransactionRequest, err := json.Marshal(newTransactionRequest)
	if err != nil {
		pool.logger.Debug(fmt.Errorf("failed to marshal transaction request: %w", err).Error())
		return
	}
	senders := pool.sendersManager.Senders()
	for _, sender := range senders {
		go func(sender network.Sender) {
			_ = sender.AddTransaction(marshaledTransactionRequest)
		}(sender)
	}
}

func (pool *TransactionsPool) Transactions() []*protocol.Transaction {
	return pool.transactions
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
		fee, err := utxosManagerCopy.CalculateFee(transaction, timestamp)
		if err != nil {
			pool.logger.Warn(fmt.Errorf("transaction removed from the transactions pool, failed to calculate fee, transaction: %v\n %w", transaction, err).Error())
			rejectedTransactions = append(rejectedTransactions, transaction)
			continue
		}
		if err = utxosManagerCopy.UpdateUtxos([]*protocol.Transaction{transaction}, nextBlockTimestamp); err != nil {
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
	err = pool.blocksManager.AddBlock(timestamp, transactions, newAddresses)
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
	if err := transaction.VerifySignatures(); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}
	utxoManagerCopy := pool.utxosManager.Copy()
	lastBlockTransactions := pool.blocksManager.LastBlockTransactions()
	if err := utxoManagerCopy.UpdateUtxos(lastBlockTransactions, nextBlockTimestamp); err != nil {
		return fmt.Errorf("failed to update UTXOs: %w", err)
	}
	if err := utxoManagerCopy.UpdateUtxos(pool.transactions, nextBlockTimestamp); err != nil {
		return fmt.Errorf("failed to update UTXOs: %w", err)
	}
	_, err := utxoManagerCopy.CalculateFee(transaction, nextBlockTimestamp)
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
