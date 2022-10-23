package protocol

import "github.com/my-cloud/ruthenium/src/api/node/network"

type Verifiable interface {
	Verify(neighbors []network.Requestable)
	StartVerification(validatable Validatable, synchronizable network.Synchronizable)
	Blocks() []*network.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	AddBlock(blockResponse *network.BlockResponse)
}
