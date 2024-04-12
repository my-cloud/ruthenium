package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	senderMock.TargetFunc = func() string { return "" }
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, urlTarget, nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) != 0
			test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_UndecipherableTransaction_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	marshalledData, _ := json.Marshal("")
	body := bytes.NewReader(marshalledData)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	target := "0.0.0.0:0"
	senderMock.TargetFunc = func() string { return target }
	senderMock.AddTransactionFunc = func([]byte) error { return errors.New("") }
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	transactionRequest, _ := protocol.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidTransaction_NeighborMethodCalled(t *testing.T) {
	// Arrange
	senderMock := new(network.SenderMock)
	target := "0.0.0.0:0"
	senderMock.TargetFunc = func() string { return target }
	senderMock.AddTransactionFunc = func([]byte) error { return nil }
	logger := log.NewLoggerMock()
	handler := NewHandler(senderMock, logger)
	transactionRequest, _ := protocol.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 201
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
