package mine

import (
	"encoding/json"
	"fmt"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/neighborhood"
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

func NewTransaction(senderAddress string, recipientAddress string, value float32) *Transaction {
	return &Transaction{
		senderAddress,
		recipientAddress,
		value,
	}
}

func NewTransactionFromDto(transaction *neighborhood.TransactionResponse) *Transaction {
	return &Transaction{
		transaction.SenderAddress,
		transaction.RecipientAddress,
		transaction.Value,
	}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.senderAddress,
		Recipient: transaction.recipientAddress,
		Value:     transaction.value,
	})
}

func (transaction *Transaction) Value() float32 {
	return transaction.value
}

func (transaction *Transaction) SenderAddress() string {
	return transaction.senderAddress
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) Sign(privateKey *authentication.PrivateKey) (signature *authentication.Signature, err error) {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	return authentication.NewSignature(marshaledTransaction, privateKey)
}

func (transaction *Transaction) GetDto() *neighborhood.TransactionResponse {
	return &neighborhood.TransactionResponse{
		SenderAddress:    transaction.senderAddress,
		RecipientAddress: transaction.recipientAddress,
		Value:            transaction.value,
	}
}
