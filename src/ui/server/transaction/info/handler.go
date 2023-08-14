package info

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
	"sort"
	"strconv"
)

type Handler struct {
	host                  network.Neighbor
	halfLifeInNanoseconds float64
	minimalTransactionFee uint64
	particlesCount        uint64
	validationTimestamp   int64
	watch                 clock.Watch
	logger                log.Logger
}

func NewHandler(host network.Neighbor, halfLifeInNanoseconds float64, minimalTransactionFee uint64, particlesCount uint64, validationTimestamp int64, watch clock.Watch, logger log.Logger) *Handler {
	return &Handler{host, halfLifeInNanoseconds, minimalTransactionFee, particlesCount, validationTimestamp, watch, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		jsonWriter := server.NewIoWriter(writer, handler.logger)
		address := req.URL.Query().Get("address")
		if address == "" {
			errorMessage := "address is missing in amount request"
			handler.logger.Error(errorMessage)
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write(errorMessage)
			return
		}
		requestValue := req.URL.Query().Get("value")
		parsedValue, err := strconv.Atoi(requestValue)
		if err != nil {
			errorMessage := "failed to parse transaction value"
			handler.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write(errorMessage)
			return
		}
		requestConsolidation := req.URL.Query().Get("consolidation")
		isConsolidationRequired, err := strconv.ParseBool(requestConsolidation)
		if err != nil {
			errorMessage := "failed to parse consolidation value"
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
		var selectedUtxos []*UtxoResponse
		var inputsValue uint64
		now := handler.watch.Now().UnixNano()
		nextBlockHeight := (now-genesisBlock.Timestamp)/handler.validationTimestamp + 1
		nextBlockTimestamp := genesisBlock.Timestamp + nextBlockHeight*handler.validationTimestamp
		// TODO avoid creating NewOutputFromUtxoResponse twice
		sort.Slice(utxos, func(i, j int) bool {
			firstOutput := validation.NewOutputFromUtxoResponse(utxos[i], handler.halfLifeInNanoseconds, handler.validationTimestamp, genesisBlock.Timestamp)
			firstOutputValue := firstOutput.Value(int(nextBlockHeight), nextBlockTimestamp)
			secondOutput := validation.NewOutputFromUtxoResponse(utxos[j], handler.halfLifeInNanoseconds, handler.validationTimestamp, genesisBlock.Timestamp)
			secondOutputValue := secondOutput.Value(int(nextBlockHeight), nextBlockTimestamp)
			return firstOutputValue < secondOutputValue
		})
		value := uint64(parsedValue)
		for _, utxo := range utxos {
			output := validation.NewOutputFromUtxoResponse(utxo, handler.halfLifeInNanoseconds, handler.validationTimestamp, genesisBlock.Timestamp)
			outputValue := output.Value(int(nextBlockHeight), nextBlockTimestamp)
			if isConsolidationRequired {
				inputsValue += outputValue
				selectedUtxos = append(selectedUtxos, &UtxoResponse{
					OutputIndex:   utxo.OutputIndex,
					TransactionId: utxo.TransactionId,
				})
			} else if inputsValue < value {
				inputsValue += outputValue
				selectedUtxos = append(selectedUtxos, &UtxoResponse{
					OutputIndex:   utxo.OutputIndex,
					TransactionId: utxo.TransactionId,
				})
			} else {
				break
			}
		}
		if inputsValue < value {
			errorMessage := "insufficient wallet balance"
			handler.logger.Error(errors.New(errorMessage).Error())
			writer.WriteHeader(http.StatusMethodNotAllowed)
			jsonWriter.Write(errorMessage)
			return
		}
		rest := inputsValue - value - handler.minimalTransactionFee
		response := &TransactionInfoResponse{
			Rest:  rest,
			Utxos: selectedUtxos,
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
