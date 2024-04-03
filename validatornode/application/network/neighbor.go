package network

import (
	"encoding/json"
	"fmt"

	"github.com/my-cloud/ruthenium/validatornode/presentation"
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
	return neighbor.sendRequest(presentation.BlocksEndpoint, startingBlockHeight)
}

func (neighbor *Neighbor) GetFirstBlockTimestamp() (int64, error) {
	res, err := neighbor.sender.Send(presentation.FirstBlockTimestampEndpoint, []byte{})
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
	return neighbor.sendRequest(presentation.SettingsEndpoint, []byte{})
}

func (neighbor *Neighbor) SendTargets(targets []string) error {
	_, err := neighbor.sendRequest(presentation.TargetsEndpoint, targets)
	return err
}

func (neighbor *Neighbor) AddTransaction(transaction []byte) error {
	_, err := neighbor.sender.Send(presentation.TransactionEndpoint, transaction)
	return err
}

func (neighbor *Neighbor) GetTransactions() ([]byte, error) {
	return neighbor.sender.Send(presentation.TransactionsEndpoint, []byte{})
}

func (neighbor *Neighbor) GetUtxos(address string) ([]byte, error) {
	return neighbor.sendRequest(presentation.UtxosEndpoint, address)
}

func (neighbor *Neighbor) sendRequest(topic string, request interface{}) ([]byte, error) {
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return neighbor.sender.Send(topic, bytes)
}
