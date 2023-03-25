package network

type TransactionResponse struct {
	Id        [32]byte
	Inputs    []*InputResponse
	Outputs   []*OutputResponse
	Timestamp int64
}
