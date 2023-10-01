package verification

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/config"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
)

type transactionDto struct {
	Id        string    `json:"id"`
	Inputs    []*Input  `json:"inputs"`
	Outputs   []*Output `json:"outputs"`
	Timestamp int64     `json:"timestamp"`
}

type Transaction struct {
	id                     string
	inputs                 []*Input
	outputs                []*Output
	timestamp              int64
	hasReward              bool
	rewardRecipientAddress string
	rewardValue            uint64
}

func NewRewardTransaction(address string, hasIncome bool, timestamp int64, value uint64) (*Transaction, error) {
	outputs := []*Output{NewOutput(address, hasIncome, true, value)}
	return newTransaction([]*Input{}, outputs, timestamp)
}

func NewTransactionFromRequest(transactionRequest *network.TransactionRequest) (*Transaction, error) {
	var inputs []*Input
	for _, inputRequest := range *transactionRequest.Inputs {
		input, err := NewInput(*inputRequest.OutputIndex, *inputRequest.TransactionId, *inputRequest.PublicKey, *inputRequest.Signature)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}
	var outputs []*Output
	for _, output := range *transactionRequest.Outputs {
		outputs = append(outputs, NewOutput(*output.Address, *output.HasIncome, *output.HasReward, *output.Value))
	}
	transaction, err := newTransaction(inputs, outputs, *transactionRequest.Timestamp)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func newTransaction(inputs []*Input, outputs []*Output, timestamp int64) (transaction *Transaction, err error) {
	transaction = &Transaction{"", inputs, outputs, timestamp, false, "", 0}
	if transaction.id, err = transaction.generateId(); err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}
	if err = transaction.findReward(); err != nil {
		return nil, fmt.Errorf("failed to find reward: %w", err)
	}
	return transaction, nil
}

func (transaction *Transaction) Equals(other *Transaction) bool {
	return transaction.id == other.Id()
}

func (transaction *Transaction) UnmarshalJSON(data []byte) error {
	var dto *transactionDto
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}
	transactionFromDto, err := newTransaction(dto.Inputs, dto.Outputs, dto.Timestamp)
	if err != nil {
		return err
	}
	if transactionFromDto.Id() != dto.Id {
		return fmt.Errorf("wrong transaction ID, provided: %s, calculated: %s", dto.Id, transaction.id)
	}
	transaction.id = transactionFromDto.Id()
	transaction.inputs = transactionFromDto.Inputs()
	transaction.outputs = transactionFromDto.Outputs()
	transaction.timestamp = transactionFromDto.Timestamp()
	transaction.hasReward = transactionFromDto.HasReward()
	transaction.rewardRecipientAddress = transactionFromDto.RewardRecipientAddress()
	transaction.rewardValue = transactionFromDto.RewardValue()
	return nil
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(transactionDto{
		Id:        transaction.id,
		Inputs:    transaction.inputs,
		Outputs:   transaction.outputs,
		Timestamp: transaction.timestamp,
	})
}

func (transaction *Transaction) marshalJSONWithoutId() ([]byte, error) {
	return json.Marshal(struct {
		Inputs    []*Input  `json:"inputs"`
		Outputs   []*Output `json:"outputs"`
		Timestamp int64     `json:"timestamp"`
	}{
		Inputs:    transaction.inputs,
		Outputs:   transaction.outputs,
		Timestamp: transaction.timestamp,
	})
}

func (transaction *Transaction) VerifySignatures() error {
	for _, input := range transaction.inputs {
		if err := input.VerifySignature(); err != nil {
			return fmt.Errorf("failed to verify signature for input: %v\n %w", input, err)
		}
	}
	return nil
}

func (transaction *Transaction) Fee(genesisTimestamp int64, settings config.Settings, timestamp int64, validationTimestamp int64, utxoFinder protocol.UtxoFinder) (uint64, error) {
	incomeBase := settings.IncomeBaseInParticles
	incomeLimit := settings.IncomeLimitInParticles
	var inputsValue uint64
	var outputsValue uint64
	for _, input := range transaction.inputs {
		utxo, err := utxoFinder(input)
		if err != nil {
			return 0, err
		}
		value := utxo.Value(timestamp, genesisTimestamp, settings.HalfLifeInNanoseconds, incomeBase, incomeLimit, validationTimestamp)
		inputsValue += value
	}
	for _, output := range transaction.outputs {
		outputsValue += output.InitialValue()
	}
	if inputsValue < outputsValue {
		return 0, errors.New("transaction fee is negative")
	}
	fee := inputsValue - outputsValue
	minimalTransactionFee := settings.MinimalTransactionFee
	if fee < minimalTransactionFee {
		return 0, fmt.Errorf("transaction fee is too low, fee: %d, minimal fee: %d", fee, minimalTransactionFee)
	}
	return fee, nil
}

func (transaction *Transaction) Id() string {
	return transaction.id
}

func (transaction *Transaction) Inputs() []*Input {
	return transaction.inputs
}

func (transaction *Transaction) Outputs() []*Output {
	return transaction.outputs
}

func (transaction *Transaction) HasReward() bool {
	return transaction.hasReward
}

func (transaction *Transaction) RewardRecipientAddress() string {
	return transaction.rewardRecipientAddress
}

func (transaction *Transaction) RewardValue() uint64 {
	return transaction.rewardValue
}

func (transaction *Transaction) Timestamp() int64 {
	return transaction.timestamp
}

func (transaction *Transaction) generateId() (string, error) {
	marshaledTransaction, err := transaction.marshalJSONWithoutId()
	if err != nil {
		return "", errors.New("failed to marshal transaction")
	}
	transactionHash := sha256.Sum256(marshaledTransaction)
	id := fmt.Sprintf("%x", transactionHash)
	return id, nil
}

func (transaction *Transaction) findReward() error {
	for _, output := range transaction.outputs {
		if output == nil {
			return errors.New("an output is nil")
		}
		if output.HasReward() {
			if transaction.hasReward {
				return errors.New("multiple rewards attempt for the same transaction")
			}
			transaction.hasReward = true
			transaction.rewardRecipientAddress = output.Address()
			transaction.rewardValue = output.InitialValue()
		}
	}
	return nil
}
