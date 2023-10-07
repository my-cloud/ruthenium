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
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	var genesisValue uint64 = 0
	settings := new(protocoltest.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+2, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

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
	logger := logtest.NewLoggerMock()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	var genesisValue uint64 = 0
	settings := new(protocoltest.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

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
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	settings := new(protocoltest.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey2, publicKey, now, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

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
	neighborMock.AddTransactionFunc = func([]byte) error { return nil }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return []network.Neighbor{neighborMock} }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.UtxoFunc = func(input protocol.InputInfo) (protocol.Utxo, error) {
		inputInfo := verification.NewInputInfo(0, "")
		return verification.NewUtxo(inputInfo, &verification.Output{}, 0), nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	settings := new(protocoltest.SettingsMock)
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	var outputIndex uint16 = 0
	transactionId := "0"
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, outputIndex, walletAAddress, privateKey, publicKey, now, transactionId, genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*verification.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_BlockAlreadyExist_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	var now int64 = 2
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now }
	settings := new(protocoltest.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"unable to create block, a block with the same timestamp is already in the blockchain"}, logger.ErrorCalls())
}

func Test_Validate_BlockIsMissing_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	var now int64 = 3
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	settings := new(protocoltest.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"unable to create block, a block is missing in the blockchain"}, logger.ErrorCalls())
}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.UtxoFunc = func(input protocol.InputInfo) (protocol.Utxo, error) {
		inputInfo := verification.NewInputInfo(0, "")
		return verification.NewUtxo(inputInfo, &verification.Output{}, 0), nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	settings := new(protocoltest.SettingsMock)
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+1, "0", genesisValue, false)
	pool.AddTransaction(transactionRequest, "0")
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"transaction removed from the transactions pool, the transaction timestamp is too far in the future"}, logger.WarnCalls())
}

func Test_Validate_TransactionTimestampIsTooOld_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 3
	var transactionFee uint64 = 0
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.UtxoFunc = func(input protocol.InputInfo) (protocol.Utxo, error) {
		inputInfo := verification.NewInputInfo(0, "")
		return verification.NewUtxo(inputInfo, &verification.Output{}, 0), nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	settings := new(protocoltest.SettingsMock)
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue, false)
	pool.AddTransaction(transactionRequest, "0")
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"transaction removed from the transactions pool, the transaction timestamp is too old"}, logger.WarnCalls())
}

func Test_Validate_ValidTransaction_TransactionsValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := logtest.NewLoggerMock()
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	blockchainMock.UtxoFunc = func(input protocol.InputInfo) (protocol.Utxo, error) {
		inputInfo := verification.NewInputInfo(0, "")
		return verification.NewUtxo(inputInfo, &verification.Output{}, 0), nil
	}
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var genesisValue uint64 = 0
	settings := new(protocoltest.SettingsMock)
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := validation.NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now, "0", genesisValue, false)
	pool.AddTransaction(transactionRequest, "0")

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
	var expectedTransaction *validation.TransactionRequest
	_ = json.Unmarshal(transactionRequest, &expectedTransaction)
	test.Assert(t, actualTransaction.Equals(expectedTransaction.Transaction()), "The first validated transaction is not the expected one.")
	rewardTransaction := transactions[1]
	isRewardTransaction := rewardTransaction.HasReward()
	test.Assert(t, isRewardTransaction, "The second validated transaction should be the reward.")
}
