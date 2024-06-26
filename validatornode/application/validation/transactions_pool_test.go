package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_AddTransaction_TransactionTimestampIsInTheFuture_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return nil }
	sendersManagerMock.IncentiveFunc = func(string) {}
	watchMock := new(application.TimeProviderMock)
	var now int64 = 2
	watchMock.NowFunc = func() time.Time { return time.Unix(0, now) }
	logger := log.NewLoggerMock()
	transactionFee := 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, validatorWalletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+2, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transaction, "0", "0")

	// Assert
	expectedTransactionsLength := 0
	actualTransactionsLength := len(pool.Transactions())
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, logger.DebugCalls(), "failed to add transaction: the transaction timestamp is too far in the future")
}

func Test_AddTransaction_TransactionTimestampIsTooOld_TransactionNotAdded(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return nil }
	sendersManagerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	logger := log.NewLoggerMock()
	transactionFee := 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, validatorWalletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue, false)

	// Act
	pool.AddTransaction(transaction, "0", "0")

	// Assert
	expectedTransactionsLength := 0
	actualTransactionsLength := len(pool.Transactions())
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, logger.DebugCalls(), "failed to add transaction: the transaction timestamp is too old")
}

func Test_AddTransaction_InvalidSignature_TransactionNotAdded(t *testing.T) {
	// Arrange
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return nil }
	sendersManagerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	transactionFee := 0
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, walletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, outputIndex, "A", privateKey2, publicKey, now, transactionId, genesisValue, false)

	// Act
	pool.AddTransaction(transaction, "0", "0")

	// Assert
	expectedTransactionsLength := 0
	actualTransactionsLength := len(pool.Transactions())
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	test.AssertThatMessageIsLogged(t, logger.DebugCalls(), "failed to add transaction: failed to verify signature")
}

func Test_AddTransaction_ValidTransaction_TransactionAdded(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.AddTransactionFunc = func([]byte) error { return nil }
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return []application.Sender{senderMock} }
	sendersManagerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	transactionFee := 0
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() application.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*ledger.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(*ledger.Transaction, int64) (uint64, error) { return 0, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	var outputIndex uint16 = 0
	transactionId := ""
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, walletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, outputIndex, walletAddress, privateKey, publicKey, now, transactionId, genesisValue, false)

	// Act
	pool.AddTransaction(transaction, "0", "0")

	// Assert
	expectedTransactionsLength := 1
	actualTransactionsLength := len(pool.Transactions())
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func Test_Validate_BlockAlreadyExist_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	sendersManagerMock := new(application.SendersManagerMock)
	var now int64 = 2
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, validatorWalletAddress, logger)

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.ErrorCalls(), "unable to create block, a block with the same timestamp is already in the blockchain")
}

func Test_Validate_BlockIsMissing_TransactionsNotValidated(t *testing.T) {
	// Arrange
	validatorWalletAddress := test.Address
	sendersManagerMock := new(application.SendersManagerMock)
	var now int64 = 3
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, validatorWalletAddress, logger)

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.ErrorCalls(), "unable to create block, a block is missing in the blockchain")
}

func Test_Validate_TransactionTimestampIsInTheFuture_TransactionsNotValidated(t *testing.T) {
	// Arrange
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return nil }
	sendersManagerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	transactionFee := 0
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() application.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*ledger.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(*ledger.Transaction, int64) (uint64, error) { return 0, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, walletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now+1, "0", genesisValue, false)
	pool.AddTransaction(transaction, "0", "0")
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.WarnCalls(), "transaction removed from the transactions pool, the transaction timestamp is too far in the future")
}

func Test_Validate_TransactionTimestampIsTooOld_TransactionsNotValidated(t *testing.T) {
	// Arrange
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return nil }
	sendersManagerMock.IncentiveFunc = func(string) {}
	var now int64 = 3
	transactionFee := 0
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 2 }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() application.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*ledger.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(*ledger.Transaction, int64) (uint64, error) { return 0, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, walletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now-2, "0", genesisValue, false)
	pool.AddTransaction(transaction, "0", "0")
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }

	// Act
	pool.Validate(now)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.WarnCalls(), "transaction removed from the transactions pool, the transaction timestamp is too old")
}

func Test_Validate_ValidTransaction_TransactionsValidated(t *testing.T) {
	// Arrange
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []application.Sender { return nil }
	sendersManagerMock.IncentiveFunc = func(string) {}
	var now int64 = 2
	transactionFee := 0
	logger := log.NewLoggerMock()
	blocksManagerMock := new(application.BlocksManagerMock)
	blocksManagerMock.LastBlockTransactionsFunc = func() []*ledger.Transaction { return nil }
	blocksManagerMock.LastBlockTimestampFunc = func() int64 { return now - 1 }
	blocksManagerMock.AddBlockFunc = func(int64, []*ledger.Transaction, []string) error { return nil }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	utxosManagerMock := new(application.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() application.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*ledger.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(*ledger.Transaction, int64) (uint64, error) { return 0, nil }
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	walletAddress := publicKey.Address()
	pool := NewTransactionsPool(blocksManagerMock, settings, sendersManagerMock, utxosManagerMock, walletAddress, logger)
	var genesisValue uint64 = 0
	transaction := ledger.NewSignedTransaction(genesisValue, transactionFee, 0, "A", privateKey, publicKey, now, "0", genesisValue, false)
	pool.AddTransaction(transaction, "0", "0")

	// Act
	pool.Validate(now)

	// Assert
	addBlockCalls := blocksManagerMock.AddBlockCalls()
	test.Assert(t, len(addBlockCalls) == 1, fmt.Sprintf("AddBlock method should be called only once whereas it's called %d times", len(addBlockCalls)))
	transactions := addBlockCalls[0].Transactions
	isTwoTransactions := len(transactions) == 2
	test.Assert(t, isTwoTransactions, "Validated transactions pool should contain exactly 2 transactions.")
	actualTransaction := transactions[0]
	test.Assert(t, actualTransaction.Equals(transaction), "The first validated transaction is not the expected one.")
	rewardTransaction := transactions[1]
	isRewardTransaction := rewardTransaction.HasReward()
	test.Assert(t, isRewardTransaction, "The second validated transaction should be the reward.")
}
