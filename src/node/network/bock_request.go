package network

type BlockRequest struct {
	BlockHeight *uint64
}

func (blocksRequest BlockRequest) IsInvalid() bool {
	return blocksRequest.BlockHeight == nil
}
