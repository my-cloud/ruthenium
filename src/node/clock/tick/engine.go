package tick

import (
	"time"

	"github.com/my-cloud/ruthenium/src/node/clock"
)

type Engine struct {
	function func(timestamp int64)

	watch              clock.Watch
	timer              time.Duration
	ticker             *time.Ticker
	occurrences        int
	skippedOccurrences int
	started            bool
	requested          bool
}

func NewEngine(function func(timestamp int64), watch clock.Watch, timer time.Duration, occurrences int, skippedOccurrences int) *Engine {
	subTimer := timer
	if occurrences > 0 {
		subTimer = time.Duration(timer.Nanoseconds() / int64(occurrences))
	}
	ticker := time.NewTicker(subTimer)
	return &Engine{function, watch, subTimer, ticker, occurrences, skippedOccurrences, false, false}
}

func (engine *Engine) Do() {
	if engine.started || engine.requested {
		return
	}
	now := engine.watch.Now()
	startTime := now.Truncate(engine.timer).Add(engine.timer)
	deadline := startTime.Sub(now)
	engine.ticker.Reset(deadline)
	engine.requested = true
	<-engine.ticker.C
	engine.function(startTime.UnixNano())
	engine.requested = false
	if engine.started {
		newParsedStartDate := startTime.Add(engine.timer)
		newDeadline := newParsedStartDate.Sub(startTime)
		engine.ticker.Reset(newDeadline)
	} else {
		engine.ticker.Stop()
	}
}

func (engine *Engine) Start() {
	if engine.started {
		return
	}
	engine.started = true
	initialTime := engine.watch.Now()
	startTime := initialTime.Truncate(engine.timer).Add(engine.timer)
	deadline := startTime.Sub(initialTime)
	engine.ticker.Reset(deadline)
	<-engine.ticker.C
	engine.ticker.Reset(engine.timer)
	for {
		for i := 0; i < engine.occurrences; i++ {
			if i < engine.occurrences-engine.skippedOccurrences {
				if !engine.started {
					engine.ticker.Stop()
					return
				}
				now := engine.watch.Now().Round(engine.timer)
				engine.function(now.UnixNano())
			}
			<-engine.ticker.C
		}
	}
}

func (engine *Engine) Stop() {
	engine.started = false
	engine.ticker.Reset(time.Nanosecond)
}
