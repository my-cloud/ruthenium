package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
	"time"
)

type Output struct {
	address   string
	hasReward bool
	hasIncome bool
	value     uint64
	// TODO add blockheight

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
	return uint64(math.Floor(float64(amount) * math.Exp(-float64(elapsedTimestamp)*math.Log(2)/output.halfLifeInNanoseconds)))
}

func (output *Output) calculateValue(blockHeight int, currentTimestamp int64) uint64 {
	totalValue := output.value
	blockTimestamp := output.genesisTimestamp + int64(blockHeight)*output.validationTimestamp
	elapsedTimestamp := currentTimestamp - blockTimestamp
	g := output.g(float64(elapsedTimestamp), float64(totalValue))
	v := uint64(g)
	return v
}

func (output *Output) g(x, y float64) float64 {
	//unit := 1.
	unit := 24 * float64(time.Hour.Nanoseconds())
	k := 2.3
	R := output.halfLifeInNanoseconds / unit
	var M float64 = 10000000000000
	h := output.h(y, M, k, R)
	x2 := x / unit
	//a := x2/(R*r(M, y)+r(y, M)) + h
	//b := k*r(M, y) + r(y, M)
	//exp := -math.Pow(a, b)/R
	//g := (r(y, M)*y-M)*math.Exp(exp) + M
	//return g
	if y < M {
		a := x2*math.Log(2)/R + h
		exp := -math.Pow(a, k) / R
		return -M*math.Exp(exp) + M
	} else if M < y {
		exp := -x2 * math.Log(2) / R
		return (y-M)*math.Exp(exp) + M
	} else {
		return M
	}
}

func (output *Output) h(y float64, M float64, k float64, R float64) float64 {
	//l := c(M, y)/M + r(y, M)
	//x := -R * math.Log(l)
	//h := math.Pow(x, 1/k)
	//return h
	if y < M {
		return math.Pow(-R*math.Log((M-y)/M), 1/k)
	} else {
		return 0
	}
}

func r(x, y float64) float64 {
	r := c(x, y) / (x - y)
	return r
}

func c(x float64, y float64) float64 {
	c := (math.Abs(x-y) + x - y) / 2
	return c
}
