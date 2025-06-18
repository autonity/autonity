package backend

import (
	"fmt"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
)

const cleanUpInterval = 60 // 60s

func (sb *Backend) startRateLimiterGCRoutine() {
	sb.cleanupTicker = time.NewTicker(time.Second * cleanUpInterval)
	sb.wg.Add(1)
	go func() {
		defer sb.wg.Done()
		defer sb.cleanupTicker.Stop()
		for {
			select {
			case <-sb.cleanupTicker.C:
				sb.askSyncRateLimiter.Cleanup()
			case <-sb.stopped:
				return
			}
		}
	}()
}

// syncPeer process the ask sync msg from a lost liveness node.
func (sb *Backend) syncPeer(payload []byte, sender common.Address, errCh chan<- error) {
	err := sb.handleAskSyncEvent(payload, sender)
	if err != nil {
		sb.logger.Error("syncPeer", "error", err, "sender", sender)
		// the errors return from handler could freeze the peer connection for 30 seconds by according to dev p2p protocol.
		select {
		case errCh <- err:
		default: // do nothing
		}
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

	msgCount := 0
	// prioritize the sending of missing proposals.
	for _, m := range proposals {
		msgCount++
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}
	if msgCount > 0 {
		sb.logger.Debug("sending missing proposal to remote peer",  "peer", sender, "count", msgCount)
	}

	msgCount = 0
	// then sends the missing precommits, as precommits could trigger round rotation or a commitment of a value.
	for _, m := range precommits {
		msgCount++
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}
	if msgCount > 0 {
		sb.logger.Debug("sending missing precommits to remote peer",  "peer", sender, "count", msgCount)
	}

	msgCount = 0
	for _, m := range prevotes {
		msgCount++
		go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
	}
	if msgCount > 0 {
		sb.logger.Debug("sending missing prevotes to remote peer",  "peer", sender, "count", msgCount)
	}

	return nil
}

// missingProposals collects all the missing proposals of a consensus instance base on the asker's view.
func (sb *Backend) missingProposals(askSync *message.AskSyncMsg) []*message.Propose {
	hasProposalMap := askSync.HasProposal()
	missingProposals := sb.MsgStore.GetProposals(askSync.Height, func(m *message.Propose) bool {
		_, hasProposal := hasProposalMap[uint64(m.R())]
		return !hasProposal
	})

	return missingProposals
}

// missingPrevotes collects all the missing prevotes of a consensus instance base on the asker's view.
func (sb *Backend) missingPrevotes(askSync *message.AskSyncMsg) []*message.Prevote {
	prevoteSigners := askSync.Prevotes()

	missingPrevotes := sb.MsgStore.GetPrevotes(askSync.Height, func(m *message.Prevote) bool {
		// return all prevotes if the remote node doesn't know this round
		msgRound := uint64(m.R())
		_, knownRound := prevoteSigners[msgRound]
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

	precommitSigners := askSync.Precommits()

	missingPrecommits := sb.MsgStore.GetPrecommits(askSync.Height, func(m *message.Precommit) bool {
		// return all precommits if the remote node doesn't know this round
		msgRound := uint64(m.R())
		_, knownRound := precommitSigners[msgRound]
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
