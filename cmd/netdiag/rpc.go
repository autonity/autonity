package main

import (
	"errors"

	probing "github.com/prometheus-community/pro-bing"
)

var (
	errInvalidRpcArg = errors.New("invalid RPC argument")
)

type ArgTarget struct {
	target int
}
type ArgEmpty struct {
}

// P2P represents p2p operation commands
type P2P struct {
	engine *engine
}

func (p *P2P) ConnectedPeers(_ ArgEmpty, reply *int) error {
	connected := 0
	for _, p := range p.engine.peers {
		if p.connected {
			connected++
		}
	}
	*reply = connected
	return nil
}

func (p *P2P) PingIcmp(args *ArgTarget, reply *probing.Statistics) error {
	if args.target < 0 || args.target >= len(p.engine.peers) {
		return errInvalidRpcArg
	}
	*reply = *<-pingIcmp(p.engine.peers[args.target].address)
	return nil
}

func (p *P2P) PingIcmpAll(args *ArgTarget, reply *probing.Statistics) error {
	for i := 0; i < len(p.engine.peers); i++ {
		*reply = *<-pingIcmp(p.engine.peers[args.target].address)
	}
	return nil
}

func (p *P2P) Ping(args *ArgTarget, reply *probing.Statistics) error {
	if args.target < 0 || args.target >= len(p.engine.peers) {
		return errInvalidRpcArg
	}
	*reply = *<-pingIcmp(p.engine.peers[args.target].address)
	return nil
}
