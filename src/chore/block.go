package chore

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		nonce,
		previousHash,
		time.Now().UnixNano(),
		transactions,
	}
}

func (block *Block) Hash() [32]byte {
	serializedBlock, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return sha256.Sum256(serializedBlock)
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TimeStamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		TimeStamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: block.previousHash,
		Transactions: block.transactions,
	})
}

func (block *Block) Print() {
	fmt.Printf("timestamp       %d\n", block.timestamp)
	fmt.Printf("nonce           %d\n", block.nonce)
	fmt.Printf("previous_hash   %x\n", block.previousHash)
	for _, transaction := range block.transactions {
		transaction.Print()
	}
}
