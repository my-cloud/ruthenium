package network

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock/clockmock"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2pmock"
	"testing"
	"time"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	timeMock := new(clockmock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2pmock.ClientFactoryMock)
	client := new(p2pmock.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, _ := p2p.NewSynchronizer(0, timeMock, clientFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize()

	// Assert
	synchronizer.Wait()
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
