package clock

type Engine interface {
	Start()
	Stop()
	Do()
	Wait()
}
