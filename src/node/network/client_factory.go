package network

type ClientFactory interface {
	CreateClient(ip string, port uint16, target string) (Client, error)
}
