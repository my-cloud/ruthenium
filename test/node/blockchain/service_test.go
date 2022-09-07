package blockchain

import (
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/log"
	"gitlab.com/coinsmaster/ruthenium/src/node/blockchain"
	"gitlab.com/coinsmaster/ruthenium/src/node/encryption"
	"gitlab.com/coinsmaster/ruthenium/test"
	"sync"
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
	service := blockchain.NewService(minerWalletAddress, "", 0, time.Nanosecond, logger)
	wg := service.WaitGroup()
	var amount1 uint64 = 40 * blockchain.ParticlesCount
	for service.CalculateTotalAmount(0, minerWallet.Address()) < amount1 {
		service.Mine()
		wg.Wait()
	}

	// Act
	transaction1 := blockchain.NewTransaction(0, minerWalletAddress, walletAAddress, amount1)
	signature1, _ := transaction1.Sign(minerWallet.PrivateKey())
	addTransaction(service, minerWallet.PublicKey(), transaction1, signature1)

	var amount2 uint64 = 10 * blockchain.ParticlesCount
	transaction2 := blockchain.NewTransaction(0, walletAAddress, walletBAddress, amount2)
	signature2, _ := transaction2.Sign(walletA.PrivateKey())
	addTransaction(service, walletA.PublicKey(), transaction2, signature2)

	// Assert
	expectedWalletAAmount := amount1 - amount2
	actualWalletAAmount := service.CalculateTotalAmount(0, walletAAddress)
	test.Assert(t, expectedWalletAAmount == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", expectedWalletAAmount, actualWalletAAmount))
	expectedWalletBAmount := amount2
	actualWalletBAmount := service.CalculateTotalAmount(0, walletBAddress)
	test.Assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func Test_AddTransaction_Allowed_TransactionNotAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := encryption.NewWallet()
	minerWalletAddress := minerWallet.Address()
	walletA, _ := encryption.NewWallet()
	walletAAddress := walletA.Address()
	walletB, _ := encryption.NewWallet()
	walletBAddress := walletB.Address()
	logger := log.NewLogger(log.Fatal)
	service := blockchain.NewService(minerWalletAddress, "", 0, time.Nanosecond, logger)
	wg := service.WaitGroup()
	var amount1 uint64 = 40 * blockchain.ParticlesCount
	for service.CalculateTotalAmount(0, minerWallet.Address()) < amount1 {
		service.Mine()
		wg.Wait()
	}

	// Act
	transaction1 := blockchain.NewTransaction(0, minerWalletAddress, walletAAddress, amount1)
	signature1, _ := transaction1.Sign(minerWallet.PrivateKey())
	addTransaction(service, minerWallet.PublicKey(), transaction1, signature1)

	var amount2 uint64 = 10 * blockchain.ParticlesCount
	transaction2 := blockchain.NewTransaction(0, walletAAddress, walletBAddress, amount2)
	signature2, _ := transaction2.Sign(walletB.PrivateKey())
	addTransaction(service, walletA.PublicKey(), transaction2, signature2)

	// Assert
	actualWalletAAmount := service.CalculateTotalAmount(0, walletAAddress)
	test.Assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	var expectedWalletBAmount uint64 = 0
	actualWalletBAmount := service.CalculateTotalAmount(0, walletBAddress)
	test.Assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func addTransaction(blockchain *blockchain.Service, senderWalletPublicKey *encryption.PublicKey, transaction *blockchain.Transaction, signature *encryption.Signature) *sync.WaitGroup {
	blockchain.AddTransaction(transaction, senderWalletPublicKey, signature)
	wg := blockchain.WaitGroup()
	wg.Wait()
	blockchain.Mine()
	wg.Wait()
	return wg
}
