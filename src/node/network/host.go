package network

import (
	"context"
	"fmt"
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
)

const (
	// TODO rename and extract ParticlesCount to a config file
	ParticlesCount = 100000000
)

type Host struct {
	server       Server
	blockchain   Blockchain
	pool         TransactionsPool
	validator    Validator
	synchronizer *Synchronizer
	watch        clock.Time
	logger       *log.Logger
}

func NewHost(
	server Server,
	blockchain Blockchain,
	pool TransactionsPool,
	validator Validator,
	synchronizer *Synchronizer,
	watch clock.Time,
	logger *log.Logger,
) *Host {
	return &Host{server, blockchain, pool, validator, synchronizer, watch, logger}
}

func (host *Host) GetBlocks() (res p2p.Data) {
	blockResponses := host.blockchain.Blocks()
	err := res.SetGob(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
	}
	return
}

func (host *Host) PostTargets(request []neighborhood.TargetRequest) {
	host.synchronizer.AddTargets(request)
}

func (host *Host) GetTransactions() (res p2p.Data) {
	transactionResponses := host.pool.Transactions()
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
	neighbors := host.synchronizer.Neighbors()
	host.pool.AddTransaction(request, host.blockchain, neighbors)
}

func (host *Host) Amount(request *neighborhood.AmountRequest) (res p2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in amount request")
		return
	}
	blockchainAddress := *request.Address
	amount := host.blockchain.CalculateTotalAmount(host.watch.Now().UnixNano(), blockchainAddress)
	amountResponse := &neighborhood.AmountResponse{Amount: amount}
	if err := res.SetGob(amountResponse); err != nil {
		host.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
	}
	return
}

func (host *Host) Run() error {
	go host.startBlockchain()
	host.server.SetHandle("dialog", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) { return host.handle(req) })
	return host.startServer()
}

func (host *Host) startBlockchain() {
	host.logger.Info("updating the blockchain...")
	host.synchronizer.StartSynchronization()
	host.synchronizer.Wait()
	host.blockchain.Verify()
	host.logger.Info("the blockchain is now up to date")
	host.validator.StartValidation()
	host.blockchain.StartVerification()
}

func (host *Host) handle(req p2p.Data) (res p2p.Data, err error) {
	var unknownRequest bool
	var requestString string
	var transactionRequest neighborhood.TransactionRequest
	var amountRequest neighborhood.AmountRequest
	var targetsRequest []neighborhood.TargetRequest
	res = p2p.Data{}
	if err = req.GetGob(&requestString); err == nil {
		switch requestString {
		case GetBlocksRequest:
			res = host.GetBlocks()
		case GetTransactionsRequest:
			res = host.GetTransactions()
		case MineRequest:
			host.validator.Validate()
		case StartMiningRequest:
			host.validator.StartValidation()
		case StopMiningRequest:
			host.validator.StopValidation()
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
}

func (host *Host) startServer() error {
	host.logger.Info("host server is running...")
	return host.server.Serve()
}