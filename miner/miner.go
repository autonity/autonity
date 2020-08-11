// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package miner implements Ethereum block creation and mining.
package miner

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/hexutil"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
)

// Backend wraps all methods required for mining.
type Backend interface {
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
}

// Config is the configuration parameters of mining.
type Config struct {
	Etherbase common.Address `toml:",omitempty"` // Public address for block mining rewards (default = first account)
	Notify    []string       `toml:",omitempty"` // HTTP URL list to be notified of new work packages(only useful in ethash).
	ExtraData hexutil.Bytes  `toml:",omitempty"` // Block extra data set by the miner
	GasFloor  uint64         // Target gas floor for mined blocks.
	GasCeil   uint64         // Target gas ceiling for mined blocks.
	GasPrice  *big.Int       // Minimum gas price for mining a transaction
	Recommit  time.Duration  // The time interval for miner to re-create mining work.
	Noverify  bool           // Disable remote mining solution verification(only useful in ethash).
}

// Miner creates blocks and searches for proof-of-work values.
type Miner struct {
	mux        *event.TypeMux
	worker     *worker
	coinbase   common.Address
	coinbaseMu sync.RWMutex
	eth        Backend
	engine     consensus.Engine
	exitCh     chan struct{}

	startStopMutex sync.Mutex
	canStart       bool // can start indicates whether we can start the mining operation
	shouldStart    bool // should start indicates whether we should start after sync
}

func New(eth Backend, config *Config, chainConfig *params.ChainConfig, mux *event.TypeMux, engine consensus.Engine, isLocalBlock func(block *types.Block) bool) *Miner {
	miner := &Miner{
		eth:      eth,
		mux:      mux,
		engine:   engine,
		exitCh:   make(chan struct{}),
		worker:   newWorker(config, chainConfig, engine, eth, mux, isLocalBlock, false),
		canStart: true,
	}
	go miner.update()

	return miner
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `Done` or `Failed` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your mining operation for as long as the DOS continues.
func (miner *Miner) update() {
	events := miner.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
	defer events.Unsubscribe()

	for {
		select {
		case ev := <-events.Chan():
			if ev == nil {
				return
			}
			switch ev.Data.(type) {
			case downloader.StartEvent:
				// When syncing begins we pause the miner and set a flag to
				// ensure calls to Start do not start the miner before sync is
				// finished.

				// We need to lock over setting canStart, checking Mining and
				// the call to stop, to prevent the race condition where a
				// prior concurrent call to Start which has passed all its
				// checks then calls start after stop is called here. Resulting
				// in the miner starting immediately after this call to stop
				// whilst canstart is set to false.
				miner.startStopMutex.Lock()
				miner.canStart = false
				if miner.Mining() {
					miner.stop()
					log.Info("Mining aborted due to sync")
				}
				miner.startStopMutex.Unlock()
			case downloader.DoneEvent, downloader.FailedEvent:
				// When syncing completes, we start the miner if it is expected
				// to start, we also unset the flag preventing calls to Start
				// from starting the miner.

				// We need to lock over both the check on shouldStart and the
				// call to start, to prevent the race condition where a
				// concurrent call to Stop occurs between the check of
				// shouldStart and the call to start resulting in the miner
				// being started after Stop has been called.
				miner.startStopMutex.Lock()
				miner.canStart = true
				if miner.shouldStart {
					miner.start()
				}
				miner.startStopMutex.Unlock()
				// stop immediately and ignore all further pending events
				return
			}
		case <-miner.exitCh:
			return
		}
	}
}

// Start starts the miner mining, unless it has been paused by the downloader
// during sync, in which case it will start mining once the sync has completed.
func (miner *Miner) Start(coinbase common.Address) {
	miner.startStopMutex.Lock()
	defer miner.startStopMutex.Unlock()
	miner.SetEtherbase(coinbase)
	miner.shouldStart = true
	if !miner.canStart {
		log.Info("Network syncing, will start miner afterwards")
		return
	}
	miner.start()
}

// start performs the action of starting without managing mutexes or state
// flags.
func (miner *Miner) start() {
	miner.worker.start()
}

// Stop stops the miner from mining.
func (miner *Miner) Stop() {
	miner.startStopMutex.Lock()
	defer miner.startStopMutex.Unlock()
	miner.shouldStart = false
	miner.stop()
}

// stop performs the action of stopping without managing mutexes or state
// flags.
func (miner *Miner) stop() {
	miner.worker.stop()
}

// Close stops the miner and releases any resources associated with it.
func (miner *Miner) Close() {
	miner.Stop()
	miner.worker.close()
	close(miner.exitCh)
}

func (miner *Miner) Mining() bool {
	return miner.worker.isRunning()
}

func (miner *Miner) IsMining() bool {
	miner.startStopMutex.Lock()
	defer miner.startStopMutex.Unlock()
	return miner.shouldStart
}

func (miner *Miner) HashRate() uint64 {
	if pow, ok := miner.engine.(consensus.PoW); ok {
		return uint64(pow.Hashrate())
	}
	return 0
}

func (miner *Miner) SetExtra(extra []byte) error {
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra exceeds max length. %d > %v", len(extra), params.MaximumExtraDataSize)
	}
	miner.worker.setExtra(extra)
	return nil
}

// SetRecommitInterval sets the interval for sealing work resubmitting.
func (miner *Miner) SetRecommitInterval(interval time.Duration) {
	miner.worker.setRecommitInterval(interval)
}

// Pending returns the currently pending block and associated state.
func (miner *Miner) Pending() (*types.Block, *state.StateDB) {
	return miner.worker.pending()
}

// PendingBlock returns the currently pending block.
//
// Note, to access both the pending block and the pending state
// simultaneously, please use Pending(), as the pending state can
// change between multiple method calls
func (miner *Miner) PendingBlock() *types.Block {
	return miner.worker.pendingBlock()
}

func (miner *Miner) SetEtherbase(addr common.Address) {
	miner.coinbase = addr
	miner.worker.setEtherbase(addr)
}

// EnablePreseal turns on the preseal mining feature. It's enabled by default.
// Note this function shouldn't be exposed to API, it's unnecessary for users
// (miners) to actually know the underlying detail. It's only for outside project
// which uses this library.
func (miner *Miner) EnablePreseal() {
	miner.worker.enablePreseal()
}

// DisablePreseal turns off the preseal mining feature. It's necessary for some
// fake consensus engine which can seal blocks instantaneously.
// Note this function shouldn't be exposed to API, it's unnecessary for users
// (miners) to actually know the underlying detail. It's only for outside project
// which uses this library.
func (miner *Miner) DisablePreseal() {
	miner.worker.disablePreseal()
}

// SubscribePendingLogs starts delivering logs from pending transactions
// to the given channel.
func (miner *Miner) SubscribePendingLogs(ch chan<- []*types.Log) event.Subscription {
	return miner.worker.pendingLogsFeed.Subscribe(ch)
}

func (miner *Miner) Coinbase() common.Address {
	miner.coinbaseMu.RLock()
	defer miner.coinbaseMu.RUnlock()
	return miner.coinbase
}

func (miner *Miner) SetCoinbase(addr common.Address) {
	miner.coinbaseMu.Lock()
	miner.coinbase = addr
	miner.coinbaseMu.Unlock()
}
