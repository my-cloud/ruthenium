package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction/output"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"github.com/my-cloud/ruthenium/validatornode/presentation"
)

const urlTarget = "/url-target"

func Test_ServeHTTP_InvalidHttpMethod_BadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	watchMock := new(protocol.TimeProviderMock)
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	settings := new(transaction.SettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	watchMock := new(protocol.TimeProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(protocol.TimeProviderMock)
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return nil, errors.New("") }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledEmptyArray, nil }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	rewardTransaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	transactionId := rewardTransaction.Id()
	inputInfo := ledger.NewInputInfo(0, transactionId)
	utxo := ledger.NewUtxo(inputInfo, &ledger.Output{}, 0)
	marshalledUtxos, _ := json.Marshal([]*ledger.Utxo{utxo})
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
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
	neighborMock := new(presentation.NeighborCallerMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	rewardTransaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, []*ledger.Transaction{rewardTransaction})}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	outputIndex := 0
	utxo := ledger.NewUtxo(ledger.NewInputInfo(uint16(outputIndex), rewardTransaction.Id()), rewardTransaction.Outputs()[outputIndex], rewardTransaction.Timestamp())
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
	neighborMock := new(presentation.NeighborCallerMock)
	marshalledEmptyArray := []byte{91, 93}
	neighborMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	neighborMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*ledger.Block{ledger.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	rewardTransaction, _ := ledger.NewRewardTransaction("", false, 0, 0)
	transactions := []*ledger.Transaction{rewardTransaction}
	marshalledTransactions, _ := json.Marshal(transactions)
	neighborMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledTransactions, nil }
	watchMock := new(protocol.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(transaction.SettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	handler := NewHandler(neighborMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	outputIndex := 0
	utxo := ledger.NewUtxo(ledger.NewInputInfo(uint16(outputIndex), rewardTransaction.Id()), rewardTransaction.Outputs()[outputIndex], rewardTransaction.Timestamp())
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
