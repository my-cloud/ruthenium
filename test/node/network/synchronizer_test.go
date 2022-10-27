package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	watch := clock.NewWatch()
	senderFactoryMock := new(SenderFactoryMock)
	sender := new(SenderMock)
	sender.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	senderFactoryMock.CreateSenderFunc = func(string, uint16, string) (network.Sender, error) { return sender, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer := network.NewSynchronizer("", 0, watch, senderFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize()

	// Assert
	synchronizer.Wait()
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
