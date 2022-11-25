package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"net/http"
)

type Handler struct {
	host   network.Neighbor
	logger *log.Logger
}

func NewHandler(host network.Neighbor, logger *log.Logger) *Handler {
	return &Handler{host, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
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
