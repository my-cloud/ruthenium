package network

import (
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"github.com/my-cloud/ruthenium/test/node/clock"
	"testing"
	"time"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(mock.ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleFunc = func(topic string, handler p2p.Handler) {}
	blockchainMock := new(mock.BlockchainMock)
	transactionsPoolMock := new(mock.TransactionsPoolMock)
	engineMock := new(mock.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	timeMock := new(clock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Now() }
	client := new(mock.ClientMock)
	client.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	clientFactoryMock := new(mock.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (network.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, err := network.NewSynchronizer(0, timeMock, clientFactoryMock, configurationPath, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	host := network.NewHost(serverMock, blockchainMock, transactionsPoolMock, engineMock, engineMock, synchronizer, timeMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
