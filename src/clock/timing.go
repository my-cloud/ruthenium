package clock

import "time"

type Timing interface {
	Now() time.Time
}
