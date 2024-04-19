package api

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api/history"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api/network"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api/payment"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api/protocol"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api/wallet"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/configuration"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
)

type Host struct {
	*gp2p.Server
	blocksController       *history.BlocksController
	sendersController      *network.SendersController
	settingsController     *protocol.SettingsController
	transactionsController *payment.TransactionsController
	utxosController        *wallet.UtxosController
}

func NewHost(settings *configuration.Settings,
	blocksManager application.BlocksManager,
	sendersManager application.SendersManager,
	transactionsManager application.TransactionsManager,
	utxosManager application.UtxosManager) (*Host, error) {
	port := settings.Port()
	tcp := gp2p.NewTCP("0.0.0.0", settings.Port())
	server, err := gp2p.NewServer(tcp)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host on port %d: %w", port, err)
	}
	server.SetLogger(console.NewLogger(console.Fatal))
	serverSettings := gp2p.NewServerSettings()
	serverSettings.SetConnTimeout(settings.ValidationTimeout())
	server.SetSettings(serverSettings)
	blocksController := history.NewBlocksController(blocksManager)
	sendersController := network.NewSendersController(sendersManager)
	settingsController := protocol.NewSettingsController(settings.Bytes())
	transactionsController := payment.NewTransactionsController(sendersManager, transactionsManager)
	utxosController := wallet.NewUtxosController(utxosManager)
	return &Host{server, blocksController, sendersController, settingsController, transactionsController, utxosController}, err
}

func (host *Host) SetHandleBlocksRequest(endpoint string) {
	host.SetHandle(endpoint, host.blocksController.HandleBlocksRequest)
}

func (host *Host) SetHandleFirstBlockTimestampRequest(endpoint string) {
	host.SetHandle(endpoint, host.blocksController.HandleFirstBlockTimestampRequest)
}

func (host *Host) SetHandleSettingsRequest(endpoint string) {
	host.SetHandle(endpoint, host.settingsController.HandleSettingsRequest)
}

func (host *Host) SetHandleTargetsRequest(endpoint string) {
	host.SetHandle(endpoint, host.sendersController.HandleTargetsRequest)
}

func (host *Host) SetHandleTransactionRequest(endpoint string) {
	host.SetHandle(endpoint, host.transactionsController.HandleTransactionRequest)
}

func (host *Host) SetHandleTransactionsRequest(endpoint string) {
	host.SetHandle(endpoint, host.transactionsController.HandleTransactionsRequest)
}

func (host *Host) SetHandleUtxosRequest(endpoint string) {
	host.SetHandle(endpoint, host.utxosController.HandleUtxosRequest)
}
