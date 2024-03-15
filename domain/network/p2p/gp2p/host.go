package gp2p

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/domain/network/p2p"
	"github.com/my-cloud/ruthenium/infrastructure/log/console"
	"strconv"
	"time"
)

type Host struct {
	*gp2p.Server
	handler p2p.Handler
}

func NewHost(port int, handler p2p.Handler, timeout time.Duration) (*Host, error) {
	tcp := gp2p.NewTCP("0.0.0.0", strconv.Itoa(port))
	host, err := gp2p.NewServer(tcp)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate host on port %d: %w", port, err)
	}
	host.SetLogger(console.NewLogger(console.Fatal))
	settings := gp2p.NewServerSettings()
	settings.SetConnTimeout(timeout)
	host.SetSettings(settings)
	return &Host{host, handler}, err
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
