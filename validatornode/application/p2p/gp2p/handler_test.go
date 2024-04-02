package gp2p

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"testing"
	"time"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/p2p"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_HandleTargetsRequest_AddInvalidTargets_AddTargetsNotCalled(t *testing.T) {
	// Arrange
	neighborsManagerMock := new(p2p.NeighborsManagerMock)
	neighborsManagerMock.AddTargetsFunc = func([]string) {}
	handler := NewHandler(new(protocol.BlocksManagerMock), nil, neighborsManagerMock, new(protocol.TransactionsManagerMock), new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	targets := []string{"target"}
	marshalledTargets, _ := json.Marshal(targets)
	req := gp2p.Data{Bytes: marshalledTargets}

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(neighborsManagerMock.AddTargetsCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

func Test_HandleTargetsRequest_AddValidTargets_AddTargetsCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	neighborsManagerMock := new(p2p.NeighborsManagerMock)
	neighborsManagerMock.AddTargetsFunc = func([]string) { waitGroup.Done() }
	handler := NewHandler(new(protocol.BlocksManagerMock), nil, neighborsManagerMock, new(protocol.TransactionsManagerMock), new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	targets := []string{"target"}
	marshalledTargets, _ := json.Marshal(targets)
	req := gp2p.Data{Bytes: marshalledTargets}
	waitGroup.Add(1)

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(neighborsManagerMock.AddTargetsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleSettingsRequest_ValidRequest_SettingsCalled(t *testing.T) {
	// Arrange
	expectedSettings := []byte{0}
	handler := NewHandler(new(protocol.BlocksManagerMock), expectedSettings, new(p2p.NeighborsManagerMock), new(protocol.TransactionsManagerMock), new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	req := gp2p.Data{}

	// Act
	data, _ := handler.HandleSettingsRequest(context.TODO(), req)

	// Assert
	actualSettings := data.GetBytes()
	test.Assert(t, reflect.DeepEqual(expectedSettings, actualSettings), "Settings are not the expected ones.")
}

func Test_HandleFirstBlockTimestampRequest_ValidRequest_FirstBlockTimestampCalled(t *testing.T) {
	// Arrange
	blocksManagerMock := new(protocol.BlocksManagerMock)
	blocksManagerMock.FirstBlockTimestampFunc = func() int64 { return 0 }
	handler := NewHandler(blocksManagerMock, nil, new(p2p.NeighborsManagerMock), new(protocol.TransactionsManagerMock), new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	req := gp2p.Data{}

	// Act
	_, _ = handler.HandleFirstBlockTimestampRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blocksManagerMock.FirstBlockTimestampCalls()) != 0
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleTransactionRequest_AddValidTransaction_AddTransactionCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	transactionsManagerMock := new(protocol.TransactionsManagerMock)
	transactionsManagerMock.AddTransactionFunc = func([]byte, string) { waitGroup.Done() }
	neighborsManagerMock := new(p2p.NeighborsManagerMock)
	neighborsManagerMock.HostTargetFunc = func() string { return "" }
	handler := NewHandler(new(protocol.BlocksManagerMock), nil, neighborsManagerMock, transactionsManagerMock, new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	req := gp2p.Data{}
	waitGroup.Add(1)

	// Act
	_, _ = handler.HandleTransactionRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(transactionsManagerMock.AddTransactionCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleUtxosRequest_ValidUtxosRequest_UtxosByAddressCalled(t *testing.T) {
	// Arrange
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.UtxosFunc = func(string) []byte { return nil }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	handler := NewHandler(new(protocol.BlocksManagerMock), nil, new(p2p.NeighborsManagerMock), new(protocol.TransactionsManagerMock), utxosManagerMock, watchMock, log.NewLoggerMock())
	address := "address"
	marshalledAddress, _ := json.Marshal(&address)
	req := gp2p.Data{Bytes: marshalledAddress}

	// Act
	_, _ = handler.HandleUtxosRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(utxosManagerMock.UtxosCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleBlocksRequest_ValidBlocksRequest_LastBlocksCalled(t *testing.T) {
	// Arrange
	blocksManagerMock := new(protocol.BlocksManagerMock)
	blocksManagerMock.BlocksFunc = func(uint64) []byte { return nil }
	handler := NewHandler(blocksManagerMock, nil, new(p2p.NeighborsManagerMock), new(protocol.TransactionsManagerMock), new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	var height uint64 = 0
	marshalledHeight, _ := json.Marshal(&height)
	req := gp2p.Data{Bytes: marshalledHeight}

	// Act
	_, _ = handler.HandleBlocksRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blocksManagerMock.BlocksCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleTransactionsRequest_ValidTransactionsRequest_TransactionsCalled(t *testing.T) {
	// Arrange
	transactionsManagerMock := new(protocol.TransactionsManagerMock)
	transactionsManagerMock.TransactionsFunc = func() []byte { return nil }
	handler := NewHandler(new(protocol.BlocksManagerMock), nil, new(p2p.NeighborsManagerMock), transactionsManagerMock, new(protocol.UtxosManagerMock), new(protocol.TimeProviderMock), log.NewLoggerMock())
	req := gp2p.Data{}

	// Act
	_, _ = handler.HandleTransactionsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsManagerMock.TransactionsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
