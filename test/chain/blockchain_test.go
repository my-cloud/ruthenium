package chain

import (
	"fmt"
	"path/filepath"
	"runtime"
	"ruthenium/src/chain"
	"testing"
)

func Test_Blockchain(t *testing.T) {
	// Wallet
	walletA := chain.NewWallet()
	walletB := chain.NewWallet()
	walletC := chain.NewWallet()
	minerWallet := chain.NewWallet()

	// Block chain
	blockChain := chain.NewBlockchain(minerWallet.Address(), 0)

	transaction1 := chain.NewTransaction(walletA, walletB.Address(), 100.)
	signature1 := chain.NewSignature(transaction1, walletA.PrivateKey())
	isAdded1 := blockChain.AddTransaction(transaction1, signature1)
	assert(t, isAdded1, "Failed to add first transaction")
	blockChain.Mine()

	transaction2 := chain.NewTransaction(walletB, walletC.Address(), 10.)
	signature2 := chain.NewSignature(transaction2, walletB.PrivateKey())
	isAdded2 := blockChain.AddTransaction(transaction2, signature2)
	assert(t, isAdded2, "Failed to add second transaction")
	blockChain.Mine()

	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount(walletA.Address()))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount(walletB.Address()))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount(walletC.Address()))
	fmt.Printf("mine %.1f\n", blockChain.CalculateTotalAmount(minerWallet.Address()))
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
