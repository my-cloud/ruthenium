package transactions

import (
	"fmt"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type Handler struct {
	neighbor network.NeighborCaller
	logger   log.Logger
}

func NewHandler(neighbor network.NeighborCaller, logger log.Logger) *Handler {
	return &Handler{neighbor, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := handler.neighbor.GetTransactions()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		io.NewIoWriter(writer, handler.logger).Write(string(transactions[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
