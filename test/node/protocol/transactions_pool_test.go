package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
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
	blockchain.AddBlock(NewGenesisBlockResponse(validatorWalletAddress))
	blockchain.AddBlock(NewEmptyBlockResponse(now - 1))
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
	blockchain.AddBlock(NewGenesisBlockResponse(validatorWalletAddress))
	blockchain.AddBlock(NewEmptyBlockResponse(now - 2))
	blockchain.AddBlock(NewEmptyBlockResponse(now - 1))
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
	blockchain.AddBlock(NewGenesisBlockResponse(validatorWalletAddress))
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	transaction, _ := protocol.NewTransactionFromRequest(&invalidTransactionRequest)
	blockchain.AddBlock(NewBlockResponse(now-1, [32]byte{}, transaction))

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
	blockchain.AddBlock(NewGenesisBlockResponse(validatorWalletAddress))

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
	blockchain.AddBlock(NewGenesisBlockResponse(validatorWalletAddress))

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
