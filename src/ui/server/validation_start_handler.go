package server

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"net/http"
)

type ValidationStartHandler struct {
	host   neighborhood.Neighbor
	logger *log.Logger
}

func NewValidationStartHandler(host neighborhood.Neighbor, logger *log.Logger) *ValidationStartHandler {
	return &ValidationStartHandler{host, logger}
}

func (handler *ValidationStartHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := handler.host.StartMining()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to start mining: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
