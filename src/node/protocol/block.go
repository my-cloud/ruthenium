package protocol

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
)

type Block struct {
	timestamp           int64
	previousHash        [32]byte
	transactions        []*Transaction
	registeredAddresses []string
}

func NewBlock(timestamp int64, previousHash [32]byte, transactions []*Transaction, registeredAddresses []string) *Block {
	return &Block{
		timestamp,
		previousHash,
		transactions,
		registeredAddresses,
	}
}

func NewBlockFromResponse(block *neighborhood.BlockResponse) (*Block, error) {
	var transactions []*Transaction
	for _, transactionResponse := range block.Transactions {
		transaction, err := NewTransactionFromResponse(transactionResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return &Block{
		block.Timestamp,
		block.PreviousHash,
		transactions,
		block.RegisteredAddresses,
	}, nil
}

func (block *Block) Hash() (hash [32]byte, err error) {
	marshaledBlock, err := block.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("failed to marshal block: %w", err)
		return
	}
	hash = sha256.Sum256(marshaledBlock)
	return
}

func (block *Block) Timestamp() int64 {
	return block.timestamp
}

func (block *Block) PreviousHash() [32]byte {
	return block.previousHash
}

func (block *Block) Transactions() []*Transaction {
	return block.transactions
}

func (block *Block) RegisteredAddresses() []string {
	return block.registeredAddresses
}

func (block *Block) ValidatorAddress() string {
	var validatorAddress string
	for i := len(block.transactions) - 1; i >= 0; i-- {
		if block.transactions[i].IsReward() {
			validatorAddress = block.transactions[i].RecipientAddress()
			break
		}
	}
	return validatorAddress
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp           int64          `json:"timestamp"`
		PreviousHash        string         `json:"previous_hash"`
		Transactions        []*Transaction `json:"transactions"`
		RegisteredAddresses []string       `json:"registered_addresses"`
	}{
		Timestamp:           block.timestamp,
		PreviousHash:        fmt.Sprintf("%x", block.previousHash),
		Transactions:        block.transactions,
		RegisteredAddresses: block.registeredAddresses,
	})
}

func (block *Block) GetResponse() *neighborhood.BlockResponse {
	var transactions []*neighborhood.TransactionResponse
	for _, transaction := range block.transactions {
		transactions = append(transactions, transaction.GetResponse())
	}
	return &neighborhood.BlockResponse{
		Timestamp:           block.timestamp,
		PreviousHash:        block.previousHash,
		Transactions:        transactions,
		RegisteredAddresses: block.registeredAddresses,
	}
}
