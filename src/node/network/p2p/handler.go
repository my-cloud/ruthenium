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

func (handler *Handler) HandleBlocksRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	var startingBlockHeight uint64
	res := gp2p.Data{}
	data := req.GetBytes()
	if err := json.Unmarshal(data, &startingBlockHeight); err != nil {
		handler.logger.Debug(BadRequest)
		return res, err
	}
	blocks := handler.blockchain.Blocks(startingBlockHeight)
	res.SetBytes(blocks)
	return res, nil
}

func (handler *Handler) HandleFirstBlockTimestampRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	timestamp := handler.blockchain.FirstBlockTimestamp()
	timestampBytes, err := json.Marshal(timestamp)
	if err != nil {
		return res, err
	}
	res.SetBytes(timestampBytes)
	return res, nil
}

func (handler *Handler) HandleTargetsRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	var targets []string
	res := gp2p.Data{}
	data := req.GetBytes()
	if err := json.Unmarshal(data, &targets); err != nil {
		handler.logger.Debug(BadRequest)
		return res, err
	}
	go handler.synchronizer.AddTargets(targets)
	return res, nil
}

func (handler *Handler) HandleTransactionRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	var transactionRequest network.TransactionRequest
	res := gp2p.Data{}
	data := req.GetBytes()
	if err := json.Unmarshal(data, &transactionRequest); err != nil {
		handler.logger.Debug(BadRequest)
		return res, err
	}
	if transactionRequest.IsInvalid() {
		handler.logger.Debug(BadRequest)
		err := errors.New(BadRequest)
		return res, err
	}
	go handler.transactionsPool.AddTransaction(&transactionRequest, handler.synchronizer.HostTarget())
	return res, nil
}

func (handler *Handler) HandleTransactionsRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	transactions := handler.transactionsPool.Transactions()
	res.SetBytes(transactions)
	return res, nil
}

func (handler *Handler) HandleUtxosRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	var address string
	data := req.GetBytes()
	if err := json.Unmarshal(data, &address); err != nil {
		handler.logger.Debug(BadRequest)
		return res, err
	}
	utxosByAddress := handler.blockchain.UtxosByAddress(address)
	marshaledUtxos, err := json.Marshal(utxosByAddress)
	if err != nil {
		handler.logger.Error(err.Error())
		return res, err
	}
	res.SetBytes(marshaledUtxos)
	return res, nil
}
