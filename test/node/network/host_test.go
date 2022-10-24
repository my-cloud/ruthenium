package network

import (
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/connection"
	network_interface "github.com/my-cloud/ruthenium/src/api/node/network"
	"github.com/my-cloud/ruthenium/src/api/node/protocol"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test/node"
	"testing"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	watchMock := node.NewWatchMock()
	verifiableMock := new(VerifiableMock)
	verifiableMock.VerifyFunc = func(neighbors []network_interface.Requestable) {}
	verifiableMock.StartVerificationFunc = func(protocol.Validatable, network_interface.Synchronizable) {}
	validatableMock := new(ValidatableMock)
	bootableMock := new(BootableMock)
	bootableMock.StartFunc = func() {}
	senderProviderMock := new(SenderProviderMock)
	sender := new(SenderMock)
	sender.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	senderProviderMock.CreateSenderFunc = func(string, uint16, string) (connection.Sender, error) { return sender, nil }
	neighborhood := network.NewNeighborhood("", 0, watchMock, senderProviderMock, configurationPath, logger)
	servableMock := new(ServableMock)
	servableMock.ServeFunc = func() error { return nil }
	servableMock.SetHandleFunc = func(topic string, handler p2p.Handler) {}
	host := network.NewHost(servableMock, verifiableMock, validatableMock, bootableMock, neighborhood, watchMock, logger)
	host.Run()
}
