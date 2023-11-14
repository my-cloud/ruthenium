package tick

import "time"

type Watch struct{}

func NewWatch() *Watch {
	return &Watch{}
}

func (watch *Watch) Now() time.Time {
	return time.Now()
}
