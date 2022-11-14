package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
)

type Neighbor struct {
	target *Target
	client Client
	logger *log.Logger
}

func NewNeighbor(target *Target, clientFactory ClientFactory, logger *log.Logger) (*Neighbor, error) {
	client, err := clientFactory.CreateClient(target.Ip(), target.Port(), target.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to create client reaching %s: %w", target.Value(), err)
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

func (neighbor *Neighbor) GetBlocks() (blockResponses []*network.BlockResponse, err error) {
	res, err := neighbor.sendRequest(GetBlocksRequest)
	if err == nil {
		err = res.GetGob(&blockResponses)
	}
	return
}

func (neighbor *Neighbor) SendTargets(request []network.TargetRequest) (err error) {
	_, err = neighbor.sendRequest(request)
	return
}

func (neighbor *Neighbor) AddTransaction(request network.TransactionRequest) (err error) {
	_, err = neighbor.sendRequest(request)
	return
}

func (neighbor *Neighbor) GetTransactions() (transactionResponses []network.TransactionResponse, err error) {
	res, err := neighbor.sendRequest(GetTransactionsRequest)
	if err != nil {
		return
	}
	err = res.GetGob(&transactionResponses)
	if transactionResponses == nil {
		return []network.TransactionResponse{}, err
	}
	return
}

func (neighbor *Neighbor) GetAmount(request network.AmountRequest) (amountResponse *network.AmountResponse, err error) {
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
	res, err = neighbor.client.Send("dialog", req)
	return
}
