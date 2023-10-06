package verification

import (
	"encoding/json"
	"math"
)

type detailedOutputDto struct {
	Address       string `json:"address"`
	BlockHeight   int    `json:"block_height"`
	HasReward     bool   `json:"has_reward"`
	HasIncome     bool   `json:"has_income"`
	OutputIndex   uint16 `json:"output_index"`
	TransactionId string `json:"transaction_id"`
	Value         uint64 `json:"value"`
}

type DetailedOutput struct {
	*Output
	blockHeight   int
	outputIndex   uint16
	transactionId string
}

func NewDetailedOutput(output *Output, blockHeight int, outputIndex uint16, transactionId string) *DetailedOutput {
	return &DetailedOutput{output, blockHeight, outputIndex, transactionId}
}

func (detailedOutput *DetailedOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(detailedOutputDto{
		Address:       detailedOutput.Address(),
		BlockHeight:   detailedOutput.blockHeight,
		HasIncome:     detailedOutput.HasIncome(),
		HasReward:     detailedOutput.HasReward(),
		OutputIndex:   detailedOutput.outputIndex,
		TransactionId: detailedOutput.transactionId,
		Value:         detailedOutput.InitialValue(),
	})
}

func (detailedOutput *DetailedOutput) UnmarshalJSON(data []byte) error {
	var dto *detailedOutputDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	detailedOutput.Output = NewOutput(dto.Address, dto.HasIncome, dto.HasReward, dto.Value)
	detailedOutput.blockHeight = dto.BlockHeight
	detailedOutput.outputIndex = dto.OutputIndex
	detailedOutput.transactionId = dto.TransactionId
	return nil
}

func (detailedOutput *DetailedOutput) OutputIndex() uint16 {
	return detailedOutput.outputIndex
}

func (detailedOutput *DetailedOutput) TransactionId() string {
	return detailedOutput.transactionId
}

func (detailedOutput *DetailedOutput) Value(currentTimestamp int64, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, validationTimestamp int64) uint64 {
	outputTimestamp := genesisTimestamp + int64(detailedOutput.blockHeight)*validationTimestamp
	if currentTimestamp == outputTimestamp {
		return detailedOutput.InitialValue()
	}
	x := float64(currentTimestamp - outputTimestamp)
	if detailedOutput.HasIncome() {
		return detailedOutput.g(halfLifeInNanoseconds, incomeBase, incomeLimit, x)
	} else {
		return detailedOutput.f(halfLifeInNanoseconds, x)
	}
}

func (detailedOutput *DetailedOutput) f(h float64, x float64) uint64 {
	y := float64(detailedOutput.InitialValue())
	result := y * math.Exp(-x*math.Log(2)/h)
	return uint64(result)
}

func (detailedOutput *DetailedOutput) g(h float64, incomeBase uint64, incomeLimit uint64, x float64) uint64 {
	y := float64(detailedOutput.InitialValue())
	l := float64(incomeLimit)
	if detailedOutput.InitialValue() < incomeLimit {
		k1 := k1(incomeBase, incomeLimit)
		k2 := k2(incomeBase, incomeLimit, k1)
		exp := -math.Pow(x*math.Log(2)/(k2*h)+math.Pow(-math.Log((l-y)/l), 1/k1), k1)
		result := math.Floor(-l*math.Exp(exp)) + l
		return uint64(result)
	} else if incomeLimit < detailedOutput.InitialValue() {
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
