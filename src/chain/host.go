package chain

import (
	"context"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"io/ioutil"
	"net/http"
	"ruthenium/src/log"
	"strconv"
	"time"
)

const (
	GetBlocksRequest       = "GET BLOCKS REQUEST"
	GetTransactionsRequest = "GET TRANSACTIONS REQUEST"
	MineRequest            = "MINE REQUEST"
	StartMiningRequest     = "START MINING REQUEST"
	StopMiningRequest      = "STOP MINING REQUEST"
	ConsensusRequest       = "CONSENSUS REQUEST"
)

type Host struct {
	publicKey  string
	privateKey string
	ip         string
	port       uint16
	blockchain *Blockchain
	logger     *log.Logger
}

func NewHost(publicKey string, privateKey string, port uint16, logLevel log.Level) *Host {
	host := new(Host)
	host.publicKey = publicKey
	host.privateKey = privateKey
	host.logger = log.NewLogger(logLevel)
	host.port = port
	ip, err := host.findPublicIp()
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to find the public IP: %w", err).Error())
	}
	host.ip = ip
	wallet, err := NewWallet(host.publicKey, host.privateKey)
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to create wallet: %w", err).Error())
	} else {
		host.blockchain = NewBlockchain(wallet.Address(), host.ip, host.port, host.logger)
	}
	return host
}

func (host *Host) GetBlocks() (res p2p.Data) {
	blockResponses := host.blockchain.Blocks()
	err := res.SetGob(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
	}
	return
}

func (host *Host) PostTargets(request []TargetRequest) {
	host.blockchain.AddTargets(request)
}

func (host *Host) GetTransactions() (res p2p.Data) {
	var transactionResponses []*TransactionResponse
	for _, transaction := range host.blockchain.Transactions() {
		transactionResponses = append(transactionResponses, transaction.GetDto())
	}
	if err := res.SetGob(transactionResponses); err != nil {
		host.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
	}
	return
}

func (host *Host) PostTransactions(request *TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in transaction request")
		return
	}
	publicKey, err := NewPublicKey(*request.SenderPublicKey)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to decode transaction public key: %w", err).Error())
		return
	}
	signature, err := DecodeSignature(*request.Signature)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to decode transaction signature: %w", err).Error())
		return
	}
	transaction := NewTransaction(publicKey, *request.SenderAddress, *request.RecipientAddress, *request.Value, host.logger)
	host.blockchain.CreateTransaction(transaction, signature)
}

func (host *Host) PutTransactions(request *TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in transaction request")
		return
	}
	publicKey, err := NewPublicKey(*request.SenderPublicKey)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to decode transaction public key: %w", err).Error())
		return
	}
	signature, err := DecodeSignature(*request.Signature)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to decode transaction signature: %w", err).Error())
		return
	}
	transaction := NewTransaction(publicKey, *request.SenderAddress, *request.RecipientAddress, *request.Value, host.logger)
	host.blockchain.AddTransaction(transaction, signature)
}

func (host *Host) Mine() {
	host.blockchain.Mine()
}

func (host *Host) StartMining() {
	host.blockchain.StartMining()
}

func (host *Host) StopMining() {
	host.blockchain.StopMining()
}

func (host *Host) Amount(request *AmountRequest) (res p2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in amount request")
		return
	}
	blockchainAddress := *request.Address
	amount := host.blockchain.CalculateTotalAmount(blockchainAddress)
	amountResponse := &AmountResponse{amount}
	if err := res.SetGob(amountResponse); err != nil {
		host.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
	}
	return
}

func (host *Host) Consensus() {
	host.blockchain.ResolveConflicts()
}

func (host *Host) Run() {
	host.blockchain.Run()
	host.startServer()
}

func (host *Host) findPublicIp() (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(body), resp.Body.Close()
}

func (host *Host) startServer() {
	tcp := p2p.NewTCP("0.0.0.0", strconv.Itoa(int(host.port)))
	server, err := p2p.NewServer(tcp)
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to create server: %w", err).Error())
		return
	}
	server.SetLogger(log.NewLogger(log.Fatal))
	settings := p2p.NewServerSettings()
	settings.SetConnTimeout(HostConnectionTimeoutSecond * time.Second)
	server.SetSettings(settings)
	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		var unknownRequest bool
		var requestString string
		var transactionRequest TransactionRequest
		var amountRequest AmountRequest
		var targetsRequest []TargetRequest
		res = p2p.Data{}
		if err = req.GetGob(&requestString); err == nil {
			switch requestString {
			case GetBlocksRequest:
				res = host.GetBlocks()
			case GetTransactionsRequest:
				res = host.GetTransactions()
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
		} else if err = req.GetGob(&targetsRequest); err == nil {
			host.PostTargets(targetsRequest)
		} else {
			unknownRequest = true
		}

		if unknownRequest {
			host.logger.Error("unknown request")
		}
		return
	})
	err = server.Serve()
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to start server: %w", err).Error())
	}
}
