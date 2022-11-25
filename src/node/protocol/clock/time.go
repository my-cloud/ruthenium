package clock

import "time"

type Time struct{}

func NewTime() *Time {
	return &Time{}
}

func (watch *Time) Now() time.Time {
	return time.Now()
}
