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
	"testing"
	"time"
)

func Test_HandleTargetsRequest_AddInvalidTargets_AddTargetsNotCalled(t *testing.T) {
	// Arrange
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.AddTargetsFunc = func([]string) {}
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), synchronizerMock, new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal("")
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.HandleTargetsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(synchronizerMock.AddTargetsCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

// TODO use wait group
//func Test_HandleTargetsRequest_AddValidTargets_AddTargetsCalled(t *testing.T) {
//	// Arrange
//	synchronizerMock := new(networktest.SynchronizerMock)
//	synchronizerMock.AddTargetsFunc = func([]network.TargetRequest) {}
//	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), synchronizerMock, new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
//	target := "target"
//	data, err := json.Marshal([]network.TargetRequest{{Target: &target}})
//	if err != nil {
//		return
//	}
//	req := gp2p.Data{}
//	req.SetBytes(data)
//
//	// Act
//	_, _ = handler.HandleTargetsRequest(context.TODO(), req)
//
//	// Assert
//	isMethodCalled := len(synchronizerMock.AddTargetsCalls()) == 1
//	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
//}

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

func Test_HandleTransactionRequest_AddInvalidTransaction_AddTransactionNotCalled(t *testing.T) {
	// Arrange
	transactionsPoolMock := new(protocoltest.TransactionsPoolMock)
	transactionsPoolMock.AddTransactionFunc = func([]byte, string) {}
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), transactionsPoolMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal(network.TransactionRequest{})
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.HandleTransactionRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsPoolMock.AddTransactionCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

// TODO use wait group
//func Test_HandleTransactionRequest_AddValidTransaction_AddTransactionCalled(t *testing.T) {
//	// Arrange
//	transactionsPoolMock := new(protocoltest.TransactionsPoolMock)
//	transactionsPoolMock.AddTransactionFunc = func(*network.TransactionRequest, string) {}
//	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), transactionsPoolMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
//	privateKey, _ := encryption.DecodePrivateKey(test.PrivateKey)
//	transaction := server.NewTransaction(0, "A", "B", encryption.NewPublicKey(privateKey), 0, 0)
//	_ = transaction.Sign(privateKey)
//	data, err := json.Marshal(transaction.GetRequest())
//	if err != nil {
//		return
//	}
//	req := gp2p.Data{}
//	req.SetBytes(data)
//
//	// Act
//	_, _ = handler.HandleTransactionRequest(context.TODO(), req)
//
//	// Assert
//	isMethodCalled := len(transactionsPoolMock.AddTransactionCalls()) == 1
//	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
//}

func Test_HandleUtxosRequest_ValidUtxosRequest_UtxosByAddressCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.UtxosFunc = func(string) []byte { return nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), watchMock, logtest.NewLoggerMock())
	address := "address"
	data, err := json.Marshal(&address)
	if err != nil {
		return
	}
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
	data, err := json.Marshal(&height)
	if err != nil {
		return
	}
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
