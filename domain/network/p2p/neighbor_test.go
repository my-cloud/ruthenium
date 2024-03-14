package p2p

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/infrastructure/test"
	"testing"
)

func Test_Target_NoError_ReturnTarget(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	expectedTargetValue := "ip:port"
	target, _ := NewTargetFromValue(expectedTargetValue)
	neighbor, _ := NewNeighbor(target, clientFactoryMock)

	// Act
	actualTargetString := neighbor.Target()

	// Assert
	test.Assert(t, actualTargetString == expectedTargetValue, fmt.Sprintf("Wrong target value. expected: %s actual: %s", expectedTargetValue, actualTargetString))
}

func Test_GetBlocks_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)
	var startingBlockHeight uint64 = 0

	// Act
	_, err := neighbor.GetBlocks(startingBlockHeight)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&startingBlockHeight)
	test.Assert(t, bytes.Equal(req, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetBlocks_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)
	var startingBlockHeight uint64 = 0

	// Act
	_, err := neighbor.GetBlocks(startingBlockHeight)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&startingBlockHeight)
	test.Assert(t, bytes.Equal(req, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetFirstBlockTimestamp_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	responseBytes, _ := json.Marshal(0)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return responseBytes, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)
	targets := []string{"target"}

	// Act
	err := neighbor.SendTargets(targets)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(targets)
	test.Assert(t, bytes.Equal(req, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_SendTargets_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

	// Act
	err := neighbor.SendTargets([]string{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal([]string{})
	test.Assert(t, bytes.Equal(req, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetSettings_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetSettings()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetSettings_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetSettings()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_AddTransaction_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)

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
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, nil }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)
	expectedAddress := "expected address"

	// Act
	_, err := neighbor.GetUtxos(expectedAddress)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&expectedAddress)
	test.Assert(t, bytes.Equal(req, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err == nil, "Error is not nil whereas it should be.")
}

func Test_GetUtxos_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(ClientMock)
	clientMock.SendFunc = func(string, []byte) ([]byte, error) { return []byte{}, errors.New("") }
	clientFactoryMock := new(ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (Client, error) { return clientMock, nil }
	neighbor, _ := NewNeighbor(new(Target), clientFactoryMock)
	expectedAddress := "expected address"

	// Act
	_, err := neighbor.GetUtxos(expectedAddress)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedRequestBytes, _ := json.Marshal(&expectedAddress)
	test.Assert(t, bytes.Equal(req, expectedRequestBytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}
