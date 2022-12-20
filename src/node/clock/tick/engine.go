package tick

import (
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"sync"
	"time"
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

	waitGroup *sync.WaitGroup
	logger    log.Logger
}

func NewEngine(function func(timestamp int64), watch clock.Watch, timer time.Duration, occurrences int, skippedOccurrences int, logger log.Logger) *Engine {
	subTimer := timer
	if occurrences > 0 {
		subTimer = time.Duration(timer.Nanoseconds() / int64(occurrences))
	}
	ticker := time.NewTicker(subTimer)
	var waitGroup sync.WaitGroup
	return &Engine{function, watch, subTimer, ticker, occurrences, skippedOccurrences, false, false, &waitGroup, logger}
}

func (engine *Engine) Do() {
	if engine.started || engine.requested {
		return
	}
	startTime := engine.watch.Now()
	parsedStartDate := startTime.Truncate(engine.timer).Add(engine.timer)
	deadline := parsedStartDate.Sub(startTime)
	engine.ticker.Reset(deadline)
	engine.requested = true
	engine.waitGroup.Add(1)
	go func() {
		defer engine.waitGroup.Done()
		<-engine.ticker.C
		now := engine.watch.Now().Round(engine.timer)
		engine.function(now.UnixNano())
		engine.requested = false
		if engine.started {
			newParsedStartDate := now.Add(engine.timer)
			newDeadline := newParsedStartDate.Sub(now)
			engine.ticker.Reset(newDeadline)
		} else {
			engine.ticker.Stop()
		}
	}()
}

func (engine *Engine) Start() {
	if engine.started {
		return
	}
	engine.started = true
	startTime := engine.watch.Now()
	parsedStartDate := startTime.Truncate(engine.timer).Add(engine.timer)
	deadline := parsedStartDate.Sub(startTime)
	engine.ticker.Reset(deadline)
	<-engine.ticker.C
	engine.ticker.Reset(engine.timer)
	go func() {
		for {
			for i := 0; i < engine.occurrences; i++ {
				if i >= engine.skippedOccurrences {
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
	}()
}

func (engine *Engine) Stop() {
	engine.started = false
	engine.ticker.Reset(time.Nanosecond)
}

func (engine *Engine) Wait() {
	engine.waitGroup.Wait()
}
