package humanity

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	networkId                  = "mainnet"
	infuraKey                  = "ac46e51cf15e45e0a4c00c35fa780f1b"
	pohSmartContractAddressHex = "0xC5E9dDebb09Cd64DfaCab4011A0D5cEDaf7c9BDb"
)

type Registry struct {
	clientUrl string
}

func NewRegistry() *Registry {
	clientUrl := fmt.Sprintf("https://%s.infura.io/v3/%s", networkId, infuraKey)
	return &Registry{clientUrl}
}

func (registry *Registry) IsRegistered(address string) (isRegistered bool, err error) {
	client, err := ethclient.Dial(registry.clientUrl)
	if err != nil {
		return
	}
	pohSmartContractAddress := common.HexToAddress(pohSmartContractAddressHex)
	proofOfHumanity, err := NewPoh(pohSmartContractAddress, client)
	if err != nil {
		return
	}
	isRegistered, err = proofOfHumanity.PohCaller.IsRegistered(nil, common.HexToAddress(address))
	if err != nil {
		return
	}
	return
}
