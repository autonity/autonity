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
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/eth/protocols/tm"
	"github.com/autonity/autonity/p2p/enode"
)

// tmHandler implements the tm.Backend interface to handle the various network
// packets that are sent as replies or broadcasts.
type tmHandler handler

func (h *tmHandler) Chain() *core.BlockChain { return h.chain }

// RunPeer is invoked when a peer joins on the `snap` protocol.
func (h *tmHandler) RunPeer(peer *tm.Peer, hand tm.Handler) error {
	return (*handler)(h).runConsensusPeer(peer, hand)
}

// PeerInfo retrieves all known `snap` information about a peer.
func (h *tmHandler) PeerInfo(id enode.ID) interface{} {
	//TODO
	return nil
}

// Handle is invoked from a peer's message handler when it receives a new remote
// message that the handler couldn't consume and serve itself.
func (h *tmHandler) Handle(peer *tm.Peer, packet tm.Packet) error {
	//TODO: what to do here, may be remove this
	return nil
}
func (h *tmHandler) TxPool() tm.TxPool { return h.txpool }
