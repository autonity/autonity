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
		go func() { // todo: TEST WITHOUT PROTORW !!!
			err := peer.DisseminateRequest(p.Code, packetId, 0, p.State.Id, uint64(maxPeers), data)
			if err != nil {
				log.Error("DisseminateRequest err:", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (p *Broadcast) HandlePacket(requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte) error {
	// Simple broadcast - nothing to propagate.
	return nil
}

func (p *BroadcastBlocking) Execute(packetId uint64, data []byte, maxPeers int) error {
	for i := 0; i < maxPeers; i++ {
		peer := p.Peers(i)
		if peer == nil {
			continue
		}
		err := peer.DisseminateRequest(p.Code, packetId, 0, p.State.Id, uint64(maxPeers), data)
		if err != nil {
			log.Error("DisseminateRequest err:", err)
		}
	}
	return nil
}

func (p *BroadcastBlocking) HandlePacket(requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte) error {
	return nil
}
