package backend

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/consensus/tendermint/helpers"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

func setupMocks(backend *Backend, ctrl *gomock.Controller, t *testing.T) {
	if err := backend.Close(); err != nil { // close engine to avoid race while updating the broadcaster
		t.Fatalf("can't stop the engine")
	}
	member := backend.blockchain.Genesis().Header().Epoch.Committee.Members[0]
	sender := consensus.NewMockPeer(ctrl)
	committeeMember := consensus.NewMockPeer(ctrl)
	broadcaster := consensus.NewMockBroadcaster(ctrl)
	addressCache := fixsizecache.New[common.Hash, bool](1997, 10, fixsizecache.HashKey[common.Hash])
	sender.EXPECT().Cache().Return(addressCache).AnyTimes()
	committeeMember.EXPECT().Cache().Return(addressCache).AnyTimes()
	broadcaster.EXPECT().FindPeer(testAddress).Return(sender, true).AnyTimes()
	broadcaster.EXPECT().FindPeer(member.Address).Return(committeeMember, true).AnyTimes()
	committeeMember.EXPECT().SendRaw(gomock.Any(), gomock.Any()).AnyTimes()

	backend.SetBroadcaster(broadcaster)

	if err := backend.Start(context.Background()); err != nil {
		t.Fatalf("could not restart core")
	}
}

func TestTendermintMessage(t *testing.T) {
	_, backend := newBlockChain(1)
	// generate one msg
	data := message.NewPrevote(1, 2, common.Hash{}, testSigner, testCommitteeMember, 1)
	msg := p2p.Msg{Code: message.PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())} // #nosec

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	setupMocks(backend, ctrl, t)

	// 1. this message should not be in cache
	// for peers
	if peer, ok := backend.Broadcaster.FindPeer(testAddress); ok {
		if peer.Cache().Contains(data.Hash()) {
			t.Fatalf("the cache of messages for this peer should be empty")
		}
	}

	// for self
	if _, ok := backend.knownMessages.Get(data.Hash()); ok {
		t.Fatalf("the cache of messages should be nil")
	}

	// 2. this message should be in cache after we handle it
	errCh := make(chan error, 1)
	_, err := backend.HandleMsg(testAddress, msg, errCh)
	if err != nil {
		t.Fatalf("handle message failed: %v", err)
	}
	// for peers
	if peer, ok := backend.Broadcaster.FindPeer(testAddress); ok {
		cache := peer.Cache()
		if !cache.Contains(data.Hash()) {
			t.Fatalf("the cache of messages for this peer cannot be found")
		}
	}

	// for self
	if _, ok := backend.knownMessages.Get(data.Hash()); !ok {
		t.Fatalf("the cache of messages cannot be found")
	}
}
func TestSynchronisationMessage(t *testing.T) {
	t.Run("engine not running, ignored", func(t *testing.T) {
		b := &Backend{
			database: rawdb.NewMemoryDatabase(),
			logger:   log.New("backend", "test", "id", 0),
		}
		msg := makeMsg(message.SyncNetworkMsg, []byte{})
		addr := common.BytesToAddress([]byte("address"))
		errCh := make(chan error, 1)
		if res, err := b.HandleMsg(addr, msg, errCh); !res || err != nil {
			t.Fatalf("HandleMsg unexpected return")
		}
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-errCh:
			t.Fatalf("not expected message")
		case <-timer.C:
		}
	})

	t.Run("engine running, msg cannot be decoded", func(t *testing.T) {
		b := &Backend{
			database:           rawdb.NewMemoryDatabase(),
			logger:             log.New("backend", "test", "id", 0),
			askSyncRateLimiter: helpers.NewTimeWindowLimiter(constants.AskSyncInterval, 2),
		}
		b.coreStarting.Store(true)
		b.coreRunning.Store(true)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockedPeer := consensus.NewMockPeer(ctrl)
		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeer(testAddress).Return(mockedPeer, true).AnyTimes()
		b.Broadcaster = broadcaster

		b.coreStarting.Store(true)
		b.coreRunning.Store(true)
		msg := makeMsg(message.SyncNetworkMsg, []byte{})
		errCh := make(chan error, 1)
		if res, err := b.HandleMsg(testAddress, msg, errCh); !res || err != nil {
			t.Fatalf("HandleMsg unexpected return")
		}
		select {
		case err := <-errCh:
			require.NotNil(t, err)
		}
	})
}

func TestNewChainHead(t *testing.T) {
	t.Run("engine not started, error returned", func(t *testing.T) {
		b := &Backend{
			database: rawdb.NewMemoryDatabase()}

		err := b.NewChainHead()
		if err != ErrStoppedEngine {
			t.Fatalf("expected %v, got %v", ErrStoppedEngine, err)
		}
	})

	t.Run("engine is running, no errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		coreEventCh := make(chan events.CoreEvent, 10)
		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Start(gomock.Any(), gomock.Any()).MaxTimes(1)
		tendermintC.EXPECT().Height().Return(common.Big1).AnyTimes()
		tendermintC.EXPECT().EventCh().Return(coreEventCh).AnyTimes()
		evDispathcer := interfaces.NewMockEventDispatcher(ctrl)
		evDispathcer.EXPECT().Post(gomock.Any()).MaxTimes(1)
		chain, _ := newBlockChain(1)
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())
		mockRouter := interfaces.NewMockRouter(ctrl)
		mockRouter.EXPECT().Start(gomock.Any(), gomock.Any()).MaxTimes(1)

		b := &Backend{
			database:            rawdb.NewMemoryDatabase(),
			core:                tendermintC,
			coreEventDispatcher: evDispathcer,
			gossiper:            g,
			blockchain:          chain,
			eventMux:            event.NewTypeMuxSilent(nil, log.Root()),
			askSyncRateLimiter:  helpers.NewTimeWindowLimiter(constants.AskSyncInterval, 2),
			logger:              log.Root(),
			router:              mockRouter,
		}
		b.aggregator = &aggregator{logger: log.Root(), backend: b, core: tendermintC}
		b.Start(ctx)

		err := b.NewChainHead()
		if err != nil {
			t.Fatalf("expected <nil>, got %v", err)
		}
	})
}
func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	var buff bytes.Buffer
	io.Copy(&buff, r)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: bytes.NewReader(buff.Bytes())}
}

func TestSignerJailed(t *testing.T) {
	chain, backend := newBlockChain(1)

	member := chain.Genesis().Header().Epoch.Committee.Members[0]

	// generate one msg
	data := message.NewPrevote(0, 1, common.Hash{}, testSigner, &member, 1)
	msg := p2p.Msg{Code: message.PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())} // #nosec

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	setupMocks(backend, ctrl, t)

	backend.jailed.Lock()
	backend.jailed.validators[member.Address] = 0
	backend.jailed.Unlock()

	errCh := make(chan error, 1)
	_, err := backend.HandleMsg(testAddress, msg, errCh)
	require.Equal(t, ErrJailed, err)

	// same should happen for an aggregate containing a single jailed signer

	data = message.NewPrevote(0, 1, common.Hash{0xca, 0xfe}, testSigner, &member, 2)
	data.Signers().AddMember(makeBogusMember(1))
	msg = p2p.Msg{Code: message.PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())} // #nosec
	errCh = make(chan error, 1)
	_, err = backend.HandleMsg(testAddress, msg, errCh)
	require.Equal(t, ErrJailed, err)
}

func TestFutureHeightMessage(t *testing.T) {
	t.Run("received future height message is buffered", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		member := chain.Genesis().Header().Epoch.Committee.Members[0]

		// generate one msg
		futureHeight := uint64(20)
		data := message.NewPrevote(0, futureHeight, common.Hash{}, testSigner, &member, 1)
		msg := p2p.Msg{Code: message.PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())} // #nosec

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		setupMocks(backend, ctrl, t)

		errCh := make(chan error, 1)
		_, err := backend.HandleMsg(testAddress, msg, errCh)
		require.NoError(t, err)

		backend.future.RLock()
		defer backend.future.RUnlock()
		require.Equal(t, 1, len(backend.future.messages[futureHeight]))
		require.Equal(t, data.Hash(), backend.future.messages[futureHeight][0].Message.Hash())
		require.Equal(t, futureHeight, backend.future.maxHeight)
		require.Equal(t, uint64(1), backend.future.size)
	})
	t.Run("if future message buffer is full, messages farther in the future are dropped", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		setupMocks(backend, ctrl, t)

		member := chain.Genesis().Header().Epoch.Committee.Members[0]

		for h := maxFutureMsgs + 100; h > 0; h-- {
			data := message.NewPrevote(0, uint64(h), common.Hash{}, testSigner, &member, 1)
			msg := p2p.Msg{Code: message.PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())} // #nosec
			errCh := make(chan error, 1)
			_, err := backend.HandleMsg(testAddress, msg, errCh)
			require.NoError(t, err)
		}

		backend.future.RLock()
		defer backend.future.RUnlock()
		require.Equal(t, maxFutureMsgs, len(backend.future.messages))
		require.Equal(t, uint64(maxFutureMsgs), backend.future.size) // works because we send only one message per height
	})
	t.Run("When processing future height messages, future height messages are re-injected", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		setupMocks(backend, ctrl, t)

		member := chain.Genesis().Header().Epoch.Committee.Members[0]

		vote := message.NewPrevote(0, 1, common.Hash{}, testSigner, &member, 1)
		errCh := make(chan error, 1)
		backend.saveFutureMsg(vote, errCh, common.Address{})
		backend.saveFutureMsg(vote, errCh, common.Address{})
		backend.saveFutureMsg(vote, errCh, common.Address{})
		backend.saveFutureMsg(vote, errCh, common.Address{})

		backend.future.RLock()
		require.Equal(t, uint64(4), backend.future.size)
		backend.future.RUnlock()

		backend.ProcessFutureMsgs(1)

		backend.future.RLock()
		require.Equal(t, uint64(0), backend.future.size)
		backend.future.RUnlock()
	})
}

// invalid proposal should be caught at handler level
func TestInvalidProposal(t *testing.T) {
	t.Run("handler rejects invalid proposal", func(t *testing.T) {
		_, backend := newBlockChain(1)

		propose := message.NewFakePropose(message.Fake{
			FakeSignatureInput: common.Hash{0xca, 0xfe},
			FakeSignerKey:      testKey.PublicKey(),
			FakeSignature:      testKey.Sign([]byte{0xff, 0xff}), // signature is not on FakeSignatureInput --> invalid
			FakeVerified:       false,
			FakeHash:           common.Hash{0xee, 0xee},
		})

		handled, err := backend.handleDecodedMsg(propose, nil, common.Address{})
		require.True(t, handled)
		t.Logf("proposal failed with error: %v", err)
		require.Error(t, err)
		require.True(t, errors.Is(err, message.ErrBadSignature))
	})
}
