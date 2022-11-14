package server

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"net/http"
)

type ValidationHandler struct {
	host   network.Neighbor
	logger *log.Logger
}

func NewValidationHandler(host network.Neighbor, logger *log.Logger) *ValidationHandler {
	return &ValidationHandler{host, logger}
}

func (handler *ValidationHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := handler.host.Mine()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to mine: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
