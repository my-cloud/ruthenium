package info

import "github.com/my-cloud/ruthenium/src/node/protocol/verification"

type TransactionInfo struct {
	Inputs    []*verification.InputInfo `json:"inputs"`
	Rest      uint64                    `json:"rest"`
	Timestamp int64                     `json:"timestamp"`
}
