package payment

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type TransactionsController struct {
	sender application.Sender
	logger log.Logger
}

func NewTransactionsController(sender application.Sender, logger log.Logger) *TransactionsController {
	return &TransactionsController{sender, logger}
}

func (controller *TransactionsController) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	response := io.NewResponse(writer, controller.logger)
	transactions, err := controller.sender.GetTransactions()
	if err != nil {
		errorMessage := "failed to get transactions"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	response.Write(http.StatusOK, string(transactions[:]))
}
