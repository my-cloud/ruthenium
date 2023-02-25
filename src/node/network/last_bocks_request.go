package network

type LastBlocksRequest struct {
	StartingBlockHeight *uint64
}

func (lastBlocksRequest LastBlocksRequest) IsInvalid() bool {
	return lastBlocksRequest.StartingBlockHeight == nil
}
