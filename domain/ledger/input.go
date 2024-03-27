package ledger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/domain/encryption"
)

type inputDto struct {
	OutputIndex   uint16 `json:"output_index"`
	TransactionId string `json:"transaction_id"`
	PublicKey     string `json:"public_key"`
	Signature     string `json:"signature"`
}

type Input struct {
	*InputInfo
	publicKey *encryption.PublicKey
	signature *encryption.Signature
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
	return &Input{NewInputInfo(outputIndex, transactionId), publicKey, signature}, nil
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
		OutputIndex:   input.InputInfo.OutputIndex(),
		TransactionId: input.InputInfo.TransactionId(),
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
	input.InputInfo = NewInputInfo(dto.OutputIndex, dto.TransactionId)
	input.publicKey = publicKey
	input.signature = signature
	return nil
}

func (input *Input) VerifySignature() error {
	marshaledInputInfo, err := json.Marshal(input.InputInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal input, %w", err)
	}
	if !input.signature.Verify(marshaledInputInfo, input.publicKey) {
		return errors.New("signature is invalid")
	}
	return nil
}
