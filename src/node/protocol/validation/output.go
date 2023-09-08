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
	k                     float64
	timestamp             int64
}

func NewOutputFromUtxoResponse(response *network.UtxoResponse, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, validationTimestamp int64) *Output {
	k := math.Log(2) / math.Sqrt(-math.Log(1-float64(incomeBase)/float64(incomeLimit)))
	timestamp := genesisTimestamp + int64(response.BlockHeight)*validationTimestamp
	return &Output{response.Address, response.HasReward, response.HasIncome, response.Value, halfLifeInNanoseconds, incomeLimit, k, timestamp}
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
	return uint64(math.Floor(float64(output.value) * math.Exp(-float64(elapsedTimestamp)*math.Log(2)/output.halfLifeInNanoseconds)))
}

func (output *Output) calculateValue(currentTimestamp int64) uint64 {
	elapsedTimestamp := currentTimestamp - output.timestamp
	g := output.g(float64(elapsedTimestamp), output.value)
	v := uint64(g)
	return v
}

func (output *Output) g(x float64, y uint64) float64 {
	k := output.k
	h := output.halfLifeInNanoseconds
	incomeLimit := float64(output.incomeLimit)
	if y < output.incomeLimit {
		exp := -math.Pow(x*math.Log(2)/(k*h)+math.Sqrt(-math.Log((incomeLimit-float64(y))/incomeLimit)), 2)
		return -incomeLimit*math.Exp(exp) + incomeLimit
	} else if output.incomeLimit < y {
		exp := -x * math.Log(2) / h
		return (float64(y)-incomeLimit)*math.Exp(exp) + incomeLimit
	} else {
		return incomeLimit
	}
}
