package strats

import (
	"math"
	rand2 "math/rand"
)

// *********************************************
// ***** SIMPLE DISSEMINATION STRATEGY *********
// *********************************************
// Explanations:
// Assume m = sqrt(n) where n is the network size. This mechanism works by splitting N into m groups of m  nodes.
// From each group we pick up a random node, we call it the group leader.
// We broadcast initially only to the group leaders our message. Upon reception, the group leader is
// responsible to broadcast to the other members of his group the message.
// For added redundancy, after this first phase, we can randomly select some other nodes for a second round.

type Simple struct {
	BaseStrategy
}

func init() {
	registerStrategy("Simple Dissemination Tree", func(base BaseStrategy) Strategy {
		return &Simple{base}
	})
}

func (p *Simple) Execute(packetId uint64, data []byte, maxPeers int) error {
	// SENDER Should be excluded from maxPeers !
	groupSize := int(math.Sqrt(float64(maxPeers)))
	groupCount := groupSize
	if groupSize*groupCount < maxPeers {
		groupCount++
	}
	for i := 0; i < groupCount; i++ {
		var (
			target Peer
			peerId int
		)
		for target == nil {
			l := rand2.Intn(groupSize)
			peerId = i*groupSize + l
			target = p.Peers(peerId)
			// edge cases:
			// - no suitable target found in the group to deal with
			// - last group size
		}
		err := target.DisseminateRequest(p.Code, packetId, 1, uint64(p.State.Id), uint64(maxPeers), data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Simple) HandlePacket(requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data any) error {
	if hop == 1 {
		// need to disseminate in the group
		allPeers := make([]Peer, maxPeers)
		for i := range allPeers {
			allPeers[i] = p.Peers(i)
		}
		group := disseminationGroup(int(p.State.Id), allPeers)
		for i := range group {
			if group[i] == nil {
				continue
			}
			err := group[i].DisseminateRequest(p.Code, requestId, 0, originalSender, maxPeers, data)
			if err != nil {
				return err
			}
		}
	}
	if hop == 0 {
		// todo: include random peer selection logic - maybe set it as a parameter?
	}
	return nil
}

func disseminationGroup(id int, peers []Peer) []Peer {
	groupSize := int(math.Sqrt(float64(len(peers))))
	groupCount := groupSize
	if groupSize*groupCount < len(peers) {
		groupCount++
	}
	group := make([]Peer, groupCount)
	myGroup := id % groupSize
	for i := range group {
		group[i] = peers[myGroup*groupSize+i]
	}
	return group // we are returning ourselves so caller be aware
}
