package clock

import "time"

type Timeable interface {
	Now() time.Time
}
