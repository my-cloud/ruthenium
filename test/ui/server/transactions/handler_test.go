package transaction

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server/transactions"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"net/http"
	"net/http/httptest"
	"testing"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	logger := logtest.NewLoggerMock()
	handler := transactions.NewHandler(neighborMock, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, urlTarget, nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			isNeighborMethodCalled := len(neighborMock.GetTransactionsCalls()) != 0
			test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetTransactionsFunc = func() ([]network.TransactionResponse, error) { return nil, errors.New("") }
	logger := logtest.NewLoggerMock()
	handler := transactions.NewHandler(neighborMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", urlTarget, nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetTransactionsCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidRequest_NeighborMethodCalled(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetTransactionsFunc = func() ([]network.TransactionResponse, error) { return nil, nil }
	logger := logtest.NewLoggerMock()
	handler := transactions.NewHandler(neighborMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", urlTarget, nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetTransactionsCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
