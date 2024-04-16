package application

type SenderCreator interface {
	CreateSender(ip string, port string) (Sender, error)
}
