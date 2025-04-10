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
	preCommits, err := fd.missingPrecommits(presentedRounds, &lostSync)
	if err != nil {
		fd.logger.Error("Going to suspend peer connection", "err", err, "peer", sender)
		return err
	}

	preVotes, err := fd.missingPrevotes(presentedRounds, &lostSync)
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

func (fd *FaultDetector) missingProposals(lostSync *message.LostSyncMsg) []*message.Propose {
	var missingProposals []*message.Propose

	for _, roundView := range lostSync.RoundsViews {
		// get missing proposals.
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

func (fd *FaultDetector) missingPrevotes(presentedRounds map[uint64]struct{}, lostSync *message.LostSyncMsg) ([]*message.Prevote, error) {
	var missingPrevotes []*message.Prevote
	for _, roundView := range lostSync.RoundsViews {
		// get missing prevotes of the round, they could have different value and different signers.
		presentedPrevoteVal := make(map[common.Hash]struct{})
		for i, value := range roundView.Prevotes {

			if _, ok := presentedPrevoteVal[value]; ok {
				return nil, errInvalidLostSyncMsg
			} else {
				presentedPrevoteVal[value] = struct{}{}
			}

			presentedSigners := roundView.PrevotesSigners[i]
			if presentedSigners == nil {
				return nil, errInvalidLostSyncMsg
			}

			// select prevotes of the same round with same value but with different presentedSigners
			prevotes := fd.msgStore.GetPrevotes(lostSync.Height, func(m *message.Prevote) bool {
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

			if len(prevotes) > 0 {
				missingPrevotes = append(missingPrevotes, prevotes...)
			}
		}
		// get prevotes of not presented values.
		prevotes := fd.msgStore.GetPrevotes(lostSync.Height, func(m *message.Prevote) bool {
			if uint64(m.R()) == roundView.Round {
				if _, ok := presentedPrevoteVal[m.Value()]; !ok {
					return true
				}
			}
			return false
		})
		if len(prevotes) > 0 {
			missingPrevotes = append(missingPrevotes, prevotes...)
		}
	}

	// select prevotes of not presented rounds.
	prevotes := fd.msgStore.GetPrevotes(lostSync.Height, func(m *message.Prevote) bool {
		if _, ok := presentedRounds[uint64(m.R())]; !ok {
			return true
		}
		return false
	})

	if len(prevotes) > 0 {
		missingPrevotes = append(missingPrevotes, prevotes...)
	}

	return missingPrevotes, nil
}

func (fd *FaultDetector) missingPrecommits(presentedRounds map[uint64]struct{}, lostSync *message.LostSyncMsg) ([]*message.Precommit, error) {
	var missingPrecommits []*message.Precommit

	for _, roundView := range lostSync.RoundsViews {
		// get missing precommits.
		presentedPrecommitVal := make(map[common.Hash]struct{})
		for i, value := range roundView.Precommits {

			if _, ok := presentedPrecommitVal[value]; ok {
				return nil, errInvalidLostSyncMsg
			} else {
				presentedPrecommitVal[value] = struct{}{}
			}

			presentedSigners := roundView.PrecommitsSigners[i]
			if presentedSigners == nil {
				return nil, errInvalidLostSyncMsg
			}

			// select precommits of the same round with same value but with different presentedSigners
			precommits := fd.msgStore.GetPrecommits(lostSync.Height, func(m *message.Precommit) bool {
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

			if len(precommits) > 0 {
				missingPrecommits = append(missingPrecommits, precommits...)
			}
		}
		// get precommits of not presented values.
		precommits := fd.msgStore.GetPrecommits(lostSync.Height, func(m *message.Precommit) bool {
			if uint64(m.R()) == roundView.Round {
				if _, ok := presentedPrecommitVal[m.Value()]; !ok {
					return true
				}
			}
			return false
		})
		if len(precommits) > 0 {
			missingPrecommits = append(missingPrecommits, precommits...)
		}
	}

	// select precommits of not presented rounds.
	precommits := fd.msgStore.GetPrecommits(lostSync.Height, func(m *message.Precommit) bool {
		if _, ok := presentedRounds[uint64(m.R())]; !ok {
			return true
		}
		return false
	})

	if len(precommits) > 0 {
		missingPrecommits = append(missingPrecommits, precommits...)
	}

	return missingPrecommits, nil
}
