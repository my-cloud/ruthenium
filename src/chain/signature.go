package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
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

func NewSignature(transaction *Transaction, privateKey *ecdsa.PrivateKey) *Signature {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		log.Println("ERROR: Failed to marshal transaction")
	}
	hash := sha256.Sum256(marshaledTransaction)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Println("ERROR: Failed to generate signature")
	}
	return &Signature{r, s}
}

func DecodeSignature(signatureString string) *Signature {
	r, s, err := string2BigIntTuple(signatureString)
	if err != nil {
		log.Println("ERROR: Failed to decode signature")
	}
	return &Signature{&r, &s}
}

func NewPublicKey(publicKeyString string) *ecdsa.PublicKey {
	x, y, err := string2BigIntTuple(publicKeyString)
	if err != nil {
		log.Println("ERROR: Failed to decode public key")
	}
	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}
}

func NewPrivateKey(privateKeyString string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, err := hex.DecodeString(privateKeyString[:])
	if err != nil {
		log.Println("ERROR: Failed to decode private key")
	}
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{*publicKey, &bi}
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
