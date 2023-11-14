package domain

import (
	"encoding/json"
)

type inputInfoDto struct {
	OutputIndex   uint16 `json:"output_index"`
	TransactionId string `json:"transaction_id"`
}

type InputInfo struct {
	outputIndex   uint16
	transactionId string
}

func NewInputInfo(outputIndex uint16, transactionId string) *InputInfo {
	return &InputInfo{outputIndex, transactionId}
}

func (inputInfo *InputInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(inputInfoDto{
		OutputIndex:   inputInfo.outputIndex,
		TransactionId: inputInfo.transactionId,
	})
}

func (inputInfo *InputInfo) UnmarshalJSON(data []byte) error {
	var dto *inputInfoDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	inputInfo.outputIndex = dto.OutputIndex
	inputInfo.transactionId = dto.TransactionId
	return nil
}

func (inputInfo *InputInfo) OutputIndex() uint16 {
	return inputInfo.outputIndex
}

func (inputInfo *InputInfo) TransactionId() string {
	return inputInfo.transactionId
}
