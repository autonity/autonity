package tendermint

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"sync"
)

type blockAwaiter struct {
	valueCond   *sync.Cond
	latestValue *types.Block
	allValues   map[common.Hash]*types.Block
	quit        chan struct{}
}

// Create new blockAwaiter per height
func newBlockAwaiter() *blockAwaiter {
	return &blockAwaiter{
		valueCond: sync.NewCond(&sync.Mutex{}),
		allValues: make(map[common.Hash]*types.Block),
		quit:      make(chan struct{}),
	}
}

func (a *blockAwaiter) setValue(b *types.Block) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	a.latestValue = b
	a.allValues[b.Hash()] = b
	// Wake a go routine, if any, waiting on valueCond
	a.valueCond.Signal()
	println("setting value", a.latestValue.Hash().String()[2:8], "value height", a.latestValue.Number().String())
}

func (a blockAwaiter) value(h common.Hash) *types.Block {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()
	return a.allValues[h]
}

// awaitValue will return the latest value set by setValue for the current height.
func (a *blockAwaiter) awaitValue(height *big.Int) (*types.Block, error) {
	a.valueCond.L.Lock()
	defer a.valueCond.L.Unlock()

	for {
		select {
		case <-a.quit:
			return nil, errStopped
		default:
			if a.latestValue == nil || a.latestValue.Number().Cmp(height) != 0 {
				a.latestValue = nil
				if a.latestValue == nil {
					println("awaiting value", "valueIsNil")
				} else {
					println("awaiting value", "value height", a.latestValue.Number().String(), "awaited height", height.String())
				}
				a.valueCond.Wait()
			} else {
				println("received awaited value", a.latestValue.Hash().String()[2:8], "value height", a.latestValue.Number().String(), "awaited height", height.String())
				v := a.latestValue
				a.latestValue = nil
				return v, nil
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
