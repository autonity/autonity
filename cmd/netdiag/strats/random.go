package strats

import (
	rand2 "math/rand"
)

// ********** RANDOM DISSEMINATION STRATEGY *********
// Explanations:
// Original sender and recipients chose a fixed random number of peers to
// propagate the block:
// Limited version: Recipients of Hop 0 do not propagate.
// Full version: Recipients of Hop 0 do propagate.
// Todo: find analytic probabilistic formulation.

type Random struct {
	BaseStrategy
	RandomPC int
	Hop0     bool
}

func init() {
	registerStrategy("10% Random Broadcast - 10% Random Relay - Single Hop", func(base BaseStrategy) Strategy {
		return &Random{base, 10, false}
	})
	registerStrategy("25% Random Broadcast - 25% Random Relay - Single Hop", func(base BaseStrategy) Strategy {
		return &Random{base, 25, false}
	})
	registerStrategy("50% Random Broadcast - 50% Random Relay - Single Hop", func(base BaseStrategy) Strategy {
		return &Random{base, 50, false}
	})
	registerStrategy("10% Random Broadcast - 10% Random Relay - Multi Hops", func(base BaseStrategy) Strategy {
		return &Random{base, 10, true}
	})
	registerStrategy("25% Random Broadcast - 25% Random Relay - Multi Hops", func(base BaseStrategy) Strategy {
		return &Random{base, 25, true}
	})
	registerStrategy("50% Random Broadcast - 50% Random Relay - Multi Hops", func(base BaseStrategy) Strategy {
		return &Random{base, 50, true}
	})
}

func (p *Random) Execute(packetId uint64, data []byte, maxPeers int) error {
	return p.randomDissemination(p.RandomPC, packetId, data, maxPeers, uint64(p.State.Id), 1)
}

func (p *BaseStrategy) randomDissemination(randomPc int, packetId uint64, data []byte, maxPeers int, originalSender uint64, hop int) error {
	sent := map[int]struct{}{}
	numRecipients := (maxPeers * randomPc) / 100
	for i := 0; i < numRecipients; i++ {
		var (
			target Peer
			peerId int
		)
		for target == nil {
			peerId = rand2.Intn(maxPeers)
			if _, already := sent[peerId]; already {
				continue
			}
			target = p.Peers(peerId)
		}
		sent[peerId] = struct{}{}
		// TODO : test async!
		if err := target.DisseminateRequest(p.Code, packetId, uint8(hop), originalSender, uint64(maxPeers), data); err != nil {
			return err
		}
	}
	return nil
}

func (p *Random) HandlePacket(packetId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte) error {
	if hop == 1 {
		return p.randomDissemination(p.RandomPC, packetId, data, int(maxPeers), originalSender, 0)
	}
	if hop == 0 && p.Hop0 {
		return p.randomDissemination(p.RandomPC, packetId, data, int(maxPeers), originalSender, 0)
	}
	return nil
}
