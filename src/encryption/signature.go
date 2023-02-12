package encryption

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func NewSignature(marshaledTransaction []byte, privateKey *PrivateKey) (*Signature, error) {
	hash := sha256.Sum256(marshaledTransaction)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey.PrivateKey, hash[:])
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

func (signature *Signature) String() string {
	return fmt.Sprintf("%064x%064x", signature.r, signature.s)
}

func (signature *Signature) Verify(marshaledTransaction []byte, publicKey *PublicKey, transactionSenderAddress string) bool {
	hash := sha256.Sum256(marshaledTransaction)
	isSignatureValid := ecdsa.Verify(publicKey.PublicKey, hash[:], signature.r, signature.s)
	var isTransactionValid bool
	if isSignatureValid {
		publicKeyAddress := publicKey.Address()
		isTransactionValid = transactionSenderAddress == publicKeyAddress
	}
	return isTransactionValid
}

func string2BigIntTuple(s string) (big.Int, big.Int, error) {
	if len(s) != 128 {
		return big.Int{}, big.Int{}, errors.New("signature length is invalid")
	}
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
