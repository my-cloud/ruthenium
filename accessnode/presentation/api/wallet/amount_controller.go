package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/accessnode/presentation/wallet"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type AmountController struct {
	sender   application.Sender
	settings wallet.SettingsProvider
	watch    application.TimeProvider
	logger   log.Logger
}

func NewAmountController(sender application.Sender, settings wallet.SettingsProvider, watch application.TimeProvider, logger log.Logger) *AmountController {
	return &AmountController{sender, settings, watch, logger}
}

func (controller *AmountController) GetWalletAmount(writer http.ResponseWriter, req *http.Request) {
	response := io.NewResponse(writer, controller.logger)
	address := req.URL.Query().Get("address")
	if address == "" {
		errorMessage := "address is missing in amount request"
		controller.logger.Error(errorMessage)
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	utxosBytes, err := controller.sender.GetUtxos(address)
	if err != nil {
		errorMessage := "failed to get UTXOs"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var utxos []*protocol.Utxo
	err = json.Unmarshal(utxosBytes, &utxos)
	if err != nil {
		errorMessage := "failed to unmarshal UTXOs"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var balance uint64
	for _, utxo := range utxos {
		now := controller.watch.Now().UnixNano()
		balance += utxo.Value(now, controller.settings.HalfLifeInNanoseconds(), controller.settings.IncomeBase(), controller.settings.IncomeLimit())
	}
	amount := float64(balance) / float64(controller.settings.SmallestUnitsPerCoin())
	response.WriteJson(http.StatusOK, amount)
}
