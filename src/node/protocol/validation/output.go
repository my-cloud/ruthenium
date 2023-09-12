package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
)

type Output struct {
	address   string
	hasReward bool
	hasIncome bool
	value     uint64

	halfLifeInNanoseconds float64
	incomeLimit           uint64
	k1                    float64
	k2                    float64
	timestamp             int64
}

func NewOutputFromUtxoResponse(response *network.UtxoResponse, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, validationTimestamp int64) *Output {
	var k1 float64
	var k2 float64
	if incomeLimit > incomeBase {
		k1 = 3 - 2*math.Log(2*float64(incomeBase))/math.Log(float64(incomeLimit))
		k2 = math.Log(2) / math.Pow(-math.Log(1-float64(incomeBase)/float64(incomeLimit)), 1/k1)
	} else {
		k1 = 1
		k2 = 1
	}
	timestamp := genesisTimestamp + int64(response.BlockHeight)*validationTimestamp
	return &Output{response.Address, response.HasReward, response.HasIncome, response.Value, halfLifeInNanoseconds, incomeLimit, k1, k2, timestamp}
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

func (output *Output) Value(currentTimestamp int64) uint64 {
	if output.hasIncome {
		return output.calculateValue(currentTimestamp)
	} else {
		return output.decay(currentTimestamp)
	}
}

func (output *Output) decay(newTimestamp int64) uint64 {
	elapsedTimestamp := newTimestamp - output.timestamp
	result := float64(output.value) * math.Exp(-float64(elapsedTimestamp)*math.Log(2)/output.halfLifeInNanoseconds)
	return uint64(result)
}

func (output *Output) calculateValue(currentTimestamp int64) uint64 {
	x := float64(currentTimestamp - output.timestamp)
	k1 := output.k1
	k2 := output.k2
	h := output.halfLifeInNanoseconds
	y := float64(output.value)
	l := float64(output.incomeLimit)
	if output.value < output.incomeLimit {
		exp := -math.Pow(x*math.Log(2)/(k2*h)+math.Pow(-math.Log((l-y)/l), 1/k1), k1)
		result := math.Floor(-l*math.Exp(exp)) + l
		return uint64(result)
	} else if output.incomeLimit < output.value {
		exp := -x * math.Log(2) / h
		result := math.Floor((y-l)*math.Exp(exp)) + l
		return uint64(result)
	} else {
		return output.incomeLimit
	}
}
