package encryption

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type PublicKey struct {
	*ecdsa.PublicKey
}

func NewPublicKey(privateKey *PrivateKey) *PublicKey {
	return &PublicKey{privateKey.Public().(*ecdsa.PublicKey)}
}

func NewPublicKeyFromHex(publicKeyString string) (*PublicKey, error) {
	bytes, err := hexutil.Decode(publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	ecdsaPublicKey, err := crypto.UnmarshalPubkey(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %w", err)
	}
	return &PublicKey{ecdsaPublicKey}, nil
}

func (publicKey *PublicKey) Address() string {
	return crypto.PubkeyToAddress(*publicKey.PublicKey).Hex()
}

func (publicKey *PublicKey) String() string {
	publicKeyBytes := crypto.FromECDSAPub(publicKey.PublicKey)
	return hexutil.Encode(publicKeyBytes)
}
