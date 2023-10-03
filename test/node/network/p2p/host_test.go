package p2p

import (
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(p2ptest.ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleBlocksRequestFunc = func(string) {}
	serverMock.SetHandleFirstBlockTimestampRequestFunc = func(string) {}
	serverMock.SetHandleTargetsRequestFunc = func(string) {}
	serverMock.SetHandleTransactionRequestFunc = func(string) {}
	serverMock.SetHandleTransactionsRequestFunc = func(string) {}
	serverMock.SetHandleUtxosRequestFunc = func(string) {}
	engineMock := new(clocktest.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	engineMock.WaitFunc = func() {}
	logger := logtest.NewLoggerMock()
	host := p2p.NewHost(serverMock, engineMock, engineMock, engineMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
