package network

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/api/connection"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node"
	"testing"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	configurationPath := "../../"
	watch := node.NewWatchMock()
	senderProviderMock := new(SenderProviderMock)
	senderProviderMock.CreateSenderFunc = func(string, uint16, string) (connection.Sender, error) { return new(SenderMock), nil }
	logger := log.NewLogger(log.Fatal)
	neighborhood := network.NewNeighborhood("", 0, watch, senderProviderMock, configurationPath, logger)

	// Act
	neighborhood.Synchronize()

	// Assert
	neighborhood.Wait()
	neighbors := neighborhood.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
