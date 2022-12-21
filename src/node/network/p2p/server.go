package p2p

import gp2p "github.com/leprosus/golang-p2p"

type Server interface {
	SetHandle(topic string, handler gp2p.Handler)
	Serve() (err error)
}
