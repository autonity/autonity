// Copyright 2020 The go-ethereum Authors
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
	autonity "github.com/autonity/autonity"
	"github.com/autonity/autonity/eth/protocols/tm"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/p2p"
)

// ConsensusPeerSet represents the collection of active peers currently participating in
// the `eth` protocol, with or without the `snap` extension.
type consensusPeerSet struct {
	peers  map[string]*consensusPeer // Peers connected on the `eth` protocol
	lock   sync.RWMutex
	closed bool
}

// Voters creates a new peer set to track the active participants.
func newConsensusPeerSet() *consensusPeerSet {
	return &consensusPeerSet{
		peers: make(map[string]*consensusPeer),
	}
}

// registerPeer injects a new `eth` peer into the working set, or returns an error
// if the peer is already known.
func (ps *consensusPeerSet) registerPeer(peer *tm.Peer) error {
	// Start tracking the new peer
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errPeerSetClosed
	}
	id := peer.ID()
	if _, ok := ps.peers[id]; ok {
		return errPeerAlreadyRegistered
	}
	ps.peers[id] = &consensusPeer{Peer: peer}
	return nil
}

// unregisterPeer removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *consensusPeerSet) unregisterPeer(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.peers[id]
	if !ok {
		return errPeerNotRegistered
	}
	delete(ps.peers, id)
	return nil
}

// peer retrieves the registered peer with the given id.
func (ps *consensusPeerSet) peer(id string) *consensusPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}
func (ps *consensusPeerSet) findPeers(targets map[common.Address]struct{}) map[common.Address]autonity.Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	m := make(map[common.Address]autonity.Peer)
	for _, p := range ps.peers {
		addr := p.Address()
		if _, ok := targets[addr]; ok {
			m[addr] = p
		}
	}
	return m
}

// peersWithoutBlock retrieves a list of peers that do not have a given block in
// their set of known hashes, so it might be propagated to them.
func (ps *consensusPeerSet) peersWithoutBlock(hash common.Hash) []*consensusPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*consensusPeer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.KnownBlock(hash) {
			list = append(list, p)
		}
	}
	return list
}

// len returns if the current number of `eth` peers in the set. Since the `snap`
// peers are tied to the existence of an `eth` connection, that will always be a
// subset of `eth`.
func (ps *consensusPeerSet) len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// close disconnects all peers.
func (ps *consensusPeerSet) close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
