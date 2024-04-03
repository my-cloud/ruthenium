package p2p

import (
	"context"
	"encoding/json"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol"
)

type Handler struct {
	blocksManager       protocol.BlocksManager
	settings            []byte
	neighborsManager    network.NeighborsManager
	transactionsManager protocol.TransactionsManager
	utxosManager        protocol.UtxosManager
	watch               protocol.TimeProvider
}

func NewHandler(blocksManager protocol.BlocksManager,
	settings []byte,
	neighborsManager network.NeighborsManager,
	transactionsManager protocol.TransactionsManager,
	utxosManager protocol.UtxosManager,
	watch protocol.TimeProvider) *Handler {
	return &Handler{blocksManager, settings, neighborsManager, transactionsManager, utxosManager, watch}
}

func (handler *Handler) HandleBlocksRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	var startingBlockHeight uint64
	res := gp2p.Data{}
	data := req.GetBytes()
	if err := json.Unmarshal(data, &startingBlockHeight); err != nil {
		return res, err
	}
	blocks := handler.blocksManager.Blocks(startingBlockHeight)
	res.SetBytes(blocks)
	return res, nil
}

func (handler *Handler) HandleFirstBlockTimestampRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	timestamp := handler.blocksManager.FirstBlockTimestamp()
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
		return res, err
	}
	go handler.neighborsManager.AddTargets(targets)
	return res, nil
}

func (handler *Handler) HandleTransactionRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	data := req.GetBytes()
	go handler.transactionsManager.AddTransaction(data, handler.neighborsManager.HostTarget())
	return gp2p.Data{}, nil
}

func (handler *Handler) HandleTransactionsRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	transactions := handler.transactionsManager.Transactions()
	res.SetBytes(transactions)
	return res, nil
}

func (handler *Handler) HandleUtxosRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	var address string
	data := req.GetBytes()
	if err := json.Unmarshal(data, &address); err != nil {
		return res, err
	}
	utxosByAddress := handler.utxosManager.Utxos(address)
	res.SetBytes(utxosByAddress)
	return res, nil
}
