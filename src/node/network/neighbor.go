package network

type Neighbor interface {
	Target() string
	GetBlocks(startingBlockHeight uint64) ([]byte, error)
	GetFirstBlockTimestamp() (int64, error)
	SendTargets(targets []string) (err error)
	AddTransaction(transaction []byte) (err error)
	GetTransactions() (transactions []byte, err error)
	GetUtxos(address string) (utxos []byte, err error)
}
