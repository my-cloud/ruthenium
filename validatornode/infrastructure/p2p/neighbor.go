package p2p

import (
	"encoding/json"
	"time"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

const (
	BlocksEndpoint              = "blocks"
	FirstBlockTimestampEndpoint = "first-block-timestamp"
	SettingsEndpoint            = "settings"
	TargetsEndpoint             = "targets"
	TransactionEndpoint         = "transaction"
	TransactionsEndpoint        = "transactions"
	UtxosEndpoint               = "utxos"
)

type Neighbor struct {
	*gp2p.Client
	target *network.Target
}

func NewNeighbor(ip string, port string, connectionTimeout time.Duration, logger log.Logger) (*Neighbor, error) {
	tcp := gp2p.NewTCP(ip, port)
	client, err := gp2p.NewClient(tcp)
	if err != nil {
		return nil, err
	}
	settings := gp2p.NewClientSettings()
	settings.SetRetry(1, time.Nanosecond)
	settings.SetConnTimeout(connectionTimeout)
	client.SetSettings(settings)
	client.SetLogger(logger)
	target := network.NewTarget(ip, port)
	return &Neighbor{client, target}, err
}

func (neighbor *Neighbor) Target() string {
	return neighbor.target.Value()
}

func (neighbor *Neighbor) GetBlocks(startingBlockHeight uint64) ([]byte, error) {
	return neighbor.sendRequest(BlocksEndpoint, startingBlockHeight)
}

func (neighbor *Neighbor) GetFirstBlockTimestamp() (int64, error) {
	res, err := neighbor.sendRequestBytes(FirstBlockTimestampEndpoint, []byte{})
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
	return neighbor.sendRequestBytes(SettingsEndpoint, []byte{})
}

func (neighbor *Neighbor) SendTargets(targets []string) error {
	_, err := neighbor.sendRequest(TargetsEndpoint, targets)
	return err
}

func (neighbor *Neighbor) AddTransaction(transaction []byte) error {
	_, err := neighbor.sendRequestBytes(TransactionEndpoint, transaction)
	return err
}

func (neighbor *Neighbor) GetTransactions() ([]byte, error) {
	return neighbor.sendRequestBytes(TransactionsEndpoint, []byte{})
}

func (neighbor *Neighbor) GetUtxos(address string) ([]byte, error) {
	return neighbor.sendRequest(UtxosEndpoint, address)
}

func (neighbor *Neighbor) sendRequest(topic string, request interface{}) ([]byte, error) {
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return neighbor.sendRequestBytes(topic, bytes)
}

func (neighbor *Neighbor) sendRequestBytes(topic string, request []byte) ([]byte, error) {
	data, err := neighbor.Client.Send(topic, gp2p.Data{Bytes: request})
	if err != nil {
		return []byte{}, err
	}
	return data.Bytes, err
}
