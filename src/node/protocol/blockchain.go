package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Blockchain interface {
	AddBlock(timestamp int64, transactions []byte, newRegisteredAddresses []string) error
	Blocks(startingBlockHeight uint64) []byte
	Copy() Blockchain // TODO remove
	FirstBlockTimestamp() int64
	LastBlockTimestamp() int64
	Utxos(address string) []byte
	UtxosById() map[string][]*network.UtxoResponse
}
