package p2p

import (
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock"
	"github.com/my-cloud/ruthenium/test/node/network"
	"github.com/my-cloud/ruthenium/test/node/protocol"
	"testing"
	"time"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleFunc = func(topic string, handler gp2p.Handler) {}
	synchronizerMock := new(network.SynchronizerMock)
	blockchainMock := new(protocol.BlockchainMock)
	blockchainMock.VerifyFunc = func(int64) {}
	transactionsPoolMock := new(protocol.TransactionsPoolMock)
	engineMock := new(clock.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	engineMock.WaitFunc = func() {}
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Now() }
	client := new(ClientMock)
	client.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (p2p.Client, error) { return client, nil }
	logger := console.NewLogger(console.Fatal)
	host := p2p.NewHost(serverMock, synchronizerMock, blockchainMock, transactionsPoolMock, engineMock, engineMock, engineMock, watchMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
