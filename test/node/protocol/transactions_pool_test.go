package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node"
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
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, NewEmptyBlockResponse(now-1))
	blockchainMock := new(node.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchainMock, nil)
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
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, NewEmptyBlockResponse(now-2))
	blockResponses = append(blockResponses, NewEmptyBlockResponse(now-1))
	blockchainMock := new(node.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchainMock, nil)
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
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	transaction, _ := protocol.NewTransactionFromRequest(&invalidTransactionRequest)
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, NewBlockResponse(now-1, [32]byte{}, transaction))
	blockchainMock := new(node.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, blockchainMock, nil)
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
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(node.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)

	var amount uint64 = 1
	transaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(walletA.PrivateKey())
	transactionRequest := transaction.GetRequest()

	// Act
	pool.AddTransaction(&transactionRequest, blockchainMock, nil)
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
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(node.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	pool := protocol.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)

	var amount uint64 = 1
	transaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(validatorWallet.PrivateKey())
	transactionRequest := transaction.GetRequest()
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return *transactionRequest.Value + *transactionRequest.Fee }

	// Act
	pool.AddTransaction(&transactionRequest, blockchainMock, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}
