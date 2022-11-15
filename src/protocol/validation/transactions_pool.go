package validation

import (
	"github.com/my-cloud/ruthenium/src/network"
)

type TransactionsPool interface {
	AddTransaction(transactionRequest *network.TransactionRequest, neighbors []network.Neighbor)
	Transactions() []*network.TransactionResponse
}
