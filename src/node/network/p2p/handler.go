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

const GetTransactions = "GET TRANSACTIONS"

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
	var blockRequest network.BlockRequest
	var blocksRequest network.BlocksRequest
	var transactionRequest network.TransactionRequest
	var utxosRequest network.UtxosRequest
	var targetsRequest []network.TargetRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &requestString); err == nil {
		switch requestString {
		case GetTransactions:
			res = handler.transactions()
		default:
			unknownRequest = true
		}
	} else if err = json.Unmarshal(data, &blockRequest); err == nil && !blockRequest.IsInvalid() {
		res = handler.block(&blockRequest)
	} else if err = json.Unmarshal(data, &transactionRequest); err == nil && !transactionRequest.IsInvalid() {
		go handler.transactionsPool.AddTransaction(&transactionRequest, handler.synchronizer.HostTarget())
	} else if err = json.Unmarshal(data, &blocksRequest); err == nil && !blocksRequest.IsInvalid() {
		res = handler.blocks(&blocksRequest)
	} else if err = json.Unmarshal(data, &utxosRequest); err == nil && !utxosRequest.IsInvalid() {
		res = handler.utxos(&utxosRequest)
	} else if err = json.Unmarshal(data, &targetsRequest); err == nil {
		handler.addTargets(targetsRequest, unknownRequest)
	} else {
		unknownRequest = true
	}

	if unknownRequest {
		handler.logger.Debug("unknown request")
	}
	return
}

func (handler *Handler) addTargets(targetsRequest []network.TargetRequest, invalidRequest bool) {
	for _, request := range targetsRequest {
		if request.IsInvalid() {
			invalidRequest = true
			break
		}
	}
	if !invalidRequest {
		go handler.synchronizer.AddTargets(targetsRequest)
	}
}

func (handler *Handler) utxos(request *network.UtxosRequest) (res gp2p.Data) {
	utxosByAddress := handler.blockchain.UtxosByAddress(*request.Address)
	data, err := json.Marshal(utxosByAddress)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
		return
	}
	res.SetBytes(data)
	return
}

func (handler *Handler) block(request *network.BlockRequest) (res gp2p.Data) {
	blockResponse := handler.blockchain.Block(*request.BlockHeight)
	data, err := json.Marshal(blockResponse)
	if err != nil {
		handler.logger.Error(fmt.Errorf("failed to get block: %w", err).Error())
		return
	}
	res.SetBytes(data)
	return
}

func (handler *Handler) blocks(request *network.BlocksRequest) (res gp2p.Data) {
	blockResponses := handler.blockchain.Blocks(*request.StartingBlockHeight)
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
