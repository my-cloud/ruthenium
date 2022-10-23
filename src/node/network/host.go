package network

import (
	"context"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/api/node/network"
	"github.com/my-cloud/ruthenium/src/api/node/protocol"
	"github.com/my-cloud/ruthenium/src/log"
	"strconv"
	"time"
)

const (
	// TODO rename and extract ParticlesCount to a config file
	ParticlesCount             = 100000000
	connectionTimeoutInSeconds = 10
)

type Host struct {
	ip           string
	port         uint16
	blockchain   protocol.Verifiable
	pool         protocol.Validatable
	validation   protocol.Controllable
	neighborhood *Neighborhood
	logger       *log.Logger
}

func NewHost(
	ip string,
	port uint16,
	blockchain protocol.Verifiable,
	pool protocol.Validatable,
	validation protocol.Controllable,
	neighborhood *Neighborhood,
	logger *log.Logger,
) *Host {
	return &Host{ip, port, blockchain, pool, validation, neighborhood, logger}
}

func (host *Host) GetBlocks() (res p2p.Data) {
	blockResponses := host.blockchain.Blocks()
	err := res.SetGob(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
	}
	return
}

func (host *Host) PostTargets(request []network.TargetRequest) {
	host.neighborhood.AddTargets(request)
}

func (host *Host) GetTransactions() (res p2p.Data) {
	transactionResponses := host.pool.Transactions()
	if err := res.SetGob(transactionResponses); err != nil {
		host.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
	}
	return
}

func (host *Host) AddTransactions(request *network.TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in transaction request")
		return
	}
	neighbors := host.neighborhood.Neighbors()
	host.pool.AddTransaction(request, host.blockchain, neighbors)
}

func (host *Host) Amount(request *network.AmountRequest) (res p2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in amount request")
		return
	}
	blockchainAddress := *request.Address
	amount := host.blockchain.CalculateTotalAmount(time.Now().UnixNano(), blockchainAddress)
	amountResponse := &network.AmountResponse{Amount: amount}
	if err := res.SetGob(amountResponse); err != nil {
		host.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
	}
	return
}

func (host *Host) Run() {
	go func() {
		host.logger.Info("updating the blockchain...")
		host.neighborhood.StartSynchronization()
		host.neighborhood.Wait()
		neighbors := host.neighborhood.Neighbors()
		host.blockchain.Verify(neighbors)
		host.logger.Info("the blockchain is now up to date")
		host.validation.Start()
		host.blockchain.StartVerification(host.pool, host.neighborhood)
	}()
	host.startServer()
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
		var transactionRequest network.TransactionRequest
		var amountRequest network.AmountRequest
		var targetsRequest []network.TargetRequest
		res = p2p.Data{}
		if err = req.GetGob(&requestString); err == nil {
			switch requestString {
			case GetBlocksRequest:
				res = host.GetBlocks()
			case GetTransactionsRequest:
				res = host.GetTransactions()
			case MineRequest:
				host.validation.Do()
			case StartMiningRequest:
				host.validation.Start()
			case StopMiningRequest:
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
