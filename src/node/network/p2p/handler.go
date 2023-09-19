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
	GetTransactions = "GET TRANSACTIONS"
	BadRequest      = "bad request"
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

func (handler *Handler) HandleBlockRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var blockRequest network.BlockRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &blockRequest); err == nil && !blockRequest.IsInvalid() {
		res = handler.block(&blockRequest)
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) HandleBlocksRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var blocksRequest network.BlocksRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &blocksRequest); err == nil && !blocksRequest.IsInvalid() {
		res = handler.blocks(&blocksRequest)
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) HandleTargetsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var targetsRequest []network.TargetRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &targetsRequest); err == nil {
		handler.addTargets(targetsRequest)
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) HandleTransactionRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var transactionRequest network.TransactionRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &transactionRequest); err == nil && !transactionRequest.IsInvalid() {
		go handler.transactionsPool.AddTransaction(&transactionRequest, handler.synchronizer.HostTarget())
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) HandleTransactionsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var requestString string
	res = gp2p.Data{}
	data := req.GetBytes()
	err = json.Unmarshal(data, &requestString) // TODO remove useless request string
	if err == nil && requestString == GetTransactions {
		res = handler.transactions()
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) HandleUtxosRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var utxosRequest network.UtxosRequest
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &utxosRequest); err == nil && !utxosRequest.IsInvalid() {
		res = handler.utxos(&utxosRequest)
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) addTargets(targetsRequest []network.TargetRequest) {
	for _, request := range targetsRequest {
		if request.IsInvalid() {
			handler.logger.Error("invalid targets request")
			return
		}
	}
	go handler.synchronizer.AddTargets(targetsRequest)
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
