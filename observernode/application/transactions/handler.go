package transactions

import (
	"fmt"
	"net/http"

	"github.com/my-cloud/ruthenium/common/infrastructure/log"
	"github.com/my-cloud/ruthenium/observernode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/presentation/network"
)

type Handler struct {
	host   network.NeighborController
	logger log.Logger
}

func NewHandler(host network.NeighborController, logger log.Logger) *Handler {
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
		io.NewIoWriter(writer, handler.logger).Write(string(transactions[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
