package gp2p

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"strconv"
	"time"
)

const (
	clientConnectionTimeoutInSeconds = 10
	neighborFindingTimeoutSecond     = 5
)

type ClientFactory struct {
	ipFinder network.IpFinder
	logger   log.Logger
}

func NewClientFactory(ipFinder network.IpFinder) *ClientFactory {
	return &ClientFactory{ipFinder, console.NewLogger(console.Fatal)}
}

func (factory *ClientFactory) CreateClient(ip string, port uint16, target string) (p2p.Client, error) {
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
	tcp := gp2p.NewTCP(ip, strconv.Itoa(int(port)))
	var client *gp2p.Client
	client, err = gp2p.NewClient(tcp)
	if err != nil {
		return nil, err
	}
	settings := gp2p.NewClientSettings()
	settings.SetRetry(1, time.Nanosecond)
	settings.SetConnTimeout(clientConnectionTimeoutInSeconds * time.Second)
	client.SetSettings(settings)
	client.SetLogger(factory.logger)
	return client, err
}
