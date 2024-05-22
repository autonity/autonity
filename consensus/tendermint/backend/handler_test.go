package backend

import (
	"bytes"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"

	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"

	"github.com/autonity/autonity/common"
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

	// 1. this message should not be in cache
	// for peers
	if _, ok := backend.recentMessages.Get(testAddress); ok {
		t.Fatalf("the cache of messages for this peer should be nil")
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
	if ms, ok := backend.recentMessages.Get(testAddress); ms == nil || !ok {
		t.Fatalf("the cache of messages for this peer cannot be nil")
	} else if m, ok := ms.(*lru.ARCCache); !ok {
		t.Fatalf("the cache of messages for this peer cannot be casted")
	} else if _, ok := m.Get(data.Hash()); !ok {
		t.Fatalf("the cache of messages for this peer cannot be found")
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
		b := &Backend{
			eventMux: event.NewTypeMuxSilent(nil, log.New("backend", "test", "id", 0)),
		}
		b.coreRunning.Store(true)

		err := b.NewChainHead()
		if err != nil {
			t.Fatalf("expected <nil>, got %v", err)
		}
	})
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
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
