package p2p

import (
	"fmt"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"
)

type SenderFactory struct {
	ipFinder          net.IpFinder
	connectionTimeout time.Duration
}

func NewSenderFactory(ipFinder net.IpFinder, connectionTimeout time.Duration) *SenderFactory {
	return &SenderFactory{ipFinder, connectionTimeout}
}

func (factory *SenderFactory) CreateSender(ip string, port string) (network.Sender, error) {
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
