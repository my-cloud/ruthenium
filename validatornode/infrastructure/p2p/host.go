package p2p

import (
	"fmt"
	"strconv"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/config"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
)

type Host struct {
	*gp2p.Server
	handler *Handler
}

func NewHost(port int,
	settings *config.Settings,
	blocksManager protocol.BlocksManager,
	neighborsManager network.NeighborsManager,
	transactionsManager protocol.TransactionsManager,
	utxosManager protocol.UtxosManager,
	watch protocol.TimeProvider) (*Host, error) {
	tcp := gp2p.NewTCP("0.0.0.0", strconv.Itoa(port))
	server, err := gp2p.NewServer(tcp)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host on port %d: %w", port, err)
	}
	server.SetLogger(console.NewLogger(console.Fatal))
	serverSettings := gp2p.NewServerSettings()
	serverSettings.SetConnTimeout(settings.ValidationTimeout())
	server.SetSettings(serverSettings)
	handler := NewHandler(blocksManager, settings.Bytes(), neighborsManager, transactionsManager, utxosManager, watch)
	return &Host{server, handler}, err
}

func (host *Host) SetHandleBlocksRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleBlocksRequest)
}

func (host *Host) SetHandleFirstBlockTimestampRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleFirstBlockTimestampRequest)
}

func (host *Host) SetHandleSettingsRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleSettingsRequest)
}

func (host *Host) SetHandleTargetsRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleTargetsRequest)
}

func (host *Host) SetHandleTransactionRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleTransactionRequest)
}

func (host *Host) SetHandleTransactionsRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleTransactionsRequest)
}

func (host *Host) SetHandleUtxosRequest(endpoint string) {
	host.SetHandle(endpoint, host.handler.HandleUtxosRequest)
}
