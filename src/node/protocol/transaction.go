package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/api/node"
	"github.com/my-cloud/ruthenium/src/node/encryption"
)

const (
	transactionFee      = 1000
	rewardSenderAddress = "REWARD SENDER ADDRESS"
)

type Transaction struct {
	recipientAddress string
	senderAddress    string
	senderPublicKey  *encryption.PublicKey
	signature        *encryption.Signature
	timestamp        int64
	value            uint64
	fee              uint64
}

func NewRewardTransaction(recipientAddress string, timestamp int64, value uint64) *Transaction {
	return &Transaction{
		recipientAddress: recipientAddress,
		senderAddress:    rewardSenderAddress,
		senderPublicKey:  nil,
		timestamp:        timestamp,
		value:            value,
		fee:              transactionFee,
	}
}

func NewTransaction(recipientAddress string, senderAddress string, senderPublicKey *encryption.PublicKey, timestamp int64, value uint64) *Transaction {
	return &Transaction{
		recipientAddress: recipientAddress,
		senderAddress:    senderAddress,
		senderPublicKey:  senderPublicKey,
		timestamp:        timestamp,
		value:            value,
		fee:              transactionFee,
	}
}

func NewTransactionFromRequest(transactionRequest *node.TransactionRequest) (*Transaction, error) {
	senderPublicKey, err := encryption.DecodePublicKey(*transactionRequest.SenderPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction public key: %w", err)
	}
	signature, err := encryption.DecodeSignature(*transactionRequest.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction signature: %w", err)
	}
	return &Transaction{
		*transactionRequest.RecipientAddress,
		*transactionRequest.SenderAddress,
		senderPublicKey,
		signature,
		*transactionRequest.Timestamp,
		*transactionRequest.Value,
		*transactionRequest.Fee,
	}, nil
}

func NewTransactionFromResponse(transactionResponse *node.TransactionResponse) (transaction *Transaction, err error) {
	var senderPublicKey *encryption.PublicKey
	if len(transactionResponse.SenderPublicKey) != 0 {
		senderPublicKey, err = encryption.DecodePublicKey(transactionResponse.SenderPublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transaction public key: %w", err)
		}
	}
	var signature *encryption.Signature
	if len(transactionResponse.SenderPublicKey) != 0 {
		signature, err = encryption.DecodeSignature(transactionResponse.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transaction signature: %w", err)
		}
	}
	return &Transaction{
		transactionResponse.RecipientAddress,
		transactionResponse.SenderAddress,
		senderPublicKey,
		signature,
		transactionResponse.Timestamp,
		transactionResponse.Value,
		transactionResponse.Fee,
	}, nil
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RecipientAddress string `json:"recipient_address"`
		SenderAddress    string `json:"sender_address"`
		Timestamp        int64  `json:"timestamp"`
		Value            uint64 `json:"value"`
		Fee              uint64 `json:"fee"`
	}{
		RecipientAddress: transaction.recipientAddress,
		SenderAddress:    transaction.senderAddress,
		Timestamp:        transaction.timestamp,
		Value:            transaction.value,
		Fee:              transaction.fee,
	})
}

func (transaction *Transaction) Timestamp() int64 {
	return transaction.timestamp
}

func (transaction *Transaction) SenderAddress() string {
	return transaction.senderAddress
}

func (transaction *Transaction) RecipientAddress() string {
	return transaction.recipientAddress
}

func (transaction *Transaction) Value() uint64 {
	return transaction.value
}

func (transaction *Transaction) Fee() uint64 {
	return transaction.fee
}

func (transaction *Transaction) Sign(privateKey *encryption.PrivateKey) (err error) {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}
	transaction.signature, err = encryption.NewSignature(marshaledTransaction, privateKey)
	return
}

func (transaction *Transaction) GetRequest() node.TransactionRequest {
	encodedPublicKey := transaction.senderPublicKey.String()
	encodedSignature := transaction.signature.String()
	return node.TransactionRequest{
		RecipientAddress: &transaction.recipientAddress,
		SenderAddress:    &transaction.senderAddress,
		SenderPublicKey:  &encodedPublicKey,
		Signature:        &encodedSignature,
		Timestamp:        &transaction.timestamp,
		Value:            &transaction.value,
		Fee:              &transaction.fee,
	}
}

func (transaction *Transaction) GetResponse() *node.TransactionResponse {
	var encodedPublicKey string
	if transaction.senderPublicKey != nil {
		encodedPublicKey = transaction.senderPublicKey.String()
	}
	var encodedSignature string
	if transaction.signature != nil {
		encodedSignature = transaction.signature.String()
	}
	return &node.TransactionResponse{
		RecipientAddress: transaction.recipientAddress,
		SenderAddress:    transaction.senderAddress,
		SenderPublicKey:  encodedPublicKey,
		Signature:        encodedSignature,
		Timestamp:        transaction.timestamp,
		Value:            transaction.value,
		Fee:              transaction.fee,
	}
}

func (transaction *Transaction) Equals(other *Transaction) bool {
	return transaction.recipientAddress == other.recipientAddress &&
		transaction.senderAddress == other.senderAddress &&
		transaction.timestamp == other.timestamp &&
		transaction.value == other.value &&
		transaction.fee == other.fee
}

func (transaction *Transaction) VerifySignature() error {
	marshaledTransaction, err := transaction.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal transaction, %w", err)
	}
	if !transaction.signature.Verify(marshaledTransaction, transaction.senderPublicKey, transaction.senderAddress) {
		return errors.New("failed to verify signature")
	}
	return nil
}

func (transaction *Transaction) IsReward() bool {
	return transaction.SenderAddress() == rewardSenderAddress
}
