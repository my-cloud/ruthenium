package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Input struct {
	outputIndex   uint16
	transactionId string
	publicKey     *encryption.PublicKey
	signature     *encryption.Signature
}

func NewInput(outputIndex uint16, transactionId string, publicKeyString string, signatureString string) (*Input, error) {
	publicKey, err := encryption.NewPublicKeyFromHex(publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	signature, err := encryption.DecodeSignature(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	return &Input{outputIndex, transactionId, publicKey, signature}, nil
}

func (input *Input) MarshalJSON() ([]byte, error) {
	var encodedPublicKey string
	if input.publicKey != nil {
		encodedPublicKey = input.publicKey.String()
	}
	var encodedSignature string
	if input.signature != nil {
		encodedSignature = input.signature.String()
	}
	return json.Marshal(struct {
		OutputIndex   uint16 `json:"output_index"`
		TransactionId string `json:"transaction_id"`
		PublicKey     string `json:"public_key"`
		Signature     string `json:"signature"`
	}{
		OutputIndex:   input.outputIndex,
		TransactionId: input.transactionId,
		PublicKey:     encodedPublicKey,
		Signature:     encodedSignature,
	})
}

func (input *Input) UnmarshalJSON(data []byte) error {
	var inputDto network.InputResponse
	err := json.Unmarshal(data, &inputDto)
	if err != nil {
		return err
	}
	input.outputIndex = inputDto.OutputIndex
	input.transactionId = inputDto.TransactionId
	publicKey, err := encryption.NewPublicKeyFromHex(inputDto.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}
	signature, err := encryption.DecodeSignature(inputDto.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}
	input.publicKey = publicKey
	input.signature = signature
	return nil
}

func (input *Input) OutputIndex() uint16 {
	return input.outputIndex
}

func (input *Input) TransactionId() string {
	return input.transactionId
}

func (input *Input) VerifySignature() error {
	marshaledInput, err := input.marshalJSONWithoutKeyAndSignature()
	if err != nil {
		return fmt.Errorf("failed to marshal input, %w", err)
	}
	if !input.signature.Verify(marshaledInput, input.publicKey) {
		return errors.New("signature is invalid")
	}
	return nil
}

func (input *Input) marshalJSONWithoutKeyAndSignature() ([]byte, error) {
	return json.Marshal(struct {
		OutputIndex   uint16 `json:"output_index"`
		TransactionId string `json:"transaction_id"`
	}{
		OutputIndex:   input.outputIndex,
		TransactionId: input.transactionId,
	})
}
