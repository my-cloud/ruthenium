package payment

import (
	"context"
	"encoding/json"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"sync"
	"testing"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_HandleTransactionRequest_AddValidTransaction_AddTransactionCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	transactionsManagerMock := new(application.TransactionsManagerMock)
	transactionsManagerMock.AddTransactionFunc = func(*ledger.Transaction, string, string) { waitGroup.Done() }
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.HostTargetFunc = func() string { return "" }
	controller := NewTransactionsController(sendersManagerMock, transactionsManagerMock)
	req := gp2p.Data{}
	transaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	transactionBytes, _ := json.Marshal(transaction)
	req.SetBytes(transactionBytes)
	waitGroup.Add(1)

	// Act
	_, _ = controller.HandleTransactionRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(transactionsManagerMock.AddTransactionCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleTransactionRequest_AddInvalidValidTransaction_AddTransactionNotCalled(t *testing.T) {
	// Arrange
	transactionsManagerMock := new(application.TransactionsManagerMock)
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.HostTargetFunc = func() string { return "" }
	controller := NewTransactionsController(sendersManagerMock, transactionsManagerMock)
	req := gp2p.Data{}

	// Act
	_, err := controller.HandleTransactionRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsManagerMock.AddTransactionCalls()) == 0
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
	test.Assert(t, isMethodCalled, "Method is called whereas it should not.")
}

func Test_HandleTransactionsRequest_ValidTransactionsRequest_TransactionsCalled(t *testing.T) {
	// Arrange
	transactionsManagerMock := new(application.TransactionsManagerMock)
	transactionsManagerMock.TransactionsFunc = func() []*ledger.Transaction { return nil }
	controller := NewTransactionsController(nil, transactionsManagerMock)
	req := gp2p.Data{}

	// Act
	_, _ = controller.HandleTransactionsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(transactionsManagerMock.TransactionsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
