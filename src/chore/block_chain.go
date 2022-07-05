package chore

import (
	"fmt"
	"log"
	"strings"
)

const (
	MiningDifficulty = 3
	MiningSender     = "THE BLOCK CHAIN"
	MiningReward     = 1.
)

type BlockChain struct {
	transactions []*Transaction
	blocks       []*Block
	address      string
}

func NewBlockChain(address string) *BlockChain {
	blockchain := new(BlockChain)
	blockchain.address = address
	blockchain.createBlock(0, new(Block).Hash())
	return blockchain
}

func (blockChain *BlockChain) Print() {
	for i, block := range blockChain.blocks {
		fmt.Printf("%s Block  %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

func (blockChain *BlockChain) AddTransaction(sender string, recipient string, value float32) {
	transaction := NewTransaction(sender, recipient, value)
	blockChain.transactions = append(blockChain.transactions, transaction)
}

func (blockChain *BlockChain) Mining() bool {
	blockChain.AddTransaction(MiningSender, blockChain.address, MiningReward)
	nonce := blockChain.proofOfWork()
	lastHash := blockChain.lastBlock().Hash()
	blockChain.createBlock(nonce, lastHash)
	log.Println("action=mining, status=success")
	return true
}

func (blockChain *BlockChain) createBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, blockChain.transactions)
	blockChain.blocks = append(blockChain.blocks, block)
	blockChain.transactions = []*Transaction{}
	return block
}

func (blockChain *BlockChain) lastBlock() *Block {
	return blockChain.blocks[len(blockChain.blocks)-1]
}

func (blockChain *BlockChain) copyTransactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range blockChain.transactions {
		transactions = append(transactions,
			NewTransaction(transaction.sender,
				transaction.recipient,
				transaction.value))
	}
	return transactions
}

func (blockChain *BlockChain) validProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (blockChain *BlockChain) proofOfWork() int {
	transactions := blockChain.copyTransactions()
	lastHash := blockChain.lastBlock().Hash()
	var nonce int
	for !blockChain.validProof(nonce, lastHash, transactions, MiningDifficulty) {
		nonce++
	}
	return nonce
}

func (blockChain *BlockChain) CalculateTotalAmount(address string) float32 {
	var totalAmount float32 = 0.0
	for _, block := range blockChain.blocks {
		for _, transaction := range block.transactions {
			value := transaction.value
			if address == transaction.recipient {
				totalAmount += value
			}

			if address == transaction.sender {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}
