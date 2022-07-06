package chain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	sender           *Wallet
	recipientAddress string
	value            float32
}

func NewTransaction(sender *Wallet, recipientAddress string, value float32) *Transaction {
	return &Transaction{sender, recipientAddress, value}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.sender.Address(),
		Recipient: transaction.recipientAddress,
		Value:     transaction.value,
	})
}

func (transaction *Transaction) Value() float32 {
	return transaction.value
}

func (transaction *Transaction) Sender() *Wallet {
	return transaction.sender
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address %s\n", transaction.sender.Address())
	fmt.Printf(" recipient_address %s\n", transaction.recipientAddress)
	fmt.Printf(" value %.1f\n", transaction.value)
}
