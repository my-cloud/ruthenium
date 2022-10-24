package connection

type SenderProvider interface {
	CreateSender(ip string, port uint16, target string) (Sender, error)
}
