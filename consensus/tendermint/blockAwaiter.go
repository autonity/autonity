package tendermint

import (
	"sync"
	time "time"

	"github.com/clearmatics/autonity/core/types"
)

type blockAwaiter struct {
	valueCond *sync.Cond
	lastValue *types.Block
	dlog      *debugLog // for debug
	stopped   bool
}

func newBlockAwaiter(dlog *debugLog) *blockAwaiter {
	return &blockAwaiter{
		valueCond: sync.NewCond(&sync.Mutex{}),
		dlog:      dlog,
	}
}

func (a *blockAwaiter) setValue(b *types.Block) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	a.lastValue = b
	// Wake a go routine, if any, waiting on valueCond
	a.valueCond.Signal()
	a.dlog.print("setting value", a.lastValue.Hash().String()[2:8], "value height", a.lastValue.Number().String())
}

// value will return the lastValue set by setValue for the current height. If lastValue is nil or is of a previous
// height then the function will wait until setValue is called or signal to quit is received.
func (a *blockAwaiter) value(height uint64) (*types.Block, error) {
	start := time.Now()
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	a.dlog.print("beginning awaiting value", height)
	for {
		secondsWaited := time.Since(start) / time.Second
		if a.stopped {
			return nil, errStopped
		}
		if a.lastValue == nil || a.lastValue.Number().Uint64() != height {
			if a.lastValue == nil {
				a.dlog.print("awaiting value", "valueIsNil", "awaited height", height, "waiting for", secondsWaited)
			} else {
				a.dlog.print("awaiting value", a.lastValue.Hash().String()[2:8], "value height", a.lastValue.Number().String(), "awaited height", height, "waiting for", secondsWaited)
			}
			a.lastValue = nil
			a.valueCond.Wait()
		} else {
			a.dlog.print("received awaited value", a.lastValue.Hash().String()[2:8], "value height", a.lastValue.Number().String(), "awaited height", height, "waited for", secondsWaited)
			return a.lastValue, nil
		}
	}
}

func (a *blockAwaiter) stop() {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()
	a.stopped = true
	a.valueCond.Signal()
}

func (a *blockAwaiter) start() {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()
	a.stopped = false
}
