package ledger

import (
	"encoding/json"

	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
)

func NewSignedTransaction(inputsValue uint64, fee uint64, outputIndex uint16, recipientAddress string, privateKey *encryption.PrivateKey, publicKey *encryption.PublicKey, timestamp int64, transactionId string, value uint64, isYielding bool) *Transaction {
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
