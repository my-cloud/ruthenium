package domain

type Blockchain interface {
	AddBlock(timestamp int64, transactionsBytes []byte, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []byte
	Copy() Blockchain // TODO remove
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	Utxos(address string) []byte
	Utxo(input InputInfo) (Utxo, error)
}
