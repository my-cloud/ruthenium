package network

type Neighbor interface {
	Target() string
	GetBlocks() (blockResponses []*BlockResponse, err error)
	GetLastBlocks(startingBlockHeight uint64) ([]*BlockResponse, error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []TransactionResponse, err error)
	GetAmount(address string) (amount uint64, err error)
}
