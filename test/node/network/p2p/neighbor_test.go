package p2p

import (
	"bytes"
	"encoding/json"
	"errors"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
)

func Test_GetBlocks_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, _ = neighbor.GetBlocks()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(p2p.GetBlocks)
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
}

func Test_GetBlocks_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)

	// Act
	_, err := neighbor.GetBlocks()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(p2p.GetBlocks)
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetLastBlocks_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	var startingBlockHeight uint64 = 0

	// Act
	_, _ = neighbor.GetLastBlocks(startingBlockHeight)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(network.LastBlocksRequest{StartingBlockHeight: &startingBlockHeight})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
}

func Test_GetLastBlocks_Error_ReturnsError(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, errors.New("") }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	var startingBlockHeight uint64 = 0

	// Act
	_, err := neighbor.GetLastBlocks(startingBlockHeight)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(network.LastBlocksRequest{StartingBlockHeight: &startingBlockHeight})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
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

	// Act
	_ = neighbor.SendTargets([]network.TargetRequest{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal([]network.TargetRequest{})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
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
	err := neighbor.SendTargets([]network.TargetRequest{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal([]network.TargetRequest{})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
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
	_ = neighbor.AddTransaction(network.TransactionRequest{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(network.TransactionRequest{})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
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
	err := neighbor.AddTransaction(network.TransactionRequest{})

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(network.TransactionRequest{})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
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
	_, _ = neighbor.GetTransactions()

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(p2p.GetTransactions)
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
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
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(p2p.GetTransactions)
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_GetAmount_NoError_ClientCalled(t *testing.T) {
	// Arrange
	clientMock := new(p2ptest.ClientMock)
	clientMock.SetSettingsFunc = func(*gp2p.ClientSettings) {}
	clientMock.SendFunc = func(string, gp2p.Data) (gp2p.Data, error) { return gp2p.Data{}, nil }
	clientFactoryMock := new(p2ptest.ClientFactoryMock)
	clientFactoryMock.CreateClientFunc = func(string, string) (p2p.Client, error) { return clientMock, nil }
	neighbor, _ := p2p.NewNeighbor(new(p2p.Target), clientFactoryMock)
	expectedAddress := "expected address"

	// Act
	_, _ = neighbor.GetUtxos(expectedAddress)

	// Assert
	sendCalls := clientMock.SendCalls()
	isSendCalledOnce := len(sendCalls) == 1
	test.Assert(t, isSendCalledOnce, "Client is not called a single time whereas it should be.")
	req := sendCalls[0].Req
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(network.UtxosRequest{Address: &expectedAddress})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
}

func Test_GetAmount_Error_ReturnsError(t *testing.T) {
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
	expectedReq := gp2p.Data{}
	data, _ := json.Marshal(network.UtxosRequest{Address: &expectedAddress})
	expectedReq.SetBytes(data)
	test.Assert(t, bytes.Equal(req.Bytes, expectedReq.Bytes), "Client is not called with the good parameter.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}
