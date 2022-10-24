package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/connection"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	configurationPath := "../../"
	watch := clock.NewWatch()
	senderProviderMock := new(SenderProviderMock)
	sender := new(SenderMock)
	sender.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	senderProviderMock.CreateSenderFunc = func(string, uint16, string) (connection.Sender, error) { return sender, nil }
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
