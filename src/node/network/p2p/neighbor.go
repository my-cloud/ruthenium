package p2p

import (
	"encoding/json"
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
	client, err := clientFactory.CreateClient(target.Ip(), target.Port())
	if err != nil {
		return nil, fmt.Errorf("failed to create client reaching %s: %w", target.Value(), err)
	}
	settings := gp2p.NewClientSettings()
	settings.SetRetry(1, time.Nanosecond)
	settings.SetConnTimeout(commonConnectionTimeoutInSeconds * time.Second)
	client.SetSettings(settings)
	return &Neighbor{target, client, settings}, nil
}

func (neighbor *Neighbor) Target() string {
	return neighbor.target.Value()
}

func (neighbor *Neighbor) GetBlock(blockHeight uint64) (blockResponse *network.BlockResponse, err error) {
	request := network.BlockRequest{BlockHeight: &blockHeight}
	res, err := neighbor.sendRequest(request)
	if err == nil {
		data := res.GetBytes()
		err = json.Unmarshal(data, &blockResponse)
	}
	return
}

func (neighbor *Neighbor) GetBlocks() (blockResponses []*network.BlockResponse, err error) {
	neighbor.settings.SetConnTimeout(initializationConnectionTimeoutInSeconds * time.Second)
	neighbor.client.SetSettings(neighbor.settings)
	res, err := neighbor.sendRequest(GetBlocks)
	if err == nil {
		data := res.GetBytes()
		err = json.Unmarshal(data, &blockResponses)
	}
	neighbor.settings.SetConnTimeout(commonConnectionTimeoutInSeconds * time.Second)
	neighbor.client.SetSettings(neighbor.settings)
	return
}

func (neighbor *Neighbor) GetLambda() (lambda float64, err error) {
	res, err := neighbor.sendRequest(GetLambda)
	if err == nil {
		data := res.GetBytes()
		err = json.Unmarshal(data, &lambda)
	}
	return
}

func (neighbor *Neighbor) GetLastBlocks(startingBlockHeight uint64) (blockResponses []*network.BlockResponse, err error) {
	request := network.LastBlocksRequest{StartingBlockHeight: &startingBlockHeight}
	res, err := neighbor.sendRequest(request)
	if err == nil {
		data := res.GetBytes()
		err = json.Unmarshal(data, &blockResponses)
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
	res, err := neighbor.sendRequest(GetTransactions)
	if err != nil {
		return
	}
	data := res.GetBytes()
	err = json.Unmarshal(data, &transactionResponses)
	if transactionResponses == nil {
		return []network.TransactionResponse{}, err
	}
	return
}

func (neighbor *Neighbor) GetUtxos(address string) (utxos []*network.WalletOutputResponse, err error) {
	request := network.UtxosRequest{Address: &address}
	res, err := neighbor.sendRequest(request)
	if err != nil {
		return
	}
	data := res.GetBytes()
	err = json.Unmarshal(data, &utxos)
	if err != nil {
		return
	}
	return
}

func (neighbor *Neighbor) sendRequest(request interface{}) (res gp2p.Data, err error) {
	req := gp2p.Data{}
	data, err := json.Marshal(request)
	if err != nil {
		return
	}
	req.SetBytes(data)
	res, err = neighbor.client.Send("dialog", req)
	return
}
