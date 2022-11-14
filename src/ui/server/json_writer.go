package server

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"io"
	"net/http"
)

type IoWriter struct {
	writer io.Writer
	logger *log.Logger
}

func NewIoWriter(writer http.ResponseWriter, logger *log.Logger) *IoWriter {
	return &IoWriter{writer, logger}
}

func (writer *IoWriter) Write(message string) {
	i, err := io.WriteString(writer.writer, message)
	if err != nil || i == 0 {
		writer.logger.Error(fmt.Sprintf("failed to write message: %s", message))
	}
}
