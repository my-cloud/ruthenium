package gp2p

import (
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
)

type ServerFactory struct{}

func NewServerFactory() *ServerFactory {
	return &ServerFactory{}
}

func (factory *ServerFactory) CreateServer(port int, handler p2p.Handler) (p2p.Server, error) {
	server, err := NewServer(port, handler)
	if err != nil {
		return nil, err
	}
	return server, nil
}
