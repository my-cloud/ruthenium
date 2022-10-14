package blockchain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/my-cloud/ruthenium/src/poh"
)

const (
	networkId                  = "mainnet"
	infuraKey                  = "ac46e51cf15e45e0a4c00c35fa780f1b"
	pohSmartContractAddressHex = "0xC5E9dDebb09Cd64DfaCab4011A0D5cEDaf7c9BDb"
)

type Human struct {
	address common.Address
}

func NewHuman(hex string) *Human {
	return &Human{common.HexToAddress(hex)}
}

func (human *Human) IsRegistered() (isRegistered bool, err error) {
	clientUrl := fmt.Sprintf("https://%s.infura.io/v3/%s", networkId, infuraKey)
	client, err := ethclient.Dial(clientUrl)
	if err != nil {
		return
	}
	proofOfHumanity, err := poh.NewPoh(common.HexToAddress(pohSmartContractAddressHex), client)
	if err != nil {
		return
	}
	isRegistered, err = proofOfHumanity.PohCaller.IsRegistered(nil, human.address)
	if err != nil {
		return
	}
	return
}
