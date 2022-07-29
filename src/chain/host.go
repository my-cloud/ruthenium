package chain

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"io/ioutil"
	"net/http"
	"ruthenium/src/log"
	"strconv"
	"time"
)

var cachedBlockchain = make(map[string]*Blockchain)

const (
	GetBlocksRequest       = "GET BLOCKS REQUEST"
	PostTargetRequest      = "POST IP REQUEST"
	GetTransactionsRequest = "GET TRANSACTIONS REQUEST"
	MineRequest            = "MINE REQUEST"
	StartMiningRequest     = "START MINING REQUEST"
	StopMiningRequest      = "STOP MINING REQUEST"
	ConsensusRequest       = "CONSENSUS REQUEST"
)

type Host struct {
	ip     string
	port   uint16
	logger *log.Logger
}

func NewHost(port uint16, logLevel log.Level) *Host {
	host := new(Host)
	host.logger = log.NewLogger(logLevel)
	ip, err := host.getPublicIp()
	if err != nil {
		panic(err)
	}
	host.ip = ip
	host.port = port
	return host
}

func (host *Host) getPublicIp() (ip string, err error) {
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
		host.logger.Error(fmt.Sprintf("ERROR: Failed to close body after getting the public IP\n%v", bodyCloseError))
	}
	return
}

func (host *Host) GetBlockchain() *Blockchain {
	blockchain, ok := cachedBlockchain["blockchain"]
	if !ok {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			host.logger.Fatal(fmt.Sprintf("ERROR: Failed to generate private key\n%v", err))
		} else {
			hostWallet := NewWallet(privateKey)
			blockchain = NewBlockchain(hostWallet.Address(), host.ip, host.port, host.logger)
			//TODO remove fmt
			fmt.Println("host address: " + hostWallet.Address())
			cachedBlockchain["blockchain"] = blockchain
		}
	}
	return blockchain
}

func (host *Host) GetBlocks() (res p2p.Data) {
	blockchain := host.GetBlockchain()
	blockResponses := blockchain.Blocks()
	err := res.SetGob(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Sprintf("ERROR: Failed to get blocks\n%v", err))
	}
	return
}

func (host *Host) PostTarget(request *TargetRequest) {
	if request.IsInvalid() {
		host.logger.Error("ERROR: Field(s) are missing in target request")
	} else {
		blockchain := host.GetBlockchain()
		blockchain.AddTarget(*request.Ip, *request.Port)
	}
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

func (host *Host) PostTransactions(request *TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("ERROR: Field(s) are missing in transaction request")
	} else {
		publicKey := NewPublicKey(*request.SenderPublicKey)
		signature := DecodeSignature(*request.Signature)
		blockchain := host.GetBlockchain()
		blockchain.CreateTransaction(*request.SenderAddress, *request.RecipientAddress, publicKey, *request.Value, signature)
	}
}

func (host *Host) PutTransactions(request *TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("ERROR: Field(s) are missing in transaction request")
	} else {
		publicKey := NewPublicKey(*request.SenderPublicKey)
		signature := DecodeSignature(*request.Signature)
		blockchain := host.GetBlockchain()
		blockchain.UpdateTransaction(*request.SenderAddress, *request.RecipientAddress, publicKey, *request.Value, signature)
	}
}

func (host *Host) Mine() {
	blockchain := host.GetBlockchain()
	blockchain.Mine()
}

func (host *Host) StartMining() {
	blockchain := host.GetBlockchain()
	blockchain.StartMining()
}

func (host *Host) StopMining() {
	blockchain := host.GetBlockchain()
	blockchain.StopMining()
}

func (host *Host) Amount(request *AmountRequest) (res p2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("ERROR: Field(s) are missing in amount request")
	} else {
		blockchainAddress := *request.Address
		amount := host.GetBlockchain().CalculateTotalAmount(blockchainAddress)
		amountResponse := &AmountResponse{amount}
		if err := res.SetGob(amountResponse); err != nil {
			host.logger.Error(fmt.Sprintf("ERROR: Failed to get amount\n%v", err))
		}
	}
	return
}

func (host *Host) Consensus() {
	blockchain := host.GetBlockchain()
	blockchain.ResolveConflicts()
}

func (host *Host) Run() {
	host.GetBlockchain().Run()
	host.startHost()
}

func (host *Host) startHost() {
	tcp := p2p.NewTCP("0.0.0.0", strconv.Itoa(int(host.port)))

	server, err := p2p.NewServer(tcp)
	if err != nil {
		host.logger.Fatal(err.Error())
	} else {
		server.SetLogger(host.logger)
		settings := p2p.NewServerSettings()
		settings.SetConnTimeout(HostConnectionTimeoutSecond * time.Second)
		server.SetSettings(settings)

		server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
			var unknownRequest bool
			var requestString string
			var transactionRequest TransactionRequest
			var amountRequest AmountRequest
			var targetRequest TargetRequest
			res = p2p.Data{}
			if err = req.GetGob(&requestString); err == nil {
				switch requestString {
				case GetBlocksRequest:
					res = host.GetBlocks()
				case GetTransactionsRequest:
					if res, err = host.GetTransactions(); err != nil {
						host.logger.Error(fmt.Sprintf("ERROR: Failed to get transactions, error: %v", err))
					}
				case MineRequest:
					host.Mine()
				case StartMiningRequest:
					host.StartMining()
				case StopMiningRequest:
					host.StopMining()
				case ConsensusRequest:
					host.Consensus()
				default:
					unknownRequest = true
				}
			} else if err = req.GetGob(&transactionRequest); err == nil {
				if *transactionRequest.Verb == POST {
					host.PostTransactions(&transactionRequest)
				} else if *transactionRequest.Verb == PUT {
					host.PutTransactions(&transactionRequest)
				}
			} else if err = req.GetGob(&amountRequest); err == nil {
				res = host.Amount(&amountRequest)
			} else if err = req.GetGob(&targetRequest); err == nil {
				switch *targetRequest.Kind {
				case PostTargetRequest:
					host.PostTarget(&targetRequest)
				default:
					unknownRequest = true
				}
			} else {
				unknownRequest = true
			}

			if unknownRequest {
				host.logger.Error("ERROR: Unknown request")
			}
			return
		})

		err = server.Serve()
		if err != nil {
			host.logger.Fatal(err.Error())
		}
	}
}
