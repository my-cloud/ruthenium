package node

import "time"

type WatchMock struct {
	counter int64
}

func NewWatch() *WatchMock {
	return &WatchMock{}
}

func (watch *WatchMock) Now() time.Time {
	counter := watch.counter
	watch.counter++
	return time.Unix(0, counter)
}
