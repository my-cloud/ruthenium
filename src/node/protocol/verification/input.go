package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
)

type inputDto struct {
	OutputIndex   uint16 `json:"output_index"`
	TransactionId string `json:"transaction_id"`
	PublicKey     string `json:"public_key"`
	Signature     string `json:"signature"`
}

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
	return json.Marshal(inputDto{
		OutputIndex:   input.outputIndex,
		TransactionId: input.transactionId,
		PublicKey:     encodedPublicKey,
		Signature:     encodedSignature,
	})
}

func (input *Input) UnmarshalJSON(data []byte) error {
	var dto *inputDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	publicKey, err := encryption.NewPublicKeyFromHex(dto.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}
	signature, err := encryption.DecodeSignature(dto.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}
	input.outputIndex = dto.OutputIndex
	input.transactionId = dto.TransactionId
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
