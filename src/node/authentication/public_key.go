package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
)

type PublicKey struct {
	*ecdsa.PublicKey
}

func NewPublicKey(publicKeyString string) (*PublicKey, error) {
	x, y, err := string2BigIntTuple(publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	return &PublicKey{&ecdsa.PublicKey{elliptic.P256(), &x, &y}}, nil
}

func (publicKey *PublicKey) String() string {
	return fmt.Sprintf("%064x%064x", publicKey.X.Bytes(), publicKey.Y.Bytes())
}
