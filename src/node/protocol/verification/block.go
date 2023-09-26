package verification

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
)

type Block struct {
	timestamp                  int64
	previousHash               [32]byte
	transactions               []*validation.Transaction
	addedRegisteredAddresses   []string
	removedRegisteredAddresses []string
}

func NewBlock(timestamp int64, previousHash [32]byte, transactions []*validation.Transaction, addedRegisteredAddresses []string, removedRegisteredAddresses []string) *Block {
	return &Block{timestamp, previousHash, transactions, addedRegisteredAddresses, removedRegisteredAddresses}
}

func (block *Block) AddedRegisteredAddresses() []string {
	return block.addedRegisteredAddresses
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

func (block *Block) UnmarshalJSON(data []byte) error {
	blockDto := struct {
		Timestamp                  int64                     `json:"timestamp"`
		PreviousHash               [32]byte                  `json:"previous_hash"`
		Transactions               []*validation.Transaction `json:"transactions"`
		AddedRegisteredAddresses   []string                  `json:"added_registered_addresses"`
		RemovedRegisteredAddresses []string                  `json:"removed_registered_addresses"`
	}{}
	err := json.Unmarshal(data, &blockDto)
	if err != nil {
		return err
	}
	block.timestamp = blockDto.Timestamp
	block.previousHash = blockDto.PreviousHash
	block.transactions = blockDto.Transactions
	block.addedRegisteredAddresses = blockDto.AddedRegisteredAddresses
	block.removedRegisteredAddresses = blockDto.RemovedRegisteredAddresses
	return nil
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp                  int64                     `json:"timestamp"`
		PreviousHash               [32]byte                  `json:"previous_hash"`
		Transactions               []*validation.Transaction `json:"transactions"`
		AddedRegisteredAddresses   []string                  `json:"added_registered_addresses"`
		RemovedRegisteredAddresses []string                  `json:"removed_registered_addresses"`
	}{
		Timestamp:                  block.timestamp,
		PreviousHash:               block.previousHash,
		Transactions:               block.transactions,
		AddedRegisteredAddresses:   block.addedRegisteredAddresses,
		RemovedRegisteredAddresses: block.removedRegisteredAddresses,
	})
}

func (block *Block) PreviousHash() [32]byte {
	return block.previousHash
}

func (block *Block) RemovedRegisteredAddresses() []string {
	return block.removedRegisteredAddresses
}

func (block *Block) Timestamp() int64 {
	return block.timestamp
}

func (block *Block) Transactions() []*validation.Transaction {
	return block.transactions
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
