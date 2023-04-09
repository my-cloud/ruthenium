package network

type TransactionResponse struct {
	Id        string
	Inputs    []*InputResponse
	Outputs   []*OutputResponse
	Timestamp int64
}
