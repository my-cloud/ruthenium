package configuration

import "time"

type RegistrySettingsProvider interface {
	SynchronizationTimer() time.Duration
}
