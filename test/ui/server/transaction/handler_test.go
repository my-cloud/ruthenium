package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
)

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
	neighborMock := new(networktest.NeighborMock)
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
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
	neighborMock := new(networktest.NeighborMock)
	neighborMock.TargetFunc = func() string { return "0.0.0.0:0" }
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	transactionRequest := protocoltest.NewTransactionRequest(0, "", "", 0, 0)
	b, _ := json.Marshal(&transactionRequest)
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
	neighborMock := new(networktest.NeighborMock)
	neighborMock.TargetFunc = func() string { return "0.0.0.0:0" }
	neighborMock.AddTransactionFunc = func(network.TransactionRequest) error { return errors.New("") }
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey, _ := encryption.DecodePublicKey(test.PublicKey)
	transactionRequest := protocoltest.NewSignedTransactionRequest(0, "RecipientAddress", "SenderAddress", privateKey, publicKey, 0, 0)
	b, _ := json.Marshal(&transactionRequest)
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
	neighborMock := new(networktest.NeighborMock)
	neighborMock.TargetFunc = func() string { return "0.0.0.0:0" }
	neighborMock.AddTransactionFunc = func(network.TransactionRequest) error { return nil }
	logger := logtest.NewLoggerMock()
	handler := transaction.NewHandler(neighborMock, logger)
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey, _ := encryption.DecodePublicKey(test.PublicKey)
	transactionRequest := protocoltest.NewSignedTransactionRequest(0, "RecipientAddress", "SenderAddress", privateKey, publicKey, 0, 0)
	b, _ := json.Marshal(&transactionRequest)
	body := bytes.NewReader(b)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/transaction", body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 201
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
