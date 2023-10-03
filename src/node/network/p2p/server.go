package p2p

type Server interface {
	Serve() (err error)
	SetHandleBlocksRequest(endpoint string)
	SetHandleFirstBlockTimestampRequest(endpoint string)
	SetHandleTargetsRequest(endpoint string)
	SetHandleTransactionRequest(endpoint string)
	SetHandleTransactionsRequest(endpoint string)
	SetHandleUtxosRequest(endpoint string)
}
