package gp2p

import (
	"fmt"
	"github.com/my-cloud/ruthenium/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/infrastructure/network"
	"github.com/my-cloud/ruthenium/infrastructure/network/p2p"
	"time"
)

type ClientFactory struct {
	ipFinder          network.IpFinder
	connectionTimeout time.Duration
}

func NewClientFactory(ipFinder network.IpFinder, connectionTimeout time.Duration) *ClientFactory {
	return &ClientFactory{ipFinder, connectionTimeout}
}

func (factory *ClientFactory) CreateClient(ip string, port string) (p2p.Client, error) {
	lookedUpIp, err := factory.ipFinder.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to look up IP on addresse %s: %w", ip, err)
	}
	client, err := NewClient(lookedUpIp, port, factory.connectionTimeout, console.NewLogger(console.Fatal))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate client for address %s: %w", ip, err)
	}
	return client, err
}
