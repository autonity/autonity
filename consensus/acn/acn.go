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
	"context"
	"math"
	"sync"

	"github.com/autonity/autonity/consensus/acn/protocol"
	"github.com/autonity/autonity/eth"

	autonity "github.com/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/forkid"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/enode"
)

type ACN struct {
	networkID  uint64
	peers      *peerSet
	chain      *core.BlockChain
	wg         sync.WaitGroup
	forkFilter forkid.Filter // Fork ID filter, constant across the lifetime of the node
	server     *p2p.Server
	log        log.Logger
	address    common.Address
	cancel     context.CancelFunc
}

func New(stack *node.Node, backend *eth.Ethereum, netID uint64) *ACN {
	nodeKey, _ := stack.Config().AutonityKeys()
	acn := &ACN{
		peers:      newPeerSet(),
		chain:      backend.BlockChain(),
		networkID:  netID,
		forkFilter: forkid.NewFilter(backend.BlockChain()),
		server:     stack.ConsensusServer(),
		log:        log.New(),
		address:    crypto.PubkeyToAddress(nodeKey.PublicKey),
	}

	acn.server.MaxPeers = math.MaxInt
	stack.RegisterConsensusProtocols(acn.Protocols())
	stack.RegisterLifecycle(acn)
	if handler, ok := acn.chain.Engine().(consensus.Handler); ok {
		handler.SetBroadcaster(acn)
	}
	// once p2p protocol handler is initialized, set it for accountability module for the off-chain accountability protocol.
	backend.FD().SetBroadcaster(acn)
	return acn
}

func (acn *ACN) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	acn.watchCommittee(ctx)
	acn.cancel = cancel
	return nil
}

func (acn *ACN) Protocols() []p2p.Protocol {
	protos := protocol.MakeProtocols(acn, acn.networkID)
	return protos
}

func (acn *ACN) FindPeers(targets map[common.Address]struct{}) map[common.Address]autonity.Peer {
	return acn.peers.find(targets)
}

// runConsensusPeer registers a `consensus` peer into the consensus peerset and
// starts handling inbound messages.
func (acn *ACN) runConsensusPeer(peer *protocol.Peer, handler protocol.HandlerFunc) error {
	acn.wg.Add(1)
	defer acn.wg.Done()

	// Execute the Consensus handshake
	var (
		genesis = acn.chain.Genesis()
		head    = acn.chain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = acn.chain.GetTd(hash, number)
	)
	forkID := forkid.NewID(acn.chain.Config(), acn.chain.Genesis().Hash(), acn.chain.CurrentHeader().Number.Uint64())
	if err := peer.Handshake(acn.networkID, td, hash, genesis.Hash(), forkID, acn.forkFilter); err != nil {
		peer.Log().Debug("Consensus handshake failed", "err", err)
		return err
	}

	if err := acn.peers.register(peer); err != nil {
		peer.Log().Error("Snapshot extension registration failed", "err", err)
		return err
	}
	//TODO: checkpoint hash and required blocks check not done on consensus channel
	// Do we need it here
	defer acn.peers.unregister(peer.ID())
	return handler(peer)
}

func (acn *ACN) Stop() error {
	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to h.peers yet
	// will exit when they try to register.
	acn.cancel()
	acn.peers.close()
	acn.wg.Wait()

	return nil
}

func (acn *ACN) Chain() *core.BlockChain { return acn.chain }

// RunPeer is invoked when a peer joins on the `snap` protocol.
func (acn *ACN) RunPeer(peer *protocol.Peer, hand protocol.HandlerFunc) error {
	return acn.runConsensusPeer(peer, hand)
}

// PeerInfo retrieves all known `acn` information about a peer.
func (acn *ACN) PeerInfo(_ enode.ID) interface{} {
	//TODO
	return nil
}
