package gp2p

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
)

type ClientFactory struct {
	ipFinder network.IpFinder
	settings p2p.Settings
}

func NewClientFactory(ipFinder network.IpFinder, settings p2p.Settings) *ClientFactory {
	return &ClientFactory{ipFinder, settings}
}

func (factory *ClientFactory) CreateClient(ip string, port string) (p2p.Client, error) {
	lookedUpIp, err := factory.ipFinder.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to look up IP on addresse %s: %w", ip, err)
	}
	client, err := NewClient(lookedUpIp, port, factory.settings.ValidationTimeout(), console.NewLogger(console.Fatal))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate client for address %s: %w", ip, err)
	}
	return client, err
}
