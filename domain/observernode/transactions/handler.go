package transactions

import (
	"fmt"
	"github.com/my-cloud/ruthenium/domain/observernode"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"github.com/my-cloud/ruthenium/infrastructure/network"
	"net/http"
)

type Handler struct {
	host   network.Neighbor
	logger log.Logger
}

func NewHandler(host network.Neighbor, logger log.Logger) *Handler {
	return &Handler{host, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := handler.host.GetTransactions()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		observernode.NewIoWriter(writer, handler.logger).Write(string(transactions[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
