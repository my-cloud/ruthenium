package node

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"log"
	"net"
	"os"
	"ruthenium/src/chain"
	"strconv"
)

var cachedBlockchain = make(map[string]*Blockchain)

const (
	GetBlocksRequest          = "GET BLOCKS REQUEST"
	GetTransactionsRequest    = "GET TRANSACTIONS REQUEST"
	DeleteTransactionsRequest = "DELETE TRANSACTIONS REQUEST"
	MineRequest               = "MINE REQUEST"
	StartMiningRequest        = "START MINING REQUEST"
	ConsensusRequest          = "CONSENSUS REQUEST"
)

type Host struct {
	port uint16
	ip   string
}

func NewHost(port uint16) *Host {
	host := new(Host)
	host.port = port
	return host
}

func (host *Host) Port() uint16 {
	return host.port
}

func (host *Host) GetBlockchain() *Blockchain {
	blockchain, ok := cachedBlockchain["blockchain"]
	if !ok {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(fmt.Sprintf("ERROR: Failed to generate private key, err%v\n", err))
		} else {
			hostWallet := chain.NewWallet(privateKey)
			blockchain = NewBlockchain(hostWallet.Address(), host.Port())
			//TODO remove fmt
			fmt.Println("host address: " + hostWallet.Address())
			cachedBlockchain["blockchain"] = blockchain
		}
	}
	return blockchain
}

func (host *Host) GetChain() (res p2p.Data, err error) {
	res = p2p.Data{}
	err = res.SetGob(host.GetBlockchain().Blocks())
	return res, err
}

func (host *Host) GetTransactions() (res p2p.Data, err error) {
	res = p2p.Data{}
	if err = res.SetGob(host.GetBlockchain().Transactions()); err != nil {
		return
	}
	return
}

func (host *Host) PostTransactions(request *chain.PostTransactionRequest) (res p2p.Data, err error) {
	if request.IsInvalid() {
		log.Println("ERROR: Field(s) are missing in transaction request")
		err = errors.New("fail")
		return
	}
	publicKey := chain.NewPublicKey(*request.SenderPublicKey)
	signature := chain.DecodeSignature(*request.Signature)
	blockchain := host.GetBlockchain()
	isCreated := blockchain.CreateTransaction(*request.SenderAddress, *request.RecipientAddress, publicKey, *request.Value, signature)

	if !isCreated {
		err = errors.New("fail")
		return
	}
	res = p2p.Data{}
	if err = res.SetGob(true); err != nil {
		return
	}
	return
}

func (host *Host) PutTransactions(request *chain.PutTransactionRequest) (res p2p.Data, err error) {
	if request.IsInvalid() {
		log.Println("ERROR: Field(s) are missing in transaction request")
		err = errors.New("fail")
		return
	}
	publicKey := chain.NewPublicKey(*request.SenderPublicKey)
	signature := chain.DecodeSignature(*request.Signature)
	blockchain := host.GetBlockchain()
	isCreated := blockchain.UpdateTransaction(*request.SenderAddress, *request.RecipientAddress, publicKey, *request.Value, signature)

	if !isCreated {
		err = errors.New("fail")
		return
	}
	res = p2p.Data{}
	if err = res.SetGob(true); err != nil {
		return
	}
	return
}

func (host *Host) DeleteTransactions() (res p2p.Data, err error) {
	blockchain := host.GetBlockchain()
	blockchain.ClearTransactions()
	res = p2p.Data{}
	if err = res.SetGob(true); err != nil {
		return
	}
	return
}

func (host *Host) Mine() (res p2p.Data, err error) {
	blockchain := host.GetBlockchain()
	isMined := blockchain.Mine()
	if isMined {
		res = p2p.Data{}
		if err = res.SetGob(true); err != nil {
			return
		}
	} else {
		err = errors.New("fail")
		return
	}
	return
}

func (host *Host) StartMining() (res p2p.Data, err error) {
	blockchain := host.GetBlockchain()
	blockchain.StartMining()
	res = p2p.Data{}
	if err = res.SetGob(true); err != nil {
		return
	}
	return
}

func (host *Host) Amount(request *chain.AmountRequest) (res p2p.Data, err error) {
	if request.IsInvalid() {
		log.Println("ERROR: Field(s) are missing in amount request")
		err = errors.New("fail")
		return
	}
	blockchainAddress := *request.Address
	amount := host.GetBlockchain().CalculateTotalAmount(blockchainAddress)
	amountResponse := chain.NewAmountResponse(amount)
	res = p2p.Data{}
	if err = res.SetGob(amountResponse); err != nil {
		return
	}
	return
}

func (host *Host) Consensus() (res p2p.Data, err error) {
	blockchain := host.GetBlockchain()
	isSelfBlockchainReplacedByNeighborsOne := blockchain.ResolveConflicts()
	res = p2p.Data{}
	if err = res.SetGob(isSelfBlockchainReplacedByNeighborsOne); err != nil {
		return
	}
	return
}

func (host *Host) Run() {
	host.GetBlockchain().Run()
	host.startHost()
}

func (host *Host) startHost() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "127.0.0.1"
	}
	ips, err := net.LookupHost(hostname)
	if err != nil {
		ips[0] = "127.0.0.1"
	}
	// FIXME host[3] = 192.168.1.90
	host.ip = "localhost"

	tcp := p2p.NewTCP(host.ip, strconv.Itoa(int(host.port)))

	server, err := p2p.NewServer(tcp)
	if err != nil {
		log.Panicln(err)
	}

	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		var requestString string
		if err = req.GetGob(&requestString); err == nil {
			switch requestString {
			case GetBlocksRequest:
				if res, err = host.GetChain(); err != nil {
					return
				}
			case GetTransactionsRequest:
				if res, err = host.GetTransactions(); err != nil {
					return
				}
			case DeleteTransactionsRequest:
				if res, err = host.DeleteTransactions(); err != nil {
					return
				}
			case MineRequest:
				if res, err = host.Mine(); err != nil {
					return
				}
			case StartMiningRequest:
				if res, err = host.StartMining(); err != nil {
					return
				}
			case ConsensusRequest:
				if res, err = host.Consensus(); err != nil {
					return
				}
			}
			return
		}
		var requestPostTransaction chain.PostTransactionRequest
		if err = req.GetGob(&requestPostTransaction); err == nil {
			if res, err = host.PostTransactions(&requestPostTransaction); err != nil {
				return
			}
			return
		}
		var requestPutTransaction chain.PutTransactionRequest
		if err = req.GetGob(&requestPutTransaction); err == nil {
			if res, err = host.PutTransactions(&requestPutTransaction); err != nil {
				return
			}
			return
		}
		var requestAmount chain.AmountRequest
		if err = req.GetGob(&requestAmount); err == nil {
			if res, err = host.Amount(&requestAmount); err != nil {
				return
			}
			return
		}

		log.Println("ERROR: Unknown request")
		return
	})

	err = server.Serve()
	if err != nil {
		log.Panicln(err)
	}
}
