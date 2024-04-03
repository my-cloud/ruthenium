package presentation

const (
	BlocksEndpoint              = "blocks"
	FirstBlockTimestampEndpoint = "first-block-timestamp"
	SettingsEndpoint            = "settings"
	TargetsEndpoint             = "targets"
	TransactionEndpoint         = "transaction"
	TransactionsEndpoint        = "transactions"
	UtxosEndpoint               = "utxos"
)

type Node struct {
	server  Server
	engines []Pulser
}

func NewNode(server Server, engines ...Pulser) *Node {
	server.SetHandleBlocksRequest(BlocksEndpoint)
	server.SetHandleFirstBlockTimestampRequest(FirstBlockTimestampEndpoint)
	server.SetHandleSettingsRequest(SettingsEndpoint)
	server.SetHandleTargetsRequest(TargetsEndpoint)
	server.SetHandleTransactionRequest(TransactionEndpoint)
	server.SetHandleTransactionsRequest(TransactionsEndpoint)
	server.SetHandleUtxosRequest(UtxosEndpoint)
	return &Node{server, engines}
}

func (node *Node) Run() error {
	for _, engine := range node.engines {
		go engine.Start()
	}
	return node.server.Serve()
}
