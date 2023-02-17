package address

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/amount"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
)

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	handler := amount.NewHandler(neighborMock, 1, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, "/wallet/address", nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_InvalidAddress_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetAmountFunc = func(network.AmountRequest) (*network.AmountResponse, error) { return nil, errors.New("") }
	handler := amount.NewHandler(neighborMock, 1, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/wallet/amount", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetAmountCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_NodeError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetAmountFunc = func(network.AmountRequest) (*network.AmountResponse, error) { return nil, errors.New("") }
	handler := amount.NewHandler(neighborMock, 1, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/wallet/amount?address=address", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetAmountCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidRequest_ReturnsAmount(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetAmountFunc = func(network.AmountRequest) (*network.AmountResponse, error) {
		return &network.AmountResponse{Amount: 0}, nil
	}
	handler := amount.NewHandler(neighborMock, 1, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/wallet/amount?address=address", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetAmountCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
