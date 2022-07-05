package main

import (
	"fmt"
	"ruthenium/src/chore"
)

func main() {
	myBlockChainAddress := "my block chain address"
	blockChain := chore.NewBlockChain(myBlockChainAddress)
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
}
