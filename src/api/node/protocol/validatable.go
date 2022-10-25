package protocol

import (
	"github.com/my-cloud/ruthenium/src/api/node"
	"github.com/my-cloud/ruthenium/src/api/node/network"
)

type Validatable interface {
	Validate(timestamp int64, verifiable Verifiable, address string)
	AddTransaction(transactionRequest *node.TransactionRequest, blockchain Verifiable, neighbors []network.Requestable)
	Transactions() []*node.TransactionResponse
	Clear()
}
