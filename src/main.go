package main

import (
	"fmt"
	"ruthenium/src/chain"
	"ruthenium/src/wallet"
)

func main() {
	// Block chain
	myBlockChainAddress := "my block chain address"
	blockChain := chain.NewBlockChain(myBlockChainAddress)
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 100.)
	blockChain.Mining()
	blockChain.AddTransaction("B", "C", 10.)
	blockChain.Mining()
	blockChain.Print()
	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount("A"))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount("B"))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount("C"))
	fmt.Printf("mine %.1f\n", blockChain.CalculateTotalAmount("my block chain address"))

	// Wallet
	w := wallet.NewWallet()
	fmt.Println("Private key: ", w.PrivateKeyStr())
	fmt.Println("Public key: ", w.PublicKeyStr())
	fmt.Println("Address: ", w.Address())
}
