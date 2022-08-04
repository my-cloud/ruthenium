package chain

import (
	"errors"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"net"
	"ruthenium/src/log"
	"strconv"
	"time"
)

type Node struct {
	ip     string
	port   uint16
	logger *log.Logger
}

func NewNode(ip string, port uint16, logger *log.Logger) *Node {
	node := new(Node)
	node.ip = ip
	node.port = port
	node.logger = logger
	return node
}

func (node *Node) Ip() string {
	return node.ip
}

func (node *Node) Port() uint16 {
	return node.port
}

func (node *Node) Target() string {
	return net.JoinHostPort(node.ip, strconv.Itoa(int(node.port)))
}

func (node *Node) IsFound() bool {
	target := fmt.Sprintf("%s:%d", node.ip, node.port)
	_, err := net.DialTimeout("tcp", target, NeighborFindingTimeoutSecond*time.Second)
	return err == nil
}

func (node *Node) GetBlocks() (blockResponses []*BlockResponse, err error) {
	res, err := node.sendRequest(GetBlocksRequest)
	if err == nil {
		err = res.GetGob(&blockResponses)
	}

	return
}

func (node *Node) SendTargets(request []TargetRequest) (err error) {
	_, err = node.sendRequest(request)
	return
}

func (node *Node) Consensus() (err error) {
	_, err = node.sendRequest(ConsensusRequest)
	return
}

func (node *Node) UpdateTransactions(request TransactionRequest) (err error) {
	_, err = node.sendRequest(request)
	return
}

func (node *Node) GetAmount(request AmountRequest) (amountResponse *AmountResponse, err error) {
	res, err := node.sendRequest(request)
	if err == nil {
		err = res.GetGob(&amountResponse)
	}

	return
}

func (node *Node) Mine() (err error) {
	_, err = node.sendRequest(MineRequest)
	return
}

func (node *Node) StartMining() (err error) {
	_, err = node.sendRequest(StartMiningRequest)
	return
}

func (node *Node) StopMining() (err error) {
	_, err = node.sendRequest(StopMiningRequest)
	return
}

func (node *Node) sendRequest(request interface{}) (res p2p.Data, err error) {
	req := p2p.Data{}
	err = req.SetGob(request)
	if err == nil {
		if node.IsFound() {
			tcp := p2p.NewTCP(node.ip, strconv.Itoa(int(node.port)))
			var client *p2p.Client
			client, err = p2p.NewClient(tcp)
			if err == nil {
				client.SetLogger(log.NewLogger(log.Fatal))
			}
			if err != nil {
				node.logger.Error(fmt.Sprintf("Failed to start client for target %s\n%v", node.Target(), err))
			} else {
				res, err = client.Send("dialog", req)
			}
		} else {
			err = errors.New("unable to find node")
		}
	}

	return
}
