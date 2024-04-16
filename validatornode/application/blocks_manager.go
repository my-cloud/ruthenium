package application

import "github.com/my-cloud/ruthenium/validatornode/domain/protocol"

type BlocksManager interface {
	AddBlock(timestamp int64, transactions []*protocol.Transaction, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []*protocol.Block
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	LastBlockTransactions() []*protocol.Transaction
}
