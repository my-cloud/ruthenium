package blockchain

import (
	"fmt"
	"path/filepath"
	"runtime"
	"ruthenium/src/log"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/blockchain"
	"ruthenium/src/node/blockchain/mining"
	"testing"
	"time"
)

func Test_Blockchain(t *testing.T) {
	// Arrange
	minerWallet, _ := authentication.NewWallet("", "")
	minerWalletAddress := minerWallet.Address()
	walletA, _ := authentication.NewWallet("", "")
	walletAAddress := walletA.Address()
	walletB, _ := authentication.NewWallet("", "")
	walletBAddress := walletB.Address()
	logger := log.NewLogger(log.Error)
	blockChain := blockchain.NewService(minerWalletAddress, "", 8106, time.Nanosecond, logger)
	wg := blockChain.WaitGroup()

	// Act
	var value1 uint64 = 40 * blockchain.ParticlesCount
	transaction1 := mining.NewTransaction(0, minerWalletAddress, walletAAddress, value1)
	signature1, _ := transaction1.Sign(minerWallet.PrivateKey())
	blockChain.AddTransaction(transaction1, minerWallet.PublicKey(), signature1)
	wg.Wait()
	blockChain.Mine()
	wg.Wait()

	var value2 uint64 = 10 * blockchain.ParticlesCount
	transaction2 := mining.NewTransaction(0, walletAAddress, walletBAddress, value2)
	signature2, _ := transaction2.Sign(walletA.PrivateKey())
	blockChain.AddTransaction(transaction2, walletA.PublicKey(), signature2)
	wg.Wait()
	blockChain.Mine()
	wg.Wait()

	// Assert
	expectedWalletAAmount := value1 - value2
	actualWalletAAmount := blockChain.CalculateTotalAmount(0, walletAAddress)
	assert(t, expectedWalletAAmount == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %d - Actual: %d", expectedWalletAAmount, actualWalletAAmount))
	expectedWalletBAmount := value2
	actualWalletBAmount := blockChain.CalculateTotalAmount(0, walletBAddress)
	assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %d - Actual: %d", expectedWalletBAmount, actualWalletBAmount))
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
