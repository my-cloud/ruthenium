package gp2p

import (
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"strconv"
	"time"
)

const serverConnectionTimeoutInSeconds = 10

type ServerFactory struct {
	logger log.Logger
}

func NewServerFactory() *ServerFactory {
	return &ServerFactory{console.NewLogger(console.Fatal)}
}

func (factory *ServerFactory) CreateServer(port int) (p2p.Server, error) {
	tcp := gp2p.NewTCP("0.0.0.0", strconv.Itoa(port))
	server, err := gp2p.NewServer(tcp)
	if err != nil {
		return nil, err
	}
	server.SetLogger(console.NewLogger(console.Fatal))
	settings := gp2p.NewServerSettings()
	settings.SetConnTimeout(serverConnectionTimeoutInSeconds * time.Second)
	server.SetSettings(settings)
	return server, nil
}
