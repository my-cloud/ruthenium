package info

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction/status"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/ui/server/servertest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	watchMock := new(clocktest.WatchMock)
	settings := new(servertest.SettingsMock)
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	invalidHttpMethods := []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
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

func Test_ServeHTTP_UndecipherableTransaction_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(networktest.NeighborMock)
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	watchMock := new(clocktest.WatchMock)
	logger := logtest.NewLoggerMock()
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	data := ""
	marshalledData, _ := json.Marshal(data)
	body := bytes.NewReader(marshalledData)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
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
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

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
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_GetBlocksError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*verification.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1 && len(neighborMock.GetBlocksCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_GetTransactionsError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*verification.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*verification.Block{verification.NewBlock(0, [32]byte{}, nil, nil, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return nil, errors.New("") }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1 && len(neighborMock.GetBlocksCalls()) == 1 && len(neighborMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_ServeHTTP_TransactionNotFound_ReturnsRejected(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*verification.Utxo{})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*verification.Block{verification.NewBlock(0, [32]byte{}, nil, nil, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledEmptyArray, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(servertest.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseInParticlesFunc = func() uint64 { return 0 }
	settings.IncomeLimitInParticlesFunc = func() uint64 { return 0 }
	settings.ParticlesPerTokenFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := status.NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	transactionRequest, _ := verification.NewRewardTransaction("", false, 0, 0)
	marshalledTransaction, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(marshalledTransaction)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1 && len(neighborMock.GetBlocksCalls()) == 1 && len(neighborMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}
