package network

type Neighbor interface {
	Target() string
	GetBlocks() (blockResponses []*BlockResponse, err error)
	GetLastBlocks(lastBlocksRequest LastBlocksRequest) ([]*BlockResponse, error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []TransactionResponse, err error)
	GetAmount(request AmountRequest) (amountResponse *AmountResponse, err error)
}
