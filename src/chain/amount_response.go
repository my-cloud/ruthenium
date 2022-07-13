package chain

import "encoding/json"

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func NewAmountResponse(amount float32) *AmountResponse {
	return &AmountResponse{amount}
}

func (amountResponse *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: amountResponse.Amount,
	})
}
