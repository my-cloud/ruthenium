package network

type SenderFactory interface {
	CreateSender(ip string, port uint16, target string) (Sender, error)
}
