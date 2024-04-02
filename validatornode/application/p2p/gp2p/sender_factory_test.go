package gp2p

import (
	"errors"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_CreateSender_IpFinderError_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(net.IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", errors.New("") }
	senderFactory := NewSenderFactory(ipFinder, 0)

	// Act
	client, _ := senderFactory.CreateSender("", "0")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateSender_ValidIp_ReturnsClient(t *testing.T) {
	// Arrange
	ipFinder := new(net.IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", nil }
	senderFactory := NewSenderFactory(ipFinder, 0)

	// Act
	client, _ := senderFactory.CreateSender("", "0")

	// Assert
	test.Assert(t, client != nil, "client is nil whereas it should not")
}
