package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"testing"
	"time"
)

func Test_AddTransaction_TransactionTimestampIsInTheFuture_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 2
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	addEmptyBlock(blockchain, now-1)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_TransactionTimestampIsOlderThan2Blocks_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 3
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	addEmptyBlock(blockchain, now-2)
	addEmptyBlock(blockchain, now-1)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}
func Test_AddTransaction_TransactionIsAlreadyInTheBlockchain_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 2
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	transaction, _ := protocol.NewTransactionFromRequest(&invalidTransactionRequest)
	addBlock(blockchain, now-1, transaction)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_InvalidSignature_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	walletA, _ := encryption.DecodeWallet(test.Mnemonic2, test.DerivationPath, "", "")
	walletAAddress := walletA.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)

	var amount uint64 = 1
	transaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(walletA.PrivateKey())
	transactionRequest := transaction.GetRequest()

	// Act
	pool.AddTransaction(&transactionRequest, blockchain, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_ValidTransaction_TransactionAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	walletA, _ := encryption.DecodeWallet(test.Mnemonic2, test.DerivationPath, "", "")
	walletAAddress := walletA.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)

	var amount uint64 = 1
	transaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(validatorWallet.PrivateKey())
	transactionRequest := transaction.GetRequest()

	// Act
	pool.AddTransaction(&transactionRequest, blockchain, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_InvalidSignature_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	walletA, _ := encryption.DecodeWallet(test.Mnemonic2, test.DerivationPath, "", "")
	walletAAddress := walletA.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewValidation(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(walletA.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Act
	validation.Validate()
	validation.Wait()

	// Assert
	assertTransactionNotValidated(t, blockchain)
}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewValidation(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Act
	validation.Validate()
	validation.Wait()

	// Assert
	assertTransactionNotValidated(t, blockchain)
}

func Test_Validate_TransactionTimestampIsOlderThan2Blocks_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 3
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewValidation(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()
	addEmptyBlock(blockchain, now-2)
	addEmptyBlock(blockchain, now-1)

	// Act
	validation.Validate()
	validation.Wait()

	// Assert
	assertTransactionNotValidated(t, blockchain)
}

func Test_Validate_TransactionIsAlreadyInTheBlockchain_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 2
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewValidation(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()
	transaction, _ := protocol.NewTransactionFromRequest(&invalidTransactionRequest)
	addBlock(blockchain, now-1, transaction)

	// Act
	validation.Validate()
	validation.Wait()

	// Assert
	assertTransactionNotValidated(t, blockchain)
}

func Test_Validate_ValidTransaction_TransactionValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	synchronizerMock := new(SynchronizerMock)
	blockchain := protocol.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewValidation(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	addGenesisBlock(validatorWalletAddress, blockchain)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Act
	validation.Validate()
	validation.Wait()

	// Assert
	blocks := blockchain.Blocks()
	lastBlock := blocks[len(blocks)-1]
	expectedTransactionsLength := 2
	actualTransactionsLength := len(lastBlock.Transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func assertTransactionNotValidated(t *testing.T, blockchain *protocol.Blockchain) {
	blocks := blockchain.Blocks()
	lastBlock := blocks[len(blocks)-1]
	expectedTransactionsLength := 1
	actualTransactionsLength := len(lastBlock.Transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	actualTransaction, _ := protocol.NewTransactionFromResponse(lastBlock.Transactions[0])
	test.Assert(t, actualTransaction.IsReward(), "The single expected transaction is not a reward whereas it should be.")
}

func addGenesisBlock(validatorWalletAddress string, blockchain *protocol.Blockchain) {
	genesisTransaction := protocol.NewRewardTransaction(validatorWalletAddress, 0, 1e13)
	blockchain.AddBlock(&neighborhood.BlockResponse{
		Timestamp:           0,
		PreviousHash:        [32]byte{},
		Transactions:        []*neighborhood.TransactionResponse{genesisTransaction.GetResponse()},
		RegisteredAddresses: nil,
	})
}

func addEmptyBlock(blockchain *protocol.Blockchain, timestamp int64) {
	blockchain.AddBlock(&neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        nil,
		RegisteredAddresses: nil,
	})
}

func addBlock(blockchain *protocol.Blockchain, timestamp int64, transaction *protocol.Transaction) {
	blockchain.AddBlock(&neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        []*neighborhood.TransactionResponse{transaction.GetResponse()},
		RegisteredAddresses: nil,
	})
}
