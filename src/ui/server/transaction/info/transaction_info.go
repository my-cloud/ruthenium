package info

import "github.com/my-cloud/ruthenium/src/node/protocol/verification"

type TransactionInfo struct {
	Rest      uint64                    `json:"rest"`
	Utxos     []*verification.InputInfo `json:"utxos"`
	Timestamp int64                     `json:"timestamp"`
}
