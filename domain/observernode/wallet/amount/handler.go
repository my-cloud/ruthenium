package amount

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/domain"
	"github.com/my-cloud/ruthenium/domain/ledger"
	"github.com/my-cloud/ruthenium/domain/network"
	"github.com/my-cloud/ruthenium/domain/observernode"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"net/http"
)

type Handler struct {
	host     network.Neighbor
	settings observernode.SettingsProvider
	watch    domain.TimeProvider
	logger   log.Logger
}

func NewHandler(host network.Neighbor, settings observernode.SettingsProvider, watch domain.TimeProvider, logger log.Logger) *Handler {
	return &Handler{host, settings, watch, logger}
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
		utxosBytes, err := handler.host.GetUtxos(address)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var utxos []*ledger.Utxo
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
		observernode.NewIoWriter(writer, handler.logger).Write(string(marshaledAmount[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
