package p2p

type ClientFactory interface {
	CreateClient(ip string, port uint16) (Client, error)
}
