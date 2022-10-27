package server

import (
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
)

type Amount struct {
	address string
}

func NewAmount(address string) *Amount {
	return &Amount{address}
}

func (amount *Amount) GetRequest() *neighborhood.AmountRequest {
	return &neighborhood.AmountRequest{Address: &amount.address}
}
