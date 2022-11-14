package server

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"net/http"
)

type TransactionsHandler struct {
	host   network.Neighbor
	logger *log.Logger
}

func NewTransactionsHandler(host network.Neighbor, logger *log.Logger) *TransactionsHandler {
	return &TransactionsHandler{host, logger}
}

func (handler *TransactionsHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := handler.host.GetTransactions()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		marshaledTransactions, err := json.Marshal(transactions)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		NewIoWriter(writer, handler.logger).Write(string(marshaledTransactions[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
