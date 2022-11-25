package wallet

import (
	network2 "github.com/my-cloud/ruthenium/src/node/network"
)

type Amount struct {
	address string
}

func NewAmount(address string) *Amount {
	return &Amount{address}
}

func (amount *Amount) GetRequest() *network2.AmountRequest {
	return &network2.AmountRequest{Address: &amount.address}
}
