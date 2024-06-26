package network

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_AddTargets_MoreThanOneTarget_IncentiveTargetsSender(t *testing.T) {
	// Arrange
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	senderCreatorMock := new(application.SenderCreatorMock)
	senderMock := new(application.SenderMock)
	senderMock.TargetFunc = func() string { return "0.0.0.0:1" }
	senderMock.SendTargetsFunc = func([]string) error { return nil }
	senderCreatorMock.CreateSenderFunc = func(string, string) (application.Sender, error) { return senderMock, nil }
	scoresBySeedTarget := map[string]int{}
	neighborhood := NewNeighborhood(senderCreatorMock, "0.0.0.0", "0", 1, scoresBySeedTarget, watchMock)
	target1 := "0.0.0.0:1"
	target2 := "0.0.0.0:0"
	targetRequests := []string{target1, target2}

	// Act
	neighborhood.AddTargets(targetRequests)

	// Assert
	neighborhood.Synchronize(0)
	neighbors := neighborhood.Senders()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
	neighborTarget := neighbors[0].Target()
	test.Assert(t, neighborTarget == target1, fmt.Sprintf("Wrong neighbor. Expected: %s - Actual: %s", target1, neighborTarget))
}

func Test_Incentive_TargetIsNotKnown_TargetIncentive(t *testing.T) {
	// Arrange
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	senderCreatorMock := new(application.SenderCreatorMock)
	senderMock := new(application.SenderMock)
	senderMock.TargetFunc = func() string { return "0.0.0.0:1" }
	senderMock.SendTargetsFunc = func([]string) error { return nil }
	senderCreatorMock.CreateSenderFunc = func(string, string) (application.Sender, error) { return senderMock, nil }
	scoresBySeedTarget := map[string]int{}
	neighborhood := NewNeighborhood(senderCreatorMock, "0.0.0.0", "0", 1, scoresBySeedTarget, watchMock)
	expectedTarget := "0.0.0.0:1"

	// Act
	neighborhood.Incentive(expectedTarget)

	// Assert
	neighborhood.Synchronize(0)
	neighbors := neighborhood.Senders()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
	target := neighbors[0].Target()
	test.Assert(t, target == expectedTarget, fmt.Sprintf("Wrong target. Expected: %s - Actual: %s", expectedTarget, target))
}

func Test_Synchronize_OneSeed_NeighborAdded(t *testing.T) {
	// Arrange
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	senderCreatorMock := new(application.SenderCreatorMock)
	senderMock := new(application.SenderMock)
	senderMock.TargetFunc = func() string { return "0.0.0.0:1" }
	senderMock.SendTargetsFunc = func([]string) error { return nil }
	senderCreatorMock.CreateSenderFunc = func(string, string) (application.Sender, error) { return senderMock, nil }
	scoresBySeedTarget := map[string]int{"0.0.0.0:1": 0}
	neighborhood := NewNeighborhood(senderCreatorMock, "0.0.0.0", "0", 1, scoresBySeedTarget, watchMock)

	// Act
	neighborhood.Synchronize(0)

	// Assert
	neighbors := neighborhood.Senders()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
