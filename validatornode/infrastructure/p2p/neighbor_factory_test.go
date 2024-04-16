package p2p

import (
	"errors"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_CreateSender_IpFinderError_ReturnsNil(t *testing.T) {
	// Arrange
	ipFinder := new(IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", errors.New("") }
	logger := log.NewLoggerMock()
	neighborFactory := NewNeighborFactory(ipFinder, 0, logger)

	// Act
	client, _ := neighborFactory.CreateSender("", "0")

	// Assert
	test.Assert(t, client == nil, "client is not nil whereas it should be")
}

func Test_CreateSender_ValidIp_ReturnsClient(t *testing.T) {
	// Arrange
	ipFinder := new(IpFinderMock)
	ipFinder.LookupIPFunc = func(string) (string, error) { return "", nil }
	logger := log.NewLoggerMock()
	neighborFactory := NewNeighborFactory(ipFinder, 0, logger)

	// Act
	client, _ := neighborFactory.CreateSender("", "0")

	// Assert
	test.Assert(t, client != nil, "client is nil whereas it should not")
}
