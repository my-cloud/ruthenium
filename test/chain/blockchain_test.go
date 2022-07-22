package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"runtime"
	"ruthenium/src/chain"
	"testing"
)

func Test_Blockchain(t *testing.T) {
	// Wallet
	privateKey1, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	walletA := chain.NewWallet(privateKey1)
	privateKey2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	walletB := chain.NewWallet(privateKey2)
	minerPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	minerWallet := chain.NewWallet(minerPrivateKey)

	// Blockchain
	blockChain := chain.NewBlockchain(minerWallet.Address(), "", 5000)
	var value1 float32 = 100.
	for blockChain.CalculateTotalAmount(minerWallet.Address()) < value1 {
		blockChain.Mine()
	}

	transaction1 := chain.NewTransaction(minerWallet.Address(), minerWallet.PublicKey(), walletA.Address(), value1)
	signature1 := chain.NewSignature(transaction1, minerWallet.PrivateKey())
	isAdded1 := blockChain.UpdateTransaction(minerWallet.Address(), walletA.Address(), minerWallet.PublicKey(), value1, signature1)
	assert(t, isAdded1, "Failed to add first transaction")
	blockChain.Mine()

	var value2 float32 = 10.
	transaction2 := chain.NewTransaction(walletA.Address(), walletA.PublicKey(), walletB.Address(), value2)
	signature2 := chain.NewSignature(transaction2, walletA.PrivateKey())
	isAdded2 := blockChain.UpdateTransaction(walletA.Address(), walletB.Address(), walletA.PublicKey(), value2, signature2)
	assert(t, isAdded2, "Failed to add second transaction")
	blockChain.Mine()

	fmt.Printf("mine %.1f\n", blockChain.CalculateTotalAmount(minerWallet.Address()))
	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount(walletA.Address()))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount(walletB.Address()))
}

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
