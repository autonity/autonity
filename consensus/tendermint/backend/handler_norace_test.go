// we're disabling race flag due to oom issues with travis CI
//go:build !race
// +build !race

package backend

import (
	"context"
	message "github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/stretchr/testify/require"
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
			msg := makeMsg(TendermintMsgProposal, append(counter, []byte("data")...))
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
			if savedMsg.(UnhandledMsg).msg.Code != TendermintMsgProposal {
				t.Fatalf("wrong msg code")
			}
			var payload []byte
			var decodedMsg = message.Message{
				ConsensusMsg: &message.Proposal{},
			}
			if err := savedMsg.(UnhandledMsg).msg.Decode(&decodedMsg); err != nil {
				t.Fatalf("couldnt decode payload")
			}

			payload = decodedMsg.Payload
			expectedPayload := append(counter, []byte("data")...)

			require.Equal(t, expectedAddr, addr)
			require.Equal(t, expectedPayload, payload)

			//if !reflect.DeepEqual(addr, expectedAddr) || !reflect.DeepEqual(payload, expectedPayload) {
			//	t.Fatalf("message lost or not expected")
			//}
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

			consensusMsg := &message.Vote{
				Height: big.NewInt(2137 + i),
			}

			msg := makeMsgVote(TendermintMsgVote, append(counter, []byte("data")...), consensusMsg)

			addr := common.BytesToAddress(append(counter, []byte("addr")...))
			if result, err := backend.HandleMsg(addr, msg, nil); !result || err != nil {
				t.Fatalf("handleMsg should have been successful")
			}
		}
		sub := backend.eventMux.Subscribe(events.NewMessageEvent{})
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
				payload := eve.Data.(events.NewMessageEvent).Message.Payload
				if !reflect.DeepEqual(payload[len(payload)-4:], []byte("data")) {
					t.Fatalf("message not expected")
				}
				i++
				received[new(big.Int).SetBytes(payload[:len(payload)-4]).Uint64()] = true

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
