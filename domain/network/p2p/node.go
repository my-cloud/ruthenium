package p2p

import (
	"github.com/my-cloud/ruthenium/domain"
	"github.com/my-cloud/ruthenium/infrastructure/log"
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
	server                Server
	synchronizationEngine domain.Pulser
	validationEngine      domain.Pulser
	verificationEngine    domain.Pulser
	logger                log.Logger
}

func NewNode(server Server, synchronizationEngine domain.Pulser, validationEngine domain.Pulser, verificationEngine domain.Pulser, logger log.Logger) *Node {
	return &Node{server, synchronizationEngine, validationEngine, verificationEngine, logger}
}

func (node *Node) Run() error {
	go node.synchronizationEngine.Start()
	go node.validationEngine.Start()
	go node.verificationEngine.Start()
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
