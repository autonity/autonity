// we're disabling race flag due to oom issues with travis CI
//go:build !race
// +build !race

package backend

import (
	"bytes"
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/p2p"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
)

func TestUnhandledMsgs(t *testing.T) {
	t.Run("core not running, unhandled messages are saved", func(t *testing.T) {
		blockchain, backend := newBlockChain(1)
		engine := blockchain.Engine().(consensus.BFT)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		broadcaster := consensus.NewMockBroadcaster(ctrl)
		backend.SetBroadcaster(broadcaster)
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
			if !reflect.DeepEqual(addr, expectedAddr) || !bytes.Equal(payload, expectedPayload) {
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockedPeer := consensus.NewMockPeer(ctrl)
		broadcaster := consensus.NewMockBroadcaster(ctrl)
		addressCache := fixsizecache.New[common.Hash, bool](1997, 10, 0, fixsizecache.HashKey[common.Hash])
		mockedPeer.EXPECT().Cache().Return(addressCache).AnyTimes()
		broadcaster.EXPECT().FindPeer(gomock.Any()).Return(mockedPeer, true).AnyTimes()
		backend.SetBroadcaster(broadcaster)

		dis := interfaces.NewMockEventDispatcher(ctrl)
		i := 0
		var received [ringCapacity]bool
		dis.EXPECT().Post(gomock.Any()).Do(func(eve events.MessageEvent) {
			message := eve.Message
			if message.R() != 1 || message.H() != 2 {
				t.Fatalf("message not expected")
			}
			i++
			received[message.Value().Big().Uint64()] = true
		}).AnyTimes()
		backend.evDispatcher = dis
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
		if err := backend.Start(context.Background()); err != nil {
			t.Fatalf("could not restart core")
		}

		backend.HandleUnhandledMsgs(context.Background())

		time.Sleep(time.Second)
		if i != ringCapacity {
			t.Fatalf("could not receiving events")
		}

		for _, msg := range received {
			if !msg {
				t.Fatalf("message lost")
			}
		}
	})
}
