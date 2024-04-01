package gp2p

import (
	"fmt"
	"time"

	"github.com/my-cloud/ruthenium/common/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/common/infrastructure/net"
	"github.com/my-cloud/ruthenium/validatornode/application/p2p"
)

type SenderFactory struct {
	ipFinder          net.IpFinder
	connectionTimeout time.Duration
}

func NewSenderFactory(ipFinder net.IpFinder, connectionTimeout time.Duration) *SenderFactory {
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
