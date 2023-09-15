package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Blockchain interface {
	AddBlock(timestamp int64, transactions []*network.TransactionResponse, newRegisteredAddresses []string) error
	AllBlocks() []*network.BlockResponse
	Block(blockHeight uint64) *network.BlockResponse
	Blocks(startingBlockHeight uint64) []*network.BlockResponse
	Copy() Blockchain
	FindFee(transaction *network.TransactionResponse, timestamp int64) (uint64, error)
	UtxosByAddress(address string) []*network.UtxoResponse
}
