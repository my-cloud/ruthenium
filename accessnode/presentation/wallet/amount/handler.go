package amount

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/accessnode/presentation/wallet"
	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/presentation"
)

type Handler struct {
	neighbor presentation.NeighborCaller
	settings wallet.SettingsProvider
	watch    ledger.TimeProvider
	logger   log.Logger
}

func NewHandler(neighbor presentation.NeighborCaller, settings wallet.SettingsProvider, watch ledger.TimeProvider, logger log.Logger) *Handler {
	return &Handler{neighbor, settings, watch, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		if address == "" {
			handler.logger.Error("address is missing in amount request")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		utxosBytes, err := handler.neighbor.GetUtxos(address)
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
		var balance uint64
		for _, utxo := range utxos {
			now := handler.watch.Now().UnixNano()
			balance += utxo.Value(now, handler.settings.HalfLifeInNanoseconds(), handler.settings.IncomeBase(), handler.settings.IncomeLimit())
		}
		marshaledAmount, err := json.Marshal(float64(balance) / float64(handler.settings.SmallestUnitsPerCoin()))
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		io.NewIoWriter(writer, handler.logger).Write(string(marshaledAmount[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
