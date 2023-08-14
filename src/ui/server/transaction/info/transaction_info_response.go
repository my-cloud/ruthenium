package info

type TransactionInfoResponse struct {
	Rest  uint64          `json:"rest"`
	Utxos []*UtxoResponse `json:"utxos"`
}
