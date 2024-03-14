package info

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/domain/clock"
	"github.com/my-cloud/ruthenium/domain/ledger"
	"github.com/my-cloud/ruthenium/domain/network"
	"github.com/my-cloud/ruthenium/domain/observernode"
	"github.com/my-cloud/ruthenium/infrastructure/log"
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

func Test_ServeHTTP_InvalidValue_ReturnsBadRequest(t *testing.T) {
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
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_IsRegisteredNotProvided_ReturnsBadRequest(t *testing.T) {
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
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0", urlTarget), nil)

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
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetUtxosCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_GetFirstBlockTimestampError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*ledger.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
	watchMock := new(clock.WatchMock)
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_InsufficientWalletBalance_ReturnsMethodNotAllowed(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*ledger.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 405
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_ConsolidationNotRequired_ReturnsSomeUtxos(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)

	inputInfo1 := ledger.NewInputInfo(0, "")
	inputInfo2 := ledger.NewInputInfo(1, "")
	inputInfo3 := ledger.NewInputInfo(2, "")
	output1 := ledger.NewOutput("", false, 1)
	output2 := ledger.NewOutput("", false, 2)
	output3 := ledger.NewOutput("", false, 0)
	utxos := []*ledger.Utxo{
		ledger.NewUtxo(inputInfo1, output1, 1),
		ledger.NewUtxo(inputInfo2, output2, 1),
		ledger.NewUtxo(inputInfo3, output3, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=1&consolidation=false", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo *TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 1
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}

func Test_ServeHTTP_ConsolidationRequired_ReturnsAllUtxos(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	inputInfo1 := ledger.NewInputInfo(0, "")
	inputInfo2 := ledger.NewInputInfo(2, "")
	output1 := ledger.NewOutput("", false, 1)
	output2 := ledger.NewOutput("", false, 2)
	utxos := []*ledger.Utxo{
		ledger.NewUtxo(inputInfo1, output1, 1),
		ledger.NewUtxo(inputInfo2, output2, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=1&consolidation=true", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 2
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}
