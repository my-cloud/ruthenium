package chain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"log"
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

func NewTransactionFromDto(transaction *TransactionResponse) *Transaction {
	return &Transaction{transaction.SenderAddress, nil, transaction.RecipientAddress, transaction.Value}
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

func (transaction *Transaction) SenderPublicKey() *ecdsa.PublicKey {
	return transaction.senderPublicKey
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) VerifySignature(signature *Signature) bool {
	marshaledBlockchain, err := json.Marshal(transaction)
	if err != nil {
		log.Println("ERROR: Failed to marshal blockchain")
	}
	hash := sha256.Sum256(marshaledBlockchain)
	return ecdsa.Verify(transaction.SenderPublicKey(), hash[:], signature.r, signature.s)
}

func (transaction *Transaction) GetDto() *TransactionResponse {
	return &TransactionResponse{
		SenderAddress:    transaction.senderAddress,
		RecipientAddress: transaction.recipientAddress,
		Value:            transaction.value,
	}
}
