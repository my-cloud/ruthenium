package network

import "time"

type SettingsProvider interface {
	ValidationTimeout() time.Duration
}
