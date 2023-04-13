package utxos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
	"strconv"
)

type Handler struct {
	host                  network.Neighbor
	lambda                float64
	minimalTransactionFee uint64
	particlesCount        uint64
	genesisTimestamp      int64
	validationTimestamp   int64
	watch                 clock.Watch
	logger                log.Logger
}

func NewHandler(host network.Neighbor, lambda float64, minimalTransactionFee uint64, particlesCount uint64, genesisTimestamp int64, validationTimestamp int64, watch clock.Watch, logger log.Logger) *Handler {
	return &Handler{host, lambda, minimalTransactionFee, particlesCount, genesisTimestamp, validationTimestamp, watch, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		jsonWriter := server.NewIoWriter(writer, handler.logger)
		address := req.URL.Query().Get("address")
		requestValue := req.URL.Query().Get("value")
		parsedValue, err := strconv.Atoi(requestValue)
		if err != nil {
			errorMessage := "failed to parse transaction value"
			handler.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write(errorMessage)
			return
		}
		utxos, err := handler.host.GetUtxos(address)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get UTXOs: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var selectedUtxos []*network.WalletOutputResponse
		var inputsValue uint64
		var inputsValueForIncome uint64
		now := handler.watch.Now().UnixNano()
		nextBlockHeight := (now-handler.genesisTimestamp)/handler.validationTimestamp + 1
		nextBlockTimestamp := handler.genesisTimestamp + nextBlockHeight*handler.validationTimestamp
		value := uint64(parsedValue)
		var hasIncome bool
		for _, utxo := range utxos {
			output := validation.NewOutputFromWalletResponse(utxo, handler.lambda, handler.validationTimestamp, handler.genesisTimestamp)
			outputValue := output.Value(nextBlockTimestamp)
			inputsValueForIncome += outputValue
			if inputsValue < value {
				inputsValue += outputValue
				selectedUtxos = append(selectedUtxos, utxo)
			}
			if utxo.HasIncome {
				hasIncome = true
			}
		}
		if inputsValue < value {
			handler.logger.Error(errors.New("insufficient wallet balance to send transaction").Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("insufficient wallet balance to send transaction")
			return
		}
		if hasIncome {
			selectedUtxos = utxos
		}
		response := &Response{
			BlockHeight: int(nextBlockHeight),
			HasIncome:   hasIncome,
			Rest:        inputsValue - value - handler.minimalTransactionFee,
			Utxos:       selectedUtxos,
		}
		marshaledResponse, err := json.Marshal(response)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledResponse[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
