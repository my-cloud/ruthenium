package p2p

import gp2p "github.com/leprosus/golang-p2p"

type Client interface {
	Send(topic string, req gp2p.Data) (res gp2p.Data, err error)
	SetSettings(stg *gp2p.ClientSettings)
}
