package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"testing"
	"time"
)

func Test_AddTransaction_TransactionTimestampIsInTheFuture_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 2
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, mock.NewEmptyBlockResponse(now-1))
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, nil)
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
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 3
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, mock.NewEmptyBlockResponse(now-2))
	blockResponses = append(blockResponses, mock.NewEmptyBlockResponse(now-1))
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, nil)
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
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 2
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	transaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	var transactionResponses []*network.TransactionResponse
	transactionResponses = append(transactionResponses, transaction.GetResponse())
	blockResponses = append(blockResponses, mock.NewBlockResponse(now-1, [32]byte{}, transactionResponses, nil))
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, nil)
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
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)

	var amount uint64 = 1
	transaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(walletA.PrivateKey())
	transactionRequest := transaction.GetRequest()

	// Act
	pool.AddTransaction(&transactionRequest, nil)
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
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)

	var amount uint64 = 1
	transaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(validatorWallet.PrivateKey())
	transactionRequest := transaction.GetRequest()
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return *transactionRequest.Value + *transactionRequest.Fee }

	// Act
	pool.AddTransaction(&transactionRequest, nil)
	pool.Wait()

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_BlockchainIsEmpty_GenesisTransactionValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.IsEmptyFunc = func() bool { return true }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) {}
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)

	// Act
	pool.Validate(timeMock.Now().UnixNano())

	// Assert
	isTransactionsPoolValidated := len(blockchainMock.AddBlockCalls()) == 1
	test.Assert(t, isTransactionsPoolValidated, "Transactions pool is not validated whereas it should be.")
}

//
//func Test_Validate_InvalidSignature_TransactionNotValidated(t *testing.T) {
//	// Arrange
//	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
//	validatorWalletAddress := validatorWallet.Address()
//	walletA, _ := encryption.DecodeWallet(test.Mnemonic2, test.DerivationPath, "", "")
//	walletAAddress := walletA.Address()
//	registryMock := new(mock.RegistryMock)
//	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	timeMock := new(clock.TimeMock)
//	var now int64 = 1
//	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
//	validationTimer := time.Nanosecond
//	logger := log.NewLogger(log.Fatal)
//	blockchainMock := new(mock.BlockchainMock)
//	blockchainMock.CopyFunc = func() verification.Blockchain { return blockchainMock }
//	blockchainMock.IsEmptyFunc = func() bool { return false }
//	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) {}
//	var blockResponses []*network.BlockResponse
//	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
//	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
//	invalidTransaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
//	_ = invalidTransaction.Sign(walletA.PrivateKey())
//	invalidTransactionRequest := invalidTransaction.GetRequest()
//	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)
//	pool.AddTransaction(&invalidTransactionRequest, nil)
//	pool.Wait()
//
//	// Act
//	pool.Validate(timeMock.Now().UnixNano())
//
//	// Assert
//	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
//}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.IsEmptyFunc = func() bool { return false }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) {}
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	genesisAmount := *invalidTransactionRequest.Value + *invalidTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)
	pool.AddTransaction(&invalidTransactionRequest, nil)
	pool.Wait()

	// Act
	pool.Validate(timeMock.Now().UnixNano())

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_TransactionTimestampIsOlderThan2Blocks_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 3
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.IsEmptyFunc = func() bool { return false }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) {}
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	genesisAmount := *invalidTransactionRequest.Value + *invalidTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)
	pool.AddTransaction(&invalidTransactionRequest, nil)
	pool.Wait()
	blockResponses = append(blockResponses, mock.NewEmptyBlockResponse(now-2))
	blockResponses = append(blockResponses, mock.NewEmptyBlockResponse(now-1))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }

	// Act
	pool.Validate(timeMock.Now().UnixNano())

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_TransactionIsAlreadyInTheBlockchain_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 2
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.IsEmptyFunc = func() bool { return false }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) {}
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	genesisAmount := *invalidTransactionRequest.Value + *invalidTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)
	pool.AddTransaction(&invalidTransactionRequest, nil)
	pool.Wait()
	transaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	var transactionResponses []*network.TransactionResponse
	transactionResponses = append(transactionResponses, transaction.GetResponse())
	blockResponses = append(blockResponses, mock.NewBlockResponse(now-1, [32]byte{}, transactionResponses, nil))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }

	// Act
	pool.Validate(timeMock.Now().UnixNano())

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_ValidTransaction_TransactionValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(mock.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	var now int64 = 1
	timeMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.IsEmptyFunc = func() bool { return false }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) {}
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, mock.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	validTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = validTransaction.Sign(validatorWallet.PrivateKey())
	validTransactionRequest := validTransaction.GetRequest()
	genesisAmount := *validTransactionRequest.Value + *validTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, registryMock, validatorWalletAddress, 0, validationTimer, timeMock, logger)
	pool.AddTransaction(&validTransactionRequest, nil)
	pool.Wait()

	// Act
	pool.Validate(timeMock.Now().UnixNano())

	// Assert
	validatedPool := blockchainMock.AddBlockCalls()
	isTransactionsPoolValidated := len(validatedPool) == 1
	test.Assert(t, isTransactionsPoolValidated, "Transactions pool should be validated only once.")
	validatedTransactions := validatedPool[0].Transactions
	isTwoTransactions := len(validatedTransactions) == 2
	test.Assert(t, isTwoTransactions, "Validated transactions pool should contain exactly 2 transactions.")
	actualTransaction := validatedTransactions[0]
	expectedTransaction, _ := validation.NewTransactionFromRequest(&validTransactionRequest)
	test.Assert(t, expectedTransaction.Equals(actualTransaction), "The first validated transaction is not the expected one.")
	rewardTransaction, _ := validation.NewTransactionFromResponse(validatedTransactions[1])
	isRewardTransaction := rewardTransaction.IsReward()
	test.Assert(t, isRewardTransaction, "The second validated transaction should be the reward.")
}

func assertAddBlockCalledWithRewardTransactionOnly(t *testing.T, blockchainMock *mock.BlockchainMock) {
	validatedPool := blockchainMock.AddBlockCalls()
	isTransactionsPoolValidated := len(validatedPool) == 1
	test.Assert(t, isTransactionsPoolValidated, "Transactions pool should be validated only once.")
	validatedTransactions := validatedPool[0].Transactions
	isSingleTransaction := len(validatedTransactions) == 1
	test.Assert(t, isSingleTransaction, "Validated transactions pool should contain only one transaction.")
	transaction, _ := validation.NewTransactionFromResponse(validatedTransactions[0])
	isRewardTransaction := transaction.IsReward()
	test.Assert(t, isRewardTransaction, "The single validated transaction should be the reward.")
}
