package chain

import (
	"fmt"
	"math"
	"path/filepath"
	"runtime"
	"ruthenium/src/log"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/blockchain"
	"ruthenium/src/node/blockchain/mine"
	"testing"
)

func Test_Blockchain(t *testing.T) {
	// Arrange
	walletA, _ := authentication.NewWallet("", "")
	walletB, _ := authentication.NewWallet("", "")
	minerWallet, _ := authentication.NewWallet("", "")

	// Act
	logger := log.NewLogger(log.Error)
	blockChain := blockchain.NewService(minerWallet.Address(), "", 8106, logger)
	wg := blockChain.WaitGroup()
	var value1 float32 = 40.
	for blockChain.CalculateTotalAmount(minerWallet.Address()) < value1 {
		blockChain.Mine()
		wg.Wait()
	}

	transaction1 := mine.NewTransaction(minerWallet.Address(), walletA.Address(), value1)
	signature1, _ := transaction1.Sign(minerWallet.PrivateKey())
	blockChain.AddTransaction(transaction1, minerWallet.PublicKey(), signature1)
	wg.Wait()
	blockChain.Mine()
	wg.Wait()

	var value2 float32 = 10.
	transaction2 := mine.NewTransaction(walletA.Address(), walletB.Address(), value2)
	signature2, _ := transaction2.Sign(walletA.PrivateKey())
	blockChain.AddTransaction(transaction2, walletA.PublicKey(), signature2)
	wg.Wait()
	blockChain.Mine()
	wg.Wait()

	// Assert
	reward := blockchain.MiningReward
	mineOperationsCount := float32(math.Ceil(float64(value1 / reward)))
	expectedMinerWalletAmount := mineOperationsCount*reward - value1 + 2*reward
	actualMinerWalletAmount := blockChain.CalculateTotalAmount(minerWallet.Address())
	assert(t, expectedMinerWalletAmount == actualMinerWalletAmount, fmt.Sprintf("Wrong miner wallet amount. Expected: %f - Actual: %f", expectedMinerWalletAmount, actualMinerWalletAmount))
	expectedWalletAAmount := value1 - value2
	actualWalletAAmount := blockChain.CalculateTotalAmount(walletA.Address())
	assert(t, expectedWalletAAmount == actualWalletAAmount, fmt.Sprintf("Wrong wallet A amount. Expected: %f - Actual: %f", expectedWalletAAmount, actualWalletAAmount))
	expectedWalletBAmount := value2
	actualWalletBAmount := blockChain.CalculateTotalAmount(walletB.Address())
	assert(t, expectedWalletBAmount == actualWalletBAmount, fmt.Sprintf("Wrong wallet B amount. Expected: %f - Actual: %f", expectedWalletBAmount, actualWalletBAmount))
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
