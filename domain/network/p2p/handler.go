package p2p

import (
	"context"
	gp2p "github.com/leprosus/golang-p2p"
)

type Handler interface {
	HandleBlocksRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleFirstBlockTimestampRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleSettingsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleTargetsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleTransactionRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleTransactionsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleUtxosRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
}
