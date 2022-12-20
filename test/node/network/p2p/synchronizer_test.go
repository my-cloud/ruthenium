package p2p

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
	"time"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	client := new(p2ptest.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := "../../../"
	logger := logtest.NewLoggerMock()
	synchronizer, _ := p2p.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize(0)

	// Assert
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
