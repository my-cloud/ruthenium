package network

type LastBlocksRequest struct {
	StartingBlockHeight *int64
}

func (lastBlocksRequest LastBlocksRequest) IsInvalid() bool {
	return lastBlocksRequest.StartingBlockHeight == nil
}
