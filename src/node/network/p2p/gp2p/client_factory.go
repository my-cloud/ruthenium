package gp2p

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
)

type ClientFactory struct {
	ipFinder network.IpFinder
	logger   log.Logger
}

func NewClientFactory(ipFinder network.IpFinder) *ClientFactory {
	return &ClientFactory{ipFinder, console.NewLogger(console.Fatal)}
}

func (factory *ClientFactory) CreateClient(ip string, port string) (p2p.Client, error) {
	lookedUpIp, err := factory.ipFinder.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to look up IP on addresse %s: %w", ip, err)
	}
	client, err := NewClient(lookedUpIp, port, factory.logger)
	if err != nil {
		return nil, err
	}
	return client, err
}
