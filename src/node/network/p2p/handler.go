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
)

type Handler struct {
	blockchain       protocol.Blockchain
	synchronizer     network.Synchronizer
	transactionsPool protocol.TransactionsPool
	watch            clock.Watch
	logger           log.Logger
}

func NewHandler(blockchain protocol.Blockchain,
	synchronizer network.Synchronizer,
	transactionsPool protocol.TransactionsPool,
	watch clock.Watch,
	logger log.Logger) *Handler {
	return &Handler{blockchain, synchronizer, transactionsPool, watch, logger}
}

func (handler *Handler) Handle(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
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
			res = handler.blocks()
		case GetTransactions:
			res = handler.transactions()
		default:
			unknownRequest = true
		}
	} else if err = json.Unmarshal(data, &amountRequest); err == nil && !amountRequest.IsInvalid() {
		res = handler.amount(&amountRequest)
	} else if err = json.Unmarshal(data, &transactionRequest); err == nil && !transactionRequest.IsInvalid() {
		go handler.transactionsPool.AddTransaction(&transactionRequest, handler.synchronizer.HostTarget())
	} else if err = json.Unmarshal(data, &lastBlocksRequest); err == nil && !lastBlocksRequest.IsInvalid() {
		res = handler.lastBlocks(&lastBlocksRequest)
	} else if err = json.Unmarshal(data, &targetsRequest); err == nil {
		for _, request := range targetsRequest {
			if request.IsInvalid() {
				unknownRequest = true
			}
		}
		if !unknownRequest {
			go handler.synchronizer.AddTargets(targetsRequest)
		}
	} else {
		unknownRequest = true
	}

	if unknownRequest {
		handler.logger.Error("unknown request")
	}
	return
}

func (handler *Handler) amount(request *network.AmountRequest) (res gp2p.Data) {
	blockchainAddress := *request.Address
	amount := handler.blockchain.Copy().CalculateTotalAmount(handler.watch.Now().UnixNano(), blockchainAddress)
	amountResponse := &network.AmountResponse{Amount: amount}
	data, err := json.Marshal(amountResponse)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
		return
	}
	res.SetBytes(data)
	return
}

func (handler *Handler) blocks() (res gp2p.Data) {
	blockResponses := handler.blockchain.Blocks()
	data, err := json.Marshal(blockResponses)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to get blocks: %w", err).Error())
		return
	}
	res.SetBytes(data)
	return
}

func (handler *Handler) lastBlocks(request *network.LastBlocksRequest) (res gp2p.Data) {
	blockResponses := handler.blockchain.LastBlocks(*request.StartingBlockNonce)
	data, err := json.Marshal(blockResponses)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to get last blocks: %w", err).Error())
		return
	}
	res.SetBytes(data)
	return
}

func (handler *Handler) transactions() (res gp2p.Data) {
	transactionResponses := handler.transactionsPool.Transactions()
	data, err := json.Marshal(transactionResponses)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
		return
	}
	res.SetBytes(data)
	return
}
