package chain

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var cachedBlockchain = make(map[string]*Blockchain)

const (
	GetBlocksRequest          = "GET BLOCKS REQUEST"
	PostTargetRequest         = "POST IP REQUEST"
	GetTransactionsRequest    = "GET TRANSACTIONS REQUEST"
	DeleteTransactionsRequest = "DELETE TRANSACTIONS REQUEST"
	MineRequest               = "MINE REQUEST"
	StartMiningRequest        = "START MINING REQUEST"
	StopMiningRequest         = "STOP MINING REQUEST"
	ConsensusRequest          = "CONSENSUS REQUEST"
)

type Host struct {
	ip   string
	port uint16
}

func NewHost(port uint16) *Host {
	host := new(Host)
	ip, err := getPublicIp()
	if err != nil {
		panic(err)
	}
	host.ip = ip
	host.port = port
	return host
}

func getPublicIp() (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		ip = string(body)
	}
	bodyCloseError := resp.Body.Close()
	if err != nil {
		log.Printf("ERROR: Failed to close body after getting the public IP, error: %v", bodyCloseError)
	}
	return
}

func (host *Host) GetBlockchain() *Blockchain {
	blockchain, ok := cachedBlockchain["blockchain"]
	if !ok {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(fmt.Sprintf("ERROR: Failed to generate private key, err%v\n", err))
		} else {
			hostWallet := NewWallet(privateKey)
			blockchain = NewBlockchain(hostWallet.Address(), host.ip, host.port)
			//TODO remove fmt
			fmt.Println("host address: " + hostWallet.Address())
			cachedBlockchain["blockchain"] = blockchain
		}
	}
	return blockchain
}

func (host *Host) GetBlocks() (res p2p.Data, err error) {
	res = p2p.Data{}
	var blockResponses []*BlockResponse
	for _, block := range host.GetBlockchain().Blocks() {
		blockResponses = append(blockResponses, block.GetDto())
	}
	err = res.SetGob(blockResponses)
	return
}

func (host *Host) PostTarget(ip string, port uint16) (res p2p.Data, err error) {
	host.GetBlockchain().AddTarget(ip, port)

	res = p2p.Data{}
	if err = res.SetGob(true); err != nil {
		return
	}
	return
}

// TODO unused
func (host *Host) GetTransactions() (res p2p.Data, err error) {
	res = p2p.Data{}
	var transactionResponses []*TransactionResponse
	for _, transaction := range host.GetBlockchain().Transactions() {
		transactionResponses = append(transactionResponses, transaction.GetDto())
	}
	if err = res.SetGob(transactionResponses); err != nil {
		return
	}
	return
}

func (host *Host) PostTransactions(request *TransactionRequest) (res p2p.Data, err error) {
	if request.IsInvalid() {
		log.Println("ERROR: Field(s) are missing in transaction request")
		err = errors.New("fail")
		return
	}
	publicKey := NewPublicKey(*request.SenderPublicKey)
	signature := DecodeSignature(*request.Signature)
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

func (host *Host) PutTransactions(request *TransactionRequest) (res p2p.Data, err error) {
	if request.IsInvalid() {
		log.Println("ERROR: Field(s) are missing in transaction request")
		err = errors.New("fail")
		return
	}
	publicKey := NewPublicKey(*request.SenderPublicKey)
	signature := DecodeSignature(*request.Signature)
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

func (host *Host) StopMining() (res p2p.Data, err error) {
	blockchain := host.GetBlockchain()
	blockchain.StopMining()
	res = p2p.Data{}
	if err = res.SetGob(true); err != nil {
		return
	}
	return
}

func (host *Host) Amount(request *AmountRequest) (res p2p.Data, err error) {
	if request.IsInvalid() {
		log.Println("ERROR: Field(s) are missing in amount request")
		err = errors.New("fail")
		return
	}
	blockchainAddress := *request.Address
	amount := host.GetBlockchain().CalculateTotalAmount(blockchainAddress)
	amountResponse := &AmountResponse{amount}
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
	tcp := p2p.NewTCP("127.0.0.1", strconv.Itoa(int(host.port)))

	server, err := p2p.NewServer(tcp)
	if err != nil {
		log.Panicln(err)
	}

	settings := p2p.NewServerSettings()
	settings.SetConnTimeout(HostConnectionTimeoutSecond * time.Second)
	server.SetSettings(settings)

	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		var requestString string
		if err = req.GetGob(&requestString); err == nil {
			switch requestString {
			case GetBlocksRequest:
				if res, err = host.GetBlocks(); err != nil {
					log.Println("ERROR: Failed to get blocks")
					return
				}
			case GetTransactionsRequest:
				if res, err = host.GetTransactions(); err != nil {
					log.Println("ERROR: Failed to get transactions")
					return
				}
			case DeleteTransactionsRequest:
				if res, err = host.DeleteTransactions(); err != nil {
					log.Println("ERROR: Failed to delete transactions")
					return
				}
			case MineRequest:
				if res, err = host.Mine(); err != nil {
					log.Println("ERROR: Failed to mine")
					return
				}
			case StartMiningRequest:
				if res, err = host.StartMining(); err != nil {
					log.Println("ERROR: Failed to start mining")
					return
				}
			case StopMiningRequest:
				if res, err = host.StopMining(); err != nil {
					log.Println("ERROR: Failed to stop mining")
					return
				}
			case ConsensusRequest:
				if res, err = host.Consensus(); err != nil {
					log.Println("ERROR: Failed to get consensus")
					return
				}
			}
			return
		}
		var transactionRequest TransactionRequest
		if err = req.GetGob(&transactionRequest); err == nil {
			if *transactionRequest.Verb == POST {
				if res, err = host.PostTransactions(&transactionRequest); err != nil {
					log.Println("ERROR: Failed to post transactions")
					return
				}
			} else if *transactionRequest.Verb == PUT {
				if res, err = host.PutTransactions(&transactionRequest); err != nil {
					log.Println("ERROR: Failed to put transactions")
					return
				}
			}
			return
		}
		var amountRequest AmountRequest
		if err = req.GetGob(&amountRequest); err == nil {
			if res, err = host.Amount(&amountRequest); err != nil {
				log.Println("ERROR: Failed to get amount")
				return
			}
			return
		}
		var request Request
		if err = req.GetGob(&request); err == nil {
			switch *request.Kind {
			case PostTargetRequest:
				fields := *request.Fields
				port, parseError := strconv.ParseUint(*fields[1], 10, 16)
				if parseError != nil {
					log.Println("ERROR: Field port is invalid in TargetRequest")
					return
				}
				if res, err = host.PostTarget(*fields[0], uint16(port)); err != nil {
					log.Println("ERROR: Failed to post IP")
					return
				}
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
