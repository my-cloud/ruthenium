package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
	"time"
)

type Handler struct {
	host                network.Neighbor
	validationTimestamp int64
	watch               clock.Watch
	logger              log.Logger
}

func NewHandler(host network.Neighbor, validationTimestamp int64, watch clock.Watch, logger log.Logger) *Handler {
	return &Handler{host, validationTimestamp, watch, logger}
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
		genesisBlock, err := handler.host.GetBlock(0)
		if err != nil || genesisBlock == nil {
			handler.logger.Error(fmt.Errorf("failed to get genesis block: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactionTimestamp := *transactionRequest.Timestamp
		expectedNextBlockHeight := (transactionTimestamp-genesisBlock.Timestamp)/handler.validationTimestamp + 1
		expectedNextBlockTimestamp := genesisBlock.Timestamp + expectedNextBlockHeight*handler.validationTimestamp
		transactionReceptionTimestamp := handler.watch.Now().UnixNano() + time.Second.Nanoseconds()
		nextBlockHeight := (transactionReceptionTimestamp-genesisBlock.Timestamp)/handler.validationTimestamp + 1
		nextBlockTimestamp := genesisBlock.Timestamp + nextBlockHeight*handler.validationTimestamp
		if nextBlockTimestamp > expectedNextBlockTimestamp {
			handler.logger.Error("a new block was created during the transaction")
			writer.WriteHeader(http.StatusConflict)
			return
		}
		err = handler.host.AddTransaction(transactionRequest)
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
