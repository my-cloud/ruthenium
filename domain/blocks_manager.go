package domain

type BlocksManager interface {
	AddBlock(timestamp int64, transactionsBytes []byte, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []byte
	Copy() BlocksManager // TODO remove
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	Utxos(address string) []byte
	Utxo(input InputInfoProvider) (UtxoInfoProvider, error)
}
