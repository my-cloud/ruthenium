package protocol

type Bootable interface {
	Start()
	Stop()
	Do()
}
