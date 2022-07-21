package chain

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	ip     string
	port   uint16
	client *p2p.Client
	mutex  sync.Mutex
}

func NewNode(ip string, port uint16) *Node {
	node := new(Node)
	node.ip = ip
	node.port = port
	return node
}

func (node *Node) StartClient() {
	tcp := p2p.NewTCP(node.ip, strconv.Itoa(int(node.port)))
	client, err := p2p.NewClient(tcp)
	if err != nil {
		log.Println(err)
	}

	node.client = client
}

func (node *Node) IpAndPort() string {
	return fmt.Sprintf("%s:%d", node.ip, node.port)
}

func (node *Node) IsFound() bool {
	target := fmt.Sprintf("%s:%d", node.ip, node.port)
	_, err := net.DialTimeout("tcp", target, time.Second)
	return err == nil
}

func (node *Node) GetBlocks() []*Block {
	res, err := node.sendRequest(GetBlocksRequest)
	if err != nil {
		log.Println(err)
		return nil
	}

	var blockResponses []*BlockResponse
	err = res.GetGob(&blockResponses)
	if err != nil {
		log.Println(err)
		return nil
	}

	var blocks []*Block
	for _, block := range blockResponses {
		blocks = append(blocks, NewBlockFromDto(block))
	}

	return blocks
}

func (node *Node) DeleteTransactions() (deleted bool) {
	res, err := node.sendRequest(DeleteTransactionsRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&deleted)
	return
}

func (node *Node) Consensus() (consented bool) {
	res, err := node.sendRequest(ConsensusRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&consented)
	return
}

func (node *Node) PostTransactions(request PostTransactionRequest) (created bool) {
	res, err := node.sendRequest(request)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&created)
	return
}

func (node *Node) PutTransactions(request PutTransactionRequest) (updated bool) {
	res, err := node.sendRequest(request)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&updated)
	return
}

func (node *Node) GetAmount(request AmountRequest) *AmountResponse {
	res, err := node.sendRequest(request)
	if err != nil {
		log.Println(err)
		return nil
	}

	var amount *AmountResponse
	err = res.GetGob(&amount)
	return amount
}

func (node *Node) Mine() (mined bool) {
	res, err := node.sendRequest(MineRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&mined)
	return
}

func (node *Node) StartMining() (miningStarted bool) {
	res, err := node.sendRequest(StartMiningRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&miningStarted)
	return
}

func (node *Node) StopMining() (miningStopped bool) {
	res, err := node.sendRequest(StopMiningRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	err = res.GetGob(&miningStopped)
	return
}

func (node *Node) sendRequest(request interface{}) (res p2p.Data, err error) {
	req := p2p.Data{}
	err = req.SetGob(request)
	if err != nil {
		log.Println(err)
		return
	}

	res = p2p.Data{}
	node.mutex.Lock()
	defer node.mutex.Unlock()
	res, err = node.client.Send("dialog", req)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
