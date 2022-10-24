package network

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/node/network"
	"github.com/my-cloud/ruthenium/src/log"
	"strconv"
)

const (
	GetBlocksRequest       = "GET BLOCKS REQUEST"
	GetTransactionsRequest = "GET TRANSACTIONS REQUEST"
	MineRequest            = "MINE REQUEST"
	StartMiningRequest     = "START MINING REQUEST"
	StopMiningRequest      = "STOP MINING REQUEST"

	NeighborFindingTimeoutSecond = 5
)

type Neighbor struct {
	ip     string
	port   uint16
	target *Target
	logger *log.Logger
}

func NewNeighbor(ip string, port uint16, logger *log.Logger) *Neighbor {
	target := NewTarget(ip, port)
	return &Neighbor{ip, port, target, logger}
}

func (neighbor *Neighbor) Ip() string {
	return neighbor.ip
}

func (neighbor *Neighbor) Port() uint16 {
	return neighbor.port
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
	if err = neighbor.target.Reach(); err != nil {
		err = fmt.Errorf("unable to find node for target %s", neighbor.Target())
		return
	}
	tcp := p2p.NewTCP(neighbor.ip, strconv.Itoa(int(neighbor.port)))
	var c *p2p.Client
	c, err = p2p.NewClient(tcp)
	if err != nil {
		err = fmt.Errorf("failed to start client for target %s: %w", neighbor.Target(), err)
	} else {
		c.SetLogger(log.NewLogger(log.Fatal))
		res, err = c.Send("dialog", req)
	}
	return
}
