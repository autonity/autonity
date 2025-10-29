package tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"

	"github.com/stretchr/testify/require"
)

var (
	reporter = *params.TestAutonityContractConfig.Validators[0].NodeAddress
)

func TestAccusation(t *testing.T) {
	r := Setup(t, nil)

	// load the accountability precompiles into the EVM
	accountability.LoadPrecompiles()

	// setup current height
	currentHeight := uint64(1024)
	// finalize blocks from 1st block to current height block to construct internal state includes lastFinalizedBlock, etc...
	r.WaitNBlocks(int(currentHeight - 1))
	lastCommittedHeight := currentHeight - 1

	config, _, err := r.Accountability.GetConfig(nil)
	require.NoError(t, err)

	r.Run("PVN accusation with prevote nil should revert", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Delta.Uint64()
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{0x1} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("PVO accusation with prevote nil should revert", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Delta.Uint64()
		r.Evm.Context.GetHash = func(_ uint64) common.Hash { return common.Hash{0x1} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{}, reporter, 0, autonity.PVO))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("C1 accusation with prevote nil should revert", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Delta.Uint64()
		r.Evm.Context.GetHash = func(_ uint64) common.Hash { return common.Hash{0x1} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{}, reporter, 0, autonity.C1))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("accusation for committed value should revert", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Delta.Uint64()
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{0xca, 0xfe} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})

	r.Run("reporting right tests", func(r *Runner) {
		// reporting should be reverted since reporter is not in current committee and last committee
		accusationHeight := lastCommittedHeight - config.Range.Uint64() + (config.Range.Uint64() / 4) + 1
		noAccessor := common.Address{}
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, noAccessor, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		require.Equal(r.T, "execution reverted: function restricted to a committee member", err.Error())
		// set committee with reporter
		committee, _, err := r.Autonity.GetCommittee(nil)
		require.NoError(r.T, err)
		var newCommittee []common.Address
		for _, c := range committee {
			newCommittee = append(newCommittee, c.Addr)
		}
		newCommittee = append(newCommittee, noAccessor)
		// set the new committee that contains reporter account, then it is allowed for reporting.
		_, err = r.Accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee)
		require.NoError(r.T, err)
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{0xca, 0xfe} }
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}, noAccessor, 0, autonity.PVN))
		require.NoError(r.T, err)
		// now set new committee, it will set current committee as last committee, the reporter is still allowed for reporting.
		_, err = r.Accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee[0:len(newCommittee)-1])
		require.NoError(r.T, err)
		// report same accusation should be reverted since the accusation is pending now.
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}, noAccessor, 0, autonity.PVN))
		require.Equal(t, "execution reverted: already processing an accusation", err.Error())
		// now set new committee without having the reporter, then it is not allowed for reporting.
		_, err = r.Accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee[0:len(newCommittee)-1])
		require.NoError(r.T, err)
		// report same accusation should be reverted since the accusation is pending now.
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}, noAccessor, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		require.Equal(r.T, "execution reverted: function restricted to a committee member", err.Error())
	})
}

func TestAccusationTiming(t *testing.T) {
	r := Setup(t, nil)

	accountability.LoadPrecompiles()

	currentHeight := uint64(1024) // height of current consensus instance
	// finalize blocks from 1st block to current height block to construct internal state includes lastFinalizedBlock, etc...
	r.WaitNBlocks(int(currentHeight - 1))
	r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
	lastCommittedHeight := currentHeight - 1 // height of last committed block

	config, _, err := r.Accountability.GetConfig(nil)
	require.NoError(t, err)

	r.Run("submit accusation at height = lastCommittedHeight - delta (valid)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Delta.Uint64()
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.NoError(r.T, err)
	})
	r.Run("submit accusation at height = lastCommittedHeight - delta + 1 (too recent)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Delta.Uint64() + 1

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight (too recent)", func(r *Runner) {
		accusationHeight := lastCommittedHeight

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight + 5 (future)", func(r *Runner) {
		accusationHeight := lastCommittedHeight + 5

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange (too old)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Range.Uint64()

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4)  (too old)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Range.Uint64() + (config.Range.Uint64() / 4)

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})

	r.Run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4) + 1  (valid)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - config.Range.Uint64() + (config.Range.Uint64() / 4) + 1
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
		require.NoError(r.T, err)
	})
}

// a validator pausing at epoch x can still be accused in the first blocks epoch x+1
// but accusation needs to be for blocks of epoch x
func TestCrossEpochAccusation(t *testing.T) {
	r := Setup(t, nil)
	accountability.LoadPrecompiles()

	// pause the offender
	offender := *params.TestAutonityContractConfig.Validators[0].NodeAddress
	offenderTreasury := params.TestAutonityContractConfig.Validators[0].Treasury
	_, err := r.Autonity.PauseValidator(FromSender(offenderTreasury, nil), offender)
	require.NoError(t, err)

	epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
	require.NoError(t, err)

	config, _, err := r.Accountability.GetConfig(nil)
	require.NoError(t, err)

	r.WaitNBlocks(int(epochPeriod.Uint64() + config.Delta.Uint64() - 2))

	epochID, _, err := r.Autonity.GetEpochID(nil)
	require.NoError(t, err)

	require.Equal(t, uint64(1), epochID.Uint64())

	// offender should not be in committee anymore
	committee, _, err := r.Autonity.GetCommittee(nil)
	require.NoError(t, err)
	for _, member := range committee {
		if member.Addr == offender {
			t.Fatalf("offender is still in committee for epoch 1")
		}
	}

	// accusation should be for a block of past epoch
	accusationHeight := r.Evm.Context.BlockNumber.Uint64() - config.Delta.Uint64() - 1
	epochID, _, err = r.Autonity.GetEpochFromBlock(nil, new(big.Int).SetUint64(accusationHeight))
	require.NoError(t, err)

	require.Equal(t, uint64(0), epochID.Uint64())

	r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
	_, err = r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
	require.NoError(r.T, err)
}

// increasing the height range through the protocol contracts
// potentially opens up a window where undefendable accusation could be raised
// (undefendable because the defender already discarded the old messages).
// these accusation still need to be rejected
func TestUndefendableAccusation(t *testing.T) {
	r := Setup(t, nil)

	accountability.LoadPrecompiles()

	currentHeight := uint64(999) // height of current consensus instance
	r.WaitNBlocks(int(currentHeight - 1))
	require.Equal(t, r.Evm.Context.BlockNumber.Uint64(), currentHeight)
	t.Logf("current core height: %d", currentHeight)

	config, _, err := r.Accountability.GetConfig(nil)
	require.NoError(t, err)

	// accusationHeight < currentHeight - heightRange + heightRange / 4 --> too old
	accusationHeight := currentHeight - config.Range.Uint64() + (config.Range.Uint64() / 4) - 1
	_, err = r.Accountability.CallHandleAccusation(r, &runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
	t.Log(err)
	require.ErrorIs(r.T, err, vm.ErrExecutionReverted)

	// accusationHeight >= currentHeight - heightRange + heightRange / 4 --> valid
	accusationHeight = currentHeight - config.Range.Uint64() + (config.Range.Uint64() / 4)
	_, err = r.Accountability.CallHandleAccusation(r, &runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
	require.NoError(r.T, err)

	// now let's increase the range of 12 blocks
	newRange := new(big.Int).SetUint64(config.Range.Uint64() + 12)
	_, err = r.Accountability.SetRange(r.Operator, newRange)
	require.NoError(r.T, err)
	t.Logf("new range applied: %d", newRange.Uint64())

	// finalize 1 more block to make the change take effect
	r.WaitNBlocks(1)
	currentHeight += 1
	require.Equal(t, r.Evm.Context.BlockNumber.Uint64(), currentHeight)
	t.Logf("current core height: %d", currentHeight)

	newConfig, _, err := r.Accountability.GetConfig(nil)
	require.NoError(t, err)
	require.Equal(t, newRange.Uint64(), newConfig.Range.Uint64())

	// we need an initial grace period to allow for nodes to adapt their buffers to the new range
	deltaRange := newConfig.Range.Uint64() - config.Range.Uint64()
	gracePeriod := deltaRange - deltaRange/4
	t.Logf("delta range %d, gracePeriod %d", deltaRange, gracePeriod)

	// Normally the following accusation should be valid,
	// however since we are in a range transition period they
	// are considered invalid, since validators might have already
	// garbage collected those messages.
	previousOldestValidHeight := currentHeight - config.Range.Uint64() + (config.Range.Uint64() / 4)
	for i := uint64(0); i < (gracePeriod * 2); i++ {
		oldestValidHeight := currentHeight - newConfig.Range.Uint64() + (newConfig.Range.Uint64() / 4)
		t.Logf("currentHeight: %d, oldestValidHeight: %d, previousOldestValidHeight: %d", currentHeight, oldestValidHeight, previousOldestValidHeight)
		for j := int(-gracePeriod); j < int(gracePeriod*2); j++ {
			accusationHeight := uint64(int(oldestValidHeight) + j)
			if accusationHeight < oldestValidHeight {
				t.Logf("currentHeight: %d, accusation for %d should fail (too old)", currentHeight, accusationHeight)
				_, err = r.Accountability.CallHandleAccusation(r, &runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
				t.Log(err)
				require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
				continue
			}
			if accusationHeight < previousOldestValidHeight {
				t.Logf("currentHeight: %d, accusation for %d should fail (grace period)", currentHeight, accusationHeight)
				_, err = r.Accountability.CallHandleAccusation(r, &runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
				t.Log(err)
				require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
				continue
			}
			t.Logf("currentHeight: %d, accusation for %d should be fine", currentHeight, accusationHeight)
			_, err = r.Accountability.CallHandleAccusation(r, &runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN))
			require.NoError(t, err)
		}
		t.Log("mining one block")
		r.WaitNBlocks(1)
		currentHeight += 1
	}

}
