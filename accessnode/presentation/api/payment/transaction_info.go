package payment

import (
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
)

type TransactionInfo struct {
	Inputs    []*protocol.InputInfo `json:"inputs"`
	Rest      uint64                `json:"rest"`
	Timestamp int64                 `json:"timestamp"`
}
