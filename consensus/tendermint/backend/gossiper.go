package backend

import (
	"errors"
	"github.com/autonity/autonity/rlp"
	"math/big"
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
	Route(broadcaster consensus.Broadcaster, committee *types.Committee, msg message.Msg, from common.Address) ([]types.CommitteeMember, error)
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

func (g *Gossiper) Gossip(committee *types.Committee, msg message.Msg) {
	hash := msg.Hash()
	if !g.knownMessages.Contains(hash) {
		g.knownMessages.Add(hash, true)
	}
	if g.broadcaster == nil {
		return
	}
	code := message.NetworkCodes[msg.Code()]
	payload := msg.Payload()

	recipients, err := g.router.Route(g.broadcaster, committee, msg, g.address)
	if err != nil {
		if !errors.Is(err, consensus.ErrFutureEpochMessage) {
			log.Debug("No recipients for proposal", "error", err, "height", msg.H())
			return
		}
		// forward future epoch proposal to all the committee members, as most of them are still in the committee.
		recipients = committee.Members
	}

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

func (g *Gossiper) AskSync(committee *types.Committee, _ uint64, _ int64, syncMsg *message.LostSyncMsg) {
	encoded, err := rlp.EncodeToBytes(syncMsg)
	if err != nil {
		log.Error("Error encoding sync msg", "err", err)
		return
	}

	targets := make([]common.Address, 0, committee.Len())
	for _, val := range committee.Members {
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
				// todo: double check if quorum nodes are sufficient for state recovery?
				// ask to a quorum nodes to sync, 1 must then be honest and updated
				if count.Cmp(bft.Quorum(committee.TotalVotingPower())) >= 0 {
					break
				}
				g.logger.Debug("Asking sync to", "addr", addr)
				go p.Send(message.SyncNetworkMsg, encoded) //nolint

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
