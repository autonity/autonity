package core

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
	Config

	Id     int
	server *p2p.Server

	Peers  []*Peer // nil never connected, can do probably cleaner
	Enodes []*enode.Node

	State      *strats.State
	Strategies []strats.Strategy

	sync.RWMutex
}

func NewEngine(cfg Config, id int, key *ecdsa.PrivateKey, networkMode string) *Engine {
	e := &Engine{
		State: strats.NewState(uint64(id), len(cfg.Nodes)),
	}
	for _, s := range strats.StrategyRegistry {
		e.Strategies = append(e.Strategies, s.Constructor(e.peer, e.State))
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
			handler, ok := protocolHandlers[msg.Code]
			if !ok {
				return errInvalidMsgCode
			}
			if err = handler(e, node, msg.Payload); err != nil {
				log.Debug("Peer handler error", "id", peer.String(), "error", err)
				return err
			}
		}
	}
	transport := p2p.TCP
	if networkMode == "udp" {
		transport = p2p.UDP
	} else if networkMode == "quic" {
		transport = p2p.QUIC
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
	e.Config = cfg
	e.server = &p2p.Server{Net: p2p.Consensus, Config: p2pConfig, Transport: transport}

	enodesToResolve := make([]string, len(e.Config.Nodes))
	for i := range enodesToResolve {
		enodesToResolve[i] = e.Config.Nodes[i].Enode
	}
	e.Enodes = types.NewNodes(enodesToResolve, true).List
	e.Peers = make([]*Peer, len(e.Enodes))
	e.Id = id
	return e
}

func (e *Engine) Start() error {
	// attempt to connect to everyone. Use our logic.
	if err := e.server.Start(); err != nil {
		log.Error("error starting p2p server", "err", err)
		return err
	}
	e.server.UpdateConsensusEnodes(e.Enodes, e.Enodes)
	return nil
}

func (e *Engine) addPeer(node *p2p.Peer, rw p2p.MsgReadWriter) (*Peer, error) {
	e.Lock()
	defer e.Unlock()
	p := &Peer{
		Peer:          node,
		MsgReadWriter: rw,
		Connected:     true,
		requests:      make(map[uint64]chan any),
	}
	for i := 0; i < len(e.Config.Nodes); i++ {
		if e.Enodes[i].ID() == node.ID() {
			p.Ip = e.Enodes[i].IP().String()
			p.id = i
			e.Peers[i] = p
			break
		}
	}

	return p, nil
}

func (e *Engine) peerCount() int {
	e.Lock()
	defer e.Unlock()
	return len(e.Peers)
}

func (e *Engine) peer(i int) strats.Peer {
	if i >= len(e.Peers) {
		return nil
	}
	if e.Peers[i] == nil {
		return nil
	}
	return e.Peers[i]
}
