package main

import (
	"crypto/ecdsa"
	"sync"

	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/enode"
)

type Engine struct {
	config
	id     int
	server *p2p.Server

	peers  []*peer // nil never connected, can do probably cleaner
	enodes []*enode.Node
	sync.RWMutex
}

func newEngine(cfg config, key *ecdsa.PrivateKey) *Engine {
	e := new(Engine)
	runner := func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
		node, err := e.addPeer(peer, rw)
		if err != nil {
			return err
		}
		for {
			msg, err := rw.ReadMsg()
			if err != nil {
				return err
			}
			handler := protocolHandlers[msg.Code]
			if err = handler(node, msg.Payload); err != nil {
				return err
			}
		}
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
			Length:  5,
			Run:     runner,
		}},
		ListenAddr:      "0.0.0.0:20203",
		NAT:             nil,
		Dialer:          nil, // nil is default TCP, have UDP supported at one point
		NoDial:          false,
		EnableMsgEvents: false,
		Logger:          log.Root(),
	}
	e.config = cfg
	e.server = &p2p.Server{Net: p2p.Consensus, Config: p2pConfig}

	enodesToResolve := make([]string, len(e.config.Nodes))
	for i := range enodesToResolve {
		enodesToResolve[i] = e.config.Nodes[i].Enode
	}
	e.enodes = types.NewNodes(enodesToResolve, true).List
	e.peers = make([]*peer, len(e.enodes))
	return e
}

func (e *Engine) start() error {
	// attempt to connect to everyone. Use our logic.
	if err := e.server.Start(); err != nil {
		log.Error("error starting p2p server", "err", err)
		return err
	}
	e.server.UpdateConsensusEnodes(e.enodes, e.enodes)

	/*
		// We need first a  good view on local time for remote peers.
		// Basic responsiveness -
		// establish first a baseline with regular IMCP ping messages Point to Point Sequentially
		e.peersDo(func(p *peer) {
			pingIcmp(p)
		})

		// Now wait for the connection to get established.
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			if e.peerCount() == len(enodes.List) {
				ticker.Stop()
				break
			}
		}

		// EACH DIAGNOSYS PHASE REQUIRE TEST, HYPOTHESIS, CONFRONTATION.

		// NOW do it at devp2p protocol level still sequentially
		// CONFIRM PTP RTT and get estimation of time delta

		e.peersDo(func() {
			req := sendPing(n)
			// wait for response before moving on !
		})

		//------------- BANDWITH TEST
		// HERE we are testing for a long stream of data - reception time

		// --- INDIVIDUAL NODE
		e.peersDo(func() {
			sendData(n)
			// wait for response before moving on !
		})
		// --- P2P peers | increment by 5 every run the number of peers receiving data.

		e.peersDoParrallel(func() {
			sendData(n)
			// wait for response before moving on !
		})
		//------------- LATENCY TEST

		// Now test sending one small packet (less than MTU)

		// 10kb 50kb 200kb . CONFIRM AGAINST TIME
	*/
	return nil
}

func (e *Engine) addPeer(node *p2p.Peer, rw p2p.MsgReadWriter) (*peer, error) {
	e.Lock()
	defer e.Unlock()
	p := &peer{
		Peer:          node,
		MsgReadWriter: rw,
		connected:     true,
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

func (e *Engine) peersDo(f func(p *peer)) {
	e.Lock()
	defer e.Unlock()
	for _, p := range e.peers {
		if p != nil {
			f(p)
		}
	}
}
