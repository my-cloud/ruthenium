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
	transactions, err := controller.sender.GetTransactions()
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	io.NewIoWriter(writer, controller.logger).Write(string(transactions[:]))
}
