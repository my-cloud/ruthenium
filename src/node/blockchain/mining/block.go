package mining

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"ruthenium/src/node/neighborhood"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(previousHash [32]byte, transactions []*Transaction) (block *Block, err error) {
	var nonce int
	for {
		block = &Block{
			time.Now().UnixNano(),
			nonce,
			previousHash,
			transactions,
		}
		var pow *ProofOfWork
		if pow, err = block.ProofOfWork(); err != nil {
			err = fmt.Errorf("failed to get proof of work: %w", err)
			return
		}
		if pow.IsInValid() {
			nonce++
		} else {
			return
		}
	}
}

func NewBlockFromDto(block *neighborhood.BlockResponse) *Block {
	var transactions []*Transaction
	for _, transaction := range block.Transactions {
		transactions = append(transactions, NewTransactionFromResponse(transaction))
	}
	return &Block{
		block.Timestamp,
		block.Nonce,
		block.PreviousHash,
		transactions,
	}
}

func (block *Block) Hash() (hash [32]byte, err error) {
	marshaledBlock, err := json.Marshal(block)
	if err != nil {
		err = fmt.Errorf("failed to marshal block: %w", err)
		return
	}
	hash = sha256.Sum256(marshaledBlock)
	return
}

func (block *Block) ProofOfWork() (pow *ProofOfWork, err error) {
	hash, err := block.Hash()
	if err != nil {
		err = fmt.Errorf("failed to calculate block hash: %w", err)
		return
	}
	pow = NewProofOfWork(hash)
	return
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

func (block *Block) GetDto() *neighborhood.BlockResponse {
	var transactions []*neighborhood.TransactionResponse
	for _, transaction := range block.transactions {
		transactions = append(transactions, transaction.GetDto())
	}
	return &neighborhood.BlockResponse{
		Timestamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: block.previousHash,
		Transactions: transactions,
	}
}
