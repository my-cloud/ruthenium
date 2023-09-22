package p2p

import (
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
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
	host.server.SetHandle("block", host.handler.HandleBlockRequest)
	host.server.SetHandle("blocks", host.handler.HandleBlocksRequest)
	host.server.SetHandle("targets", host.handler.HandleTargetsRequest)
	host.server.SetHandle("transaction", host.handler.HandleTransactionRequest)
	host.server.SetHandle("transactions", host.handler.HandleTransactionsRequest)
	host.server.SetHandle("utxos", host.handler.HandleUtxosRequest)
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
