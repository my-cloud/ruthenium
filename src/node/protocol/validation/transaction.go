package validation

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Transaction struct {
	id                     [32]byte
	inputs                 []*network.InputResponse
	outputs                []*network.OutputResponse
	timestamp              int64
	hasReward              bool
	rewardRecipientAddress string
	rewardValue            uint64
}

func NewRewardTransaction(address string, blockHeight int, timestamp int64, value uint64) (*network.TransactionResponse, error) {
	outputs := []*network.OutputResponse{
		{
			Address:     address,
			BlockHeight: blockHeight,
			HasReward:   true,
			HasIncome:   true,
			Value:       value,
		},
	}
	transaction, err := newTransaction([]*network.InputResponse{}, outputs, timestamp)
	if err != nil {
		return nil, err
	}
	return transaction.GetResponse(), nil
}

func NewTransactionFromRequest(transactionRequest *network.TransactionRequest) (*Transaction, error) {
	//address := *transactionRequest.SenderAddress
	//isRegistered, err := registry.IsRegistered(address)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get proof of humanity: %w", err)
	//}
	//
	//var inputs []*Input
	//var inputsValue uint64
	//// for _, utxo := range utxos {
	//// TODO if isRegistered then use all utxo, else select only some to have the smallest byte size
	//input, err := NewInput(0, [32]byte{}, *transactionRequest.SenderPublicKey, *transactionRequest.Signature)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to instantiate input: %w", err)
	//}
	//inputs = append(inputs, input)
	////inputsValue += utxo.Value
	//// }
	//
	//var outputs []*Output
	//transactionRequestValue := *transactionRequest.Value
	//output := NewOutput(*transactionRequest.RecipientAddress, blockHeight, false, false, transactionRequestValue)
	//outputs = append(outputs, output)
	//surplus := NewOutput(address, blockHeight, false, isRegistered, inputsValue-transactionRequestValue)
	//outputs = append(outputs, surplus)
	var inputs []*network.InputResponse
	for _, input := range *transactionRequest.Inputs {
		inputs = append(inputs, &network.InputResponse{OutputIndex: *input.OutputIndex, TransactionId: *input.TransactionId, PublicKey: *input.PublicKey, Signature: *input.Signature})
	}
	var outputs []*network.OutputResponse
	for _, output := range *transactionRequest.Outputs {
		outputs = append(outputs, &network.OutputResponse{Address: *output.Address, BlockHeight: *output.BlockHeight, HasReward: *output.HasReward, HasIncome: *output.HasIncome, Value: *output.Value})
	}
	transaction, err := newTransaction(inputs, outputs, *transactionRequest.Timestamp)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func newTransaction(inputs []*network.InputResponse, outputs []*network.OutputResponse, timestamp int64) (*Transaction, error) {
	transaction := &Transaction{[32]byte{}, inputs, outputs, timestamp, false, "", 0}
	if err := transaction.generateId(); err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}
	if err := transaction.findReward(); err != nil {
		return nil, fmt.Errorf("failed to find reward: %w", err)
	}
	return transaction, nil
}

func NewTransactionFromResponse(transactionResponse *network.TransactionResponse) (*Transaction, error) {
	transaction, err := newTransaction(transactionResponse.Inputs, transactionResponse.Outputs, transactionResponse.Timestamp)
	if err != nil {
		return nil, err
	}
	if !transaction.Equals(transactionResponse) {
		return nil, fmt.Errorf("wrong transaction ID, provided: %s, calculated: %s", transactionResponse.Id, transaction.id)
	}
	return transaction, nil
}

func (transaction *Transaction) Equals(other *network.TransactionResponse) bool {
	return transaction.id == other.Id
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Inputs    []*network.InputResponse  `json:"inputs"`
		Outputs   []*network.OutputResponse `json:"outputs"`
		Timestamp int64                     `json:"timestamp"`
	}{
		Inputs:    transaction.inputs,
		Outputs:   transaction.outputs,
		Timestamp: transaction.timestamp,
	})
}

func (transaction *Transaction) GetResponse() *network.TransactionResponse {
	return &network.TransactionResponse{
		Id:        transaction.id,
		Inputs:    transaction.inputs,
		Outputs:   transaction.outputs,
		Timestamp: transaction.timestamp,
	}
}

func (transaction *Transaction) VerifySignatures() error {
	for _, inputResponse := range transaction.inputs {
		input, err := NewInputFromResponse(inputResponse)
		if err != nil {
			return fmt.Errorf("failed to instantiate input: %v: %w", input, err)
		}
		if err = input.VerifySignature(); err != nil {
			return fmt.Errorf("failed to verify signature for input: %v: %w", input, err)
		}
	}
	return nil
}

func (transaction *Transaction) Id() [32]byte {
	return transaction.id
}

func (transaction *Transaction) Inputs() []*network.InputResponse {
	return transaction.inputs
}

func (transaction *Transaction) Outputs() []*network.OutputResponse {
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

func (transaction *Transaction) generateId() error {
	marshaledTransaction, err := transaction.MarshalJSON()
	if err != nil {
		return errors.New("failed to marshal transaction")
	}
	transaction.id = sha256.Sum256(marshaledTransaction)
	return nil
}

func (transaction *Transaction) findReward() error {
	for _, output := range transaction.outputs {
		if output.HasReward {
			if transaction.hasReward {
				return errors.New("multiple rewards attempt for the same transaction")
			}
			transaction.hasReward = true
			transaction.rewardRecipientAddress = output.Address
			transaction.rewardValue = output.Value
		}
	}
	return nil
}
