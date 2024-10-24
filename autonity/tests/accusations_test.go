package tests

import (
	"testing"

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

	// TODO(lorenzo) add similar tests for PVO and C1
	r.Run("PVN accusation with prevote nil should revert", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.DeltaBlocks
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{0x1} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("accusation for committed value should revert", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.DeltaBlocks
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{0xca, 0xfe} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})

	r.Run("reporting right tests", func(r *Runner) {

		// reporting should be reverted since reporter is not in current committee and last committee
		accusationHeight := lastCommittedHeight - accountability.HeightRange + (accountability.HeightRange / 4) + 1
		noAccessor := common.Address{}
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, noAccessor))
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
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}, noAccessor))
		require.NoError(r.T, err)
		// now set new committee, it will set current committee as last committee, the reporter is still allowed for reporting.
		_, err = r.Accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee[0:len(newCommittee)-1])
		require.NoError(r.T, err)
		// report same accusation should be reverted since the accusation is pending now.
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}, noAccessor))
		require.Equal(t, "execution reverted: already processing an accusation", err.Error())
		// nwo set new committee without having the reporter, then it is not allowed for reporting.
		_, err = r.Accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee[0:len(newCommittee)-1])
		require.NoError(r.T, err)
		// report same accusation should be reverted since the accusation is pending now.
		_, err = r.Accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}, noAccessor))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		require.Equal(r.T, "execution reverted: function restricted to a committee member", err.Error())
	})
}

func TestAccusationTiming(t *testing.T) {
	r := Setup(t, nil)

	// no more dependency of blockchain now.
	accountability.LoadPrecompiles()

	currentHeight := uint64(1024) // height of current consensus instance
	// finalize blocks from 1st block to current height block to construct internal state includes lastFinalizedBlock, etc...
	r.WaitNBlocks(int(currentHeight - 1))
	r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
	lastCommittedHeight := currentHeight - 1 // height of last committed block

	r.Run("submit accusation at height = lastCommittedHeight - delta (valid)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.DeltaBlocks
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.NoError(r.T, err)
	})
	r.Run("submit accusation at height = lastCommittedHeight - delta + 1 (too recent)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.DeltaBlocks + 1

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight (too recent)", func(r *Runner) {
		accusationHeight := lastCommittedHeight

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight + 5 (future)", func(r *Runner) {
		accusationHeight := lastCommittedHeight + 5

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange (too old)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.HeightRange

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4)  (too old)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.HeightRange + (accountability.HeightRange / 4)

		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})

	r.Run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4) + 1  (valid)", func(r *Runner) {
		accusationHeight := lastCommittedHeight - accountability.HeightRange + (accountability.HeightRange / 4) + 1
		r.Evm.Context.GetHash = func(n uint64) common.Hash { return common.Hash{} }
		_, err := r.Accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
		require.NoError(r.T, err)
	})
}
