package protocoltest

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
)

func NewSignedTransactionRequest(inputsValue uint64, fee uint64, recipientAddress string, utxoTransaction *network.TransactionResponse, utxoIndex uint16, senderPrivateKey *encryption.PrivateKey, senderPublicKey *encryption.PublicKey, timestamp int64, value uint64, blockHeight int) network.TransactionRequest {
	utxo := NewUtxoFromOutput(utxoTransaction, utxoIndex)
	marshalledInput, _ := json.Marshal(struct {
		OutputIndex   uint16 `json:"output_index"`
		TransactionId string `json:"transaction_id"`
	}{
		OutputIndex:   utxo.OutputIndex,
		TransactionId: utxo.TransactionId,
	})
	signature, _ := encryption.NewSignature(marshalledInput, senderPrivateKey)
	hexPublicKey := senderPublicKey.String()
	hexSignature := signature.String()
	input := network.InputRequest{
		OutputIndex:   &utxo.OutputIndex,
		TransactionId: &utxo.TransactionId,
		PublicKey:     &hexPublicKey,
		Signature:     &hexSignature,
	}
	var b bool
	sent := network.OutputRequest{
		Address:     &recipientAddress,
		BlockHeight: &blockHeight,
		HasReward:   &b,
		HasIncome:   &b,
		Value:       &value,
	}
	restValue := inputsValue - value - fee
	rest := network.OutputRequest{
		Address:     &recipientAddress,
		BlockHeight: &blockHeight,
		HasReward:   &b,
		HasIncome:   &b,
		Value:       &restValue,
	}
	broadcasterTarget := "0"
	return network.TransactionRequest{
		Inputs:                       &[]network.InputRequest{input},
		Outputs:                      &[]network.OutputRequest{sent, rest},
		Timestamp:                    &timestamp,
		TransactionBroadcasterTarget: &broadcasterTarget,
	}
}

func NewUtxoFromOutput(utxoTransaction *network.TransactionResponse, utxoIndex uint16) *network.UtxoResponse {
	outputResponse := utxoTransaction.Outputs[utxoIndex]
	return &network.UtxoResponse{
		Address:       outputResponse.Address,
		BlockHeight:   outputResponse.BlockHeight,
		HasReward:     outputResponse.HasReward,
		HasIncome:     outputResponse.HasIncome,
		OutputIndex:   utxoIndex,
		TransactionId: utxoTransaction.Id,
		Value:         outputResponse.Value,
	}
}

func NewTransactionRequest(address string, blockHeight int, value uint64, timestamp int64, target string) network.TransactionRequest {
	b := false
	output := network.OutputRequest{
		Address:     &address,
		BlockHeight: &blockHeight,
		HasReward:   &b,
		HasIncome:   &b,
		Value:       &value,
	}
	transactionRequest := network.TransactionRequest{
		Inputs:                       nil,
		Outputs:                      &[]network.OutputRequest{output},
		Timestamp:                    &timestamp,
		TransactionBroadcasterTarget: &target,
	}
	return transactionRequest
}
