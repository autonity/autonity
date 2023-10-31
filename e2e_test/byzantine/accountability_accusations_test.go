package byzantine

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"testing"
)

type AccusationPO struct {
	*core.Core
}

// simulate an old proposal which refer to less quorum preVotes to trigger the accusation of rule PO
func (s *AccusationPO) Broadcast(ctx context.Context, msg message.Message) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(ctx, msg)
		return
	}
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	vR := nPR - 1
	invalidProposal := message.NewPropose(nPR, msg.H(), vR, proposal.Block(), s.Backend().Sign)

	s.Logger().Info("Accusation of PO rule is simulated")
	s.BroadcastAll(ctx, proposal)
	s.BroadcastAll(ctx, invalidProposal)
}

type AccusationPVN struct {
	*core.Core
}

// simulate an accusation context that node preVote for a value that the corresponding proposal is missing.
func (s *AccusationPVN) Broadcast(ctx context.Context, msg message.Message) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal || !s.IsProposer() {
		s.BroadcastAll(ctx, msg)
		return
	}
	preVote := message.NewVote[message.Prevote](msg.R()+1, msg.H(), e2e.NonNilValue, s.Backend().Sign)

	s.Logger().Info("Accusation of PVN rule is simulated.")
	s.BroadcastAll(ctx, proposal)
	s.BroadcastAll(ctx, preVote)
}

type AccusationPVO struct {
	*core.Core
}

// simulate an accusation context that an old proposal have less quorum preVotes for the value at the valid round.
func (s *AccusationPVO) Broadcast(ctx context.Context, msg message.Message) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(ctx, msg)
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
	prevote := message.NewVote[message.Prevote](nPR, msg.H(), proposal.Block().Hash(), s.Backend().Sign)

	s.Logger().Info("Accusation of PVO rule is simulated.")
	s.BroadcastAll(ctx, proposal)
	s.BroadcastAll(ctx, newProposal)
	s.BroadcastAll(ctx, prevote)
}

type AccusationC1 struct {
	*core.Core
}

// simulate an accusation context that node preCommit for a value that have less quorum of preVote for the value.
func (s *AccusationC1) Broadcast(ctx context.Context, msg message.Message) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(ctx, msg)
		return
	}
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	if s.IsProposer() { // youssef: probably not needed
		preCommit := message.NewVote[message.Precommit](nPR, msg.H(), proposal.Block().Hash(), s.Backend().Sign)
		s.Logger().Info("Accusation of C1 rule is simulated.")
		s.BroadcastAll(ctx, preCommit)
	}
	s.BroadcastAll(ctx, proposal)
}

func TestAccusationFlow(t *testing.T) {
	t.Run("AccusationRulePO", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: &AccusationPO{}}
		tp := autonity.Accusation
		rule := autonity.PO
		runTest(t, handler, tp, rule, 100)
	})
	t.Run("AccusationRulePVN", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: &AccusationPVN{}}
		tp := autonity.Accusation
		rule := autonity.PVN
		runTest(t, handler, tp, rule, 100)
	})
	/*
		Not supported, require more complicated setup
		we need to be able to handle more than one byzantine validator
		t.Run("AccusationRulePVO", func(t *testing.T) {
			handler := &node.TendermintServices{Broadcaster: &AccusationPVO{}}
			tp := autonity.Accusation
			rule := autonity.PVO
			runTest(t, handler, tp, rule, 60)
		})
	*/
	t.Run("AccusationRuleC1", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: &AccusationC1{}}
		tp := autonity.Accusation
		rule := autonity.C1
		runTest(t, handler, tp, rule, 60)
	})
}
