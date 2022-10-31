package p2p

import (
	p2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"strconv"
	"time"
)

const serverConnectionTimeoutInSeconds = 10

type ServerFactory struct {
	logger p2p.Logger
}

func NewServerFactory() *ServerFactory {
	return &ServerFactory{log.NewLogger(log.Fatal)}
}

func (factory *ServerFactory) CreateServer(port int) (network.Server, error) {
	tcp := p2p.NewTCP("0.0.0.0", strconv.Itoa(port))
	server, err := p2p.NewServer(tcp)
	if err != nil {
		return nil, err
	}
	server.SetLogger(log.NewLogger(log.Fatal))
	settings := p2p.NewServerSettings()
	settings.SetConnTimeout(serverConnectionTimeoutInSeconds * time.Second)
	server.SetSettings(settings)
	return server, nil
}
