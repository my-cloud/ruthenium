package p2p

import (
	"testing"

	"github.com/my-cloud/ruthenium/common/application"
	"github.com/my-cloud/ruthenium/common/infrastructure/log"
	"github.com/my-cloud/ruthenium/common/infrastructure/test"
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
	engineMock := new(application.PulserMock)
	engineMock.StartFunc = func() {}
	engineMock.PulseFunc = func() {}
	logger := log.NewLoggerMock()
	node := NewNode(serverMock, logger, engineMock)

	// Act
	_ = node.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
