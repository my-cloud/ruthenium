package protocol

import "github.com/my-cloud/ruthenium/src/api/node/network"

type Validatable interface {
	Validate(timestamp int64, verifiable Verifiable, address string)
	AddTransaction(transactionRequest *network.TransactionRequest, blockchain Verifiable, neighbors []network.Requestable)
	Transactions() []*network.TransactionResponse
	Clear()
}
