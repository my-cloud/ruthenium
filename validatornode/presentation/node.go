package presentation

import "github.com/my-cloud/ruthenium/validatornode/infrastructure/p2p"

type Node struct {
	server  Server
	engines []Pulser
}

func NewNode(server Server, engines ...Pulser) *Node {
	server.SetHandleBlocksRequest(p2p.BlocksEndpoint)
	server.SetHandleFirstBlockTimestampRequest(p2p.FirstBlockTimestampEndpoint)
	server.SetHandleSettingsRequest(p2p.SettingsEndpoint)
	server.SetHandleTargetsRequest(p2p.TargetsEndpoint)
	server.SetHandleTransactionRequest(p2p.TransactionEndpoint)
	server.SetHandleTransactionsRequest(p2p.TransactionsEndpoint)
	server.SetHandleUtxosRequest(p2p.UtxosEndpoint)
	return &Node{server, engines}
}

func (node *Node) Run() error {
	for _, engine := range node.engines {
		go engine.Start()
	}
	return node.server.Serve()
}
