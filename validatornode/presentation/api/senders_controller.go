package api

import (
	"context"
	"encoding/json"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
)

type SendersController struct {
	sendersManager network.SendersManager
}

func NewSendersController(sendersManager network.SendersManager) *SendersController {
	return &SendersController{sendersManager}
}

func (controller *SendersController) HandleTargetsRequest(_ context.Context, req gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	var targets []string
	data := req.GetBytes()
	if err := json.Unmarshal(data, &targets); err != nil {
		return res, err
	}
	go controller.sendersManager.AddTargets(targets)
	return res, nil
}
