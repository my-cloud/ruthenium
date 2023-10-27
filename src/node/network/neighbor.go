package network

type Neighbor interface {
	Target() string
	GetBlocks(startingBlockHeight uint64) (blocks []byte, err error)
	GetFirstBlockTimestamp() (firstBlockTimestamp int64, err error)
	GetSettings() (settings []byte, err error)
	SendTargets(targets []string) error
	AddTransaction(transaction []byte) error
	GetTransactions() (transactions []byte, err error)
	GetUtxos(address string) (utxos []byte, err error)
}
