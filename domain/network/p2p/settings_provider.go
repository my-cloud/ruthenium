package p2p

import "time"

type SettingsProvider interface {
	ValidationTimeout() time.Duration
}
