package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
)

const incomeExponent = 0.54692829

type Output struct {
	address     string
	blockHeight int
	hasReward   bool
	hasIncome   bool
	value       uint64

	genesisTimestamp    int64
	lambda              float64
	validationTimestamp int64
}

func NewOutputFromResponse(response *network.OutputResponse, lambda float64, validationTimestamp int64, genesisTimestamp int64) *Output {
	return &Output{response.Address, response.BlockHeight, response.HasReward, response.HasIncome, response.Value, genesisTimestamp, lambda, validationTimestamp}
}

func NewOutputFromUtxoResponse(response *network.UtxoResponse, lambda float64, validationTimestamp int64, genesisTimestamp int64) *Output {
	return &Output{response.Address, response.BlockHeight, response.HasReward, response.HasIncome, response.Value, genesisTimestamp, lambda, validationTimestamp}
}

func (output *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address     string `json:"address"`
		BlockHeight bool   `json:"block_height"`
		HasReward   bool   `json:"has_reward"`
		HasIncome   bool   `json:"has_income"`
		Value       uint64 `json:"value"`
	}{
		Address:   output.address,
		HasReward: output.hasReward,
		HasIncome: output.hasIncome,
		Value:     output.value,
	})
}

func (output *Output) GetResponse() *network.OutputResponse {
	return &network.OutputResponse{
		Address:     output.address,
		BlockHeight: output.blockHeight,
		HasReward:   output.hasReward,
		HasIncome:   output.hasIncome,
		Value:       output.value,
	}
}

func (output *Output) Value(currentTimestamp int64) uint64 {
	if output.hasIncome {
		return output.calculateValue(currentTimestamp)
	} else {
		timestamp := output.genesisTimestamp + int64(output.blockHeight)*output.validationTimestamp
		return output.decay(timestamp, currentTimestamp, output.value)
	}
}

func (output *Output) decay(lastTimestamp int64, newTimestamp int64, amount uint64) uint64 {
	elapsedTimestamp := newTimestamp - lastTimestamp
	return uint64(math.Floor(float64(amount) * math.Exp(-output.lambda*float64(elapsedTimestamp))))
}

func (output *Output) calculateValue(currentTimestamp int64) uint64 {
	totalValue := output.value
	var timestamp int64
	blockTimestamp := output.genesisTimestamp + int64(output.blockHeight)*output.validationTimestamp
	for timestamp = blockTimestamp; timestamp < currentTimestamp; timestamp += output.validationTimestamp {
		if totalValue > 0 {
			totalValue = output.decay(timestamp, timestamp+output.validationTimestamp, totalValue)
			income := uint64(math.Round(math.Pow(float64(totalValue), incomeExponent)))
			totalValue += income
		}
	}
	return output.decay(timestamp, currentTimestamp, totalValue)
}
