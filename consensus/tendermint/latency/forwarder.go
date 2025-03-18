package latency

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/log"
)

// proposeNetworkMsg is redefined here to avoid circular dependencies
var proposeNetworkMsg uint64 = 0x11

func forward(router *Router, bc *core.BlockChain, m message.Msg, sender common.Address) {
	if router == nil {
		log.Error("Forwarder: Router is nil")
		return
	}

	// no need to forward if it's not a proposal
	if m.Code() != message.ProposalCode {
		return
	}

	// cant forward if we don't have clustering for this height
	if !router.ClusteringActive(m.H()) {
		log.Debug("Forwarder: Clustering is not active at height", "height", m.H())
		return
	}

	committee, err := bc.CommitteeByHeight(m.H())
	if err != nil {
		log.Info("Forwarder: Failed to get committee", "error", err, "height", m.H())
		return
	}

	recipients := router.Route(committee, m, sender)
	if len(recipients) == 0 {
		log.Debug("No recipients for proposal", "proposal", m)
		return
	}

	for _, recipient := range recipients {
		if recipient.Address == sender {
			continue
		}
		if p, ok := router.broadcaster.FindPeer(recipient.Address); ok {
			if p.Cache().Contains(m.Hash()) {
				// This peer had this event, skip it
				continue
			}
			p.Cache().Add(m.Hash(), true)
			go p.SendRaw(proposeNetworkMsg, m.Payload()) //nolint
		} else {
			//todo: shall we select other backups for live ness?
		}
	}
}
