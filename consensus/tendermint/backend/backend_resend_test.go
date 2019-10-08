package backend

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/log"
	lru "github.com/hashicorp/golang-lru"
)

func TestSendToConnectedPeers(t *testing.T) {
	t.Run("no peers, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errConnectedPeers = make([]common.Address, 0, 0)

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(make(map[common.Address]struct{}))

		b := &Backend{}
		b.SetBroadcaster(broadcaster)

		res := b.sendToConnectedPeers(context.Background(), messageToPeers{})

		if !reflect.DeepEqual(errConnectedPeers, res) {
			t.Fatalf("Expected %v, got %v", errConnectedPeers, res)
		}
	})

	t.Run("peers found, messages sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errConnectedPeers = make([]common.Address, 0, 0)

		peerAddr1 := common.HexToAddress("0x0123456789")
		msg := messageToPeers{
			msg: message{},
			peers: []common.Address{
				peerAddr1,
			},
			startTime: time.Time{},
			lastTry:   time.Time{},
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		peer1Mock := consensus.NewMockPeer(ctrl)
		peer1Mock.EXPECT().Send(uint64(tendermintMsg), msg.msg.payload)

		peers := make(map[common.Address]consensus.Peer)
		peers[peerAddr1] = peer1Mock

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			recentMessages: recentMessages,
		}
		b.SetBroadcaster(broadcaster)

		res := b.sendToConnectedPeers(context.Background(), msg)

		if !reflect.DeepEqual(errConnectedPeers, res) {
			t.Fatalf("Expected %v, got %v", errConnectedPeers, res)
		}
	})

	t.Run("peers found, messages sent, don't send again", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errConnectedPeers = make([]common.Address, 0, 0)

		peerAddr1 := common.HexToAddress("0x0123456789")
		msg := messageToPeers{
			msg: message{
				hash: common.HexToHash("0x0123456789"),
			},
			peers: []common.Address{
				peerAddr1,
			},
			startTime: time.Time{},
			lastTry:   time.Time{},
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		peers := make(map[common.Address]consensus.Peer)
		peers[peerAddr1] = consensus.NewMockPeer(ctrl)

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		perrCache, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}
		perrCache.Add(msg.msg.hash, msg.msg.payload)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		recentMessages.Add(peerAddr1, perrCache)

		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			recentMessages: recentMessages,
		}
		b.SetBroadcaster(broadcaster)

		res := b.sendToConnectedPeers(context.Background(), msg)

		if !reflect.DeepEqual(errConnectedPeers, res) {
			t.Fatalf("Expected %v, got %v", errConnectedPeers, res)
		}
	})

	t.Run("peers found, messages sent on second retry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errConnectedPeers = make([]common.Address, 0, 0)

		peerAddr1 := common.HexToAddress("0x0123456789")
		msg := messageToPeers{
			msg: message{},
			peers: []common.Address{
				peerAddr1,
			},
			startTime: time.Time{},
			lastTry:   time.Time{},
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		peer1Mock := consensus.NewMockPeer(ctrl)
		peer1Mock.EXPECT().Send(uint64(tendermintMsg), msg.msg.payload).Return(errors.New("some error"))
		peer1Mock.EXPECT().Send(uint64(tendermintMsg), msg.msg.payload).Return(nil)

		peers := make(map[common.Address]consensus.Peer)
		peers[peerAddr1] = peer1Mock

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			recentMessages: recentMessages,
		}
		b.SetBroadcaster(broadcaster)

		res := b.sendToConnectedPeers(context.Background(), msg)

		if !reflect.DeepEqual(errConnectedPeers, res) {
			t.Fatalf("Expected %v, got %v", errConnectedPeers, res)
		}
	})

	t.Run("context cancelled, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		peerAddr1 := common.HexToAddress("0x0123456789")
		var errConnectedPeers = make([]common.Address, 0, 1)
		errConnectedPeers = append(errConnectedPeers, peerAddr1)

		msg := messageToPeers{
			msg: message{},
			peers: []common.Address{
				peerAddr1,
			},
			startTime: time.Time{},
			lastTry:   time.Time{},
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		peers := make(map[common.Address]consensus.Peer)
		peers[peerAddr1] = consensus.NewMockPeer(ctrl)

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			recentMessages: recentMessages,
		}
		b.SetBroadcaster(broadcaster)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		res := b.sendToConnectedPeers(ctx, msg)

		if !reflect.DeepEqual(errConnectedPeers, res) {
			t.Fatalf("Expected %v, got %v", errConnectedPeers, res)
		}
	})
}

func TestTrySend(t *testing.T) {
	t.Run("TTL expired, nothing done", func(t *testing.T) {
		b := &Backend{
			logger: log.New("backend", "test", "id", 0),
		}

		b.trySend(context.Background(), messageToPeers{})
	})

	t.Run("retryInterval exceeded, message resent", func(t *testing.T) {
		msg := messageToPeers{
			startTime: time.Now().Add(-(retryInterval + 1) * time.Millisecond),
			lastTry:   time.Now(),
		}

		b := &Backend{
			logger: log.New("backend", "test", "id", 0),
			resend: make(chan messageToPeers, 1),
		}

		b.trySend(context.Background(), msg)

		m := <-b.resend
		if !reflect.DeepEqual(msg, m) {
			t.Fatalf("Expected %v, got %v", msg, m)
		}
	})

	t.Run("context done, message not resent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		peerAddr1 := common.HexToAddress("0x0123456789")
		var errConnectedPeers = make([]common.Address, 0, 1)
		errConnectedPeers = append(errConnectedPeers, peerAddr1)

		msg := messageToPeers{
			msg: message{},
			peers: []common.Address{
				peerAddr1,
			},
			startTime: time.Now(),
			lastTry:   time.Now(),
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		peer1Mock := consensus.NewMockPeer(ctrl)

		peers := make(map[common.Address]consensus.Peer)
		peers[peerAddr1] = peer1Mock

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			recentMessages: recentMessages,
		}
		b.SetBroadcaster(broadcaster)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		b.trySend(ctx, msg)
	})
}
