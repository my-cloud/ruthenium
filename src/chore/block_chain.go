package chore

import (
	"fmt"
	"strings"
)

const MiningDifficulty = 3

type BlockChain struct {
	transactions []*Transaction
	blocks       []*Block
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
	block := NewBlock(nonce, previousHash, blockChain.transactions)
	blockChain.blocks = append(blockChain.blocks, block)
	blockChain.transactions = []*Transaction{}
	return block
}

func (blockChain *BlockChain) AddTransaction(sender string, recipient string, value float32) {
	transaction := NewTransaction(sender, recipient, value)
	blockChain.transactions = append(blockChain.transactions, transaction)
}

func (blockChain *BlockChain) CopyTransactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range blockChain.transactions {
		transactions = append(transactions,
			NewTransaction(transaction.sender,
				transaction.recipient,
				transaction.value))
	}
	return transactions
}

func (blockChain *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (blockChain *BlockChain) ProofOfWork() int {
	transactions := blockChain.CopyTransactions()
	lastHash := blockChain.LastBlock().Hash()
	var nonce int
	for !blockChain.ValidProof(nonce, lastHash, transactions, MiningDifficulty) {
		nonce++
	}
	return nonce
}
