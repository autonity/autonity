package byzantine

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/e2e_test"
	"testing"
)

type AccusationPO struct {
	*core.Core
}

// simulate an old proposal which refer to less quorum preVotes to trigger the accusation of rule PO
func (s *AccusationPO) Broadcast(ctx context.Context, msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(msg)
		return
	}
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	vR := nPR - 1
	invalidProposal := message.NewPropose(nPR, msg.H(), vR, proposal.Block(), s.Backend().Sign)

	s.Logger().Info("PO Accusation rule simulation")
	s.BroadcastAll(proposal)
	s.BroadcastAll(invalidProposal)
}

type AccusationPVN struct {
	*core.Core
}

// simulate an accusation context that node preVote for a value that the corresponding proposal is missing.
func (s *AccusationPVN) Broadcast(ctx context.Context, msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal || !s.IsProposer() {
		s.BroadcastAll(msg)
		return
	}
	preVote := message.NewPrevote(msg.R()+1, msg.H(), e2e.NonNilValue, s.Backend().Sign)

	s.Logger().Info("PVN Accusation rule simulation")
	s.BroadcastAll(proposal)
	s.BroadcastAll(preVote)
}

type AccusationPVO struct {
	*core.Core
}

// simulate an accusation context that an old proposal have less quorum preVotes for the value at the valid round.
func (s *AccusationPVO) Broadcast(ctx context.Context, msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(msg)
		return
	}
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	// set a valid round.
	validRound := nPR - 2
	if validRound < 0 {
		nPR = e2e.NextProposeRound(nPR, s.Core)
		validRound = nPR - 2
	}

	// simulate a proposal at round: nPR, and with a valid round: nPR-2
	newProposal := message.NewPropose(nPR, msg.H(), validRound, proposal.Block(), s.Backend().Sign)

	// simulate a preVote at round nPR, for value v, this preVote for new value break PVO1.
	prevote := message.NewPrevote(nPR, msg.H(), proposal.Block().Hash(), s.Backend().Sign)

	s.Logger().Info("PVO accusation rule simulation")
	s.BroadcastAll(proposal)
	s.BroadcastAll(newProposal)
	s.BroadcastAll(prevote)
}

type AccusationC1 struct {
	*core.Core
}

// simulate an accusation context that node preCommit for a value that have less quorum of preVote for the value.
func (s *AccusationC1) Broadcast(ctx context.Context, msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(msg)
		return
	}
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	if s.IsProposer() { // youssef: probably not needed
		preCommit := message.NewPrecommit(nPR, msg.H(), proposal.Block().Hash(), s.Backend().Sign)
		s.Logger().Info("C1 accusation rule simulation ,.")
		s.BroadcastAll(preCommit)
	}
	s.BroadcastAll(proposal)
}

func TestAccusationFlow(t *testing.T) {
	t.Run("AccusationRulePO", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: &AccusationPO{}}
		tp := autonity.Accusation
		rule := autonity.PO
		runTest(t, handler, tp, rule, 100)
	})
	t.Run("AccusationRulePVN", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: &AccusationPVN{}}
		tp := autonity.Accusation
		rule := autonity.PVN
		runTest(t, handler, tp, rule, 100)
	})
	/*
		Not supported, require more complicated setup
		we need to be able to handle more than one byzantine validator
		t.Run("AccusationRulePVO", func(t *testing.T) {
			handler := &interfaces.Services{Broadcaster: &AccusationPVO{}}
			tp := autonity.Accusation
			rule := autonity.PVO
			runTest(t, handler, tp, rule, 60)
		})
	*/
	t.Run("AccusationRuleC1", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: &AccusationC1{}}
		tp := autonity.Accusation
		rule := autonity.C1
		runTest(t, handler, tp, rule, 60)
	})
}
