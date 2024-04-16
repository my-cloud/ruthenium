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
	response := io.NewResponse(writer, controller.logger)
	decoder := json.NewDecoder(req.Body)
	var transaction *protocol.Transaction
	err := decoder.Decode(&transaction)
	if err != nil {
		errorMessage := "failed to decode transaction"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	transactionRequest := protocol.NewTransactionRequest(transaction, controller.sender.Target())
	marshaledTransaction, err := json.Marshal(transactionRequest)
	if err != nil {
		errorMessage := "failed to marshal transaction request"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	err = controller.sender.AddTransaction(marshaledTransaction)
	if err != nil {
		errorMessage := "failed to add transaction"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	response.Write(http.StatusCreated, "success")
}
