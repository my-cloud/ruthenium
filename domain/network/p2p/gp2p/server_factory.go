package gp2p

import (
	"github.com/my-cloud/ruthenium/domain/network/p2p"
)

type ServerFactory struct {
	handler  p2p.Handler
	settings p2p.Settings
}

func NewServerFactory(handler p2p.Handler, settings p2p.Settings) *ServerFactory {
	return &ServerFactory{handler, settings}
}

func (factory *ServerFactory) CreateServer(port int) (p2p.Server, error) {
	server, err := NewServer(port, factory.handler, factory.settings.ValidationTimeout())
	if err != nil {
		return nil, err
	}
	return server, nil
}
