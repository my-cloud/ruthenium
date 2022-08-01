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

func (node *Node) GetBlocks() (blockResponses []*BlockResponse, err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	res, err := node.sendRequest(GetBlocksRequest)
	if err == nil {
		err = res.GetGob(&blockResponses)
	}

	return
}

func (node *Node) SendTargets(request []TargetRequest) (err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	_, err = node.sendRequest(request)
	return
}

func (node *Node) Consensus() (err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	_, err = node.sendRequest(ConsensusRequest)
	return
}

func (node *Node) UpdateTransactions(request TransactionRequest) (err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	_, err = node.sendRequest(request)
	return
}

func (node *Node) GetAmount(request AmountRequest) (amountResponse *AmountResponse, err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	res, err := node.sendRequest(request)
	if err == nil {
		err = res.GetGob(&amountResponse)
	}

	return
}

func (node *Node) Mine() (err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	_, err = node.sendRequest(MineRequest)
	return
}

func (node *Node) StartMining() (err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	_, err = node.sendRequest(StartMiningRequest)
	return
}

func (node *Node) StopMining() (err error) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	_, err = node.sendRequest(StopMiningRequest)
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
			//settings := p2p.NewClientSettings()
			//settings.SetRetry(10, p2p.DefaultDelayTimeout)
			//client.SetSettings(settings)
			res, err = client.Send("dialog", req)
		}
	}

	return
}
