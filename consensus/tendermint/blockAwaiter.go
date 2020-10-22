package tendermint

import (
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"sync"
)

type blockAwaiter struct {
	valueCond      *sync.Cond
	lastAddedValue *types.Block
	quit           chan struct{}
}

// Create new blockAwaiter per height
func newBlockAwaiter() *blockAwaiter {
	return &blockAwaiter{
		valueCond: sync.NewCond(&sync.Mutex{}),
		quit:      make(chan struct{}),
	}
}

func (a *blockAwaiter) addValue(b *types.Block) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	a.lastAddedValue = b
	// Wake a go routine, if any, waiting on valueCond
	a.valueCond.Signal()
	println("setting value", a.lastAddedValue.Hash().String()[2:8], "value height", a.lastAddedValue.Number().String())
}

// latestValue will return the lastAddedValue set by addValue for the current height. If lastAddedValue is nil or is of
// a previous height then the function will wait until addValue is called or signal to quit is received.
func (a *blockAwaiter) latestValue(height *big.Int) (*types.Block, error) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	for {
		select {
		case <-a.quit:
			return nil, errStopped
		default:
			if a.lastAddedValue == nil || a.lastAddedValue.Number().Cmp(height) != 0 {
				a.lastAddedValue = nil
				if a.lastAddedValue == nil {
					println("awaiting value", "valueIsNil")
				} else {
					println("awaiting value", "value height", a.lastAddedValue.Number().String(), "awaited height", height.String())
				}
				a.valueCond.Wait()
			} else {
				println("received awaited value", a.lastAddedValue.Hash().String()[2:8], "value height", a.lastAddedValue.Number().String(), "awaited height", height.String())
				return a.lastAddedValue, nil
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
