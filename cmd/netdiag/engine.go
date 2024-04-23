package main

import (
	"crypto/ecdsa"
	"sync"

	"github.com/autonity/autonity/cmd/netdiag/strats"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/enode"
)

type Engine struct {
	config
	id     int
	server *p2p.Server

	peers  []*Peer // nil never connected, can do probably cleaner
	enodes []*enode.Node

	state      *strats.State
	strategies []strats.Strategy

	sync.RWMutex
}

func newEngine(cfg config, id int, key *ecdsa.PrivateKey, networkMode string) *Engine {
	e := &Engine{
		state: strats.NewState(uint64(id)),
	}
	for _, s := range strats.StrategyRegistry {
		e.strategies = append(e.strategies, s.Constructor(e.peer, e.state))
	}

	runner := func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
		node, err := e.addPeer(peer, rw)
		defer log.Debug("Reading loop broken")
		if err != nil {
			return err
		}
		for {
			msg, err := rw.ReadMsg()
			if err != nil {
				log.Error("Reading loop fatal error", "error", err)
				return err
			}
			handler := protocolHandlers[msg.Code]
			if err = handler(e, node, msg.Payload); err != nil {
				log.Debug("Peer handler error", "id", peer.String(), "error", err)
				return err
			}
		}
	}
	transport := p2p.TCP
	if networkMode == "udp" {
		transport = p2p.UDP
	}
	p2pConfig := p2p.Config{
		PrivateKey:      key,
		MaxPeers:        1000,
		MaxPendingPeers: 25,
		DialRatio:       0,
		NoDiscovery:     true,
		Name:            "diag",
		NodeDatabase:    "", // use memory
		Protocols: []p2p.Protocol{{
			Name:    "diag",
			Version: 1,
			Length:  ProtocolMessages + 1,
			Run:     runner,
		}},
		ListenAddr:      "0.0.0.0:20203",
		NAT:             nil,
		NoDial:          false,
		EnableMsgEvents: false,
		Logger:          log.Root(),
	}
	e.config = cfg
	e.server = &p2p.Server{Net: p2p.Consensus, Config: p2pConfig, Transport: transport}

	enodesToResolve := make([]string, len(e.config.Nodes))
	for i := range enodesToResolve {
		enodesToResolve[i] = e.config.Nodes[i].Enode
	}
	e.enodes = types.NewNodes(enodesToResolve, true).List
	e.peers = make([]*Peer, len(e.enodes))
	e.id = id
	return e
}

func (e *Engine) start() error {
	// attempt to connect to everyone. Use our logic.
	if err := e.server.Start(); err != nil {
		log.Error("error starting p2p server", "err", err)
		return err
	}
	e.server.UpdateConsensusEnodes(e.enodes, e.enodes)
	return nil
}

func (e *Engine) addPeer(node *p2p.Peer, rw p2p.MsgReadWriter) (*Peer, error) {
	e.Lock()
	defer e.Unlock()
	p := &Peer{
		Peer:          node,
		MsgReadWriter: rw,
		connected:     true,
		requests:      make(map[uint64]chan any),
	}
	for i := 0; i < len(e.config.Nodes); i++ {
		if e.enodes[i].ID() == node.ID() {
			p.ip = e.enodes[i].IP().String()
			e.peers[i] = p
			break
		}
	}

	return p, nil
}

func (e *Engine) peerCount() int {
	e.Lock()
	defer e.Unlock()
	return len(e.peers)
}

func (e *Engine) peer(i int) strats.Peer {
	return e.peers[i]
}

func (e *Engine) peerToId(peer *Peer) int {
	for i := range e.peers {
		if e.peers[i] == nil {
			continue
		}
		if e.peers[i].ID() == peer.ID() {
			return i
		}
	}
	return 0
}
