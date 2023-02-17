package protocoltest

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
)

func NewSignedTransactionRequest(fee uint64, recipientAddress string, senderAddress string, senderPrivateKey *encryption.PrivateKey, senderPublicKey *encryption.PublicKey, timestamp int64, value uint64) network.TransactionRequest {
	marshaledTransaction, _ := json.Marshal(struct {
		Fee              uint64 `json:"fee"`
		RecipientAddress string `json:"recipient_address"`
		SenderAddress    string `json:"sender_address"`
		Timestamp        int64  `json:"timestamp"`
		Value            uint64 `json:"value"`
	}{
		Fee:              fee,
		RecipientAddress: recipientAddress,
		SenderAddress:    senderAddress,
		Timestamp:        timestamp,
		Value:            value,
	})
	signature, _ := encryption.NewSignature(marshaledTransaction, senderPrivateKey)
	encodedPublicKey := senderPublicKey.String()
	encodedSignature := signature.String()
	return network.TransactionRequest{
		Fee:              &fee,
		RecipientAddress: &recipientAddress,
		SenderAddress:    &senderAddress,
		SenderPublicKey:  &encodedPublicKey,
		Signature:        &encodedSignature,
		Timestamp:        &timestamp,
		Value:            &value,
	}
}

func NewTransactionRequest(fee uint64, recipientAddress string, senderAddress string, timestamp int64, value uint64) network.TransactionRequest {
	return network.TransactionRequest{
		Fee:              &fee,
		RecipientAddress: &recipientAddress,
		SenderAddress:    &senderAddress,
		Timestamp:        &timestamp,
		Value:            &value,
	}
}
