package poh

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/my-cloud/ruthenium/src/log"
)

const (
	pohSmartContractAddressHex = "0xC5E9dDebb09Cd64DfaCab4011A0D5cEDaf7c9BDb"
	clientUrl                  = "https://mainnet.infura.io/v3/"
)

type Registry struct {
	infuraKey string
}

func NewRegistry(infuraKey string, logger log.Logger) *Registry {
	if infuraKey == "" {
		logger.Warn("infura key not provided")
	}
	return &Registry{infuraKey}
}

func (registry *Registry) IsRegistered(address string) (isRegistered bool, err error) {
	if registry.infuraKey == "" {
		return true, nil
	}
	client, err := ethclient.Dial(fmt.Sprint(clientUrl, registry.infuraKey))
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
