package console

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_NewLogger_TraceLevel_LevelIsTrace(t *testing.T) {
	// Arrange
	// Act
	logger := NewLogger("TRACE")

	// Assert
	test.Assert(t, logger.level == trace, "level is %d whereas it should be %d", logger.level, trace)
}

func Test_NewLogger_UnknownLevel_LevelIsInfo(t *testing.T) {
	// Arrange
	// Act
	logger := NewLogger("unknown")

	// Assert
	test.Assert(t, logger.level == info, "level is %d whereas it should be %d", logger.level, info)
}

func Test_Trace_TraceLevel_MessageIsLogged(t *testing.T) {
	// Arrange
	output := captureOutput(t)
	logger := NewLogger("trace")
	expectedMessage := "trace message"

	// Act
	logger.Trace(expectedMessage)

	// Assert
	test.Assert(t, strings.Contains(output.String(), "TRACE: "+expectedMessage), "message is not logged whereas it should be")
}

func Test_Trace_DebugLevel_MessageIsNotLogged(t *testing.T) {
	// Arrange
	output := captureOutput(t)
	logger := NewLogger("debug")

	// Act
	logger.Trace("trace message")

	// Assert
	test.Assert(t, output.Len() == 0, "message is logged whereas it should not be")
}

func Test_Debug_TraceLevel_MessageIsLogged(t *testing.T) {
	// Arrange
	output := captureOutput(t)
	logger := NewLogger("trace")
	expectedMessage := "debug message"

	// Act
	logger.Debug(expectedMessage)

	// Assert
	test.Assert(t, strings.Contains(output.String(), "DEBUG: "+expectedMessage), "message is not logged whereas it should be")
}

func captureOutput(t *testing.T) *bytes.Buffer {
	output := new(bytes.Buffer)
	log.SetOutput(output)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return output
}
