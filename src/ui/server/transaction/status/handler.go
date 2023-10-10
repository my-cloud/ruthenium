package status

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	host     network.Neighbor
	settings server.Settings
	watch    clock.Watch
	logger   log.Logger
}

func NewHandler(host network.Neighbor, settings server.Settings, watch clock.Watch, logger log.Logger) *Handler {
	return &Handler{host, settings, watch, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		jsonWriter := server.NewIoWriter(writer, handler.logger)
		decoder := json.NewDecoder(req.Body)
		var transaction *verification.Transaction
		err := decoder.Decode(&transaction)
		if err != nil || len(transaction.Outputs()) != 2 {
			handler.logger.Error(fmt.Errorf("failed to decode transaction: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid transaction")
			return
		}
		outputIndex := 1
		rest := transaction.Outputs()[outputIndex]
		utxosBytes, err := handler.host.GetUtxos(rest.Address())
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var utxos []*verification.Utxo
		err = json.Unmarshal(utxosBytes, &utxos)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to unmarshal UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		genesisTimestamp, err := handler.host.GetFirstBlockTimestamp()
		now := handler.watch.Now().UnixNano()
		currentBlockHeight := (now - genesisTimestamp) / handler.settings.ValidationTimestamp()
		currentBlockTimestamp := genesisTimestamp + currentBlockHeight*handler.settings.ValidationTimestamp()
		progress := &Progress{
			CurrentBlockTimestamp: currentBlockTimestamp,
			ValidationTimestamp:   handler.settings.ValidationTimestamp(),
		}
		for _, utxo := range utxos {
			if utxo.TransactionId() == transaction.Id() && utxo.OutputIndex() == uint16(outputIndex) {
				progress.TransactionStatus = "confirmed"
				handler.sendResponse(writer, progress)
				return
			}
		}
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get genesis timestamp: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		blocksBytes, err := handler.host.GetBlocks(uint64(currentBlockHeight))
		var blocks []*verification.Block
		err = json.Unmarshal(blocksBytes, &blocks)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to unmarshal blocks: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(blocks) == 0 {
			handler.logger.Error("failed to get last block")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, validatedTransaction := range blocks[0].Transactions() {
			if validatedTransaction.Id() == transaction.Id() {
				progress.TransactionStatus = "validated"
				handler.sendResponse(writer, progress)
				return
			}
		}
		transactionsBytes, err := handler.host.GetTransactions()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var transactions []*verification.Transaction
		err = json.Unmarshal(transactionsBytes, &transactions)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to unmarshal transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, pendingTransaction := range transactions {
			if pendingTransaction.Id() == transaction.Id() {
				progress.TransactionStatus = "sent"
				handler.sendResponse(writer, progress)
				return
			}
		}
		progress.TransactionStatus = "rejected"
		handler.sendResponse(writer, progress)
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (handler *Handler) sendResponse(writer http.ResponseWriter, progress *Progress) {
	marshaledResponse, err := json.Marshal(progress)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	server.NewIoWriter(writer, handler.logger).Write(string(marshaledResponse[:]))
}
