package verification

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
)

type Utxo struct {
	*Output
	blockHeight   int
	outputIndex   uint16
	transactionId string
}

func NewUtxo(output *Output, blockHeight int, outputIndex uint16, transactionId string) *Utxo {
	return &Utxo{output, blockHeight, outputIndex, transactionId}
}

func (utxo *Utxo) MarshalJSON() ([]byte, error) {
	return json.Marshal(network.UtxoResponse{
		Address:       utxo.Address(),
		BlockHeight:   utxo.blockHeight,
		HasIncome:     utxo.HasIncome(),
		HasReward:     utxo.HasReward(),
		OutputIndex:   utxo.outputIndex,
		TransactionId: utxo.transactionId,
		Value:         utxo.InitialValue(),
	})
}

func (utxo *Utxo) UnmarshalJSON(data []byte) error {
	var utxoDto network.UtxoResponse
	err := json.Unmarshal(data, &utxoDto)
	if err != nil {
		return err
	}
	utxo.Output = NewOutput(utxoDto.Address, utxoDto.HasIncome, utxoDto.HasReward, utxoDto.Value)
	utxo.blockHeight = utxoDto.BlockHeight
	utxo.outputIndex = utxoDto.OutputIndex
	utxo.transactionId = utxoDto.TransactionId
	return nil
}

func (utxo *Utxo) OutputIndex() uint16 {
	return utxo.outputIndex
}

func (utxo *Utxo) TransactionId() string {
	return utxo.transactionId
}

func (utxo *Utxo) Value(currentTimestamp int64, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, validationTimestamp int64) uint64 {
	outputTimestamp := genesisTimestamp + int64(utxo.blockHeight)*validationTimestamp
	if currentTimestamp == outputTimestamp {
		return utxo.InitialValue()
	}
	x := float64(currentTimestamp - outputTimestamp)
	if utxo.HasIncome() {
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
