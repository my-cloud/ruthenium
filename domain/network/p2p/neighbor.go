package p2p

import (
	"encoding/json"
	"fmt"
)

type Neighbor struct {
	target *Target
	sender Sender
}

func NewNeighbor(target *Target, clientFactory SenderCreator) (*Neighbor, error) {
	sender, err := clientFactory.CreateSender(target.Ip(), target.Port())
	if err != nil {
		return nil, fmt.Errorf("failed to create sender reaching %s: %w", target.Value(), err)
	}
	return &Neighbor{target, sender}, nil
}

func (neighbor *Neighbor) Target() string {
	return neighbor.target.Value()
}

func (neighbor *Neighbor) GetBlocks(startingBlockHeight uint64) ([]byte, error) {
	return neighbor.sendRequest(blocksEndpoint, startingBlockHeight)
}

func (neighbor *Neighbor) GetFirstBlockTimestamp() (int64, error) {
	res, err := neighbor.sender.Send(firstBlockTimestampEndpoint, []byte{})
	var timestamp int64
	if err != nil {
		return timestamp, err
	}
	err = json.Unmarshal(res, &timestamp)
	if err != nil {
		return timestamp, err
	}
	return timestamp, err
}

func (neighbor *Neighbor) GetSettings() ([]byte, error) {
	return neighbor.sendRequest(settingsEndpoint, []byte{})
}

func (neighbor *Neighbor) SendTargets(targets []string) error {
	_, err := neighbor.sendRequest(targetsEndpoint, targets)
	return err
}

func (neighbor *Neighbor) AddTransaction(transaction []byte) error {
	_, err := neighbor.sender.Send(transactionEndpoint, transaction)
	return err
}

func (neighbor *Neighbor) GetTransactions() ([]byte, error) {
	return neighbor.sender.Send(transactionsEndpoint, []byte{})
}

func (neighbor *Neighbor) GetUtxos(address string) ([]byte, error) {
	return neighbor.sendRequest(utxosEndpoint, address)
}

func (neighbor *Neighbor) sendRequest(topic string, request interface{}) ([]byte, error) {
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return neighbor.sender.Send(topic, bytes)
}
