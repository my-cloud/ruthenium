package network

import (
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
)

type TransactionsPool interface {
	AddTransaction(transactionRequest *neighborhood.TransactionRequest, blockchain Blockchain, neighbors []neighborhood.Neighbor)
	Transactions() []*neighborhood.TransactionResponse
}
