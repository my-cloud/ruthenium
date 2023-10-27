package gp2p

import (
	"errors"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"testing"
)

func Test_CreateClient_IpFinderError_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(networktest.IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", errors.New("") }
	clientFactory := gp2p.NewClientFactory(ipFinder, 0)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_ValidIp_ReturnsClient(t *testing.T) {
	// Arrange
	ipFinder := new(networktest.IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", nil }
	clientFactory := gp2p.NewClientFactory(ipFinder, 0)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client != nil, "client is nil whereas it should not")
}
