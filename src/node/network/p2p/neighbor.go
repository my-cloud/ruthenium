package p2p

import (
	"encoding/json"
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/node/network"
	"time"
)

const connectionTimeoutInSeconds = 5 // TODO calculate from validation timestamp

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

func (neighbor *Neighbor) GetBlocks(startingBlockHeight uint64) (blocks []byte, err error) {
	res, err := neighbor.sendRequest(blocksEndpoint, startingBlockHeight)
	if err == nil {
		return res.GetBytes(), nil
	}
	return
}

func (neighbor *Neighbor) GetFirstBlockTimestamp() (timestamp int64, err error) {
	res, err := neighbor.client.Send(firstBlockTimestampEndpoint, gp2p.Data{})
	if err == nil {
		timestampBytes := res.GetBytes()
		err = json.Unmarshal(timestampBytes, &timestamp)
		if err != nil {
			return
		}
	}
	return
}

func (neighbor *Neighbor) SendTargets(targets []string) (err error) {
	_, err = neighbor.sendRequest(targetsEndpoint, targets)
	return
}

func (neighbor *Neighbor) AddTransaction(request network.TransactionRequest) (err error) {
	_, err = neighbor.sendRequest(transactionEndpoint, request)
	return
}

func (neighbor *Neighbor) GetTransactions() (transactionResponses []byte, err error) {
	res, err := neighbor.client.Send(transactionsEndpoint, gp2p.Data{})
	if err == nil {
		return res.GetBytes(), nil
	}
	return
}

func (neighbor *Neighbor) GetUtxos(address string) (utxos []byte, err error) {
	res, err := neighbor.sendRequest(utxosEndpoint, address)
	if err == nil {
		return res.GetBytes(), nil
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
	res, err = neighbor.client.Send(topic, req) // FIXME panic
	return
}
