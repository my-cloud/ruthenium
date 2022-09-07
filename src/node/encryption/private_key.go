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

func NewPrivateKey() (*PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return &PrivateKey{privateKey}, err
}

func DecodePrivateKey(privateKeyString string) (*PrivateKey, error) {
	bytes, _ := hexutil.Decode(privateKeyString)
	ecdsaPrivateKey, err := crypto.ToECDSA(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	return &PrivateKey{ecdsaPrivateKey}, nil
}

func (privateKey *PrivateKey) String() string {
	privateKeyBytes := crypto.FromECDSA(privateKey.PrivateKey)
	return hexutil.Encode(privateKeyBytes)
}
