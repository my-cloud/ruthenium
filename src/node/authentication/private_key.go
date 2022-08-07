package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

func NewPrivateKey(privateKeyString string, publicKey *PublicKey) (*PrivateKey, error) {
	b, err := hex.DecodeString(privateKeyString[:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	var bi big.Int
	_ = bi.SetBytes(b)
	return &PrivateKey{&ecdsa.PrivateKey{*publicKey.PublicKey, &bi}}, nil
}

func GeneratePrivateKey() (*PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{privateKey}, err
}

func (privateKey *PrivateKey) ExtractPublicKey() *PublicKey {
	return &PublicKey{&privateKey.PublicKey}
}

func (privateKey *PrivateKey) String() string {
	return fmt.Sprintf("%x", privateKey.D.Bytes())
}
