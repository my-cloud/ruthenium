package server

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Transaction struct {
	fee              uint64
	recipientAddress string
	senderAddress    string
	senderPublicKey  *encryption.PublicKey
	signature        *encryption.Signature
	timestamp        int64
	value            uint64
}

func NewTransaction(fee uint64, recipientAddress string, senderAddress string, senderPublicKey *encryption.PublicKey, timestamp int64, value uint64) *Transaction {
	return &Transaction{
		fee:              fee,
		recipientAddress: recipientAddress,
		senderAddress:    senderAddress,
		senderPublicKey:  senderPublicKey,
		timestamp:        timestamp,
		value:            value,
	}
}

func (transaction *Transaction) Sign(privateKey *encryption.PrivateKey) (err error) {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}
	transaction.signature, err = encryption.NewSignature(marshaledTransaction, privateKey)
	return
}

func (transaction *Transaction) GetRequest() network.TransactionRequest {
	encodedPublicKey := transaction.senderPublicKey.String()
	encodedSignature := transaction.signature.String()
	return network.TransactionRequest{
		Fee:              &transaction.fee,
		RecipientAddress: &transaction.recipientAddress,
		SenderAddress:    &transaction.senderAddress,
		SenderPublicKey:  &encodedPublicKey,
		Signature:        &encodedSignature,
		Timestamp:        &transaction.timestamp,
		Value:            &transaction.value,
	}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Fee              uint64 `json:"fee"`
		RecipientAddress string `json:"recipient_address"`
		SenderAddress    string `json:"sender_address"`
		Timestamp        int64  `json:"timestamp"`
		Value            uint64 `json:"value"`
	}{
		Fee:              transaction.fee,
		RecipientAddress: transaction.recipientAddress,
		SenderAddress:    transaction.senderAddress,
		Timestamp:        transaction.timestamp,
		Value:            transaction.value,
	})
}
