package connection

import p2p "github.com/leprosus/golang-p2p"

type Servable interface {
	SetHandle(topic string, handler p2p.Handler)
	Serve() (err error)
}
