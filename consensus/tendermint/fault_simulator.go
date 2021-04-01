package tendermint

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

// createMisbehaviourContext create misbehaviour context of msgs by according to the configuration and innocent msg.
func (b *Bridge) createMisbehaviourContext(innocentMsg *algorithm.ConsensusMessage) (msgs [][]byte) {
	nonNilValue := algorithm.ValueID(common.Hash{0x1})
	msg := func(s algorithm.Step, h uint64, r int64, v algorithm.ValueID, vr int64) *algorithm.ConsensusMessage {
		return &algorithm.ConsensusMessage{
			MsgType:    s,
			Height:     h,
			Round:      r,
			Value:      v,
			ValidRound: vr,
		}
	}

	nextProposeRound := func(currentRound int64) (int64, error) {
		r := currentRound + 1
		for ; ; r++ {
			p, err := b.proposerAddr(b.lastHeader, r)
			if err != nil {
				return 0, err
			}
			if p == b.address {
				break
			}
		}
		return r, nil
	}

	// simulate a context of msgs that node propose a new proposal rather than the one it locked at previous round.
	maliciousContextPN := func() [][]byte {
		// find a next proposing round.
		nPR, err := nextProposeRound(innocentMsg.Round)
		if err != nil {
			return nil
		}
		// simulate a preCommit msg that locked a value at previous round than next proposing round.
		msgEvidence := msg(algorithm.Precommit, innocentMsg.Height, nPR-1, nonNilValue, 0)
		mE, err := EncodeSignedMessage(msgEvidence, b.key, nil)
		if err != nil {
			return nil
		}
		// simulate a proposal that propose a new value with -1 as the valid round.
		msgPN := msg(algorithm.Propose, innocentMsg.Height, nPR, innocentMsg.Value, -1)
		mPN, err := EncodeSignedMessage(msgPN, b.key, b.msgStore.value(common.Hash(msgPN.Value)))
		if err != nil {
			return nil
		}
		return append(msgs, mE, mPN)
	}

	// simulate a context of msgs that node propose a proposal that proposed a value for which it is not the one it
	// locked on.
	maliciousContextPO := func() [][]byte {
		// find a next proposing round.
		nPR, err := nextProposeRound(innocentMsg.Round)
		if err != nil {
			return nil
		}
		vR := nPR - 1
		// simulate a preCommit msg that locked a value at vR.
		msgEvidence := msg(algorithm.Precommit, innocentMsg.Height, vR, nonNilValue, 0)
		mE, err := EncodeSignedMessage(msgEvidence, b.key, nil)
		if err != nil {
			return nil
		}
		// simulate a proposal that node propose for an old value which it is not the one it locked.
		msgPO := msg(algorithm.Propose, innocentMsg.Height, nPR, innocentMsg.Value, vR)
		mPO, err := EncodeSignedMessage(msgPO, b.key, b.msgStore.value(common.Hash(msgPO.Value)))
		if err != nil {
			return nil
		}
		return append(msgs, mE, mPO)
	}

	// simulate a context of msgs that a node preVote for a new value rather than the one it locked on.
	// preCommit (h, r, v1)
	// propose   (h, r+1, v2)
	// preVote   (h, r+1, v2)
	maliciousContextPVN := func() [][]byte {
		// find a next proposing round.
		nPR, err := nextProposeRound(innocentMsg.Round)
		if err != nil {
			return nil
		}
		r := nPR - 1
		// simulate a preCommit at round r, for value v1.
		msgEvidence := msg(algorithm.Precommit, innocentMsg.Height, r, nonNilValue, 0)
		mE, err := EncodeSignedMessage(msgEvidence, b.key, nil)
		if err != nil {
			return nil
		}
		// simulate a proposal at round r+1, for value v2.
		msgProposal := msg(algorithm.Propose, innocentMsg.Height, nPR, innocentMsg.Value, -1)
		mP, err := EncodeSignedMessage(msgProposal, b.key, b.msgStore.value(common.Hash(innocentMsg.Value)))
		if err != nil {
			return nil
		}
		// simulate a preVote at round r+1, for value v2, this preVote for new value break PVN.
		msgPVN := msg(algorithm.Prevote, innocentMsg.Height, nPR, innocentMsg.Value, 0)
		mPVN, err := EncodeSignedMessage(msgPVN, b.key, nil)
		if err != nil {
			return nil
		}

		return append(msgs, mE, mP, mPVN)
	}

	// simulate a context of msgs that node preCommit at a value V of the round where exist quorum preVotes
	// for not V.
	maliciousContextC := func() [][]byte {
		msgC := msg(algorithm.Precommit, innocentMsg.Height, innocentMsg.Round, nonNilValue, 0)
		mC, err := EncodeSignedMessage(msgC, b.key, nil)
		if err != nil {
			return nil
		}

		return append(msgs, mC)
	}

	// simulate an invalid proposal.
	invalidProposal := func() [][]byte {
		header := &types.Header{Number: new(big.Int).SetUint64(innocentMsg.Height)}
		block := types.NewBlockWithHeader(header)
		msgP := msg(algorithm.Propose, innocentMsg.Height, innocentMsg.Round, algorithm.ValueID(block.Hash()),
			innocentMsg.ValidRound)
		mP, err := EncodeSignedMessage(msgP, b.key, block)
		if err != nil {
			return nil
		}
		return append(msgs, mP)
	}

	// simulate a non proposer node proposing a proposal.
	invalidProposer := func() [][]byte {
		msgP := msg(algorithm.Propose, innocentMsg.Height, innocentMsg.Round, innocentMsg.Value, -1)
		mP, err := EncodeSignedMessage(msgP, b.key, b.msgStore.value(common.Hash(innocentMsg.Value)))
		if err != nil {
			return nil
		}
		return append(msgs, mP)
	}

	// simulate an equivocation over preVote.
	equivocation := func() [][]byte {
		if innocentMsg.MsgType == algorithm.Prevote {
			msgEq := msg(algorithm.Prevote, innocentMsg.Height, innocentMsg.Round, nonNilValue, 0)
			mE, err := EncodeSignedMessage(msgEq, b.key, nil)
			if err != nil {
				return nil
			}
			return append(msgs, mE)
		}
		return nil
	}

	type Rule uint8
	const (
		PN Rule = iota
		PO
		PVN
		PVO
		C
		C1
		InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
		InvalidProposer // A proposal sent from none proposer nodes of the committee.
		Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
		UnknownRule
	)

	r := Rule(*b.MaliciousRuleID)
	if r == PN && innocentMsg.MsgType == algorithm.Propose {
		return maliciousContextPN()
	}

	if r == PO && innocentMsg.MsgType == algorithm.Propose {
		return maliciousContextPO()
	}

	if r == PVN && innocentMsg.MsgType == algorithm.Prevote {
		return maliciousContextPVN()
	}

	if r == C && innocentMsg.MsgType == algorithm.Precommit {
		return maliciousContextC()
	}

	if r == InvalidProposal && innocentMsg.MsgType == algorithm.Propose {
		return invalidProposal()
	}

	if r == InvalidProposer && b.proposer != b.address {
		return invalidProposer()
	}

	if r == Equivocation {
		return equivocation()
	}

	return msgs
}
