package core

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
)

type FaultSimulatorConfig struct {
	MisbehaviourRuleID *uint8
	AccusationRuleID   *uint8
}

var FaultSimulatorConfigs *FaultSimulatorConfig

func SetFaultSimulatorConfig(t uint8, rule uint8) {
	FaultSimulatorConfigs = new(FaultSimulatorConfig)
	if t == autonity.Misbehaviour {
		FaultSimulatorConfigs.MisbehaviourRuleID = new(uint8)
		*FaultSimulatorConfigs.MisbehaviourRuleID = rule
	}
	if t == autonity.Accusation {
		FaultSimulatorConfigs.AccusationRuleID = new(uint8)
		*FaultSimulatorConfigs.AccusationRuleID = rule
	}
}

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

	nextProposeRound := func(currentRound int64) int64 {
		r := currentRound + 1
		for ; ; r++ {
			p := c.committeeSet().GetProposer(r)
			if p.Address == c.address {
				break
			}
		}
		return r
	}

	// simulate a context of msgs that node propose a new proposal rather than the one it locked at previous round.
	maliciousContextPN := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		// simulate a preCommit msg that locked a value at previous round than next proposing round.
		msgEvidence := msgVote(msgPrecommit, innocentMsg.H(), nPR-1, nonNilValue)
		mE, err := c.finalizeMessage(msgEvidence)
		if err != nil {
			return msgs
		}

		var proposal Proposal
		err = innocentMsg.Decode(&proposal)
		if err != nil {
			return msgs
		}

		// simulate a proposal that propose a new value with -1 as the valid round.
		msgPN := msgPropose(proposal.ProposalBlock, innocentMsg.H(), nPR, -1)
		mPN, err := c.finalizeMessage(msgPN)
		if err != nil {
			return msgs
		}
		c.logger.Info("Misbehaviour of PN rule is simulated.")
		return append(msgs, mE, mPN)
	}

	// simulate a context of msgs that node propose a proposal that proposed a value for which it is not the one it
	// locked on.
	maliciousContextPO := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		vR := nPR - 1
		// simulate a preCommit msg that locked a value at vR.
		msgEvidence := msgVote(msgPrecommit, innocentMsg.H(), vR, nonNilValue)
		mE, err := c.finalizeMessage(msgEvidence)
		if err != nil {
			return msgs
		}

		var proposal Proposal
		err = innocentMsg.Decode(&proposal)
		if err != nil {
			return msgs
		}

		// simulate a proposal that node propose for an old value which it is not the one it locked.
		msgPO := msgPropose(proposal.ProposalBlock, innocentMsg.H(), nPR, vR)
		mPO, err := c.finalizeMessage(msgPO)
		if err != nil {
			return msgs
		}
		c.logger.Info("Misbehaviour of PO rule is simulated.")
		return append(msgs, mE, mPO)
	}

	// simulate a context of msgs that a node preVote for a new value rather than the one it locked on.
	// preCommit (h, r, v1)
	// propose   (h, r+1, v2)
	// preVote   (h, r+1, v2)
	maliciousContextPVN := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		r := nPR - 1
		// simulate a preCommit at round r, for value v1.
		msgEvidence := msgVote(msgPrecommit, innocentMsg.H(), r, nonNilValue)
		mE, err := c.finalizeMessage(msgEvidence)
		if err != nil {
			return msgs
		}

		var p Proposal
		err = innocentMsg.Decode(&p)
		if err != nil {
			return msgs
		}
		// simulate a proposal at round r+1, for value v2.
		msgProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, -1)
		mP, err := c.finalizeMessage(msgProposal)
		if err != nil {
			return msgs
		}
		// simulate a preVote at round r+1, for value v2, this preVote for new value break PVN.
		msgPVN := msgVote(msgPrevote, innocentMsg.H(), nPR, p.GetValue())
		mPVN, err := c.finalizeMessage(msgPVN)
		if err != nil {
			return msgs
		}
		c.logger.Info("Misbehaviour of PVN rule is simulated.")
		return append(msgs, mE, mP, mPVN)
	}

	// simulate a context of msgs that a node preVote for a value that is not the one it precommitted at previous round.
	// create a proposal: (h, r:3, vr: 0, with v.)
	// preCommit (h, r:0, v)
	// proCommit (h, r:1, v)
	// preCommit (h, r:2, not v)
	// preVote   (h, r:3, v)
	maliciousContextPVO1 := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		// set a valid round.
		currentRound := nPR
		validRound := nPR - 2
		if validRound < 0 {
			nPR = nextProposeRound(nPR)
			currentRound = nPR
			validRound = nPR - 2
		}

		// simulate a proposal at round: nPR, and with a valid round: nPR-2
		var p Proposal
		err := innocentMsg.Decode(&p)
		if err != nil {
			return msgs
		}

		msgProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, validRound)
		mP, err := c.finalizeMessage(msgProposal)
		if err != nil {
			return msgs
		}

		// simulate preCommits at each round between [validRound, current)
		var messages [][]byte
		for i := validRound; i < currentRound; i++ {
			if i == currentRound-1 {
				msgPC := msgVote(msgPrecommit, innocentMsg.H(), i, nonNilValue)
				mPC, err := c.finalizeMessage(msgPC)
				if err != nil {
					return msgs
				}
				messages = append(messages, mPC)
			} else {
				msgPC := msgVote(msgPrecommit, innocentMsg.H(), i, p.GetValue())
				mPC, err := c.finalizeMessage(msgPC)
				if err != nil {
					return msgs
				}
				messages = append(messages, mPC)
			}
		}

		// simulate a preVote at round 3, for value v, this preVote for new value break PVO1.
		msgPVO1 := msgVote(msgPrevote, innocentMsg.H(), nPR, p.GetValue())
		mPVO1, err := c.finalizeMessage(msgPVO1)
		if err != nil {
			return msgs
		}
		c.logger.Info("Misbehaviour of PVO1 rule is simulated.")
		return append(append(msgs, mP, mPVO1), messages...)
	}

	// simulate a context of msgs that a node preVote for a value that is not the one it precommitted at previous round.
	// create a proposal: (h, r:3, vr: 0, with v.)
	// preCommit (h, r:0, not v)
	// proCommit (h, r:1, not v)
	// preCommit (h, r:2, not v)
	// preVote   (h, r:3, v)
	maliciousContextPVO2 := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())

		// set a valid round.
		currentRound := nPR
		validRound := nPR - 2
		if validRound < 0 {
			nPR = nextProposeRound(nPR)
			currentRound = nPR
			validRound = nPR - 2
		}

		// simulate a proposal at round: nPR, and with a valid round: nPR-2
		var p Proposal
		err := innocentMsg.Decode(&p)
		if err != nil {
			return msgs
		}

		msgProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, validRound)
		mP, err := c.finalizeMessage(msgProposal)
		if err != nil {
			return msgs
		}

		// simulate preCommits of not V at each round between [validRound, current)
		var messages [][]byte
		for i := validRound; i < currentRound; i++ {
			msgPC := msgVote(msgPrecommit, innocentMsg.H(), i, nonNilValue)
			mPC, err := c.finalizeMessage(msgPC)
			if err != nil {
				return msgs
			}
			messages = append(messages, mPC)
		}

		// simulate a preVote at round 3, for value v, this preVote for new value break PVO2.
		msgPVO2 := msgVote(msgPrevote, innocentMsg.H(), nPR, p.GetValue())
		mPVO2, err := c.finalizeMessage(msgPVO2)
		if err != nil {
			return msgs
		}

		c.logger.Info("Misbehaviour of PVO2 rule is simulated.")
		return append(append(msgs, messages...), mP, mPVO2)
	}

	// simulate a context of msgs that node preCommit at a value V of the round where exist quorum preVotes
	// for not V, in this case, we simulate quorum prevotes for not V, to trigger the fault of breaking of C.
	maliciousContextC := func() [][]byte {
		if innocentMsg.H() == uint64(5) && innocentMsg.R() == 0 {
			msgPV := msgVote(msgPrevote, innocentMsg.H(), innocentMsg.R(), nonNilValue)
			mPV, err := c.finalizeMessage(msgPV)
			if err != nil {
				return msgs
			}

			c.logger.Info("Misbehaviour of C rule is simulated.")
			return append(msgs, mPV)
		}
		return msgs
	}

	// simulate an invalid proposal.
	invalidProposal := func() [][]byte {
		// send an invalid proposal happens at height 3.
		if innocentMsg.H() == uint64(3) {
			nextPR := nextProposeRound(innocentMsg.R())
			header := &types.Header{Number: new(big.Int).SetUint64(innocentMsg.H())}
			block := types.NewBlockWithHeader(header)
			msgP := msgPropose(block, innocentMsg.H(), nextPR, innocentMsg.ValidRound())
			mP, err := c.finalizeMessage(msgP)
			if err != nil {
				return msgs
			}
			c.logger.Info("Misbehaviour of invalid proposal rule is simulated.")
			return append(msgs, mP)
		}
		return msgs
	}

	// simulate a non proposer node proposing a proposal.
	invalidProposer := func() [][]byte {
		header := &types.Header{Number: new(big.Int).SetUint64(innocentMsg.H())}
		block := types.NewBlockWithHeader(header)
		msgP := msgPropose(block, innocentMsg.H(), innocentMsg.R(), -1)
		mP, err := c.finalizeMessage(msgP)
		if err != nil {
			return msgs
		}
		c.logger.Info("Misbehaviour of invalid proposer rule is simulated.")
		return append(msgs, mP)
	}

	// simulate an equivocation over preVote.
	equivocation := func() [][]byte {
		// let proposer of the round send equivocated preVote.
		if innocentMsg.Code == msgPrevote && c.isProposer() {
			msgEq := msgVote(msgPrevote, innocentMsg.H(), innocentMsg.R(), nonNilValue)
			mE, err := c.finalizeMessage(msgEq)
			if err != nil {
				return msgs
			}
			c.logger.Info("Misbehaviour of equivocation rule is simulated.")
			return append(msgs, mE)
		}
		return nil
	}

	// simulate an old proposal which refer to less quorum preVotes to trigger the accusation of rule PO
	accusationContextPO := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		vR := nPR - 1
		var p Proposal
		err := innocentMsg.Decode(&p)
		if err != nil {
			return msgs
		}

		invalidProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, vR)
		mP, err := c.finalizeMessage(invalidProposal)
		if err != nil {
			return msgs
		}
		c.logger.Info("Accusation of PO rule is simulated.")
		return append(msgs, mP)
	}

	// simulate an accusation context that node preVote for a value that the corresponding proposal is missing.
	accusationContextPVN := func() [][]byte {
		// let the proposer at height 3, send the preVote for a value that the corresponding proposal is missing.
		if c.isProposer() && innocentMsg.H() == uint64(3) {
			preVote := msgVote(msgPrevote, innocentMsg.H(), innocentMsg.R()+1, nonNilValue)
			m, err := c.finalizeMessage(preVote)
			if err != nil {
				return msgs
			}
			c.logger.Info("Accusation of PVN rule is simulated.")
			return append(msgs, m)
		}
		return msgs
	}

	// simulate an accusation context that an old proposal have less quorum preVotes for the value at the valid round.
	accusationContextPVO := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		// set a valid round.
		validRound := nPR - 2
		if validRound < 0 {
			nPR = nextProposeRound(nPR)
			validRound = nPR - 2
		}

		// simulate a proposal at round: nPR, and with a valid round: nPR-2
		var p Proposal
		err := innocentMsg.Decode(&p)
		if err != nil {
			return msgs
		}

		msgProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, validRound)
		mP, err := c.finalizeMessage(msgProposal)
		if err != nil {
			return msgs
		}

		// simulate a preVote at round 3, for value v, this preVote for new value break PVO1.
		msgPVO1 := msgVote(msgPrevote, innocentMsg.H(), nPR, p.GetValue())
		mPVO1, err := c.finalizeMessage(msgPVO1)
		if err != nil {
			return msgs
		}

		c.logger.Info("Accusation of PVO rule is simulated.")
		return append(msgs, mP, mPVO1)
	}

	// simulate an accusation context that node preCommit for a value that the corresponding proposal is missing.
	accusationContextC := func() [][]byte {
		preCommit := msgVote(msgPrecommit, innocentMsg.H(), innocentMsg.R(), nonNilValue)
		m, err := c.finalizeMessage(preCommit)
		if err != nil {
			return msgs
		}
		c.logger.Info("Accusation of C rule is simulated.")
		return append(msgs, m)
	}

	// simulate an accusation context that node preCommit for a value that have less quorum of preVote for the value.
	accusationContextC1 := func() [][]byte {
		// find a next proposing round.
		nPR := nextProposeRound(innocentMsg.R())
		var p Proposal
		err := innocentMsg.Decode(&p)
		if err != nil {
			return msgs
		}
		invalidProposal := msgPropose(p.ProposalBlock, innocentMsg.H(), nPR, -1)
		mP, err := c.finalizeMessage(invalidProposal)
		if err != nil {
			return msgs
		}

		if c.isProposer() {
			preCommit := msgVote(msgPrecommit, innocentMsg.H(), nPR, p.GetValue())
			m, err := c.finalizeMessage(preCommit)
			if err != nil {
				return msgs
			}
			c.logger.Info("Accusation of C1 rule is simulated.")
			return append(msgs, mP, m)
		}
		return msgs
	}

	type Rule uint8
	const (
		PN Rule = iota
		PO
		PVN
		PVO
		PVO1
		PVO2
		C
		C1
		GarbageMessage  // message was signed by valid member, but it cannot be decoded.
		InvalidProposal // The value proposed by proposer cannot pass the blockchain's validation.
		InvalidProposer // A proposal sent from none proposer nodes of the committee.
		Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
		UnknownRule
	)
	if FaultSimulatorConfigs.MisbehaviourRuleID != nil {
		r := Rule(*FaultSimulatorConfigs.MisbehaviourRuleID)
		if r == PN && innocentMsg.Code == msgProposal {
			return maliciousContextPN()
		}

		if r == PO && innocentMsg.Code == msgProposal {
			return maliciousContextPO()
		}

		if r == PVN && innocentMsg.Code == msgProposal {
			return maliciousContextPVN()
		}

		if r == PVO1 && innocentMsg.Code == msgProposal {
			return maliciousContextPVO1()
		}

		if r == PVO2 && innocentMsg.Code == msgProposal {
			return maliciousContextPVO2()
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
	}

	if FaultSimulatorConfigs.AccusationRuleID != nil {
		r := Rule(*FaultSimulatorConfigs.AccusationRuleID)
		if r == PO && innocentMsg.Code == msgProposal {
			return accusationContextPO()
		}

		if r == PVN && innocentMsg.Code == msgPrevote {
			return accusationContextPVN()
		}

		if r == PVO && innocentMsg.Code == msgProposal {
			return accusationContextPVO()
		}

		if r == C && innocentMsg.Code == msgPrecommit {
			return accusationContextC()
		}

		if r == C1 && innocentMsg.Code == msgProposal {
			return accusationContextC1()
		}
	}
	return msgs
}
