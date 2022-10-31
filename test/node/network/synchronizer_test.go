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
	clientFactoryMock := new(ClientFactoryMock)
	client := new(ClientMock)
	client.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (network.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, _ := network.NewSynchronizer(0, watch, clientFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize()

	// Assert
	synchronizer.Wait()
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
