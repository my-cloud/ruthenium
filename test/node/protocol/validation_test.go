package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	verification2 "github.com/my-cloud/ruthenium/test/node/protocol/verification"
	"testing"
	"time"
)

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
	blockchain := verification.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := validation.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewEngine(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	blockchain.AddBlock(verification2.NewGenesisBlockResponse(validatorWalletAddress))
	invalidTransaction := server.NewTransaction(walletAAddress, validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(walletA.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Act
	validation.Do()
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
	blockchain := verification.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := validation.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewEngine(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	blockchain.AddBlock(verification2.NewGenesisBlockResponse(validatorWalletAddress))
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now+2, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Act
	validation.Do()
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
	blockchain := verification.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := validation.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewEngine(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	blockchain.AddBlock(verification2.NewGenesisBlockResponse(validatorWalletAddress))
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now-3, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()
	blockchain.AddBlock(verification2.NewEmptyBlockResponse(now - 2))
	blockchain.AddBlock(verification2.NewEmptyBlockResponse(now - 1))

	// Act
	validation.Do()
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
	blockchain := verification.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := validation.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewEngine(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	blockchain.AddBlock(verification2.NewGenesisBlockResponse(validatorWalletAddress))
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()
	transaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	blockchain.AddBlock(verification2.NewBlockResponse(now-1, [32]byte{}, transaction))

	// Act
	validation.Do()
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
	blockchain := verification.NewBlockchain(registryMock, validationTimer, timeMock, synchronizerMock, logger)
	pool := validation.NewTransactionsPool(registryMock, validationTimer, timeMock, logger)
	validation := protocol.NewEngine(validatorWalletAddress, blockchain, pool, timeMock, validationTimer, logger)
	blockchain.AddBlock(verification2.NewGenesisBlockResponse(validatorWalletAddress))
	invalidTransaction := server.NewTransaction("A", validatorWalletAddress, validatorWallet.PublicKey(), now, 1)
	_ = invalidTransaction.Sign(validatorWallet.PrivateKey())
	invalidTransactionRequest := invalidTransaction.GetRequest()
	pool.AddTransaction(&invalidTransactionRequest, blockchain, nil)
	pool.Wait()

	// Act
	validation.Do()
	validation.Wait()

	// Assert
	blocks := blockchain.Blocks()
	lastBlock := blocks[len(blocks)-1]
	expectedTransactionsLength := 2
	actualTransactionsLength := len(lastBlock.Transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
}

func assertTransactionNotValidated(t *testing.T, blockchain *verification.Blockchain) {
	blocks := blockchain.Blocks()
	lastBlock := blocks[len(blocks)-1]
	expectedTransactionsLength := 1
	actualTransactionsLength := len(lastBlock.Transactions)
	test.Assert(t, actualTransactionsLength == expectedTransactionsLength, fmt.Sprintf("Wrong transactions count. Expected: %d - Actual: %d", expectedTransactionsLength, actualTransactionsLength))
	actualTransaction, _ := validation.NewTransactionFromResponse(lastBlock.Transactions[0])
	test.Assert(t, actualTransaction.IsReward(), "The single expected transaction is not a reward whereas it should be.")
}
