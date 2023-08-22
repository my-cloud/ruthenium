package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
)

const incomeExponent = 0.54692829

type Output struct {
	address   string
	hasReward bool
	hasIncome bool
	value     uint64

	genesisTimestamp      int64
	halfLifeInNanoseconds float64
	validationTimestamp   int64
}

func NewOutputFromUtxoResponse(response *network.UtxoResponse, halfLifeInNanoseconds float64, validationTimestamp int64, genesisTimestamp int64) *Output {
	return &Output{response.Address, response.HasReward, response.HasIncome, response.Value, genesisTimestamp, halfLifeInNanoseconds, validationTimestamp}
}

func (output *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address   string `json:"address"`
		HasReward bool   `json:"has_reward"`
		HasIncome bool   `json:"has_income"`
		Value     uint64 `json:"value"`
	}{
		Address:   output.address,
		HasReward: output.hasReward,
		HasIncome: output.hasIncome,
		Value:     output.value,
	})
}

func (output *Output) GetResponse() *network.OutputResponse {
	return &network.OutputResponse{
		Address:   output.address,
		HasReward: output.hasReward,
		HasIncome: output.hasIncome,
		Value:     output.value,
	}
}

func (output *Output) Value(blockHeight int, currentTimestamp int64) uint64 {
	if output.hasIncome {
		return output.calculateValue(blockHeight, currentTimestamp)
	} else {
		timestamp := output.genesisTimestamp + int64(blockHeight)*output.validationTimestamp
		return output.decay(timestamp, currentTimestamp, output.value)
	}
}

func (output *Output) decay(lastTimestamp int64, newTimestamp int64, amount uint64) uint64 {
	elapsedTimestamp := newTimestamp - lastTimestamp
	return uint64(math.Floor(float64(amount) * math.Exp(-float64(elapsedTimestamp)/output.halfLifeInNanoseconds)))
}

func (output *Output) calculateValue(blockHeight int, currentTimestamp int64) uint64 {
	totalValue := output.value
	var timestamp int64
	blockTimestamp := output.genesisTimestamp + int64(blockHeight)*output.validationTimestamp
	for timestamp = blockTimestamp; timestamp < currentTimestamp; timestamp += output.validationTimestamp {
		if totalValue > 0 {
			totalValue = output.decay(timestamp, timestamp+output.validationTimestamp, totalValue)
			income := uint64(math.Round(math.Pow(float64(totalValue), incomeExponent)))
			totalValue += income
		}
	}
	return output.decay(timestamp, currentTimestamp, totalValue)
}
