package strats

import (
	"sync"

	"github.com/autonity/autonity/log"
)

type Broadcast struct {
	BaseStrategy
}

func init() {
	registerStrategy("Simple Broadcast", func(base BaseStrategy) Strategy {
		return &Broadcast{base}
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
