package backend

import (
	"context"
	"errors"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

func TestPrepare(t *testing.T) {
	chain, engine := newBlockChain(1)

	header := makeHeader(chain.Genesis(), chain)
	err := engine.Prepare(chain, header)
	require.NoError(t, err)

	header.ParentHash = common.BytesToHash([]byte("1234567890"))
	err = engine.Prepare(chain, header)
	require.True(t, errors.Is(err, consensus.ErrUnknownAncestor))
}

func TestSealCommitted(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	expectedBlock, err := engine.AddSeal(block)
	require.NoError(t, err)

	resultCh := make(chan *types.Block)
	engine.SetResultChan(resultCh)
	err = engine.Seal(chain, block, resultCh, nil)
	require.NoError(t, err)

	finalBlock := <-resultCh
	require.Equal(t, expectedBlock.Hash(), finalBlock.Hash())
}

func TestVerifyHeader(t *testing.T) {
	chain, engine := newBlockChain(1)

	// errEmptyQuorumCertificate case
	block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	block, err = engine.AddSeal(block)
	require.NoError(t, err)

	err = engine.VerifyHeader(chain, block.Header(), false)
	require.True(t, errors.Is(err, types.ErrEmptyQuorumCertificate))

	header := block.Header()

	// non zero MixDigest
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	header = block.Header()
	header.MixDigest = common.BytesToHash([]byte("123456789"))
	err = engine.VerifyHeader(chain, header, false)
	require.True(t, errors.Is(err, errInvalidMixDigest))

	// invalid uncles hash
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	header = block.Header()
	header.UncleHash = common.BytesToHash([]byte("123456789"))
	err = engine.VerifyHeader(chain, header, false)
	require.True(t, errors.Is(err, errInvalidUncleHash))

	// invalid difficulty
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	require.True(t, errors.Is(err, errInvalidDifficulty))

	// invalid timestamp
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	header = block.Header()
	header.Time = 0
	err = engine.VerifyHeader(chain, header, false)
	require.True(t, errors.Is(err, errInvalidTimestamp))

	// future block
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	header = block.Header()
	header.Time = new(big.Int).Add(big.NewInt(now().Unix()), new(big.Int).SetUint64(10)).Uint64()
	err = engine.VerifyHeader(chain, header, false)
	require.True(t, errors.Is(err, consensus.ErrFutureTimestampBlock))

	// invalid nonce
	block, err = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	require.NoError(t, err)
	header = block.Header()
	copy(header.Nonce[:], hexutil.MustDecode("0x111111111111"))
	err = engine.VerifyHeader(chain, header, false)
	require.True(t, errors.Is(err, errInvalidNonce))
}

// insert block with valid quorum certificate in the chain.
// It also add the precommit to the msgStore so that we can successfully create activity proof for the following blocks
// It assumes that we have a single committee member
// This is needed for `makeBlockWithoutSeal` to generate another block correctly.
// It returns the block with the quorum certificate
func insertBlock(t *testing.T, chain *core.BlockChain, engine *Backend, b *types.Block) *types.Block {
	self := &chain.Genesis().Header().Epoch.Committee.Members[0]

	header := b.Header()
	precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), engine.Sign, self, 1)
	header.QuorumCertificate = types.NewAggregateSignature(precommit.Signature().(*blst.BlsSignature), precommit.Signers())
	blockWithCertificate := b.WithSeal(header) // improper use, we use the WithSeal function to substitute the header with the one with quorumCertificate set
	time.Sleep(1 * time.Second)                // wait a couple seconds so that the block has not future timestamp anymore and the block import is done
	_, err := chain.InsertChain(types.Blocks{blockWithCertificate})
	require.NoError(t, err)

	engine.MsgStore.Save(precommit)
	return blockWithCertificate
}

// The logic of this needs to change with respect of Autonity contact
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
		require.NoError(t, err)

		b, err = engine.AddSeal(b)
		require.NoError(t, err)

		b = insertBlock(t, chain, engine, b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	// reset the chain to simulate receiving new headers
	err = chain.Reset()
	require.NoError(t, err)

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}
	defer func() { now = time.Now }() // if not reassigned, it will influence future tests

	_, results := engine.VerifyHeaders(chain, headers, nil)

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			require.NoError(t, err)
			index++
			if index == size {
				break OUT1
			}
		case <-timeout.C:
			t.Fatal("timeout expired")
		}
	}
	// avoid data race for the re-assignment of now to time.Now
	chain.Stop()
}

// The logic of this needs to change with respect of Autonity contact
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
		require.NoError(t, err)

		b, err = engine.AddSeal(b)
		require.NoError(t, err)

		b = insertBlock(t, chain, engine, b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	// reset the chain to simulate receiving new headers
	err = chain.Reset()
	require.NoError(t, err)

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}
	defer func() { now = time.Now }() // if not reassigned, it will influence future tests

	const timeoutDura = 2 * time.Second

	// abort cases
	abort, results := engine.VerifyHeaders(chain, headers, nil)
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT2:
	for {
		select {
		case err := <-results:
			require.NoError(t, err)
			index++
			if index == 5 {
				abort <- struct{}{}
			}
			if index >= 15 {
				t.Errorf("verifyheaders should be aborted rapidly")
				break OUT2
			}
		case <-timeout.C:
			break OUT2
		}
	}
	t.Log(index)
	// avoid data race for the re-assignment of now to time.Now
	chain.Stop()
}

// The logic of this needs to change with respect of Autonity contact
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
		require.NoError(t, err)

		b, err = engine.AddSeal(b)
		require.NoError(t, err)

		b = insertBlock(t, chain, engine, b)

		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	// reset the chain to simulate receiving new headers
	err = chain.Reset()
	require.NoError(t, err)

	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}
	defer func() { now = time.Now }() // if not reassigned, it will influence future tests

	const timeoutDura = 2 * time.Second

	// error header cases
	headers[2].Number = big.NewInt(100)
	_, results := engine.VerifyHeaders(chain, headers, nil)
	timeout := time.NewTimer(timeoutDura)
	index := 0
	errorCount := 0
	// header[2] out of epoch range | header[3] != header[2]+1 | header[12] invalid activity proof
	expectedErrors := 3

OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				t.Logf("received error: %v", err)
				errorCount++
			}
			index++
			if index == size {
				require.Equal(t, expectedErrors, errorCount)
				break OUT3
			}
		case <-timeout.C:
			t.Fatal("timeout expired")
		}
	}
	// avoid data race for the re-assignment of now to time.Now
	chain.Stop()
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

// needed because backend.Close() also stops the aggregator. It checks that Stop() is called at maximum once
func fakeAggregator() *aggregator {
	mux := new(event.TypeMux)
	stopped := false
	fakeAggregator := &aggregator{
		logger: log.Root(),
		cancel: func() {
			if !stopped {
				stopped = true
			} else {
				// already stopped once
				panic("aggregator stopped two times")
			}
		},
		coreSub: mux.Subscribe(),
	}
	return fakeAggregator
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
			core:       tendermintC,
			aggregator: fakeAggregator(),
			stopped:    make(chan struct{}),
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
			core:       tendermintC,
			aggregator: fakeAggregator(),
			stopped:    make(chan struct{}),
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
			core:       tendermintC,
			aggregator: fakeAggregator(),
			stopped:    make(chan struct{}),
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
		tendermintC.EXPECT().Height().Return(common.Big1).AnyTimes()
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:       tendermintC,
			gossiper:   g,
			blockchain: chain,
			eventMux:   event.NewTypeMuxSilent(nil, log.Root()),
		}
		b.aggregator = &aggregator{logger: log.Root(), backend: b, core: tendermintC}

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
		tendermintC.EXPECT().Height().Return(common.Big1).AnyTimes()
		chain, _ := newBlockChain(1)
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:       tendermintC,
			gossiper:   g,
			blockchain: chain,
			eventMux:   event.NewTypeMuxSilent(nil, log.Root()),
		}
		b.aggregator = &aggregator{logger: log.Root(), backend: b, core: tendermintC}
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
		tendermintC.EXPECT().Height().Return(common.Big1).AnyTimes()
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:       tendermintC,
			gossiper:   g,
			blockchain: chain,
			eventMux:   event.NewTypeMuxSilent(nil, log.Root()),
		}
		b.aggregator = &aggregator{logger: log.Root(), backend: b, core: tendermintC}
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
	tendermintC.EXPECT().Height().Return(common.Big1).AnyTimes()
	chain, _ := newBlockChain(1)
	g := interfaces.NewMockGossiper(ctrl)
	g.EXPECT().UpdateStopChannel(gomock.Any()).MaxTimes(5)

	b := &Backend{
		core:       tendermintC,
		gossiper:   g,
		blockchain: chain,
		eventMux:   event.NewTypeMuxSilent(nil, log.Root()),
	}
	b.aggregator = &aggregator{logger: log.Root(), backend: b, core: tendermintC}
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
