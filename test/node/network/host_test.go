package network

import (
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	p2p2 "github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"testing"
	"time"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(mock.ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleFunc = func(topic string, handler p2p.Handler) {}
	blockchainMock := new(mock.BlockchainMock)
	blockchainMock.VerifyFunc = func(int64) {}
	transactionsPoolMock := new(mock.TransactionsPoolMock)
	engineMock := new(mock.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Now() }
	client := new(mock.ClientMock)
	client.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	clientFactoryMock := new(mock.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p2.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, err := p2p2.NewSynchronizer(0, timeMock, clientFactoryMock, configurationPath, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	host := p2p2.NewHost(serverMock, blockchainMock, transactionsPoolMock, engineMock, engineMock, synchronizer, timeMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
