package byzantine

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/stretchr/testify/require"
)

func runTest(t *testing.T, services *node.TendermintServices, eventType autonity.AccountabilityEventType, rule autonity.Rule, period uint64) {

	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set Malicious validators
	faultyNode := 0
	validators[faultyNode].TendermintServices = services
	// creates a network of 4 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(period, 500) // nolint

	// check if the misbehaviour is presented for faulty node #0
	faultyAddress := network[faultyNode].Address
	detected := e2e.AccountabilityEventDetected(t, faultyAddress, eventType, rule, network)
	require.Equal(t, true, detected)
}

type PN struct {
	*core.Core
	done bool
}

// simulate a context of msgs that node proposes a new proposal rather than the one it locked at previous rounds.
func (s *PN) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	if s.done || s.Height().Uint64() < 10 || msg.Code != consensus.MsgProposal {
		e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
		return
	}
	_ = msg.DecodePayload()
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	s.Logger().Info("Simulating PN fault", "h", s.Core.Height(), "r", s.Core.Round(), "npr", nPR)

	// simulate a preCommit msg that locked a value at previous round than next proposing round.
	msgEvidence := e2e.NewVoteMsg(consensus.MsgPrecommit, msg.H(), nPR-1, e2e.NonNilValue, s.Core)
	printMessage(msgEvidence)
	// simulate a proposal that propose a new value with -1 as the valid round.
	//msgPN := e2e.NewProposeMsg(proposal.ProposalBlock, decodedMsg.H(), nPR, -1, s.Core)
	msgPN := e2e.NewProposeMsg(s.Core.Address(), msg.ConsensusMsg.(*message.Proposal).ProposalBlock, msg.H(), nPR, -1, s.Core.Backend().Sign)
	printMessage(msgPN)

	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msgEvidence)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msgPN)

	s.done = true
}

type PO struct {
	*core.Core
	done bool
}

// simulate a context of msgs that node proposes a value for which was not the one it locked on.
func (s *PO) SignAndBroadcast(ctx context.Context, proposal *message.Message) {
	if s.done || s.Height().Uint64() < 10 || proposal.Code != consensus.MsgProposal {
		e2e.DefaultSignAndBroadcast(ctx, s.Core, proposal)
		return
	}
	// start to simulate malicious context to break rule PO.
	_ = proposal.DecodePayload()
	nPR := e2e.NextProposeRound(proposal.R(), s.Core)
	vR := nPR - 1
	s.Logger().Info("Simulating PO fault", "h", s.Core.Height(), "r", s.Core.Round(), "npr", nPR)
	// simulate a preCommit proposal that locked a value at vR.
	msgEvidence := e2e.NewVoteMsg(consensus.MsgPrecommit, proposal.H(), vR, e2e.NonNilValue, s.Core)
	printMessage(msgEvidence)
	// simulate a proposal that node propose for an old value which it is not the one it locked.
	msgPO := e2e.NewProposeMsg(s.Address(), proposal.ConsensusMsg.(*message.Proposal).ProposalBlock, proposal.H(), nPR, vR, s.Core.Backend().Sign)
	printMessage(msgPO)

	e2e.DefaultSignAndBroadcast(ctx, s.Core, proposal)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msgEvidence)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msgPO)
	s.done = true
}

type PVN struct {
	*core.Core
	done bool
}

// simulate a context of msgs that a node preVote for a new value rather than the one it locked on.
// An example context like below:
// preCommit (h, r, v1)
// preVote   (h, r+1, v2)
func (s *PVN) SignAndBroadcast(ctx context.Context, proposal *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, proposal)
	if s.done || s.Height().Uint64() < 10 || proposal.Code != consensus.MsgProposal {
		return
	}

	_ = proposal.DecodePayload()
	nPR := e2e.NextProposeRound(proposal.R(), s.Core)
	// Create a block with a new value, the hash should be different
	newHeader := proposal.ConsensusMsg.(*message.Proposal).ProposalBlock.Header()
	newHeader.Time = 1337
	newBlock := types.NewBlockWithHeader(newHeader)

	newProposal := e2e.NewProposeMsg(s.Address(), newBlock, proposal.H(), nPR, -1, s.Core.Backend().Sign)
	fmt.Println("BYZ PROPOSAL HASH", "old", proposal.Value(), "new", newProposal.Value())
	e2e.DefaultSignAndBroadcast(ctx, s.Core, newProposal)
	// simulate a preCommit at round r, for value v1.
	precommit := e2e.NewVoteMsg(consensus.MsgPrecommit, proposal.H(), proposal.R(), proposal.ConsensusMsg.V(), s.Core)
	// simulate nil precommits until nPr to get contiguous precommits
	for i := proposal.R() + 1; i < nPR; i++ {
		nilPrecommit := e2e.NewVoteMsg(consensus.MsgPrecommit, proposal.H(), i, core.NilValue, s.Core)
		e2e.DefaultSignAndBroadcast(ctx, s.Core, nilPrecommit)
	}
	// simulate a preVote at round nPR, for value v2, this preVote for new value break PVN.
	evidence := e2e.NewVoteMsg(consensus.MsgPrevote, proposal.H(), nPR, newProposal.Value(), s.Core)

	e2e.DefaultSignAndBroadcast(ctx, s.Core, precommit)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, evidence)
	s.done = true
	// TODO:(youssef) We need to test the accusation flow when we have an evidence for it too !
}

// PVO rule requires coordination from multiple agents otherwise only the proposal for "PO" will be submitted on-chain.
type PVO1 struct {
	*core.Core
	done bool
}

// simulate a context of msgs that a node preVote for a value that is not the one it precommitted at previous round.
func (s *PVO1) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)

	if s.done || s.Height().Uint64() < 10 || msg.Code != consensus.MsgProposal {
		return
	}

	_ = msg.DecodePayload()

	// we pickup a round far ahead so we don't generate equivocations
	round := e2e.NextProposeRound(20, s.Core)
	// set a valid round.

	validRound := round - 5

	proposal := e2e.NewProposeMsg(s.Address(), msg.ConsensusMsg.(*message.Proposal).ProposalBlock, msg.H(), round, validRound, s.Backend().Sign)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, proposal)

	for r := validRound; r < round; r++ {
		// send precommit for nil and one not for vr
		val := common.Hash{}
		if r == round-1 {
			val = e2e.NonNilValue
		}
		precommit := e2e.NewVoteMsg(consensus.MsgPrecommit, proposal.H(), r, val, s.Core)
		e2e.DefaultSignAndBroadcast(ctx, s.Core, precommit)
	}
	evidence := e2e.NewVoteMsg(consensus.MsgPrevote, proposal.H(), round, proposal.Value(), s.Core)
	e2e.DefaultSignAndBroadcast(ctx, s.Core, evidence)

	s.done = true
}

type InvalidProposal struct {
	*core.Core
}

func (s *InvalidProposal) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	if msg.Code != consensus.MsgProposal {
		e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
		return
	}
	_ = msg.DecodePayload()
	nextPR := e2e.NextProposeRound(msg.R(), s.Core)
	// a proposal with invalid header of missing metas.
	header := &types.Header{Number: new(big.Int).SetUint64(msg.H())}
	block := types.NewBlockWithHeader(header)
	msgP := e2e.NewProposeMsg(s.Address(), block, msg.H(), nextPR, msg.ConsensusMsg.(*message.Proposal).ValidRound, s.Backend().Sign)
	mP, err := s.SignMessage(msgP)
	if err != nil {
		s.Logger().Warn("Cannot simulate Misbehaviour of invalid proposal rule", err)
	}
	s.Logger().Info("Misbehaviour of invalid proposal rule is simulated.")
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
}

type InvalidProposer struct {
	*core.Core
}

func (s *InvalidProposer) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	_ = msg.DecodePayload()
	// if current node is the proposer of current round, skip and return.
	if s.CommitteeSet().GetProposer(msg.R()).Address == s.Address() {
		e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
		return
	}
	// current node is not the proposer of current round, propose a proposal.
	header := &types.Header{Number: new(big.Int).SetUint64(msg.H())}
	block := types.NewBlockWithHeader(header)
	msgP := e2e.NewProposeMsg(s.Address(), block, msg.H(), msg.R(), -1, s.Backend().Sign)
	mP, err := s.SignMessage(msgP)
	if err != nil {
		s.Logger().Crit("Cannot simulate Misbehaviour of invalid proposer rule", err)
	}
	s.Logger().Info("Misbehaviour of invalid proposer rule is simulated.")
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mP)
}

type Equivocation struct {
	*core.Core
}

func (s *Equivocation) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	_ = msg.DecodePayload()
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	// let proposer of the round send equivocated preVote.
	if msg.Code == consensus.MsgPrevote && s.IsProposer() {
		msgEq := e2e.NewVoteMsg(consensus.MsgPrevote, msg.H(), msg.R(), e2e.NonNilValue, s.Core)
		mE, err := s.SignMessage(msgEq)
		if err != nil {
			s.Logger().Warn("Cannot simulate Misbehaviour of equivocation rule", err)
		}
		s.Logger().Info("Misbehaviour of equivocation rule is simulated.")
		_ = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), mE)
	}
}

func TestFaultProofs(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		broadcasters interfaces.Broadcaster
		rule         autonity.Rule
	}{
		{"PN", &PN{}, autonity.PN}, // Pass with 120
		{"PO", &PO{}, autonity.PO}, // Pass with 120
		// {"PVN", &PVN{}, autonity.PVN}, //Not supported, need multiple byzantine validators
		// {"PVO1", &PVO1{}, autonity.PVO12}, Not supported currently, need multiple byzantine validators to generate.
		// {"InvalidProposal", &InvalidProposal{}, autonity.InvalidProposal}, Invalid proposals are not currently supported
		{"InvalidProposer", &InvalidProposer{}, autonity.InvalidProposer}, // Pass with 120
		{"Equivocation", &Equivocation{}, autonity.Equivocation},          // Pass with 120
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			handler := &node.TendermintServices{Broadcaster: test.broadcasters}
			runTest(t, handler, autonity.Misbehaviour, test.rule, 120)
		})
	}

}

func printMessage(message *message.Message) {
	marshalled, _ := json.MarshalIndent(message, "", "\t")
	fmt.Println(string(marshalled))
}
