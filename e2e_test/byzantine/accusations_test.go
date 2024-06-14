package byzantine

import (
	"math/rand"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
)

func selfAndCsize(c *core.Core, h uint64) (*types.CommitteeMember, int) {
	committee, err := c.Backend().BlockChain().CommitteeOfHeight(h)
	if err != nil {
		panic(err)
	}

	return committee.MemberByAddress(c.Address()), committee.Len()
}

type AccusationPO struct {
	*core.Core
}

func newAccusationPO(c interfaces.Core) interfaces.Broadcaster {
	return &AccusationPO{c.(*core.Core)}
}

// simulate an old proposal not backed by a quorum preVotes to trigger the accusation of rule PO
func (s *AccusationPO) Broadcast(msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(msg)
		return
	}
	// find a next proposing round.
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	vR := nPR - 1

	self, _ := selfAndCsize(s.Core, msg.H())

	// change header nonce to a random value to have a different block hash
	header := proposal.Block().Header()
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256)) //nolint
	}
	header.Nonce = nonce
	block := types.NewBlockWithHeader(header)
	invalidProposal := message.NewPropose(nPR, msg.H(), vR, block, s.Backend().Sign, self)

	s.Logger().Info("PO Accusation rule simulation")
	s.BroadcastAll(proposal)
	s.BroadcastAll(invalidProposal)
}

type AccusationPVN struct {
	*core.Core
}

func newAccusationPVN(c interfaces.Core) interfaces.Broadcaster {
	return &AccusationPVN{c.(*core.Core)}
}

// simulate an accusation where the node preVotes for a value but the corresponding proposal is missing.
func (s *AccusationPVN) Broadcast(msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal || !s.IsProposer() {
		s.BroadcastAll(msg)
		return
	}
	self, csize := selfAndCsize(s.Core, msg.H())
	preVote := message.NewPrevote(msg.R()+1, msg.H(), e2e.NonNilValue, s.Backend().Sign, self, csize)

	s.Logger().Info("PVN Accusation rule simulation")
	s.BroadcastAll(proposal)
	s.BroadcastAll(preVote)
}

type AccusationPVO struct {
	*core.Core
}

/*
To be uncommented when PVO is fixed
func newAccusationPVO(c interfaces.Core) interfaces.Broadcaster {
	return &AccusationPVO{c.(*core.Core)}
}
*/

// simulate an accusation context that an old proposal have less quorum preVotes for the value at the valid round.
func (s *AccusationPVO) Broadcast(msg message.Msg) {
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
	self, csize := selfAndCsize(s.Core, msg.H())
	newProposal := message.NewPropose(nPR, msg.H(), validRound, proposal.Block(), s.Backend().Sign, self)

	// simulate a preVote at round nPR, for value v, this preVote for new value break PVO1.
	prevote := message.NewPrevote(nPR, msg.H(), proposal.Block().Hash(), s.Backend().Sign, self, csize)

	s.Logger().Info("PVO accusation rule simulation")
	s.BroadcastAll(proposal)
	s.BroadcastAll(newProposal)
	s.BroadcastAll(prevote)
}

type AccusationC1 struct {
	*core.Core
}

func newAccusationC1(c interfaces.Core) interfaces.Broadcaster {
	return &AccusationC1{c.(*core.Core)}
}

// simulate an accusation context that node preCommit for a value that have less quorum of preVote for the value.
func (s *AccusationC1) Broadcast(msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(msg)
		return
	}

	self, csize := selfAndCsize(s.Core, msg.H())

	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	if s.IsProposer() { // youssef: probably not needed
		preCommit := message.NewPrecommit(nPR, msg.H(), common.Hash{0xca, 0xfe}, s.Backend().Sign, self, csize)
		s.Logger().Info("C1 accusation rule simulation")
		s.BroadcastAll(preCommit)
	}
	s.BroadcastAll(proposal)
}

func TestAccusationFlow(t *testing.T) {
	t.Run("AccusationRulePO", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newAccusationPO}
		tp := autonity.Accusation
		rule := autonity.PO
		runTest(t, handler, tp, rule, 100)
	})
	t.Run("AccusationRulePVN", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newAccusationPVN}
		tp := autonity.Accusation
		rule := autonity.PVN
		runTest(t, handler, tp, rule, 100)
	})
	t.Run("AccusationRuleC1", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newAccusationC1}
		tp := autonity.Accusation
		rule := autonity.C1
		runTest(t, handler, tp, rule, 60)
	})
}
