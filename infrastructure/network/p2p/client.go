package p2p

type Client interface {
	Send(topic string, req []byte) (res []byte, err error)
}
