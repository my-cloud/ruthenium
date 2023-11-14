package gp2p

import (
	"errors"
	"github.com/my-cloud/ruthenium/infrastructure/network"
	"github.com/my-cloud/ruthenium/infrastructure/test"
	"testing"
)

func Test_CreateClient_IpFinderError_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(network.IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", errors.New("") }
	clientFactory := NewClientFactory(ipFinder, 0)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateClient_ValidIp_ReturnsClient(t *testing.T) {
	// Arrange
	ipFinder := new(network.IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", nil }
	clientFactory := NewClientFactory(ipFinder, 0)

	// Act
	client, _ := clientFactory.CreateClient("", "0")

	// Assert
	test.Assert(t, client != nil, "client is nil whereas it should not")
}
