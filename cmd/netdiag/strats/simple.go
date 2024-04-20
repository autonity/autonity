package strats

import (
	"math"
	rand2 "math/rand"
)

// *******************************
// ***** SIMPLE STRATEGY *********
// *******************************
// Explanations:
// Assume m = sqrt(n) where n is the network size. This mechanism works by splitting N into m groups of m  nodes.
// From each group we pick up a random node, we call it the group leader.
// We broadcast initially only to the group leaders our message. Upon reception, the group leader is
// responsible to broadcast to the other members of his group the message.
// For added redundancy, after this first phase, we can randomly select some other nodes for a second round.

type Simple struct {
	BaseStrategy
}

type ResultDisseminate struct {
	BaseResult
}

func (p *Simple) Execute(packetId uint64, data []byte, maxPeers int) error {
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
			// no suitable target found in the group to deal with
			// last group size
		}
		target.Send(SimpleCode, DisseminatePacket{packetId, uint64(p.State.Id), 1, data})
	}
	return nil
}

func (p *Simple) HandlePacket() {

}

func disseminationGroup(id int, peers []*Peer) []*Peer {
	// todo: create a special object for each propagation strategy to not overload state
	groupSize := int(math.Sqrt(float64(len(peers))))
	groupCount := groupSize
	if groupSize*groupCount < len(peers) {
		groupCount++
	}
	group := make([]*Peer, groupCount)
	myGroup := id % groupSize
	for i := range group {
		group[i] = peers[myGroup*groupSize+i]
	}
	return group // we are returning ourselves so caller be aware
}
