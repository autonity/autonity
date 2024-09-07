package strats

import (
	"sync"

	"github.com/autonity/autonity/log"
)

type Broadcast struct {
	BaseStrategy
}
type BroadcastBlocking struct {
	BaseStrategy
}

func init() {
	registerStrategy("Simple Broadcast (Non-Blocking) - No Relays", func(base BaseStrategy) Strategy {
		return &Broadcast{base}
	})
	registerStrategy("Simple Broadcast (Blocking) - No Relays", func(base BaseStrategy) Strategy {
		return &BroadcastBlocking{base}
	})
}

func (p *Broadcast) Execute(packetId uint64, data []byte, maxPeers int) error {
	var wg sync.WaitGroup
	for i := 0; i < maxPeers; i++ {
		peer := p.Peers(i)
		if peer == nil {
			continue
		}
		wg.Add(1)
		go func() {
			err := peer.DisseminateRequest(p.Code, packetId, 0, p.State.Id, uint64(maxPeers), data, false, 0, 0)
			if err != nil {
				log.Error("DisseminateRequest err:", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (p *Broadcast) HandlePacket(packetId uint64, hop uint8, originalSender uint64, receivedFrom uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	// Simple broadcast - nothing to propagate.
	return nil
}

func (p *Broadcast) ConstructGraph(maxPeers int) error {
	return nil
}

func (p *Broadcast) GraphReadyForPeer(peerID int) {}

func (p *Broadcast) IsGraphReadyForPeer(peerID int) bool {
	return true
}

func (p *BroadcastBlocking) Execute(packetId uint64, data []byte, maxPeers int) error {
	for i := 0; i < maxPeers; i++ {
		peer := p.Peers(i)
		if peer == nil {
			continue
		}
		err := peer.DisseminateRequest(p.Code, packetId, 0, p.State.Id, uint64(maxPeers), data, false, 0, 0)
		if err != nil {
			log.Error("DisseminateRequest err:", err)
		}
	}
	return nil
}

func (p *BroadcastBlocking) HandlePacket(requestId uint64, hop uint8, originalSender uint64, from uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error {
	return nil
}

func (p *BroadcastBlocking) ConstructGraph(maxPeers int) error {
	return nil
}

func (p *BroadcastBlocking) GraphReadyForPeer(peerID int) {}

func (p *BroadcastBlocking) IsGraphReadyForPeer(peerID int) bool {
	return true
}
