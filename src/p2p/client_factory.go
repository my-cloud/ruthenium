package p2p

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"strconv"
	"time"
)

const (
	clientConnectionTimeoutInSeconds = 10
	neighborFindingTimeoutSecond     = 5
)

type ClientFactory struct {
	ipFinder IpFinder
	logger   p2p.Logger
}

func NewClientFactory(ipFinder IpFinder) *ClientFactory {
	return &ClientFactory{ipFinder, log.NewLogger(log.Fatal)}
}

func (factory *ClientFactory) CreateClient(ip string, port uint16, target string) (network.Client, error) {
	lookedUpIps, err := factory.ipFinder.LookupIP(ip)
	if err != nil {
		return nil, fmt.Errorf("DNS discovery failed on addresse %s: %w", ip, err)
	}

	ipsCount := len(lookedUpIps)
	if ipsCount != 1 {
		return nil, fmt.Errorf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", ipsCount, ip)
	}
	_, err = factory.ipFinder.DialTimeout("tcp", target, neighborFindingTimeoutSecond*time.Second)
	if err != nil {
		return nil, err
	}
	tcp := p2p.NewTCP(ip, strconv.Itoa(int(port)))
	var client *p2p.Client
	client, err = p2p.NewClient(tcp)
	if err != nil {
		return nil, err
	}
	settings := p2p.NewClientSettings()
	settings.SetRetry(1, time.Nanosecond)
	settings.SetConnTimeout(clientConnectionTimeoutInSeconds * time.Second)
	client.SetSettings(settings)
	client.SetLogger(factory.logger)
	return client, err
}
