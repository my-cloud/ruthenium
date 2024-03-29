package p2p

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
	"time"
)

func Test_AddTargets_MoreThanOneTarget_IncentiveTargetsSender(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientMock := new(p2ptest.ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	scoresBySeedTarget := map[string]int{}
	synchronizer := p2p.NewSynchronizer(clientFactoryMock, "0.0.0.0", "0", 1, scoresBySeedTarget, watchMock)
	target1 := "0.0.0.0:1"
	target2 := "0.0.0.0:0"
	targetRequests := []string{target1, target2}

	// Act
	synchronizer.AddTargets(targetRequests)

	// Assert
	synchronizer.Synchronize(0)
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
	neighborTarget := neighbors[0].Target()
	test.Assert(t, neighborTarget == target1, fmt.Sprintf("Wrong neighbor. Expected: %s - Actual: %s", target1, neighborTarget))
}

func Test_Incentive_TargetIsNotKnown_TargetIncentive(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientMock := new(p2ptest.ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	scoresBySeedTarget := map[string]int{}
	synchronizer := p2p.NewSynchronizer(clientFactoryMock, "0.0.0.0", "0", 1, scoresBySeedTarget, watchMock)
	expectedTarget := "0.0.0.0:1"

	// Act
	synchronizer.Incentive(expectedTarget)

	// Assert
	synchronizer.Synchronize(0)
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
	target := neighbors[0].Target()
	test.Assert(t, target == expectedTarget, fmt.Sprintf("Wrong target. Expected: %s - Actual: %s", expectedTarget, target))
}

func Test_Synchronize_OneSeed_NeighborAdded(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientMock := new(p2ptest.ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	scoresBySeedTarget := map[string]int{"0.0.0.0:1": 0}
	synchronizer := p2p.NewSynchronizer(clientFactoryMock, "0.0.0.0", "0", 1, scoresBySeedTarget, watchMock)

	// Act
	synchronizer.Synchronize(0)

	// Assert
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
