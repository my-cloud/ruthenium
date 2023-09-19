package p2p

import (
	"context"
	"encoding/json"
	"errors"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
)

const BadRequest = "bad request"

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
	var blockHeight uint64
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &blockHeight); err == nil {
		blockResponse := handler.blockchain.Block(blockHeight)
		marshaledBlock, err := json.Marshal(blockResponse)
		if err != nil {
			handler.logger.Error(err.Error())
			return
		}
		res.SetBytes(marshaledBlock)
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}

func (handler *Handler) HandleBlocksRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	var startingBlockHeight uint64
	res = gp2p.Data{}
	data := req.GetBytes()
	if err = json.Unmarshal(data, &startingBlockHeight); err == nil {
		blockResponses := handler.blockchain.Blocks(startingBlockHeight)
		marshaledBlocks, err := json.Marshal(blockResponses)
		if err != nil {
			handler.logger.Error(err.Error())
			return
		}
		res.SetBytes(marshaledBlocks)
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
		for _, request := range targetsRequest {
			if request.IsInvalid() {
				err = errors.New("invalid targets request")
				handler.logger.Error(err.Error())
				return
			}
		}
		go handler.synchronizer.AddTargets(targetsRequest)
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

func (handler *Handler) HandleTransactionsRequest(_ context.Context, _ gp2p.Data) (res gp2p.Data, err error) {
	transactionResponses := handler.transactionsPool.Transactions()
	data, err := json.Marshal(transactionResponses)
	if err != nil {
		handler.logger.Error(err.Error())
		return
	}
	res.SetBytes(data)
	return
}

func (handler *Handler) HandleUtxosRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error) {
	res = gp2p.Data{}
	var address string
	data := req.GetBytes()
	if err = json.Unmarshal(data, &address); err == nil {
		utxosByAddress := handler.blockchain.UtxosByAddress(address)
		marshaledUtxos, err := json.Marshal(utxosByAddress)
		if err != nil {
			handler.logger.Error(err.Error())
			return
		}
		res.SetBytes(marshaledUtxos)
	} else {
		handler.logger.Debug(BadRequest)
	}
	return
}
