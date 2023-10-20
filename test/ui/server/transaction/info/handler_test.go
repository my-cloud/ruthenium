package info

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction/info"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/ui/server/servertest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*verification.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*verification.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
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
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)

	inputInfo1 := verification.NewInputInfo(0, "")
	inputInfo2 := verification.NewInputInfo(1, "")
	inputInfo3 := verification.NewInputInfo(2, "")
	output1 := verification.NewOutput("", false, 1)
	output2 := verification.NewOutput("", false, 2)
	output3 := verification.NewOutput("", false, 0)
	utxos := []*verification.Utxo{
		verification.NewUtxo(inputInfo1, output1, 1),
		verification.NewUtxo(inputInfo2, output2, 1),
		verification.NewUtxo(inputInfo3, output3, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=1&consolidation=false", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo *info.TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 1
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}

func Test_ServeHTTP_ConsolidationRequired_ReturnsAllUtxos(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	inputInfo1 := verification.NewInputInfo(0, "")
	inputInfo2 := verification.NewInputInfo(2, "")
	output1 := verification.NewOutput("", false, 1)
	output2 := verification.NewOutput("", false, 2)
	utxos := []*verification.Utxo{
		verification.NewUtxo(inputInfo1, output1, 1),
		verification.NewUtxo(inputInfo2, output2, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := info.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=1&consolidation=true", urlTarget), nil)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo info.TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 2
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}
