package api

import (
	"context"
	"encoding/json"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
)

type BlocksController struct {
	blocksManager ledger.BlocksManager
}

func NewBlocksController(blocksManager ledger.BlocksManager) *BlocksController {
	return &BlocksController{blocksManager}
}

func (controller *BlocksController) HandleBlocksRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	var startingBlockHeight uint64
	res := gp2p.Data{}
	data := req.GetBytes()
	if err := json.Unmarshal(data, &startingBlockHeight); err != nil {
		return res, err
	}
	blocks := controller.blocksManager.Blocks(startingBlockHeight)
	blocksBytes, err := json.Marshal(blocks)
	if err != nil {
		return res, err
	}
	res.SetBytes(blocksBytes)
	return res, nil
}

func (controller *BlocksController) HandleFirstBlockTimestampRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	timestamp := controller.blocksManager.FirstBlockTimestamp()
	timestampBytes, err := json.Marshal(timestamp)
	if err != nil {
		return res, err
	}
	res.SetBytes(timestampBytes)
	return res, nil
}
