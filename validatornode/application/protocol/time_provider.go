package protocol

import "time"

type TimeProvider interface {
	Now() time.Time
}
