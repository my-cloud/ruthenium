package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		time.Now().UnixNano(),
		nonce,
		previousHash,
		transactions,
	}
}

func (block *Block) Print() {
	fmt.Printf("timestamp       %d\n", block.timestamp)
	fmt.Printf("nonce           %d\n", block.nonce)
	fmt.Printf("previous_hash   %x\n", block.previousHash)
	for _, transaction := range block.transactions {
		transaction.Print()
	}
}

func (block *Block) Hash() [32]byte {
	marshaledBlock, err := json.Marshal(block)
	if err != nil {
		log.Println("ERROR: Failed to marshal block")
	}
	return sha256.Sum256(marshaledBlock)
}

func (block *Block) PreviousHash() [32]byte {
	return block.previousHash
}

func (block *Block) Nonce() int {
	return block.nonce
}

func (block *Block) Transactions() []*Transaction {
	return block.transactions
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: fmt.Sprintf("%x", block.previousHash),
		Transactions: block.transactions,
	})
}

func (block *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	blockDto := struct {
		Timestamp    *int64          `json:"timestamp"`
		Nonce        *int            `json:"nonce"`
		PreviousHash *string         `json:"previous_hash"`
		Transactions *[]*Transaction `json:"transactions"`
	}{
		Timestamp:    &block.timestamp,
		Nonce:        &block.nonce,
		PreviousHash: &previousHash,
		Transactions: &block.transactions,
	}
	if err := json.Unmarshal(data, &blockDto); err != nil {
		return err
	}
	pr, err := hex.DecodeString(*blockDto.PreviousHash)
	if err != nil {
		log.Println("ERROR: Failed to decode previous hash")
	}
	copy(block.previousHash[:], pr[:32])
	return nil
}
