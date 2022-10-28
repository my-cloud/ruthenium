package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"testing"
	"time"
)

func Test_AddTransaction_Allowed_TransactionAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := encryption.NewWallet()
	minerWalletAddress := minerWallet.Address()
	walletA, _ := encryption.NewWallet()
	walletAAddress := walletA.Address()
	walletB, _ := encryption.NewWallet()
	walletBAddress := walletB.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := NewTimeMock()
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, timeMock, logger)
	validation := protocol.NewValidation(minerWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	validation.Validate()
	validation.Wait()

	var amount1 uint64 = 40 * network.ParticlesCount
	transaction1 := server.NewTransaction(walletAAddress, minerWalletAddress, minerWallet.PublicKey(), 1, amount1)
	_ = transaction1.Sign(minerWallet.PrivateKey())

	var amount2 uint64 = 10 * network.ParticlesCount
	transaction2 := server.NewTransaction(walletBAddress, walletAAddress, walletA.PublicKey(), 2, amount2)
	_ = transaction2.Sign(walletA.PrivateKey())

	// Act
	addTransaction(pool, blockchain, validation, transaction1)
	actualWalletAAmount := blockchain.CalculateTotalAmount(0, walletAAddress)

	addTransaction(pool, blockchain, validation, transaction2)
	actualWalletBAmount := blockchain.CalculateTotalAmount(0, walletBAddress)

	// Assert
	test.Assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	test.Assert(t, amount2 == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", amount2, actualWalletBAmount))
}

func Test_AddTransaction_NotAllowed_TransactionNotAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := encryption.NewWallet()
	minerWalletAddress := minerWallet.Address()
	walletA, _ := encryption.NewWallet()
	walletAAddress := walletA.Address()
	walletB, _ := encryption.NewWallet()
	walletBAddress := walletB.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := NewTimeMock()
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, timeMock, logger)
	validation := protocol.NewValidation(minerWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	validation.Validate()
	validation.Wait()

	var amount1 uint64 = 40 * network.ParticlesCount
	transaction1 := server.NewTransaction(walletAAddress, minerWalletAddress, minerWallet.PublicKey(), 0, amount1)
	_ = transaction1.Sign(minerWallet.PrivateKey())

	var amount2 uint64 = 10 * network.ParticlesCount
	transaction2 := server.NewTransaction(walletBAddress, walletAAddress, walletA.PublicKey(), 0, amount2)
	_ = transaction2.Sign(walletB.PrivateKey())

	// Act
	addTransaction(pool, blockchain, validation, transaction1)
	actualWalletAAmount := blockchain.CalculateTotalAmount(0, walletAAddress)

	addTransaction(pool, blockchain, validation, transaction2)
	actualWalletBAmount := blockchain.CalculateTotalAmount(0, walletBAddress)

	// Assert
	test.Assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	var expectedWalletBAmount uint64 = 0
	test.Assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func addTransaction(pool *protocol.TransactionsPool, blockchain *protocol.Blockchain, validation *protocol.Validation, transaction *server.Transaction) {
	transactionRequest := transaction.GetRequest()
	pool.AddTransaction(&transactionRequest, blockchain, nil)
	pool.Wait()
	validation.Validate()
	validation.Wait()
}

func Test_AddTransaction_TransactionTimestampIsOlderThan2Blocks_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.NewWallet()
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := NewTimeMock()
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, timeMock, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	addEmptyBlock(blockchain, 1)
	addEmptyBlock(blockchain, 2)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), 0, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions length. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_TransactionTimestampIsOlderThan2Blocks_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.NewWallet()
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := NewTimeMock()
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, timeMock, logger)
	validation := protocol.NewValidation(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), 0, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()
	addEmptyBlock(blockchain, 1)
	addEmptyBlock(blockchain, 2)

	// Act
	validation.Validate()
	validation.Wait()

	// Assert
	blocks := blockchain.Blocks()
	lastBlock := blocks[len(blocks)-1]
	expectedTransactionsLength := 1
	actualTransactionsLength := len(lastBlock.Transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions length. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	actualTransaction, _ := protocol.NewTransactionFromResponse(lastBlock.Transactions[0])
	test.Assert(t, actualTransaction.IsReward(), "The single expected transaction is not a reward whereas it should be.")
}

func addEmptyBlock(blockchain *protocol.Blockchain, timestamp int64) {
	blockchain.AddBlock(&neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        nil,
		RegisteredAddresses: nil,
	})
}

func addGenesisBlock(validatorWalletAddress string, blockchain *protocol.Blockchain) {
	genesisTransaction := protocol.NewRewardTransaction(validatorWalletAddress, 0, network.ParticlesCount)
	blockchain.AddBlock(&neighborhood.BlockResponse{
		Timestamp:           0,
		PreviousHash:        [32]byte{},
		Transactions:        []*neighborhood.TransactionResponse{genesisTransaction.GetResponse()},
		RegisteredAddresses: nil,
	})
}
