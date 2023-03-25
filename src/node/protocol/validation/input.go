package validation

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Input struct {
	outputIndex           uint16
	previousTransactionId [32]byte
	publicKey             *encryption.PublicKey
	signature             *encryption.Signature
}

func NewInput(outputIndex uint16, previousTransactionId [32]byte, publicKeyString string, signatureString string) (*Input, error) {
	publicKey, err := encryption.NewPublicKeyFromHex(publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	signature, err := encryption.DecodeSignature(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	return &Input{outputIndex, previousTransactionId, publicKey, signature}, nil
}

func (input *Input) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		OutputIndex           uint16 `json:"output_index"`
		PreviousTransactionId string `json:"previous_transaction_id"`
	}{
		OutputIndex:           input.outputIndex,
		PreviousTransactionId: fmt.Sprintf("%x", input.previousTransactionId),
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
		OutputIndex:           input.outputIndex,
		PreviousTransactionId: input.previousTransactionId,
		PublicKey:             encodedPublicKey,
		Signature:             encodedSignature,
	}
}

func (input *Input) VerifySignature() error {
	// FIXME implement
	//marshaledTransaction, err := input.MarshalJSON()
	//if err != nil {
	//	return fmt.Errorf("failed to marshal transaction, %w", err)
	//}
	//if !input.signature.Verify(marshaledTransaction, input.publicKey, input.address) {
	//	return errors.New("failed to verify signature")
	//}
	return nil
}
