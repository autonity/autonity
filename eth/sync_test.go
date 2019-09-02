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

package eth

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/p2p"
)

// Tests that fast sync gets disabled as soon as a real block is successfully
// imported into the blockchain.
func TestFastSyncDisabling(t *testing.T) {
	// Create a pristine protocol manager, check that fast sync is left enabled
	p2pPeerEmpty := newTestP2PPeer("peerEmpty")
	p2pPeerFull := newTestP2PPeer("peerFull")

	pmEmpty, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil, []string{p2pPeerEmpty.Info().Enode, p2pPeerFull.Info().Enode})

	if atomic.LoadUint32(&pmEmpty.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}

	// Create a full protocol manager, check that fast sync gets disabled
	pmFull, _ := newTestProtocolManagerMust(t, downloader.FastSync, 1024, nil, nil, []string{p2pPeerEmpty.Info().Enode, p2pPeerFull.Info().Enode})

	if atomic.LoadUint32(&pmFull.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}
	// Sync up the two peers
	io1, io2 := p2p.MsgPipe()

	go pmFull.handle(pmFull.newPeer(63, p2pPeerEmpty, io2))
	go pmEmpty.handle(pmEmpty.newPeer(63, p2pPeerFull, io1))

	time.Sleep(250 * time.Millisecond)
	pmEmpty.synchronise(pmEmpty.peers.BestPeer())

	// Check that fast sync was disabled
	if atomic.LoadUint32(&pmEmpty.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}
}

// Tests that fast sync gets disabled as soon as a real block is successfully
// imported into the blockchain.
func TestFastSyncDisablingMany(t *testing.T) {
	// Create a pristine protocol manager, check that fast sync is left enabled
	p2pPeerEmpty := newTestP2PPeer("peerEmpty")
	p2pPeerEmpty1 := newTestP2PPeer("peerEmpty1")
	p2pPeerFull := newTestP2PPeer("peerFull")

	pmEmpty, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil, []string{p2pPeerEmpty.Info().Enode, p2pPeerFull.Info().Enode, p2pPeerEmpty1.Info().Enode})

	if atomic.LoadUint32(&pmEmpty.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}

	pmEmpty1, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil, []string{p2pPeerEmpty.Info().Enode, p2pPeerFull.Info().Enode, p2pPeerEmpty1.Info().Enode})

	if atomic.LoadUint32(&pmEmpty1.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain 1")
	}

	// Create a full protocol manager, check that fast sync gets disabled
	pmFull, _ := newTestProtocolManagerMust(t, downloader.FastSync, 1024, nil, nil, []string{p2pPeerEmpty.Info().Enode, p2pPeerFull.Info().Enode, p2pPeerEmpty1.Info().Enode})

	if atomic.LoadUint32(&pmFull.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}
	// Sync up the two peers
	// p2pPeerEmpty <-> p2pPeerFull
	io1, io2 := p2p.MsgPipe()
	go func() {
		_ = pmFull.handle(pmFull.newPeer(63, p2pPeerEmpty, io2))
	}()
	go func() {
		_ = pmEmpty.handle(pmEmpty.newPeer(63, p2pPeerFull, io1))
	}()

	// p2pPeerEmpty1 <-> p2pPeerFull
	io3, io4 := p2p.MsgPipe()
	go func() {
		_ = pmFull.handle(pmFull.newPeer(63, p2pPeerEmpty1, io3))
	}()
	go func() {
		_ = pmEmpty1.handle(pmEmpty.newPeer(63, p2pPeerFull, io4))
	}()

	// p2pPeerEmpty1 <-> p2pPeerEmpty
	io5, io6 := p2p.MsgPipe()
	go func() {
		_ = pmEmpty.handle(pmFull.newPeer(63, p2pPeerEmpty1, io5))
	}()
	go func() {
		_ = pmEmpty1.handle(pmEmpty.newPeer(63, p2pPeerEmpty, io6))
	}()

	time.Sleep(250 * time.Millisecond)
	pmEmpty.synchronise(pmEmpty.peers.BestPeer())

	// Check that fast sync was disabled
	if atomic.LoadUint32(&pmEmpty.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}

	pmEmpty1.synchronise(pmEmpty1.peers.BestPeer())
	// Check that fast sync was disabled
	if atomic.LoadUint32(&pmEmpty1.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}
}
