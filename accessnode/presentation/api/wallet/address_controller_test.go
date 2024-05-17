package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_GetWalletAddress_InvalidPublicKey_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	controller := NewAddressController(logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?publicKey=invalidPublicKey", "/"), nil)

	// Act
	controller.GetWalletAddress(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetWalletAddress_ValidRequest_ReturnsAddress(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	controller := NewAddressController(logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?publicKey=%s", "/", test.PublicKey), nil)

	// Act
	controller.GetWalletAddress(recorder, request)

	// Assert
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
