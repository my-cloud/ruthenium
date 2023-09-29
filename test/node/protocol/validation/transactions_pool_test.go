package validation

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
	"testing"
	"time"
)

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
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	var genesisValue uint64 = 0
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+2, "0", genesisValue)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*verification.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: the transaction timestamp is too far in the future"}, logger.DebugCalls())
}

func Test_AddTransaction_TransactionTimestampIsTooOld_TransactionNotAdded(t *testing.T) {
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
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	var genesisValue uint64 = 0
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*verification.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: the transaction timestamp is too old"}, logger.DebugCalls())
}

func Test_AddTransaction_InvalidSignature_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey2, publicKey, now, "0", genesisValue)

	// Act
	pool.AddTransaction(&invalidTransactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*verification.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: failed to verify transaction: failed to verify signature"}, logger.DebugCalls())
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
	var now int64 = 2
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	var outputIndex uint16 = 0
	transactionId := "0"
	blockchainMock.UtxosFunc = func(string) []byte {
		utxos := []*network.UtxoResponse{
			{
				Address:       walletAAddress,
				HasReward:     true,
				HasIncome:     true,
				OutputIndex:   outputIndex,
				TransactionId: transactionId,
				Value:         genesisValue,
			},
		}
		marshalledUtxos, _ := json.Marshal(utxos)
		return marshalledUtxos
	}
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, outputIndex, walletAAddress, privateKey, publicKey, now, transactionId, genesisValue)

	// Act
	pool.AddTransaction(&transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*verification.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
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
	var now int64 = 2
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+1, "0", genesisValue)
	pool.AddTransaction(&invalidTransactionRequest, "0")
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"transaction removed from the transactions pool, the transaction timestamp is too far in the future"}, logger.WarnCalls())
}

func Test_Validate_TransactionTimestampIsTooOld_TransactionNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 3
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue)
	pool.AddTransaction(&invalidTransactionRequest, "0")
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"transaction removed from the transactions pool, the transaction timestamp is too old"}, logger.WarnCalls())
}

func Test_Validate_ValidTransaction_TransactionValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	validationTimer := time.Nanosecond
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.FindFeeFunc = func([]*network.InputResponse, []*network.OutputResponse, int64) (uint64, error) {
		return transactionFee, nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	pool := validation.NewTransactionsPool(blockchainMock, genesisValue, transactionFee, synchronizerMock, validatorWalletAddress, validationTimer, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now, "0", genesisValue)
	pool.AddTransaction(&transactionRequest, "0")

	// Act
	pool.Validate(now)

	// Assert
	validatedPool := blockchainMock.AddBlockCalls()
	expectedCallsCount := 7
	isTransactionsPoolValidated := len(validatedPool) == expectedCallsCount
	test.Assert(t, isTransactionsPoolValidated, fmt.Sprintf("AddBlock method should be called only %d times whereas it's called %d times", expectedCallsCount, len(validatedPool)))
	transactionsBytes := validatedPool[expectedCallsCount-1].Transactions
	var transactions []*verification.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	isTwoTransactions := len(transactions) == 2
	test.Assert(t, isTwoTransactions, "Validated transactions pool should contain exactly 2 transactions.")
	actualTransaction := transactions[0]
	expectedTransaction, _ := verification.NewTransactionFromRequest(&transactionRequest)
	test.Assert(t, expectedTransaction.Equals(actualTransaction), "The first validated transaction is not the expected one.")
	rewardTransaction := transactions[1]
	isRewardTransaction := rewardTransaction.HasReward()
	test.Assert(t, isRewardTransaction, "The second validated transaction should be the reward.")
}
