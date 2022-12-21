package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type TransactionsPool interface {
	AddTransaction(transactionRequest *network.TransactionRequest, neighbors []network.Neighbor)
	Transactions() []*network.TransactionResponse
}
