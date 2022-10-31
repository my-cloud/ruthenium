package p2p

import (
	"errors"
	"github.com/my-cloud/ruthenium/src/p2p"
	"github.com/my-cloud/ruthenium/test"
	"net"
	"testing"
	"time"
)

func Test_CreateClient_NoSuchHost_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return nil, errors.New("no such host")
	}
	clientFactory := p2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", 0, "")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_NoIpAddress_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return []net.IP{}, nil
	}
	clientFactory := p2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", 0, "")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_NoReachableIpAddress_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return []net.IP{{}}, nil
	}
	ipFinder.DialTimeoutFunc = func(string, string, time.Duration) (net.Conn, error) {
		return nil, errors.New("missing address")
	}
	clientFactory := p2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", 0, "")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_ValidIp_ReturnsClient(t *testing.T) {
	// Arrange
	ipFinder := new(IpFinderMock)
	ipFinder.LookupIPFunc = func(ip string) ([]net.IP, error) {
		return []net.IP{{}}, nil
	}
	ipFinder.DialTimeoutFunc = func(string, string, time.Duration) (net.Conn, error) {
		return nil, nil
	}
	clientFactory := p2p.NewClientFactory(ipFinder)

	// Act
	client, _ := clientFactory.CreateClient("", 0, "")

	// Assert
	test.Assert(t, client != nil, "client is nil whereas it should not")
}
