package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock"
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
	logger := log.NewLogger(log.Fatal)
	watch := clock.NewWatch()
	validationTimer := time.Nanosecond
	blockchain := protocol.NewBlockchain(validationTimer.Nanoseconds(), watch, logger)
	pool := protocol.NewPool(watch, logger)
	validation := protocol.NewValidation(minerWalletAddress, blockchain, pool, watch, validationTimer, logger)
	validation.Do()

	// Act
	var amount1 uint64 = 40 * protocol.ParticlesCount
	transaction1 := protocol.NewTransaction(walletAAddress, minerWalletAddress, minerWallet.PublicKey(), 0, amount1)
	_ = transaction1.Sign(minerWallet.PrivateKey())
	addTransaction(pool, blockchain, validation, transaction1)

	var amount2 uint64 = 10 * protocol.ParticlesCount
	transaction2 := protocol.NewTransaction(walletBAddress, walletAAddress, walletA.PublicKey(), 0, amount2)
	_ = transaction2.Sign(walletA.PrivateKey())
	addTransaction(pool, blockchain, validation, transaction2)

	// Assert
	expectedWalletAAmount := amount1 - amount2 - transaction2.Fee()
	actualWalletAAmount := blockchain.CalculateTotalAmount(4, walletAAddress)
	test.Assert(t, expectedWalletAAmount == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", expectedWalletAAmount, actualWalletAAmount))
	expectedWalletBAmount := amount2
	actualWalletBAmount := blockchain.CalculateTotalAmount(4, walletBAddress)
	test.Assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func Test_AddTransaction_NotAllowed_TransactionNotAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := encryption.NewWallet()
	minerWalletAddress := minerWallet.Address()
	walletA, _ := encryption.NewWallet()
	walletAAddress := walletA.Address()
	walletB, _ := encryption.NewWallet()
	walletBAddress := walletB.Address()
	logger := log.NewLogger(log.Fatal)
	watch := clock.NewWatch()
	validationTimer := time.Nanosecond
	blockchain := protocol.NewBlockchain(validationTimer.Nanoseconds(), watch, logger)
	pool := protocol.NewPool(watch, logger)
	validation := protocol.NewValidation(minerWalletAddress, blockchain, pool, watch, validationTimer, logger)
	validation.Do()

	// Act
	var amount1 uint64 = 40 * protocol.ParticlesCount
	transaction1 := protocol.NewTransaction(walletAAddress, minerWalletAddress, minerWallet.PublicKey(), 0, amount1)
	_ = transaction1.Sign(minerWallet.PrivateKey())
	addTransaction(pool, blockchain, validation, transaction1)

	var amount2 uint64 = 10 * protocol.ParticlesCount
	transaction2 := protocol.NewTransaction(walletBAddress, walletAAddress, walletA.PublicKey(), 0, amount2)
	_ = transaction2.Sign(walletB.PrivateKey())
	addTransaction(pool, blockchain, validation, transaction2)

	// Assert
	actualWalletAAmount := blockchain.CalculateTotalAmount(4, walletAAddress)
	test.Assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	var expectedWalletBAmount uint64 = 0
	actualWalletBAmount := blockchain.CalculateTotalAmount(4, walletBAddress)
	test.Assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func addTransaction(pool *protocol.Pool, blockchain *protocol.Blockchain, validation *protocol.Validation, transaction *protocol.Transaction) {
	pool.AddTransaction(transaction, blockchain, nil)
	pool.Wait()
	validation.Do()
	validation.Wait()
}
