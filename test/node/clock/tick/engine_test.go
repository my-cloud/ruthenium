package tick

import (
	"fmt"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
)

func Test_Do(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var calls int
	engine := tick.NewEngine(func(int64) { calls++ }, watchMock, 1, 0, 0)

	// Act
	engine.Do()

	// Assert
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}

func Test_StartAndStop(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var calls int
	timer := time.Duration(1000)
	engine := tick.NewEngine(func(int64) { calls++ }, watchMock, timer, 1, 0)

	// Act
	go engine.Start()
	time.Sleep(timer)
	engine.Stop()

	// Assert
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}
