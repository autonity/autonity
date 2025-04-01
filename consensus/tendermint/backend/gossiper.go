package backend

import (
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

type router interface {
	Route(committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error)
	SetBroadcaster(broadcaster consensus.Broadcaster)
}

type Gossiper struct {
	knownMessages *fixsizecache.Cache[common.Hash, bool] // the cache of self messages
	address       common.Address                         // address of the local peer
	broadcaster   consensus.Broadcaster
	logger        log.Logger
	stopped       chan struct{}
	router        router
}

func NewGossiper(
	knownMessages *fixsizecache.Cache[common.Hash, bool],
	address common.Address,
	logger log.Logger,
	stopped chan struct{},
	router router,
) *Gossiper {
	return &Gossiper{
		knownMessages: knownMessages,
		address:       address,
		logger:        logger,
		stopped:       stopped,
		router:        router,
	}
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
	g.router.SetBroadcaster(broadcaster)
}

func (g *Gossiper) Broadcaster() consensus.Broadcaster {
	return g.broadcaster
}

func (g *Gossiper) KnownMessages() *fixsizecache.Cache[common.Hash, bool] {
	return g.knownMessages
}

func (g *Gossiper) Address() common.Address {
	return g.address
}

func (g *Gossiper) UpdateStopChannel(stopCh chan struct{}) {
	g.stopped = stopCh
}

func (g *Gossiper) SlowGossip(committee *types.Committee, msg message.Msg) {
	// only gossip to very small committee
	numTargets := len(committee.Members)
	if numTargets > 10 { // todo: minimum nodes to start slow gossip
		numTargets = int(math.Sqrt(float64(len(committee.Members))))
	} else if numTargets == 0 {
		log.Error("no target to slow gossip", "num", len(committee.Members), "numTargets", numTargets, "committee", committee.Members)
		return
	}
	targetIndices := rand.Perm(numTargets)
	recipients := make([]types.CommitteeMember, numTargets)
	log.Debug("total committee members", "num", len(committee.Members), "numTargets", numTargets, "committee", committee.Members)
	for i := 0; i < numTargets; i++ {
		recipients[i] = committee.Members[targetIndices[i]]
	}
	g.gossip(msg, recipients)
}

func (g *Gossiper) gossip(msg message.Msg, recipients []types.CommitteeMember) {
	hash := msg.Hash()
	if !g.knownMessages.Contains(hash) {
		g.knownMessages.Add(hash, true)
	}
	if g.broadcaster == nil {
		return
	}
	code := message.NetworkCodes[msg.Code()]
	payload := msg.Payload()
	lostPeers := make([]common.Address, 0)
	for _, val := range recipients {
		if val.Address == g.address {
			continue
		}
		if p, ok := g.broadcaster.FindPeer(val.Address); ok {
			if p.Cache().Contains(hash) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(hash, true)
			go p.SendRaw(code, payload) //nolint
		} else {
			// todo: Jason, shall we select other backups for liveness?
			lostPeers = append(lostPeers, val.Address)
		}
	}
	if len(lostPeers) > 0 {
		g.logger.Debug("Gossiper: peers not found", "len", len(lostPeers), "peers", lostPeers)
	}

}

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg) {
	recipients, err := g.router.Route(committee, msg, g.address)
	if err != nil {
		//if !errors.Is(err, consensus.ErrFutureEpochMessage) {
		//	log.Debug("No recipients for message", "error", err, "height", msg.H(), "message type", msg.Code())
		//	return
		//}
		// forward future epoch proposal to all the committee members, as most of them are still in the committee.
		recipients = committee.Members
	}
	g.gossip(msg, recipients)
}

func (g *Gossiper) AskSync(committee *types.Committee, coreHeight uint64, round int64) {
	f := message.Fake{FakeHeight: coreHeight, FakeRound: uint64(round), FakeCode: uint8(message.SyncNetworkMsg)}
	recipients, err := g.router.Route(committee, f, g.address)
	if err != nil {
		log.Error("Error selecting peers members to broadcast sync", "error", err)
		recipients = committee.Members
	}

	targets := make([]common.Address, 0, committee.Len())
	for _, val := range recipients {
		if val.Address != g.address {
			targets = append(targets, val.Address)
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
				if count.Cmp(bft.Quorum(committee.TotalVotingPower())) >= 0 {
					break
				}
				g.logger.Debug("Asking sync to", "addr", addr)
				go p.Send(message.SyncNetworkMsg, []byte{}) //nolint

				member := committee.MemberByAddress(addr)
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
