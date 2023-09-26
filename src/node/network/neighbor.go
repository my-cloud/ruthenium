package network

type Neighbor interface {
	Target() string
	GetFirstBlockTimestamp() (int64, error)
	GetBlocks(startingBlockHeight uint64) ([]byte, error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []byte, err error)
	GetUtxos(address string) (utxos []*UtxoResponse, err error)
}
