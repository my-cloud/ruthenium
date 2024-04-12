package p2p

import (
	"fmt"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"
)

type NeighborFactory struct {
	ipFinder          net.IpFinder
	connectionTimeout time.Duration
}

func NewNeighborFactory(ipFinder net.IpFinder, connectionTimeout time.Duration) *NeighborFactory {
	return &NeighborFactory{ipFinder, connectionTimeout}
}

func (factory *NeighborFactory) CreateSender(ip string, port string) (network.Sender, error) {
	lookedUpIp, err := factory.ipFinder.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to look up IP on addresse %s: %w", ip, err)
	}
	neighbor, err := NewNeighbor(lookedUpIp, port, factory.connectionTimeout, console.NewLogger(console.Fatal))
	if err != nil {
		return nil, fmt.Errorf("failed to create neighbor for address %s: %w", ip, err)
	}
	return neighbor, err
}
