package backend

import (
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
)

// SyncPeer is a backend interface to be called by core, it assumes the askSync message was checked by the core handler.
func (sb *Backend) SyncPeer(askSync *message.AskSyncMsg, sender common.Address) {
	if sb.Broadcaster == nil {
		sb.logger.Warn("p2p protocol handler is not ready yet")
		return
	}
	peer, ok := sb.Broadcaster.FindPeer(sender)
	if !ok {
		sb.logger.Debug("no peer connection for sender", "peer", sender)
		return
	}

	// fetch remote's peer missing messages
	proposals := sb.missingProposals(askSync)
	prevotes := sb.missingPrevotes(askSync)
	precommits := sb.missingPrecommits(askSync)

	// prioritize the sending of missing proposals.
	for _, m := range proposals {
		sb.logger.Debug("sending missing proposal to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "VR", m.ValidRound(), "from", sb.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	// then sends the missing precommits, as precommits could trigger round rotation or a commitment of a value.
	for _, m := range precommits {
		sb.logger.Debug("sending missing precommits to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", sb.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	for _, m := range prevotes {
		sb.logger.Debug("sending missing prevotes to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", sb.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}
}

// handleAskSyncEvent handles the ask sync request from a lost sync validator or from a rebooting validator.
// Any error return from this function will drop the remote peer.
func (sb *Backend) handleAskSyncEvent(payload []byte, sender common.Address) error {
	if sb.Broadcaster == nil {
		sb.logger.Warn("p2p protocol handler is not ready yet")
		return nil
	}
	peer, ok := sb.Broadcaster.FindPeer(sender)
	if !ok {
		sb.logger.Debug("no peer connection for sender", "peer", sender)
		return nil
	}

	if err := sb.askSyncRateLimiter.Allow(sender); err != nil {
		return err
	}

	askSync := new(message.AskSyncMsg)
	if err := rlp.DecodeBytes(payload, askSync); err != nil {
		return fmt.Errorf("cannot decode ask sync msg: %w", err)
	}

	if err := askSync.Validate(); err != nil {
		return fmt.Errorf("ask sync msg sanity check failed: %w", err)
	}

	// fetch remote's peer missing messages
	proposals := sb.missingProposals(askSync)
	prevotes := sb.missingPrevotes(askSync)
	precommits := sb.missingPrecommits(askSync)

	// prioritize the sending of missing proposals.
	for _, m := range proposals {
		sb.logger.Debug("sending missing proposal to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "VR", m.ValidRound(), "from", sb.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	// then sends the missing precommits, as precommits could trigger round rotation or a commitment of a value.
	for _, m := range precommits {
		sb.logger.Debug("sending missing precommits to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", sb.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	for _, m := range prevotes {
		sb.logger.Debug("sending missing prevotes to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", sb.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	return nil
}

// missingProposals collects all the missing proposals of a consensus instance base on the asker's view.
func (sb *Backend) missingProposals(askSync *message.AskSyncMsg) []*message.Propose {
	rounds := askSync.Rounds()
	nilProposal := askSync.NilProposal()
	missingProposals := sb.MsgStore.GetProposals(askSync.Height, func(m *message.Propose) bool {
		_, knownRound := rounds[uint64(m.R())]
		_, unknownProposal := nilProposal[uint64(m.R())]
		return !knownRound || unknownProposal
	})

	return missingProposals
}

// missingPrevotes collects all the missing prevotes of a consensus instance base on the asker's view.
func (sb *Backend) missingPrevotes(askSync *message.AskSyncMsg) []*message.Prevote {

	rounds := askSync.Rounds()
	prevoteSigners := askSync.Prevotes()

	missingPrevotes := sb.MsgStore.GetPrevotes(askSync.Height, func(m *message.Prevote) bool {
		// return all prevotes if the remote node doesn't know this round
		msgRound := uint64(m.R())
		_, knownRound := rounds[msgRound]
		if !knownRound {
			return true
		}
		// if the round is known, return prevotes that have a value that the remote node didn't see
		signers, knownValue := prevoteSigners[msgRound][m.Value()]
		if !knownValue {
			return true
		}
		// if both the round and value are known, return prevotes that have signers that the remote node didn't see
		for _, idx := range m.Signers().FlattenUniq() {
			if signers.Bit(idx) == 0 {
				return true
			}
		}
		// otherwise, the node already has this message
		return false
	})

	return missingPrevotes
}

// missingPrecommits collects all the missing precommits of a consensus instance base on the asker's view.
func (sb *Backend) missingPrecommits(askSync *message.AskSyncMsg) []*message.Precommit {

	rounds := askSync.Rounds()
	precommitSigners := askSync.Precommits()

	missingPrecommits := sb.MsgStore.GetPrecommits(askSync.Height, func(m *message.Precommit) bool {
		// return all precommits if the remote node doesn't know this round
		msgRound := uint64(m.R())
		_, knownRound := rounds[msgRound]
		if !knownRound {
			return true
		}
		// if the round is known, return precommits that have a value that the remote node didn't see
		signers, knownValue := precommitSigners[msgRound][m.Value()]
		if !knownValue {
			return true
		}
		// if both the round and value are known, return precommits that have signers that the remote node didn't see
		for _, idx := range m.Signers().FlattenUniq() {
			if signers.Bit(idx) == 0 {
				return true
			}
		}
		// otherwise, the node already has this message
		return false
	})

	return missingPrecommits
}
