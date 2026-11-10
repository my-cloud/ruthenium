package payment

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"
	"strconv"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type ProgressController struct {
	sender   application.Sender
	settings application.ProtocolSettingsProvider
	watch    application.TimeProvider
	logger   log.Logger
}

func NewProgressController(sender application.Sender, settings application.ProtocolSettingsProvider, watch application.TimeProvider, logger log.Logger) *ProgressController {
	return &ProgressController{sender, settings, watch, logger}
}

func (controller *ProgressController) GetTransactionProgress(writer http.ResponseWriter, req *http.Request) {
	response := io.NewResponse(writer, controller.logger)
	address := req.URL.Query().Get("address")
	if address == "" {
		errorMessage := "address is missing in progress request"
		controller.logger.Error(errorMessage)
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	transactionId := req.URL.Query().Get("transactionId")
	if transactionId == "" {
		errorMessage := "transaction ID is missing in progress request"
		controller.logger.Error(errorMessage)
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	requestOutputIndex := req.URL.Query().Get("outputIndex")
	parsedOutputIndex, err := strconv.ParseUint(requestOutputIndex, 10, 16)
	if err != nil {
		errorMessage := "failed to parse output index"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	outputIndex := uint16(parsedOutputIndex)
	utxosBytes, err := controller.sender.GetUtxos(address)
	if err != nil {
		errorMessage := "failed to get UTXOs"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var utxos []*ledger.Utxo
	err = json.Unmarshal(utxosBytes, &utxos)
	if err != nil {
		errorMessage := "failed to unmarshal UTXOs"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
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
		if utxo.TransactionId() == transactionId && utxo.OutputIndex() == outputIndex {
			progressInfo.TransactionStatus = "confirmed"
			response.WriteJson(http.StatusOK, progressInfo)
			return
		}
	}
	if err != nil {
		errorMessage := "failed to get genesis timestamp"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	blocksBytes, err := controller.sender.GetBlocks(uint64(currentBlockHeight))
	if err != nil {
		errorMessage := "failed to get blocks"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var blocks []*ledger.Block
	err = json.Unmarshal(blocksBytes, &blocks)
	if err != nil {
		errorMessage := "failed to unmarshal blocks"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	if len(blocks) == 0 {
		errorMessage := "failed to get last block, get blocks returned an empty list"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	for _, validatedTransaction := range blocks[0].Transactions() {
		if validatedTransaction.Id() == transactionId {
			progressInfo.TransactionStatus = "validated"
			response.WriteJson(http.StatusOK, progressInfo)
			return
		}
	}
	transactionsBytes, err := controller.sender.GetTransactions()
	if err != nil {
		errorMessage := "failed to get transactions"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var transactions []*ledger.Transaction
	err = json.Unmarshal(transactionsBytes, &transactions)
	if err != nil {
		errorMessage := "failed to unmarshal transactions"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	for _, pendingTransaction := range transactions {
		if pendingTransaction.Id() == transactionId {
			progressInfo.TransactionStatus = "sent"
			response.WriteJson(http.StatusOK, progressInfo)
			return
		}
	}
	progressInfo.TransactionStatus = "rejected"
	response.WriteJson(http.StatusOK, progressInfo)
}
