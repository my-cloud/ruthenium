package network

import p2p "github.com/leprosus/golang-p2p"

type Server interface {
	SetHandle(topic string, handler p2p.Handler)
	Serve() (err error)
}
