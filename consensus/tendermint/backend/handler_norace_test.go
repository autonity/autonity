// we're disabling race flag due to oom issues with travis CI
//go:build !race
// +build !race

package backend

import (
	"bytes"
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/p2p"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/events"

	"github.com/autonity/autonity/common"
)

func TestUnhandledMsgs(t *testing.T) {
	t.Run("core not running, unhandled messages are saved", func(t *testing.T) {
		blockchain, backend := newBlockChain(1)
		engine := blockchain.Engine().(consensus.BFT)
		// we close the engine for enabling cache storing
		if err := engine.Close(); err != nil {
			t.Fatalf("can't stop the engine")
		}
		//we generate a bunch of messages overflowing max capacity
		for i := int64(0); i < 2*ringCapacity; i++ {
			counter := big.NewInt(i).Bytes()
			msg := makeMsg(PrevoteNetworkMsg, append(counter, []byte("data")...))
			addr := common.BytesToAddress(append(counter, []byte("addr")...))
			if result, err := backend.HandleMsg(addr, msg, nil); !result || err != nil {
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
			if savedMsg.(UnhandledMsg).msg.Code != PrevoteNetworkMsg {
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

	t.Run("core running, unhandled messages are processed", func(t *testing.T) {
		blockchain, backend := newBlockChain(1)
		engine := blockchain.Engine().(consensus.BFT)
		// we close the engine for enabling cache storing
		if err := engine.Close(); err != nil {
			t.Fatalf("can't stop the engine")
		}

		for i := int64(0); i < ringCapacity; i++ {
			counter := big.NewInt(i).Bytes()
			vote := message.NewPrevote(1, 2, common.BigToHash(big.NewInt(i)), dummySigner)
			msg := p2p.Msg{Code: PrevoteNetworkMsg, Size: uint32(len(vote.Payload())), Payload: bytes.NewReader(vote.Payload())}
			addr := common.BytesToAddress(append(counter, []byte("addr")...))
			if result, err := backend.HandleMsg(addr, msg, nil); !result || err != nil {
				t.Fatalf("handleMsg should have been successful")
			}
		}
		sub := backend.eventMux.Subscribe(events.MessageEvent{})
		if err := backend.Start(context.Background()); err != nil {
			t.Fatalf("could not restart core")
		}
		backend.HandleUnhandledMsgs(context.Background())
		timer := time.NewTimer(time.Second)
		i := 0
		var received [ringCapacity]bool
		// events can come out of order so we track them using an array.
	LOOP:
		for {
			select {
			case eve := <-sub.Chan():
				message := eve.Data.(events.MessageEvent).Message
				if message.R() != 1 || message.H() != 2 {
					t.Fatalf("message not expected")
				}
				i++
				received[message.Value().Big().Uint64()] = true

			case <-timer.C:
				if i == ringCapacity {
					break LOOP
				}
				t.Fatalf("timeout receiving events")
			}
		}

		for _, msg := range received {
			if !msg {
				t.Fatalf("message lost")
			}
		}
	})
}
