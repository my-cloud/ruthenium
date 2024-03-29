package gp2p

import (
	gp2p "github.com/leprosus/golang-p2p"
	"github.com/my-cloud/ruthenium/src/log"
	"time"
)

type Client struct {
	*gp2p.Client
}

func NewClient(ip string, port string, connectionTimeout time.Duration, logger log.Logger) (*Client, error) {
	tcp := gp2p.NewTCP(ip, port)
	client, err := gp2p.NewClient(tcp)
	if err != nil {
		return nil, err
	}
	settings := gp2p.NewClientSettings()
	settings.SetRetry(1, time.Nanosecond)
	settings.SetConnTimeout(connectionTimeout)
	client.SetSettings(settings)
	client.SetLogger(logger)
	return &Client{client}, err
}

func (client *Client) Send(topic string, req []byte) (res []byte, err error) {
	data, err := client.Client.Send(topic, gp2p.Data{Bytes: req})
	if err != nil {
		return []byte{}, err
	}
	return data.Bytes, err
}
