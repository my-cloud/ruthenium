package ledger

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
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

func NewRewardTransaction(address string, isYielding bool, timestamp int64, value uint64) (*Transaction, error) {
	outputs := []*Output{NewOutput(address, isYielding, value)}
	var inputs []*Input
	id, err := generateId(inputs, outputs, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}
	return &Transaction{id, inputs, outputs, timestamp, true, address, value}, nil
}

func (transaction *Transaction) Equals(other *Transaction) bool {
	return transaction.id == other.Id()
}

func (transaction *Transaction) UnmarshalJSON(data []byte) error {
	var dto *transactionDto
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}
	id, err := generateId(dto.Inputs, dto.Outputs, dto.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to generate id: %w", err)
	}
	if id != dto.Id {
		return fmt.Errorf("wrong transaction ID, provided: %s, calculated: %s", dto.Id, id)
	}
	if len(dto.Inputs) == 0 {
		if len(dto.Outputs) > 1 {
			return errors.New("multiple rewards attempt for the same transaction")
		} else if len(dto.Outputs) == 0 {
			return errors.New("reward not found whereas the transaction has no input")
		}
		transaction.hasReward = true
		transaction.rewardRecipientAddress = dto.Outputs[0].Address()
		transaction.rewardValue = dto.Outputs[0].InitialValue()
	}
	transaction.id = dto.Id
	transaction.inputs = dto.Inputs
	transaction.outputs = dto.Outputs
	transaction.timestamp = dto.Timestamp
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

func (transaction *Transaction) VerifySignatures() error {
	for _, input := range transaction.inputs {
		if err := input.VerifySignature(); err != nil {
			return fmt.Errorf("failed to verify signature of an input: %w", err)
		}
	}
	return nil
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

func generateId(inputs []*Input, outputs []*Output, timestamp int64) (string, error) {
	marshaledTransaction, err := json.Marshal(struct {
		Inputs    []*Input  `json:"inputs"`
		Outputs   []*Output `json:"outputs"`
		Timestamp int64     `json:"timestamp"`
	}{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: timestamp,
	})
	if err != nil {
		return "", errors.New("failed to marshal transaction")
	}
	transactionHash := sha256.Sum256(marshaledTransaction)
	id := fmt.Sprintf("%x", transactionHash)
	return id, nil
}
