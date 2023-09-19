package network

import (
	"context"
	gp2p "github.com/leprosus/golang-p2p"
)

type Handler interface {
	HandleBlockRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleBlocksRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleTargetsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleTransactionRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleTransactionsRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
	HandleUtxosRequest(_ context.Context, req gp2p.Data) (res gp2p.Data, err error)
}
