package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

// createMisbehaviourContext create misbehaviour context of msgs by according to the configuration and innocent msg.
func (c *core) createMisbehaviourContext(innocentMsg *Message) (msgs [][]byte) {
	nonNilValue := common.Hash{0x1}
	msgPropose := func(block *types.Block, h uint64, r int64, vr int64) *Message {
		proposal := NewProposal(r, new(big.Int).SetUint64(h), vr, block)
		v, err := Encode(proposal)
		if err != nil {
			return nil
		}

		return &Message{
			Code:          msgProposal,
			Msg:           v,
			Address:       c.address,
			CommittedSeal: []byte{},
		}
	}

	msgVote := func(code uint8, h uint64, r int64, v common.Hash) *Message {

		var preVote = Vote{
			Round:             r,
			Height:            new(big.Int).SetUint64(h),
			ProposedBlockHash: v,
		}

		encodedVote, err := Encode(&preVote)
		if err != nil {
			return nil
		}
		return &Message{
			Code:          code,
			Msg:           encodedVote,
			Address:       c.address,
			CommittedSeal: []byte{},
		}
	}

	nextProposeRound := func(currentRound int64) (int64, error) {
		r := currentRound + 1
		for ; ; r++ {
			p := c.committeeSet().GetProposer(r)
			if p.Address == c.address {
				break
			}
		}
		return r, nil
	}

	// simulate a context of msgs that node propose a new proposal rather than the one it locked at previous round.
	maliciousContextPN := func() [][]byte {
		// find a next proposing round.
		nPR, err := nextProposeRound(innocentMsg.R())
		if err != nil {
			return nil
		}
		// simulate a preCommit msg that locked a value at previous round than next proposing round.
		msgEvidence := msgVote(msgPrecommit, innocentMsg.H(), nPR-1, nonNilValue)
		mE, err := c.finalizeMessage(msgEvidence)
		if err != nil {
			return nil
		}

		var proposal Proposal
		err = innocentMsg.Decode(&proposal)
		if err != nil {
			return nil
		}

		// simulate a proposal that propose a new value with -1 as the valid round.
		msgPN := msgPropose(proposal.ProposalBlock, innocentMsg.H(), nPR, -1)
		mPN, err := c.finalizeMessage(msgPN)
		if err != nil {
			return nil
		}
		return append(msgs, mE, mPN)
	}

	// simulate a context of msgs that node propose a proposal that proposed a value for which it is not the one it
	// locked on.
	maliciousContextPO := func() [][]byte {
		// find a next proposing round.
		nPR, err := nextProposeRound(innocentMsg.R())
		if err != nil {
			return nil
		}
		vR := nPR - 1
		// simulate a preCommit msg that locked a value at vR.
		msgEvidence := msgVote(msgPrecommit, innocentMsg.H(), vR, nonNilValue)
		mE, err := c.finalizeMessage(msgEvidence)
		if err != nil {
			return nil
		}

		var proposal Proposal
		err = innocentMsg.Decode(&proposal)
		if err != nil {
			return nil
		}

		// simulate a proposal that node propose for an old value which it is not the one it locked.
		msgPO := msgPropose(proposal.ProposalBlock, innocentMsg.H(), nPR, vR)
		mPO, err := c.finalizeMessage(msgPO)
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
		nPR, err := nextProposeRound(innocentMsg.R())
		if err != nil {
			return nil
		}
		r := nPR - 1
		// simulate a preCommit at round r, for value v1.
		msgEvidence := msgVote(msgPrecommit, innocentMsg.H(), r, nonNilValue)
		mE, err := c.finalizeMessage(msgEvidence)
		if err != nil {
			return nil
		}

		var p Proposal
		err = innocentMsg.Decode(&p)
		if err != nil {
			return nil
		}
		// simulate a proposal at round r+1, for value v2.
		msgProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, -1)
		mP, err := c.finalizeMessage(msgProposal)
		if err != nil {
			return nil
		}
		// simulate a preVote at round r+1, for value v2, this preVote for new value break PVN.
		msgPVN := msgVote(msgPrevote, innocentMsg.H(), nPR, p.GetValue())
		mPVN, err := c.finalizeMessage(msgPVN)
		if err != nil {
			return nil
		}
		return append(msgs, mE, mP, mPVN)
	}

	// simulate a context of msgs that node preCommit at a value V of the round where exist quorum preVotes
	// for not V, in this case, we simulate quorum prevotes for not V, to trigger the fault of breaking of C.
	maliciousContextC := func() [][]byte {
		if innocentMsg.H() == uint64(5) && innocentMsg.R() == 0 {
			msgPV := msgVote(msgPrevote, innocentMsg.H(), innocentMsg.R(), nonNilValue)
			mPV, err := c.finalizeMessage(msgPV)
			if err != nil {
				return nil
			}

			return append(msgs, mPV)
		}
		return msgs
	}

	// simulate an invalid proposal.
	invalidProposal := func() [][]byte {
		header := &types.Header{Number: new(big.Int).SetUint64(innocentMsg.H())}
		block := types.NewBlockWithHeader(header)
		msgP := msgPropose(block, innocentMsg.H(), innocentMsg.R(), innocentMsg.ValidRound())
		mP, err := c.finalizeMessage(msgP)
		if err != nil {
			return nil
		}
		return append(msgs, mP)
	}

	// simulate a non proposer node proposing a proposal.
	invalidProposer := func() [][]byte {
		header := &types.Header{Number: new(big.Int).SetUint64(innocentMsg.H())}
		block := types.NewBlockWithHeader(header)
		msgP := msgPropose(block, innocentMsg.H(), innocentMsg.R(), -1)
		mP, err := c.finalizeMessage(msgP)
		if err != nil {
			return nil
		}
		return append(msgs, mP)
	}

	// simulate an equivocation over preVote.
	equivocation := func() [][]byte {
		if innocentMsg.Code == msgPrevote {
			msgEq := msgVote(msgPrevote, innocentMsg.H(), innocentMsg.R(), nonNilValue)
			mE, err := c.finalizeMessage(msgEq)
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
		GarbageMessage  // message was signed by valid member, but it cannot be decoded.
		InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
		InvalidProposer // A proposal sent from none proposer nodes of the committee.
		Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
		UnknownRule
	)

	r := Rule(c.misbehaviourConfig.MisbehaviourRuleID)
	if r == PN && innocentMsg.Code == msgProposal {
		return maliciousContextPN()
	}

	if r == PO && innocentMsg.Code == msgProposal {
		return maliciousContextPO()
	}

	if r == PVN && innocentMsg.Code == msgProposal {
		return maliciousContextPVN()
	}

	if r == C && innocentMsg.Code == msgPrecommit {
		return maliciousContextC()
	}

	if r == InvalidProposal && innocentMsg.Code == msgProposal {
		return invalidProposal()
	}

	if r == InvalidProposer && c.committeeSet().GetProposer(innocentMsg.R()).Address != c.address {
		return invalidProposer()
	}

	if r == Equivocation {
		return equivocation()
	}

	return msgs
}
