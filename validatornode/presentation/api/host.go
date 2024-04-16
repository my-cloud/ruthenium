package api

import (
	"fmt"
	"strconv"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/config"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
)

type Host struct {
	*gp2p.Server
	blocksController       *BlocksController
	sendersController      *SendersController
	settingsController     *SettingsController
	transactionsController *TransactionsController
	utxosController        *UtxosController
}

func NewHost(port int,
	settings *config.Settings,
	blocksManager ledger.BlocksManager,
	sendersManager network.SendersManager,
	transactionsManager ledger.TransactionsManager,
	utxosManager ledger.UtxosManager) (*Host, error) {
	tcp := gp2p.NewTCP("0.0.0.0", strconv.Itoa(port))
	server, err := gp2p.NewServer(tcp)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host on port %d: %w", port, err)
	}
	server.SetLogger(console.NewLogger(console.Fatal))
	serverSettings := gp2p.NewServerSettings()
	serverSettings.SetConnTimeout(settings.ValidationTimeout())
	server.SetSettings(serverSettings)
	blocksController := NewBlocksController(blocksManager)
	sendersController := NewSendersController(sendersManager)
	settingsController := NewSettingsController(settings.Bytes())
	transactionsController := NewTransactionsController(sendersManager, transactionsManager)
	utxosController := NewUtxosController(utxosManager)
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
