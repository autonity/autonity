package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/log"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestTendermintStateDB(t *testing.T) {
	logger := log.New()
	memoryDB := rawdb.NewMemoryDatabase()
	db := newTendermintStateDB(logger, memoryDB)
	require.NotNil(t, db)
	require.Equal(t, uint64(0), db.lastConsensusMsgID)
	require.Equal(t, uint64(0), db.maxMsgID)

	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	addr := committeeSet.Committee()[0].Address // round 3 - height 1 proposer
	height := uint64(1)
	round := int64(0)
	signer := makeSigner(keys[addr].consensus)
	signerMember := &committeeSet.Committee()[0]
	cSize := len(committeeSet.Committee())
	proposal := generateBlockProposal(round, new(big.Int).SetUint64(height), -1, false, signer, signerMember)
	header := &types.Header{Number: common.Big0, Committee: committeeSet.Committee()}

	t.Run("flush tendermint state", func(t *testing.T) {
		state := TendermintState{
			height:                common.Big256,
			round:                 1,
			step:                  Propose,
			decision:              proposal.Block(),
			lockedRound:           0,
			validRound:            0,
			lockedValue:           proposal.Block(),
			validValue:            proposal.Block(),
			sentProposal:          true,
			sentPrevote:           true,
			sentPrecommit:         true,
			setValidRoundAndValue: true,
		}
		err := db.UpdateLastRoundState(&state, false)
		require.NoError(t, err)

		flushedState := db.GetLastTendermintState()
		require.Equal(t, state.height.Uint64(), flushedState.Height)
		require.Equal(t, state.round, int64(flushedState.Round))
		require.Equal(t, state.step, flushedState.Step)
		require.Equal(t, state.decision.Hash(), flushedState.Decision)
		require.Equal(t, state.lockedRound, int64(flushedState.LockedRound))
		require.Equal(t, state.validRound, int64(flushedState.ValidRound))
		require.Equal(t, state.lockedValue.Hash(), flushedState.LockedValue)
		require.Equal(t, state.validValue.Hash(), flushedState.ValidValue)
		require.Equal(t, state.sentProposal, flushedState.SentProposal)
		require.Equal(t, state.sentPrevote, flushedState.SentPrevote)
		require.Equal(t, state.sentPrecommit, flushedState.SentPrecommit)
		require.Equal(t, state.setValidRoundAndValue, flushedState.SetValidRoundAndValue)
	})

	// todo: Jason, add gc collection test trigger by height rotation.
	t.Run("flush state with height rotation", func(t *testing.T) {

	})

	t.Run("flush consensus messages", func(t *testing.T) {
		// flush proposal
		err := db.AddMsg(proposal, true)
		require.NoError(t, err)
		require.Equal(t, uint64(1), db.maxMsgID)
		require.Equal(t, uint64(1), db.lastConsensusMsgID)
		flushedMaxMsgID, err := db.GetMsgID(maxMessageIDKey)
		require.NoError(t, err)
		require.Equal(t, uint64(1), flushedMaxMsgID)
		flushedLastMsgID, err := db.GetMsgID(lastTBFTInstanceMsgIDKey)
		require.NoError(t, err)
		require.Equal(t, uint64(1), flushedLastMsgID)

		msg, verified, err := db.GetMsg(db.lastConsensusMsgID)
		require.NoError(t, err)
		require.Equal(t, true, verified)

		err = msg.PreValidate(header)
		require.NoError(t, err)
		err = msg.Validate()
		require.NoError(t, err)
		actualProposal := msg.(*message.Propose)
		require.Equal(t, proposal.H(), actualProposal.H())
		require.Equal(t, proposal.R(), actualProposal.R())
		require.Equal(t, proposal.Code(), actualProposal.Code())
		require.Equal(t, proposal.ValidRound(), actualProposal.ValidRound())
		require.Equal(t, proposal.Signer(), actualProposal.Signer())
		require.Equal(t, proposal.SignerIndex(), actualProposal.SignerIndex())
		require.Equal(t, proposal.Power(), actualProposal.Power())
		require.Equal(t, proposal.SignerKey(), actualProposal.SignerKey())

		require.Equal(t, proposal.Hash(), actualProposal.Hash())
		require.Equal(t, proposal.Payload(), actualProposal.Payload())
		require.Equal(t, proposal.Validate(), actualProposal.Validate())
		require.Equal(t, proposal.Block().Hash(), actualProposal.Block().Hash())
		require.Equal(t, proposal.Block().Number().Uint64(), actualProposal.Block().Number().Uint64())

		// flush a prevote
		preVote := message.NewPrevote(round, height, common.Hash{}, signer, &committeeSet.Committee()[0], cSize)
		err = db.AddMsg(preVote, true)
		require.NoError(t, err)
		require.Equal(t, uint64(2), db.maxMsgID)
		require.Equal(t, uint64(2), db.lastConsensusMsgID)
		flushedMaxMsgID, err = db.GetMsgID(maxMessageIDKey)
		require.NoError(t, err)
		require.Equal(t, uint64(2), flushedMaxMsgID)
		flushedLastMsgID, err = db.GetMsgID(lastTBFTInstanceMsgIDKey)
		require.NoError(t, err)
		require.Equal(t, uint64(2), flushedLastMsgID)

		msg, verified, err = db.GetMsg(db.lastConsensusMsgID)
		require.NoError(t, err)
		require.Equal(t, true, verified)
		err = msg.PreValidate(header)
		require.NoError(t, err)
		err = msg.Validate()
		require.NoError(t, err)
		if !deep.Equal(preVote, msg) {
			t.Fatalf("expected: %v\n, got: %v", preVote, msg)
		}

		// flush a precommit
		precomit := message.NewPrecommit(round, height, common.Hash{}, signer, &committeeSet.Committee()[0], cSize)
		err = db.AddMsg(precomit, true)
		require.NoError(t, err)
		require.Equal(t, uint64(3), db.maxMsgID)
		require.Equal(t, uint64(3), db.lastConsensusMsgID)
		flushedMaxMsgID, err = db.GetMsgID(maxMessageIDKey)
		require.NoError(t, err)
		require.Equal(t, uint64(3), flushedMaxMsgID)
		flushedLastMsgID, err = db.GetMsgID(lastTBFTInstanceMsgIDKey)
		require.NoError(t, err)
		require.Equal(t, uint64(3), flushedLastMsgID)

		msg, verified, err = db.GetMsg(db.lastConsensusMsgID)
		require.NoError(t, err)
		require.Equal(t, true, verified)
		err = msg.PreValidate(header)
		require.NoError(t, err)
		err = msg.Validate()
		require.NoError(t, err)
		if !deep.Equal(precomit, msg) {
			t.Fatalf("expected: %v\n, got: %v", precomit, msg)
		}
	})
}
