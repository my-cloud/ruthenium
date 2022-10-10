package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/poh"
)

const (
	networkId               = "mainnet"
	infuraKey               = "ac46e51cf15e45e0a4c00c35fa780f1b"
	pohSmartContractAddress = "0xC5E9dDebb09Cd64DfaCab4011A0D5cEDaf7c9BDb"
)

type Block struct {
	timestamp    int64
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(timestamp int64, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		timestamp,
		previousHash,
		transactions,
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

func (block *Block) IsProofOfHumanityValid() (err error) {
	minerAddress := block.minerAddress()
	clientUrl := fmt.Sprintf("https://%s.infura.io/v3/%s", networkId, infuraKey)
	client, err := ethclient.Dial(clientUrl)
	if err != nil {
		return err
	}
	proofOfHumanity, err := poh.NewPoh(common.HexToAddress(pohSmartContractAddress), client)
	if err != nil {
		return err
	}
	isRegistered, err := proofOfHumanity.PohCaller.IsRegistered(nil, common.HexToAddress(minerAddress))
	if err != nil {
		return err
	}
	if !isRegistered {
		return errors.New("not registered")
	}
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

func (block *Block) minerAddress() string {
	var minerAddress string
	for i := len(block.transactions) - 1; i >= 0; i-- {
		if block.transactions[i].SenderAddress() == RewardSenderAddress {
			minerAddress = block.transactions[i].RecipientAddress()
			break
		}
	}
	return minerAddress
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		PreviousHash: fmt.Sprintf("%x", block.previousHash),
		Transactions: block.transactions,
	})
}

func (block *Block) GetResponse() *neighborhood.BlockResponse {
	var transactions []*neighborhood.TransactionResponse
	for _, transaction := range block.transactions {
		transactions = append(transactions, transaction.GetResponse())
	}
	return &neighborhood.BlockResponse{
		Timestamp:    block.timestamp,
		PreviousHash: block.previousHash,
		Transactions: transactions,
	}
}
