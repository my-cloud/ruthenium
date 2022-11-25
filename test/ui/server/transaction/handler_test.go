package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/network/networkmock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networkmock.NeighborMock)
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, "/transaction", nil)

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
	neighborMock := new(networkmock.NeighborMock)
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	transactionRequest := ""
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_InvalidTransaction_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networkmock.NeighborMock)
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	transactionRequest := newTransactionRequest("", "", "", "", "")
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_InvalidPrivateKey_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networkmock.NeighborMock)
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	transactionRequest := newTransactionRequest(
		"InvalidPrivateKey",
		test.Address,
		"A",
		test.PublicKey,
		"1",
	)
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_InvalidTransactionValue_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networkmock.NeighborMock)
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	transactionRequest := newTransactionRequest(
		test.PrivateKey,
		test.Address,
		"A",
		test.PublicKey,
		"InvalidTransactionValue",
	)
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

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
	neighborMock := new(networkmock.NeighborMock)
	neighborMock.AddTransactionFunc = func(network.TransactionRequest) error { return errors.New("") }
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	transactionRequest := newTransactionRequest(
		test.PrivateKey,
		test.Address,
		"A",
		test.PublicKey,
		"1",
	)
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

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
	neighborMock := new(networkmock.NeighborMock)
	neighborMock.AddTransactionFunc = func(network.TransactionRequest) error { return nil }
	logger := log.NewLogger(log.Fatal)
	handler := transaction.NewHandler(neighborMock, 1, logger)
	transactionRequest := newTransactionRequest(
		test.PrivateKey,
		test.Address,
		"A",
		test.PublicKey,
		"1",
	)
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func newTransactionRequest(senderPrivateKey, senderAddress, recipientAddress, senderPublicKey, value string) *server.TransactionRequest {
	return &server.TransactionRequest{
		SenderPrivateKey: &senderPrivateKey,
		SenderAddress:    &senderAddress,
		RecipientAddress: &recipientAddress,
		SenderPublicKey:  &senderPublicKey,
		Value:            &value,
	}
}
