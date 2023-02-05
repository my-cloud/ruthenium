package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type TransactionsPool interface {
	AddTransaction(transactionRequest *network.TransactionRequest, hostTarget string)
	Transactions() []*network.TransactionResponse
}
