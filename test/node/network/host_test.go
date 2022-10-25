package network

import (
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/connection"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	servableMock := new(ServableMock)
	servableMock.ServeFunc = func() error { return nil }
	servableMock.SetHandleFunc = func(topic string, handler p2p.Handler) {}
	verifiableMock := new(VerifiableMock)
	verifiableMock.VerifyFunc = func() {}
	verifiableMock.StartVerificationFunc = func() {}
	validatableMock := new(ValidatableMock)
	bootableMock := new(BootableMock)
	bootableMock.StartFunc = func() {}
	watchMock := clock.NewWatch()
	sender := new(SenderMock)
	sender.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	senderProviderMock := new(SenderProviderMock)
	senderProviderMock.CreateSenderFunc = func(string, uint16, string) (connection.Sender, error) { return sender, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	neighborhood := network.NewNeighborhood("", 0, watchMock, senderProviderMock, configurationPath, logger)
	host := network.NewHost(servableMock, verifiableMock, validatableMock, bootableMock, neighborhood, watchMock, logger)

	// Act
	host.Run()

	// Assert
	isServerStarted := len(servableMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
