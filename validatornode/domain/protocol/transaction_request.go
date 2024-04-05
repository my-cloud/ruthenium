package protocol

import (
	"encoding/json"
)

type transactionRequestDto struct {
	Transaction                  *Transaction
	TransactionBroadcasterTarget string
}

type TransactionRequest struct {
	transaction                  *Transaction
	transactionBroadcasterTarget string
}

func NewTransactionRequest(transaction *Transaction, transactionBroadcasterTarget string) *TransactionRequest {
	return &TransactionRequest{transaction, transactionBroadcasterTarget}
}

func (request *TransactionRequest) UnmarshalJSON(data []byte) error {
	var dto *transactionRequestDto
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}
	request.transaction = dto.Transaction
	request.transactionBroadcasterTarget = dto.TransactionBroadcasterTarget
	return nil
}

func (request *TransactionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(transactionRequestDto{
		Transaction:                  request.transaction,
		TransactionBroadcasterTarget: request.transactionBroadcasterTarget,
	})
}

func (request *TransactionRequest) Transaction() *Transaction {
	return request.transaction
}

func (request *TransactionRequest) TransactionBroadcasterTarget() string {
	return request.transactionBroadcasterTarget
}
