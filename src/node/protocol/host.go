package protocol

import (
	"context"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultPort                 = 8106
	ParticlesCount              = 100000000
	connectionTimeoutInSeconds  = 10
	validationIntervalInSeconds = 60
)

type Host struct {
	ip           string
	port         uint16
	network      *Network
	blockchain   *Blockchain
	pool         *Pool
	validation   *Validation
	verification *Verification
	logger       *log.Logger
}

func NewHost(mnemonic string, derivationPath string, password string, privateKey string, port uint16, logLevel log.Level) *Host {
	host := new(Host)
	host.logger = log.NewLogger(logLevel)
	host.port = port
	ip, err := host.findPublicIp()
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to find the public IP: %w", err).Error())
	}
	host.ip = ip
	wallet, err := encryption.DecodeWallet(mnemonic, derivationPath, password, privateKey)
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to create wallet: %w", err).Error())
	} else {
		watch := clock.NewWatch()
		host.network = NewNetwork(host.ip, host.port, watch, host.logger)
		validationTimer := validationIntervalInSeconds * time.Second
		host.blockchain = NewBlockchain(validationTimer.Nanoseconds(), watch, host.logger)
		host.pool = NewPool(watch, host.logger)
		host.validation = NewValidation(wallet.Address(), host.blockchain, host.pool, watch, validationTimer, host.logger)
		host.verification = NewVerification(host.blockchain, host.pool, host.network)
	}
	return host
}

func (host *Host) GetBlocks() (res p2p.Data) {
	blockResponses := host.blockchain.BlockResponses()
	err := res.SetGob(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
	}
	return
}

func (host *Host) PostTargets(request []neighborhood.TargetRequest) {
	host.network.AddTargets(request)
}

func (host *Host) GetTransactions() (res p2p.Data) {
	var transactionResponses []*neighborhood.TransactionResponse
	for _, transaction := range host.pool.Transactions() {
		transactionResponses = append(transactionResponses, transaction.GetResponse())
	}
	if err := res.SetGob(transactionResponses); err != nil {
		host.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
	}
	return
}

func (host *Host) AddTransactions(request *neighborhood.TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in transaction request")
		return
	}
	transaction, err := NewTransactionFromRequest(request)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to instantiate transaction: %w", err).Error())
		return
	}
	neighbors := host.network.Neighbors()
	host.pool.AddTransaction(transaction, host.blockchain, neighbors)
}

func (host *Host) Amount(request *neighborhood.AmountRequest) (res p2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in amount request")
		return
	}
	blockchainAddress := *request.Address
	amount := host.blockchain.CalculateTotalAmount(time.Now().UnixNano(), blockchainAddress)
	amountResponse := &neighborhood.AmountResponse{amount}
	if err := res.SetGob(amountResponse); err != nil {
		host.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
	}
	return
}

func (host *Host) Run() {
	go func() {
		host.logger.Info("updating the blockchain...")
		host.network.StartNeighborsSynchronization()
		host.network.Wait()
		neighbors := host.network.Neighbors()
		host.blockchain.Verify(neighbors)
		host.network.logger.Info("the blockchain is now up to date")
		host.validation.Start()
		host.verification.Start()
	}()
	host.startServer()
}

func (host *Host) findPublicIp() (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			host.logger.Error(fmt.Errorf("failed to close public IP request body: %w", bodyCloseError).Error())
		}
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ip = string(body)
	return
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
	settings.SetConnTimeout(connectionTimeoutInSeconds * time.Second)
	server.SetSettings(settings)
	server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		var unknownRequest bool
		var requestString string
		var transactionRequest neighborhood.TransactionRequest
		var amountRequest neighborhood.AmountRequest
		var targetsRequest []neighborhood.TargetRequest
		res = p2p.Data{}
		if err = req.GetGob(&requestString); err == nil {
			switch requestString {
			case neighborhood.GetBlocksRequest:
				res = host.GetBlocks()
			case neighborhood.GetTransactionsRequest:
				res = host.GetTransactions()
			case neighborhood.MineRequest:
				host.validation.Do()
			case neighborhood.StartMiningRequest:
				host.validation.Start()
			case neighborhood.StopMiningRequest:
				host.validation.Stop()
			default:
				unknownRequest = true
			}
		} else if err = req.GetGob(&transactionRequest); err == nil {
			host.AddTransactions(&transactionRequest)
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
	host.logger.Info("host server is running...")
	err = server.Serve()
	if err != nil {
		host.logger.Fatal(fmt.Errorf("failed to start server: %w", err).Error())
	}
}
