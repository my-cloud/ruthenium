package p2p

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
	"time"
)

func Test_Synchronize_OneSeed_NeighborAdded(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	client := new(p2ptest.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	client.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := "../../../config"
	logger := logtest.NewLoggerMock()
	synchronizer, _ := p2p.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)

	// Act
	synchronizer.Synchronize(0)

	// Assert
	neighbors := synchronizer.Neighbors()
	expectedNeighborsCount := 1
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}

func Test_Synchronize_NoConfigurationFolder_ReturnsError(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	client := new(p2ptest.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := ""
	logger := logtest.NewLoggerMock()

	// Act
	_, err := p2p.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)

	// Assert
	test.Assert(t, err != nil, "No error returned whereas it should")
}

func Test_Synchronize_EmptyConfigurationFile_ReturnsError(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	client := new(p2ptest.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := "../../../config/emptyconfigfile"
	logger := logtest.NewLoggerMock()

	// Act
	_, err := p2p.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)

	// Assert
	test.Assert(t, err != nil, "No error returned whereas it should")
}

func Test_Synchronize_WrongSeedFormat_ReturnsError(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	client := new(p2ptest.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := "../../../config/wrongseedformat"
	logger := logtest.NewLoggerMock()

	// Act
	_, err := p2p.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)

	// Assert
	test.Assert(t, err != nil, "No error returned whereas it should")
}

func Test_Synchronize_WrongSeedPortFormat_ReturnsError(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	client := new(p2ptest.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	configurationPath := "../../../config/wrongseedportformat"
	logger := logtest.NewLoggerMock()

	// Act
	_, err := p2p.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)

	// Assert
	test.Assert(t, err != nil, "No error returned whereas it should")
}
