package malicious

import (
	"context"
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	et "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
)

func runAccountabilityEventTest(t *testing.T, handler *node.CustomHandler, tp autonity.AccountabilityEventType,
	rule autonity.Rule, testPeriod uint64) {

	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	users, err := test.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set Malicious users
	faultyNode := 0
	users[faultyNode].CustHandler = handler
	// creates a network of 4 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(testPeriod, 500) // nolint

	// check if the misbehaviour is presented for faulty node #0
	faultyAddress := network[faultyNode].Address
	detected := et.AccountabilityEventDetected(t, faultyAddress, tp, rule, network)
	require.Equal(t, true, detected)
}

type MisbehaviourRulePNBroadcaster struct {
	*core.Core
}

// simulate a context of msgs that node proposes a new proposal rather than the one it locked at previous rounds.
func (s *MisbehaviourRulePNBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}

	// start to simulate malicious context to break rule PN.
	decodedMsg := et.DecodeMsg(msg, s.Core)
	nPR := et.NextProposeRound(decodedMsg.R(), s.Core)

	// simulate a preCommit msg that locked a value at previous round than next proposing round.
	msgEvidence := et.NewVoteMsg(consensus.MsgPrecommit, decodedMsg.H(), nPR-1, et.NonNilValue, s.Core)
	mE, err := s.FinalizeMessage(msgEvidence)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PN rule.", err)
	}

	var proposal mUtils.Proposal
	err = msg.Decode(&proposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PN rule.", err)
	}

	// simulate a proposal that propose a new value with -1 as the valid round.
	//msgPN := et.NewProposeMsg(proposal.ProposalBlock, decodedMsg.H(), nPR, -1, s.Core)
	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), nPR, -1,
		proposal.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PN rule", err)
	}
	msgPN := et.NewProposeMsg(s.Core.Address(), proposal.ProposalBlock, decodedMsg.H(), nPR, -1, liteSig)

	mPN, err := s.FinalizeMessage(msgPN)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PN rule.", err)
	}
	s.Logger().Info("Misbehaviour of PN rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mE)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mPN)
}

type MisbehaviourRulePOBroadcaster struct {
	*core.Core
}

// simulate a context of msgs that node proposes a value for which was not the one it locked on.
func (s *MisbehaviourRulePOBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	// start to simulate malicious context to break rule PO.
	decodedMsg := et.DecodeMsg(msg, s.Core)
	nPR := et.NextProposeRound(decodedMsg.R(), s.Core)
	vR := nPR - 1
	// simulate a preCommit msg that locked a value at vR.
	msgEvidence := et.NewVoteMsg(consensus.MsgPrecommit, decodedMsg.H(), vR, et.NonNilValue, s.Core)
	mE, err := s.FinalizeMessage(msgEvidence)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PO rule.", err)
	}

	var proposal mUtils.Proposal
	err = msg.Decode(&proposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PO rule.", err)
	}

	// simulate a proposal that node propose for an old value which it is not the one it locked.
	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), nPR, vR,
		proposal.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PO rule", err)
	}
	msgPO := et.NewProposeMsg(s.Address(), proposal.ProposalBlock, decodedMsg.H(), nPR, vR, liteSig)
	mPO, err := s.FinalizeMessage(msgPO)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PO rule.", err)
	}

	s.Logger().Info("Misbehaviour of PN rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)

	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mE)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mPO)
}

type MisbehaviourRulePVNBroadcaster struct {
	*core.Core
}

// simulate a context of msgs that a node preVote for a new value rather than the one it locked on.
// An example context like below:
// preCommit (h, r, v1)
// propose   (h, r+1, v2)
// preVote   (h, r+1, v2)
func (s *MisbehaviourRulePVNBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}

	decodedMsg := et.DecodeMsg(msg, s.Core)
	// find a next proposing round.
	nPR := et.NextProposeRound(decodedMsg.R(), s.Core)
	r := nPR - 1
	// simulate a preCommit at round r, for value v1.
	msgEvidence := et.NewVoteMsg(consensus.MsgPrecommit, decodedMsg.H(), r, et.NonNilValue, s.Core)
	mE, err := s.FinalizeMessage(msgEvidence)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVN rule.", err)
	}

	var p mUtils.Proposal
	err = msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVN rule.", err)
	}

	// simulate a proposal at round r+1, for value v2.
	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), nPR, -1,
		p.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVN rule", err)
	}
	msgProposal := et.NewProposeMsg(s.Address(), p.ProposalBlock, decodedMsg.H(), nPR, -1, liteSig)
	mP, err := s.FinalizeMessage(msgProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVN rule.", err)
	}

	// simulate a preVote at round r+1, for value v2, this preVote for new value break PVN.
	msgPVN := et.NewVoteMsg(consensus.MsgPrevote, decodedMsg.H(), nPR, p.V(), s.Core)
	mPVN, err := s.FinalizeMessage(msgPVN)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVN rule.", err)
	}

	s.Logger().Info("Misbehaviour of PVN rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mE)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mPVN)
}

type MisbehaviourRulePVO12Broadcaster struct {
	*core.Core
}

// simulate a context of msgs that a node preVote for a value that is not the one it precommitted at previous round.
// An example context like below:
// create a proposal: (h, r:3, vr: 0, with v.)
// preCommit (h, r:0, v)
// proCommit (h, r:1, v)
// preCommit (h, r:2, not v)
// preVote   (h, r:3, v)
func (s *MisbehaviourRulePVO12Broadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}

	decodedMsg := et.DecodeMsg(msg, s.Core)
	// find a next proposing round.
	nPR := et.NextProposeRound(decodedMsg.R(), s.Core)
	// set a valid round.
	currentRound := nPR
	validRound := nPR - 2
	if validRound < 0 {
		nPR = et.NextProposeRound(nPR, s.Core)
		currentRound = nPR
		validRound = nPR - 2
	}

	var p mUtils.Proposal
	err := msg.Decode(&p)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVO12 rule.", err)
	}

	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), nPR, validRound,
		p.ProposalBlock.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PV01 rule", err)
	}
	msgProposal := et.NewProposeMsg(s.Address(), p.ProposalBlock, decodedMsg.H(), nPR, validRound, liteSig)
	mP, err := s.FinalizeMessage(msgProposal)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVO12 rule.", err)
	}
	// simulate preCommits at each round between [validRound, current)
	var messages [][]byte
	for i := validRound; i < currentRound; i++ {
		if i == currentRound-1 {
			msgPC := et.NewVoteMsg(consensus.MsgPrecommit, decodedMsg.H(), i, et.NonNilValue, s.Core)
			mPC, err := s.FinalizeMessage(msgPC)
			if err != nil {
				s.Logger().Warn("Cannot simulate Misbehaviour of PVO12 rule.", err)
			}
			messages = append(messages, mPC)
		} else {
			msgPC := et.NewVoteMsg(consensus.MsgPrecommit, decodedMsg.H(), i, p.V(), s.Core)
			mPC, err := s.FinalizeMessage(msgPC)
			if err != nil {
				s.Logger().Warn("Cannot simulate Misbehaviour of PVO12 rule.", err)
			}
			messages = append(messages, mPC)
		}
	}
	// simulate a preVote at round 3, for value v, this preVote for new value break PVO12.
	msgPVO12 := et.NewVoteMsg(consensus.MsgPrevote, decodedMsg.H(), nPR, p.V(), s.Core)
	mPVO12, err := s.FinalizeMessage(msgPVO12)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of PVO12 rule.", err)
	}
	s.Logger().Info("Misbehaviour of PVO12 rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mPVO12)
	for _, m := range messages {
		_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), m)
	}
}

type MisbehaviourRuleInvalidProposalBroadcaster struct {
	*core.Core
}

func (s *MisbehaviourRuleInvalidProposalBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	if msg.Code != consensus.MsgProposal {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	decodedMsg := et.DecodeMsg(msg, s.Core)
	nextPR := et.NextProposeRound(decodedMsg.R(), s.Core)
	// a proposal with invalid header of missing metas.
	header := &types.Header{Number: new(big.Int).SetUint64(decodedMsg.H())}
	block := types.NewBlockWithHeader(header)

	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), nextPR,
		decodedMsg.ValidRound(), block.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of invalid proposal rule", err)
	}
	msgP := et.NewProposeMsg(s.Address(), block, decodedMsg.H(), nextPR, decodedMsg.ValidRound(), liteSig)
	mP, err := s.FinalizeMessage(msgP)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of invalid proposal rule", err)
	}
	s.Logger().Info("Misbehaviour of invalid proposal rule is simulated.")
	//et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
}

type MisbehaviourRuleInvalidProposerBroadcaster struct {
	*core.Core
}

func (s *MisbehaviourRuleInvalidProposerBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	decodedMsg := et.DecodeMsg(msg, s.Core)
	// if current node is the proposer of current round, skip and return.
	if s.CommitteeSet().GetProposer(decodedMsg.R()).Address == s.Address() {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	// current node is not the proposer of current round, propose a proposal.
	header := &types.Header{Number: new(big.Int).SetUint64(decodedMsg.H())}
	block := types.NewBlockWithHeader(header)
	liteSig, err := core.LiteProposalSignature(s.Backend(), new(big.Int).SetUint64(decodedMsg.H()), decodedMsg.R(), -1,
		block.Hash())
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of invalid proposer rule", err)
	}
	msgP := et.NewProposeMsg(s.Address(), block, decodedMsg.H(), decodedMsg.R(), -1, liteSig)
	mP, err := s.FinalizeMessage(msgP)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of invalid proposer rule", err)
	}
	s.Logger().Info("Misbehaviour of invalid proposer rule is simulated.")
	et.DefaultBehaviour(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
}

type MisbehaviourRuleEquivocationBroadcaster struct {
	*core.Core
}

func (s *MisbehaviourRuleEquivocationBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	decodedMsg := et.DecodeMsg(msg, s.Core)
	et.DefaultBehaviour(ctx, s.Core, msg)
	// let proposer of the round send equivocated preVote.
	if decodedMsg.Code == consensus.MsgPrevote && s.IsProposer() {
		msgEq := et.NewVoteMsg(consensus.MsgPrevote, decodedMsg.H(), decodedMsg.R(), et.NonNilValue, s.Core)
		mE, err := s.FinalizeMessage(msgEq)
		if err != nil {
			s.Logger().Warn("Cannot simulate Misbehaviour of equivocation rule", err)
		}
		s.Logger().Info("Misbehaviour of equivocation rule is simulated.")
		_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mE)
	}
}

func TestTBFTMisbehaviourTests(t *testing.T) {

	t.Run("TestTBFTMisbehaviourRulePN", func(t *testing.T) {
		t.Skip("This case is not stable at CI due to the fault simulation is not always deterministic")
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRulePNBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.PN
		runAccountabilityEventTest(t, handler, tp, rule, 45)
	})
	t.Run("TestTBFTMisbehaviourRulePO", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRulePOBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.PO
		runAccountabilityEventTest(t, handler, tp, rule, 100)
	})
	t.Run("TestTBFTMisbehaviourRulePVN", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRulePVNBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.PVN
		runAccountabilityEventTest(t, handler, tp, rule, 60)
	})
	t.Run("TestTBFTMisbehaviourRulePVO12", func(t *testing.T) {
		// todo: improve the simulation since the simulation is not deterministic.
		t.Skip("Skip this unstable case")
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRulePVO12Broadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.PVO12
		runAccountabilityEventTest(t, handler, tp, rule, 60)
	})
	t.Run("TestTBFTMisbehaviourRuleInvalidProposal", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRuleInvalidProposalBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.InvalidProposal
		runAccountabilityEventTest(t, handler, tp, rule, 60)
	})
	t.Run("TestTBFTMisbehaviourRuleInvalidProposer", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRuleInvalidProposerBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.InvalidProposer
		runAccountabilityEventTest(t, handler, tp, rule, 45)
	})
	t.Run("TestTBFTMisbehaviourRuleEquivocation", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &MisbehaviourRuleEquivocationBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.Equivocation
		runAccountabilityEventTest(t, handler, tp, rule, 50)
	})
}
