package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/node/neighborhood"
	"io/ioutil"
	"net/http"
	"strings"
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

func NewBlockFromResponse(block *neighborhood.BlockResponse) *Block {
	var transactions []*Transaction
	for _, transaction := range block.Transactions {
		transactions = append(transactions, NewTransactionFromResponse(transaction))
	}
	return &Block{
		block.Timestamp,
		block.PreviousHash,
		transactions,
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

func (block *Block) IsProofOfHumanityValid() (err error) {
	proofOfHumanity := block.minerAddress()
	resp, err := http.Get("https://api.poh.dev/profiles/" + proofOfHumanity)
	if err != nil {
		err = fmt.Errorf("failed to get proof of humanity: %w", err)
		return
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			// TODO extract this code or log it properly
			fmt.Println(fmt.Errorf("failed to close proof of humanity request body: %w", bodyCloseError).Error())
		}
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read proof of humanity response: %w", err)
		return
	}
	if !strings.Contains(string(body), "\"registered\":true") {
		err = fmt.Errorf("the miner is currently not registered as a real human: %w", err)
		return
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
	return block.lastTransaction().RecipientAddress()
}

func (block *Block) lastTransaction() *Transaction {
	return block.transactions[len(block.transactions)-1]
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
