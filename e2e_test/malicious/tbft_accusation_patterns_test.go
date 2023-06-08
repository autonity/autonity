package malicious

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	et "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"math/big"
	"testing"
)

type AccusationRulePOBroadcaster struct {
	*core.Core
}

// simulate an old proposal which refer to less quorum preVotes to trigger the accusation of rule PO
func (s *AccusationRulePOBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	decodedMsg := et.DecodeMsg(msg, s.Core)
	// find a next proposing round.
	nPR := et.NextProposeRound(decodedMsg.R(), s.Core)
	vR := nPR - 1
	var p mUtils.Proposal
	err := msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PO", err)
	}

	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), nPR, vR,
		p.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PO", "error", err)
	}
	invalidProposal := et.NewProposeMsg(s.Address(), p.ProposalBlock, decodedMsg.H(), nPR, vR, liteSig)
	mP, err := s.FinalizeMessage(invalidProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PO", err)
	}
	et.DefaultBehaviour(ctx, s.Core, msg)
	s.Logger().Info("Accusation of PO rule is simulated")
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
}

type AccusationRulePVNBroadcaster struct {
	*core.Core
}

// simulate an accusation context that node preVote for a value that the corresponding proposal is missing.
func (s *AccusationRulePVNBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal || s.IsProposer() == false {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	decodeMsg := et.DecodeMsg(msg, s.Core)
	preVote := et.NewVoteMsg(consensus.MsgPrevote, decodeMsg.H(), decodeMsg.R()+1, et.NonNilValue, s.Core)
	m, err := s.FinalizeMessage(preVote)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVN", err)
	}
	s.Logger().Info("Accusation of PVN rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), m)
}

type AccusationRulePVOBroadcaster struct {
	*core.Core
}

// simulate an accusation context that an old proposal have less quorum preVotes for the value at the valid round.
func (s *AccusationRulePVOBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	decodeMsg := et.DecodeMsg(msg, s.Core)
	// find a next proposing round.
	nPR := et.NextProposeRound(decodeMsg.R(), s.Core)
	// set a valid round.
	validRound := nPR - 2
	if validRound < 0 {
		nPR = et.NextProposeRound(nPR, s.Core)
		validRound = nPR - 2
	}

	// simulate a proposal at round: nPR, and with a valid round: nPR-2
	var p mUtils.Proposal
	err := msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", err)
	}
	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodeMsg.H()), nPR, validRound,
		p.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", "error", err)
	}
	msgProposal := et.NewProposeMsg(s.Address(), p.ProposalBlock, decodeMsg.H(), nPR, validRound, liteSig)
	mP, err := s.FinalizeMessage(msgProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", err)
	}
	// simulate a preVote at round nPR, for value v, this preVote for new value break PVO1.
	msgPVO1 := et.NewVoteMsg(consensus.MsgPrevote, decodeMsg.H(), nPR, p.V(), s.Core)
	mPVO1, err := s.FinalizeMessage(msgPVO1)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", err)
	}
	s.Logger().Info("Accusation of PVO rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mPVO1)
}

type AccusationRuleC1Broadcaster struct {
	*core.Core
}

// simulate an accusation context that node preCommit for a value that have less quorum of preVote for the value.
func (s *AccusationRuleC1Broadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	decodeMsg := et.DecodeMsg(msg, s.Core)
	// find a next proposing round.
	nPR := et.NextProposeRound(decodeMsg.R(), s.Core)
	var p mUtils.Proposal
	err := decodeMsg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule C1", err)
	}

	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodeMsg.H()), nPR, -1,
		p.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule C1", "error", err)
	}
	invalidProposal := et.NewProposeMsg(s.Address(), p.ProposalBlock, decodeMsg.H(), nPR, -1, liteSig)
	mP, err := s.FinalizeMessage(invalidProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule C1", err)
	}

	if s.IsProposer() {
		preCommit := et.NewVoteMsg(consensus.MsgPrecommit, decodeMsg.H(), nPR, p.V(), s.Core)
		m, err := s.FinalizeMessage(preCommit)
		if err != nil {
			s.Logger().Warn("Cannot simulate accusation for rule C1", err)
		}
		s.Logger().Info("Accusation of C1 rule is simulated.")
		_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
		_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), m)
	}
	et.DefaultBehaviour(ctx, s.Core, msg)
}

func TestTBFTAccusationTests(t *testing.T) {
	t.Run("TestTBFTAccusationRulePO", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &AccusationRulePOBroadcaster{}}
		tp := autonity.Accusation
		rule := autonity.PO
		runAccountabilityEventTest(t, handler, tp, rule, 100)
	})
	t.Run("TestTBFTAccusationRulePVN", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &AccusationRulePVNBroadcaster{}}
		tp := autonity.Accusation
		rule := autonity.PVN
		runAccountabilityEventTest(t, handler, tp, rule, 100)
	})
	t.Run("TestTBFTAccusationRulePVO", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &AccusationRulePVOBroadcaster{}}
		tp := autonity.Accusation
		rule := autonity.PVO
		runAccountabilityEventTest(t, handler, tp, rule, 60)
	})
	t.Run("TestTBFTAccusationRuleC1", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &AccusationRuleC1Broadcaster{}}
		tp := autonity.Accusation
		rule := autonity.C1
		runAccountabilityEventTest(t, handler, tp, rule, 60)
	})
}
