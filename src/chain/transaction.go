package chain

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	senderAddress    string
	senderPublicKey  *ecdsa.PublicKey
	recipientAddress string
	value            float32
}

func NewTransaction(senderAddress string, senderPublicKey *ecdsa.PublicKey, recipientAddress string, value float32) *Transaction {
	return &Transaction{senderAddress, senderPublicKey, recipientAddress, value}
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

func (transaction *Transaction) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &struct {
		Sender    *string  `json:"sender_address"`
		Recipient *string  `json:"recipient_address"`
		Value     *float32 `json:"value"`
	}{
		Sender:    &transaction.senderAddress,
		Recipient: &transaction.recipientAddress,
		Value:     &transaction.value,
	}); err != nil {
		return err
	}
	return nil
}

func (transaction *Transaction) Value() float32 {
	return transaction.value
}

func (transaction *Transaction) SenderAddress() string {
	return transaction.senderAddress
}

func (transaction *Transaction) SenderPublicKey() *ecdsa.PublicKey {
	return transaction.senderPublicKey
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address %s\n", transaction.senderAddress)
	fmt.Printf(" recipient_address %s\n", transaction.recipientAddress)
	fmt.Printf(" value %.1f\n", transaction.value)
}
