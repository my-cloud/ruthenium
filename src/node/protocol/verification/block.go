package verification

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
)

type Block struct {
	timestamp                  int64
	previousHash               [32]byte
	transactions               []*validation.Transaction
	addedRegisteredAddresses   []string
	removedRegisteredAddresses []string
	registeredAddresses        []string
}

func NewBlockResponse(timestamp int64, previousHash [32]byte, transactions []*network.TransactionResponse, addedRegisteredAddresses []string, removedRegisteredAddresses []string) *network.BlockResponse {
	return &network.BlockResponse{
		Timestamp:                  timestamp,
		PreviousHash:               previousHash,
		Transactions:               transactions,
		AddedRegisteredAddresses:   addedRegisteredAddresses,
		RemovedRegisteredAddresses: removedRegisteredAddresses,
	}
}

func NewBlockFromResponse(block *network.BlockResponse, lastRegisteredAddresses []string) (*Block, error) {
	var transactions []*validation.Transaction
	for _, transactionResponse := range block.Transactions {
		transaction, err := validation.NewTransactionFromResponse(transactionResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	cleanedAddresses := lastRegisteredAddresses
	for _, address := range block.RemovedRegisteredAddresses {
		removeAddress(cleanedAddresses, address)
	}
	registeredAddressesMap := make(map[string]bool)
	for _, address := range append(lastRegisteredAddresses, block.AddedRegisteredAddresses...) {
		if _, ok := registeredAddressesMap[address]; !ok {
			registeredAddressesMap[address] = false
		}
	}
	var registeredAddresses []string
	for address := range registeredAddressesMap {
		registeredAddresses = append(registeredAddresses, address)
	}
	return &Block{
		block.Timestamp,
		block.PreviousHash,
		transactions,
		block.AddedRegisteredAddresses,
		block.RemovedRegisteredAddresses,
		registeredAddresses,
	}, nil
}

func (block *Block) AddedRegisteredAddresses() []string {
	return block.addedRegisteredAddresses
}

func (block *Block) GetResponse() *network.BlockResponse {
	var transactions []*network.TransactionResponse
	for _, transaction := range block.transactions {
		transactions = append(transactions, transaction.GetResponse())
	}
	return &network.BlockResponse{
		Timestamp:                  block.timestamp,
		PreviousHash:               block.previousHash,
		Transactions:               transactions,
		AddedRegisteredAddresses:   block.registeredAddresses,
		RemovedRegisteredAddresses: block.registeredAddresses,
	}
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

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp           int64                     `json:"timestamp"`
		PreviousHash        string                    `json:"previous_hash"`
		Transactions        []*validation.Transaction `json:"transactions"`
		RegisteredAddresses []string                  `json:"registered_addresses"`
	}{
		Timestamp:           block.timestamp,
		PreviousHash:        fmt.Sprintf("%x", block.previousHash),
		Transactions:        block.transactions,
		RegisteredAddresses: block.registeredAddresses,
	})
}

func (block *Block) PreviousHash() [32]byte {
	return block.previousHash
}

func (block *Block) RegisteredAddresses() []string {
	return block.registeredAddresses
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
