package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
)

type Transaction struct {
	sender           *Wallet
	recipientAddress string
	value            float32
}

func NewTransaction(sender *Wallet, recipientAddress string, value float32) *Transaction {
	return &Transaction{sender, recipientAddress, value}
}

func (transaction *Transaction) GenerateSignature() *Signature {
	serializedTransaction, _ := json.Marshal(transaction)
	hash := sha256.Sum256(serializedTransaction)
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.sender.privateKey, hash[:])
	return &Signature{r, s}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.sender.address,
		Recipient: transaction.recipientAddress,
		Value:     transaction.value,
	})
}
