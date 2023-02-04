package p2p

type ClientFactory interface {
	CreateClient(ip string, port string) (Client, error)
}
