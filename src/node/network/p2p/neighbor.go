package p2p

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network"
	"time"
)

const (
	initializationConnectionTimeoutInSeconds = 600
	commonConnectionTimeoutInSeconds         = 5
)

type Neighbor struct {
	target   *Target
	client   Client
	settings *gp2p.ClientSettings
}

func NewNeighbor(target *Target, clientFactory ClientFactory) (*Neighbor, error) {
	client, err := clientFactory.CreateClient(target.Ip(), target.Port(), target.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to create client reaching %s: %w", target.Value(), err)
	}
	settings := gp2p.NewClientSettings()
	settings.SetRetry(1, time.Nanosecond)
	settings.SetConnTimeout(commonConnectionTimeoutInSeconds * time.Second)
	client.SetSettings(settings)
	return &Neighbor{target, client, settings}, nil
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
	neighbor.settings.SetConnTimeout(initializationConnectionTimeoutInSeconds * time.Second)
	neighbor.client.SetSettings(neighbor.settings)
	res, err := neighbor.sendRequest(GetBlocks)
	if err == nil {
		err = res.GetGob(&blockResponses)
	}
	neighbor.settings.SetConnTimeout(commonConnectionTimeoutInSeconds * time.Second)
	neighbor.client.SetSettings(neighbor.settings)
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
	res, err := neighbor.sendRequest(GetTransactions)
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

func (neighbor *Neighbor) StartValidation() (err error) {
	_, err = neighbor.sendRequest(StartValidation)
	return
}

func (neighbor *Neighbor) StopValidation() (err error) {
	_, err = neighbor.sendRequest(StopValidation)
	return
}

func (neighbor *Neighbor) sendRequest(request interface{}) (res gp2p.Data, err error) {
	req := gp2p.Data{}
	err = req.SetGob(request)
	if err != nil {
		return
	}
	res, err = neighbor.client.Send("dialog", req)
	return
}
