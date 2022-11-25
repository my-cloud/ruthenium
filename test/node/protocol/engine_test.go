package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock/clockmock"
	"testing"
	"time"
)

func Test_Do(t *testing.T) {
	// Arrange
	timeMock := new(clockmock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var calls int
	engine := clock.NewEngine(func(int64) { calls++ }, timeMock, 1, 0, 0, nil)

	// Act
	engine.Do()

	// Assert
	engine.Wait()
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}
