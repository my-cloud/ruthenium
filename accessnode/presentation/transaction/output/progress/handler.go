package progress

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction/output"
	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type Handler struct {
	neighbor network.NeighborCaller
	settings SettingsProvider
	watch    ledger.TimeProvider
	logger   log.Logger
}

func NewHandler(neighbor network.NeighborCaller, settings SettingsProvider, watch ledger.TimeProvider, logger log.Logger) *Handler {
	return &Handler{neighbor, settings, watch, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		jsonWriter := io.NewIoWriter(writer, handler.logger)
		decoder := json.NewDecoder(req.Body)
		var searchedUtxo *protocol.Utxo
		err := decoder.Decode(&searchedUtxo)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to decode utxo: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid utxo")
			return
		}
		utxosBytes, err := handler.neighbor.GetUtxos(searchedUtxo.Address())
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var utxos []*protocol.Utxo
		err = json.Unmarshal(utxosBytes, &utxos)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to unmarshal UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		genesisTimestamp, err := handler.neighbor.GetFirstBlockTimestamp()
		now := handler.watch.Now().UnixNano()
		currentBlockHeight := (now - genesisTimestamp) / handler.settings.ValidationTimestamp()
		currentBlockTimestamp := genesisTimestamp + currentBlockHeight*handler.settings.ValidationTimestamp()
		progressInfo := &output.ProgressInfo{
			CurrentBlockTimestamp: currentBlockTimestamp,
			ValidationTimestamp:   handler.settings.ValidationTimestamp(),
		}
		for _, utxo := range utxos {
			if utxo.TransactionId() == searchedUtxo.TransactionId() && utxo.OutputIndex() == searchedUtxo.OutputIndex() {
				progressInfo.TransactionStatus = "confirmed"
				handler.sendResponse(writer, progressInfo)
				return
			}
		}
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get genesis timestamp: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		blocksBytes, err := handler.neighbor.GetBlocks(uint64(currentBlockHeight))
		if err != nil {
			handler.logger.Error("failed to get blocks")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var blocks []*protocol.Block
		err = json.Unmarshal(blocksBytes, &blocks)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to unmarshal blocks: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(blocks) == 0 {
			handler.logger.Error("failed to get last block, get blocks returned an empty list")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, validatedTransaction := range blocks[0].Transactions() {
			if validatedTransaction.Id() == searchedUtxo.TransactionId() {
				progressInfo.TransactionStatus = "validated"
				handler.sendResponse(writer, progressInfo)
				return
			}
		}
		transactionsBytes, err := handler.neighbor.GetTransactions()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var transactions []*protocol.Transaction
		err = json.Unmarshal(transactionsBytes, &transactions)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to unmarshal transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, pendingTransaction := range transactions {
			if pendingTransaction.Id() == searchedUtxo.TransactionId() {
				progressInfo.TransactionStatus = "sent"
				handler.sendResponse(writer, progressInfo)
				return
			}
		}
		progressInfo.TransactionStatus = "rejected"
		handler.sendResponse(writer, progressInfo)
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (handler *Handler) sendResponse(writer http.ResponseWriter, progress *output.ProgressInfo) {
	marshaledResponse, err := json.Marshal(progress)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to marshal progress: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	io.NewIoWriter(writer, handler.logger).Write(string(marshaledResponse[:]))
}
