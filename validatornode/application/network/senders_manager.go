package network

type SendersManager interface {
	AddTargets(targets []string)
	HostTarget() string
	Incentive(target string)
	Senders() []Sender
}
