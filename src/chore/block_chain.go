package chore

import (
	"fmt"
	"strings"
)

type BlockChain struct {
	transactionPool []string
	blocks          []*Block
}

func NewBlockChain() *BlockChain {
	blockchain := new(BlockChain)
	blockchain.CreateBlock(0, new(Block).Hash())
	return blockchain
}

func (blockChain *BlockChain) LastBlock() *Block {
	return blockChain.blocks[len(blockChain.blocks)-1]
}

func (blockChain *BlockChain) Print() {
	for i, block := range blockChain.blocks {
		fmt.Printf("%s Block  %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

func (blockChain *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash)
	blockChain.blocks = append(blockChain.blocks, block)
	return block
}
