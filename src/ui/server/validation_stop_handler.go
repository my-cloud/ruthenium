package server

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"net/http"
)

type ValidationStopHandler struct {
	host   neighborhood.Neighbor
	logger *log.Logger
}

func NewValidationStopHandler(host neighborhood.Neighbor, logger *log.Logger) *ValidationStopHandler {
	return &ValidationStopHandler{host, logger}
}

func (handler *ValidationStopHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := handler.host.StopMining()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to stop mining: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
