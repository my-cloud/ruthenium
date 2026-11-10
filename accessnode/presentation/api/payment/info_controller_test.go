package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_GetTransactionInfo_InvalidAddress_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_InvalidValue_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_IsRegisteredNotProvided_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_GetUtxosError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return nil, errors.New("") }
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetUtxosCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_GetFirstBlockTimestampError_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*ledger.Utxo{})
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_InsufficientWalletBalance_ReturnsMethodNotAllowed(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	marshalledEmptyUtxos, _ := json.Marshal([]*ledger.Utxo{})
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledEmptyUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 405
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_ConsolidationNotRequiredAndOneUtxoIsGreater_ReturnsOneUtxo(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)

	inputInfo1 := ledger.NewInputInfo(0, "")
	inputInfo2 := ledger.NewInputInfo(1, "")
	inputInfo3 := ledger.NewInputInfo(2, "")
	output1 := ledger.NewOutput("", false, 1)
	output2 := ledger.NewOutput("", false, 3)
	output3 := ledger.NewOutput("", false, 7)
	utxos := []*ledger.Utxo{
		ledger.NewUtxo(inputInfo1, output1, 1),
		ledger.NewUtxo(inputInfo2, output2, 1),
		ledger.NewUtxo(inputInfo3, output3, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=3&consolidation=false", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo *TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 1
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}

func Test_GetTransactionInfo_ConsolidationNotRequiredAndNoUtxoIsGreater_ReturnsSomeUtxos(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)

	inputInfo1 := ledger.NewInputInfo(0, "")
	inputInfo2 := ledger.NewInputInfo(1, "")
	inputInfo3 := ledger.NewInputInfo(2, "")
	output1 := ledger.NewOutput("", false, 1)
	output2 := ledger.NewOutput("", false, 2)
	output3 := ledger.NewOutput("", false, 2)
	utxos := []*ledger.Utxo{
		ledger.NewUtxo(inputInfo1, output1, 1),
		ledger.NewUtxo(inputInfo2, output2, 1),
		ledger.NewUtxo(inputInfo3, output3, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=3&consolidation=false", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo *TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 2
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}

func Test_GetTransactionInfo_ConsolidationRequired_ReturnsAllUtxos(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	inputInfo1 := ledger.NewInputInfo(0, "")
	inputInfo2 := ledger.NewInputInfo(2, "")
	output1 := ledger.NewOutput("", false, 1)
	output2 := ledger.NewOutput("", false, 2)
	utxos := []*ledger.Utxo{
		ledger.NewUtxo(inputInfo1, output1, 1),
		ledger.NewUtxo(inputInfo2, output2, 1),
	}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=1&consolidation=true", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	areNeighborMethodsCalled := len(senderMock.GetUtxosCalls()) == 1 && len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, areNeighborMethodsCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedInputsCount := 2
	actualInputsCount := len(transactionInfo.Inputs)
	test.Assert(t, actualInputsCount == expectedInputsCount, fmt.Sprintf("Wrong inputs count. expected: %d actual: %d", expectedInputsCount, actualInputsCount))
}

func Test_GetTransactionInfo_InvalidFee_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false&fee=-1", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetUtxosCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_FeeIsLowerThanMinimalFee_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	watchMock := new(application.TimeProviderMock)
	settings := new(application.ProtocolSettingsProviderMock)
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1000 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 100000000 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=0&consolidation=false&fee=999", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetUtxosCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedMessage := "transaction fee is too low, minimal fee: 0.00001 coins"
	actualMessage := recorder.Body.String()
	test.Assert(t, actualMessage == expectedMessage, fmt.Sprintf("Wrong response message. expected: %s actual: %s", expectedMessage, actualMessage))
}

func Test_GetTransactionInfo_FeeIsNotProvided_ReturnsRestWithMinimalFee(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	inputInfo := ledger.NewInputInfo(0, "")
	output := ledger.NewOutput("", false, 10)
	utxos := []*ledger.Utxo{ledger.NewUtxo(inputInfo, output, 1)}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=3&consolidation=true", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo *TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedFee := uint64(1)
	test.Assert(t, transactionInfo.Fee == expectedFee, fmt.Sprintf("Wrong fee. expected: %d actual: %d", expectedFee, transactionInfo.Fee))
	expectedRest := uint64(6)
	test.Assert(t, transactionInfo.Rest == expectedRest, fmt.Sprintf("Wrong rest. expected: %d actual: %d", expectedRest, transactionInfo.Rest))
}

func Test_GetTransactionInfo_FeeIsProvided_ReturnsRestWithProvidedFee(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	inputInfo := ledger.NewInputInfo(0, "")
	output := ledger.NewOutput("", false, 10)
	utxos := []*ledger.Utxo{ledger.NewUtxo(inputInfo, output, 1)}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=3&consolidation=true&fee=4", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	var transactionInfo *TransactionInfo
	_ = json.Unmarshal(recorder.Body.Bytes(), &transactionInfo)
	expectedFee := uint64(4)
	test.Assert(t, transactionInfo.Fee == expectedFee, fmt.Sprintf("Wrong fee. expected: %d actual: %d", expectedFee, transactionInfo.Fee))
	expectedRest := uint64(3)
	test.Assert(t, transactionInfo.Rest == expectedRest, fmt.Sprintf("Wrong rest. expected: %d actual: %d", expectedRest, transactionInfo.Rest))
}

func Test_GetTransactionInfo_FeeExceedsWalletBalance_ReturnsMethodNotAllowed(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	inputInfo := ledger.NewInputInfo(0, "")
	output := ledger.NewOutput("", false, 10)
	utxos := []*ledger.Utxo{ledger.NewUtxo(inputInfo, output, 1)}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=3&consolidation=true&fee=8", "/"), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 405
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetTransactionInfo_ValueAndFeeOverflow_ReturnsBadRequest(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	senderMock := new(application.SenderMock)
	inputInfo := ledger.NewInputInfo(0, "")
	output := ledger.NewOutput("", false, 10)
	utxos := []*ledger.Utxo{ledger.NewUtxo(inputInfo, output, 1)}
	marshalledUtxos, _ := json.Marshal(utxos)
	senderMock.GetUtxosFunc = func(string) ([]byte, error) { return marshalledUtxos, nil }
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	settings := new(application.ProtocolSettingsProviderMock)
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.SmallestUnitsPerCoinFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	controller := NewInfoController(senderMock, settings, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?address=address&value=3&consolidation=true&fee=%d", "/", uint64(math.MaxUint64)), nil)

	// Act
	controller.GetTransactionInfo(recorder, request)

	// Assert
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedMessage := "transaction fee is too high"
	actualMessage := recorder.Body.String()
	test.Assert(t, actualMessage == expectedMessage, fmt.Sprintf("Wrong response message. expected: %s actual: %s", expectedMessage, actualMessage))
}
