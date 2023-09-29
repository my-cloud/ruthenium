package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Blockchain interface {
	AddBlock(timestamp int64, transactions []byte, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []byte
	Copy() Blockchain // TODO remove
	FindFee(inputs []*network.InputResponse, outputs []*network.OutputResponse, timestamp int64) (uint64, error)
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	Utxos(address string) []byte
}
