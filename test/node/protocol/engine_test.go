package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/clock/node"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock/clocktest"
	"testing"
	"time"
)

func Test_Do(t *testing.T) {
	// Arrange
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 0) }
	var calls int
	engine := node.NewEngine(func(int64) { calls++ }, watchMock, 1, 0, 0, nil)

	// Act
	engine.Do()

	// Assert
	engine.Wait()
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}
