package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"github.com/my-cloud/ruthenium/test/node/clock"
	"testing"
	"time"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	timeMock := new(clock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(mock.ClientFactoryMock)
	client := new(mock.ClientMock)
	client.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (network.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, _ := network.NewSynchronizer(0, timeMock, clientFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize()

	// Assert
	synchronizer.Wait()
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
