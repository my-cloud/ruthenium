package server

import (
	"github.com/my-cloud/ruthenium/src/api/node"
)

type Amount struct {
	address string
}

func NewAmount(address string) *Amount {
	return &Amount{address}
}

func (amount *Amount) GetRequest() *node.AmountRequest {
	return &node.AmountRequest{Address: &amount.address}
}
