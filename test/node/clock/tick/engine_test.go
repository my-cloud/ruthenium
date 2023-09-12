package tick

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
)

func Test_Do_NoError_FunctionCalled(t *testing.T) {
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

func Test_Start_NotStarted_Started(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var waitGroup sync.WaitGroup
	engine := tick.NewEngine(func(int64) { waitGroup.Done() }, watchMock, 1, 1, 0)

	// Act
	waitGroup.Add(1)
	go engine.Start()
	isFunctionCalled := waitTimeout(&waitGroup, time.Second.Nanosecond())
	engine.Stop()

	// Assert
	test.Assert(t, isFunctionCalled, "The function is not called whereas it should be.")
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}
