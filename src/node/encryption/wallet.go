package encryption

import (
	"encoding/json"
	"fmt"
)

type Wallet struct {
	privateKey *PrivateKey
	publicKey  *PublicKey
	address    string
}

func NewEmptyWallet() *Wallet {
	return &Wallet{nil, nil, ""}
}

func DecodeWallet(mnemonicString string, derivationPath string, password string, privateKeyString string) (*Wallet, error) {
	var privateKey *PrivateKey
	var publicKey *PublicKey
	var address string
	var err error
	if mnemonicString != "" {
		mnemonic := NewMnemonic(mnemonicString)
		privateKey, err = mnemonic.PrivateKey(derivationPath, password)
	} else if privateKeyString != "" {
		privateKey, err = DecodePrivateKey(privateKeyString)
	} else {
		return NewEmptyWallet(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}
	publicKey = NewPublicKey(privateKey)
	address = publicKey.Address()
	return &Wallet{privateKey, publicKey, address}, nil
}

func (wallet *Wallet) MarshalJSON() ([]byte, error) {
	var privateKey string
	if wallet.privateKey != nil {
		privateKey = wallet.privateKey.String()
	}
	var publicKey string
	if wallet.publicKey != nil {
		publicKey = wallet.publicKey.String()
	}
	return json.Marshal(struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Address    string `json:"address"`
	}{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    wallet.address,
	})
}

func (wallet *Wallet) PrivateKey() *PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) PublicKey() *PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) Address() string {
	return wallet.address
}
