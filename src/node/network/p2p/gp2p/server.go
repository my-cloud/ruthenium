package gp2p

import (
	"fmt"
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"strconv"
	"time"
)

const serverConnectionTimeoutInSeconds = 10 // TODO calculate from validation timestamp

type Server struct {
	*gp2p.Server
	handler p2p.Handler
}

func NewServer(port int, handler p2p.Handler) (*Server, error) {
	tcp := gp2p.NewTCP("0.0.0.0", strconv.Itoa(port))
	server, err := gp2p.NewServer(tcp)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate server on port %d: %w", port, err)
	}
	server.SetLogger(console.NewLogger(console.Fatal))
	settings := gp2p.NewServerSettings()
	settings.SetConnTimeout(serverConnectionTimeoutInSeconds * time.Second)
	server.SetSettings(settings)
	return &Server{server, handler}, err
}

func (server *Server) SetHandleBlocksRequest(endpoint string) {
	server.SetHandle(endpoint, server.handler.HandleBlocksRequest)
}

func (server *Server) SetHandleFirstBlockTimestampRequest(endpoint string) {
	server.SetHandle(endpoint, server.handler.HandleFirstBlockTimestampRequest)
}

func (server *Server) SetHandleTargetsRequest(endpoint string) {
	server.SetHandle(endpoint, server.handler.HandleTargetsRequest)
}

func (server *Server) SetHandleTransactionRequest(endpoint string) {
	server.SetHandle(endpoint, server.handler.HandleTransactionRequest)
}

func (server *Server) SetHandleTransactionsRequest(endpoint string) {
	server.SetHandle(endpoint, server.handler.HandleTransactionsRequest)
}

func (server *Server) SetHandleUtxosRequest(endpoint string) {
	server.SetHandle(endpoint, server.handler.HandleUtxosRequest)
}
