package network

import (
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock/clockmock"
	"github.com/my-cloud/ruthenium/test/node/network/networkmock"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2pmock"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocolmock"
	"testing"
	"time"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(p2pmock.ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleFunc = func(topic string, handler gp2p.Handler) {}
	synchronizerMock := new(networkmock.SynchronizerMock)
	blockchainMock := new(protocolmock.BlockchainMock)
	blockchainMock.VerifyFunc = func(int64) {}
	transactionsPoolMock := new(protocolmock.TransactionsPoolMock)
	engineMock := new(clockmock.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	engineMock.WaitFunc = func() {}
	timeMock := new(clockmock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Now() }
	client := new(p2pmock.ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2pmock.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	logger := log.NewLogger(log.Fatal)
	host := p2p.NewHost(serverMock, synchronizerMock, blockchainMock, transactionsPoolMock, engineMock, engineMock, engineMock, timeMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
