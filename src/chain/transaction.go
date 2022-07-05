package chain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

func NewTransaction(senderAddress string, recipientAddress string, value float32) *Transaction {
	return &Transaction{senderAddress, recipientAddress, value}
}

func (transaction *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address %s\n", transaction.senderAddress)
	fmt.Printf(" recipient_address %s\n", transaction.recipientAddress)
	fmt.Printf(" value %.1f\n", transaction.value)
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
