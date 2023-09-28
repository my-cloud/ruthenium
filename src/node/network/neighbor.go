package network

type Neighbor interface {
	Target() string
	GetBlocks(startingBlockHeight uint64) ([]byte, error)
	GetFirstBlockTimestamp() (int64, error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []byte, err error)
	GetUtxos(address string) (utxos []*UtxoResponse, err error)
}
