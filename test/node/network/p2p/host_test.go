package p2p

import (
	"context"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	handlerMock := new(networktest.HandlerMock)
	handlerMock.HandleBlocksRequestFunc = func(context.Context, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	handlerMock.HandleFirstBlockTimestampRequestFunc = func(context.Context, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	handlerMock.HandleTargetsRequestFunc = func(context.Context, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	handlerMock.HandleTransactionRequestFunc = func(context.Context, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	handlerMock.HandleTransactionsRequestFunc = func(context.Context, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	handlerMock.HandleUtxosRequestFunc = func(context.Context, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	serverMock := new(p2ptest.ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleFunc = func(topic string, handler gp2p.Handler) {}
	engineMock := new(clocktest.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	engineMock.WaitFunc = func() {}
	logger := logtest.NewLoggerMock()
	host := p2p.NewHost(handlerMock, serverMock, engineMock, engineMock, engineMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
