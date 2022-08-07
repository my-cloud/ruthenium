package authentication

import (
	"encoding/json"
	"fmt"
)

type Wallet struct {
	publicKey  *PublicKey
	privateKey *PrivateKey
	address    string
}

func NewWallet(publicKeyString string, privateKeyString string) (*Wallet, error) {
	var publicKey *PublicKey
	var privateKey *PrivateKey
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
	} else if privateKey, err = GeneratePrivateKey(); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	} else {
		publicKey = privateKey.ExtractPublicKey()
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
		PublicKey:  wallet.publicKey.String(),
		PrivateKey: wallet.privateKey.String(),
		Address:    wallet.Address(),
	})
}

func (wallet *Wallet) PublicKey() *PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) PrivateKey() *PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) Address() string {
	return wallet.address
}
