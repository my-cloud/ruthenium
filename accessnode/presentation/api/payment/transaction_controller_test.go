package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_PostTransaction_UndecipherableTransaction_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	logger := log.NewLoggerMock()
	controller := NewTransactionController(senderMock, logger)
	marshalledData, _ := json.Marshal("")
	body := bytes.NewReader(marshalledData)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", body)

	// Act
	controller.PostTransaction(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_PostTransaction_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	target := "0.0.0.0:0"
	senderMock.TargetFunc = func() string { return target }
	senderMock.AddTransactionFunc = func([]byte) error { return errors.New("") }
	logger := log.NewLoggerMock()
	controller := NewTransactionController(senderMock, logger)
	transactionRequest, _ := ledger.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", body)

	// Act
	controller.PostTransaction(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_PostTransaction_ValidTransaction_NeighborMethodCalled(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	target := "0.0.0.0:0"
	senderMock.TargetFunc = func() string { return target }
	senderMock.AddTransactionFunc = func([]byte) error { return nil }
	logger := log.NewLoggerMock()
	controller := NewTransactionController(senderMock, logger)
	transactionRequest, _ := ledger.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", body)

	// Act
	controller.PostTransaction(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 201
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
