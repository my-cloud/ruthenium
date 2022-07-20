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

func (node *Node) DeleteTransactions() bool {
	res, err := node.sendRequest(DeleteTransactionsRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	var deleted bool
	err = res.GetGob(&deleted)
	return deleted
}

func (node *Node) Consensus() bool {
	res, err := node.sendRequest(ConsensusRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	var consented bool
	err = res.GetGob(&consented)
	return consented
}

func (node *Node) PostTransactions(request PostTransactionRequest) bool {
	res, err := node.sendRequest(request)
	if err != nil {
		log.Println(err)
		return false
	}

	var created bool
	err = res.GetGob(&created)
	return created
}

func (node *Node) PutTransactions(request PutTransactionRequest) bool {
	res, err := node.sendRequest(request)
	if err != nil {
		log.Println(err)
		return false
	}

	var updated bool
	err = res.GetGob(&updated)
	return updated
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
