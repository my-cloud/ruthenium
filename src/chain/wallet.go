package chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
)

type Wallet struct {
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
	address    string
}

func NewWallet(publicKeyString string, privateKeyString string) (*Wallet, error) {
	var publicKey *ecdsa.PublicKey
	var privateKey *ecdsa.PrivateKey
	var address string
	var err error
	if publicKeyString != "" && privateKeyString != "" {
		publicKey, err = NewPublicKey(publicKeyString)
		if err != nil {
			return nil, fmt.Errorf("failed to decode wallet public key: %w", err)
		}
		privateKey, err = NewPrivateKey(privateKeyString, publicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transaction private key: %w", err)
		}
		address = NewAddress(publicKey)
	} else if privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	} else {
		publicKey = &privateKey.PublicKey
		address = NewAddress(publicKey)
	}
	return &Wallet{publicKey, privateKey, address}, nil
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

func (wallet *Wallet) PublicKey() *ecdsa.PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
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
