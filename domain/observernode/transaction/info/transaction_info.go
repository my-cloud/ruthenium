package info

import (
	"github.com/my-cloud/ruthenium/domain"
)

type TransactionInfo struct {
	Inputs    []*domain.InputInfo `json:"inputs"`
	Rest      uint64              `json:"rest"`
	Timestamp int64               `json:"timestamp"`
}
