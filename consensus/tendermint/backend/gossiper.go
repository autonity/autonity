package backend

import (
	"math/big"
	"time"

	lru "github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
)

type Gossiper struct {
	recentMessages *lru.LRU[common.Address, *lru.LRU[common.Hash, bool]] // the cache of peer's messages
	knownMessages  *lru.LRU[common.Hash, bool]                           // the cache of self messages
	address        common.Address                                        // address of the local peer
	broadcaster    consensus.Broadcaster
	logger         log.Logger
	stopped        chan struct{}
}

var (
	GossiperStepOneBg   = metrics.NewRegisteredBufferedGauge("gossiper/one.bg", nil, getIntPointer(256))   // time between round start and precommit sent
	GossiperStepTwoBg   = metrics.NewRegisteredBufferedGauge("gossiper/two.bg", nil, getIntPointer(256))   // time between round start and precommit sent
	GossiperStepThreeBg = metrics.NewRegisteredBufferedGauge("gossiper/three.bg", nil, getIntPointer(256)) // time between round start and precommit sent
	GossiperFullBg      = metrics.NewRegisteredBufferedGauge("gossiper/total.bg", nil, getIntPointer(256)) // time between round start and precommit sent
)

func getIntPointer(val int) *int {
	return &val
}
func NewGossiper(recentMessages *lru.LRU[common.Address, *lru.LRU[common.Hash, bool]], knownMessages *lru.LRU[common.Hash, bool], address common.Address, logger log.Logger, stopped chan struct{}) *Gossiper {
	return &Gossiper{
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		address:        address,
		logger:         logger,
		stopped:        stopped,
	}
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
}

func (g *Gossiper) Broadcaster() consensus.Broadcaster {
	return g.broadcaster
}

func (g *Gossiper) RecentMessages() *lru.LRU[common.Address, *lru.LRU[common.Hash, bool]] {
	return g.recentMessages
}

func (g *Gossiper) KnownMessages() *lru.LRU[common.Hash, bool] {
	return g.knownMessages
}

func (g *Gossiper) Address() common.Address {
	return g.address
}

func (g *Gossiper) UpdateStopChannel(stopCh chan struct{}) {
	g.stopped = stopCh
}

func (g *Gossiper) Gossip(committee types.Committee, message message.Msg) {
	n1 := time.Now()
	hash := message.Hash()
	if !g.knownMessages.Contains(hash) {
		g.knownMessages.Add(hash, true)
	}
	GossiperStepOneBg.Add(time.Since(n1).Nanoseconds())
	n2 := time.Now()
	targets := make(map[common.Address]struct{}, len(committee))
	for _, val := range committee {
		if val.Address != g.address {
			targets[val.Address] = struct{}{}
		}
	}
	GossiperStepTwoBg.Add(time.Since(n2).Nanoseconds())
	if g.broadcaster != nil && len(targets) > 0 {
		ps := g.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			n3 := time.Now()
			ms, ok := g.recentMessages.Get(addr)
			if ok {
				if ms.Contains(hash) {
					// This peer had this event, skip it
					continue
				}
				ms.Add(hash, true)
			} else {
				ms = lru.NewLRU[common.Hash, bool](0, nil, time.Second*10)
				ms.Add(hash, true)
				g.recentMessages.Add(addr, ms)
			}

			go p.SendRaw(NetworkCodes[message.Code()], message.Payload()) //nolint
			GossiperStepThreeBg.Add(time.Since(n3).Nanoseconds())
		}
	}
	GossiperFullBg.Add(time.Since(n1).Nanoseconds())
}

func (g *Gossiper) AskSync(header *types.Header) {

	targets := make(map[common.Address]struct{})
	for _, val := range header.Committee {
		if val.Address != g.address {
			targets[val.Address] = struct{}{}
		}
	}

	if g.broadcaster != nil && len(targets) > 0 {
		for {
			ps := g.broadcaster.FindPeers(targets)
			// If we didn't find any peers try again in 10ms or exit if we have
			// been stopped.
			if len(ps) == 0 {
				t := time.NewTimer(retryPeriod * time.Millisecond)
				select {
				case <-t.C:
					continue
				case <-g.stopped:
					return
				}
			}
			count := new(big.Int)
			for addr, p := range ps {
				//ask to a quorum nodes to sync, 1 must then be honest and updated
				if count.Cmp(bft.Quorum(header.TotalVotingPower())) >= 0 {
					break
				}
				g.logger.Debug("Asking sync to", "addr", addr)
				go p.Send(SyncNetworkMsg, []byte{}) //nolint

				member := header.CommitteeMember(addr)
				if member == nil {
					g.logger.Error("could not retrieve member from address")
					continue
				}
				count.Add(count, member.VotingPower)
			}
			break
		}
	}
}
