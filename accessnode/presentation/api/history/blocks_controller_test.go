package history

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_GetBlocks_HeightIsMissing_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, nil, nil, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blocks", nil)

	// Act
	controller.GetBlocks(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetBlocksCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetBlocks_HeightIsNegative_BadRequest(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, nil, nil, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blocks?height=-1", nil)

	// Act
	controller.GetBlocks(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetBlocksCalls()) != 0
	test.Assert(t, !isNeighborMethodCalled, "Neighbor method is called whereas it should not.")
	expectedStatusCode := 400
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetBlocks_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return nil, errors.New("") }
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, nil, nil, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blocks?height=0", nil)

	// Act
	controller.GetBlocks(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetBlocksCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetBlocks_ValidRequest_BlocksReturned(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	expectedBlocks := `[{"previous_hash":[],"added_registered_addresses":null,"removed_registered_addresses":null,"timestamp":0,"transactions":[]}]`
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) { return []byte(expectedBlocks), nil }
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, nil, nil, logger)
	recorder := httptest.NewRecorder()
	expectedHeight := uint64(42)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/blocks?height=%d", expectedHeight), nil)

	// Act
	controller.GetBlocks(recorder, request)

	// Assert
	getBlocksCalls := senderMock.GetBlocksCalls()
	isNeighborMethodCalled := len(getBlocksCalls) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	actualHeight := getBlocksCalls[0].StartingBlockHeight
	test.Assert(t, actualHeight == expectedHeight, fmt.Sprintf("Wrong starting block height. expected: %d actual: %d", expectedHeight, actualHeight))
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	actualBlocks := recorder.Body.String()
	test.Assert(t, actualBlocks == expectedBlocks, fmt.Sprintf("Wrong blocks. expected: %s actual: %s", expectedBlocks, actualBlocks))
}

func Test_GetCurrentBlockHeight_NodeError_InternalServerError(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, errors.New("") }
	settingsMock := new(application.ProtocolSettingsProviderMock)
	settingsMock.ValidationTimestampFunc = func() int64 { return 3 }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, settingsMock, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blocks/current-height", nil)

	// Act
	controller.GetCurrentBlockHeight(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 500
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
}

func Test_GetCurrentBlockHeight_EmptyBlockchain_ReturnsZero(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 0, nil }
	settingsMock := new(application.ProtocolSettingsProviderMock)
	settingsMock.ValidationTimestampFunc = func() int64 { return 3 }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 10) }
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, settingsMock, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blocks/current-height", nil)

	// Act
	controller.GetCurrentBlockHeight(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedHeight := "0"
	actualHeight := recorder.Body.String()
	test.Assert(t, actualHeight == expectedHeight, fmt.Sprintf("Wrong current block height. expected: %s actual: %s", expectedHeight, actualHeight))
}

func Test_GetCurrentBlockHeight_ValidRequest_ReturnsEstimatedHeight(t *testing.T) {
	// Arrange
	senderMock := new(application.SenderMock)
	senderMock.GetFirstBlockTimestampFunc = func() (int64, error) { return 2, nil }
	settingsMock := new(application.ProtocolSettingsProviderMock)
	settingsMock.ValidationTimestampFunc = func() int64 { return 3 }
	watchMock := new(application.TimeProviderMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 11) }
	logger := log.NewLoggerMock()
	controller := NewBlocksController(senderMock, settingsMock, watchMock, logger)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/blocks/current-height", nil)

	// Act
	controller.GetCurrentBlockHeight(recorder, request)

	// Assert
	isNeighborMethodCalled := len(senderMock.GetFirstBlockTimestampCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
	expectedStatusCode := 200
	test.Assert(t, recorder.Code == expectedStatusCode, fmt.Sprintf("Wrong response status code. expected: %d actual: %d", expectedStatusCode, recorder.Code))
	expectedHeight := "3"
	actualHeight := recorder.Body.String()
	test.Assert(t, actualHeight == expectedHeight, fmt.Sprintf("Wrong current block height. expected: %s actual: %s", expectedHeight, actualHeight))
}
