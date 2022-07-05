package main

import (
	"ruthenium/src/chore"
)

func main() {
	blockChain := chore.NewBlockChain()
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.)

	previousHash := blockChain.LastBlock().Hash()
	blockChain.CreateBlock(5, previousHash)
	blockChain.Print()

	previousHash = blockChain.LastBlock().Hash()
	blockChain.CreateBlock(2, previousHash)
	blockChain.Print()
}
