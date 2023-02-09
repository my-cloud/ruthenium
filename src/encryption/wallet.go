package encryption

import (
	"encoding/json"
	"fmt"
)

type Wallet struct {
	publicKey        *PublicKey
	address          string
	privateKeyString string
}

func NewEmptyWallet() *Wallet {
	return &Wallet{nil, "", ""}
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
	return &Wallet{publicKey, address, privateKeyString}, nil
}

func (wallet *Wallet) MarshalJSON() ([]byte, error) {
	var publicKey string
	if wallet.publicKey != nil {
		publicKey = wallet.publicKey.String()
	}
	return json.Marshal(struct {
		PublicKey string `json:"public_key"`
		Address   string `json:"address"`
	}{
		PublicKey: publicKey,
		Address:   wallet.address,
	})
}

func (wallet *Wallet) PrivateKeyString() string {
	return wallet.privateKeyString
}

func (wallet *Wallet) Address() string {
	return wallet.address
}
