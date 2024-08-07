package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/log"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestTendermintStateManagement(t *testing.T) {
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	addr := committeeSet.Committee()[0].Address // round 3 - height 1 proposer
	height := uint64(1)
	round := int64(0)
	signer := makeSigner(keys[addr].consensus)
	signerMember := &committeeSet.Committee()[0]
	cSize := len(committeeSet.Committee())
	proposal := generateBlockProposal(round, new(big.Int).SetUint64(height), -1, false, signer, signerMember)

	logger := log.New()
	state := newTendermintState(logger, nil, nil)

	expectedState := newStateWithoutWAL(logger)
	if !deep.Equal(expectedState, state) {
		t.Fatalf("expected: %v, got: %v", expectedState, state)
	}

	// start new height.
	newHeight := common.Big1
	state.StartNewHeight(newHeight)
	expectedState.height = newHeight
	if !deep.Equal(expectedState, state) {
		t.Fatalf("expected: %v, got: %v", expectedState, state)
	}

	// start new round
	newRound := int64(1)
	state.StartNewRound(newRound)
	expectedState.round = newRound
	if !deep.Equal(expectedState.TendermintState, state.(*TendermintStateImpl).TendermintState) {
		t.Fatalf("expected: %v\n, got: %v", expectedState, state)
	}

	roundMsgs := state.GetOrCreate(newRound)
	require.NotNil(t, roundMsgs)
	require.Equal(t, 0, len(roundMsgs.AllMessages()))

	// set new step.
	state.SetStep(Prevote)
	require.Equal(t, Prevote, state.Step())

	// set decision
	v, r := generateBlock(common.Big1), int64(1)
	state.SetDecision(v, r)
	require.Equal(t, v, state.Decision())
	require.Equal(t, r, state.DecisionRound())

	// setLockedRoundAndValue
	lockedV, lockedR := v, int64(0)
	state.SetLockedRoundAndValue(lockedR, lockedV)
	require.Equal(t, lockedV, state.LockedValue())
	require.Equal(t, lockedR, state.LockedRound())

	// setValidRoundAndValue
	validV, validR := v, int64(0)
	state.SetValidRoundAndValue(validR, validV)
	require.Equal(t, validV, state.ValidValue())
	require.Equal(t, validR, state.ValidRound())

	// setSentProposal
	state.SetSentProposal()
	require.Equal(t, true, state.SentProposal())

	// sentPrevote
	state.SetSentPrevote()
	require.Equal(t, true, state.SentPrevote())

	// sentPrecommit
	state.SetSentPrecommit()
	require.Equal(t, true, state.SentPrecommit())

	// add proposal
	state.SetProposal(proposal, true)
	rMsgs := state.GetOrCreate(proposal.R())
	actualProposal := rMsgs.Proposal()
	if !deep.Equal(proposal, actualProposal) {
		t.Fatalf("expected: %v\n, got: %v", proposal, actualProposal)
	}

	// add prevote
	preVote := message.NewPrevote(round, height, common.Hash{}, signer, &committeeSet.Committee()[0], cSize)
	state.AddPrevote(preVote)
	preVotes := state.GetOrCreate(round).AllPrevotes()
	require.Equal(t, 1, len(preVotes))
	if !deep.Equal(preVote, preVotes[0]) {
		t.Fatalf("expected: %v\n, got: %v", preVote, preVotes[0])
	}

	// add precommit
	preCommit := message.NewPrecommit(round, height, common.Hash{}, signer, &committeeSet.Committee()[0], cSize)
	state.AddPrecommit(preCommit)
	preCommits := state.GetOrCreate(round).AllPrecommits()
	require.Equal(t, 1, len(preCommits))
	if !deep.Equal(preCommit, preCommits[0]) {
		t.Fatalf("expected: %v\n, got: %v", preCommit, preCommits[0])
	}
}
