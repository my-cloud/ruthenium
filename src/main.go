package main

import (
	"fmt"
	"ruthenium/src/chain"
	"ruthenium/src/wallet"
)

func main() {
	// Wallet
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	walletC := wallet.NewWallet()
	minerWallet := wallet.NewWallet()

	transaction1 := wallet.NewTransaction(walletA, walletB.Address(), 100.)

	// Block chain
	blockChain := chain.NewBlockchain(minerWallet.Address())
	isAdded := blockChain.AddTransaction(walletA.Address(), walletB.Address(), 100., walletA.PublicKey(), transaction1.GenerateSignature())
	fmt.Println("Added? ", isAdded)
	blockChain.Mining()

	transaction2 := wallet.NewTransaction(walletB, walletC.Address(), 10.)
	blockChain.AddTransaction(walletB.Address(), walletC.Address(), 10., walletB.PublicKey(), transaction2.GenerateSignature())
	blockChain.Mining()

	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount(walletA.Address()))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount(walletB.Address()))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount(walletC.Address()))
	fmt.Printf("mine %.1f\n", blockChain.CalculateTotalAmount(minerWallet.Address()))
	fmt.Printf("Signature: %s \n", transaction1.GenerateSignature())
}
