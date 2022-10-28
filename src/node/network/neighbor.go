package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
)

const (
	GetBlocksRequest       = "GET BLOCKS REQUEST"
	GetTransactionsRequest = "GET TRANSACTIONS REQUEST"
	MineRequest            = "MINE REQUEST"
	StartMiningRequest     = "START MINING REQUEST"
	StopMiningRequest      = "STOP MINING REQUEST"
)

type Neighbor struct {
	target *Target
	sender Sender
	logger *log.Logger
}

func NewNeighbor(target *Target, senderFactory SenderFactory, logger *log.Logger) (*Neighbor, error) {
	client, err := senderFactory.CreateSender(target.Ip(), target.Port(), target.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to start client reaching %s: %w", target.Value(), err)
	}
	return &Neighbor{target, client, logger}, nil
}

func (neighbor *Neighbor) Ip() string {
	return neighbor.target.ip
}

func (neighbor *Neighbor) Port() uint16 {
	return neighbor.target.port
}

func (neighbor *Neighbor) Target() string {
	return neighbor.target.Value()
}

func (neighbor *Neighbor) GetBlocks() (blockResponses []*neighborhood.BlockResponse, err error) {
	res, err := neighbor.sendRequest(GetBlocksRequest)
	if err == nil {
		err = res.GetGob(&blockResponses)
	}
	return
}

func (neighbor *Neighbor) SendTargets(request []neighborhood.TargetRequest) (err error) {
	_, err = neighbor.sendRequest(request)
	return
}

func (neighbor *Neighbor) AddTransaction(request neighborhood.TransactionRequest) (err error) {
	_, err = neighbor.sendRequest(request)
	return
}

func (neighbor *Neighbor) GetTransactions() (transactionResponses []neighborhood.TransactionResponse, err error) {
	res, err := neighbor.sendRequest(GetTransactionsRequest)
	if err != nil {
		return
	}
	err = res.GetGob(&transactionResponses)
	if transactionResponses == nil {
		return []neighborhood.TransactionResponse{}, err
	}
	return
}

func (neighbor *Neighbor) GetAmount(request neighborhood.AmountRequest) (amountResponse *neighborhood.AmountResponse, err error) {
	res, err := neighbor.sendRequest(request)
	if err == nil {
		err = res.GetGob(&amountResponse)
	}
	return
}

func (neighbor *Neighbor) Mine() (err error) {
	_, err = neighbor.sendRequest(MineRequest)
	return
}

func (neighbor *Neighbor) StartMining() (err error) {
	_, err = neighbor.sendRequest(StartMiningRequest)
	return
}

func (neighbor *Neighbor) StopMining() (err error) {
	_, err = neighbor.sendRequest(StopMiningRequest)
	return
}

func (neighbor *Neighbor) sendRequest(request interface{}) (res p2p.Data, err error) {
	req := p2p.Data{}
	err = req.SetGob(request)
	if err != nil {
		return
	}
	res, err = neighbor.sender.Send("dialog", req)
	return
}