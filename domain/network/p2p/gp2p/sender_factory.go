package gp2p

import (
	"fmt"
	"github.com/my-cloud/ruthenium/domain/network"
	"github.com/my-cloud/ruthenium/domain/network/p2p"
	"github.com/my-cloud/ruthenium/infrastructure/log/console"
	"time"
)

type SenderFactory struct {
	ipFinder          network.IpFinder
	connectionTimeout time.Duration
}

func NewSenderFactory(ipFinder network.IpFinder, connectionTimeout time.Duration) *SenderFactory {
	return &SenderFactory{ipFinder, connectionTimeout}
}

func (factory *SenderFactory) CreateSender(ip string, port string) (p2p.Sender, error) {
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