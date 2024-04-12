package p2p

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_HandleTargetsRequest_AddInvalidTargets_AddTargetsNotCalled(t *testing.T) {
	// Arrange
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.AddTargetsFunc = func([]string) {}
	handler := NewHandler(new(ledger.BlocksManagerMock), nil, sendersManagerMock, new(ledger.TransactionsManagerMock), new(ledger.UtxosManagerMock))
	targets := []string{"target"}
	marshalledTargets, _ := json.Marshal(targets)
	req := gp2p.Data{Bytes: marshalledTargets}

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(sendersManagerMock.AddTargetsCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

func Test_HandleTargetsRequest_AddValidTargets_AddTargetsCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.AddTargetsFunc = func([]string) { waitGroup.Done() }
	handler := NewHandler(new(ledger.BlocksManagerMock), nil, sendersManagerMock, new(ledger.TransactionsManagerMock), new(ledger.UtxosManagerMock))
	targets := []string{"target"}
	marshalledTargets, _ := json.Marshal(targets)
	req := gp2p.Data{Bytes: marshalledTargets}
	waitGroup.Add(1)

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(sendersManagerMock.AddTargetsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleSettingsRequest_ValidRequest_SettingsCalled(t *testing.T) {
	// Arrange
	expectedSettings := []byte{0}
	handler := NewHandler(new(ledger.BlocksManagerMock), expectedSettings, new(network.SendersManagerMock), new(ledger.TransactionsManagerMock), new(ledger.UtxosManagerMock))
	req := gp2p.Data{}

	// Act
	data, _ := handler.HandleSettingsRequest(context.TODO(), req)

	// Assert
	actualSettings := data.GetBytes()
	test.Assert(t, reflect.DeepEqual(expectedSettings, actualSettings), "Settings are not the expected ones.")
}

func Test_HandleFirstBlockTimestampRequest_ValidRequest_FirstBlockTimestampCalled(t *testing.T) {
	// Arrange
	blocksManagerMock := new(ledger.BlocksManagerMock)
	blocksManagerMock.FirstBlockTimestampFunc = func() int64 { return 0 }
	handler := NewHandler(blocksManagerMock, nil, new(network.SendersManagerMock), new(ledger.TransactionsManagerMock), new(ledger.UtxosManagerMock))
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
	transactionsManagerMock := new(ledger.TransactionsManagerMock)
	transactionsManagerMock.AddTransactionFunc = func(*protocol.Transaction, string, string) { waitGroup.Done() }
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.HostTargetFunc = func() string { return "" }
	handler := NewHandler(new(ledger.BlocksManagerMock), nil, sendersManagerMock, transactionsManagerMock, new(ledger.UtxosManagerMock))
	req := gp2p.Data{}
	transaction, _ := protocol.NewRewardTransaction("", false, 0, 0)
	transactionBytes, _ := json.Marshal(transaction)
	req.SetBytes(transactionBytes)
	waitGroup.Add(1)

	// Act
	_, _ = handler.HandleTransactionRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(transactionsManagerMock.AddTransactionCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleTransactionRequest_AddInvalidValidTransaction_AddTransactionNotCalled(t *testing.T) {
	// Arrange
	transactionsManagerMock := new(ledger.TransactionsManagerMock)
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.HostTargetFunc = func() string { return "" }
	handler := NewHandler(new(ledger.BlocksManagerMock), nil, sendersManagerMock, transactionsManagerMock, new(ledger.UtxosManagerMock))
	req := gp2p.Data{}

	// Act
	_, err := handler.HandleTransactionRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsManagerMock.AddTransactionCalls()) == 0
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
	test.Assert(t, isMethodCalled, "Method is called whereas it should not.")
}

func Test_HandleUtxosRequest_ValidUtxosRequest_UtxosByAddressCalled(t *testing.T) {
	// Arrange
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UtxosFunc = func(string) []*protocol.Utxo { return nil }
	handler := NewHandler(new(ledger.BlocksManagerMock), nil, new(network.SendersManagerMock), new(ledger.TransactionsManagerMock), utxosManagerMock)
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
	blocksManagerMock := new(ledger.BlocksManagerMock)
	blocksManagerMock.BlocksFunc = func(uint64) []*protocol.Block { return nil }
	handler := NewHandler(blocksManagerMock, nil, new(network.SendersManagerMock), new(ledger.TransactionsManagerMock), new(ledger.UtxosManagerMock))
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
	transactionsManagerMock := new(ledger.TransactionsManagerMock)
	transactionsManagerMock.TransactionsFunc = func() []*protocol.Transaction { return nil }
	handler := NewHandler(new(ledger.BlocksManagerMock), nil, new(network.SendersManagerMock), transactionsManagerMock, new(ledger.UtxosManagerMock))
	req := gp2p.Data{}

	// Act
	_, _ = handler.HandleTransactionsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsManagerMock.TransactionsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
