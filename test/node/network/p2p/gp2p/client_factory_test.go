package gp2p

import (
	"errors"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"net"
	"testing"
)

func Test_CreateClient_NoSuchHost_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(networktest.IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return nil, errors.New("no such host")
	}
	clientFactory := gp2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_NoIpAddress_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(networktest.IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return []net.IP{}, nil
	}
	clientFactory := gp2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_ValidIp_ReturnsClient(t *testing.T) {
	// Arrange
	ipFinder := new(networktest.IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return []net.IP{{}}, nil
	}
	clientFactory := gp2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client != nil, "client is nil whereas it should not")
}
