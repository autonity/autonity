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
	"reflect"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

const maxSyncFailures = 100

// Backend wraps all methods required for mining. Only full node is capable
// to offer all the functions here.
type Backend interface {
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
	StateAtBlock(block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error)
	Logger() log.Logger
}

// Config is the configuration parameters of mining.
type Config struct {
	Etherbase  common.Address `toml:",omitempty"` // Public address for block mining rewards (default = first account)
	Notify     []string       `toml:",omitempty"` // HTTP URL list to be notified of new work packages (only useful in ethash).
	NotifyFull bool           `toml:",omitempty"` // Notify with pending block headers instead of work packages
	ExtraData  hexutil.Bytes  `toml:",omitempty"` // Block extra data set by the miner
	GasFloor   uint64         // Target gas floor for mined blocks.
	GasCeil    uint64         // Target gas ceiling for mined blocks.
	GasPrice   *big.Int       // Minimum gas price for mining a transaction
	Recommit   time.Duration  // The time interval for miner to re-create mining work.
	Noverify   bool           // Disable remote mining solution verification(only useful in ethash).
}

// Miner creates blocks and searches for proof-of-work values.
type Miner struct {
	mux          *event.TypeMux
	worker       *worker
	eth          Backend
	engine       consensus.Engine
	exitCh       chan struct{}
	startCh      chan struct{}
	stopCh       chan struct{}
	forceStartCh chan struct{}

	wg sync.WaitGroup

	// used in the miner update loop
	canStart    bool
	shouldStart bool
}

func New(eth Backend, config *Config, chainConfig *params.ChainConfig, mux *event.TypeMux, engine consensus.Engine, isLocalBlock func(header *types.Header) bool) *Miner {
	miner := &Miner{
		eth:          eth,
		mux:          mux,
		engine:       engine,
		exitCh:       make(chan struct{}),
		startCh:      make(chan struct{}),
		stopCh:       make(chan struct{}),
		forceStartCh: make(chan struct{}),
		worker:       newWorker(config, chainConfig, engine, eth, mux, isLocalBlock, true),
		shouldStart:  false,
		canStart:     false,
	}
	miner.wg.Add(1)
	go miner.update()
	return miner
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `DoneEvent` or `SyncedEvent` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your mining operation for as long as the DOS continues.
func (miner *Miner) update() {
	defer miner.wg.Done()

	events := miner.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.SyncedEvent{}, downloader.FailedEvent{})
	defer func() {
		if !events.Closed() {
			events.Unsubscribe()
		}
	}()

	// miner.shouldStart is set at initialization to false
	// it will be true when node is a committee member, therefore consensus engine should be started if possible.

	/* miner.canStart is set at initialization to false
	* miner.canStart = true when consensus engine can be safely started OR if we forced mining start.
	* It is safe to start the consensus engine when we are reasonably sure to be synced with the head of the chain:
	* - the first chain sync with our peers terminates, and we conclude that we are synced to the chain head.
	* - the first chain sync with our peer terminates, and we successfully imported blocks till the head of the chain
	*
	* NOTE: This mechanism does not give 100% guarantee that we are synced to the head of the chain before starting consensus. There is always the possibility that we are connected to peers which are behind the "global" chain head.
	 */

	// keeps track of subsequent sync failures at startup sync
	syncFailures := 0

	dlEventCh := events.Chan()
	for {
		select {
		case ev := <-dlEventCh:
			if ev == nil {
				// Unsubscription done, stop listening
				dlEventCh = nil
				continue
			}
			switch ev.Data.(type) {
			case downloader.StartEvent:
				miner.eth.Logger().Info("Chain syncing started, waiting for completion to start consensus engine", "shouldStart", miner.shouldStart, "canStart", miner.canStart)
			case downloader.FailedEvent:
				syncFailures++
				miner.eth.Logger().Info("Chain syncing failed", "#failures", syncFailures, "shouldStart", miner.shouldStart, "canStart", miner.canStart)
				// if we fail more than maxSyncFailures times consequently, assume we are under attack
				if syncFailures >= maxSyncFailures {
					miner.eth.Logger().Warn("************************** PROBLEM DETECTED ******************************")
					miner.eth.Logger().Warn("Multiple sequential chain sync failures detected", "sync failures", syncFailures)
					miner.eth.Logger().Warn("Either you have a network connectivity issue")
					miner.eth.Logger().Warn("Or your node is under attack by malicious peers, which are preventing sync to complete")
					miner.eth.Logger().Warn("Try restarting your node and connecting to a trusted set of peers")
					miner.eth.Logger().Warn("Reach out to Autonity social media channels for support and additional informations")
					miner.eth.Logger().Warn("**************************************************************************")
				}
			// `DoneEvent` deals with the normal scenario:
			// - when starting the node we have some blocks to sync
			// - once finished syncing we can start consensus
			// `SyncedEvent` is needed to start mining if we are already synced with the chain head
			// Example scenario:
			// The chain halts, and we are restarting our offline validator to make it un-halt.
			case downloader.DoneEvent, downloader.SyncedEvent:
				miner.canStart = true
				miner.eth.Logger().Info("Chain syncing completed, consensus engine can start", "event", reflect.TypeOf(ev.Data), "shouldStart", miner.shouldStart, "canStart", miner.canStart)
				miner.startWorker()
				// Stop reacting to downloader events
				if !events.Closed() {
					events.Unsubscribe()
				}
			}
		case <-miner.forceStartCh:
			miner.eth.Logger().Info("Forcing consensus engine start")
			miner.canStart = true
			// don't need to react to downloader events anymore, we don't care about sync status
			if !events.Closed() {
				events.Unsubscribe()
			}
			miner.shouldStart = true
			miner.startWorker()
		// the committeeWatcher in the Ethereum backend will trigger these codepaths, depending on whether we enter/exit the committee
		case <-miner.startCh:
			miner.shouldStart = true
			miner.startWorker()
		case <-miner.stopCh:
			miner.shouldStart = false
			miner.worker.stop()
		case <-miner.exitCh:
			miner.worker.close()
			return
		}
	}
}

func (miner *Miner) startWorker() {
	if !(miner.shouldStart && miner.canStart) {
		return
	}
	miner.worker.start()
}

// force the start of the worker, without waiting for chain sync completion
// this is useful if you want to run a single node network and still do consensus
func (miner *Miner) ForceStart() {
	miner.forceStartCh <- struct{}{}
}

// Start signals that the mining should start
// mining will actually start once we are synced with the network (unless forcing start)
func (miner *Miner) Start() {
	miner.startCh <- struct{}{}
}

func (miner *Miner) Stop() {
	miner.stopCh <- struct{}{}
}

func (miner *Miner) Close() {
	close(miner.exitCh)
	miner.wg.Wait()
}

func (miner *Miner) Mining() bool {
	return miner.worker.isRunning()
}

func (miner *Miner) Hashrate() uint64 {
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

// PendingBlockAndReceipts returns the currently pending block and corresponding receipts.
func (miner *Miner) PendingBlockAndReceipts() (*types.Block, types.Receipts) {
	return miner.worker.pendingBlockAndReceipts()
}

// SetGasCeil sets the gaslimit to strive for when mining blocks post 1559.
// For pre-1559 blocks, it sets the ceiling.
func (miner *Miner) SetGasCeil(ceil uint64) {
	miner.worker.setGasCeil(ceil)
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

// GetSealingBlock retrieves a sealing block based on the given parameters.
// The returned block is not sealed but all other fields should be filled.
func (miner *Miner) GetSealingBlock(parent common.Hash, timestamp uint64, coinbase common.Address, random common.Hash) (*types.Block, error) {
	return miner.worker.getSealingBlock(parent, timestamp, coinbase, random)
}

// SubscribePendingLogs starts delivering logs from pending transactions
// to the given channel.
func (miner *Miner) SubscribePendingLogs(ch chan<- []*types.Log) event.Subscription {
	return miner.worker.pendingLogsFeed.Subscribe(ch)
}
