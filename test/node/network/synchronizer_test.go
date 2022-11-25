package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	p2p2 "github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"testing"
	"time"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(mock.ClientFactoryMock)
	client := new(mock.ClientMock)
	client.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p2.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, _ := p2p2.NewSynchronizer(0, timeMock, clientFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize()

	// Assert
	synchronizer.Wait()
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
