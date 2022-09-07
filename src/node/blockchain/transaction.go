package blockchain

import (
	"encoding/json"
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/node/encryption"
	"gitlab.com/coinsmaster/ruthenium/src/node/neighborhood"
)

type Transaction struct {
	timestamp        int64
	senderAddress    string
	recipientAddress string
	value            uint64
}

func NewTransaction(timestamp int64, senderAddress string, recipientAddress string, value uint64) *Transaction {
	return &Transaction{
		timestamp,
		senderAddress,
		recipientAddress,
		value,
	}
}

func NewTransactionFromRequest(transaction *neighborhood.TransactionRequest) *Transaction {
	return &Transaction{
		*transaction.Timestamp,
		*transaction.SenderAddress,
		*transaction.RecipientAddress,
		*transaction.Value,
	}
}

func NewTransactionFromResponse(transaction *neighborhood.TransactionResponse) *Transaction {
	return &Transaction{
		transaction.Timestamp,
		transaction.SenderAddress,
		transaction.RecipientAddress,
		transaction.Value,
	}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp int64  `json:"timestamp"`
		Sender    string `json:"sender_address"`
		Recipient string `json:"recipient_address"`
		Value     uint64 `json:"value"`
	}{
		Timestamp: transaction.timestamp,
		Sender:    transaction.senderAddress,
		Recipient: transaction.recipientAddress,
		Value:     transaction.value,
	})
}

func (transaction *Transaction) Timestamp() int64 {
	return transaction.timestamp
}

func (transaction *Transaction) SenderAddress() string {
	return transaction.senderAddress
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) Value() uint64 {
	return transaction.value
}

func (transaction *Transaction) Sign(privateKey *encryption.PrivateKey) (signature *encryption.Signature, err error) {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	return encryption.NewSignature(marshaledTransaction, privateKey)
}

func (transaction *Transaction) GetResponse() *neighborhood.TransactionResponse {
	return &neighborhood.TransactionResponse{
		Timestamp:        transaction.timestamp,
		SenderAddress:    transaction.senderAddress,
		RecipientAddress: transaction.recipientAddress,
		Value:            transaction.value,
	}
}

func (transaction *Transaction) Equals(other *Transaction) bool {
	return transaction.timestamp == other.timestamp &&
		transaction.senderAddress == other.senderAddress &&
		transaction.recipientAddress == other.recipientAddress &&
		transaction.value == other.value
}