package protocol

import (
	network2 "github.com/my-cloud/ruthenium/src/node/network"
)

type TransactionsPool interface {
	AddTransaction(transactionRequest *network2.TransactionRequest, neighbors []network2.Neighbor)
	Transactions() []*network2.TransactionResponse
}
