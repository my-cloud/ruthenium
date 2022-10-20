package protocol

import (
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"sync"
	"time"
)

const genesisAmount uint64 = 100000 * ParticlesCount

type Validation struct {
	address    string
	blockchain *Blockchain
	pool       *Pool

	timing    clock.Timing
	timer     time.Duration
	ticker    *time.Ticker
	started   bool
	requested bool

	waitGroup *sync.WaitGroup
	logger    *log.Logger
}

func NewValidation(address string, blockchain *Blockchain, pool *Pool, timing clock.Timing, timer time.Duration, logger *log.Logger) *Validation {
	ticker := time.NewTicker(timer)
	var waitGroup sync.WaitGroup
	return &Validation{address, blockchain, pool, timing, timer, ticker, false, false, &waitGroup, logger}
}

func (validation *Validation) Do() {
	if validation.started || validation.requested {
		return
	}
	startTime := validation.timing.Now()
	parsedStartDate := startTime.Truncate(validation.timer).Add(validation.timer)
	deadline := parsedStartDate.Sub(startTime)
	validation.ticker.Reset(deadline)
	validation.requested = true
	validation.waitGroup.Add(1)
	go func() {
		defer validation.waitGroup.Done()
		<-validation.ticker.C
		now := validation.timing.Now().Round(validation.timer)
		validation.do(now.UnixNano())
		validation.requested = false
		if validation.started {
			newParsedStartDate := now.Add(validation.timer)
			newDeadline := newParsedStartDate.Sub(now)
			validation.ticker.Reset(newDeadline)
		} else {
			validation.ticker.Stop()
		}
	}()
}

func (validation *Validation) Start() {
	if validation.started {
		return
	}
	validation.started = true
	startTime := validation.timing.Now()
	parsedStartDate := startTime.Truncate(validation.timer).Add(validation.timer)
	deadline := parsedStartDate.Sub(startTime)
	validation.ticker.Reset(deadline)
	<-validation.ticker.C
	validation.ticker.Reset(validation.timer)
	go func() {
		for {
			if !validation.started {
				validation.ticker.Stop()
				return
			}
			now := validation.timing.Now().Round(validation.timer)
			validation.do(now.UnixNano())
			<-validation.ticker.C
		}
	}()
}

func (validation *Validation) Stop() {
	validation.started = false
	validation.ticker.Reset(time.Nanosecond)
}

func (validation *Validation) Wait() {
	validation.waitGroup.Wait()
}

func (validation *Validation) do(timestamp int64) {
	if validation.blockchain.IsEmpty() {
		genesisTransaction := NewRewardTransaction(validation.address, timestamp, genesisAmount)
		transactions := []*Transaction{genesisTransaction}
		genesisBlock := NewBlock(timestamp, [32]byte{}, transactions, nil)
		validation.blockchain.AddBlock(genesisBlock)
		validation.logger.Debug("genesis block added")
	} else {
		validation.pool.Validate(timestamp, validation.blockchain, validation.address)
	}
}
