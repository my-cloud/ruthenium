package api

import (
	"context"
	"encoding/json"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
)

type TransactionsController struct {
	sendersManager      network.SendersManager
	transactionsManager ledger.TransactionsManager
}

func NewTransactionsController(sendersManager network.SendersManager, transactionsManager ledger.TransactionsManager) *TransactionsController {
	return &TransactionsController{sendersManager, transactionsManager}
}

func (controller *TransactionsController) HandleTransactionRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	var transactionRequest *protocol.TransactionRequest
	data := req.GetBytes()
	res := gp2p.Data{}
	if err := json.Unmarshal(data, &transactionRequest); err != nil {
		return res, err
	}
	go controller.transactionsManager.AddTransaction(transactionRequest.Transaction(), transactionRequest.TransactionBroadcasterTarget(), controller.sendersManager.HostTarget())
	return res, nil
}

func (controller *TransactionsController) HandleTransactionsRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	transactions := controller.transactionsManager.Transactions()
	transactionsBytes, err := json.Marshal(transactions)
	if err != nil {
		return res, err
	}
	res.SetBytes(transactionsBytes)
	return res, nil
}
