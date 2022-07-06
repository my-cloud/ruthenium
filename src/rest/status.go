package rest

import (
	"encoding/json"
	"io"
	"log"
)

type Status struct {
	message string
}

func NewStatus(message string) *Status {
	return &Status{message}
}

func (status *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: status.message,
	})
}

func (status *Status) Write(writer io.Writer) {
	i, err := io.WriteString(writer, status.stringValue())
	if err != nil || i == 0 {
		log.Println("ERROR: Failed to write status")
	}
}

func (status *Status) stringValue() string {
	marshaledStatus, err := status.MarshalJSON()
	if err != nil {
		log.Println("ERROR: Failed to marshal status")
	}
	return string(marshaledStatus)
}
