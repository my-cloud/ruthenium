package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/domain/clock"
	"github.com/my-cloud/ruthenium/domain/ledger"
	"github.com/my-cloud/ruthenium/domain/network"
	"github.com/my-cloud/ruthenium/domain/observernode"
	"github.com/my-cloud/ruthenium/domain/observernode/transaction/output"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"github.com/my-cloud/ruthenium/infrastructure/test"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func Test_ServeHTTP_UndecipherableUtxo_BadRequest(t *testing.T) {
	// Arrange
	neighborMock := new(network.NeighborMock)
	settings := new(observernode.SettingsMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	watchMock := new(clock.WatchMock)
	logger := log.NewLoggerMock()
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	data := ""
	marshalledData, _ := json.Marshal(data)
	body := bytes.NewReader(marshalledData)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	isNeighborMethodCalled := len(neighborMock.GetUtxosCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
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
	utxo := ledger.NewUtxo(&ledger.InputInfo{}, &ledger.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
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
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
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
	utxo := ledger.NewUtxo(&ledger.InputInfo{}, &ledger.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
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
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return nil, errors.New("") }
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
	utxo := ledger.NewUtxo(&ledger.InputInfo{}, &ledger.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
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
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return nil, errors.New("") }
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
	utxo := ledger.NewUtxo(&ledger.InputInfo{}, &ledger.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
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
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledEmptyArray, nil }
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
	utxo := ledger.NewUtxo(&ledger.InputInfo{}, &ledger.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1 && len(neighborMock.GetBlocksCalls()) == 1 && len(neighborMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "rejected"
	response := recorder.Body.Bytes()
	var progressInfo *output.ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}

func Test_ServeHTTP_UtxoFound_ReturnsConfirmed(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	transaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	transactionId := transaction.Id()
	inputInfo := ledger.NewInputInfo(0, transactionId)
	utxo := ledger.NewUtxo(inputInfo, &ledger.Output{}, 0)
	marshalledUtxos, _ := json.Marshal([]*ledger.Utxo{utxo})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
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
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "confirmed"
	response := recorder.Body.Bytes()
	var progressInfo *output.ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}

func Test_ServeHTTP_ValidatedTransactionFound_ReturnsValidated(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	transaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, []*ledger.Transaction{transaction})}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
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
	outputIndex := 0
	utxo := ledger.NewUtxo(ledger.NewInputInfo(uint16(outputIndex), transaction.Id()), transaction.Outputs()[outputIndex], transaction.Timestamp())
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1 && len(neighborMock.GetBlocksCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "validated"
	response := recorder.Body.Bytes()
	var progressInfo *output.ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}

func Test_ServeHTTP_PendingTransactionFound_ReturnsSent(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(network.NeighborMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	transaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	transactions := []*ledger.Transaction{transaction}
	marshalledTransactions, _ := json.Marshal(transactions)
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledTransactions, nil }
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
	outputIndex := 0
	utxo := ledger.NewUtxo(ledger.NewInputInfo(uint16(outputIndex), transaction.Id()), transaction.Outputs()[outputIndex], transaction.Timestamp())
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, urlTarget, body)

	// Act
	handler.ServeHTTP(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(neighborMock.GetUtxosCalls()) == 1 && len(neighborMock.GetFirstBlockTimestampCalls()) == 1 && len(neighborMock.GetBlocksCalls()) == 1 && len(neighborMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "sent"
	response := recorder.Body.Bytes()
	var progressInfo *output.ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}
