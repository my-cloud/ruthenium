package amount

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	host     network.Neighbor
	settings server.Settings
	watch    clock.Watch
	logger   log.Logger
}

func NewHandler(host network.Neighbor, settings server.Settings, watch clock.Watch, logger log.Logger) *Handler {
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
		var utxos []*verification.Utxo
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
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledAmount[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
