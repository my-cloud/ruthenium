package encryption

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

func DecodePrivateKey(privateKeyString string) (*PrivateKey, error) {
	bytes, err := hexutil.Decode(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	ecdsaPrivateKey, err := crypto.ToECDSA(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key bytes to ECDSA: %w", err)
	}
	return &PrivateKey{ecdsaPrivateKey}, nil
}

func (privateKey *PrivateKey) String() string {
	privateKeyBytes := crypto.FromECDSA(privateKey.PrivateKey)
	return hexutil.Encode(privateKeyBytes)
}
