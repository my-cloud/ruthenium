package server

import "github.com/my-cloud/ruthenium/src/api/node/network"

type Amount struct {
	address string
}

func NewAmount(address string) *Amount {
	return &Amount{address}
}

func (amount *Amount) GetRequest() *network.AmountRequest {
	return &network.AmountRequest{Address: &amount.address}
}
