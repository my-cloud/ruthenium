package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"runtime"
	"ruthenium/src"
	"ruthenium/src/chain"
	"testing"
)

func Test_Blockchain(t *testing.T) {
	// Wallet
	privateKey1, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	walletA := chain.NewWallet(privateKey1)
	privateKey2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	walletB := chain.NewWallet(privateKey2)
	privateKey3, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	walletC := chain.NewWallet(privateKey3)
	minerPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	minerWallet := chain.NewWallet(minerPrivateKey)

	// Block chain
	blockChain := src.NewBlockchain(minerWallet.Address(), 0)

	var value1 float32 = 100.
	transaction1 := chain.NewTransaction(walletA.Address(), walletA.PublicKey(), walletB.Address(), value1)
	signature1 := chain.NewSignature(transaction1, walletA.PrivateKey())
	isAdded1 := blockChain.UpdateTransaction(walletA.Address(), walletB.Address(), walletA.PublicKey(), value1, signature1)
	assert(t, isAdded1, "Failed to add first transaction")
	blockChain.Mine()

	var value2 float32 = 10.
	transaction2 := chain.NewTransaction(walletB.Address(), walletB.PublicKey(), walletC.Address(), value2)
	signature2 := chain.NewSignature(transaction2, walletB.PrivateKey())
	isAdded2 := blockChain.UpdateTransaction(walletB.Address(), walletC.Address(), walletB.PublicKey(), value2, signature2)
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
