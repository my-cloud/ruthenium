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
	"math"
	"net/http"
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
		now := handler.watch.Now().UnixNano()
		nextBlockHeight := (now-genesisBlock.Timestamp)/handler.validationTimestamp + 1
		nextBlockTimestamp := genesisBlock.Timestamp + nextBlockHeight*handler.validationTimestamp
		utxosByValue := make(map[uint64][]*UtxoResponse)
		var walletBalance uint64
		var values []uint64
		for _, utxo := range utxos {
			output := validation.NewOutputFromUtxoResponse(utxo, handler.halfLifeInNanoseconds, handler.validationTimestamp, genesisBlock.Timestamp)
			outputValue := output.Value(utxo.BlockHeight, nextBlockTimestamp)
			utxoResponse := &UtxoResponse{
				OutputIndex:   utxo.OutputIndex,
				TransactionId: utxo.TransactionId,
			}
			walletBalance += outputValue
			if isConsolidationRequired {
				selectedUtxos = append(selectedUtxos, utxoResponse)
			} else {
				if _, ok := utxosByValue[outputValue]; !ok {
					values = append(values, outputValue)
				}
				utxosByValue[outputValue] = append(utxosByValue[outputValue], utxoResponse)
			}
		}
		value := uint64(parsedValue)
		targetValue := value + handler.minimalTransactionFee
		if walletBalance < targetValue {
			errorMessage := "insufficient wallet balance"
			handler.logger.Error(errors.New(errorMessage).Error())
			writer.WriteHeader(http.StatusMethodNotAllowed)
			jsonWriter.Write(errorMessage)
			return
		}

		var inputsValue uint64
		if isConsolidationRequired {
			inputsValue = walletBalance
		} else if values != nil {
			for inputsValue < targetValue {
				closestValueIndex := findClosestValueIndex(targetValue, values)
				closestValue := values[closestValueIndex]
				values[closestValueIndex] = 0
				closestUtxos := utxosByValue[closestValue]
				for i := 0; i < len(closestUtxos) && inputsValue < targetValue; i++ {
					inputsValue += closestValue
					selectedUtxos = append(selectedUtxos, closestUtxos[i])
				}
			}
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

func findClosestValueIndex(target uint64, values []uint64) int {
	closestValueIndex := 0
	closestDifference := uint64(math.MaxUint64)
	var isTargetSmaller bool
	for i, value := range values {
		if value < target && isTargetSmaller {
			continue
		}
		var difference uint64
		if value < target {
			difference = target - value
		} else {
			isTargetSmaller = true
			difference = value - target
		}
		if difference < closestDifference {
			closestValueIndex = i
			closestDifference = difference
		}
	}
	return closestValueIndex
}
