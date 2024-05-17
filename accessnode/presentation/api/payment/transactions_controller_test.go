package payment

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_GetTransactions_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return nil, errors.New("") }
	logger := log.NewLoggerMock()
	controller := NewTransactionsController(senderMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	// Act
	controller.GetTransactions(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactions_ValidRequest_NeighborMethodCalled(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return nil, nil }
	logger := log.NewLoggerMock()
	controller := NewTransactionsController(senderMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	// Act
	controller.GetTransactions(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
