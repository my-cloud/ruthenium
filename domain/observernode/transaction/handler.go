package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/domain/ledger"
	"github.com/my-cloud/ruthenium/domain/network"
	"github.com/my-cloud/ruthenium/domain/observernode"
	"github.com/my-cloud/ruthenium/infrastructure/log"
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
	case http.MethodPost:
		jsonWriter := observernode.NewIoWriter(writer, handler.logger)
		decoder := json.NewDecoder(req.Body)
		var transaction *ledger.Transaction
		err := decoder.Decode(&transaction)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to decode transaction: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid transaction")
			return
		}
		transactionRequest := ledger.NewTransactionRequest(transaction, handler.host.Target())
		marshaledTransaction, err := json.Marshal(transactionRequest)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal transaction request: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = handler.host.AddTransaction(marshaledTransaction)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to add transaction: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusCreated)
		jsonWriter.Write("success")
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
