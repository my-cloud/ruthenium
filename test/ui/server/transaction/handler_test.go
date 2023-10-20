package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/src/ui/server/transaction"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	neighborMock.TargetFunc = func() string { return "" }
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, urlTarget, nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) != 0
			test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_UndecipherableTransaction_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	marshalledData, _ := json.Marshal("")
	body := bytes.NewReader(marshalledData)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	target := "0.0.0.0:0"
	neighborMock.TargetFunc = func() string { return target }
	neighborMock.AddTransactionFunc = func([]byte) error { return errors.New("") }
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidTransaction_NeighborMethodCalled(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	target := "0.0.0.0:0"
	neighborMock.TargetFunc = func() string { return target }
	neighborMock.AddTransactionFunc = func([]byte) error { return nil }
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 201
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
