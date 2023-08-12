package info

import "github.com/my-cloud/ruthenium/src/node/network"

type TransactionInfoResponse struct {
	Rest  uint64                  `json:"rest"`
	Utxos []*network.UtxoResponse `json:"utxos"`
}
