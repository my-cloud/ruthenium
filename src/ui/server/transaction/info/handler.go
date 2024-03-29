package info

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"math"
	"net/http"
	"strconv"
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
		genesisTimestamp, err := handler.host.GetFirstBlockTimestamp()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get genesis timestamp: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var selectedInputs []*verification.InputInfo
		now := handler.watch.Now().UnixNano()
		nextBlockHeight := (now-genesisTimestamp)/handler.settings.ValidationTimestamp() + 1
		nextBlockTimestamp := genesisTimestamp + nextBlockHeight*handler.settings.ValidationTimestamp()
		utxosByValue := make(map[uint64][]*verification.InputInfo)
		var walletBalance uint64
		var values []uint64
		for _, utxo := range utxos {
			utxoValue := utxo.Value(nextBlockTimestamp, handler.settings.HalfLifeInNanoseconds(), handler.settings.IncomeBase(), handler.settings.IncomeLimit())
			walletBalance += utxoValue
			if isConsolidationRequired {
				selectedInputs = append(selectedInputs, utxo.InputInfo)
			} else {
				if _, ok := utxosByValue[utxoValue]; !ok {
					values = append(values, utxoValue)
				}
				utxosByValue[utxoValue] = append(utxosByValue[utxoValue], utxo.InputInfo)
			}
		}
		value := uint64(parsedValue)
		targetValue := value + handler.settings.MinimalTransactionFee()
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
		} else if len(values) != 0 {
			for inputsValue < targetValue {
				closestValueIndex := findClosestValueIndex(targetValue, values)
				closestValue := values[closestValueIndex]
				values[closestValueIndex] = 0
				closestUtxos := utxosByValue[closestValue]
				for i := 0; i < len(closestUtxos) && inputsValue < targetValue; i++ {
					inputsValue += closestValue
					selectedInputs = append(selectedInputs, closestUtxos[i])
				}
			}
		}
		rest := inputsValue - targetValue
		response := &TransactionInfo{
			Rest:      rest,
			Inputs:    selectedInputs,
			Timestamp: now,
		}
		marshaledResponse, err := json.Marshal(response)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
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
