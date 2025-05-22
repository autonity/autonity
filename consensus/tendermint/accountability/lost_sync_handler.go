package accountability

import (
	"errors"
	"fmt"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
)

// lost sync handler, process the ask sync msg from a lost liveness node. As the msg store in the AFD module saves recent
// 256 blocks consensus messages, thus it provides extensive msg views for those chain head synced or un-synced nodes,
// more over that, future round messages can be synced now and the handling of AskSync msg does not block the consensus
// engine anymore.

var errAskSyncOverRated = errors.New("ask sync over rated")

type AskSyncRateLimiter struct {
	lastRequestTSs map[common.Address]int64
}

func NewAskSyncRateLimiter() *AskSyncRateLimiter {
	return &AskSyncRateLimiter{
		lastRequestTSs: make(map[common.Address]int64),
	}
}

func (r *AskSyncRateLimiter) overRated(asker common.Address) bool {
	now := time.Now().Unix()

	lastTS, exists := r.lastRequestTSs[asker]
	if !exists {
		r.lastRequestTSs[asker] = now
		return false
	}

	r.lastRequestTSs[asker] = now
	timeDiff := now - lastTS

	return timeDiff < int64(constants.AskSyncInterval)
}

func (r *AskSyncRateLimiter) resetRateLimiter() {
	for k := range r.lastRequestTSs {
		delete(r.lastRequestTSs, k)
	}
}

// handleLostSyncEvent handles the ask sync request from a lost sync validator or from a rebooting validator.
// Any error return from this function will drop the remote peer.
func (fd *FaultDetector) handleLostSyncEvent(payload []byte, sender common.Address) error {

	if fd.askSyncRateLimiter.overRated(sender) {
		return errAskSyncOverRated
	}

	lostSync := new(message.LostSyncMsg)
	if err := rlp.DecodeBytes(payload, lostSync); err != nil {
		return fmt.Errorf("cannot decode ask sync msg: %w", err)
	}

	if err := lostSync.Validate(); err != nil {
		return fmt.Errorf("ask sync msg sanity check failed: %w", err)
	}

	// fetch remote's peer missing messages
	proposals := fd.missingProposals(lostSync)
	prevotes := fd.missingPrevotes(lostSync)
	precommits := fd.missingPrecommits(lostSync)

	// broadcast them to the missing peer
	if fd.broadcaster == nil {
		fd.logger.Warn("p2p protocol handler is not ready yet")
		return nil
	}

	peer, ok := fd.broadcaster.FindPeer(sender)
	if !ok {
		fd.logger.Debug("no peer connection for sender", "peer", sender)
		return nil
	}

	// prioritize the sending of missing proposals.
	for _, m := range proposals {
		fd.logger.Debug("sending missing proposal to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "VR", m.ValidRound(), "from", fd.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	// then sends the missing precommits, as precommits could trigger round rotation or a commitment of a value.
	for _, m := range precommits {
		fd.logger.Debug("sending missing precommits to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", fd.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	for _, m := range prevotes {
		fd.logger.Debug("sending missing prevotes to remote peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", fd.address, "to", sender)
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}

	return nil
}

// missingProposals collects all the missing proposals of a consensus instance base on the asker's view.
func (fd *FaultDetector) missingProposals(lostSync *message.LostSyncMsg) []*message.Propose {
	rounds := lostSync.Rounds()
	nilProposal := lostSync.NilProposal()
	missingProposals := fd.msgStore.GetProposals(lostSync.Height, func(m *message.Propose) bool {
		_, knownRound := rounds[uint64(m.R())]
		_, unknownProposal := nilProposal[uint64(m.R())]
		return !knownRound || unknownProposal
	})

	return missingProposals
}

// missingPrevotes collects all the missing prevotes of a consensus instance base on the asker's view.
func (fd *FaultDetector) missingPrevotes(lostSync *message.LostSyncMsg) []*message.Prevote {

	rounds := lostSync.Rounds()
	prevoteSigners := lostSync.Prevotes()

	missingPrevotes := fd.msgStore.GetPrevotes(lostSync.Height, func(m *message.Prevote) bool {
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
func (fd *FaultDetector) missingPrecommits(lostSync *message.LostSyncMsg) []*message.Precommit {

	rounds := lostSync.Rounds()
	precommitSigners := lostSync.Precommits()

	missingPrecommits := fd.msgStore.GetPrecommits(lostSync.Height, func(m *message.Precommit) bool {
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
