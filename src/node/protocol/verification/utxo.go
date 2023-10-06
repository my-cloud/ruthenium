package verification

import (
	"encoding/json"
	"math"
)

type detailedUtxo struct {
	Address       string `json:"address"`
	Timestamp     int64  `json:"timestamp"`
	HasIncome     bool   `json:"has_income"`
	OutputIndex   uint16 `json:"output_index"`
	TransactionId string `json:"transaction_id"`
	Value         uint64 `json:"value"`
}

type Utxo struct {
	*InputInfo
	*Output
	timestamp int64
}

func NewUtxo(inputInfo *InputInfo, output *Output, timestamp int64) *Utxo {
	return &Utxo{inputInfo, output, timestamp}
}

func (utxo *Utxo) MarshalJSON() ([]byte, error) {
	return json.Marshal(detailedUtxo{
		Address:       utxo.Address(),
		Timestamp:     utxo.timestamp,
		HasIncome:     utxo.IsRegistered(),
		OutputIndex:   utxo.outputIndex,
		TransactionId: utxo.transactionId,
		Value:         utxo.InitialValue(),
	})
}

func (utxo *Utxo) UnmarshalJSON(data []byte) error {
	var dto *detailedUtxo
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	utxo.InputInfo = NewInputInfo(dto.OutputIndex, dto.TransactionId)
	utxo.Output = NewOutput(dto.Address, dto.HasIncome, dto.Value)
	utxo.timestamp = dto.Timestamp
	return nil
}

func (utxo *Utxo) Value(currentTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64) uint64 {
	if currentTimestamp == utxo.timestamp {
		return utxo.InitialValue()
	}
	x := float64(currentTimestamp - utxo.timestamp)
	if utxo.IsRegistered() {
		return utxo.g(halfLifeInNanoseconds, incomeBase, incomeLimit, x)
	} else {
		return utxo.f(halfLifeInNanoseconds, x)
	}
}

func (utxo *Utxo) f(h float64, x float64) uint64 {
	y := float64(utxo.InitialValue())
	result := y * math.Exp(-x*math.Log(2)/h)
	return uint64(result)
}

func (utxo *Utxo) g(h float64, incomeBase uint64, incomeLimit uint64, x float64) uint64 {
	y := float64(utxo.InitialValue())
	l := float64(incomeLimit)
	if utxo.InitialValue() < incomeLimit {
		k1 := k1(incomeBase, incomeLimit)
		k2 := k2(incomeBase, incomeLimit, k1)
		exp := -math.Pow(x*math.Log(2)/(k2*h)+math.Pow(-math.Log((l-y)/l), 1/k1), k1)
		result := math.Floor(-l*math.Exp(exp)) + l
		return uint64(result)
	} else if incomeLimit < utxo.InitialValue() {
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
