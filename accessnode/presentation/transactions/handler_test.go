package transactions

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, urlTarget, nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			isNeighborMethodCalled := len(senderMock.GetTransactionsCalls()) != 0
			test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return nil, errors.New("") }
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, urlTarget, nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidRequest_NeighborMethodCalled(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return nil, nil }
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, urlTarget, nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
