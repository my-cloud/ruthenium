package ledger

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type blockDto struct {
	PreviousHash               [32]byte       `json:"previous_hash"`
	AddedRegisteredAddresses   []string       `json:"added_registered_addresses"`
	RemovedRegisteredAddresses []string       `json:"removed_registered_addresses"`
	Timestamp                  int64          `json:"timestamp"`
	Transactions               []*Transaction `json:"transactions"`
}

type Block struct {
	previousHash               [32]byte
	addedRegisteredAddresses   []string
	removedRegisteredAddresses []string
	timestamp                  int64
	transactions               []*Transaction
}

func NewBlock(previousHash [32]byte, addedRegisteredAddresses []string, removedRegisteredAddresses []string, timestamp int64, transactions []*Transaction) *Block {
	return &Block{previousHash, addedRegisteredAddresses, removedRegisteredAddresses, timestamp, transactions}
}

func (block *Block) UnmarshalJSON(data []byte) error {
	var dto *blockDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	block.previousHash = dto.PreviousHash
	block.addedRegisteredAddresses = dto.AddedRegisteredAddresses
	block.removedRegisteredAddresses = dto.RemovedRegisteredAddresses
	block.timestamp = dto.Timestamp
	block.transactions = dto.Transactions
	return nil
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(blockDto{
		PreviousHash:               block.previousHash,
		AddedRegisteredAddresses:   block.addedRegisteredAddresses,
		RemovedRegisteredAddresses: block.removedRegisteredAddresses,
		Timestamp:                  block.timestamp,
		Transactions:               block.transactions,
	})
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

func (block *Block) ValidatorAddress() string {
	var validatorAddress string
	for i := len(block.transactions) - 1; i >= 0; i-- {
		if block.transactions[i].HasReward() {
			validatorAddress = block.transactions[i].RewardRecipientAddress()
			break
		}
	}
	return validatorAddress
}

func (block *Block) PreviousHash() [32]byte {
	return block.previousHash
}

func (block *Block) AddedRegisteredAddresses() []string {
	return block.addedRegisteredAddresses
}

func (block *Block) RemovedRegisteredAddresses() []string {
	return block.removedRegisteredAddresses
}

func (block *Block) Timestamp() int64 {
	return block.timestamp
}

func (block *Block) Transactions() []*Transaction {
	return block.transactions
}
