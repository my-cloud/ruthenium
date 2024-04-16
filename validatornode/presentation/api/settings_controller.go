package api

import (
	"context"

	gp2p "github.com/leprosus/golang-p2p"
)

type SettingsController struct {
	settings []byte
}

func NewSettingsController(settings []byte) *SettingsController {
	return &SettingsController{settings}
}

func (controller *SettingsController) HandleSettingsRequest(_ context.Context, _ gp2p.Data) (gp2p.Data, error) {
	res := gp2p.Data{}
	res.SetBytes(controller.settings)
	return res, nil
}
