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
	validationTimestamp   int64
	watch                 clock.Watch
	logger                log.Logger
}

func NewHandler(host network.Neighbor, lambda float64, minimalTransactionFee uint64, particlesCount uint64, validationTimestamp int64, watch clock.Watch, logger log.Logger) *Handler {
	return &Handler{host, lambda, minimalTransactionFee, particlesCount, validationTimestamp, watch, logger}
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
		requestIsRegistered := req.URL.Query().Get("registered")
		isRegistered, err := strconv.ParseBool(requestIsRegistered)
		if err != nil {
			errorMessage := "failed to parse registered value"
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
		genesisBlock, err := handler.host.GetBlock(0)
		if err != nil || genesisBlock == nil {
			handler.logger.Error(fmt.Errorf("failed to get genesis block: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var selectedUtxos []*network.UtxoResponse
		var inputsValue uint64
		now := handler.watch.Now().UnixNano()
		nextBlockHeight := (now-genesisBlock.Timestamp)/handler.validationTimestamp + 1
		nextBlockTimestamp := genesisBlock.Timestamp + nextBlockHeight*handler.validationTimestamp
		value := uint64(parsedValue)
		for _, utxo := range utxos {
			output := validation.NewOutputFromUtxoResponse(utxo, handler.lambda, handler.validationTimestamp, genesisBlock.Timestamp)
			outputValue := output.Value(nextBlockTimestamp)
			if isRegistered {
				inputsValue += outputValue
			} else if inputsValue < value {
				inputsValue += outputValue
				selectedUtxos = append(selectedUtxos, utxo)
			}
		}
		if inputsValue < value {
			handler.logger.Error(errors.New("insufficient wallet balance to send transaction").Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("insufficient wallet balance to send transaction")
			return
		}
		if isRegistered {
			selectedUtxos = utxos
		}
		rest := inputsValue - value - handler.minimalTransactionFee
		response := &Response{
			BlockHeight: int(nextBlockHeight),
			Rest:        rest,
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
