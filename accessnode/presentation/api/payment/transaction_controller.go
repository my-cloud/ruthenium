package payment

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type TransactionController struct {
	sender application.Sender
	logger log.Logger
}

func NewTransactionController(sender application.Sender, logger log.Logger) *TransactionController {
	return &TransactionController{sender, logger}
}

func (controller *TransactionController) PostTransaction(writer http.ResponseWriter, req *http.Request) {
	jsonWriter := io.NewIoWriter(writer, controller.logger)
	decoder := json.NewDecoder(req.Body)
	var transaction *protocol.Transaction
	err := decoder.Decode(&transaction)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to decode transaction: %w", err).Error())
		writer.WriteHeader(http.StatusBadRequest)
		jsonWriter.Write("invalid transaction")
		return
	}
	transactionRequest := protocol.NewTransactionRequest(transaction, controller.sender.Target())
	marshaledTransaction, err := json.Marshal(transactionRequest)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to marshal transaction request: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = controller.sender.AddTransaction(marshaledTransaction)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to add transaction: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	jsonWriter.Write("success")
}
