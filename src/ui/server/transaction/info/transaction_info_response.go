package info

type TransactionInfoResponse struct {
	Rest      uint64          `json:"rest"`
	Utxos     []*UtxoResponse `json:"utxos"`
	Timestamp int64           `json:"timestamp"`
}
