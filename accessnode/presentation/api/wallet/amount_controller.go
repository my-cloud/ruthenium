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
	address := req.URL.Query().Get("address")
	if address == "" {
		controller.logger.Error("address is missing in amount request")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	utxosBytes, err := controller.sender.GetUtxos(address)
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
	var balance uint64
	for _, utxo := range utxos {
		now := controller.watch.Now().UnixNano()
		balance += utxo.Value(now, controller.settings.HalfLifeInNanoseconds(), controller.settings.IncomeBase(), controller.settings.IncomeLimit())
	}
	marshaledAmount, err := json.Marshal(float64(balance) / float64(controller.settings.SmallestUnitsPerCoin()))
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	io.NewIoWriter(writer, controller.logger).Write(string(marshaledAmount[:]))
}
