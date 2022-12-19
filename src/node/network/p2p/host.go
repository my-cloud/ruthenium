package p2p

import (
	"context"
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
)

const (
	GetBlocksRequest       = "GET BLOCKS REQUEST"
	GetTransactionsRequest = "GET TRANSACTIONS REQUEST"
	MineRequest            = "MINE REQUEST"
	StartMiningRequest     = "START MINING REQUEST"
	StopMiningRequest      = "STOP MINING REQUEST"
)

type Host struct {
	server                Server
	synchronizer          network.Synchronizer
	blockchain            protocol.Blockchain
	pool                  protocol.TransactionsPool
	synchronizationEngine clock.Engine
	validationEngine      clock.Engine
	verificationEngine    clock.Engine
	time                  clock.Time
	logger                *log.Logger
}

func NewHost(
	server Server,
	synchronizer network.Synchronizer,
	blockchain protocol.Blockchain,
	pool protocol.TransactionsPool,
	synchronizationEngine clock.Engine,
	validationEngine clock.Engine,
	verificationEngine clock.Engine,
	time clock.Time,
	logger *log.Logger,
) *Host {
	return &Host{server, synchronizer, blockchain, pool, synchronizationEngine, validationEngine, verificationEngine, time, logger}
}

func (host *Host) GetBlocks() (res gp2p.Data) {
	blockResponses := host.blockchain.Blocks()
	err := res.SetGob(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
	}
	return
}

func (host *Host) PostTargets(request []network.TargetRequest) {
	host.synchronizer.AddTargets(request)
}

func (host *Host) GetTransactions() (res gp2p.Data) {
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
	neighbors := host.synchronizer.Neighbors()
	host.pool.AddTransaction(request, neighbors)
}

func (host *Host) Amount(request *network.AmountRequest) (res gp2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in amount request")
		return
	}
	blockchainAddress := *request.Address
	amount := host.blockchain.Copy().CalculateTotalAmount(host.time.Now().UnixNano(), blockchainAddress)
	amountResponse := &network.AmountResponse{Amount: amount}
	if err := res.SetGob(amountResponse); err != nil {
		host.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
	}
	return
}

func (host *Host) Run() error {
	host.startBlockchain()
	host.server.SetHandle("dialog", host.handle)
	return host.startServer()
}

func (host *Host) startBlockchain() {
	host.synchronizationEngine.Do()
	host.synchronizationEngine.Wait()
	host.synchronizationEngine.Start()
	host.validationEngine.Start()
	host.verificationEngine.Start()
}

func (host *Host) handle(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var unknownRequest bool
	var requestString string
	var transactionRequest network.TransactionRequest
	var amountRequest network.AmountRequest
	var targetsRequest []network.TargetRequest
	res = gp2p.Data{}
	if err = req.GetGob(&requestString); err == nil {
		switch requestString {
		case GetBlocksRequest:
			res = host.GetBlocks()
		case GetTransactionsRequest:
			res = host.GetTransactions()
		case MineRequest:
			host.validationEngine.Do()
		case StartMiningRequest:
			host.validationEngine.Start()
		case StopMiningRequest:
			host.validationEngine.Stop()
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
	host.logger.Info("host node started...")
	return host.server.Serve()
}
