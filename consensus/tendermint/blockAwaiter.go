package tendermint

import (
	"sync"

	"github.com/clearmatics/autonity/core/types"
)

type blockAwaiter struct {
	valueCond *sync.Cond
	lastValue *types.Block
	quit      chan struct{}
}

func newBlockAwaiter() *blockAwaiter {
	return &blockAwaiter{
		valueCond: sync.NewCond(&sync.Mutex{}),
		quit:      make(chan struct{}, 1),
	}
}

func (a *blockAwaiter) setValue(b *types.Block) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	a.lastValue = b
	// Wake a go routine, if any, waiting on valueCond
	a.valueCond.Signal()
	println("setting value", a.lastValue.Hash().String()[2:8], "value height", a.lastValue.Number().String())
}

// value will return the lastValue set by setValue for the current height. If lastValue is nil or is of a previous
// height then the function will wait until setValue is called or signal to quit is received.
func (a *blockAwaiter) value(height uint64) (*types.Block, error) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	println("beginning awaiting value", height)
	for {
		select {
		case <-a.quit:
			return nil, errStopped
		default:
			if a.lastValue == nil || a.lastValue.Number().Uint64() != height {
				a.lastValue = nil
				if a.lastValue == nil {
					println("awaiting value", "valueIsNil")
				} else {
					println("awaiting value", "value height", a.lastValue.Number().String(), "awaited height", height)
				}
				a.valueCond.Wait()
			} else {
				println("received awaited value", a.lastValue.Hash().String()[2:8], "value height", a.lastValue.Number().String(), "awaited height", height)
				return a.lastValue, nil
			}
		}
	}
}

func (a *blockAwaiter) stop() {
	a.quit <- struct{}{}
	a.valueCond.L.Lock()
	a.valueCond.Signal()
	a.valueCond.L.Unlock()
}
