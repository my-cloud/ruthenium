package network

import "github.com/my-cloud/ruthenium/src/api/node"

type Requestable interface {
	Ip() string
	Port() uint16
	Target() string
	GetBlocks() (blockResponses []*node.BlockResponse, err error)
	SendTargets(request []node.TargetRequest) (err error)
	AddTransaction(request node.TransactionRequest) (err error)
	GetTransactions() (transactionResponses []node.TransactionResponse, err error)
	GetAmount(request node.AmountRequest) (amountResponse *node.AmountResponse, err error)
	Mine() (err error)
	StartMining() (err error)
	StopMining() error
}
