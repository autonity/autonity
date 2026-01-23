// Copyright 2015 The go-ethereum Authors
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

package miner

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/holiman/uint256"

	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/core/txpool"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/eth/ethconfig"
	"github.com/autonity/autonity/metrics"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/misc"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

const (
	// resultQueueSize is the size of channel listening to sealing result.
	resultQueueSize = 10

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10

	// resubmitAdjustChanSize is the size of resubmitting interval adjustment channel.
	resubmitAdjustChanSize = 10

	// sealingLogAtDepth is the number of confirmations before logging successful sealing.
	sealingLogAtDepth = 7

	// minRecommitInterval is the minimal time interval to recreate the sealing block with
	// any newly arrived transactions.
	minRecommitInterval = 1 * time.Second

	// maxRecommitInterval is the maximum time interval to recreate the sealing block with
	// any newly arrived transactions.
	maxRecommitInterval = 15 * time.Second

	// intervalAdjustRatio is the impact a single interval adjustment has on sealing work
	// resubmitting interval.
	intervalAdjustRatio = 0.1

	// intervalAdjustBias is applied during the new resubmit interval calculation in favor of
	// increasing upper limit or decreasing lower limit so that the limit can be reachable.
	intervalAdjustBias = 200 * 1000.0 * 1000.0

	// staleThreshold is the maximum depth of the acceptable stale block.
	staleThreshold = 7
)

// environment is the worker's current environment and holds all
// information of the sealing block generation.
type environment struct {
	signer types.Signer

	state    *state.StateDB // apply state changes here
	tcount   int            // tx count in cycle
	gasPool  *core.GasPool  // available gas used to pack transactions
	coinbase common.Address

	header       *types.Header
	txs          []*types.Transaction
	receipts     []*types.Receipt
	parentHeader *types.Header

	contractsConfig *types.ContractsConfig
	evm             *vm.EVM
}

// copy creates a deep copy of environment.
func (env *environment) copy() *environment {
	cpy := &environment{
		signer:          env.signer,
		state:           env.state.Copy(),
		tcount:          env.tcount,
		coinbase:        env.coinbase,
		header:          types.CopyHeader(env.header),
		receipts:        copyReceipts(env.receipts),
		parentHeader:    types.CopyHeader(env.parentHeader),
		contractsConfig: env.contractsConfig.Copy(),
	}
	if env.gasPool != nil {
		gasPool := *env.gasPool
		cpy.gasPool = &gasPool
	}
	// The content of txs are immutable, unnecessary
	// to do the expensive deep copy for them.
	cpy.txs = make([]*types.Transaction, len(env.txs))
	copy(cpy.txs, env.txs)
	return cpy
}

// unclelist returns the contained uncles as the list format.
func (env *environment) unclelist() []*types.Header {
	var uncles []*types.Header
	return uncles
}

// discard terminates the background prefetcher go-routine. It should
// always be called for all created environment instances otherwise
// the go-routine leak can happen.
func (env *environment) discard() {
	if env.state == nil {
		return
	}
	env.state.StopPrefetcher()
}

// task contains all information for consensus engine sealing and result submitting.
type task struct {
	env       *environment
	block     *types.Block
	createdAt time.Time
}

const (
	commitInterruptNone int32 = iota
	commitInterruptNewHead
	commitInterruptResubmit
)

// newWorkReq represents a request for new sealing work submitting with relative interrupt notifier.
type newWorkReq struct {
	interrupt *int32
	noempty   bool
	timestamp int64
	parent    *types.Header
}

// getWorkReq represents a request for getting a new sealing work with provided parameters.
type getWorkReq struct {
	params *generateParams
	err    error
	result chan *types.Block
}

// intervalAdjust represents a resubmitting interval adjustment.
type intervalAdjust struct {
	ratio float64
	inc   bool
}

// worker is the main object which takes care of submitting new work to consensus engine
// and gathering the sealing result.
type worker struct {
	confMu      sync.RWMutex // The lock used to protect the config fields: GasCeil, GasTip and Extradata
	config      *ethconfig.MinerConfig
	chainConfig *params.ChainConfig
	engine      consensus.Engine
	eth         Backend
	chain       *core.BlockChain

	// Feeds
	pendingLogsFeed event.Feed

	// Subscriptions
	mux          *event.TypeMux
	txsCh        chan core.NewTxsEvent
	txsSub       event.Subscription
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription

	// Channels
	newWorkCh               chan *newWorkReq
	getWorkCh               chan *getWorkReq
	proposalVerifiedEventCh chan *types.Header
	taskCh                  chan *task
	resultCh                chan *types.Block
	startCh                 chan struct{}
	exitCh                  chan struct{}
	resubmitIntervalCh      chan time.Duration
	resubmitAdjustCh        chan *intervalAdjust

	wg sync.WaitGroup

	current *environment // An environment for current running cycle.

	mu       sync.RWMutex // The lock used to protect the coinbase and extra fields
	coinbase common.Address
	extra    []byte

	pendingMu    sync.RWMutex
	pendingTasks map[common.Hash]*task

	// atomic status counters
	running int32 // The indicator whether the consensus engine is running or not.
	newTxs  int32 // New arrival transaction count since last sealing work submitting.

	// noempty is the flag used to control whether the feature of pre-seal empty
	// block is enabled. The default value is false(pre-seal is enabled by default).
	// But in some special scenario the consensus engine will seal blocks instantaneously,
	// in this case this feature will add all empty blocks into canonical chain
	// non-stop and no real transaction will be included.
	noempty uint32

	// Test hooks
	newTaskHook  func(*task)                        // Method to call upon receiving a new sealing task.
	skipSealHook func(*task) bool                   // Method to decide whether skipping the sealing.
	fullTaskHook func()                             // Method to call before pushing the full sealing task.
	resubmitHook func(time.Duration, time.Duration) // Method to call upon updating resubmitting interval.
}

func newWorker(config *ethconfig.MinerConfig, chainConfig *params.ChainConfig, engine consensus.Engine, eth Backend, mux *event.TypeMux, init bool) *worker {
	worker := &worker{
		config:                  config,
		chainConfig:             chainConfig,
		coinbase:                config.Etherbase,
		engine:                  engine,
		eth:                     eth,
		mux:                     mux,
		chain:                   eth.BlockChain(),
		pendingTasks:            make(map[common.Hash]*task),
		txsCh:                   make(chan core.NewTxsEvent, txChanSize),
		chainHeadCh:             make(chan core.ChainHeadEvent, chainHeadChanSize),
		newWorkCh:               make(chan *newWorkReq),
		getWorkCh:               make(chan *getWorkReq),
		taskCh:                  make(chan *task),
		resultCh:                make(chan *types.Block, resultQueueSize),
		proposalVerifiedEventCh: make(chan *types.Header, chainHeadChanSize),
		exitCh:                  make(chan struct{}),
		startCh:                 make(chan struct{}, 1),
		resubmitIntervalCh:      make(chan time.Duration),
		resubmitAdjustCh:        make(chan *intervalAdjust, resubmitAdjustChanSize),
	}
	// Subscribe NewTxsEvent for tx pool
	worker.txsSub = eth.TxPool().SubscribeTransactions(worker.txsCh, true)
	// Subscribe events for blockchain
	worker.chainHeadSub = eth.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)
	worker.engine.SetResultChan(worker.resultCh)
	worker.engine.SetProposalVerifiedEventChan(worker.proposalVerifiedEventCh)

	// Sanitize recommit interval if the user-specified one is too short.
	recommit := worker.config.Recommit
	if recommit < minRecommitInterval {
		eth.Logger().Warn("Sanitizing miner recommit interval", "provided", recommit, "updated", minRecommitInterval)
		recommit = minRecommitInterval
	}

	worker.wg.Add(4)
	go worker.mainLoop()
	go worker.newWorkLoop(recommit)
	go worker.resultLoop()
	go worker.taskLoop()

	// Submit first work to initialize pending state.
	if init {
		worker.startCh <- struct{}{}
	}
	return worker
}

// setExtra sets the content used to initialize the block extra field.
func (w *worker) setExtra(extra []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.extra = extra
}

// setRecommitInterval updates the interval for miner sealing work recommitting.
func (w *worker) setRecommitInterval(interval time.Duration) {
	select {
	case w.resubmitIntervalCh <- interval:
	case <-w.exitCh:
	}
}

// disablePreseal disables pre-sealing feature
func (w *worker) disablePreseal() {
	atomic.StoreUint32(&w.noempty, 1)
}

// enablePreseal enables pre-sealing feature
func (w *worker) enablePreseal() {
	atomic.StoreUint32(&w.noempty, 0)
}

// pending returns the pending state and corresponding block.
func (w *worker) pending() (*types.Header, *state.StateDB) {
	b := w.chain.CurrentBlock()
	st, _ := w.chain.StateAt(b.Root)
	return b, st
}

// pendingBlock returns pending block.
func (w *worker) pendingBlock() *types.Header {
	return w.chain.CurrentBlock()
}

// start sets the running status as 1 and triggers new work submitting.
func (w *worker) start() {
	if pos, ok := w.engine.(consensus.BFT); ok {
		err := pos.Start(context.Background())
		if err != nil && err != backend.ErrStartedEngine {
			w.eth.Logger().Error("Error starting Consensus Engine", "block", w.chain.CurrentBlock(), "error", err)
		}
	}
	atomic.StoreInt32(&w.running, 1)
	w.startCh <- struct{}{}
}

// stop sets the running status as 0.
func (w *worker) stop() {
	atomic.StoreInt32(&w.running, 0)
	if err := w.engine.Close(); err != nil {
		w.eth.Logger().Debug("Error stopping Consensus Engine", "error", err)
	}
}

// isRunning returns an indicator whether worker is running or not.
func (w *worker) isRunning() bool {
	return atomic.LoadInt32(&w.running) == 1
}

// close terminates all background threads maintained by the worker.
// Note the worker does not support being closed multiple times.
func (w *worker) close() {
	atomic.StoreInt32(&w.running, 0)
	close(w.exitCh)
	w.wg.Wait()
}

// recalcRecommit recalculates the resubmitting interval upon feedback.
func recalcRecommit(minRecommit, prev time.Duration, target float64, inc bool) time.Duration {
	var (
		prevF = float64(prev.Nanoseconds())
		next  float64
	)
	if inc {
		next = prevF*(1-intervalAdjustRatio) + intervalAdjustRatio*(target+intervalAdjustBias)
		max := float64(maxRecommitInterval.Nanoseconds())
		if next > max {
			next = max
		}
	} else {
		next = prevF*(1-intervalAdjustRatio) + intervalAdjustRatio*(target-intervalAdjustBias)
		min := float64(minRecommit.Nanoseconds())
		if next < min {
			next = min
		}
	}
	return time.Duration(int64(next))
}

// newWorkLoop is a standalone goroutine to submit new sealing work upon received events.
func (w *worker) newWorkLoop(recommit time.Duration) {
	defer w.wg.Done()
	var (
		interrupt   *int32
		minRecommit = recommit // minimal resubmit interval specified by user.
		timestamp   int64      // timestamp for each round of sealing.
	)

	timer := time.NewTimer(0)
	defer timer.Stop()
	<-timer.C // discard the initial tick

	// commit aborts in-flight transaction execution with given signal and resubmits a new one.
	commit := func(noempty bool, s int32, parent *types.Header) {
		if interrupt != nil {
			atomic.StoreInt32(interrupt, s)
		}
		interrupt = new(int32)
		select {
		case w.newWorkCh <- &newWorkReq{interrupt: interrupt, noempty: noempty, timestamp: timestamp, parent: parent}:
		case <-w.exitCh:
			return
		}
		timer.Reset(recommit)
		atomic.StoreInt32(&w.newTxs, 0)
	}
	// clearPending cleans the stale pending tasks.
	clearPending := func(number uint64) {
		w.pendingMu.Lock()
		for h, t := range w.pendingTasks {
			if t.block.NumberU64()+staleThreshold <= number {
				delete(w.pendingTasks, h)
			}
		}
		w.pendingMu.Unlock()
	}
	var lastBlock common.Hash

	for {
		select {
		case <-w.startCh:
			// do not prepare blocks if not in the committee
			if !w.isRunning() {
				continue
			}
			clearPending(w.chain.CurrentBlock().Number.Uint64())
			timestamp = time.Now().Unix()
			commit(false, commitInterruptNewHead, nil)

		case head := <-w.chainHeadCh:
			// do not prepare blocks if not in the committee
			if !w.isRunning() {
				continue
			}
			if head.Header.Hash() == lastBlock {
				log.Debug("New chain head event - block already prepared")
				if h, ok := w.engine.(consensus.Handler); ok {
					h.NewChainHead()
				}
				continue
			}
			clearPending(head.Header.Number.Uint64())
			timestamp = time.Now().Unix()
			if h, ok := w.engine.(consensus.Handler); ok {
				h.NewChainHead()
			}
			lastBlock = head.Header.Hash()
			commit(false, commitInterruptNewHead, nil)

		case block := <-w.proposalVerifiedEventCh:
			// do not prepare blocks if not in the committee
			if !w.isRunning() {
				continue
			}
			if block.Hash() == lastBlock {
				log.Debug("block already prepared")
				continue
			}
			lastBlock = block.Hash()
			clearPending(w.chain.CurrentBlock().Number.Uint64())
			timestamp = time.Now().Unix()
			commit(false, commitInterruptNewHead, block)

		case <-timer.C:
			// If sealing is running resubmit a new work cycle periodically to pull in
			// higher priced transactions. Disable this overhead for pending blocks.
			if w.isRunning() {
				// Short circuit if no new transaction arrives.
				if atomic.LoadInt32(&w.newTxs) == 0 {
					timer.Reset(recommit)
					continue
				}
				commit(true, commitInterruptResubmit, nil)
			}

		case interval := <-w.resubmitIntervalCh:
			// Adjust resubmit interval explicitly by user.
			if interval < minRecommitInterval {
				w.eth.Logger().Warn("Sanitizing miner recommit interval", "provided", interval, "updated", minRecommitInterval)
				interval = minRecommitInterval
			}
			w.eth.Logger().Info("Miner recommit interval update", "from", minRecommit, "to", interval)
			minRecommit, recommit = interval, interval

			if w.resubmitHook != nil {
				w.resubmitHook(minRecommit, recommit)
			}

		case adjust := <-w.resubmitAdjustCh:
			// Adjust resubmit interval by feedback.
			if adjust.inc {
				before := recommit
				target := float64(recommit.Nanoseconds()) / adjust.ratio
				recommit = recalcRecommit(minRecommit, recommit, target, true)
				w.eth.Logger().Trace("Increase miner recommit interval", "from", before, "to", recommit)
			} else {
				before := recommit
				recommit = recalcRecommit(minRecommit, recommit, float64(minRecommit.Nanoseconds()), false)
				w.eth.Logger().Trace("Decrease miner recommit interval", "from", before, "to", recommit)
			}

			if w.resubmitHook != nil {
				w.resubmitHook(minRecommit, recommit)
			}

		case <-w.exitCh:
			return
		}
	}
}

// mainLoop is responsible for generating and submitting sealing work based on
// the received event. It can support two modes: automatically generate task and
// submit it or return task according to given parameters for various proposes.
func (w *worker) mainLoop() {
	defer w.wg.Done()
	defer w.txsSub.Unsubscribe()
	defer w.chainHeadSub.Unsubscribe()
	defer func() {
		if w.current != nil {
			w.current.discard()
		}
	}()

	for {
		select {
		case req := <-w.newWorkCh:
			w.commitWork(req)

		case req := <-w.getWorkCh:
			block, _, err := w.generateWork(req.params)
			if err != nil {
				req.err = err
				req.result <- nil
			} else {
				req.result <- block
			}

		case ev := <-w.txsCh:
			atomic.AddInt32(&w.newTxs, int32(len(ev.Txs)))
		// System stopped
		case <-w.exitCh:
			return
		case <-w.txsSub.Err():
			return
		case <-w.chainHeadSub.Err():
			return
		}
	}
}

// taskLoop is a standalone goroutine to fetch sealing task from the generator and
// push them to consensus engine.
func (w *worker) taskLoop() {
	defer w.wg.Done()
	var (
		stopCh chan struct{}
		prev   common.Hash
	)

	// interrupt aborts the in-flight sealing task.
	interrupt := func() {
		if stopCh != nil {
			close(stopCh)
		}
	}
	for {
		select {
		case task := <-w.taskCh:
			if w.newTaskHook != nil {
				w.newTaskHook(task)
			}
			// Reject duplicate sealing work due to resubmitting.
			sealHash := w.engine.SealHash(task.block.Header())
			if sealHash == prev {
				continue
			}

			if w.skipSealHook != nil && w.skipSealHook(task) {
				continue
			}

			// Interrupt previous sealing operation
			interrupt()
			stopCh, prev = make(chan struct{}), sealHash

			w.eth.Logger().Debug("New block Seal request", "hash", sealHash)
			w.pendingMu.Lock()
			w.pendingTasks[sealHash] = task
			w.pendingMu.Unlock()

			// go routine is needed to be able to interrupt the task, in case a more recent block arrives
			go func(stopCh chan struct{}) { // pass stopCh(creates a copy of reference) to avoid write-read race on stopCh
				sealStart := time.Now()
				if err := w.engine.Seal(task.env.parentHeader, task.block, w.resultCh, stopCh); err != nil {
					w.eth.Logger().Warn("Block sealing failed", "err", err)
					w.pendingMu.Lock()
					delete(w.pendingTasks, sealHash)
					w.pendingMu.Unlock()
					return
				}

				if metrics.Enabled() {
					now := time.Now()
					SealWorkTimer.Update(now.Sub(sealStart))
					SealWorkBg.Add(now.Sub(sealStart).Nanoseconds())
				}
			}(stopCh)
		case <-w.exitCh:
			interrupt()
			return
		}
	}
}

// resultLoop is a standalone goroutine to handle sealing result submitting
// and flush relative data to the database.
func (w *worker) resultLoop() {
	defer w.wg.Done()
	for {
		select {
		case block := <-w.resultCh:
			// Short circuit when receiving empty result.
			if block == nil {
				continue
			}
			// Short circuit when receiving duplicate result caused by resubmitting.
			if w.chain.HasBlock(block.Hash(), block.NumberU64()) {
				continue
			}
			var (
				sealhash = w.engine.SealHash(block.Header())
				hash     = block.Hash()
			)
			w.pendingMu.RLock()
			task, exist := w.pendingTasks[sealhash]
			w.pendingMu.RUnlock()
			if !exist {
				w.eth.Logger().Error("Block found but no relative pending task", "number", block.Number(), "sealhash", sealhash, "hash", hash)
				continue
			}

			// The logs inside the task.env.state may contain parent block's logs as we copied the parent state when
			// making environment for the child of an optimistic parent block for consensus pipeline optimization.
			// Thus, we just deep copy logs from the receipts of the current block execution environment.
			var logs []*types.Log

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for i, r := range task.env.receipts {
				r.BlockHash = block.Hash()
				r.BlockNumber = block.Number()
				r.TransactionIndex = uint(i)
				for _, l := range r.Logs {
					l.BlockHash = block.Hash()
				}
				logs = append(logs, r.Logs...)
			}

			// Commit block and state to database.
			persistStart := time.Now()

			// As the environment is no longer used anymore, we won't copy the receipts to announce the new block.
			_, err := w.chain.WriteBlockAndSetHead(block, task.env.receipts, logs, task.env.state, true, task.env.contractsConfig)
			if err != nil {
				w.eth.Logger().Error("Failed writing block to chain", "err", err)
				continue
			}
			if metrics.Enabled() {
				now := time.Now()
				PersistWorkTimer.Update(now.Sub(persistStart))
				PersistWorkBg.Add(now.Sub(persistStart).Nanoseconds())
				TotalTaskProcessBg.Add(now.Sub(task.createdAt).Nanoseconds())
			}
			w.eth.Logger().Info("🔨 Proposed block validated with success", "number", block.Number(), "sealhash", sealhash, "hash", hash,
				"elapsed", common.PrettyDuration(time.Since(task.createdAt)))

			// Broadcast the block and announce chain insertion event
			w.mux.Post(core.NewMinedBlockEvent{Block: block})

		case <-w.exitCh:
			return
		}
	}
}

func (w *worker) optimisticStateAndConfig(parent *types.Header) (*state.StateDB, *types.ContractsConfig, error) {
	var state *state.StateDB
	var contractsConfig *types.ContractsConfig
	var hash common.Hash
	// Making environment by coping cached optimistic parent block's state also copies the receipt logs.
	if parent.Coinbase == w.coinbase { // we were the proposer for the parent
		sealHash := w.engine.SealHash(parent)
		w.pendingMu.Lock()
		task, exist := w.pendingTasks[sealHash]
		if exist {
			state = task.env.state.Copy()
			contractsConfig = task.env.contractsConfig.Copy()
		} else {
			w.pendingMu.Unlock()
			return nil, nil, fmt.Errorf("no state cache available for optimistic block")
		}
		w.pendingMu.Unlock()
	} else {
		state, contractsConfig, hash = w.chain.LoadProposalState()
		if parent.Hash() != hash {
			return nil, nil, fmt.Errorf("no state cache available for optimistic block")
		}
	}
	return state, contractsConfig, nil

}

// makeEnv creates a new environment for the sealing block.
func (w *worker) makeEnv(parent *types.Header, header *types.Header, coinbase common.Address, state *state.StateDB) *environment {
	// Note the passed coinbase may be different with header.Coinbase.
	env := &environment{
		signer:       types.MakeSigner(w.chainConfig, header.Number),
		state:        state,
		coinbase:     coinbase,
		header:       header,
		parentHeader: parent,
		evm:          vm.NewEVM(core.NewEVMBlockContext(header, w.chain, &w.coinbase), state, w.chain.Config(), *w.chain.GetVMConfig()),
	}
	// Keep track of transactions which return errors so they can be removed
	env.tcount = 0
	return env
}

func (w *worker) commitTransaction(env *environment, tx *types.Transaction) ([]*types.Log, error) {
	snap := env.state.Snapshot()

	receipt, err := core.ApplyTransaction(env.evm, env.gasPool, env.state, env.header, tx, &env.header.GasUsed)
	if err != nil {
		env.state.RevertToSnapshot(snap)
		return nil, err
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return receipt.Logs, nil
}

func (w *worker) commitTransactions(env *environment, txs *transactionsByPriceAndNonce, interrupt *int32) bool {
	gasLimit := env.header.GasLimit
	if env.gasPool == nil {
		env.gasPool = new(core.GasPool).AddGas(gasLimit)
	}
	var coalescedLogs []*types.Log

	for {
		// In the following three cases, we will interrupt the execution of the transaction.
		// (1) new head block event arrival, the interrupt signal is 1
		// (2) worker start or restart, the interrupt signal is 1
		// (3) worker recreate the sealing block with any newly arrived transactions, the interrupt signal is 2.
		// For the first two cases, the semi-finished work will be discarded.
		// For the third case, the semi-finished work will be submitted to the consensus engine.
		if interrupt != nil && atomic.LoadInt32(interrupt) != commitInterruptNone {
			// Notify resubmit loop to increase resubmitting interval due to too frequent commits.
			if atomic.LoadInt32(interrupt) == commitInterruptResubmit {
				ratio := float64(gasLimit-env.gasPool.Gas()) / float64(gasLimit)
				if ratio < 0.1 {
					ratio = 0.1
				}
				w.resubmitAdjustCh <- &intervalAdjust{
					ratio: ratio,
					inc:   true,
				}
			}
			return atomic.LoadInt32(interrupt) == commitInterruptNewHead
		}
		// If we don't have enough gas for any further transactions then we're done
		if env.gasPool.Gas() < params.TxGas {
			w.eth.Logger().Trace("Not enough gas for further transactions", "have", env.gasPool, "want", params.TxGas)
			break
		}
		// Retrieve the next transaction and abort if all done
		ltx, _ := txs.Peek()
		if ltx == nil {
			break
		}
		// If we don't have enough space for the next transaction, skip the account.
		if env.gasPool.Gas() < ltx.Gas {
			log.Trace("Not enough gas left for transaction", "hash", ltx.Hash, "left", env.gasPool.Gas(), "needed", ltx.Gas)
			txs.Pop()
			continue
		}
		// Transaction seems to fit, pull it up from the pool
		tx := ltx.Resolve()
		if tx == nil {
			log.Trace("Ignoring evicted transaction", "hash", ltx.Hash)
			txs.Pop()
			continue
		}

		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(env.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !w.chainConfig.IsEIP155(env.header.Number) {
			w.eth.Logger().Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", w.chainConfig.EIP155Block)
			txs.Pop()
			continue
		}
		// Start executing the transaction
		env.state.SetTxContext(tx.Hash(), env.tcount)

		logs, err := w.commitTransaction(env, tx)
		switch {
		case errors.Is(err, core.ErrGasLimitReached):
			// Pop the current out-of-gas transaction without shifting in the next from the account
			w.eth.Logger().Trace("Gas limit exceeded for current block", "sender", from)
			txs.Pop()

		case errors.Is(err, core.ErrNonceTooLow):
			// New head notification data race between the transaction pool and miner, shift
			w.eth.Logger().Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce(), "hash", tx.Hash().Hex())
			txs.Shift()

		case errors.Is(err, core.ErrNonceTooHigh):
			// Reorg notification data race between the transaction pool and miner, skip account =
			w.eth.Logger().Trace("Skipping account with hight nonce", "sender", from, "nonce", tx.Nonce())
			txs.Pop()

		case errors.Is(err, nil):
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			env.tcount++
			txs.Shift()

		case errors.Is(err, core.ErrTxTypeNotSupported):
			// Pop the unsupported transaction without shifting in the next from the account
			w.eth.Logger().Trace("Skipping unsupported transaction type", "sender", from, "type", tx.Type())
			txs.Pop()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			w.eth.Logger().Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}

	if !w.isRunning() && len(coalescedLogs) > 0 {
		// We don't push the pendingLogsEvent while we are sealing. The reason is that
		// when we are sealing, the worker will regenerate a sealing block every 3 seconds.
		// In order to avoid pushing the repeated pendingLog, we disable the pending log pushing.

		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		w.pendingLogsFeed.Send(cpy)
	}
	// Notify resubmit loop to decrease resubmitting interval if current interval is larger
	// than the user-specified one.
	if interrupt != nil {
		w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	}
	return false
}

// generateParams wraps various of settings for generating sealing task.
type generateParams struct {
	timestamp  uint64         // The timstamp for sealing task
	forceTime  bool           // Flag whether the given timestamp is immutable or not
	parentHash common.Hash    // Parent block hash, empty means the latest chain head
	coinbase   common.Address // The fee recipient address for including transaction
	random     common.Hash    // The randomness generated by beacon chain, empty before the merge
	noUncle    bool           // Flag whether the uncle block inclusion is allowed
	noExtra    bool           // Flag whether the extra field assignment is allowed
}

// prepareWork constructs the sealing task according to the given parameters,
// either based on the last chain head or specified parent. In this function
// the pending transactions are not filled yet, only the empty task returned.
func (w *worker) prepareWork(genParams *generateParams, parent *types.Header) (*environment, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	optimisticCandidate := false

	// Find the parent block for sealing task
	if parent == nil {
		parent = w.chain.CurrentBlock()
	} else {
		optimisticCandidate = true
	}

	if genParams.parentHash != (common.Hash{}) {
		parent = w.chain.GetHeaderByHash(genParams.parentHash)
		if parent == nil {
			return nil, fmt.Errorf("missing parent")
		}
	}

	if parent == nil {
		return nil, fmt.Errorf("missing parent")
	}

	// Sanity check the timesw.chainConfig.Tendermint.BlockPeriod timestamp correctness, recap the timestamp
	// to parent+1 if the mutation is allowed.
	timestamp := genParams.timestamp
	if parent.Time >= timestamp {
		if genParams.forceTime {
			return nil, fmt.Errorf("invalid timestamp, parent %d given %d", parent.Time, timestamp)
		}
		timestamp = parent.Time + 1
	}

	// compute current block number
	num := new(big.Int).Add(parent.Number, common.Big1)

	// based on whether the block we are building is optimistic or not:
	// 1. fetch the correct state
	// 2. fetch the correct gas limit
	// 3. fetch the correct eip1559 parameters
	var (
		state            *state.StateDB
		protocolGasLimit uint64
		eip1559Params    *types.Eip1559Params
	)
	if !optimisticCandidate {
		var err error
		state, err = w.chain.StateAt(parent.Root)
		if err != nil {
			w.eth.Logger().Error("Failed to get parent state", "root", parent.Root, "err", err)
			return nil, err
		}
		protocolGasLimitBig, err := w.chain.GasLimitByHeight(num.Uint64())
		if err != nil {
			w.eth.Logger().Error("Failed to get gas limit", "num", num.Uint64(), "err", err)
			return nil, fmt.Errorf("failed to get gas limit: %w", err)
		}
		protocolGasLimit = protocolGasLimitBig.Uint64()
		eip1559Params, err = w.chain.Eip1559ParamsByHeight(num.Uint64())
		if err != nil {
			w.eth.Logger().Error("Failed to get gas limit bound divisor", "num", num.Uint64(), "err", err)
			return nil, fmt.Errorf("failed to get gas limit bound divisor: %w", err)
		}
	} else {
		var err error
		var contractsConfig *types.ContractsConfig
		state, contractsConfig, err = w.optimisticStateAndConfig(parent)
		if err != nil {
			w.eth.Logger().Error("Failed to get optimistic state", "err", err)
			return nil, fmt.Errorf("error while fetching optimistic state: %w", err)
		}
		protocolGasLimit = contractsConfig.GasLimit.Uint64()
		eip1559Params = &(contractsConfig.Eip1559)
	}

	// Construct the sealing block header, set the extra field if it's allowed
	header := &types.Header{
		ParentHash: parent.Hash(),
		Difficulty: common.Big0,
		Number:     new(big.Int).Add(parent.Number, common.Big1),
		GasLimit:   core.CalcGasLimit(parent.GasLimit, protocolGasLimit, eip1559Params.GasLimitBoundDivisor.Uint64()),
		Time:       timestamp,
		Coinbase:   genParams.coinbase,
	}
	if !genParams.noExtra && len(w.extra) != 0 {
		header.Extra = w.extra
	}
	// Set the randomness field from the beacon chain if it's available.
	if genParams.random != (common.Hash{}) {
		header.MixDigest = genParams.random
	}
	// Set baseFee and GasLimit if we are on an EIP-1559 chain
	if w.chainConfig.IsLondon(header.Number) {
		header.BaseFee = misc.CalcBaseFee(w.chainConfig, parent, eip1559Params)
		if !w.chainConfig.IsLondon(parent.Number) {
			// take the genesis elasticity multiplier, however this
			// should never happen on Autonity, since EIP1559 is activated from genesis
			parentGasLimit := parent.GasLimit * params.DefaultElasticityMultiplier
			header.GasLimit = core.CalcGasLimit(parentGasLimit, protocolGasLimit, eip1559Params.GasLimitBoundDivisor.Uint64())
		}
	}

	// Could potentially happen if starting to mine in an odd state.
	// Note genParams.coinbase can be different with header.Coinbase
	// since clique algorithm can modify the coinbase field in header.
	env := w.makeEnv(parent, header, genParams.coinbase, state)

	// Run the consensus preparation with the default or customized consensus engine.
	if err := w.engine.Prepare(w.chain, env.parentHeader, env.header, env.state); err != nil {
		w.eth.Logger().Error("Failed to prepare header for sealing", "err", err)
		return nil, err
	}

	// start prefetcher here to avoid go routine leak
	if !optimisticCandidate {
		env.state.StartPrefetcher("miner", nil)
	}

	return env, nil
}

// fillTransactions retrieves the pending transactions from the txpool and fills them
// into the given sealing block. The transaction selection and ordering strategy can
// be customized with the plugin in the future.
func (w *worker) fillTransactions(interrupt *int32, env *environment) {
	w.confMu.RLock()
	tip := w.config.GasPrice
	prio := w.config.Prio
	w.confMu.RUnlock()
	// Retrieve the pending transactions pre-filtered by the 1559/4844 dynamic fees
	filter := txpool.PendingFilter{
		MinTip:       uint256.MustFromBig(tip),
		OnlyPlainTxs: true,
	}
	if env.header.BaseFee != nil {
		filter.BaseFee = uint256.MustFromBig(env.header.BaseFee)
	}

	pendingPlainTxs := w.eth.TxPool().Pending(filter)

	// Split the pending transactions into locals and remotes.
	prioPlainTxs, normalPlainTxs := make(map[common.Address][]*txpool.LazyTransaction), pendingPlainTxs

	for _, account := range prio {
		if txs := normalPlainTxs[account]; len(txs) > 0 {
			delete(normalPlainTxs, account)
			prioPlainTxs[account] = txs
		}
	}
	// Fill the block with all available pending transactions.
	if len(prioPlainTxs) > 0 {
		plainTxs := newTransactionsByPriceAndNonce(env.signer, prioPlainTxs, env.header.BaseFee)
		if w.commitTransactions(env, plainTxs, interrupt) {
			return
		}
	}
	if len(normalPlainTxs) > 0 {
		plainTxs := newTransactionsByPriceAndNonce(env.signer, normalPlainTxs, env.header.BaseFee)
		if w.commitTransactions(env, plainTxs, interrupt) {
			return
		}
	}
	return
}

// generateWork generates a sealing block based on the given parameters.
func (w *worker) generateWork(params *generateParams) (*types.Block, *types.ContractsConfig, error) {
	work, err := w.prepareWork(params, nil)
	if err != nil {
		return nil, nil, err
	}
	defer work.discard()
	w.fillTransactions(nil, work)
	body := &types.Body{
		Transactions: work.txs,
		Uncles:       work.unclelist(),
	}
	return w.engine.FinalizeAndAssemble(w.chain, work.header, work.state, body, &work.receipts)
}

// commitWork generates several new sealing tasks based on the parent block
// and submit them to the sealer.
func (w *worker) commitWork(req *newWorkReq) {
	start := time.Now()

	// Set the coinbase if the worker is running or it's required
	var coinbase common.Address
	if w.isRunning() {
		if w.coinbase == (common.Address{}) {
			w.eth.Logger().Error("Refusing to mine without etherbase")
			return
		}
		coinbase = w.coinbase // Use the preset address as the fee recipient
	}

	//TODO: create block with hash
	work, err := w.prepareWork(&generateParams{
		timestamp: uint64(req.timestamp),
		coinbase:  coinbase,
	}, req.parent)
	if err != nil {
		return
	}
	if metrics.Enabled() {
		now := time.Now()
		PrepareWorkTimer.Update(now.Sub(start))
		PrepareWorkBg.Add(now.Sub(start).Nanoseconds())
	}

	// Create an empty block based on temporary copied state for
	// sealing in advance without waiting block execution finished.
	if !req.noempty && w.chainConfig.Ethash != nil && atomic.LoadUint32(&w.noempty) == 0 {
		// Use a copy of the environment to avoid double-applying block rewards
		// when the full block commit happens later
		emptyEnv := &environment{
			signer:       work.signer,
			state:        work.state.Copy(),
			coinbase:     work.coinbase,
			header:       types.CopyHeader(work.header),
			parentHeader: work.parentHeader,
			evm:          vm.NewEVM(core.NewEVMBlockContext(work.header, w.chain, &w.coinbase), work.state.Copy(), w.chain.Config(), *w.chain.GetVMConfig()),
		}
		w.commit(emptyEnv, nil, false, start)
	}

	fillTxStart := time.Now()
	// Fill pending transactions from the txpool
	w.fillTransactions(req.interrupt, work)
	if metrics.Enabled() {
		now := time.Now()
		FillWorkTimer.Update(now.Sub(fillTxStart))
		FillWorkBg.Add(now.Sub(fillTxStart).Nanoseconds())
	}

	commitWorkStart := time.Now()
	w.commit(work, w.fullTaskHook, true, start)
	if metrics.Enabled() {
		now := time.Now()
		CommitWorkTimer.Update(now.Sub(commitWorkStart))
		CommitWorkBg.Add(now.Sub(commitWorkStart).Nanoseconds())
	}

	// Swap out the old work with the new one, terminating any leftover
	// prefetcher processes in the mean time and starting a new one.
	if w.current != nil {
		w.current.discard()
	}
	w.current = work
}

// commit runs any post-transaction state modifications, assembles the final block
// and commits new work if consensus engine is running.
// Note the assumption is held that the mutation is allowed to the passed env, do
// the deep copy first.
func (w *worker) commit(env *environment, interval func(), update bool, start time.Time) error {
	if w.isRunning() {
		if interval != nil {
			interval()
		}
		finalizeStart := time.Now()
		body := &types.Body{
			Transactions: env.txs,
			Uncles:       env.unclelist(),
		}
		block, contractsConfig, err := w.engine.FinalizeAndAssemble(w.chain, env.header, env.state, body, &env.receipts)
		if err != nil {
			return err
		}
		if metrics.Enabled() {
			now := time.Now()
			FinalizeWorkTimer.Update(now.Sub(finalizeStart))
			FinalizeWorkBg.Add(now.Sub(finalizeStart).Nanoseconds())
		}

		// store contracts config in the mining environment
		env.contractsConfig = contractsConfig

		select {
		case w.taskCh <- &task{env: env, block: block, createdAt: time.Now()}:
			if metrics.Enabled() {
				TotalTaskPrepareBg.Add(time.Since(start).Nanoseconds())
			}
			w.eth.Logger().Info("Preparing new block proposal", "number", block.Number(), "sealhash", w.engine.SealHash(block.Header()),
				"txs", env.tcount,
				"gas", block.GasUsed(), "fees", totalFees(block, env.receipts),
				"elapsed", common.PrettyDuration(time.Since(start))) // Consider moving that to DEBUG level

		case <-w.exitCh:
			w.eth.Logger().Info("Worker has exited")
		}
	}
	return nil
}

// getSealingBlock generates the sealing block based on the given parameters.
func (w *worker) getSealingBlock(parent common.Hash, timestamp uint64, coinbase common.Address, random common.Hash) (*types.Block, error) {
	req := &getWorkReq{
		params: &generateParams{
			timestamp:  timestamp,
			forceTime:  true,
			parentHash: parent,
			coinbase:   coinbase,
			random:     random,
			noUncle:    true,
			noExtra:    true,
		},
		result: make(chan *types.Block, 1),
	}
	select {
	case w.getWorkCh <- req:
		block := <-req.result
		if block == nil {
			return nil, req.err
		}
		return block, nil
	case <-w.exitCh:
		return nil, errors.New("miner closed")
	}
}

// copyReceipts makes a deep copy of the given receipts.
func copyReceipts(receipts []*types.Receipt) []*types.Receipt {
	result := make([]*types.Receipt, len(receipts))
	for i, l := range receipts {
		cpy := *l
		result[i] = &cpy
	}
	return result
}

// postSideBlock fires a side chain event, only use it for testing.
func (w *worker) postSideBlock(_ core.ChainSideEvent) {}

// totalFees computes total consumed miner fees in ETH. Block transactions and receipts have to have the same order.
func totalFees(block *types.Block, receipts []*types.Receipt) *big.Float {
	feesWei := new(big.Int)
	for i, tx := range block.Transactions() {
		minerFee, _ := tx.EffectiveGasTip(block.BaseFee())
		feesWei.Add(feesWei, new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), minerFee))
	}
	return new(big.Float).Quo(new(big.Float).SetInt(feesWei), new(big.Float).SetInt(big.NewInt(params.Ether)))
}
