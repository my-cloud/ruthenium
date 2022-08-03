package chain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

func NewWallet(privateKey *ecdsa.PrivateKey) *Wallet {
	// 1. Create ECDSA public key (64 bytes)
	publicKey := &privateKey.PublicKey
	address := CreateAddress(publicKey)
	return &Wallet{privateKey, publicKey, address}
}

func CreateAddress(publicKey *ecdsa.PublicKey) string {
	// 2. Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(publicKey.X.Bytes())
	h2.Write(publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes).
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network).
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	checksum := digest6[:4]
	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], checksum[:])
	// 9. Convert the result from a byte string into base58.
	address := base58.Encode(dc8)
	return address
}

func (wallet *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PublicKey  string `json:"public_key"`
		PrivateKey string `json:"private_key"`
		Address    string `json:"address"`
	}{
		PublicKey:  wallet.publicKeyStr(),
		PrivateKey: wallet.privateKeyStr(),
		Address:    wallet.Address(),
	})
}

func (wallet *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) PublicKey() *ecdsa.PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) Address() string {
	return wallet.address
}

func (wallet *Wallet) privateKeyStr() string {
	return fmt.Sprintf("%x", wallet.privateKey.D.Bytes())
}

func (wallet *Wallet) publicKeyStr() string {
	return fmt.Sprintf("%064x%064x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}
