package connection

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/connection"
	"github.com/my-cloud/ruthenium/src/log"
	"net"
	"strconv"
	"time"
)

const neighborFindingTimeoutSecond = 5

type Peering struct {
	logger p2p.Logger
}

func NewPeering() *Peering {
	return &Peering{log.NewLogger(log.Fatal)}
}

func (peering *Peering) CreateSender(ip string, port uint16, target string) (connection.Sender, error) {
	lookedUpIps, err := net.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("DNS discovery failed on addresse %s: %w", ip, err)
	}

	ipsCount := len(lookedUpIps)
	if ipsCount != 1 {
		return nil, fmt.Errorf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", ipsCount, ip)
	}
	_, err = net.DialTimeout("tcp", target, neighborFindingTimeoutSecond*time.Second)
	if err != nil {
		return nil, err
	}
	tcp := p2p.NewTCP(ip, strconv.Itoa(int(port)))
	var client *p2p.Client
	client, err = p2p.NewClient(tcp)
	if err != nil {
		return nil, fmt.Errorf("failed to start client for target %s: %w", target, err)
	}
	client.SetLogger(peering.logger)
	return client, err
}
