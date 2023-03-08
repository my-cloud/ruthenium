package encryption

import (
	"encoding/json"
)

type Wallet struct {
	publicKey *PublicKey
	address   string
}

func NewWallet(publicKey *PublicKey) *Wallet {
	address := publicKey.Address()
	return &Wallet{publicKey, address}
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

func (wallet *Wallet) Address() string {
	return wallet.address
}
