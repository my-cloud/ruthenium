package p2p

import (
	"encoding/json"
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network"
	"time"
)

const connectionTimeoutInSeconds = 5

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
	settings.SetConnTimeout(connectionTimeoutInSeconds * time.Second)
	client.SetSettings(settings)
	return &Neighbor{target, client, settings}, nil
}

func (neighbor *Neighbor) Target() string {
	return neighbor.target.Value()
}

func (neighbor *Neighbor) GetBlock(blockHeight uint64) (blockResponse *network.BlockResponse, err error) {
	request := network.BlockRequest{BlockHeight: &blockHeight}
	res, err := neighbor.sendRequest("block", request)
	if err == nil {
		data := res.GetBytes()
		err = json.Unmarshal(data, &blockResponse)
	}
	return
}

func (neighbor *Neighbor) GetBlocks(startingBlockHeight uint64) (blockResponses []*network.BlockResponse, err error) {
	request := network.BlocksRequest{StartingBlockHeight: &startingBlockHeight}
	res, err := neighbor.sendRequest("blocks", request)
	if err == nil {
		data := res.GetBytes()
		err = json.Unmarshal(data, &blockResponses)
	}
	return
}

func (neighbor *Neighbor) SendTargets(request []network.TargetRequest) (err error) {
	_, err = neighbor.sendRequest("targets", request)
	return
}

func (neighbor *Neighbor) AddTransaction(request network.TransactionRequest) (err error) {
	_, err = neighbor.sendRequest("transaction", request)
	return
}

func (neighbor *Neighbor) GetTransactions() (transactionResponses []network.TransactionResponse, err error) {
	res, err := neighbor.sendRequest("transactions", GetTransactions)
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

func (neighbor *Neighbor) GetUtxos(address string) (utxos []*network.UtxoResponse, err error) {
	request := network.UtxosRequest{Address: &address}
	res, err := neighbor.sendRequest("utxos", request)
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

func (neighbor *Neighbor) sendRequest(topic string, request interface{}) (res gp2p.Data, err error) {
	req := gp2p.Data{}
	data, err := json.Marshal(request)
	if err != nil {
		return
	}
	req.SetBytes(data)
	res, err = neighbor.client.Send(topic, req)
	return
}
