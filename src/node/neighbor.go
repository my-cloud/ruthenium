package node

import (
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"log"
	"net"
	"ruthenium/src/chain"
	"strconv"
	"time"
)

type Neighbor struct {
	ip     string
	port   uint16
	client *p2p.Client
}

func NewNeighbor(ip string, port uint16) *Neighbor {
	neighbor := new(Neighbor)
	neighbor.ip = ip
	neighbor.port = port
	return neighbor
}

func (neighbor *Neighbor) StartClient() {
	tcp := p2p.NewTCP("localhost", strconv.Itoa(int(neighbor.port)))

	client, err := p2p.NewClient(tcp)
	if err != nil {
		log.Println(err)
	}

	neighbor.client = client
}

func (neighbor *Neighbor) IpAndPort() string {
	return fmt.Sprintf("%s:%d", neighbor.ip, neighbor.port)
}

func (neighbor *Neighbor) IsFound() bool {
	target := fmt.Sprintf("%s:%d", neighbor.ip, neighbor.port)
	_, err := net.DialTimeout("tcp", target, time.Second)
	return err == nil
}

func (neighbor *Neighbor) GetBlocks() []*chain.Block {
	res, err := neighbor.sendRequest(GetBlocksRequest)
	if err != nil {
		log.Println(err)
		return nil
	}

	var blocks []*chain.Block
	err = res.GetGob(&blocks)
	if err != nil {
		log.Println(err)
		return nil
	}

	return blocks
}

func (neighbor *Neighbor) DeleteTransactions() bool {
	res, err := neighbor.sendRequest(DeleteTransactionsRequest)
	if err != nil {
		log.Println(err)
		return false
	}

	var deleted bool
	err = res.GetGob(&deleted)
	return deleted
}

func (neighbor *Neighbor) Consensus() bool {
	res, err := neighbor.sendRequest(ConsensusRequest)
	if err != nil {
		log.Println(err)
		return false

	}

	var consented bool
	err = res.GetGob(&consented)
	return consented
}

func (neighbor *Neighbor) PostTransactions(request *chain.PostTransactionRequest) bool {
	res, err := neighbor.sendRequest(request)
	if err != nil {
		log.Println(err)
		return false

	}

	var created bool
	err = res.GetGob(&created)
	return created
}

func (neighbor *Neighbor) PutTransactions(request *chain.PutTransactionRequest) bool {
	res, err := neighbor.sendRequest(request)
	if err != nil {
		log.Println(err)
		return false

	}

	var updated bool
	err = res.GetGob(&updated)
	return updated
}

func (neighbor *Neighbor) GetAmount(request chain.AmountRequest) *chain.AmountResponse {
	res, err := neighbor.sendRequest(request)
	if err != nil {
		log.Println(err)
		return nil
	}

	var amount *chain.AmountResponse
	err = res.GetGob(&amount)
	return amount
}

func (neighbor *Neighbor) sendRequest(request interface{}) (res p2p.Data, err error) {
	req := p2p.Data{}
	err = req.SetGob(request)
	if err != nil {
		log.Println(err)
		return
	}

	res = p2p.Data{}
	res, err = neighbor.client.Send("dialog", req)
	if err != nil {
		log.Println(err)
		return
	}

	return res, err
}
