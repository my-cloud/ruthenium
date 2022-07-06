package main

import (
	"encoding/json"
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

func (status *Status) StringValue() string {
	marshaledStatus, err := status.MarshalJSON()
	if err != nil {
		log.Printf("ERROR: Failed to marshal status")
	}
	return string(marshaledStatus)
}
