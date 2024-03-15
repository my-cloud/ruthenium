package p2p

type Sender interface {
	Send(topic string, req []byte) (res []byte, err error)
}
