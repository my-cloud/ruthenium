package gp2p

import (
	"context"
	"encoding/json"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/domain/validatornode"
	"github.com/my-cloud/ruthenium/infrastructure/clock"
	"github.com/my-cloud/ruthenium/infrastructure/log"
	"github.com/my-cloud/ruthenium/infrastructure/network"
)

const BadRequest = "bad request"

type Handler struct {
	blockchain       validatornode.Blockchain
	settings         []byte
	synchronizer     network.Synchronizer
	transactionsPool validatornode.TransactionsPool
	watch            clock.Watch
	logger           log.Logger
}

func NewHandler(blockchain validatornode.Blockchain,
	settings []byte,
	synchronizer network.Synchronizer,
	transactionsPool validatornode.TransactionsPool,
	watch clock.Watch,
	logger log.Logger) *Handler {
	return &Handler{blockchain, settings, synchronizer, transactionsPool, watch, logger}
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

func (handler *Handler) HandleSettingsRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	res.SetBytes(handler.settings)
	return res, nil
}

func (handler *Handler) HandleTargetsRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	var targets []string
	data := req.GetBytes()
	if err := json.Unmarshal(data, &targets); err != nil {
		handler.logger.Debug(BadRequest)
		return res, err
	}
	go handler.synchronizer.AddTargets(targets)
	return res, nil
}

func (handler *Handler) HandleTransactionRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	data := req.GetBytes()
	go handler.transactionsPool.AddTransaction(data, handler.synchronizer.HostTarget())
	return gp2p.Data{}, nil
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
	utxosByAddress := handler.blockchain.Utxos(address)
	res.SetBytes(utxosByAddress)
	return res, nil
}
