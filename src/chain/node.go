package chain

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"net"
	"ruthenium/src/log"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	ip     string
	port   uint16
	mutex  sync.Mutex
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
	_, err := net.DialTimeout("tcp", target, NeighborClientFindingTimeoutSecond*time.Second)
	return err == nil
}

func (node *Node) GetBlocks() []*Block {
	res, err := node.sendRequest(GetBlocksRequest)
	if err != nil {
		node.logger.Error(err.Error())
		return nil
	}

	var blockResponses []*BlockResponse
	err = res.GetGob(&blockResponses)
	if err != nil {
		node.logger.Error(err.Error())
		return nil
	}

	var blocks []*Block
	for _, block := range blockResponses {
		blocks = append(blocks, NewBlockFromDto(block))
	}

	return blocks
}

func (node *Node) SendTarget(request TargetRequest) (sent bool) {
	res, err := node.sendRequest(request)
	if err != nil {
		node.logger.Error(err.Error())
		return false
	}

	err = res.GetGob(&sent)
	return
}

func (node *Node) Consensus() (consented bool) {
	res, err := node.sendRequest(ConsensusRequest)
	if err != nil {
		node.logger.Error(err.Error())
		return false
	}

	err = res.GetGob(&consented)
	return
}

func (node *Node) UpdateTransactions(request TransactionRequest) (created bool) {
	res, err := node.sendRequest(request)
	if err != nil {
		node.logger.Error(err.Error())
		return false
	}

	err = res.GetGob(&created)
	return
}

func (node *Node) GetAmount(request AmountRequest) (amountResponse *AmountResponse) {
	res, err := node.sendRequest(request)
	if err != nil {
		node.logger.Error(err.Error())
		return nil
	}

	err = res.GetGob(&amountResponse)
	return
}

func (node *Node) Mine() (mined bool) {
	res, err := node.sendRequest(MineRequest)
	if err != nil {
		node.logger.Error(err.Error())
		return false
	}

	err = res.GetGob(&mined)
	return
}

func (node *Node) StartMining() (miningStarted bool) {
	res, err := node.sendRequest(StartMiningRequest)
	if err != nil {
		node.logger.Error(err.Error())
		return false
	}

	err = res.GetGob(&miningStarted)
	return
}

func (node *Node) StopMining() (miningStopped bool) {
	res, err := node.sendRequest(StopMiningRequest)
	if err != nil {
		node.logger.Error(err.Error())
		return false
	}

	err = res.GetGob(&miningStopped)
	return
}

func (node *Node) sendRequest(request interface{}) (res p2p.Data, err error) {
	req := p2p.Data{}
	err = req.SetGob(request)
	if err == nil {
		tcp := p2p.NewTCP(node.ip, strconv.Itoa(int(node.port)))
		var client *p2p.Client
		client, err = p2p.NewClient(tcp)
		if err == nil {
			client.SetLogger(node.logger)
			client.SetLogger(node.logger)
			settings := p2p.NewClientSettings()
			settings.SetRetry(1, p2p.DefaultDelayTimeout)
			client.SetSettings(settings)

			res = p2p.Data{}
			res, err = client.Send("dialog", req)
		}
	}

	return
}
