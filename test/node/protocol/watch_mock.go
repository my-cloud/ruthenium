package protocol

import "time"

type WatchMock struct {
	counter int64
}

func NewWatchMock() *WatchMock {
	return &WatchMock{}
}

func (mock *WatchMock) Now() time.Time {
	counter := mock.counter
	mock.counter++
	return time.Unix(0, counter)
}
