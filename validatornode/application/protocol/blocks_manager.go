package protocol

type BlocksManager interface {
	AddBlock(timestamp int64, transactionsBytes []byte, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []byte
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	LastBlockTransactions() []byte
}
