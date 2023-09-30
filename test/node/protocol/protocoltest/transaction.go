package protocoltest

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
)

func NewSignedTransactionRequest(inputsValue uint64, fee uint64, outputIndex uint16, recipientAddress string, senderPrivateKey *encryption.PrivateKey, senderPublicKey *encryption.PublicKey, timestamp int64, transactionId string, value uint64) network.TransactionRequest {
	marshalledInput, _ := json.Marshal(struct {
		OutputIndex   uint16 `json:"output_index"`
		TransactionId string `json:"transaction_id"`
	}{
		OutputIndex:   outputIndex,
		TransactionId: transactionId,
	})
	signature, _ := encryption.NewSignature(marshalledInput, senderPrivateKey)
	hexPublicKey := senderPublicKey.String()
	hexSignature := signature.String()
	input := network.InputRequest{
		OutputIndex:   &outputIndex,
		TransactionId: &transactionId,
		PublicKey:     &hexPublicKey,
		Signature:     &hexSignature,
	}
	var b bool
	sent := network.OutputRequest{
		Address:   &recipientAddress,
		HasReward: &b,
		HasIncome: &b,
		Value:     &value,
	}
	restValue := inputsValue - value - fee
	rest := network.OutputRequest{
		Address:   &recipientAddress,
		HasReward: &b,
		HasIncome: &b,
		Value:     &restValue,
	}
	broadcasterTarget := "0"
	return network.TransactionRequest{
		Inputs:                       &[]network.InputRequest{input},
		Outputs:                      &[]network.OutputRequest{sent, rest},
		Timestamp:                    &timestamp,
		TransactionBroadcasterTarget: &broadcasterTarget,
	}
}

func NewTransactionRequest(address string, value uint64, timestamp int64, target string) network.TransactionRequest {
	b := false
	output := network.OutputRequest{
		Address:   &address,
		HasReward: &b,
		HasIncome: &b,
		Value:     &value,
	}
	transactionRequest := network.TransactionRequest{
		Inputs:                       &[]network.InputRequest{},
		Outputs:                      &[]network.OutputRequest{output},
		Timestamp:                    &timestamp,
		TransactionBroadcasterTarget: &target,
	}
	return transactionRequest
}
