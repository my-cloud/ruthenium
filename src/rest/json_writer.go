package rest

import (
	"encoding/json"
	"io"
	"log"
)

type status struct {
	message string
}

func newStatus(message string) *status {
	return &status{message}
}

func (status *status) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: status.message,
	})
}

func (status *status) stringValue() string {
	marshaledStatus, err := status.MarshalJSON()
	if err != nil {
		log.Println("ERROR: Failed to marshal status")
	}
	return string(marshaledStatus)
}

type JsonWriter struct {
	writer io.Writer
}

func NewJsonWriter(writer io.Writer) *JsonWriter {
	return &JsonWriter{writer}
}

func (jsonWriter *JsonWriter) WriteStatus(message string) {
	messageStatus := newStatus(message)
	i, err := io.WriteString(jsonWriter.writer, messageStatus.stringValue())
	if err != nil || i == 0 {
		log.Println("ERROR: Failed to write status")
	}
}
