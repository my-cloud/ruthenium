package network

type LastBlocksRequest struct {
	StartingBlockIndex *int64
}

func (lastBlocksRequest LastBlocksRequest) IsInvalid() bool {
	return lastBlocksRequest.StartingBlockIndex == nil
}
