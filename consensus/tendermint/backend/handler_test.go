package backend

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
)

func TestTendermintCaches(t *testing.T) {
	_, backend := newBlockChain(1)

	// generate a msg
	data := message.NewPrevote(1, 1, common.Hash{}, backend.Sign)
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
	require.NoError(t, err)

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

// future height messages should be buffered
func TestFutureMsg(t *testing.T) {
	t.Run("future height messages should be buffered and processed when we reach that height", func(t *testing.T) {
		// clean blockchain, running consensus instance for height 1
		blockchain, backend := newBlockChain(1)

		// create valid future proposal encapsulated in p2p.Message
		futureHeight := uint64(3)
		round := int64(0)
		data := message.NewPropose(round, futureHeight, -1, types.NewBlockWithHeader(&types.Header{}), backend.Sign)
		msg := p2p.Msg{Code: ProposeNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

		// send it to the backend
		errCh := make(chan error)
		handled, err := backend.HandleMsg(backend.Address(), msg, errCh)
		require.True(t, handled)
		require.NoError(t, err)

		// check if buffered
		require.Equal(t, data.Hash(), backend.future[futureHeight][0].Message.Hash())
		require.Equal(t, uint64(1), backend.futureSize)

		// advance to consensus instance for height 2
		advanceBlockchain(t, backend, blockchain)

		// should still be buffered
		require.Equal(t, data.Hash(), backend.future[futureHeight][0].Message.Hash())
		require.Equal(t, uint64(1), backend.futureSize)

		// advance to consensus instance for height 3, this should trigger the processing
		advanceBlockchain(t, backend, blockchain)

		// sleep some time to make sure that the processFutureMessage go routine has run
		time.Sleep(5 * time.Second)

		// check that not buffered anymore
		require.Equal(t, uint64(0), backend.futureSize)
		require.Len(t, backend.future[futureHeight], 0)
	})
}

func TestSaveFutureMessage(t *testing.T) {
	t.Run("save future messages", func(t *testing.T) {
		_, backend := newBlockChain(1)

		var messages []message.Msg

		for i := int64(0); i < maxFutureMsgs; i++ {
			round := i % 10
			height := uint64((i / (1 + i%10)) + 2)
			msg := message.NewPrevote(round, height, common.Hash{}, backend.Sign)
			backend.saveFutureMsg(msg, nil)
			messages = append(messages, msg)
		}
		found := 0
		for _, msg := range messages {
			for _, umsg := range backend.future[msg.H()] {
				//TODO(lorenzo) does this still work
				if deep.Equal(msg, umsg.Message) {
					found++
				}
			}
		}
		require.Equal(t, maxFutureMsgs, int(backend.futureSize))
		require.Equal(t, maxFutureMsgs, found)
	})
	t.Run("excess messages are removed from the future messages backlog", func(t *testing.T) {
		_, backend := newBlockChain(1)

		var messages []message.Msg

		// we are at consensus instance for height = 1
		for i := int64(2*maxFutureMsgs) + 1; i > 1; i-- {
			round := i % 10
			msg := message.NewPrevote(round, uint64(i), common.Hash{}, backend.Sign)
			backend.saveFutureMsg(msg, nil)
			if i <= maxFutureMsgs+1 {
				messages = append(messages, msg)
			}
		}

		found := 0
		for _, msg := range messages {
			for _, umsg := range backend.future[msg.H()] {
				//TODO(lorenzo) does this still work
				if deep.Equal(msg, umsg.Message) {
					found++
				}
			}
		}
		require.Equal(t, int(backend.futureSize), maxFutureMsgs)
		require.Equal(t, int(backend.futureSize), len(backend.future))
		require.Equal(t, int(backend.futureSize), found)
		require.Equal(t, 2, int(backend.futureMinHeight))
		require.Equal(t, maxFutureMsgs+1, int(backend.futureMaxHeight))
	})
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

func TestMsgFromJailedValidator(t *testing.T) {
	_, backend := newBlockChain(1)

	// jail the test node and subscribe to events
	backend.jailed[backend.Address()] = 0
	sub := backend.Subscribe(events.MessageEvent{}, events.OldMessageEvent{})

	// create message
	data := message.NewPrevote(0, 1, common.Hash{}, backend.Sign)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.NoError(t, err)

	// make sure message has not been processed
	timeout := time.NewTimer(2 * time.Second)
	select {
	case ev := <-sub.Chan():
		switch ev.Data.(type) {
		case events.MessageEvent:
			t.Fatal("Message from jailed validator treated as valid message")
		case events.OldMessageEvent:
			t.Fatal("Message from jailed validator treated as old valid message")
		}
	case <-timeout.C:
		// good case
	}
}

func TestMsgFromOutOfCommittee(t *testing.T) {
	_, backend := newBlockChain(1)

	// create message
	data := message.NewPrevote(0, 1, common.Hash{}, testSigner)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.Equal(t, message.ErrUnauthorizedAddress, err)
}

func TestInvalidSignatureLenMsg(t *testing.T) {
	_, backend := newBlockChain(1)

	invalidSigner := func(data common.Hash) ([]byte, common.Address, *big.Int) {
		return []byte{0xca, 0xfe}, testAddress, testPower
	}
	data := message.NewPrevote(0, 1, common.Hash{}, invalidSigner)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	_, err := backend.HandleMsg(testAddress, msg, errCh)
	require.Equal(t, message.ErrBadSignature, err)
}

func TestInvalidSignatureMsg(t *testing.T) {
	_, backend := newBlockChain(1)

	invalidSigner := func(data common.Hash) ([]byte, common.Address, *big.Int) {
		out, _ := crypto.Sign([]byte{0xca, 0xfe}, testKey)
		return out, testAddress, testPower
	}
	data := message.NewPrevote(0, 1, common.Hash{}, invalidSigner)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	_, err := backend.HandleMsg(testAddress, msg, errCh)
	require.Equal(t, message.ErrBadSignature, err)
}

func TestInvalidRound(t *testing.T) {
	_, backend := newBlockChain(1)

	sub := backend.Subscribe(events.MessageEvent{}, events.OldMessageEvent{})

	data := message.NewPrevote(1000, 1, common.Hash{}, testSigner)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.Error(t, err)

	// make sure message has not been processed
	timeout := time.NewTimer(2 * time.Second)
	select {
	case ev := <-sub.Chan():
		switch ev.Data.(type) {
		case events.MessageEvent:
			t.Fatal("Invalid inner message treated as valid message")
		case events.OldMessageEvent:
			t.Fatal("Invalid inner message treated as old valid message")
		}
	case <-timeout.C:
		// good case
	}
}

func TestValidMsg(t *testing.T) {
	_, backend := newBlockChain(1)

	sub := backend.Subscribe(events.MessageEvent{}, events.OldMessageEvent{})

	data := message.NewPrevote(0, 1, common.Hash{}, backend.Sign)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.NoError(t, err)

	// make sure message has been processed
	timeout := time.NewTimer(2 * time.Second)
	select {
	case ev := <-sub.Chan():
		switch ev.Data.(type) {
		case events.MessageEvent:
			// the good case
		case events.OldMessageEvent:
			t.Fatal("valid message treated as old message")
		}
	case <-timeout.C:
		t.Fatal("valid message not processed")
	}
}

func TestValidOldMsg(t *testing.T) {
	blockchain, backend := newBlockChain(1)

	advanceBlockchain(t, backend, blockchain)

	sub := backend.Subscribe(events.MessageEvent{}, events.OldMessageEvent{})

	data := message.NewPrevote(0, 1, common.Hash{}, backend.Sign)
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.NoError(t, err)

	// make sure message has not been processed
	timeout := time.NewTimer(2 * time.Second)
	select {
	case ev := <-sub.Chan():
		switch ev.Data.(type) {
		case events.MessageEvent:
			t.Fatal("old message treated as valid message")
		case events.OldMessageEvent:
			// the good case
		}
	case <-timeout.C:
		t.Fatal("old message not processed by accountability")
	}
}

func TestGarbagePropose(t *testing.T) {
	_, backend := newBlockChain(1)

	data := make([]byte, 1024)
	_, err := rand.Read(data)
	require.NoError(t, err, "error while generating random bytes")
	msg := p2p.Msg{Code: ProposeNetworkMsg, Size: uint32(len(data)), Payload: bytes.NewReader(data)}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.Error(t, err)
}

func TestGarbagePrevote(t *testing.T) {
	_, backend := newBlockChain(1)

	data := make([]byte, 1024)
	_, err := rand.Read(data)
	require.NoError(t, err, "error while generating random bytes")
	msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(data)), Payload: bytes.NewReader(data)}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.Error(t, err)
}

func TestGarbagePrecommit(t *testing.T) {
	_, backend := newBlockChain(1)

	data := make([]byte, 1024)
	_, err := rand.Read(data)
	require.NoError(t, err, "error while generating random bytes")
	msg := p2p.Msg{Code: PrecommitNetworkMsg, Size: uint32(len(data)), Payload: bytes.NewReader(data)}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.Error(t, err)
}

func TestGarbageInvalidCode(t *testing.T) {
	_, backend := newBlockChain(1)

	data := message.NewPrecommit(0, 1, common.Hash{}, testSigner)
	msg := p2p.Msg{Code: PrecommitNetworkMsg + 100, Size: uint32(len(data.Payload())), Payload: bytes.NewReader(data.Payload())}

	// handle message, it should not be handled
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.False(t, handled)
	require.NoError(t, err)
}

func TestGarbageOversized(t *testing.T) {
	_, backend := newBlockChain(1)

	data := make([]byte, 1024*9)
	_, err := rand.Read(data)
	require.NoError(t, err, "error while generating random bytes")

	msg := p2p.Msg{Code: PrecommitNetworkMsg, Size: uint32(len(data)), Payload: bytes.NewReader(data)}

	// handle message
	errCh := make(chan error, 1)
	handled, err := backend.HandleMsg(testAddress, msg, errCh)
	require.True(t, handled)
	require.Error(t, err)
}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}
