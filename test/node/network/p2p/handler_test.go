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

func Test_Handle_AddInvalidTargets_AddTargetsNotCalled(t *testing.T) {
	// Arrange
	synchronizerMock := new(networktest.SynchronizerMock)
	synchronizerMock.AddTargetsFunc = func([]network.TargetRequest) {}
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), synchronizerMock, new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal([]network.TargetRequest{{}})
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(synchronizerMock.AddTargetsCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

// TODO use wait group
//func Test_Handle_AddTargets_AddTargetsCalled(t *testing.T) {
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
//	_, _ = handler.Handle(context.TODO(), req)
//
//	// Assert
//	isMethodCalled := len(synchronizerMock.AddTargetsCalls()) == 1
//	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
//}

func Test_Handle_AddInvalidTransaction_AddTransactionNotCalled(t *testing.T) {
	// Arrange
	transactionsPoolMock := new(protocoltest.TransactionsPoolMock)
	transactionsPoolMock.AddTransactionFunc = func(*network.TransactionRequest, string) {}
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), transactionsPoolMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal(network.TransactionRequest{})
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsPoolMock.AddTransactionCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

// TODO use wait group
//func Test_Handle_AddTransaction_AddTransactionCalled(t *testing.T) {
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
//	_, _ = handler.Handle(context.TODO(), req)
//
//	// Assert
//	isMethodCalled := len(transactionsPoolMock.AddTransactionCalls()) == 1
//	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
//}

func Test_Handle_InvalidAmountRequest_CalculateTotalAmountNotCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return 0 }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal(network.AmountRequest{})
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.CalculateTotalAmountCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is not called whereas it should be.")
}

func Test_Handle_Amount_CalculateTotalAmountCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.CopyFunc = func() protocol.Blockchain { return blockchainMock }
	blockchainMock.CalculateTotalAmountFunc = func(int64, string) uint64 { return 0 }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), watchMock, logtest.NewLoggerMock())
	address := "address"
	data, err := json.Marshal(network.AmountRequest{Address: &address})
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.CalculateTotalAmountCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_Handle_Blocks_BlocksCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.BlocksFunc = func() []*network.BlockResponse { return nil }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal(p2p.GetBlocks)
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.BlocksCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_Handle_LastBlocks_LastBlocksCalled(t *testing.T) {
	// Arrange
	blockchainMock := new(protocoltest.BlockchainMock)
	blockchainMock.LastBlocksFunc = func(int64) []*network.BlockResponse { return nil }
	handler := p2p.NewHandler(blockchainMock, new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), new(clocktest.WatchMock), logtest.NewLoggerMock())
	var index int64 = 0
	data, err := json.Marshal(network.LastBlocksRequest{StartingBlockHeight: &index})
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(blockchainMock.LastBlocksCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_Handle_Transactions_TransactionsCalled(t *testing.T) {
	// Arrange
	transactionsPoolMock := new(protocoltest.TransactionsPoolMock)
	transactionsPoolMock.TransactionsFunc = func() []*network.TransactionResponse { return nil }
	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), transactionsPoolMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
	data, err := json.Marshal(p2p.GetTransactions)
	if err != nil {
		return
	}
	req := gp2p.Data{}
	req.SetBytes(data)

	// Act
	_, _ = handler.Handle(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsPoolMock.TransactionsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

// TODO use wait group
//func Test_Handle_StartValidation_StartCalled(t *testing.T) {
//	// Arrange
//	validationEngineMock := new(clocktest.EngineMock)
//	validationEngineMock.StartFunc = func() {}
//	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), validationEngineMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
//	data, err := json.Marshal(p2p.StartValidation)
//	if err != nil {
//		return
//	}
//	req := gp2p.Data{}
//	req.SetBytes(data)
//
//	// Act
//	_, _ = handler.Handle(context.TODO(), req)
//
//	// Assert
//	isMethodCalled := len(validationEngineMock.StartCalls()) == 1
//	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
//}

// TODO use wait group
//func Test_Handle_StopValidation_StopCalled(t *testing.T) {
//	// Arrange
//	validationEngineMock := new(clocktest.EngineMock)
//	validationEngineMock.StopFunc = func() {}
//	handler := p2p.NewHandler(new(protocoltest.BlockchainMock), new(networktest.SynchronizerMock), new(protocoltest.TransactionsPoolMock), validationEngineMock, new(clocktest.WatchMock), logtest.NewLoggerMock())
//	data, err := json.Marshal(p2p.StopValidation)
//	if err != nil {
//		return
//	}
//	req := gp2p.Data{}
//	req.SetBytes(data)
//
//	// Act
//	_, _ = handler.Handle(context.TODO(), req)
//
//	// Assert
//	isMethodCalled := len(validationEngineMock.StopCalls()) == 1
//	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
//}
