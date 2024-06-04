package backend

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
)

func TestTendermintMessage(t *testing.T) {
	_, backend := newBlockChain(1)
	// generate one msg
	data := message.NewPrevote(1, 2, common.Hash{}, testSigner, testCommitteeMember, 1)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	if err := backend.Close(); err != nil { // close engine to avoid race while updating the broadcaster
		t.Fatalf("can't stop the engine")
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedPeer := consensus.NewMockPeer(ctrl)
	broadcaster := consensus.NewMockBroadcaster(ctrl)
	addressCache := fixsizecache.New[common.Hash, bool](1997, 10, fixsizecache.HashKey[common.Hash])
	mockedPeer.EXPECT().Cache().Return(addressCache).AnyTimes()
	broadcaster.EXPECT().FindPeer(testAddress).Return(mockedPeer, true).AnyTimes()
	backend.SetBroadcaster(broadcaster)

	if err := backend.Start(context.Background()); err != nil {
		t.Fatalf("could not restart core")
	}
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
		eventMux := event.NewTypeMuxSilent(nil, log.New("backend", "test", "id", 0))
		sub := eventMux.Subscribe(events.SyncEvent{})
		b := &Backend{
			logger:   log.New("backend", "test", "id", 0),
			eventMux: eventMux,
		}
		msg := makeMsg(SyncNetworkMsg, []byte{})
		addr := common.BytesToAddress([]byte("address"))
		errCh := make(chan error, 1)
		if res, err := b.HandleMsg(addr, msg, errCh); !res || err != nil {
			t.Fatalf("HandleMsg unexpected return")
		}
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-sub.Chan():
			t.Fatalf("not expected message")
		case <-timer.C:
		}
	})

	t.Run("engine running, sync returned", func(t *testing.T) {
		eventMux := event.NewTypeMuxSilent(nil, log.New("backend", "test", "id", 0))
		sub := eventMux.Subscribe(events.SyncEvent{})
		b := &Backend{
			logger:   log.New("backend", "test", "id", 0),
			eventMux: eventMux,
		}
		b.coreStarting.Store(true)
		b.coreRunning.Store(true)
		msg := makeMsg(SyncNetworkMsg, []byte{})
		addr := common.BytesToAddress([]byte("address"))
		errCh := make(chan error, 1)
		if res, err := b.HandleMsg(addr, msg, errCh); !res || err != nil {
			t.Fatalf("HandleMsg unexpected return")
		}
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-timer.C:
			t.Fatalf("sync message not posted")
		case <-sub.Chan():
		}
	})
}

func TestNewChainHead(t *testing.T) {
	t.Run("engine not started, error returned", func(t *testing.T) {
		b := &Backend{}

		err := b.NewChainHead()
		if err != ErrStoppedEngine {
			t.Fatalf("expected %v, got %v", ErrStoppedEngine, err)
		}
	})

	t.Run("engine is running, no errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().Start(gomock.Any(), gomock.Any()).MaxTimes(1)
		tendermintC.EXPECT().Height().Return(common.Big1).AnyTimes()
		evDispathcer := interfaces.NewMockEventDispatcher(ctrl)
		evDispathcer.EXPECT().Post(gomock.Any()).MaxTimes(1)
		chain, _ := newBlockChain(1)
		g := interfaces.NewMockGossiper(ctrl)
		g.EXPECT().UpdateStopChannel(gomock.Any())

		b := &Backend{
			core:         tendermintC,
			evDispatcher: evDispathcer,
			gossiper:     g,
			blockchain:   chain,
			eventMux:     event.NewTypeMuxSilent(nil, log.Root()),
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

//TODO(lorenzo) add tests for:
// - receiving msgs from jailed validators
// - receiving and processing future height messages
// - receiving msg from non-committee member

/* TODO(lorenzo) port this tests which were in Core before. Now future height messages are in the backend
func TestStoreUncheckedBacklog(t *testing.T) {
	t.Run("save messages in the untrusted backlog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		backendMock := interfaces.NewMockBackend(ctrl)
		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			backlogUntrusted: make(map[uint64][]message.Msg),
			step:             Prevote,
			round:            1,
			height:           big.NewInt(4),
		}
		var messages []message.Msg
		for i := int64(0); i < MaxSizeBacklogUnchecked; i++ {
			msg := message.NewPrevote(
				i%10,
				uint64(i/(1+i%10)),
				common.Hash{},
				defaultSigner)
			c.storeFutureMessage(msg)
			messages = append(messages, msg)
		}
		found := 0
		for _, msg := range messages {
			for _, umsg := range c.backlogUntrusted[msg.H()] {
				if deep.Equal(msg, umsg) {
					found++
				}
			}
		}
		if found != MaxSizeBacklogUnchecked {
			t.Fatal("unchecked messages lost")
		}
	})

	t.Run("excess messages are removed from the untrusted backlog", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)

		c := &Core{
			logger:           log.New("backend", "test", "id", 0),
			backend:          backendMock,
			address:          common.HexToAddress("0x1234567890"),
			backlogs:         make(map[common.Address][]message.Msg),
			backlogUntrusted: make(map[uint64][]message.Msg),
			step:             Prevote,
			round:            1,
			height:           big.NewInt(4),
		}

		var messages []message.Msg
		uncheckedFounds := make(map[uint64]struct{})
		backendMock.EXPECT().RemoveMessageFromLocalCache(gomock.Any()).Times(MaxSizeBacklogUnchecked).Do(func(msg message.Msg) {
			if _, ok := uncheckedFounds[msg.H()]; ok {
				t.Fatal("duplicate message received")
			}
			uncheckedFounds[msg.H()] = struct{}{}
		})

		for i := int64(2 * MaxSizeBacklogUnchecked); i > 0; i-- {
			prevote := message.NewPrevote(i%10, uint64(i), common.Hash{}, defaultSigner)
			c.storeFutureMessage(prevote)
			if i < MaxSizeBacklogUnchecked {
				messages = append(messages, prevote)
			}
		}

		found := 0
		for _, msg := range messages {
			for _, umsg := range c.backlogUntrusted[msg.H()] {
				if deep.Equal(msg, umsg) {
					found++
				}
			}
		}
		if found != MaxSizeBacklogUnchecked-1 {
			t.Fatal("unchecked messages lost")
		}
		if len(uncheckedFounds) != MaxSizeBacklogUnchecked {
			t.Fatal("unchecked messages lost")
		}
	})
}*/
