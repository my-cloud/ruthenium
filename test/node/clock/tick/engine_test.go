package tick

import (
	"fmt"
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
	caller := NewActionCaller()
	function := func(int64) { caller.CallAction() }
	engine := tick.NewEngine(function, watchMock, 1, 1, 0)
	var calls int
	caller.SetAction(func() {
		calls++
		engine.Stop()
	})

	// Act
	engine.Start()

	// Assert
	test.Assert(t, calls == 1, fmt.Sprintf("The function is called %d times whereas it should be called once.", calls))
}

type ActionsCaller struct {
	action func()
}

func NewActionCaller() *ActionsCaller {
	return &ActionsCaller{}
}

func (actionsCaller *ActionsCaller) SetAction(action func()) {
	actionsCaller.action = action
}

func (actionsCaller *ActionsCaller) CallAction() {
	actionsCaller.action()
}
