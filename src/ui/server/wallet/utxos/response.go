package utxos

import "github.com/my-cloud/ruthenium/src/node/network"

type Response struct {
	BlockHeight int                             `json:"block_height"`
	HasIncome   bool                            `json:"has_income"`
	Rest        uint64                          `json:"rest"`
	Utxos       []*network.WalletOutputResponse `json:"utxos"`
}
