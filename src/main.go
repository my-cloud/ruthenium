package main

import (
	"ruthenium/src/chore"
)

func main() {
	blockChain := chore.NewBlockChain()
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.)

	lastHash := blockChain.LastBlock().Hash()
	nonce := blockChain.ProofOfWork()
	blockChain.CreateBlock(nonce, lastHash)
	blockChain.Print()

	lastHash = blockChain.LastBlock().Hash()
	nonce = blockChain.ProofOfWork()
	blockChain.CreateBlock(nonce, lastHash)
	blockChain.Print()
}
