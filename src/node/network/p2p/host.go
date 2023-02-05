package p2p

import (
	"context"
	"encoding/json"
	"fmt"

	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
)

const (
	GetBlocks       = "GET BLOCKS"
	GetTransactions = "GET TRANSACTIONS"
	StartValidation = "START VALIDATION"
	StopValidation  = "STOP VALIDATION"
)

type Host struct {
	server                Server
	synchronizer          network.Synchronizer
	blockchain            protocol.Blockchain
	pool                  protocol.TransactionsPool
	synchronizationEngine clock.Engine
	validationEngine      clock.Engine
	verificationEngine    clock.Engine
	watch                 clock.Watch
	logger                log.Logger
}

func NewHost(
	server Server,
	synchronizer network.Synchronizer,
	blockchain protocol.Blockchain,
	pool protocol.TransactionsPool,
	synchronizationEngine clock.Engine,
	validationEngine clock.Engine,
	verificationEngine clock.Engine,
	watch clock.Watch,
	logger log.Logger,
) *Host {
	return &Host{server, synchronizer, blockchain, pool, synchronizationEngine, validationEngine, verificationEngine, watch, logger}
}

func (host *Host) Run() error {
	host.startBlockchain()
	host.server.SetHandle("dialog", host.handle)
	return host.startServer()
}

func (host *Host) startBlockchain() {
	host.logger.Info("updating the blockchain...")
	host.synchronizationEngine.Do()
	host.logger.Info("neighbors are synchronized")
	go host.synchronizationEngine.Start()
	host.verificationEngine.Do()
	host.logger.Info("the blockchain is now up to date")
	host.validationEngine.Do()
	go host.validationEngine.Start()
	go host.verificationEngine.Start()
}

func (host *Host) handle(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var unknownRequest bool
	var requestString string
	var lastBlocksRequest network.LastBlocksRequest
	var transactionRequest network.TransactionRequest
	var amountRequest network.AmountRequest
	var targetsRequest []network.TargetRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &requestString); err == nil {
		switch requestString {
		case GetBlocks:
			res = host.getBlocks()
		case GetTransactions:
			res = host.getTransactions()
		case StartValidation:
			go host.validationEngine.Start()
		case StopValidation:
			go host.validationEngine.Stop()
		default:
			unknownRequest = true
		}
	} else if err = json.Unmarshal(data, &lastBlocksRequest); err == nil {
		res = host.getLastBlocks(&lastBlocksRequest)
	} else if err = json.Unmarshal(data, &transactionRequest); err == nil {
		host.addTransactions(&transactionRequest)
	} else if err = json.Unmarshal(data, &amountRequest); err == nil {
		res = host.amount(&amountRequest)
	} else if err = json.Unmarshal(data, &targetsRequest); err == nil {
		host.postTargets(targetsRequest)
	} else {
		unknownRequest = true
	}

	if unknownRequest {
		host.logger.Error("unknown request")
	}
	return
}

func (host *Host) getLastBlocks(request *network.LastBlocksRequest) (res gp2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in last blocks request")
		return
	}
	blockResponses := host.blockchain.LastBlocks(*request.StartingBlockNonce)
	data, err := json.Marshal(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get last blocks: %w", err).Error())
	}
	res.SetBytes(data)
	return
}

func (host *Host) getBlocks() (res gp2p.Data) {
	blockResponses := host.blockchain.Blocks()
	data, err := json.Marshal(blockResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
	}
	res.SetBytes(data)
	return
}

func (host *Host) addTransactions(request *network.TransactionRequest) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in transaction request")
		return
	}
	go host.pool.AddTransaction(request, host.synchronizer.HostTarget())
}

func (host *Host) amount(request *network.AmountRequest) (res gp2p.Data) {
	if request.IsInvalid() {
		host.logger.Error("field(s) are missing in amount request")
		return
	}
	blockchainAddress := *request.Address
	amount := host.blockchain.Copy().CalculateTotalAmount(host.watch.Now().UnixNano(), blockchainAddress)
	amountResponse := &network.AmountResponse{Amount: amount}
	data, err := json.Marshal(amountResponse)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
	}
	res.SetBytes(data)
	return
}

func (host *Host) getTransactions() (res gp2p.Data) {
	transactionResponses := host.pool.Transactions()
	data, err := json.Marshal(transactionResponses)
	if err != nil {
		host.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
	}
	res.SetBytes(data)
	return
}

func (host *Host) postTargets(requests []network.TargetRequest) {
	for _, request := range requests {
		if request.IsInvalid() {
			host.logger.Error("field(s) are missing in target request")
			return
		}
	}
	go host.synchronizer.AddTargets(requests)
}

func (host *Host) startServer() error {
	host.logger.Info("host node started...")
	return host.server.Serve()
}
