package ledger

import "time"

type TimeProvider interface {
	Now() time.Time
}
