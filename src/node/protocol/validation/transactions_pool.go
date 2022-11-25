package validation

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
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

	blockchain       protocol.Blockchain
	registry         protocol.Registry
	validatorAddress string
	genesisAmount    uint64

	validationTimer time.Duration
	time            clock.Time

	waitGroup *sync.WaitGroup
	logger    *log.Logger
}

func NewTransactionsPool(blockchain protocol.Blockchain, registry protocol.Registry, validatorAddress string, genesisAmount uint64, validationTimer time.Duration, time clock.Time, logger *log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.blockchain = blockchain
	pool.registry = registry
	pool.validatorAddress = validatorAddress
	pool.genesisAmount = genesisAmount
	pool.validationTimer = validationTimer
	pool.time = time
	var waitGroup sync.WaitGroup
	pool.waitGroup = &waitGroup
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transactionRequest *network.TransactionRequest, neighbors []network.Neighbor) {
	pool.waitGroup.Add(1)
	go func() {
		defer pool.waitGroup.Done()
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
	}()
}

func (pool *TransactionsPool) addTransaction(transactionRequest *network.TransactionRequest) (err error) {
	transaction, err := NewTransactionFromRequest(transactionRequest)
	if err != nil {
		err = fmt.Errorf("failed to instantiate transaction: %w", err)
		return
	}
	currentBlockchain := pool.blockchain.Copy()
	blocks := currentBlockchain.Blocks()
	if len(blocks) > 1 {
		timestamp := transaction.Timestamp()
		nextBlockTimestamp := blocks[len(blocks)-1].Timestamp + 2*pool.validationTimer.Nanoseconds()
		if nextBlockTimestamp < timestamp {
			err = fmt.Errorf("the transaction timestamp is too far in the future: %v, now: %v", time.Unix(0, timestamp), time.Unix(0, nextBlockTimestamp))
			return
		}
		previousBlockTimestamp := blocks[len(blocks)-2].Timestamp
		if timestamp < previousBlockTimestamp {
			err = fmt.Errorf("the transaction timestamp is too old: %v, previous block timestamp: %v", time.Unix(0, timestamp), time.Unix(0, previousBlockTimestamp))
			return
		}
		for i := len(blocks) - 2; i < len(blocks); i++ {
			for _, validatedTransaction := range blocks[i].Transactions {
				if transaction.Equals(validatedTransaction) {
					err = errors.New("the transaction is already in the blockchain")
					return
				}
			}
		}
	}
	for _, pendingTransaction := range pool.transactionResponses {
		if transaction.Equals(pendingTransaction) {
			err = errors.New("the transaction is already in the transactions pool")
			return
		}
	}
	if err = transaction.VerifySignature(); err != nil {
		err = errors.New("failed to verify transaction")
		return
	}
	var senderWalletAmount uint64
	if len(blocks) > 0 {
		senderWalletAmount = currentBlockchain.CalculateTotalAmount(blocks[len(blocks)-1].Timestamp, transaction.SenderAddress())
	}
	insufficientBalance := senderWalletAmount < transaction.Value()+transaction.Fee()
	if insufficientBalance {
		err = errors.New("not enough balance in the sender wallet")
		return
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.transactions = append(pool.transactions, transaction)
	pool.transactionResponses = append(pool.transactionResponses, transaction.GetResponse())
	return
}

func (pool *TransactionsPool) Transactions() []*network.TransactionResponse {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()
	return pool.transactionResponses
}

func (pool *TransactionsPool) Validate(timestamp int64) {
	currentBlockchain := pool.blockchain.Copy()
	if currentBlockchain.IsEmpty() {
		genesisTransaction := NewRewardTransaction(pool.validatorAddress, timestamp, pool.genesisAmount)
		transactions := []*network.TransactionResponse{genesisTransaction}
		pool.blockchain.AddBlock(timestamp, transactions, nil)
		pool.logger.Debug("genesis block added")
		return
	}

	blocks := currentBlockchain.Blocks()
	lastBlock := blocks[len(blocks)-1]

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
		if len(blocks) > 1 {
			if transaction.Timestamp < blocks[len(blocks)-2].Timestamp {
				pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction timestamp is too old, transaction: %v", transaction))
				rejectedTransactions = append(rejectedTransactions, transaction)
				continue
			}
			for j := len(blocks) - 2; j < len(blocks); j++ {
				for _, validatedTransaction := range blocks[j].Transactions {
					if pool.transactions[i].Equals(validatedTransaction) {
						pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, the transaction is already in the blockchain, transaction: %v", transaction))
						rejectedTransactions = append(rejectedTransactions, transaction)
						continue
					}
				}
			}
		}
	}
	for _, transaction := range rejectedTransactions {
		transactionResponses = removeTransactions(transactionResponses, transaction)
	}
	for _, transaction := range transactionResponses {
		fee := transaction.Fee
		totalTransactionsValueBySenderAddress[transaction.SenderAddress] += transaction.Value + fee
		reward += fee
	}
	registeredAddresses := lastBlock.RegisteredAddresses
	registeredAddressesMap := make(map[string]bool)
	for _, registeredAddress := range registeredAddresses {
		registeredAddressesMap[registeredAddress] = true
	}
	for senderAddress, totalTransactionsValue := range totalTransactionsValueBySenderAddress {
		senderTotalAmount := currentBlockchain.CalculateTotalAmount(timestamp, senderAddress)
		if totalTransactionsValue > senderTotalAmount {
			rejectedTransactions = nil
			rand.Seed(pool.time.Now().UnixNano())
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
				transactionResponses = removeTransactions(transactionResponses, transaction)
				pool.logger.Warn(fmt.Sprintf("transaction removed from the transactions pool, total transactions value exceeds its sender wallet amount, transaction: %v", transaction))
			}
		}
		if totalTransactionsValue > 0 {
			if _, isRegistered := registeredAddressesMap[senderAddress]; !isRegistered {
				registeredAddressesMap[senderAddress] = true
			}
		}
	}
	var newRegisteredAddresses []string
	for registeredAddress := range registeredAddressesMap {
		var isPohValid bool
		isPohValid, err := pool.registry.IsRegistered(registeredAddress)
		if err != nil {
			pool.logger.Error(fmt.Errorf("failed to get proof of humanity: %w", err).Error())
		} else if isPohValid {
			newRegisteredAddresses = append(newRegisteredAddresses, registeredAddress)
		}
	}
	pool.clear()

	if lastBlock.Timestamp == timestamp {
		pool.logger.Error("unable to create block, a block with the same timestamp is already in the blockchain")
		return
	}
	rewardTransaction := NewRewardTransaction(pool.validatorAddress, timestamp, reward)
	transactionResponses = append(transactionResponses, rewardTransaction)
	pool.blockchain.AddBlock(timestamp, transactionResponses, newRegisteredAddresses)
	pool.logger.Debug(fmt.Sprintf("reward: %d", reward))
}

func (pool *TransactionsPool) clear() {
	pool.transactions = nil
	pool.transactionResponses = nil
}

func (pool *TransactionsPool) Wait() {
	pool.waitGroup.Wait()
}

func removeTransactions(transactions []*network.TransactionResponse, removedTransaction *network.TransactionResponse) []*network.TransactionResponse {
	for i := 0; i < len(transactions); i++ {
		if transactions[i] == removedTransaction {
			transactions = append(transactions[:i], transactions[i+1:]...)
			return transactions
		}
	}
	return transactions
}
