package p2p

import (
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
)

const (
	blocksEndpoint              = "blocks"
	firstBlockTimestampEndpoint = "first-block-timestamp"
	targetsEndpoint             = "targets"
	transactionEndpoint         = "transaction"
	transactionsEndpoint        = "transactions"
	utxosEndpoint               = "utxos"
)

type Host struct {
	server                Server
	synchronizationEngine clock.Engine
	validationEngine      clock.Engine
	verificationEngine    clock.Engine
	logger                log.Logger
}

func NewHost(server Server, synchronizationEngine clock.Engine, validationEngine clock.Engine, verificationEngine clock.Engine, logger log.Logger) *Host {
	return &Host{server, synchronizationEngine, validationEngine, verificationEngine, logger}
}

func (host *Host) Run() error {
	host.startBlockchain()
	host.server.SetHandleBlocksRequest(blocksEndpoint)
	host.server.SetHandleFirstBlockTimestampRequest(firstBlockTimestampEndpoint)
	host.server.SetHandleTargetsRequest(targetsEndpoint)
	host.server.SetHandleTransactionRequest(transactionEndpoint)
	host.server.SetHandleTransactionsRequest(transactionsEndpoint)
	host.server.SetHandleUtxosRequest(utxosEndpoint)
	return host.startServer()
}

func (host *Host) startBlockchain() {
	host.synchronizationEngine.Do()
	host.logger.Info("neighbors are synchronized, updating the blockchain...")
	host.verificationEngine.Do()
	host.logger.Info("the blockchain is now up to date, starting validation...")
	go host.synchronizationEngine.Start()
	go host.validationEngine.Start()
	go host.verificationEngine.Start()
}

func (host *Host) startServer() error {
	return host.server.Serve()
}
