package backend

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/rlp"
	"github.com/hashicorp/golang-lru"
)

func TestTendermintMessage(t *testing.T) {
	_, backend := newBlockChain(1)

	// generate one msg
	data := []byte("data1")
	//msg := makeMsg(TendermintMsg, data)
	msg := makeMsgVote(data, &message.Vote{})

	// payload is a reader so we need to recreate it
	payloadBytes, err := io.ReadAll(msg.Payload)
	require.NoError(t, err)
	msg.Payload = bytes.NewReader(payloadBytes)

	hash := types.RLPHash(payloadBytes)

	fmt.Printf("expectedHash %x\n", hash)

	addr := common.BytesToAddress([]byte("address"))

	// 1. this message should not be in cache
	// for peers
	if _, ok := backend.recentMessages.Get(addr); ok {
		t.Fatalf("the cache of messages for this peer should be nil")
	}

	// for self
	if _, ok := backend.knownMessages.Get(hash); ok {
		t.Fatalf("the cache of messages should be nil")
	}

	// 2. this message should be in cache after we handle it
	errCh := make(chan error, 1)
	_, err = backend.HandleMsg(addr, msg, errCh)
	if err != nil {
		t.Fatalf("handle message failed: %v", err)
	}
	// for peers
	if ms, ok := backend.recentMessages.Get(addr); ms == nil || !ok {
		t.Fatalf("the cache of messages for this peer cannot be nil")
	} else if m, ok := ms.(*lru.ARCCache); !ok {
		t.Fatalf("the cache of messages for this peer cannot be casted")
	} else if _, ok := m.Get(hash); !ok {
		t.Fatalf("the cache of messages for this peer cannot be found")
	}

	// for self
	if _, ok := backend.knownMessages.Get(hash); !ok {
		t.Fatalf("the cache of messages cannot be found")
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
		msg := makeMsg(SyncMsg, []byte{})
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
		msg := makeMsg(SyncMsg, []byte{})
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
	if code != 2 {
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

//func makeMsg(msgcode uint64, data []byte) p2p.Msg {
//	return makeMsgVote(msgcode, data, nil)
//}

func makeMsg(msgcode uint64, data interface{}) p2p.Msg {
	size, r, _ := rlp.EncodeToReader(data)
	return p2p.Msg{Code: msgcode, Size: uint32(size), Payload: r}
}

func makeMsgVote(payload []byte, vote *message.Vote) p2p.Msg {

	msg := &message.MessageVote{
		MessageBase: message.MessageBase{
			Code:          0,
			Payload:       payload,
			Address:       common.Address{},
			Signature:     nil,
			CommittedSeal: nil,
			Power:         nil,
			Bytes:         nil,
		},
		Vote: *vote,
	}

	encoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(fmt.Sprintf("makeMsg EncodeToReader failed: %s", err))
	}

	fmt.Printf("NEW MESSAGE PAYLOAD: %s\n", hex.EncodeToString(encoded))

	r := bytes.NewReader(encoded)
	size := len(encoded)

	//size, r, err := rlp.EncodeToReader(msg)
	//if err != nil {
	//	panic(fmt.Sprintf("makeMsg EncodeToReader failed: %s", err))
	//}
	//
	return p2p.Msg{Code: TendermintMsgVote, Size: uint32(size), Payload: r}
}

func makeMsgProposal(payload []byte, proposal *message.Proposal) p2p.Msg {

	msg := &message.MessageProposal{
		MessageBase: message.MessageBase{
			Code:          0,
			Payload:       payload,
			Address:       common.Address{},
			Signature:     nil,
			CommittedSeal: nil,
			Power:         nil,
			Bytes:         nil,
		},
		Proposal: *proposal,
	}

	encoded, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(fmt.Sprintf("makeMsg EncodeToReader failed: %s", err))
	}

	fmt.Printf("NEW MESSAGE PAYLOAD: %s\n", hex.EncodeToString(encoded))

	r := bytes.NewReader(encoded)
	size := len(encoded)

	//size, r, err := rlp.EncodeToReader(msg)
	//if err != nil {
	//	panic(fmt.Sprintf("makeMsg EncodeToReader failed: %s", err))
	//}
	//
	return p2p.Msg{Code: TendermintMsgProposal, Size: uint32(size), Payload: r}
}
