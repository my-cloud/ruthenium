package p2p

import (
	"github.com/my-cloud/ruthenium/domain/clock"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"github.com/my-cloud/ruthenium/infrastructure/test"
	"testing"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleBlocksRequestFunc = func(string) {}
	serverMock.SetHandleFirstBlockTimestampRequestFunc = func(string) {}
	serverMock.SetHandleSettingsRequestFunc = func(string) {}
	serverMock.SetHandleTargetsRequestFunc = func(string) {}
	serverMock.SetHandleTransactionRequestFunc = func(string) {}
	serverMock.SetHandleTransactionsRequestFunc = func(string) {}
	serverMock.SetHandleUtxosRequestFunc = func(string) {}
	engineMock := new(clock.EngineMock)
	engineMock.StartFunc = func() {}
	engineMock.DoFunc = func() {}
	engineMock.WaitFunc = func() {}
	logger := log.NewLoggerMock()
	host := NewHost(serverMock, engineMock, engineMock, engineMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
