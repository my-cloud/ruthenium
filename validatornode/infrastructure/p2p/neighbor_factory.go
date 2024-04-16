package p2p

import (
	"fmt"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type NeighborFactory struct {
	ipFinder          IpFinder
	connectionTimeout time.Duration
	logger            log.Logger
}

func NewNeighborFactory(ipFinder IpFinder, connectionTimeout time.Duration, logger log.Logger) *NeighborFactory {
	return &NeighborFactory{ipFinder, connectionTimeout, logger}
}

func (factory *NeighborFactory) CreateSender(ip string, port string) (network.Sender, error) {
	lookedUpIp, err := factory.ipFinder.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to look up IP on addresse %s: %w", ip, err)
	}
	neighbor, err := NewNeighbor(lookedUpIp, port, factory.connectionTimeout, factory.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create neighbor for address %s: %w", ip, err)
	}
	return neighbor, err
}
