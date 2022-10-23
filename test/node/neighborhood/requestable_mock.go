package neighborhood

import (
	"github.com/my-cloud/ruthenium/src/api/node"
)

type RequestableMock struct {
}

func NewRequestableMock() *RequestableMock {
	return &RequestableMock{}
}

var IpMock func() string
var PortMock func() uint16
var TargetMock func() string
var IsFoundMock func() bool
var GetBlocksMock func() ([]*node.BlockResponse, error)
var SendTargetsMock func([]node.TargetRequest) error
var AddTransactionMock func(node.TransactionRequest) error
var GetTransactionsMock func() ([]node.TransactionResponse, error)
var GetAmountMock func(node.AmountRequest) (*node.AmountResponse, error)
var MineMock func() error
var StartMiningMock func() error
var StopMiningMock func() error

func (mock *RequestableMock) Ip() string {
	return IpMock()
}

func (mock *RequestableMock) Port() uint16 {
	return PortMock()
}

func (mock *RequestableMock) Target() string {
	return TargetMock()
}

func (mock *RequestableMock) IsFound() bool {
	return IsFoundMock()
}

func (mock *RequestableMock) GetBlocks() ([]*node.BlockResponse, error) {
	return GetBlocksMock()
}

func (mock *RequestableMock) SendTargets(request []node.TargetRequest) error {
	return SendTargetsMock(request)
}

func (mock *RequestableMock) AddTransaction(request node.TransactionRequest) error {
	return AddTransactionMock(request)
}

func (mock *RequestableMock) GetTransactions() ([]node.TransactionResponse, error) {
	return GetTransactionsMock()
}

func (mock *RequestableMock) GetAmount(request node.AmountRequest) (*node.AmountResponse, error) {
	return GetAmountMock(request)
}

func (mock *RequestableMock) Mine() error {
	return MineMock()
}

func (mock *RequestableMock) StartMining() error {
	return StartMiningMock()
}

func (mock *RequestableMock) StopMining() error {
	return StopMiningMock()
}
