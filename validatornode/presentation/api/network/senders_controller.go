package network

import (
	"context"
	"encoding/json"
	"github.com/my-cloud/ruthenium/validatornode/application"

	gp2p "github.com/leprosus/golang-p2p"
)

type SendersController struct {
	sendersManager application.SendersManager
}

func NewSendersController(sendersManager application.SendersManager) *SendersController {
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
