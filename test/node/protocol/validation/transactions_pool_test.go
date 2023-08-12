package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
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
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 1
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var invalidTransactionFee uint64 = 0
	var minimalTransactionFee uint64 = 1
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse}
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return invalidTransactionFee, nil }
	pool := validation.NewTransactionsPool(blockchainMock, minimalTransactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, invalidTransactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now, 1)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_TransactionTimestampIsInTheFuture_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clocktest.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse, protocoltest.NewEmptyBlockResponse(now - 1)}
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now+2, 1)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_TransactionTimestampIsOlderThan1Blocks_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse, protocoltest.NewEmptyBlockResponse(now - 1)}
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now-2, 1)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_InvalidSignature_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	walletAAddress := test.Address2
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 1
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse}
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	blockchainMock.UtxosByAddressFunc = func(string) []*network.UtxoResponse {
		return []*network.UtxoResponse{protocoltest.NewUtxoFromOutput(genesisTransaction, genesisOutputIndex)}
	}
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, walletAAddress, genesisTransaction, genesisOutputIndex, privateKey2, publicKey, now, 1)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_AddTransaction_ValidTransaction_TransactionAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	walletAAddress := test.Address2
	neighborMock := new(networktest.NeighborMock)
	neighborMock.AddTransactionFunc = func(network.TransactionRequest) error { return nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return []network.Neighbor{neighborMock} }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 1
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse}
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	blockchainMock.UtxosByAddressFunc = func(string) []*network.UtxoResponse {
		return []*network.UtxoResponse{protocoltest.NewUtxoFromOutput(genesisTransaction, genesisOutputIndex)}
	}
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, walletAAddress, genesisTransaction, genesisOutputIndex, privateKey, publicKey, now, 1)

	// Act
	pool.AddTransaction(&transactionRequest, "0")

	// Assert
	transactions := pool.Transactions()
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 1
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse}
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now+2, 1)
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Act
	pool.Validate(now)

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_TransactionTimestampIsOlderThan2Blocks_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 3
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse}
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now-3, 1)
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	pool.AddTransaction(&invalidTransactionRequest, "0")
	blockResponses = append(blockResponses, protocoltest.NewEmptyBlockResponse(now-2))
	blockResponses = append(blockResponses, protocoltest.NewEmptyBlockResponse(now-1))
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }

	// Act
	pool.Validate(now)

	// Assert
	assertAddBlockCalledWithRewardTransactionOnly(t, blockchainMock)
}

func Test_Validate_ValidTransaction_TransactionValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 1
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []*network.TransactionResponse, []string) error { return nil }
	genesisBlockResponse := protocoltest.NewGenesisBlockResponse(validatorWalletAddress)
	blockResponses := []*network.BlockResponse{genesisBlockResponse}
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return blockResponses }
	var transactionFee uint64 = 0
	blockchainMock.FindFeeFunc = func(*network.TransactionResponse, int, int64) (uint64, error) { return transactionFee, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	genesisTransaction := genesisBlockResponse.Transactions[0]
	var genesisOutputIndex uint16 = 0
	blockchainMock.UtxosByAddressFunc = func(string) []*network.UtxoResponse {
		return []*network.UtxoResponse{protocoltest.NewUtxoFromOutput(genesisTransaction, genesisOutputIndex)}
	}
	genesisValue := genesisTransaction.Outputs[genesisOutputIndex].Value
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now, 1)
	pool := validation.NewTransactionsPool(blockchainMock, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	pool.AddTransaction(&transactionRequest, "0")

	// Act
	pool.Validate(now)

	// Assert
	validatedPool := blockchainMock.AddBlockCalls()
	expectedCallsCount := 7
	isTransactionsPoolValidated := len(validatedPool) == expectedCallsCount
	test.Assert(t, isTransactionsPoolValidated, fmt.Sprintf("AddBlock method should be called only %d times whereas it's called %d times", expectedCallsCount, len(validatedPool)))
	validatedTransactions := validatedPool[expectedCallsCount-1].Transactions
	isTwoTransactions := len(validatedTransactions) == 2
	test.Assert(t, isTwoTransactions, "Validated transactions pool should contain exactly 2 transactions.")
	actualTransaction := validatedTransactions[0]
	expectedTransaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	test.Assert(t, expectedTransaction.Equals(actualTransaction), "The first validated transaction is not the expected one.")
	rewardTransaction, _ := validation.NewTransactionFromResponse(validatedTransactions[1])
	isRewardTransaction := rewardTransaction.HasReward()
	test.Assert(t, isRewardTransaction, "The second validated transaction should be the reward.")
}

func assertAddBlockCalledWithRewardTransactionOnly(t *testing.T, blockchainMock *protocoltest.BlockchainMock) {
	validatedPool := blockchainMock.AddBlockCalls()
	expectedCallsCount := 7
	isTransactionsPoolValidated := len(validatedPool) == expectedCallsCount
	test.Assert(t, isTransactionsPoolValidated, fmt.Sprintf("AddBlock method should be called only %d times whereas it's called %d times", expectedCallsCount, len(validatedPool)))
	validatedTransactions := validatedPool[expectedCallsCount-1].Transactions
	isSingleTransaction := len(validatedTransactions) == 1
	test.Assert(t, isSingleTransaction, "Validated transactions pool should contain only one transaction.")
	transaction, _ := validation.NewTransactionFromResponse(validatedTransactions[0])
	isRewardTransaction := transaction.HasReward()
	test.Assert(t, isRewardTransaction, "The single validated transaction should be the reward.")
}
