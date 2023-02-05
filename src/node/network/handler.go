package network

import (
	"context"
	gp2p "github.com/leprosus/golang-p2p"
)

type Handler interface {
	Handle(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
}
