package p2p

type SenderCreator interface {
	CreateSender(ip string, port string) (Sender, error)
}
