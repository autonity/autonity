package backend

import (
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
)

type Gossiper struct {
	knownMessages      *fixsizecache.Cache[common.Hash, bool] // the cache of self messages
	address            common.Address                         // address of the local peer
	broadcaster        consensus.Broadcaster
	logger             log.Logger
	stopped            chan struct{}
	concurrencyLimiter chan messageInfo
}

func NewGossiper(knownMessages *fixsizecache.Cache[common.Hash, bool], address common.Address, logger log.Logger, stopped chan struct{}) *Gossiper {
	g := &Gossiper{
		knownMessages:      knownMessages,
		address:            address,
		logger:             logger,
		stopped:            stopped,
		concurrencyLimiter: make(chan messageInfo, 64),
	}
	go g.checkConcurrencyLimiter()
	return g
}

func (g *Gossiper) SetBroadcaster(broadcaster consensus.Broadcaster) {
	g.broadcaster = broadcaster
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

func (g *Gossiper) Gossip(committee types.Committee, message message.Msg) {
	hash := message.Hash()
	if !g.knownMessages.Contains(hash) {
		g.knownMessages.Add(hash, true)
	}
	if g.broadcaster == nil {
		return
	}
	code := NetworkCodes[message.Code()]
	payload := message.Payload()
	for _, val := range committee {
		if val.Address == g.address {
			continue
		}
		if p, ok := g.broadcaster.FindPeer(val.Address); ok {
			if p.Cache().Contains(hash) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(hash, true)
			g.concurrencyLimiter <- messageInfo{
				sender: val.Address,
				height: message.H(),
				round:  message.R(),
				code:   message.Code(),
			}
			go func() {
				defer func() {
					<-g.concurrencyLimiter
				}()
				p.SendRaw(code, payload) //nolint
			}()
		}
	}
}

type messageInfo struct {
	sender common.Address
	height uint64
	round  int64
	code   uint8
}

func (g *Gossiper) checkConcurrencyLimiter() {
	t := time.NewTicker(30 * time.Second)
	b := make([]messageInfo, 64)
	for range t.C {
		i := 0
	LOOP:
		for {
			select {
			case msg := <-g.concurrencyLimiter:
				b[i] = msg
				i++
			default:
				break LOOP
			}
		}

		sb := strings.Builder{}
		for k := 0; k < i; k++ {
			sb.WriteString(
				b[k].sender.String() + "|" +
					strconv.Itoa(int(b[k].height)) + "|" +
					strconv.Itoa(int(b[k].round)) + "|" +
					strconv.Itoa(int(b[k].code)) + " - ",
			)
			g.concurrencyLimiter <- b[k]
		}
		g.logger.Warn("CCLSTATE", "size", i, "buf", sb.String())
	}

}

func (g *Gossiper) AskSync(header *types.Header) {

	targets := make([]common.Address, 0, len(header.Committee))
	for _, val := range header.Committee {
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
