package status

type Progress struct {
	CurrentBlockTimestamp int64  `json:"current_block_timestamp"`
	TransactionStatus     string `json:"transaction_status"`
	ValidationTimestamp   int64  `json:"validation_timestamp"`
}
