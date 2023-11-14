package clock

import "time"

type Watch interface {
	Now() time.Time
}
