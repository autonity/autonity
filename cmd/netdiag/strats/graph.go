package strats

import (
	"sync"

	"github.com/autonity/autonity/log"
)

type GraphConstructor interface {
	ConstructGraph(maxPeers int) error
	RouteBroadcast(originalSender int, fromNode int) ([]int, error)
	LatencyType() (LatencyType, int)
}

type GraphStrategy struct {
	BaseStrategy
	GraphConstructor
	peerGraphReady []bool
}

func createGraphStrategy(base BaseStrategy, graphConstructor GraphConstructor) *GraphStrategy {
	strategy := &GraphStrategy{base, graphConstructor, make([]bool, base.State.Peers)}
	return strategy
}

func (g *GraphStrategy) GraphReadyForPeer(peerID int) {
	g.peerGraphReady[peerID] = true
}

func (g *GraphStrategy) IsGraphReadyForPeer(peerID int) bool {
	return g.peerGraphReady[peerID]
}

func (g *GraphStrategy) Execute(packetId uint64, data []byte, _ int) error {
	return g.send(g.State.Id, g.State.Id, packetId, 1, data)
}

func (g *GraphStrategy) HandlePacket(packetId uint64, hop uint8, originalSender uint64, fromNode uint64, _ uint64, data []byte, partial bool, seqNum, total uint16) error {
	return g.send(originalSender, fromNode, packetId, 0, data)
}

func (g *GraphStrategy) send(root, from, packetId uint64, hop uint8, data []byte) error {
	log.Debug("Sending packet", "root", root, "packetId", packetId, "hop", hop, "localId", g.State.Id)
	// first collect all peers to send to
	destinationPeers, err := g.RouteBroadcast(int(root), int(from))
	if err != nil {
		log.Error("RouteBroadcast err", "error", err)
		return err
	}
	// check whether we need to send at all
	if len(destinationPeers) == 0 {
		log.Debug("No peers in destinationPeers")
		return nil
	}

	for _, peerID := range destinationPeers {
		if peer := g.Peers(peerID); peer == nil {
			log.Error("Peer not found", "peerID", peerID)
			return errPeerNotFound
		}
	}

	var wg sync.WaitGroup
	for _, peerID := range destinationPeers {
		peer := g.Peers(peerID)
		wg.Add(1)
		peerID := peerID
		go func(p Peer) {
			log.Debug("Sending packet to peer", "peerID", peerID)
			err := p.DisseminateRequest(g.Code, packetId, hop, root, uint64(0), data, false, 0, 0)
			if err != nil {
				log.Error("DisseminateRequest err:", err)
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	return nil
}
