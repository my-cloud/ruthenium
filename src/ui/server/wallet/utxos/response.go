package utxos

import "github.com/my-cloud/ruthenium/src/node/network"

type Response struct {
	Rest  uint64                  `json:"rest"`
	Utxos []*network.UtxoResponse `json:"utxos"`
}
