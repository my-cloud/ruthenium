package validation

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/domain"
	"github.com/my-cloud/ruthenium/domain/validatornode"
	"github.com/my-cloud/ruthenium/infrastructure/clock"
	"github.com/my-cloud/ruthenium/infrastructure/encryption"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"github.com/my-cloud/ruthenium/infrastructure/network"
	"github.com/my-cloud/ruthenium/infrastructure/test"
	"testing"
	"time"
)

func Test_AddTransaction_TransactionTimestampIsInTheFuture_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	watchMock := new(clock.WatchMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	logger := log.NewLoggerMock()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+2, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*domain.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: the transaction timestamp is too far in the future"}, logger.DebugCalls())
}

func Test_AddTransaction_TransactionTimestampIsTooOld_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	logger := log.NewLoggerMock()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*domain.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: the transaction timestamp is too old"}, logger.DebugCalls())
}

func Test_AddTransaction_InvalidSignature_TransactionNotAdded(t *testing.T) {
	// Arrange
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	blockchainMock.UtxoFunc = func(input validatornode.InputInfo) (validatornode.Utxo, error) {
		inputInfo := domain.NewInputInfo(outputIndex, transactionId)
		return domain.NewUtxo(inputInfo, domain.NewOutput(walletAddress, false, 0), 0), nil
	}
	var genesisValue uint64 = 0
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	settings := new(validatornode.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, walletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, outputIndex, "A", privateKey2, publicKey, now, transactionId, genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*domain.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: failed to verify transaction: failed to verify signature"}, logger.DebugCalls())
}

func Test_AddTransaction_InvalidPublicKey_TransactionNotAdded(t *testing.T) {
	// Arrange
	walletAddress2 := test.Address2
	neighborMock := new(network.NeighborMock)
	neighborMock.AddTransactionFunc = func([]byte) error { return nil }
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return []network.Neighbor{neighborMock} }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	blockchainMock.UtxoFunc = func(input validatornode.InputInfo) (validatornode.Utxo, error) {
		inputInfo := domain.NewInputInfo(outputIndex, transactionId)
		return domain.NewUtxo(inputInfo, domain.NewOutput(walletAddress2, false, 0), 0), nil
	}
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, walletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, outputIndex, walletAddress2, privateKey, publicKey, now, transactionId, genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*domain.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 0
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, []string{"failed to add transaction: failed to verify transaction: output address does not derive from input public key"}, logger.DebugCalls())
}

func Test_AddTransaction_ValidTransaction_TransactionAdded(t *testing.T) {
	// Arrange
	walletAAddress := test.Address2
	neighborMock := new(network.NeighborMock)
	neighborMock.AddTransactionFunc = func([]byte) error { return nil }
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return []network.Neighbor{neighborMock} }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	blockchainMock.UtxoFunc = func(input validatornode.InputInfo) (validatornode.Utxo, error) {
		inputInfo := domain.NewInputInfo(outputIndex, transactionId)
		return domain.NewUtxo(inputInfo, domain.NewOutput(walletAddress, false, 0), 0), nil
	}
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, walletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, outputIndex, walletAAddress, privateKey, publicKey, now, transactionId, genesisValue, false)

	// Act
	pool.AddTransaction(transactionRequest, "0")

	// Assert
	transactionsBytes := pool.Transactions()
	var transactions []*domain.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	expectedTransactionsLength := 1
	actualTransactionsLength := len(transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_BlockAlreadyExist_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(network.SynchronizerMock)
	var now int64 = 2
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now }
	settings := new(validatornode.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"unable to create block, a block with the same timestamp is already in the blockchain"}, logger.ErrorCalls())
}

func Test_Validate_BlockIsMissing_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	synchronizerMock := new(network.SynchronizerMock)
	var now int64 = 3
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	settings := new(validatornode.SettingsMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, validatorWalletAddress, logger)

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"unable to create block, a block is missing in the blockchain"}, logger.ErrorCalls())
}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionsNotValidated(t *testing.T) {
	// Arrange
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	blockchainMock.UtxoFunc = func(input validatornode.InputInfo) (validatornode.Utxo, error) {
		inputInfo := domain.NewInputInfo(outputIndex, transactionId)
		return domain.NewUtxo(inputInfo, domain.NewOutput(walletAddress, false, 0), 0), nil
	}
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, walletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+1, "0", genesisValue, false)
	pool.AddTransaction(transactionRequest, "0")
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"transaction removed from the transactions pool, the transaction timestamp is too far in the future"}, logger.WarnCalls())
}

func Test_Validate_TransactionTimestampIsTooOld_TransactionsNotValidated(t *testing.T) {
	// Arrange
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 3
	var transactionFee uint64 = 0
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	blockchainMock.UtxoFunc = func(input validatornode.InputInfo) (validatornode.Utxo, error) {
		inputInfo := domain.NewInputInfo(outputIndex, transactionId)
		return domain.NewUtxo(inputInfo, domain.NewOutput(walletAddress, false, 0), 0), nil
	}
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, walletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue, false)
	pool.AddTransaction(transactionRequest, "0")
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, []string{"transaction removed from the transactions pool, the transaction timestamp is too old"}, logger.WarnCalls())
}

func Test_Validate_ValidTransaction_TransactionsValidated(t *testing.T) {
	// Arrange
	synchronizerMock := new(network.SynchronizerMock)
	synchronizerMock.NeighborsFunc = func() []network.Neighbor { return nil }
	synchronizerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	var transactionFee uint64 = 0
	logger := log.NewLoggerMock()
	blockchainMock := new(validatornode.BlockchainMock)
	blockchainMock.CopyFunc = func() validatornode.Blockchain { return blockchainMock }
	blockchainMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blockchainMock.AddBlockFunc = func(int64, []byte, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	blockchainMock.UtxoFunc = func(input validatornode.InputInfo) (validatornode.Utxo, error) {
		inputInfo := domain.NewInputInfo(outputIndex, transactionId)
		return domain.NewUtxo(inputInfo, domain.NewOutput(walletAddress, false, 0), 0), nil
	}
	var genesisValue uint64 = 0
	settings := new(validatornode.SettingsMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	pool := NewTransactionsPool(blockchainMock, settings, synchronizerMock, walletAddress, logger)
	transactionRequest := domain.NewSignedTransactionRequest(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now, "0", genesisValue, false)
	pool.AddTransaction(transactionRequest, "0")

	// Act
	pool.Validate(now)

	// Assert
	validatedPool := blockchainMock.AddBlockCalls()
	expectedCallsCount := 7
	isTransactionsPoolValidated := len(validatedPool) == expectedCallsCount
	test.Assert(t, isTransactionsPoolValidated, fmt.Sprintf("AddBlock method should be called only %d times whereas it's called %d times", expectedCallsCount, len(validatedPool)))
	transactionsBytes := validatedPool[expectedCallsCount-1].Transactions
	var transactions []*domain.Transaction
	_ = json.Unmarshal(transactionsBytes, &transactions)
	isTwoTransactions := len(transactions) == 2
	test.Assert(t, isTwoTransactions, "Validated transactions pool should contain exactly 2 transactions.")
	actualTransaction := transactions[0]
	var expectedTransaction *TransactionRequest
	_ = json.Unmarshal(transactionRequest, &expectedTransaction)
	test.Assert(t, actualTransaction.Equals(expectedTransaction.Transaction()), "The first validated transaction is not the expected one.")
	rewardTransaction := transactions[1]
	isRewardTransaction := rewardTransaction.HasReward()
	test.Assert(t, isRewardTransaction, "The second validated transaction should be the reward.")
}
