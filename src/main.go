package main

import (
	"fmt"
	"ruthenium/src/chain"
)

func main() {
	// Wallet
	walletA := chain.NewWallet()
	walletB := chain.NewWallet()
	walletC := chain.NewWallet()
	minerWallet := chain.NewWallet()

	// Block chain
	blockChain := chain.NewBlockchain(minerWallet.Address())
	isAdded1 := blockChain.AddTransaction(walletA, walletB.Address(), 100.)
	fmt.Println("Transaction 1 added? ", isAdded1)
	blockChain.Mining()

	isAdded2 := blockChain.AddTransaction(walletB, walletC.Address(), 10.)
	fmt.Println("Transaction 2 added? ", isAdded2)
	blockChain.Mining()

	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount(walletA.Address()))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount(walletB.Address()))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount(walletC.Address()))
	fmt.Printf("mine %.1f\n", blockChain.CalculateTotalAmount(minerWallet.Address()))
}
