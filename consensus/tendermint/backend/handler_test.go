package backend

import (
	"bytes"
	"io"
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
	data := message.NewPrevote(1, 2, common.Hash{}, testSigner)
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
	tic := time.NewTicker(time.Millisecond * 100)
	maxWait := 50
	counter := 0
	peerCacheNil, peerCacheEmpty, selfCacheEmpty := false, false, false

	for {
		<-tic.C
		if ms, ok := backend.recentMessages.Get(testAddress); ms == nil || !ok {
			peerCacheNil = true
		} else if _, ok := ms.(*lru.ARCCache).Get(data.Hash()); !ok {
			peerCacheEmpty = true
		}
		if _, ok := backend.knownMessages.Get(data.Hash()); !ok {
			selfCacheEmpty = true
		}
		if !peerCacheNil && !peerCacheEmpty && !selfCacheEmpty {
			break
		}
		peerCacheNil, peerCacheEmpty, selfCacheEmpty = false, false, false
		if counter >= maxWait {
			t.Fatalf("the cache of messages cannot be found")
		}
	}
}

func TestSynchronisationMessage(t *testing.T) {
	t.Run("engine not running, ignored", func(t *testing.T) {
		eventMux := event.NewTypeMuxSilent(nil, log.New("backend", "test", "id", 0))
		sub := eventMux.Subscribe(events.SyncEvent{})
		b := &Backend{
			coreStarted: false,
			logger:      log.New("backend", "test", "id", 0),
			eventMux:    eventMux,
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
			coreStarted: true,
			logger:      log.New("backend", "test", "id", 0),
			eventMux:    eventMux,
		}
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

func TestProtocol(t *testing.T) {
	b := &Backend{}
	name, code := b.Protocol()
	if name != "tendermint" {
		t.Fatalf("expected 'tendermint', got %v", name)
	}
	if code != 5 {
		t.Fatalf("expected 2, got %v", code)
	}
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
			coreStarted: true,
			eventMux:    event.NewTypeMuxSilent(nil, log.New("backend", "test", "id", 0)),
		}

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
