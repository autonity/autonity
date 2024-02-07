package byzantine

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
)

func runTest(t *testing.T, services *interfaces.Services, eventType autonity.AccountabilityEventType, rule autonity.Rule, period uint64) {

	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// set Malicious validators
	faultyNode := 0
	validators[faultyNode].TendermintServices = services
	// creates a network of 4 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(period, 500, false)
	require.NoError(t, err)

	// check if the misbehaviour is presented for faulty node #0
	faultyAddress := network[faultyNode].Address
	detected := e2e.AccountabilityEventDetected(t, faultyAddress, eventType, rule, network)
	require.Equal(t, true, detected)
}

func newPNBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &PN{c.(*core.Core), false}
}

type PN struct {
	*core.Core
	done bool
}

// simulate a context of msgs that node proposes a new proposal rather than the one it locked at previous rounds.
func (s *PN) Broadcast(msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if s.done || s.Height().Uint64() < 10 || !isProposal {
		s.BroadcastAll(msg)
		return
	}
	nPR := e2e.NextProposeRound(msg.R(), s.Core)
	s.Logger().Info("Simulating PN fault", "h", s.Core.Height(), "r", s.Core.Round(), "npr", nPR)
	// simulate a preCommit msg that locked a value at previous round than next proposing round.
	msgEvidence := message.NewPrecommit(nPR-1, msg.H(), e2e.NonNilValue, s.Backend().Sign)
	printMessage(msgEvidence)
	// simulate a proposal that propose a new value with -1 as the valid round.
	//msgPN := message.NewPropose(proposal.ProposalBlock, decodedMsg.H(), nPR, -1, s.Core)
	msgPN := message.NewPropose(nPR, msg.H(), -1, proposal.Block(), s.Backend().Sign)
	printMessage(msgPN)
	s.BroadcastAll(msg)
	s.BroadcastAll(msgEvidence)
	s.BroadcastAll(msgPN)
	s.done = true
}

func newPOBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &PO{c.(*core.Core), false}
}

type PO struct {
	*core.Core
	done bool
}

// simulate a context of msgs that node proposes a value for which was not the one it locked on.
func (s *PO) Broadcast(msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if s.done || s.Height().Uint64() < 10 || !isProposal {
		s.BroadcastAll(msg)
		return
	}
	// start to simulate malicious context to break rule PO.
	nPR := e2e.NextProposeRound(proposal.R(), s.Core)
	vR := nPR - 1
	s.Logger().Info("Simulating PO fault", "h", s.Core.Height(), "r", s.Core.Round(), "npr", nPR)
	// simulate a preCommit proposal that locked a value at vR.
	msgEvidence := message.NewPrecommit(vR, proposal.H(), e2e.NonNilValue, s.Backend().Sign)
	printMessage(msgEvidence)
	// simulate a proposal that node propose for an old value which it is not the one it locked.
	msgPO := message.NewPropose(nPR, proposal.H(), vR, proposal.Block(), s.Core.Backend().Sign)
	printMessage(msgPO)
	s.BroadcastAll(proposal)
	s.BroadcastAll(msgEvidence)
	s.BroadcastAll(msgPO)
	s.done = true
}

/* currently not used, see later commented tests
func newPVNBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &PVN{c.(*core.Core), false}
}
*/

type PVN struct {
	*core.Core
	done bool
}

// simulate a context of msgs that a node preVote for a new value rather than the one it locked on.
// An example context like below:
// preCommit (h, r, v1)
// preVote   (h, r+1, v2)
func (s *PVN) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	proposal, isProposal := msg.(*message.Propose)
	if s.done || s.Height().Uint64() < 10 || !isProposal {
		return
	}
	nPR := e2e.NextProposeRound(proposal.R(), s.Core)
	// Create a block with a new value, the hash should be different
	newHeader := proposal.Block().Header()
	newHeader.Time = 1337
	newBlock := types.NewBlockWithHeader(newHeader)
	newProposal := message.NewPropose(nPR, proposal.H(), -1, newBlock, s.Core.Backend().Sign)
	fmt.Println("BYZ PROPOSAL HASH", "old", proposal.Value(), "new", newProposal.Value())
	s.BroadcastAll(newProposal)
	// simulate a preCommit at round r, for value v1.
	precommit := message.NewPrecommit(proposal.R(), proposal.H(), proposal.Block().Hash(), s.Backend().Sign)
	// simulate nil precommits until nPr to get contiguous precommits
	for i := proposal.R() + 1; i < nPR; i++ {
		nilPrecommit := message.NewPrecommit(i, proposal.H(), core.NilValue, s.Backend().Sign)
		s.BroadcastAll(nilPrecommit)
	}
	// simulate a preVote at round nPR, for value v2, this preVote for new value break PVN.
	evidence := message.NewPrecommit(nPR, proposal.H(), newProposal.Value(), s.Backend().Sign)
	s.BroadcastAll(precommit)
	s.BroadcastAll(evidence)
	s.done = true

}

/* currently not used, see later commented tests
func newPVO1Broadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &PVO1{c.(*core.Core), false}
}
*/

// PVO rule requires coordination from multiple agents otherwise only the proposal for "PO" will be submitted on-chain.
type PVO1 struct {
	*core.Core
	done bool
}

// simulate a context of msgs that a node preVote for a value that is not the one it precommitted at previous round.
func (s *PVO1) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	proposal, isProposal := msg.(*message.Propose)
	if s.done || s.Height().Uint64() < 10 || !isProposal {
		return
	}
	// we pickup a round far ahead so we don't generate equivocations
	round := e2e.NextProposeRound(20, s.Core)
	// set a valid round.
	validRound := round - 5

	newProposal := message.NewPropose(round, msg.H(), validRound, proposal.Block(), s.Backend().Sign)
	s.BroadcastAll(newProposal)

	for r := validRound; r < round; r++ {
		// send precommit for nil and one not for vr
		val := common.Hash{}
		if r == round-1 {
			val = e2e.NonNilValue
		}
		precommit := message.NewPrecommit(r, newProposal.H(), val, s.Backend().Sign)
		s.BroadcastAll(precommit)
	}
	evidence := message.NewPrecommit(round, newProposal.H(), newProposal.Value(), s.Backend().Sign)
	s.BroadcastAll(evidence)
	s.done = true
}

func newInvalidProposalBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &InvalidProposal{c.(*core.Core)}
}

type InvalidProposal struct {
	*core.Core
}

func (s *InvalidProposal) Broadcast(msg message.Msg) {
	proposal, isProposal := msg.(*message.Propose)
	if !isProposal {
		s.BroadcastAll(msg)
		return
	}
	nextPR := e2e.NextProposeRound(msg.R(), s.Core)
	// a proposal with invalid header of missing metas.
	header := &types.Header{Number: new(big.Int).SetUint64(msg.H())}
	block := types.NewBlockWithHeader(header)
	newProposal := message.NewPropose(nextPR, msg.H(), proposal.ValidRound(), block, s.Backend().Sign)

	s.Logger().Info("Misbehaviour of invalid proposal rule is simulated.")
	s.BroadcastAll(msg)
	s.BroadcastAll(newProposal)
}

func newInvalidProposer(c interfaces.Core) interfaces.Broadcaster {
	return &InvalidProposer{c.(*core.Core)}
}

type InvalidProposer struct {
	*core.Core
}

func (s *InvalidProposer) Broadcast(msg message.Msg) {
	// if current node is the proposer of current round, skip and return.
	if s.CommitteeSet().GetProposer(msg.R()).Address == s.Address() {
		s.BroadcastAll(msg)
		return
	}
	// current node is not the proposer of current round, propose a proposal.
	header := &types.Header{Number: new(big.Int).SetUint64(msg.H())}
	block := types.NewBlockWithHeader(header)
	msgP := message.NewPropose(msg.R(), msg.H(), -1, block, s.Backend().Sign)

	s.Logger().Info("Invalid proposer simulation")
	s.BroadcastAll(msg)
	s.BroadcastAll(msgP)
}

func newEquivocation(c interfaces.Core) interfaces.Broadcaster {
	return &Equivocation{c.(*core.Core)}
}

type Equivocation struct {
	*core.Core
}

func (s *Equivocation) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	if _, isPrevote := msg.(*message.Prevote); !isPrevote {
		return
	}
	// let proposer of the round send equivocated preVote.
	if s.IsProposer() {
		msgEq := message.NewPrevote(msg.R(), msg.H(), e2e.NonNilValue, s.Backend().Sign)
		s.Logger().Info("Equivocation simulation")
		s.BroadcastAll(msgEq)
	}
}

func TestFaultProofs(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		broadcasters func(c interfaces.Core) interfaces.Broadcaster
		rule         autonity.Rule
	}{
		{"PN", newPNBroadcaster, autonity.PN}, // Pass with 120
		{"PO", newPOBroadcaster, autonity.PO}, // Pass with 120
		// {"PVN", newPVNBroadcaster, autonity.PVN}, //Not supported, need multiple byzantine validators
		// {"PVO1", newPVO1Broadcaster, autonity.PVO12}, Not supported currently, need multiple byzantine validators to generate.
		// {"InvalidProposal", newInvalidProposalBroadcaster, autonity.InvalidProposal}, Invalid proposals are not currently supported
		{"InvalidProposer", newInvalidProposer, autonity.InvalidProposer}, // Pass with 120
		{"Equivocation", newEquivocation, autonity.Equivocation},          // Pass with 120
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			handler := &interfaces.Services{Broadcaster: test.broadcasters}
			runTest(t, handler, autonity.Misbehaviour, test.rule, 120)
		})
	}

}

func printMessage(message message.Msg) {
	marshalled, _ := json.MarshalIndent(message, "", "\t")
	fmt.Println(string(marshalled))
}
