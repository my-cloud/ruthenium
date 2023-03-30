package network

type Neighbor interface {
	Target() string
	GetBlock(blockHeight uint64) (blockResponse *BlockResponse, err error)
	GetBlocks() (blockResponses []*BlockResponse, err error)
	GetLambda() (float64, error)
	GetLastBlocks(startingBlockHeight uint64) ([]*BlockResponse, error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []TransactionResponse, err error)
	GetUtxos(address string) (utxos []*OutputResponse, err error)
}
