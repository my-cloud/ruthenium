package validation

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

func NewInputFromResponse(input *network.InputResponse) (*Input, error) {
	publicKey, err := encryption.NewPublicKeyFromHex(input.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	signature, err := encryption.DecodeSignature(input.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	return &Input{input.OutputIndex, input.TransactionId, publicKey, signature}, nil
}

func (input *Input) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		OutputIndex   uint16 `json:"output_index"`
		TransactionId string `json:"transaction_id"`
	}{
		OutputIndex:   input.outputIndex,
		TransactionId: input.transactionId,
	})
}

func (input *Input) GetResponse() *network.InputResponse {
	var encodedPublicKey string
	if input.publicKey != nil {
		encodedPublicKey = input.publicKey.String()
	}
	var encodedSignature string
	if input.signature != nil {
		encodedSignature = input.signature.String()
	}
	return &network.InputResponse{
		OutputIndex:   input.outputIndex,
		TransactionId: input.transactionId,
		PublicKey:     encodedPublicKey,
		Signature:     encodedSignature,
	}
}

func (input *Input) VerifySignature() error {
	marshaledInput, err := input.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal input, %w", err)
	}
	if !input.signature.Verify(marshaledInput, input.publicKey) {
		return errors.New("signature is invalid")
	}
	return nil
}
