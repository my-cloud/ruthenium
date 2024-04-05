package protocol

import (
	"encoding/json"

	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
)

func NewSignedTransactionRequest(inputsValue uint64, fee uint64, outputIndex uint16, recipientAddress string, privateKey *encryption.PrivateKey, publicKey *encryption.PublicKey, timestamp int64, transactionId string, value uint64, isYielding bool) []byte {
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
	sent := NewOutput(recipientAddress, false, value)
	restValue := inputsValue - value - fee
	rest := NewOutput(recipientAddress, isYielding, restValue)
	inputs := []*Input{input}
	outputs := []*Output{sent, rest}
	id, _ := generateId(inputs, outputs, timestamp)
	transaction := struct {
		Id        string
		Inputs    []*Input
		Outputs   []*Output
		Timestamp int64
	}{
		Id:        id,
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: timestamp,
	}
	transactionRequest := struct {
		Transaction struct {
			Id        string
			Inputs    []*Input
			Outputs   []*Output
			Timestamp int64
		}
		TransactionBroadcasterTarget string
	}{
		Transaction:                  transaction,
		TransactionBroadcasterTarget: "0",
	}
	marshalledTransactionRequest, _ := json.Marshal(transactionRequest)
	return marshalledTransactionRequest
}
