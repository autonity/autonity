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

package acn

import (
	"errors"
	"sync"

	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/acn/protocol"
	"github.com/autonity/autonity/p2p/enode"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/p2p"
)

var (
	// errPeerSetClosed is returned if a peer is attempted to be added or removed
	// from the peer set after it has been terminated.
	errPeerSetClosed = errors.New("peerset closed")

	// errPeerAlreadyRegistered is returned if a peer is attempted to be added
	// to the peer set, but one with the same id already exists.
	errPeerAlreadyRegistered = errors.New("peer already registered")

	// errPeerNotRegistered is returned if a peer is attempted to be removed from
	// a peer set, but no peer with the given id exists.
	errPeerNotRegistered = errors.New("peer not registered")
)

// PeerSet represents the collection of active peers currently participating in consensus
type peerSet struct {
	peers     map[common.Address]*protocol.Peer // Peers connected on the `acn` protocol by address
	peersByID map[enode.ID]*protocol.Peer       // Peers connected on the `acn` protocol by ID
	sync.RWMutex
	closed bool
}

// Voters creates a new peer set to track the active participants.
func newPeerSet() *peerSet {
	return &peerSet{
		peers:     make(map[common.Address]*protocol.Peer),
		peersByID: make(map[enode.ID]*protocol.Peer),
	}
}

// register injects a new `consensus` peer into the working set, or returns an error
// if the peer is already known.
func (ps *peerSet) register(peer *protocol.Peer) error {
	// Start tracking the new peer
	ps.Lock()
	defer ps.Unlock()

	if ps.closed {
		return errPeerSetClosed
	}
	if _, ok := ps.peers[peer.Address()]; ok {
		return errPeerAlreadyRegistered
	}
	ps.peers[peer.Address()] = peer
	ps.peersByID[peer.ID()] = peer
	return nil
}

// unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *peerSet) unregister(peer *protocol.Peer) error {
	ps.Lock()
	defer ps.Unlock()

	_, ok := ps.peers[peer.Address()]
	if !ok {
		return errPeerNotRegistered
	}
	delete(ps.peers, peer.Address())
	delete(ps.peersByID, peer.ID())
	return nil
}

// find retrieves the map of registered peer with the given map of ids.
func (ps *peerSet) find(targets []common.Address) map[common.Address]consensus.Peer {
	ps.RLock()
	defer ps.RUnlock()
	m := make(map[common.Address]consensus.Peer)
	for _, addr := range targets {
		if p, ok := ps.peers[addr]; ok {
			m[addr] = p
		}
	}
	return m
}

// peer retrieves the registered peer with the given id.
func (ps *peerSet) peer(address common.Address) (*protocol.Peer, bool) {
	ps.RLock()
	defer ps.RUnlock()

	p, ok := ps.peers[address]
	return p, ok
}

// peer retrieves the registered peer with the given id.
func (ps *peerSet) peerByID(id enode.ID) (*protocol.Peer, bool) {
	ps.RLock()
	defer ps.RUnlock()
	p, ok := ps.peersByID[id]
	return p, ok
}

// close disconnects all peers.
func (ps *peerSet) close() {
	ps.Lock()
	defer ps.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
