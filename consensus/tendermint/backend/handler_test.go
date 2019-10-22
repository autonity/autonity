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

package backend

import (
	"github.com/clearmatics/autonity/consensus"
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rlp"
	"github.com/hashicorp/golang-lru"
)

func TestUnhandledMsgs(t *testing.T) {
	t.Run("core not running, unhandled messages are saved", func(t *testing.T) {
		blockchain, backend := newBlockChain(1)
		engine := blockchain.Engine().(consensus.BFT)
		// we close the engine for enabling cache storing
		if err := engine.Close(); err != nil {
			t.Fatalf("can't stop the engine")
		}
		//we generate a bunch of messages until we reach max capacity
		for i := int64(0); i < 2*ringCapacity; i++ {
			counter := big.NewInt(i).Bytes()
			msg := makeMsg(tendermintMsg, append(counter, []byte("data")...))
			addr := common.BytesToAddress(append(counter, []byte("addr")...))
			if result, err := backend.HandleMsg(addr, msg); !result || err != nil {
				t.Fatalf("handleMsg should have been successful")
			}
		}

		for i := int64(0); i < ringCapacity; i++ {
			counter := big.NewInt(i + ringCapacity).Bytes() // messages i < ringCapacity should have been discarded
			savedMsg := backend.pendingMessages.Dequeue()
			if savedMsg == nil {
				t.Fatalf("missing message")
			}
			addr := savedMsg.(UnhandledMsg).addr
			expectedAddr := common.BytesToAddress(append(counter, []byte("addr")...))
			if savedMsg.(UnhandledMsg).msg.Code != tendermintMsg {
				t.Fatalf("wrong msg code")
			}
			var payload []byte
			if err := savedMsg.(UnhandledMsg).msg.Decode(&payload); err != nil {
				t.Fatalf("couldnt decode payload")
			}
			expectedPayload := append(counter, []byte("data")...)
			if !reflect.DeepEqual(addr, expectedAddr) || !reflect.DeepEqual(payload, expectedPayload) {
				t.Fatalf("message lost or not expected")
			}
		}
		//ring should be empty at this point
		for i := int64(0); i < 2*ringCapacity; i++ {
			payload := backend.pendingMessages.Dequeue()
			if payload != nil {
				t.Fatalf("ring not empty")
			}
		}

	})
}

func TestTendermintMessage(t *testing.T) {
	_, backend := newBlockChain(1)

	// generate one msg
	data := []byte("data1")
	hash := types.RLPHash(data)
	msg := makeMsg(tendermintMsg, data)
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
	_, err := backend.HandleMsg(addr, msg)
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
			eventMux:    event.NewTypeMuxSilent(log.New("backend", "test", "id", 0)),
		}

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
