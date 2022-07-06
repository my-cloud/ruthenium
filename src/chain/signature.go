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
		log.Println("ERROR: transaction marshal failed")
	}
	hash := sha256.Sum256(marshaledTransaction)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Println("ERROR: signature generation failed")
	}
	return &Signature{r, s}
}

func NewPublicKey(publicKeyString string) *ecdsa.PublicKey {
	x, y := string2BigIntTuple(publicKeyString)
	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}
}

func NewPrivateKey(privateKeyString string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(privateKeyString[:])
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{*publicKey, &bi}
}

func (signature *Signature) String() string {
	return fmt.Sprintf("%064x%064x", signature.r, signature.s)
}

func string2BigIntTuple(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}
