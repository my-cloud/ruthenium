package p2p

import p2p "github.com/leprosus/golang-p2p"

type Client interface {
	Send(topic string, req p2p.Data) (res p2p.Data, err error)
}
