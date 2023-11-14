package amount

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/domain"
	"github.com/my-cloud/ruthenium/domain/observernode"
	"github.com/my-cloud/ruthenium/infrastructure/clock"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"github.com/my-cloud/ruthenium/infrastructure/network"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/infrastructure/test"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	watchMock := new(clock.WatchMock)
	settings := new(observernode.SettingsMock)
	handler := NewHandler(neighborMock, settings, watchMock, logger)
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

func Test_ServeHTTP_InvalidAddress_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	watchMock := new(clock.WatchMock)
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, urlTarget, nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_GetUtxosError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(clock.WatchMock)
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetUtxosCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ValidRequest_ReturnsAmount(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*domain.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetUtxosCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
