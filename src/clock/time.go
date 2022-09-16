package clock

import "time"

type Time interface {
	Now() time.Time
}
