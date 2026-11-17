package payment

import (
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
)

type TransactionInfo struct {
	Fee       uint64              `json:"fee"`
	Inputs    []*ledger.InputInfo `json:"inputs"`
	Rest      uint64              `json:"rest"`
	Timestamp int64               `json:"timestamp"`
}
