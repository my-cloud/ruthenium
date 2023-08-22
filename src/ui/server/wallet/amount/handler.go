package amount

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	host                  network.Neighbor
	halfLifeInNanoseconds float64
	particlesCount        uint64
	validationTimestamp   int64
	watch                 clock.Watch
	logger                log.Logger
}

func NewHandler(host network.Neighbor, halfLifeInNanoseconds float64, particlesCount uint64, validationTimestamp int64, watch clock.Watch, logger log.Logger) *Handler {
	return &Handler{host, halfLifeInNanoseconds, particlesCount, validationTimestamp, watch, logger}
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
		utxos, err := handler.host.GetUtxos(address)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		genesisBlock, err := handler.host.GetBlock(0)
		if err != nil || genesisBlock == nil {
			handler.logger.Error(fmt.Errorf("failed to get genesis block: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var balance uint64
		for _, utxo := range utxos {
			balance += validation.NewOutputFromUtxoResponse(utxo, handler.halfLifeInNanoseconds, handler.validationTimestamp, genesisBlock.Timestamp).Value(utxo.BlockHeight, handler.watch.Now().UnixNano())
		}
		marshaledAmount, err := json.Marshal(float64(balance) / float64(handler.particlesCount))
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledAmount[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
