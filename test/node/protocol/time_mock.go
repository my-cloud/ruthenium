package protocol

import "time"

type TimeMock struct {
	counter int64
}

func NewTimeMock() *TimeMock {
	return &TimeMock{}
}

func (mock *TimeMock) Now() time.Time {
	counter := mock.counter
	mock.counter++
	return time.Unix(0, counter)
}
