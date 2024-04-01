package application

import "time"

type TimeProvider interface {
	Now() time.Time
}
