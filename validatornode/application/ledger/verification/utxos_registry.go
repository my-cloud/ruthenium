package verification

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/array"
)

type UtxosRegistry struct {
	mutex          sync.RWMutex
	utxosByAddress map[string][]*protocol.Utxo
	utxosById      map[string][]*protocol.Utxo
}

func NewUtxosRegistry() *UtxosRegistry {
	utxosPool := &UtxosRegistry{}
	utxosPool.utxosByAddress = make(map[string][]*protocol.Utxo)
	utxosPool.utxosById = make(map[string][]*protocol.Utxo)
	return utxosPool
}

func (registry *UtxosRegistry) Clear() {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.utxosByAddress = make(map[string][]*protocol.Utxo)
	registry.utxosById = make(map[string][]*protocol.Utxo)
}

func (registry *UtxosRegistry) Copy() ledger.UtxosManager {
	poolCopy := &UtxosRegistry{}
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	poolCopy.utxosByAddress = copyUtxosMap(registry.utxosByAddress)
	poolCopy.utxosById = copyUtxosMap(registry.utxosById)
	return poolCopy
}

func (registry *UtxosRegistry) UpdateUtxos(transactionsBytes []byte, timestamp int64) error {
	if transactionsBytes == nil {
		return nil
	}
	var transactions []*protocol.Transaction
	if err := json.Unmarshal(transactionsBytes, &transactions); err != nil {
		return fmt.Errorf("failed to unmarshal transactions: %w", err)
	}
	utxosByAddress := copyUtxosMap(registry.utxosByAddress)
	utxosById := copyUtxosMap(registry.utxosById)
	for _, transaction := range transactions {
		utxosForTransactionId, ok := utxosById[transaction.Id()]
		if ok {
			return fmt.Errorf("transaction ID already exists: %s", transaction.Id())
		}
		if len(transaction.Outputs()) > 1 || transaction.Outputs()[0].InitialValue() > 0 || transaction.Outputs()[0].IsYielding() {
			for j, output := range transaction.Outputs() {
				inputInfo := protocol.NewInputInfo(uint16(j), transaction.Id())
				utxo := protocol.NewUtxo(inputInfo, output, timestamp)
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
				if output != nil && (output.InitialValue() > 0 || output.IsYielding()) {
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
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	registry.utxosById = utxosById
	registry.utxosByAddress = utxosByAddress
	return nil
}

func (registry *UtxosRegistry) Utxo(input ledger.InputInfoProvider) (ledger.UtxoInfoProvider, error) {
	utxos, ok := registry.utxosById[input.TransactionId()]
	if !ok || int(input.OutputIndex()) > len(utxos)-1 {
		return nil, fmt.Errorf("failed to find UTXO, input: %v", input)
	}
	utxo := utxos[input.OutputIndex()]
	if utxo == nil {
		return nil, fmt.Errorf("failed to find UTXO, input: %v", input)
	}
	return utxo, nil
}

func (registry *UtxosRegistry) Utxos(address string) []byte {
	utxos, ok := registry.utxosByAddress[address]
	if !ok {
		return array.MarshalledEmptyArray
	}
	marshaledUtxos, err := json.Marshal(utxos)
	if err != nil {
		return array.MarshalledEmptyArray
	}
	return marshaledUtxos
}

func copyUtxosMap(utxosMap map[string][]*protocol.Utxo) map[string][]*protocol.Utxo {
	utxosMapCopy := make(map[string][]*protocol.Utxo, len(utxosMap))
	for address, utxos := range utxosMap {
		utxosCopy := make([]*protocol.Utxo, len(utxos))
		copy(utxosCopy, utxos)
		utxosMapCopy[address] = utxosCopy
	}
	return utxosMapCopy
}

func removeUtxo(utxos []*protocol.Utxo, transactionId string, outputIndex uint16) []*protocol.Utxo {
	for i := 0; i < len(utxos); i++ {
		if utxos[i].TransactionId() == transactionId && utxos[i].OutputIndex() == outputIndex {
			utxos = append(utxos[:i], utxos[i+1:]...)
			return utxos
		}
	}
	return utxos
}

func verifyIncomes(utxosByAddress map[string][]*protocol.Utxo) error {
	for address, utxos := range utxosByAddress {
		var isYielding bool
		for _, utxo := range utxos {
			if utxo.IsYielding() {
				if isYielding {
					return fmt.Errorf("income requested for several UTXOs for address: %s", address)
				}
				isYielding = true
			}
		}
	}
	return nil
}
