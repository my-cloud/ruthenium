package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
)

type Signature struct {
	// Public key x coordinate
	r *big.Int

	// Can be computed by referring to information
	// like the transactions hash and the temporary public key
	// for generating signature
	s *big.Int
}

func NewSignature(transaction *Transaction, privateKey *ecdsa.PrivateKey) (*Signature, error) {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	hash := sha256.Sum256(marshaledTransaction)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return &Signature{r, s}, nil
}

func DecodeSignature(signatureString string) (*Signature, error) {
	r, s, err := string2BigIntTuple(signatureString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	return &Signature{&r, &s}, nil
}

func NewPublicKey(publicKeyString string) (*ecdsa.PublicKey, error) {
	x, y, err := string2BigIntTuple(publicKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}, nil
}

func NewPrivateKey(privateKeyString string, publicKey *ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(privateKeyString[:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{*publicKey, &bi}, nil
}

func (signature *Signature) String() string {
	return fmt.Sprintf("%064x%064x", signature.r, signature.s)
}

func string2BigIntTuple(s string) (big.Int, big.Int, error) {
	bx, err := hex.DecodeString(s[:64])
	if err != nil {
		return big.Int{}, big.Int{}, err
	}
	by, err := hex.DecodeString(s[64:])
	if err != nil {
		return big.Int{}, big.Int{}, err
	}

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy, nil
}
