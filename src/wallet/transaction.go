package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"log"
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
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		log.Println("ERROR: transaction marshal failed")
	}
	hash := sha256.Sum256(marshaledTransaction)
	r, s, err := ecdsa.Sign(rand.Reader, transaction.sender.PrivateKey(), hash[:])
	if err != nil {
		log.Println("ERROR: signature generation failed")
	}
	return &Signature{r, s}
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
