package domain

import "time"

type TimeProvider interface {
	Now() time.Time
}
