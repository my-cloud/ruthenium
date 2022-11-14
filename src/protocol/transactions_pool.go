package protocol

import (
	"github.com/my-cloud/ruthenium/src/network"
)

type TransactionsPool interface {
	AddTransaction(transactionRequest *network.TransactionRequest, blockchain Blockchain, neighbors []network.Neighbor)
	Transactions() []*network.TransactionResponse
}
