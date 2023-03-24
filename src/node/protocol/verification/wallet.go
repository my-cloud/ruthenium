package verification

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"math"
	"sort"
)

const incomeExponent = 0.54692829

type Wallet struct {
	address             string
	initialTimestamp    int64
	lambda              float64
	utxosByBlockHeight  map[int][]network.OutputResponse
	validationTimestamp int64
}

func NewWallet(address string, lambda float64, validationTimestamp int64, initialTimestamp int64) *Wallet {
	return &Wallet{address, initialTimestamp, lambda, make(map[int][]network.OutputResponse), validationTimestamp}
}

func (wallet *Wallet) Amount(currentTimestamp int64) uint64 {
	blockHeights := make([]int, 0)
	for blockHeight := range wallet.utxosByBlockHeight {
		blockHeights = append(blockHeights, blockHeight)
	}
	sort.Ints(blockHeights)
	var totalAmount uint64
	var isIncomeCalculated bool
	for _, blockHeight := range blockHeights {
		for _, utxo := range wallet.utxosByBlockHeight[blockHeight] {
			if utxo.HasIncome && !isIncomeCalculated {
				wallet.calculateIncome(currentTimestamp, utxo)
				isIncomeCalculated = true
			} else if totalAmount > 0 {
				blockTimestamp := wallet.initialTimestamp + int64(blockHeight)*wallet.validationTimestamp
				totalAmount += wallet.decay(blockTimestamp, currentTimestamp, utxo.Value)
			}
		}
	}
	return totalAmount
}

func (wallet *Wallet) AddUtxo(output network.OutputResponse) {
	if _, ok := wallet.utxosByBlockHeight[output.BlockHeight]; !ok {
		wallet.utxosByBlockHeight[output.BlockHeight] = []network.OutputResponse{output}
	} else {
		wallet.utxosByBlockHeight[output.BlockHeight] = append(wallet.utxosByBlockHeight[output.BlockHeight], output)
	}
}

func (wallet *Wallet) RemoveUtxo(output network.OutputResponse) {
	if _, ok := wallet.utxosByBlockHeight[output.BlockHeight]; ok {
		wallet.utxosByBlockHeight[output.BlockHeight] = removeUtxo(wallet.utxosByBlockHeight[output.BlockHeight], output)
	}
}

func (wallet *Wallet) IsEmpty() bool {
	return len(wallet.utxosByBlockHeight) == 0
}

func calculateIncome(amount uint64) uint64 {
	return uint64(math.Round(math.Pow(float64(amount), incomeExponent)))
}

func (wallet *Wallet) decay(lastTimestamp int64, newTimestamp int64, amount uint64) uint64 {
	elapsedTimestamp := newTimestamp - lastTimestamp
	return uint64(math.Floor(float64(amount) * math.Exp(-wallet.lambda*float64(elapsedTimestamp))))
}

func (wallet *Wallet) calculateIncome(currentTimestamp int64, utxo network.OutputResponse) uint64 {
	totalIncome := utxo.Value
	var timestamp int64
	blockTimestamp := wallet.initialTimestamp + int64(utxo.BlockHeight)*wallet.validationTimestamp
	for timestamp = blockTimestamp; timestamp < currentTimestamp; timestamp += wallet.validationTimestamp {
		if totalIncome > 0 {
			totalIncome = wallet.decay(timestamp, timestamp-wallet.validationTimestamp, totalIncome)
			totalIncome += calculateIncome(totalIncome)
		}
	}
	return wallet.decay(timestamp, currentTimestamp, totalIncome)
}

func removeUtxo(utxos []network.OutputResponse, utxo network.OutputResponse) []network.OutputResponse {
	for i := 0; i < len(utxos); i++ {
		if utxos[i] == utxo {
			utxos = append(utxos[:i], utxos[i+1:]...)
			return utxos
		}
	}
	return utxos
}
