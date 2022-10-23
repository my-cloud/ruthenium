package neighborhood

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/node"
	"github.com/my-cloud/ruthenium/src/log"
	"net"
	"strconv"
	"time"
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
	logger *log.Logger
}

func NewNeighbor(ip string, port uint16, logger *log.Logger) *Neighbor {
	neighbor := new(Neighbor)
	neighbor.ip = ip
	neighbor.port = port
	neighbor.logger = logger
	return neighbor
}

func (neighbor *Neighbor) Ip() string {
	return neighbor.ip
}

func (neighbor *Neighbor) Port() uint16 {
	return neighbor.port
}

func (neighbor *Neighbor) Target() string {
	return net.JoinHostPort(neighbor.ip, strconv.Itoa(int(neighbor.port)))
}

func (neighbor *Neighbor) IsFound() bool {
	target := fmt.Sprintf("%s:%d", neighbor.ip, neighbor.port)
	_, err := net.DialTimeout("tcp", target, NeighborFindingTimeoutSecond*time.Second)
	return err == nil
}

func (neighbor *Neighbor) GetBlocks() (blockResponses []*node.BlockResponse, err error) {
	res, err := neighbor.sendRequest(GetBlocksRequest)
	if err == nil {
		err = res.GetGob(&blockResponses)
	}
	return
}

func (neighbor *Neighbor) SendTargets(request []node.TargetRequest) (err error) {
	_, err = neighbor.sendRequest(request)
	return
}

func (neighbor *Neighbor) AddTransaction(request node.TransactionRequest) (err error) {
	_, err = neighbor.sendRequest(request)
	return
}

func (neighbor *Neighbor) GetTransactions() (transactionResponses []node.TransactionResponse, err error) {
	res, err := neighbor.sendRequest(GetTransactionsRequest)
	if err != nil {
		return
	}
	err = res.GetGob(&transactionResponses)
	if transactionResponses == nil {
		return []node.TransactionResponse{}, err
	}
	return
}

func (neighbor *Neighbor) GetAmount(request node.AmountRequest) (amountResponse *node.AmountResponse, err error) {
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
	if err == nil {
		if neighbor.IsFound() {
			tcp := p2p.NewTCP(neighbor.ip, strconv.Itoa(int(neighbor.port)))
			var c *p2p.Client
			c, err = p2p.NewClient(tcp)
			if err == nil {
				c.SetLogger(log.NewLogger(log.Fatal))
				res, err = c.Send("dialog", req)
			} else {
				err = fmt.Errorf("failed to start client for target %s: %w", neighbor.Target(), err)
			}
		} else {
			err = fmt.Errorf("unable to find node for target %s", neighbor.Target())
		}
	}
	return
}
