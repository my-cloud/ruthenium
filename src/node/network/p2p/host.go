package p2p

import (
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
)

const (
	blocksEndpoint              = "blocks"
	firstBlockTimestampEndpoint = "first-block-timestamp"
	lastBlockTimestampEndpoint  = "last-block-timestamp"
	targetsEndpoint             = "targets"
	transactionEndpoint         = "transaction"
	transactionsEndpoint        = "transactions"
	utxosEndpoint               = "utxos"
)

type Host struct {
	handler               network.Handler
	server                Server
	synchronizationEngine clock.Engine
	validationEngine      clock.Engine
	verificationEngine    clock.Engine
	logger                log.Logger
}

func NewHost(
	handler network.Handler,
	server Server,
	synchronizationEngine clock.Engine,
	validationEngine clock.Engine,
	verificationEngine clock.Engine,
	logger log.Logger,
) *Host {
	return &Host{handler, server, synchronizationEngine, validationEngine, verificationEngine, logger}
}

func (host *Host) Run() error {
	host.startBlockchain()
	host.server.SetHandle(blocksEndpoint, host.handler.HandleBlocksRequest)
	host.server.SetHandle(firstBlockTimestampEndpoint, host.handler.HandleFirstBlockTimestampRequest)
	host.server.SetHandle(lastBlockTimestampEndpoint, host.handler.HandleLastBlockTimestampRequest)
	host.server.SetHandle(targetsEndpoint, host.handler.HandleTargetsRequest)
	host.server.SetHandle(transactionEndpoint, host.handler.HandleTransactionRequest)
	host.server.SetHandle(transactionsEndpoint, host.handler.HandleTransactionsRequest)
	host.server.SetHandle(utxosEndpoint, host.handler.HandleUtxosRequest)
	return host.startServer()
}

func (host *Host) startBlockchain() {
	host.synchronizationEngine.Do()
	host.logger.Info("neighbors are synchronized, updating the blockchain...")
	go host.synchronizationEngine.Start()
	host.verificationEngine.Do()
	host.logger.Info("the blockchain is now up to date, starting validation...")
	go host.validationEngine.Start()
	go host.verificationEngine.Start()
}

func (host *Host) startServer() error {
	return host.server.Serve()
}
