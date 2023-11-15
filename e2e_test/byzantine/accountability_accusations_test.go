package byzantine

import (
	"context"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
)

func newAccusationRulePOBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &AccusationRulePOBroadcaster{c.(*core.Core)}
}

type AccusationRulePOBroadcaster struct {
	*core.Core
}

// simulate an old proposal which refer to less quorum preVotes to trigger the accusation of rule PO
func (s *AccusationRulePOBroadcaster) SignAndBroadcast(msg *message.Message) {
	if msg.Code != consensus.MsgProposal {
		e2e.DefaultSignAndBroadcast(s.Core, msg)
		return
	}
	_ = msg.DecodePayload()
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	vR := nPR - 1
	var p message.Proposal
	err := msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PO", err)
	}
	invalidProposal := e2e.NewProposeMsg(s.Address(), p.ProposalBlock, msg.H(), nPR, vR, s.Backend().Sign)
	mP, err := s.SignMessage(invalidProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PO", err)
	}
	e2e.DefaultSignAndBroadcast(s.Core, msg)
	s.Logger().Info("Accusation of PO rule is simulated")
	_ = s.Backend().Broadcast(s.CommitteeSet().Committee(), mP)
}

func newAccusationRulePVNBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &AccusationRulePVNBroadcaster{c.(*core.Core)}
}

type AccusationRulePVNBroadcaster struct {
	*core.Core
}

// simulate an accusation context that node preVote for a value that the corresponding proposal is missing.
func (s *AccusationRulePVNBroadcaster) SignAndBroadcast(msg *message.Message) {
	if msg.Code != consensus.MsgProposal || s.IsProposer() == false {
		e2e.DefaultSignAndBroadcast(s.Core, msg)
		return
	}
	_ = msg.DecodePayload()
	preVote := e2e.NewVoteMsg(consensus.MsgPrevote, msg.H(), msg.R()+1, e2e.NonNilValue, s.Core)
	m, err := s.SignMessage(preVote)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVN", err)
	}
	s.Logger().Info("Accusation of PVN rule is simulated.")
	e2e.DefaultSignAndBroadcast(s.Core, msg)
	_ = s.Backend().Broadcast(s.CommitteeSet().Committee(), m)
}

/* not currently supported
func newAccusationRulePVOBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &AccusationRulePVOBroadcaster{c.(*core.Core)}
}
*/

type AccusationRulePVOBroadcaster struct {
	*core.Core
}

// simulate an accusation context that an old proposal have less quorum preVotes for the value at the valid round.
func (s *AccusationRulePVOBroadcaster) SignAndBroadcast(msg *message.Message) {
	if msg.Code != consensus.MsgProposal {
		e2e.DefaultSignAndBroadcast(s.Core, msg)
		return
	}
	_ = msg.DecodePayload()
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	// set a valid round.
	validRound := nPR - 2
	if validRound < 0 {
		nPR = e2e.NextProposeRound(nPR, s.Core)
		validRound = nPR - 2
	}

	// simulate a proposal at round: nPR, and with a valid round: nPR-2
	var p message.Proposal
	err := msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", err)
	}
	msgProposal := e2e.NewProposeMsg(s.Address(), p.ProposalBlock, msg.H(), nPR, validRound, s.Core.Backend().Sign)
	mP, err := s.SignMessage(msgProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", err)
	}
	// simulate a preVote at round nPR, for value v, this preVote for new value break PVO1.
	msgPVO1 := e2e.NewVoteMsg(consensus.MsgPrevote, msg.H(), nPR, p.V(), s.Core)
	mPVO1, err := s.SignMessage(msgPVO1)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule PVO", err)
	}
	s.Logger().Info("Accusation of PVO rule is simulated.")
	e2e.DefaultSignAndBroadcast(s.Core, msg)
	_ = s.Backend().Broadcast(s.CommitteeSet().Committee(), mP)
	_ = s.Backend().Broadcast(s.CommitteeSet().Committee(), mPVO1)
}

func newAccusationRuleC1Broadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &AccusationRuleC1Broadcaster{c.(*core.Core)}
}

type AccusationRuleC1Broadcaster struct {
	*core.Core
}

// simulate an accusation context that node preCommit for a value that have less quorum of preVote for the value.
func (s *AccusationRuleC1Broadcaster) SignAndBroadcast(msg *message.Message) {
	if msg.Code != consensus.MsgProposal {
		e2e.DefaultSignAndBroadcast(s.Core, msg)
		return
	}
	_ = msg.DecodePayload()
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	var p message.Proposal
	err := msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate accusation for rule C1", err)
	}

	if s.IsProposer() {
		preCommit := e2e.NewVoteMsg(consensus.MsgPrecommit, msg.H(), nPR, p.V(), s.Core)
		m, err := s.SignMessage(preCommit)
		if err != nil {
			s.Logger().Warn("Cannot simulate accusation for rule C1", err)
		}
		s.Logger().Info("Accusation of C1 rule is simulated.")
		_ = s.Backend().Broadcast(s.CommitteeSet().Committee(), m)
	}
	e2e.DefaultSignAndBroadcast(s.Core, msg)
}

func TestAccusationFlow(t *testing.T) {
	t.Run("AccusationRulePO", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: newAccusationRulePOBroadcaster}
		tp := autonity.Accusation
		rule := autonity.PO
		runTest(t, handler, tp, rule, 100)
	})
	t.Run("AccusationRulePVN", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: newAccusationRulePVNBroadcaster}
		tp := autonity.Accusation
		rule := autonity.PVN
		runTest(t, handler, tp, rule, 100)
	})
	/* // Not supported, require more complicated setup
	t.Run("AccusationRulePVO", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: newAccusationRulePVOBroadcaster}
		tp := autonity.Accusation
		rule := autonity.PVO
		runTest(t, handler, tp, rule, 60)
	})*/
	t.Run("AccusationRuleC1", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: newAccusationRuleC1Broadcaster}
		tp := autonity.Accusation
		rule := autonity.C1
		runTest(t, handler, tp, rule, 60)
	})
}
