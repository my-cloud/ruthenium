package p2p

import (
	"context"
	"encoding/json"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
	"sync"
	"testing"
	"time"
)

func Test_HandleTargetsRequest_AddInvalidTargets_AddTargetsNotCalled(t *testing.T) {
	// Arrange
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.AddTargetsFunc = func([]string) {}
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), synchronizerMock, new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, _ := json.Marshal("")
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(synchronizerMock.AddTargetsCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

func Test_HandleTargetsRequest_AddValidTargets_AddTargetsCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.AddTargetsFunc = func([]string) { waitGroup.Done() }
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), synchronizerMock, new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	target := []string{"target"}
	data, _ := json.Marshal(target)
	req := gp2p.Data{}
	req.SetBytes(data)
	waitGroup.Add(1)

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(synchronizerMock.AddTargetsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleFirstBlockTimestampRequest_ValidRequest_FirstBlockTimestampCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.FirstBlockTimestampFunc = func() int64 { return 0 }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	req := gp2p.Data{}

	// Act
	_, _ = handler.HandleFirstBlockTimestampRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.FirstBlockTimestampCalls()) != 0
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleTransactionRequest_AddValidTransaction_AddTransactionCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	transactionsPoolMock := new(protocoltest.TransactionsPoolMock)
	transactionsPoolMock.AddTransactionFunc = func([]byte, string) { waitGroup.Done() }
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.HostTargetFunc = func() string { return "" }
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), synchronizerMock, transactionsPoolMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, _ := json.Marshal(network.TransactionRequest{})
	req := gp2p.Data{}
	req.SetBytes(data)
	waitGroup.Add(1)

	// Act
	_, _ = handler.HandleTransactionRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(transactionsPoolMock.AddTransactionCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleUtxosRequest_ValidUtxosRequest_UtxosByAddressCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.UtxosFunc = func(string) []byte { return nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), watchMock, logtest.NewLoggerMock())
	address := "address"
	data, _ := json.Marshal(&address)
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.HandleUtxosRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.UtxosCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleBlocksRequest_ValidBlocksRequest_LastBlocksCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func(uint64) []byte { return nil }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	var height uint64 = 0
	data, _ := json.Marshal(&height)
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.HandleBlocksRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.BlocksCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleTransactionsRequest_ValidTransactionsRequest_TransactionsCalled(t *testing.T) {
	// Arrange
	transactionsPoolMock := new(protocoltest.TransactionsPoolMock)
	transactionsPoolMock.TransactionsFunc = func() []byte { return nil }
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), transactionsPoolMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
	req := gp2p.Data{}

	// Act
	_, _ = handler.HandleTransactionsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsPoolMock.TransactionsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
