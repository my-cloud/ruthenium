package chain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	sender    string
	recipient string
	value     float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (transaction *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender %s\n", transaction.sender)
	fmt.Printf(" recipient %s\n", transaction.recipient)
	fmt.Printf(" value %.1f\n", transaction.value)
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.sender,
		Recipient: transaction.recipient,
		Value:     transaction.value,
	})
}
