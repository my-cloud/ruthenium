package blockchain

import (
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/log"
	"gitlab.com/coinsmaster/ruthenium/src/node/authentication"
	"gitlab.com/coinsmaster/ruthenium/src/node/blockchain"
	"gitlab.com/coinsmaster/ruthenium/src/node/blockchain/mining"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Test_AddTransaction_Allowed_TransactionAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := authentication.NewWallet("", "")
	minerWalletAddress := minerWallet.Address()
	walletA, _ := authentication.NewWallet("", "")
	walletAAddress := walletA.Address()
	walletB, _ := authentication.NewWallet("", "")
	walletBAddress := walletB.Address()
	logger := log.NewLogger(log.Error)
	service := blockchain.NewService(minerWalletAddress, "", 0, time.Nanosecond, logger)

	// Act
	var amount1 uint64 = 40 * blockchain.ParticlesCount
	transaction1 := mining.NewTransaction(0, minerWalletAddress, walletAAddress, amount1)
	signature1, _ := transaction1.Sign(minerWallet.PrivateKey())
	addTransaction(service, minerWallet.PublicKey(), transaction1, signature1)

	var amount2 uint64 = 10 * blockchain.ParticlesCount
	transaction2 := mining.NewTransaction(0, walletAAddress, walletBAddress, amount2)
	signature2, _ := transaction2.Sign(walletA.PrivateKey())
	addTransaction(service, walletA.PublicKey(), transaction2, signature2)

	// Assert
	expectedWalletAAmount := amount1 - amount2
	actualWalletAAmount := service.CalculateTotalAmount(0, walletAAddress)
	assert(t, expectedWalletAAmount == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", expectedWalletAAmount, actualWalletAAmount))
	expectedWalletBAmount := amount2
	actualWalletBAmount := service.CalculateTotalAmount(0, walletBAddress)
	assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func Test_AddTransaction_Allowed_TransactionNotAdded(t *testing.T) {
	// Arrange
	minerWallet, _ := authentication.NewWallet("", "")
	minerWalletAddress := minerWallet.Address()
	walletA, _ := authentication.NewWallet("", "")
	walletAAddress := walletA.Address()
	walletB, _ := authentication.NewWallet("", "")
	walletBAddress := walletB.Address()
	logger := log.NewLogger(log.Error)
	service := blockchain.NewService(minerWalletAddress, "", 0, time.Nanosecond, logger)

	// Act
	var amount1 uint64 = 40 * blockchain.ParticlesCount
	transaction1 := mining.NewTransaction(0, minerWalletAddress, walletAAddress, amount1)
	signature1, _ := transaction1.Sign(minerWallet.PrivateKey())
	addTransaction(service, minerWallet.PublicKey(), transaction1, signature1)

	var amount2 uint64 = 10 * blockchain.ParticlesCount
	transaction2 := mining.NewTransaction(0, walletAAddress, walletBAddress, amount2)
	signature2, _ := transaction2.Sign(walletB.PrivateKey())
	addTransaction(service, walletA.PublicKey(), transaction2, signature2)

	// Assert
	actualWalletAAmount := service.CalculateTotalAmount(0, walletAAddress)
	assert(t, amount1 == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", amount1, actualWalletAAmount))
	var expectedWalletBAmount uint64 = 0
	actualWalletBAmount := service.CalculateTotalAmount(0, walletBAddress)
	assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func addTransaction(blockchain *blockchain.Service, senderWalletPublicKey *authentication.PublicKey, transaction *mining.Transaction, signature *authentication.Signature) *sync.WaitGroup {
	blockchain.AddTransaction(transaction, senderWalletPublicKey, signature)
	wg := blockchain.WaitGroup()
	wg.Wait()
	blockchain.Mine()
	wg.Wait()
	return wg
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
