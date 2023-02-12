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
	host.server.SetHandle("dialog", host.handler.Handle)
	return host.startServer()
}

func (host *Host) startBlockchain() {
	host.logger.Info("updating the blockchain...")
	host.synchronizationEngine.Do()
	host.logger.Info("neighbors are synchronized")
	go host.synchronizationEngine.Start()
	host.verificationEngine.Do()
	host.logger.Info("the blockchain is now up to date")
	host.validationEngine.Do()
	go host.validationEngine.Start()
	go host.verificationEngine.Start()
}

func (host *Host) startServer() error {
	host.logger.Info("host node started...")
	return host.server.Serve()
}
