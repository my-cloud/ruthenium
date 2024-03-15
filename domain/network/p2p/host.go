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

type Host struct {
	server                Server
	synchronizationEngine domain.Pulser
	validationEngine      domain.Pulser
	verificationEngine    domain.Pulser
	logger                log.Logger
}

func NewHost(server Server, synchronizationEngine domain.Pulser, validationEngine domain.Pulser, verificationEngine domain.Pulser, logger log.Logger) *Host {
	return &Host{server, synchronizationEngine, validationEngine, verificationEngine, logger}
}

func (host *Host) Run() error {
	go host.synchronizationEngine.Start()
	go host.validationEngine.Start()
	go host.verificationEngine.Start()
	host.setServerHandles()
	return host.server.Serve()
}

func (host *Host) setServerHandles() {
	host.server.SetHandleBlocksRequest(blocksEndpoint)
	host.server.SetHandleFirstBlockTimestampRequest(firstBlockTimestampEndpoint)
	host.server.SetHandleSettingsRequest(settingsEndpoint)
	host.server.SetHandleTargetsRequest(targetsEndpoint)
	host.server.SetHandleTransactionRequest(transactionEndpoint)
	host.server.SetHandleTransactionsRequest(transactionsEndpoint)
	host.server.SetHandleUtxosRequest(utxosEndpoint)
}
