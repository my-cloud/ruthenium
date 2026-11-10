package ledger

import (
	"encoding/json"

	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
)

func NewSignedTransaction(inputsValue uint64, fee int, outputIndex uint16, recipientAddress string, privateKey *encryption.PrivateKey, publicKey *encryption.PublicKey, timestamp int64, transactionId string, value uint64, isYielding bool) *Transaction {
	sent := NewOutput(recipientAddress, false, value)
	restValue := uint64(int(inputsValue) - int(value) - fee)
	rest := NewOutput(recipientAddress, isYielding, restValue)
	return newSignedTransaction(outputIndex, privateKey, publicKey, timestamp, transactionId, []*Output{sent, rest})
}

func NewSignedSingleOutputTransaction(outputIndex uint16, recipientAddress string, privateKey *encryption.PrivateKey, publicKey *encryption.PublicKey, timestamp int64, transactionId string, value uint64) *Transaction {
	sent := NewOutput(recipientAddress, false, value)
	return newSignedTransaction(outputIndex, privateKey, publicKey, timestamp, transactionId, []*Output{sent})
}

func newSignedTransaction(outputIndex uint16, privateKey *encryption.PrivateKey, publicKey *encryption.PublicKey, timestamp int64, transactionId string, outputs []*Output) *Transaction {
	marshalledInput, _ := json.Marshal(struct {
		OutputIndex   uint16 `json:"output_index"`
		TransactionId string `json:"transaction_id"`
	}{
		OutputIndex:   outputIndex,
		TransactionId: transactionId,
	})
	signature, _ := encryption.NewSignature(marshalledInput, privateKey)
	signatureString := signature.String()
	input, _ := NewInput(outputIndex, transactionId, publicKey.String(), signatureString)
	inputs := []*Input{input}
	id, _ := generateId(inputs, outputs, timestamp)
	dto := &transactionDto{
		Id:        id,
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: timestamp,
	}
	marshalledTransaction, _ := json.Marshal(dto)
	var transaction *Transaction
	_ = json.Unmarshal(marshalledTransaction, &transaction)
	return transaction
}
