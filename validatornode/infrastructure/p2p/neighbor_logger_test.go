package p2p

import (
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_Info_MessageLogged_LogsWithTraceLevel(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborLogger := NewNeighborLogger(logger)
	expectedMessage := "info message"

	// Act
	neighborLogger.Info(expectedMessage)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.TraceCalls(), expectedMessage)
	test.Assert(t, len(logger.InfoCalls()) == 0, "message is logged with info level whereas it should not")
}

func Test_Warn_MessageLogged_LogsWithTraceLevel(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborLogger := NewNeighborLogger(logger)
	expectedMessage := "warn message"

	// Act
	neighborLogger.Warn(expectedMessage)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.TraceCalls(), expectedMessage)
	test.Assert(t, len(logger.WarnCalls()) == 0, "message is logged with warn level whereas it should not")
}

func Test_Error_MessageLogged_LogsWithTraceLevel(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	neighborLogger := NewNeighborLogger(logger)
	expectedMessage := "error message"

	// Act
	neighborLogger.Error(expectedMessage)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.TraceCalls(), expectedMessage)
	test.Assert(t, len(logger.ErrorCalls()) == 0, "message is logged with error level whereas it should not")
}
