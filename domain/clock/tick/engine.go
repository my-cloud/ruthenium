package tick

import (
	"time"

	"github.com/my-cloud/ruthenium/domain/clock"
)

type Engine struct {
	function func(timestamp int64)

	watch              clock.Watch
	timer              time.Duration
	subTimer           time.Duration
	ticker             *time.Ticker
	occurrences        int64
	skippedOccurrences int
	started            bool
	requested          bool
}

func NewEngine(function func(timestamp int64), watch clock.Watch, timer time.Duration, occurrences int64, skippedOccurrences int) *Engine {
	var subTimer time.Duration
	if occurrences > 0 {
		subTimer = time.Duration(timer.Nanoseconds() / occurrences)
	} else {
		subTimer = timer
	}
	ticker := time.NewTicker(timer)
	return &Engine{function, watch, timer, subTimer, ticker, occurrences, skippedOccurrences, false, false}
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
	engine.ticker.Reset(engine.subTimer)
	occurrences := int(engine.occurrences)
	for {
		for i := 0; i < occurrences; i++ {
			if i >= engine.skippedOccurrences {
				if !engine.started {
					engine.ticker.Stop()
					return
				}
				now := engine.watch.Now().Round(engine.subTimer)
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
