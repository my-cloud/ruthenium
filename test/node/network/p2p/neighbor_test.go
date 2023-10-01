package p2p

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
)

func Test_Target_NoError_ReturnTarget(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	expectedTargetValue := "ip:port"
	target, _ := p2p.NewTargetFromValue(expectedTargetValue)
	neighbor, _ := p2p.NewNeighbor(target, clientFactoryMock)

	// Act
	actualTargetString := neighbor.Target()

	// Assert
	test.Assert(t, actualTargetString == expectedTargetValue, fmt.Sprintf("Wrong target value. expected: %s actual: %s", expectedTargetValue, actualTargetString))
}

func Test_GetBlocks_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	var startingBlockHeight uint64 = 0

	// Act
	_, err := neighbor.GetBlocks(startingBlockHeight)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&startingBlockHeight)
	test.Assert(t, bytes.Equal(req.Bytes, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetBlocks_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	var startingBlockHeight uint64 = 0

	// Act
	_, err := neighbor.GetBlocks(startingBlockHeight)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&startingBlockHeight)
	test.Assert(t, bytes.Equal(req.Bytes, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetFirstBlockTimestamp_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	responseBytes, _ := json.Marshal(0)
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{Bytes: responseBytes}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetFirstBlockTimestamp()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetFirstBlockTimestamp_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetFirstBlockTimestamp()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetFirstBlockTimestamp_UnmarshalError_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetFirstBlockTimestamp()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_SendTargets_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	targets := []string{"target"}

	// Act
	err := neighbor.SendTargets(targets)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(targets)
	test.Assert(t, bytes.Equal(req.Bytes, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_SendTargets_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	err := neighbor.SendTargets([]string{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal([]string{})
	test.Assert(t, bytes.Equal(req.Bytes, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_AddTransaction_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	err := neighbor.AddTransaction([]byte{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_AddTransaction_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	err := neighbor.AddTransaction([]byte{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetTransactions_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetTransactions()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetTransactions_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetTransactions()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetUtxos_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	expectedAddress := "expected address"

	// Act
	_, err := neighbor.GetUtxos(expectedAddress)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&expectedAddress)
	test.Assert(t, bytes.Equal(req.Bytes, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetUtxos_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	expectedAddress := "expected address"

	// Act
	_, err := neighbor.GetUtxos(expectedAddress)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&expectedAddress)
	test.Assert(t, bytes.Equal(req.Bytes, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}
