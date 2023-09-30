package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
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
		var transactionRequest network.TransactionRequest
		jsonWriter := server.NewIoWriter(writer, handler.logger)
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&transactionRequest)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to decode transaction request: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid transaction request")
			return
		}
		hostTarget := handler.host.Target()
		transactionRequest.TransactionBroadcasterTarget = &hostTarget
		if transactionRequest.IsInvalid() {
			errorMessage := "field(s) are missing in transaction request"
			handler.logger.Error(errorMessage)
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write(errorMessage)
			return
		}
		marshaledTransaction, err := json.Marshal(transactionRequest)
		err = handler.host.AddTransaction(marshaledTransaction)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
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
