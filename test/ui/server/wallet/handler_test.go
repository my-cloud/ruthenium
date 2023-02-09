package wallet

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
)

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	handler := wallet.NewHandler(encryption.NewEmptyWallet(), logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, "/wallet", nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_Get_ReturnsOk(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	handler := wallet.NewHandler(encryption.NewEmptyWallet(), logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/wallet", nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
