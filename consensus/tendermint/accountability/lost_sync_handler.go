package accountability

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
	"time"
)

// lost sync handler, process the ask sync msg from a lost liveness node. As the msg store in the AFD module saves recent
// 256 blocks consensus messages, thus it provides extensive msg views for those chain head synced or un-synced nodes,
// more over that, future round messages can be synced now and the handling of AskSync msg does not block the consensus
// engine anymore.

var askSyncInterval = 5 // 5s

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

	return timeDiff < int64(askSyncInterval)
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

	var lostSync message.LostSyncMsg
	err := rlp.DecodeBytes(payload, &lostSync)
	if err != nil {
		return err
	}

	// sanity checks: no duplicated rounds, and msg set bound checks.
	if len(lostSync.RoundsViews) > constants.MaxRound {
		return errInvalidLostSyncMsg
	}
	presentedRounds := make(map[uint64]struct{})
	for _, v := range lostSync.RoundsViews {
		if v.Round > constants.MaxRound {
			return errInvalidLostSyncMsg
		}
		if _, ok := presentedRounds[v.Round]; ok {
			return errInvalidLostSyncMsg
		} else {
			presentedRounds[v.Round] = struct{}{}
		}
		if len(v.Prevotes) != len(v.PrevotesSigners) || len(v.Precommits) != len(v.PrecommitsSigners) {
			return errInvalidLostSyncMsg
		}
	}

	proposals := fd.missingProposals(&lostSync)
	preCommits, err := fd.missingVotes(presentedRounds, message.PrecommitCode, &lostSync)
	if err != nil {
		fd.logger.Error("Going to suspend peer connection", "err", err, "peer", sender)
		return err
	}

	preVotes, err := fd.missingVotes(presentedRounds, message.PrevoteCode, &lostSync)
	if err != nil {
		fd.logger.Error("Going to suspend per connection", "err", err, "peer", sender)
		return err
	}

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
	if len(proposals) > 0 {
		for _, m := range proposals {
			fd.logger.Info("sending missing proposal to lost sync peer", "value", m.Value(), "H", m.H(), "R", m.R(), "VR", m.ValidRound(), "from", fd.address, "to", sender)
			go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
		}
	}

	// then sends the missing precommits, as precommits could trigger round rotation or a commitment of a value.
	if len(preCommits) > 0 {
		for _, m := range preCommits {
			fd.logger.Info("sending missing prevote to lost sync peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", fd.address, "to", sender)
			go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
		}
	}

	if len(preVotes) > 0 {
		for _, m := range preVotes {
			fd.logger.Info("sending missing prevote to lost sync peer", "value", m.Value(), "H", m.H(), "R", m.R(), "from", fd.address, "to", sender)
			go peer.SendRaw(message.NetworkCodes[m.Code()], m.Payload())
		}
	}

	return nil
}

// missingProposals collects all the missing proposals of a consensus instance base on the asker's view.
func (fd *FaultDetector) missingProposals(lostSync *message.LostSyncMsg) []*message.Propose {
	var missingProposals []*message.Propose

	for _, roundView := range lostSync.RoundsViews {
		// from the asker's round view, if it missed a proposal of that round, then we need to send
		// any proposal of that round, otherwise we don't send it as the asker already prevoted for a value.
		if roundView.Proposal == nilValue {
			proposals := fd.msgStore.GetProposals(lostSync.Height, func(m *message.Propose) bool {
				return uint64(m.R()) == roundView.Round
			})
			if len(proposals) > 0 {
				missingProposals = append(missingProposals, proposals...)
			}
		}
	}

	return missingProposals
}

// missingVotes collects those missing prevotes, or precommits of the asker. They include those votes of missing signers, values and rounds.
func (fd *FaultDetector) missingVotes(presentedRounds map[uint64]struct{}, step uint8, lostSync *message.LostSyncMsg) ([]message.Vote, error) {
	var missingVotes []message.Vote
	for _, roundView := range lostSync.RoundsViews {
		// get missing votes of the round, they could have different value and different signers.
		presentedValue := make(map[common.Hash]struct{})

		// query for prevotes by default,
		knownVotes := roundView.Prevotes
		knownVotesSigners := roundView.PrevotesSigners
		if message.PrecommitCode == step {
			knownVotes = roundView.Precommits
			knownVotesSigners = roundView.PrecommitsSigners
		}

		for i, value := range knownVotes {

			if _, ok := presentedValue[value]; ok {
				return nil, errInvalidLostSyncMsg
			} else {
				presentedValue[value] = struct{}{}
			}

			presentedSigners := knownVotesSigners[i]
			if presentedSigners == nil {
				return nil, errInvalidLostSyncMsg
			}

			// select votes of the same round with same value but with different presentedSigners
			votes := fd.msgStore.GetVotes(lostSync.Height, step, func(m message.Vote) bool {
				if uint64(m.R()) == roundView.Round && m.Value() == value {
					// only with signers which is not in the presentedSigners
					for _, idx := range m.Signers().FlattenUniq() {
						if presentedSigners.Bit(idx) == 0 {
							return true
						}
					}
				}
				return false
			})

			if len(votes) > 0 {
				missingVotes = append(missingVotes, votes...)
			}
		}
		// get votes of not presented values.
		votes := fd.msgStore.GetVotes(lostSync.Height, step, func(m message.Vote) bool {
			if uint64(m.R()) == roundView.Round {
				if _, ok := presentedValue[m.Value()]; !ok {
					return true
				}
			}
			return false
		})
		if len(votes) > 0 {
			missingVotes = append(missingVotes, votes...)
		}
	}

	// select votes of not presented rounds.
	votes := fd.msgStore.GetVotes(lostSync.Height, step, func(m message.Vote) bool {
		if _, ok := presentedRounds[uint64(m.R())]; !ok {
			return true
		}
		return false
	})

	if len(votes) > 0 {
		missingVotes = append(missingVotes, votes...)
	}

	return missingVotes, nil
}
