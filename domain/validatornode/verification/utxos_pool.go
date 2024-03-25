package verification

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/domain"
	"github.com/my-cloud/ruthenium/domain/ledger"
	"github.com/my-cloud/ruthenium/infrastructure/array"
)

type UtxosPool struct {
	utxosByAddress map[string][]*ledger.Utxo
	utxosById      map[string][]*ledger.Utxo
}

func NewUtxosPool() *UtxosPool {
	return &UtxosPool{make(map[string][]*ledger.Utxo), make(map[string][]*ledger.Utxo)}
}

func (pool *UtxosPool) Clear() {
	pool.utxosByAddress = make(map[string][]*ledger.Utxo)
	pool.utxosById = make(map[string][]*ledger.Utxo)
}

func (pool *UtxosPool) Copy() domain.UtxosManager {
	return &UtxosPool{copyUtxosMap(pool.utxosByAddress), copyUtxosMap(pool.utxosById)}
}

func (pool *UtxosPool) UpdateUtxos(transactionsBytes []byte, timestamp int64) error {
	if transactionsBytes == nil {
		return nil
	}
	var transactions []*ledger.Transaction
	if err := json.Unmarshal(transactionsBytes, &transactions); err != nil {
		return fmt.Errorf("failed to unmarshal transactions: %w", err)
	}
	utxosByAddress := copyUtxosMap(pool.utxosByAddress)
	utxosById := copyUtxosMap(pool.utxosById)
	for _, transaction := range transactions {
		utxosForTransactionId, ok := utxosById[transaction.Id()]
		if ok {
			return fmt.Errorf("transaction ID already exists: %s", transaction.Id())
		}
		for j, output := range transaction.Outputs() {
			if output.InitialValue() > 0 {
				inputInfo := ledger.NewInputInfo(uint16(j), transaction.Id())
				utxo := ledger.NewUtxo(inputInfo, output, timestamp)
				utxosForTransactionId = append(utxosForTransactionId, utxo)
				utxosById[transaction.Id()] = utxosForTransactionId
				utxosByAddress[output.Address()] = append(utxosByAddress[output.Address()], utxo)
			}
		}
		for _, input := range transaction.Inputs() {
			utxosForInputTransactionId := utxosById[input.TransactionId()]
			if int(input.OutputIndex()) > len(utxosForInputTransactionId)-1 {
				return fmt.Errorf("failed to find UTXO, input: %v", input)
			}
			utxo := utxosForInputTransactionId[input.OutputIndex()]
			if utxo == nil {
				return fmt.Errorf("failed to find output index, input: %v", input)
			}
			utxosForUtxoAddress := utxosByAddress[utxo.Address()]
			utxosForUtxoAddress = removeUtxo(utxosForUtxoAddress, input.TransactionId(), input.OutputIndex())
			utxosByAddress[utxo.Address()] = utxosForUtxoAddress
			utxosById[input.TransactionId()][input.OutputIndex()] = nil
			isEmpty := true
			for _, output := range utxosForInputTransactionId {
				if output != nil {
					isEmpty = false
					break
				}
			}
			if isEmpty {
				delete(utxosById, input.TransactionId())
			}
			if len(utxosForUtxoAddress) == 0 {
				delete(utxosByAddress, utxo.Address())
			}
		}
	}
	if err := verifyIncomes(utxosByAddress); err != nil {
		return err
	}
	pool.utxosById = utxosById
	pool.utxosByAddress = utxosByAddress
	return nil
}

func (pool *UtxosPool) Utxo(input domain.InputInfoProvider) (domain.UtxoInfoProvider, error) {
	utxos, ok := pool.utxosById[input.TransactionId()]
	if !ok || int(input.OutputIndex()) > len(utxos)-1 {
		return nil, fmt.Errorf("failed to find UTXO, input: %v", input)
	}
	utxo := utxos[input.OutputIndex()]
	if utxo == nil {
		return nil, fmt.Errorf("failed to find UTXO, input: %v", input)
	}
	return utxo, nil
}

func (pool *UtxosPool) Utxos(address string) []byte {
	utxos, ok := pool.utxosByAddress[address]
	if !ok {
		return array.MarshalledEmptyArray
	}
	marshaledUtxos, err := json.Marshal(utxos)
	if err != nil {
		return array.MarshalledEmptyArray
	}
	return marshaledUtxos
}

func copyUtxosMap(utxosMap map[string][]*ledger.Utxo) map[string][]*ledger.Utxo {
	utxosMapCopy := make(map[string][]*ledger.Utxo, len(utxosMap))
	for address, utxos := range utxosMap {
		utxosCopy := make([]*ledger.Utxo, len(utxos))
		copy(utxosCopy, utxos)
		utxosMapCopy[address] = utxosCopy
	}
	return utxosMapCopy
}
