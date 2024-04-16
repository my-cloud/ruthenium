package payment

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction/output/progress"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type ProgressController struct {
	sender   application.Sender
	settings progress.SettingsProvider
	watch    application.TimeProvider
	logger   log.Logger
}

func NewProgressController(sender application.Sender, settings progress.SettingsProvider, watch application.TimeProvider, logger log.Logger) *ProgressController {
	return &ProgressController{sender, settings, watch, logger}
}

func (controller *ProgressController) GetTransactionProgress(writer http.ResponseWriter, req *http.Request) {
	jsonWriter := io.NewIoWriter(writer, controller.logger)
	decoder := json.NewDecoder(req.Body)
	var searchedUtxo *protocol.Utxo
	err := decoder.Decode(&searchedUtxo)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to decode utxo: %w", err).Error())
		writer.WriteHeader(http.StatusBadRequest)
		jsonWriter.Write("invalid utxo")
		return
	}
	utxosBytes, err := controller.sender.GetUtxos(searchedUtxo.Address())
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to get UTXOs: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	var utxos []*protocol.Utxo
	err = json.Unmarshal(utxosBytes, &utxos)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to unmarshal UTXOs: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	genesisTimestamp, err := controller.sender.GetFirstBlockTimestamp()
	now := controller.watch.Now().UnixNano()
	currentBlockHeight := (now - genesisTimestamp) / controller.settings.ValidationTimestamp()
	currentBlockTimestamp := genesisTimestamp + currentBlockHeight*controller.settings.ValidationTimestamp()
	progressInfo := &ProgressInfo{
		CurrentBlockTimestamp: currentBlockTimestamp,
		ValidationTimestamp:   controller.settings.ValidationTimestamp(),
	}
	for _, utxo := range utxos {
		if utxo.TransactionId() == searchedUtxo.TransactionId() && utxo.OutputIndex() == searchedUtxo.OutputIndex() {
			progressInfo.TransactionStatus = "confirmed"
			controller.sendResponse(writer, progressInfo)
			return
		}
	}
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to get genesis timestamp: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	blocksBytes, err := controller.sender.GetBlocks(uint64(currentBlockHeight))
	if err != nil {
		controller.logger.Error("failed to get blocks")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	var blocks []*protocol.Block
	err = json.Unmarshal(blocksBytes, &blocks)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to unmarshal blocks: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(blocks) == 0 {
		controller.logger.Error("failed to get last block, get blocks returned an empty list")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, validatedTransaction := range blocks[0].Transactions() {
		if validatedTransaction.Id() == searchedUtxo.TransactionId() {
			progressInfo.TransactionStatus = "validated"
			controller.sendResponse(writer, progressInfo)
			return
		}
	}
	transactionsBytes, err := controller.sender.GetTransactions()
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	var transactions []*protocol.Transaction
	err = json.Unmarshal(transactionsBytes, &transactions)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to unmarshal transactions: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, pendingTransaction := range transactions {
		if pendingTransaction.Id() == searchedUtxo.TransactionId() {
			progressInfo.TransactionStatus = "sent"
			controller.sendResponse(writer, progressInfo)
			return
		}
	}
	progressInfo.TransactionStatus = "rejected"
	controller.sendResponse(writer, progressInfo)
}

func (controller *ProgressController) sendResponse(writer http.ResponseWriter, progress *ProgressInfo) {
	marshaledResponse, err := json.Marshal(progress)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to marshal progress: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	io.NewIoWriter(writer, controller.logger).Write(string(marshaledResponse[:]))
}
