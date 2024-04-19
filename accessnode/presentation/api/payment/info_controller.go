package payment

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"math"
	"net/http"
	"strconv"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type InfoController struct {
	sender   application.Sender
	settings SettingsProvider
	watch    application.TimeProvider
	logger   log.Logger
}

func NewInfoController(sender application.Sender, settings SettingsProvider, watch application.TimeProvider, logger log.Logger) *InfoController {
	return &InfoController{sender, settings, watch, logger}
}

func (controller *InfoController) GetTransactionInfo(writer http.ResponseWriter, req *http.Request) {
	response := io.NewResponse(writer, controller.logger)
	address := req.URL.Query().Get("address")
	if address == "" {
		errorMessage := "address is missing in amount request"
		controller.logger.Error(errorMessage)
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	requestValue := req.URL.Query().Get("value")
	parsedValue, err := strconv.Atoi(requestValue)
	if err != nil {
		errorMessage := "failed to parse transaction value"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	requestConsolidation := req.URL.Query().Get("consolidation")
	isConsolidationRequired, err := strconv.ParseBool(requestConsolidation)
	if err != nil {
		errorMessage := "failed to parse consolidation value"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	utxosBytes, err := controller.sender.GetUtxos(address)
	if err != nil {
		errorMessage := "failed to get UTXOs"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var utxos []*ledger.Utxo
	err = json.Unmarshal(utxosBytes, &utxos)
	if err != nil {
		errorMessage := "failed to unmarshal UTXOs"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	genesisTimestamp, err := controller.sender.GetFirstBlockTimestamp()
	if err != nil {
		errorMessage := "failed to get genesis timestamp"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var selectedInputs []*ledger.InputInfo
	now := controller.watch.Now().UnixNano()
	nextBlockHeight := (now-genesisTimestamp)/controller.settings.ValidationTimestamp() + 1
	nextBlockTimestamp := genesisTimestamp + nextBlockHeight*controller.settings.ValidationTimestamp()
	utxosByValue := make(map[uint64][]*ledger.InputInfo)
	var walletBalance uint64
	var values []uint64
	for _, utxo := range utxos {
		utxoValue := utxo.Value(nextBlockTimestamp, controller.settings.HalfLifeInNanoseconds(), controller.settings.IncomeBase(), controller.settings.IncomeLimit())
		if utxoValue == 0 {
			continue
		}
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
	targetValue := value + controller.settings.MinimalTransactionFee()
	if walletBalance < targetValue {
		errorMessage := "insufficient wallet balance"
		controller.logger.Error(errorMessage)
		response.Write(http.StatusMethodNotAllowed, errorMessage)
		return
	}

	var inputsValue uint64
	if isConsolidationRequired {
		inputsValue = walletBalance
	} else if len(values) != 0 {
		for inputsValue < targetValue {
			closestValueIndex := findClosestValueIndex(targetValue, values)
			closestValue := values[closestValueIndex]
			if closestValue > targetValue {
				inputsValue = closestValue
				selectedInputs = []*ledger.InputInfo{utxosByValue[closestValue][0]}
				break
			}
			values = append(values[:closestValueIndex], values[closestValueIndex+1:]...)
			closestUtxos := utxosByValue[closestValue]
			for i := 0; i < len(closestUtxos) && inputsValue < targetValue; i++ {
				inputsValue += closestValue
				selectedInputs = append(selectedInputs, closestUtxos[i])
			}
		}
	}
	rest := inputsValue - targetValue
	transactionInfo := &TransactionInfo{
		Rest:      rest,
		Inputs:    selectedInputs,
		Timestamp: now,
	}
	response.WriteJson(http.StatusOK, transactionInfo)
}

func findClosestValueIndex(target uint64, values []uint64) int {
	closestValueIndex := 0
	closestDifference := uint64(math.MaxUint64)
	var isAValueGreaterThanTarget bool
	for i, value := range values {
		if isAValueGreaterThanTarget && value < target {
			continue
		}
		var difference uint64
		if value < target {
			difference = target - value
		} else {
			if !isAValueGreaterThanTarget {
				closestDifference = uint64(math.MaxUint64)
			}
			isAValueGreaterThanTarget = true
			difference = value - target
		}
		if difference < closestDifference {
			closestValueIndex = i
			closestDifference = difference
		}
	}
	return closestValueIndex
}
