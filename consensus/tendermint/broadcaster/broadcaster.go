package broadcaster

import (
	"context"
	"crypto/ecdsa"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core/types"
	lru "github.com/hashicorp/golang-lru"
)

const (
	inmemoryPeers    = 40
	inmemoryMessages = 1024
)

func NewBroadcaster(privateKey *ecdsa.PrivateKey) Broadcaster {
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)
	return Broadcaster{
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		privateKey:     privateKey,
	}
}

type Broadcaster struct {
	knownMessages *lru.ARCCache // the cache of self messages
	//TODO: ARCChace is patented by IBM, so probably need to stop using it
	recentMessages *lru.ARCCache // the cache of peer's messages
	privateKey     *ecdsa.PrivateKey
}

func (c *Broadcaster) broadcast(ctx context.Context, msg *Message) {
	logger := c.logger.New("step", c.step)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	logger.Debug("broadcasting", "msg", msg.String())
	if err = c.Broadcaster.Broadcast(ctx, c.CommitteeSet(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *Broadcaster) finalizeMessage(msg *Message) ([]byte, error) {
	var err error

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = crypto.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (sb *Broadcaster) Broadcast(ctx context.Context, valSet *committee.Set, payload []byte) error {
	// send to others
	sb.Gossip(ctx, valSet, payload)
	// send to self
	msg := events.MessageEvent{
		Payload: payload,
	}
	sb.postEvent(msg)
	return nil
}

func (sb *Broadcaster) postEvent(event interface{}) {
	go sb.Post(event)
}

func (sb *Broadcaster) Gossip(ctx context.Context, valSet *committee.Set, payload []byte) {
	hash := types.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	targets := make(map[common.Address]struct{})
	for _, val := range valSet.Committee() {
		if val.Address != sb.Address() {
			targets[val.Address] = struct{}{}
		}
	}

	if sb.broadcaster != nil && len(targets) > 0 {
		ps := sb.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			ms, ok := sb.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)

			go p.Send(tendermintMsg, payload) //nolint
		}
	}
}
