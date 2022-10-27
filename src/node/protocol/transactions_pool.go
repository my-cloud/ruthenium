package protocol

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math/rand"
	"sync"
)

type TransactionsPool struct {
	transactions         []*Transaction
	transactionResponses []*neighborhood.TransactionResponse
	mutex                sync.RWMutex

	registrable Registry

	timeable clock.Time

	waitGroup *sync.WaitGroup
	logger    *log.Logger
}

func NewTransactionsPool(registrable Registry, timeable clock.Time, logger *log.Logger) *TransactionsPool {
	pool := new(TransactionsPool)
	pool.registrable = registrable
	pool.timeable = timeable
	var waitGroup sync.WaitGroup
	pool.waitGroup = &waitGroup
	pool.logger = logger
	return pool
}

func (pool *TransactionsPool) AddTransaction(transactionRequest *neighborhood.TransactionRequest, blockchain network.Blockchain, neighbors []neighborhood.Neighbor) {
	pool.waitGroup.Add(1)
	go func() {
		defer pool.waitGroup.Done()
		err := pool.addTransaction(transactionRequest, blockchain)
		if err != nil {
			pool.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
			return
		}
		for _, neighbor := range neighbors {
			go func(neighbor neighborhood.Neighbor) {
				_ = neighbor.AddTransaction(*transactionRequest)
			}(neighbor)
		}
	}()
}

func (pool *TransactionsPool) addTransaction(transactionRequest *neighborhood.TransactionRequest, blockchain network.Blockchain) (err error) {
	transaction, err := NewTransactionFromRequest(transactionRequest)
	if err != nil {
		err = fmt.Errorf("failed to instantiate transaction: %w", err)
		return
	}
	blocks := blockchain.Blocks()
	if len(blocks) > 2 {
		if transaction.Timestamp() < blocks[len(blocks)-2].Timestamp {
			err = errors.New("the transaction timestamp is invalid")
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
		senderWalletAmount = blockchain.CalculateTotalAmount(blocks[len(blocks)-1].Timestamp, transaction.SenderAddress())
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

func (pool *TransactionsPool) Transactions() []*neighborhood.TransactionResponse {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()
	return pool.transactionResponses
}

func (pool *TransactionsPool) Validate(timestamp int64, blockchain network.Blockchain, address string) {
	blocks := blockchain.Blocks()
	lastBlockResponse := blocks[len(blocks)-1]
	lastBlock, err := NewBlockFromResponse(lastBlockResponse)
	if err != nil {
		pool.logger.Error(fmt.Errorf("unable to create block: %w", err).Error())
		return
	}
	if lastBlock.Timestamp() == timestamp {
		pool.logger.Error("unable to create block, a block with the same timestamp is already in the blockchain")
		return
	}
	lastBlockHash, err := lastBlock.Hash()
	if err != nil {
		pool.logger.Error(fmt.Errorf("failed calculate last block hash: %w", err).Error())
		return
	}
	isValidatorPohValid, err := pool.registrable.IsRegistered(address)
	if err != nil {
		pool.logger.Error(fmt.Errorf("failed to get proof of humanity: %w", err).Error())
		return
	} else if !isValidatorPohValid {
		pool.logger.Error("validator proof of humanity is invalid")
		return
	}

	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	totalTransactionsValueBySenderAddress := make(map[string]uint64)
	transactions := pool.transactions
	var reward uint64
	for _, transaction := range transactions {
		fee := transaction.Fee()
		totalTransactionsValueBySenderAddress[transaction.SenderAddress()] += transaction.Value() + fee
		reward += fee
	}
	registeredAddresses := lastBlock.RegisteredAddresses()
	registeredAddressesMap := make(map[string]bool)
	for _, registeredAddress := range registeredAddresses {
		registeredAddressesMap[registeredAddress] = true
	}
	for senderAddress, totalTransactionsValue := range totalTransactionsValueBySenderAddress {
		senderTotalAmount := blockchain.CalculateTotalAmount(timestamp, senderAddress)
		if totalTransactionsValue > senderTotalAmount {
			var rejectedTransactions []*Transaction
			rand.Seed(pool.timeable.Now().UnixNano())
			rand.Shuffle(len(transactions), func(i, j int) { transactions[i], transactions[j] = transactions[j], transactions[i] })
			for _, transaction := range transactions {
				if transaction.SenderAddress() == senderAddress {
					rejectedTransactions = append(rejectedTransactions, transaction)
					fee := transaction.Fee()
					totalTransactionsValue -= transaction.Value() + fee
					reward -= fee
					if totalTransactionsValue <= senderTotalAmount {
						break
					}
				}
			}
			for _, transaction := range rejectedTransactions {
				transactions = removeTransactions(transactions, transaction)
			}
			pool.logger.Warn("transactions removed from the transactions pool, total transactions value exceeds its sender wallet amount")
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
		isPohValid, err = pool.registrable.IsRegistered(registeredAddress)
		if err != nil {
			pool.logger.Error(fmt.Errorf("failed to get proof of humanity: %w", err).Error())
		} else if isPohValid {
			newRegisteredAddresses = append(newRegisteredAddresses, registeredAddress)
		}
	}
	rewardTransaction := NewRewardTransaction(address, timestamp, reward)
	transactions = append(transactions, rewardTransaction)
	block := NewBlock(timestamp, lastBlockHash, transactions, newRegisteredAddresses)
	blockchain.AddBlock(block.GetResponse())
	pool.clear()
	pool.logger.Debug(fmt.Sprintf("reward: %d", reward))
}

func (pool *TransactionsPool) clear() {
	pool.transactions = nil
	pool.transactionResponses = nil
}

func (pool *TransactionsPool) Wait() {
	pool.waitGroup.Wait()
}
