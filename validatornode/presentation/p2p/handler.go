package p2p

import (
	"context"
	"encoding/json"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
)

type Handler struct {
	blocksManager       ledger.BlocksManager
	settings            []byte
	neighborsManager    network.NeighborsManager
	transactionsManager ledger.TransactionsManager
	utxosManager        ledger.UtxosManager
	watch               ledger.TimeProvider
}

func NewHandler(blocksManager ledger.BlocksManager,
	settings []byte,
	neighborsManager network.NeighborsManager,
	transactionsManager ledger.TransactionsManager,
	utxosManager ledger.UtxosManager,
	watch ledger.TimeProvider) *Handler {
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
	blocksBytes, err := json.Marshal(blocks)
	if err != nil {
		return res, err
	}
	res.SetBytes(blocksBytes)
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
	var transactionRequest *protocol.TransactionRequest
	data := req.GetBytes()
	res := gp2p.Data{}
	if err := json.Unmarshal(data, &transactionRequest); err != nil {
		return res, err
	}
	go handler.transactionsManager.AddTransaction(transactionRequest.Transaction(), transactionRequest.TransactionBroadcasterTarget(), handler.neighborsManager.HostTarget())
	return res, nil
}

func (handler *Handler) HandleTransactionsRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	transactions := handler.transactionsManager.Transactions()
	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		return res, err
	}
	res.SetBytes(transactionsBytes)
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
	utxosByAddressBytes, err := json.Marshal(utxosByAddress)
	if err != nil {
		return res, err
	}
	res.SetBytes(utxosByAddressBytes)
	return res, nil
}
