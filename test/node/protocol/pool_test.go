package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node"
	"testing"
	"time"
)

func Test_AddTransaction_Allowed_TransactionAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := encryption.NewWallet()
	minerWalletAddress := minerWallet.Address()
	walletA, _ := encryption.NewWallet()
	walletAAddress := walletA.Address()
	walletB, _ := encryption.NewWallet()
	walletBAddress := walletB.Address()
	registrable := NewRegistrableMock()
	IsRegisteredMock = func(address string) (bool, error) { return true, nil }
	watch := node.NewWatchMock()
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchain := protocol.NewBlockchain(registrable, validationTimer, watch, logger)
	pool := protocol.NewPool(registrable, watch, logger)
	validation := protocol.NewValidation(minerWalletAddress, blockchain, pool, watch, validationTimer, logger)
	validation.Do()
	validation.Wait()

	var amount1 uint64 = 40 * network.ParticlesCount
	transaction1 := server.NewTransaction(walletAAddress, minerWalletAddress, minerWallet.PublicKey(), 0, amount1)
	_ = transaction1.Sign(minerWallet.PrivateKey())

	var amount2 uint64 = 10 * network.ParticlesCount
	transaction2 := server.NewTransaction(walletBAddress, walletAAddress, walletA.PublicKey(), 0, amount2)
	_ = transaction2.Sign(walletA.PrivateKey())

	// Act
	addTransaction(pool, blockchain, validation, transaction1)
	actualWalletAAmount := blockchain.CalculateTotalAmount(0, walletAAddress)

	addTransaction(pool, blockchain, validation, transaction2)
	actualWalletBAmount := blockchain.CalculateTotalAmount(0, walletBAddress)

	// Assert
	test.Assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	test.Assert(t, amount2 == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", amount2, actualWalletBAmount))
}

func Test_AddTransaction_NotAllowed_TransactionNotAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := encryption.NewWallet()
	minerWalletAddress := minerWallet.Address()
	walletA, _ := encryption.NewWallet()
	walletAAddress := walletA.Address()
	walletB, _ := encryption.NewWallet()
	walletBAddress := walletB.Address()
	registrable := NewRegistrableMock()
	IsRegisteredMock = func(address string) (bool, error) { return true, nil }
	watch := node.NewWatchMock()
	validationTimer := time.Nanosecond
	logger := log.NewLogger(log.Fatal)
	blockchain := protocol.NewBlockchain(registrable, validationTimer, watch, logger)
	pool := protocol.NewPool(registrable, watch, logger)
	validation := protocol.NewValidation(minerWalletAddress, blockchain, pool, watch, validationTimer, logger)
	validation.Do()
	validation.Wait()

	var amount1 uint64 = 40 * network.ParticlesCount
	transaction1 := server.NewTransaction(walletAAddress, minerWalletAddress, minerWallet.PublicKey(), 0, amount1)
	_ = transaction1.Sign(minerWallet.PrivateKey())

	var amount2 uint64 = 10 * network.ParticlesCount
	transaction2 := server.NewTransaction(walletBAddress, walletAAddress, walletA.PublicKey(), 0, amount2)
	_ = transaction2.Sign(walletB.PrivateKey())

	// Act
	addTransaction(pool, blockchain, validation, transaction1)
	actualWalletAAmount := blockchain.CalculateTotalAmount(0, walletAAddress)

	addTransaction(pool, blockchain, validation, transaction2)
	actualWalletBAmount := blockchain.CalculateTotalAmount(0, walletBAddress)

	// Assert
	test.Assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	var expectedWalletBAmount uint64 = 0
	test.Assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func addTransaction(pool *protocol.Pool, blockchain *protocol.Blockchain, validation *protocol.Validation, transaction *server.Transaction) {
	transactionRequest := transaction.GetRequest()
	pool.AddTransaction(&transactionRequest, blockchain, nil)
	pool.Wait()
	validation.Do()
	validation.Wait()
}
