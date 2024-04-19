package strats

import (
	"fmt"
	"math"
	rand2 "math/rand"
	"sort"
	"strings"
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

type PeersAccessor interface {
	Peer(int i)
	SendToPeer(i int)
}

type Simple struct {
	// Disseminator
	// Peer retriever
	//
}

type ResultDisseminate struct {
	ResultBase
}

func (p *Simple) Execute(data []byte, state *State, maxPeers uint64) (uint64, error) {
	groupSize := int(math.Sqrt(float64(len(p.engine.peers))))
	groupCount := groupSize
	if groupSize*groupCount < len(p.engine.peers) {
		groupCount++
	}
	for i := 0; i < groupCount; i++ {
		var (
			target *Peer
			peerId int
		)
		for target == nil {
			l := rand2.Intn(groupSize)
			peerId = i*groupSize + l
			target = p.engine.peers[peerId]
			// edge cases:
			// no suitable target found in the group to deal with
			// last group size
		}
		fmt.Println("TARGET FOUND", target.ip, "id", peerId)
		target.sendDisseminate(packetId, buff, uint64(p.engine.id), 1)
	}
	return nil
}

func (r *ResultDisseminate) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Disseminate Results \n")
	var results []*IndividualDisseminateResult
	for i, res := range r.IndividualResults {
		if res.ErrorTimeout {
			continue
		}
		results = append(results, r.IndividualResults[i])
		fmt.Fprintf(&builder, "Peer %d Duration: %s Hops: %d Relay: %d\n", i, res.ReceptionTime.Sub(r.StartTime), res.Hop, res.Relay)
	}
	sort.Slice(results, func(a, b int) bool {
		return results[a].ReceptionTime.Before(results[b].ReceptionTime)
	})
	n := len(results)
	fmt.Fprintf(&builder, "min: %s, median:%s 2/3rd:%s max: %s\n", results[0].ReceptionTime.Sub(r.StartTime), results[n/2].ReceptionTime.Sub(r.StartTime), results[(2*n)/3].ReceptionTime.Sub(r.StartTime), results[n-1].ReceptionTime.Sub(r.StartTime))

	return builder.String()
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
