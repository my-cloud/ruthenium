package io

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type Response struct {
	writer http.ResponseWriter
	logger log.Logger
}

func NewResponse(writer http.ResponseWriter, logger log.Logger) *Response {
	return &Response{writer, logger}
}

func (response *Response) Write(statusCode int, message string) {
	response.writer.WriteHeader(statusCode)
	i, err := io.WriteString(response.writer, message)
	if err != nil || i == 0 {
		response.logger.Error(fmt.Sprintf("failed to write message: %s", message))
	}
}

func (response *Response) WriteJson(statusCode int, object interface{}) {
	marshaledObject, err := json.Marshal(object)
	if err != nil {
		errorMessage := "failed to marshal response message"
		response.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	response.writer.Header().Add("Content-Type", "application/json")
	message := string(marshaledObject[:])
	response.Write(statusCode, message)
}
