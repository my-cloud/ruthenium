package p2p

import (
	"github.com/my-cloud/ruthenium/common/application"
	"github.com/my-cloud/ruthenium/common/infrastructure/log"
)

const (
	blocksEndpoint              = "blocks"
	firstBlockTimestampEndpoint = "first-block-timestamp"
	settingsEndpoint            = "settings"
	targetsEndpoint             = "targets"
	transactionEndpoint         = "transaction"
	transactionsEndpoint        = "transactions"
	utxosEndpoint               = "utxos"
)

type Node struct {
	server  Server
	logger  log.Logger
	engines []application.Pulser
}

func NewNode(server Server, logger log.Logger, engines ...application.Pulser) *Node {
	return &Node{server, logger, engines}
}

func (node *Node) Run() error {
	for _, engine := range node.engines {
		go engine.Start()
	}
	node.setServerHandles()
	return node.server.Serve()
}

func (node *Node) setServerHandles() {
	node.server.SetHandleBlocksRequest(blocksEndpoint)
	node.server.SetHandleFirstBlockTimestampRequest(firstBlockTimestampEndpoint)
	node.server.SetHandleSettingsRequest(settingsEndpoint)
	node.server.SetHandleTargetsRequest(targetsEndpoint)
	node.server.SetHandleTransactionRequest(transactionEndpoint)
	node.server.SetHandleTransactionsRequest(transactionsEndpoint)
	node.server.SetHandleUtxosRequest(utxosEndpoint)
}
