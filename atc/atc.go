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

package atc

import (
	"context"
	"math"
	"sync"

	"github.com/autonity/autonity/eth"

	autonity "github.com/autonity/autonity"
	"github.com/autonity/autonity/atc/protocol"
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

type ATC struct {
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

func New(stack *node.Node, backend *eth.Ethereum, netID uint64) *ATC {
	atc := &ATC{
		peers:      newPeerSet(),
		chain:      backend.BlockChain(),
		networkID:  netID,
		forkFilter: forkid.NewFilter(backend.BlockChain()),
		server:     stack.ConsensusServer(),
		log:        log.New(),
		address:    crypto.PubkeyToAddress(stack.Config().NodeKey().PublicKey),
	}

	atc.server.MaxPeers = math.MaxInt
	stack.RegisterConsensusProtocols(atc.ConsensusProtocols())
	stack.RegisterLifecycle(atc)
	if handler, ok := atc.chain.Engine().(consensus.Handler); ok {
		handler.SetBroadcaster(atc)
	}
	// once p2p protocol handler is initialized, set it for accountability module for the off-chain accountability protocol.
	backend.FD().SetBroadcaster(atc)
	return atc
}

func (atc *ATC) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	atc.watchCommittee(ctx)
	protocol.StartENRUpdater(ctx, atc.chain, atc.server.LocalNode())
	atc.cancel = cancel
	return nil
}

func (atc *ATC) ConsensusProtocols() []p2p.Protocol {
	protos := protocol.MakeProtocols(atc, atc.networkID)
	return protos
}

func (atc *ATC) FindPeers(targets map[common.Address]struct{}) map[common.Address]autonity.Peer {
	return atc.peers.find(targets)
}

// runConsensusPeer registers a `consensus` peer into the consensus peerset and
// starts handling inbound messages.
func (atc *ATC) runConsensusPeer(peer *protocol.Peer, handler protocol.HandlerFunc) error {
	atc.wg.Add(1)
	defer atc.wg.Done()

	// Execute the Consensus handshake
	var (
		genesis = atc.chain.Genesis()
		head    = atc.chain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = atc.chain.GetTd(hash, number)
	)
	forkID := forkid.NewID(atc.chain.Config(), atc.chain.Genesis().Hash(), atc.chain.CurrentHeader().Number.Uint64())
	if err := peer.Handshake(atc.networkID, td, hash, genesis.Hash(), forkID, atc.forkFilter); err != nil {
		peer.Log().Debug("Consensus handshake failed", "err", err)
		return err
	}

	if err := atc.peers.registerPeer(peer); err != nil {
		peer.Log().Error("Snapshot extension registration failed", "err", err)
		return err
	}
	//TODO: checkpoint hash and required blocks check not done on consensus channel
	// Do we need it here
	defer atc.peers.unregisterPeer(peer.ID())
	return handler(peer)
}

func (atc *ATC) Stop() error {
	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to h.peers yet
	// will exit when they try to register.
	atc.cancel()
	atc.peers.close()
	atc.wg.Wait()

	return nil
}

func (atc *ATC) Chain() *core.BlockChain { return atc.chain }

// RunPeer is invoked when a peer joins on the `snap` protocol.
func (atc *ATC) RunPeer(peer *protocol.Peer, hand protocol.HandlerFunc) error {
	return atc.runConsensusPeer(peer, hand)
}

// PeerInfo retrieves all known `atc` information about a peer.
func (atc *ATC) PeerInfo(_ enode.ID) interface{} {
	//TODO
	return nil
}
