package p2p

import "time"

type Settings interface {
	ValidationTimeout() time.Duration
}
