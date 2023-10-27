package protocoltest

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
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
	input, _ := verification.NewInput(outputIndex, transactionId, publicKey.String(), signatureString)
	sent := verification.NewOutput(recipientAddress, false, value)
	restValue := inputsValue - value - fee
	rest := verification.NewOutput(recipientAddress, isYielding, restValue)
	inputs := []*verification.Input{input}
	outputs := []*verification.Output{sent, rest}
	id := generateId(inputs, outputs, timestamp)
	transaction := struct {
		Id        string
		Inputs    []*verification.Input
		Outputs   []*verification.Output
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
			Inputs    []*verification.Input
			Outputs   []*verification.Output
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

func generateId(inputs []*verification.Input, outputs []*verification.Output, timestamp int64) string {
	marshaledTransaction, _ := json.Marshal(struct {
		Inputs    []*verification.Input  `json:"inputs"`
		Outputs   []*verification.Output `json:"outputs"`
		Timestamp int64                  `json:"timestamp"`
	}{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: timestamp,
	})
	transactionHash := sha256.Sum256(marshaledTransaction)
	id := fmt.Sprintf("%x", transactionHash)
	return id
}
