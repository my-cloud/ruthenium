package tick

import (
	"fmt"
	"github.com/my-cloud/ruthenium/domain/clock"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/infrastructure/test"
)

func Test_Do_NoError_FunctionCalled(t *testing.T) {
	// Arrange
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var calls int
	engine := NewEngine(func(int64) { calls++ }, watchMock, 1, 0, 0)

	// Act
	engine.Do()

	// Assert
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}

func Test_Start_NotStarted_Started(t *testing.T) {
	// Arrange
	watchMock := new(clock.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var engine = &Engine{}
	var calls int
	engine = NewEngine(func(int64) { calls++; engine.Stop() }, watchMock, 1, 1, 0)

	// Act
	engine.Start()

	// Assert
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}
