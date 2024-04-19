package application

import "github.com/my-cloud/ruthenium/validatornode/domain/ledger"

type BlocksManager interface {
	AddBlock(timestamp int64, transactions []*ledger.Transaction, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []*ledger.Block
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	LastBlockTransactions() []*ledger.Transaction
}
