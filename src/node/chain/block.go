package chain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"ruthenium/src/log"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/neighborhood"
	"strings"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*authentication.Transaction
	logger       *log.Logger
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*authentication.Transaction, logger *log.Logger) *Block {
	return &Block{
		time.Now().UnixNano(),
		nonce,
		previousHash,
		transactions,
		logger,
	}
}

func NewBlockFromDto(block *neighborhood.BlockResponse, logger *log.Logger) *Block {
	var transactions []*authentication.Transaction
	for _, transaction := range block.Transactions {
		transactions = append(transactions, authentication.NewTransactionFromDto(transaction, logger))
	}
	return &Block{
		block.Timestamp,
		block.Nonce,
		block.PreviousHash,
		transactions,
		logger,
	}
}

func (block *Block) Hash() [32]byte {
	marshaledBlock, err := json.Marshal(block)
	if err != nil {
		block.logger.Error(fmt.Errorf("failed to marshal block: %w", err).Error())
		return [32]byte{}
	}
	return sha256.Sum256(marshaledBlock)
}

func (block *Block) PreviousHash() [32]byte {
	return block.previousHash
}

func (block *Block) Nonce() int {
	return block.nonce
}

func (block *Block) Transactions() []*authentication.Transaction {
	return block.transactions
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                         `json:"timestamp"`
		Nonce        int                           `json:"nonce"`
		PreviousHash string                        `json:"previous_hash"`
		Transactions []*authentication.Transaction `json:"transactions"`
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

func (block *Block) IsInValid(difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	hashStr := fmt.Sprintf("%x", block.Hash())
	return hashStr[:difficulty] != zeros
}
