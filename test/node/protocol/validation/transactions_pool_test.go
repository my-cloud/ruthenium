package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
	"testing"
	"time"
)

func Test_AddTransaction_TransactionFeeIsTooLow_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var invalidTransactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(invalidTransactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	var minimalTransactionFee uint64 = 1
	pool := validation.NewTransactionsPool(blockchainMock, minimalTransactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_TransactionTimestampIsInTheFuture_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, protocoltest.NewEmptyBlockResponse(now-1))
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_TransactionTimestampIsOlderThan1Blocks_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now-2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockResponses = append(blockResponses, protocoltest.NewEmptyBlockResponse(now-1))
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "")

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
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	transaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	var transactionResponses []*network.TransactionResponse
	transactionResponses = append(transactionResponses, transaction.GetResponse())
	blockResponses = append(blockResponses, verification.NewBlockResponse(now-1, [32]byte{}, transactionResponses, nil, nil))
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "")

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
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 1
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	var amount uint64 = 1
	var transactionFee uint64 = 0
	transaction := server.NewTransaction(transactionFee, walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(walletA.PrivateKey())
	invalidTransactionRequest := transaction.GetRequest()
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "")

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
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 1
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	var amount uint64 = 1
	var transactionFee uint64 = 0
	transaction := server.NewTransaction(transactionFee, walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, amount)
	_ = transaction.Sign(validatorWallet.PrivateKey())
	transactionRequest := transaction.GetRequest()
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return *transactionRequest.Value + *transactionRequest.Fee }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)

	// Act
	pool.AddTransaction(&transactionRequest, "")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 1
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	genesisAmount := *invalidTransactionRequest.Value + *invalidTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)
	pool.AddTransaction(&invalidTransactionRequest, "")

	// Act
	pool.Validate(watchMock.Now().UnixNano())

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_TransactionTimestampIsOlderThan2Blocks_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 3
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	genesisAmount := *invalidTransactionRequest.Value + *invalidTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)
	pool.AddTransaction(&invalidTransactionRequest, "")
	blockResponses = append(blockResponses, protocoltest.NewEmptyBlockResponse(now-2))
	blockResponses = append(blockResponses, protocoltest.NewEmptyBlockResponse(now-1))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }

	// Act
	pool.Validate(watchMock.Now().UnixNano())

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_TransactionIsAlreadyInTheBlockchain_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	invalidTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	genesisAmount := *invalidTransactionRequest.Value + *invalidTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)
	pool.AddTransaction(&invalidTransactionRequest, "")
	transaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	var transactionResponses []*network.TransactionResponse
	transactionResponses = append(transactionResponses, transaction.GetResponse())
	blockResponses = append(blockResponses, verification.NewBlockResponse(now-1, [32]byte{}, transactionResponses, nil, nil))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }

	// Act
	pool.Validate(watchMock.Now().UnixNano())

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_ValidTransaction_TransactionValidated(t *testing.T) {
	// Arrange
	validatorWallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	validatorWalletAddress := validatorWallet.Address()
	registryMock := new(protocoltest.RegistryMock)
	registryMock.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 1
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	var blockResponses []*network.BlockResponse
	blockResponses = append(blockResponses, protocoltest.NewGenesisBlockResponse(validatorWalletAddress))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	validTransaction := server.NewTransaction(transactionFee, "A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = validTransaction.Sign(validatorWallet.PrivateKey())
	validTransactionRequest := validTransaction.GetRequest()
	genesisAmount := *validTransactionRequest.Value + *validTransactionRequest.Fee
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return genesisAmount }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, registryMock, synchronizerMock, validatorWalletAddress, validationTimer, watchMock, logger)
	pool.AddTransaction(&validTransactionRequest, "")

	// Act
	pool.Validate(watchMock.Now().UnixNano())

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

func assertAddBlockCalledWithRewardTransactionOnly(t *testing.T, blockchainMock *protocoltest.BlockchainMock) {
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
