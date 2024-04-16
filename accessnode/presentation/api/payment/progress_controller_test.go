package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_GetTransactionProgress_UndecipherableUtxo_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	watchMock := new(application.TimeProviderMock)
	logger := log.NewLoggerMock()
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	data := ""
	marshalledData, _ := json.Marshal(data)
	body := bytes.NewReader(marshalledData)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetUtxosCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionProgress_GetUtxosError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(application.TimeProviderMock)
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	utxo := protocol.NewUtxo(&protocol.InputInfo{}, &protocol.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetUtxosCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionProgress_GetFirstBlockTimestampError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyArray := []byte{91, 93}
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	utxo := protocol.NewUtxo(&protocol.InputInfo{}, &protocol.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionProgress_GetBlocksError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyArray := []byte{91, 93}
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	utxo := protocol.NewUtxo(&protocol.InputInfo{}, &protocol.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1 && len(senderMock.GetBlocksCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionProgress_GetTransactionsError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyArray := []byte{91, 93}
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*protocol.Block{protocol.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return nil, errors.New("") }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	utxo := protocol.NewUtxo(&protocol.InputInfo{}, &protocol.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1 && len(senderMock.GetBlocksCalls()) == 1 && len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionProgress_TransactionNotFound_ReturnsRejected(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyArray := []byte{91, 93}
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*protocol.Block{protocol.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledEmptyArray, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	utxo := protocol.NewUtxo(&protocol.InputInfo{}, &protocol.Output{}, 0)
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1 && len(senderMock.GetBlocksCalls()) == 1 && len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "rejected"
	response := recorder.Body.Bytes()
	var progressInfo *ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}

func Test_GetTransactionProgress_UtxoFound_ReturnsConfirmed(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	transaction, _ := protocol.NewRewardTransaction("", false, 0, 0)
	transactionId := transaction.Id()
	inputInfo := protocol.NewInputInfo(0, transactionId)
	utxo := protocol.NewUtxo(inputInfo, &protocol.Output{}, 0)
	marshalledUtxos, _ := json.Marshal([]*protocol.Utxo{utxo})
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "confirmed"
	response := recorder.Body.Bytes()
	var progressInfo *ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}

func Test_GetTransactionProgress_ValidatedTransactionFound_ReturnsValidated(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyArray := []byte{91, 93}
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	transaction, _ := protocol.NewRewardTransaction("", false, 0, 0)
	blocks := []*protocol.Block{protocol.NewBlock([32]byte{}, nil, nil, 0, []*protocol.Transaction{transaction})}
	marshalledBlocks, _ := json.Marshal(blocks)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	outputIndex := 0
	utxo := protocol.NewUtxo(protocol.NewInputInfo(uint16(outputIndex), transaction.Id()), transaction.Outputs()[outputIndex], transaction.Timestamp())
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1 && len(senderMock.GetBlocksCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "validated"
	response := recorder.Body.Bytes()
	var progressInfo *ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}

func Test_GetTransactionProgress_PendingTransactionFound_ReturnsSent(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyArray := []byte{91, 93}
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyArray, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	blocks := []*protocol.Block{protocol.NewBlock([32]byte{}, nil, nil, 0, nil)}
	marshalledBlocks, _ := json.Marshal(blocks)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return marshalledBlocks, nil }
	transaction, _ := protocol.NewRewardTransaction("", false, 0, 0)
	transactions := []*protocol.Transaction{transaction}
	marshalledTransactions, _ := json.Marshal(transactions)
	senderMock.GetTransactionsFunc = func() ([]byte, error) { return marshalledTransactions, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewProgressController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	outputIndex := 0
	utxo := protocol.NewUtxo(protocol.NewInputInfo(uint16(outputIndex), transaction.Id()), transaction.Outputs()[outputIndex], transaction.Timestamp())
	marshalledUtxo, _ := json.Marshal(utxo)
	body := bytes.NewReader(marshalledUtxo)
	request := httptest.NewRequest(http.MethodPut, "/", body)

	// Act
	controller.GetTransactionProgress(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1 && len(senderMock.GetBlocksCalls()) == 1 && len(senderMock.GetTransactionsCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedStatus := "sent"
	response := recorder.Body.Bytes()
	var progressInfo *ProgressInfo
	err := json.Unmarshal(response, &progressInfo)
	fmt.Println(err)
	actualStatus := progressInfo.TransactionStatus
	test.Assert(t, actualStatus == expectedStatus, fmt.Sprintf("Wrong response. expected: %s actual: %s", expectedStatus, actualStatus))
}
