package api

import (
	"context"
	"encoding/json"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
)

type UtxosController struct {
	utxosManager ledger.UtxosManager
}

func NewUtxosController(utxosManager ledger.UtxosManager) *UtxosController {
	return &UtxosController{utxosManager}
}

func (controller *UtxosController) HandleUtxosRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	var address string
	data := req.GetBytes()
	if err := json.Unmarshal(data, &address); err != nil {
		return res, err
	}
	utxosByAddress := controller.utxosManager.Utxos(address)
	utxosByAddressBytes, err := json.Marshal(utxosByAddress)
	if err != nil {
		return res, err
	}
	res.SetBytes(utxosByAddressBytes)
	return res, nil
}
