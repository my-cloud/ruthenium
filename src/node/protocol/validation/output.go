package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
)

type Output struct {
	address     string
	blockHeight int
	hasReward   bool
	hasIncome   bool
	value       uint64
}

func NewOutputFromUtxoResponse(response *network.UtxoResponse) *Output {
	return &Output{response.Address, response.BlockHeight, response.HasReward, response.HasIncome, response.Value}
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

func (output *Output) Value(currentTimestamp int64, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, validationTimestamp int64) uint64 {
	outputTimestamp := genesisTimestamp + int64(output.blockHeight)*validationTimestamp
	x := float64(currentTimestamp - outputTimestamp)
	if output.hasIncome {
		return output.g(halfLifeInNanoseconds, incomeBase, incomeLimit, x)
	} else {
		return output.f(halfLifeInNanoseconds, x)
	}
}

func (output *Output) f(h float64, x float64) uint64 {
	y := float64(output.value)
	result := y * math.Exp(-x*math.Log(2)/h)
	return uint64(result)
}

func (output *Output) g(h float64, incomeBase uint64, incomeLimit uint64, x float64) uint64 {
	y := float64(output.value)
	l := float64(incomeLimit)
	if output.value < incomeLimit {
		k1 := k1(incomeBase, incomeLimit)
		k2 := k2(incomeBase, incomeLimit, k1)
		exp := -math.Pow(x*math.Log(2)/(k2*h)+math.Pow(-math.Log((l-y)/l), 1/k1), k1)
		result := math.Floor(-l*math.Exp(exp)) + l
		return uint64(result)
	} else if incomeLimit < output.value {
		exp := -x * math.Log(2) / h
		result := math.Floor((y-l)*math.Exp(exp)) + l
		return uint64(result)
	} else {
		return incomeLimit
	}
}

func k1(incomeBase uint64, incomeLimit uint64) float64 {
	if incomeLimit > incomeBase {
		return 3 - 2*math.Log(2*float64(incomeBase))/math.Log(float64(incomeLimit))
	} else {
		return 1
	}
}

func k2(incomeBase uint64, incomeLimit uint64, k1 float64) float64 {
	if incomeLimit > incomeBase {
		return math.Log(2) / math.Pow(-math.Log(1-float64(incomeBase)/float64(incomeLimit)), 1/k1)
	} else {
		return 1
	}
}
