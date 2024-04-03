package presentation

type Server interface {
	Serve() (err error)
	SetHandleBlocksRequest(endpoint string)
	SetHandleFirstBlockTimestampRequest(endpoint string)
	SetHandleSettingsRequest(endpoint string)
	SetHandleTargetsRequest(endpoint string)
	SetHandleTransactionRequest(endpoint string)
	SetHandleTransactionsRequest(endpoint string)
	SetHandleUtxosRequest(endpoint string)
}
