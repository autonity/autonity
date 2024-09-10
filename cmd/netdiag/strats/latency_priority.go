package strats

import (
	"sort"

	"github.com/autonity/autonity/log"
)

type LowRTT struct {
	BaseStrategy
	RandomRatio int
}

func init() {
	registerStrategy("Low RTT Priority - Relay 10% Random", func(base BaseStrategy) Strategy {
		return &LowRTT{base, 10}
	})
	registerStrategy("Low RTT Priority - Relay 50% Random", func(base BaseStrategy) Strategy {
		return &LowRTT{base, 50}
	})
}

func (l *LowRTT) Execute(packetId uint64, data []byte, maxPeers int) error {
	sortedPeers := make([]Peer, 0)
	for i := 0; i < maxPeers; i++ {
		if p := l.Peers(i); p != nil {
			sortedPeers = append(sortedPeers, p)
		}
	}

	sort.Slice(sortedPeers, func(i, j int) bool {
		return sortedPeers[i].RTT() < sortedPeers[j].RTT()
	})

	for _, p := range sortedPeers {
		err := p.DisseminateRequest(l.Code, packetId, 0, l.State.Id, uint64(maxPeers), data, false, 0, 0)
		if err != nil {
			log.Error("DisseminateRequest err:", err)
		}
	}
	return nil
}

func (l *LowRTT) HandlePacket(requestId uint64, hop uint8, originalSender uint64, _ uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	// randomDissemination is defined in random.go
	return l.randomDissemination(l.RandomRatio, requestId, data, int(maxPeers), originalSender, int(hop+1), partial, seqNum, total)
}

func (l *LowRTT) ConstructGraph(maxPeers int) error {
	return nil
}

func (l *LowRTT) GraphReadyForPeer(peerID int) {}

func (l *LowRTT) IsGraphReadyForPeer(peerID int) bool {
	return true
}

func (l *LowRTT) LatencyType() (LatencyType, int) {
	return LatencyTypeRelative, l.State.Peers
}
