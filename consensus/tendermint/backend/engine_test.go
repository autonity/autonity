package backend

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
)

func TestPrepare(t *testing.T) {
	chain, engine := newBlockChain(1)
	header := makeHeader(chain.Genesis(), chain)
	err := engine.Prepare(chain, header)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	header.ParentHash = common.BytesToHash([]byte("1234567890"))
	err = engine.Prepare(chain, header)
	if err != consensus.ErrUnknownAncestor {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrUnknownAncestor)
	}
}

func TestSealCommittedOtherHash(t *testing.T) {
	chain, engine := newBlockChain(4)

	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	otherBlock, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	eventSub := engine.Subscribe(events.CommitEvent{})
	eventLoop := func() {
		ev := <-eventSub.Chan()
		_, ok := ev.Data.(events.CommitEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		err = engine.Commit(otherBlock, 0, make(types.Signatures))
		if err != nil {
			t.Error("commit should not return error", err.Error())
		}

		eventSub.Unsubscribe()
	}
	go eventLoop()
	seal := func() {
		resultCh := make(chan *types.Block)
		engine.SetResultChan(resultCh)
		err = engine.Seal(chain, block, resultCh, nil)
		if err != nil {
			t.Error("seal should not return error", err.Error())
		}

		<-resultCh
		t.Error("seal should not be completed")
	}
	go seal()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	<-timeout.C
	// wait 2 seconds to ensure we cannot get any blocks from Istanbul
}

func TestSealCommitted(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	expectedBlock, _ := engine.AddSeal(block)

	resultCh := make(chan *types.Block)
	engine.SetResultChan(resultCh)
	err = engine.Seal(chain, block, resultCh, nil)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	finalBlock := <-resultCh
	if finalBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
	}
}

func TestVerifyHeader(t *testing.T) {
	chain, engine := newBlockChain(1)

	// errEmptyCommittedSeals case
	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	block, _ = engine.AddSeal(block)
	err = engine.VerifyHeader(chain, block.Header(), false)
	if err != types.ErrEmptyCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrEmptyCommittedSeals)
	}

	header := block.Header()

	// non zero MixDigest
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.MixDigest = common.BytesToHash([]byte("123456789"))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidMixDigest {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMixDigest)
	}

	// invalid uncles hash
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.UncleHash = common.BytesToHash([]byte("123456789"))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidUncleHash {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidUncleHash)
	}

	// invalid difficulty
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidDifficulty {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidDifficulty)
	}

	// invalid timestamp
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.Time = 0
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	header.Time = new(big.Int).Add(big.NewInt(now().Unix()), new(big.Int).SetUint64(10)).Uint64()
	err = engine.VerifyHeader(chain, header, false)
	if err != consensus.ErrFutureTimestampBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureTimestampBlock)
	}

	// invalid nonce
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	header = block.Header()
	copy(header.Nonce[:], hexutil.MustDecode("0x111111111111"))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidNonce {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidNonce)
	}
}

/* The logic of this needs to change with respect of Autonity contact */
func TestVerifyHeaders(t *testing.T) {
	chain, engine := newBlockChain(1)

	// success case
	var headers []*types.Header
	var blocks []*types.Block
	size := 100

	var err error
	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
		} else {
			b, err = makeBlockWithoutSeal(chain, engine, blocks[i-1])
		}
		if err != nil {
			t.Fatal(err)
		}

		b, _ = engine.AddSeal(b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}

	_, results := engine.VerifyHeaders(chain, headers, nil)

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			if err != nil {
				/*  The two following errors mean that the processing has gone right */
				if !errors.Is(err, types.ErrEmptyCommittedSeals) && !errors.Is(err, types.ErrInvalidCommittedSeals) {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT1
				}
			}
			index++
			if index == size {
				break OUT1
			}
		case <-timeout.C:
			break OUT1
		}
	}
}

/* The logic of this needs to change with respect of Autonity contact */
func TestVerifyHeadersAbortValidation(t *testing.T) {
	chain, engine := newBlockChain(1)

	// success case
	var headers []*types.Header
	var blocks []*types.Block
	size := 100

	var err error
	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
		} else {
			b, err = makeBlockWithoutSeal(chain, engine, blocks[i-1])
		}
		if err != nil {
			t.Fatal(err)
		}

		b, _ = engine.AddSeal(b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}

	const timeoutDura = 2 * time.Second

	// abort cases
	abort, results := engine.VerifyHeaders(chain, headers, nil)
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT2:
	for {
		select {
		case err := <-results:
			if err != nil {
				if !errors.Is(err, types.ErrEmptyCommittedSeals) && !errors.Is(err, types.ErrInvalidCommittedSeals) {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT2
				}
			}
			index++
			if index == 5 {
				abort <- struct{}{}
			}
			if index >= size {
				t.Errorf("verifyheaders should be aborted")
				break OUT2
			}
		case <-timeout.C:
			break OUT2
		}
	}
}

/* The logic of this needs to change with respect of Autonity contact */
func TestVerifyErrorHeaders(t *testing.T) {
	chain, engine := newBlockChain(1)

	// success case
	var headers []*types.Header
	var blocks []*types.Block
	size := 100

	var err error
	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
		} else {
			b, err = makeBlockWithoutSeal(chain, engine, blocks[i-1])
		}
		if err != nil {
			t.Fatal(err)
		}

		b, _ = engine.AddSeal(b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}

	const timeoutDura = 2 * time.Second

	// error header cases
	headers[2].Number = big.NewInt(100)
	_, results := engine.VerifyHeaders(chain, headers, nil)
	timeout := time.NewTimer(timeoutDura)
	index := 0
	errorCount := 0
	expectedErrors := 2

OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				if !errors.Is(err, types.ErrEmptyCommittedSeals) && !errors.Is(err, types.ErrInvalidCommittedSeals) {
					errorCount++
				}
			}
			index++
			if index == size {
				if errorCount != expectedErrors {
					t.Errorf("error mismatch: have %v, want %v", err, expectedErrors)
				}
				break OUT3
			}
		case <-timeout.C:
			break OUT3
		}
	}
}

func TestWriteCommittedSeals(t *testing.T) {
	expectedCommittedSeal := make(types.Signatures)
	expectedCommittedSeal[testAddress] = testSignature.(*blst.BlsSignature)

	h := &types.Header{}

	// normal case
	err := types.WriteCommittedSeals(h, expectedCommittedSeal)
	if err != nil {
		t.Errorf("error mismatch: have %v, want %v", err, nil)
	}

	if !reflect.DeepEqual(h.CommittedSeals, expectedCommittedSeal) {
		t.Errorf("extra data mismatch: have %v, want %v", h.CommittedSeals, expectedCommittedSeal)
	}

	// invalid seal
	err = types.WriteCommittedSeals(h, nil)
	if err != types.ErrInvalidCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrInvalidCommittedSeals)
	}
}

func TestAPIs(t *testing.T) {
	b := &Backend{}

	APIS := b.APIs(nil)
	if len(APIS) < 1 {
		t.Fatalf("expected non empty slice")
	}

	if APIS[0].Namespace != "tendermint" {
		t.Fatalf("expected 'tendermint', got %v", APIS[0].Namespace)
	}
}

func TestClose(t *testing.T) {
	t.Run("engine is not running, error returned", func(t *testing.T) {
		b := &Backend{}

		err := b.Close()
		assertError(t, ErrStoppedEngine, err)
		assertNotCoreStarted(t, b)
	})

	t.Run("engine is running, no errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Stop().MaxTimes(1)

		b := &Backend{
			core:    tendermintC,
			stopped: make(chan struct{}),
		}
		b.coreStarting.Store(true)
		b.coreRunning.Store(true)

		err := b.Close()
		assertNilError(t, err)
		assertNotCoreStarted(t, b)
	})

	t.Run("engine is running, stopped twice", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Stop().MaxTimes(1)

		b := &Backend{
			core:    tendermintC,
			stopped: make(chan struct{}),
		}
		b.coreStarting.Store(true)
		b.coreRunning.Store(true)

		err := b.Close()
		assertNilError(t, err)
		assertNotCoreStarted(t, b)

		err = b.Close()
		assertError(t, ErrStoppedEngine, err)
		assertNotCoreStarted(t, b)
	})

	t.Run("engine is running, stopped from multiple goroutines", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Stop().MaxTimes(1)

		b := &Backend{
			core:    tendermintC,
			stopped: make(chan struct{}),
		}
		b.coreStarting.Store(true)
		b.coreRunning.Store(true)

		var wg sync.WaitGroup
		stop := 10
		errC := make(chan error, stop)

		for i := 0; i < stop; i++ {
			wg.Add(1)

			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				errC <- b.Close()
			}(&wg)

		}

		wg.Wait()
		close(errC)

		assertNotCoreStarted(t, b)

		var sawNil bool

		// Want nil once and ErrStoppedEngine 9 times
		for e := range errC {
			if e == nil {
				if sawNil {
					t.Fatalf("<nil> returned more than once, b.Close() should have only returned nil the first time it was closed")
				}
				sawNil = true
			} else if e != ErrStoppedEngine {
				assertError(t, ErrStoppedEngine, e)
			}
		}
	})
}

func TestStart(t *testing.T) {
	t.Run("engine is not running, no errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chain, _ := newBlockChain(1)
		ctx := context.Background()
		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Start(gomock.Any(), gomock.Any()).MaxTimes(1)
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:       tendermintC,
			gossiper:   g,
			blockchain: chain,
		}

		err := b.Start(ctx)
		assertNilError(t, err)
		assertCoreStarted(t, b)
	})

	t.Run("engine is running, error returned", func(t *testing.T) {
		b := &Backend{}
		b.coreStarting.Store(true)
		b.coreRunning.Store(true)

		err := b.Start(context.Background())
		assertError(t, ErrStartedEngine, err)
		assertCoreStarted(t, b)
	})

	t.Run("engine is not running, started twice", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Start(gomock.Any(), gomock.Any()).MaxTimes(1)
		chain, _ := newBlockChain(1)
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:       tendermintC,
			gossiper:   g,
			blockchain: chain,
		}
		b.coreStarting.Store(false)

		err := b.Start(ctx)
		assertNilError(t, err)
		assertCoreStarted(t, b)

		err = b.Start(ctx)
		assertError(t, ErrStartedEngine, err)
		assertCoreStarted(t, b)
	})

	t.Run("engine is not running, started from multiple goroutines", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		chain, _ := newBlockChain(1)
		ctx := context.Background()
		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Start(gomock.Any(), gomock.Any()).AnyTimes()
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:       tendermintC,
			gossiper:   g,
			blockchain: chain,
		}
		b.coreStarting.Store(false)

		var wg sync.WaitGroup
		stop := 10
		errC := make(chan error, stop)

		for i := 0; i < stop; i++ {
			wg.Add(1)

			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				errC <- b.Start(ctx)
			}(&wg)

		}

		wg.Wait()
		close(errC)

		assertCoreStarted(t, b)

		var sawNil bool

		// Want nil once and ErrStartedEngine 9 times
		for e := range errC {
			if e == nil {
				if sawNil {
					t.Fatalf("<nil> returned more than once, b.Start() should have only returned nil the first time it was started")
				}
				sawNil = true
			} else if e != ErrStartedEngine {
				assertError(t, ErrStartedEngine, e)
			}
		}
	})
}

func TestMultipleRestart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	times := 5
	ctx := context.Background()
	tendermintC := interfaces.NewMockCore(ctrl)
	tendermintC.EXPECT().Start(gomock.Any(), gomock.Any()).MaxTimes(times)
	tendermintC.EXPECT().Stop().MaxTimes(5)
	chain, _ := newBlockChain(1)
	g := interfaces.NewMockGossiper(ctrl)
	g.EXPECT().UpdateStopChannel(gomock.Any()).MaxTimes(5)

	b := &Backend{
		core:       tendermintC,
		gossiper:   g,
		blockchain: chain,
	}
	b.coreStarting.Store(false)

	for i := 0; i < times; i++ {
		err := b.Start(ctx)
		assertNilError(t, err)
		assertCoreStarted(t, b)

		err = b.Close()
		assertNilError(t, err)
		assertNotCoreStarted(t, b)
	}
}

func assertError(t *testing.T, expected, got error) {
	t.Helper()
	if expected != got {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func assertNilError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected <nil>, got %v", err)
	}
}

func assertCoreStarted(t *testing.T, b *Backend) {
	t.Helper()
	if !b.coreRunning.Load() {
		t.Fatal("expected core to have started")
	}
}

func assertNotCoreStarted(t *testing.T, b *Backend) {
	t.Helper()
	if b.coreRunning.Load() {
		t.Fatal("expected core to have stopped")
	}
}

func TestBackendSealHash(t *testing.T) {
	b := &Backend{}

	res := b.SealHash(&types.Header{})
	if res.Hex() == "" {
		t.Fatalf("expected not empty string")
	}
}
