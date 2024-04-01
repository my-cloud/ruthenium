package address

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/common/infrastructure/log"
	"github.com/my-cloud/ruthenium/common/infrastructure/test"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	handler := NewHandler(logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	for _, method := range invalidHttpMethods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, urlTarget, nil)

			// Act
			handler.ServeHTTP(recorder, request)

			// Assert
			expectedStatusCode := 400
			test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
		})
	}
}

func Test_ServeHTTP_InvalidPublicKey_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	handler := NewHandler(logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?publicKey=invalidPublicKey", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidRequest_ReturnsAddress(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	handler := NewHandler(logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?publicKey=%s", urlTarget, test.PublicKey), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
