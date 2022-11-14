package network

import (
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node"
	"testing"
)

func Test_Run_NoError_ServerStarted(t *testing.T) {
	// Arrange
	serverMock := new(ServerMock)
	serverMock.ServeFunc = func() error { return nil }
	serverMock.SetHandleFunc = func(topic string, handler p2p.Handler) {}
	blockchainMock := new(node.BlockchainMock)
	blockchainMock.VerifyFunc = func() {}
	blockchainMock.StartVerificationFunc = func() {}
	transactionsPoolMock := new(TransactionsPoolMock)
	validatorMock := new(ValidatorMock)
	validatorMock.StartValidationFunc = func() {}
	watchMock := clock.NewWatch()
	client := new(ClientMock)
	client.SendFunc = func(string, p2p.Data) (p2p.Data, error) { return p2p.Data{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, uint16, string) (network.Client, error) { return client, nil }
	configurationPath := "../../"
	logger := log.NewLogger(log.Fatal)
	synchronizer, err := network.NewSynchronizer(0, watchMock, clientFactoryMock, configurationPath, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	host := network.NewHost(serverMock, blockchainMock, transactionsPoolMock, validatorMock, synchronizer, watchMock, logger)

	// Act
	_ = host.Run()

	// Assert
	isServerStarted := len(serverMock.ServeCalls()) == 1
	test.Assert(t, isServerStarted, "Server is not started whereas it should be.")
}
