package tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

func TestAccountabilityEvents(t *testing.T) {
	r := Setup(t, nil)

	// load the accountability precompiles into the EVM
	accountability.LoadPrecompiles()

	offender := *params.TestAutonityContractConfig.Validators[1].NodeAddress

	r.Run("event with nil rawproof should not cause panics or issues", func(r *Runner) {
		misbehaviourEvent := IAccountabilityEvent{
			EventType:      uint8(autonity.Misbehaviour),
			Rule:           uint8(autonity.PVN),
			Reporter:       reporter,
			Offender:       offender,
			RawProof:       nil,
			Id:             common.Big0, // assigned contract-side
			Block:          common.Big0, // assigned contract-side
			ReportingBlock: common.Big0, // assigned contract-side
			Epoch:          common.Big0, // assigned contract-side
			MessageHash:    common.Big0, // assigned contract-side
		}
		_, err := r.Accountability.HandleMisbehaviour(&RunOptions{origin: reporter}, misbehaviourEvent)
		require.Error(r.T, err)
	})
	r.Run("accusation with nil rawproof should not cause panics or issues", func(r *Runner) {
		accusationEvent := NewAccusationEvent(10, common.Hash{0xca, 0xfe}, reporter, 0, autonity.PVN)
		accusationEvent.RawProof = nil
		_, err := r.Accountability.HandleAccusation(&RunOptions{origin: reporter}, accusationEvent)
		require.Error(r.T, err)
	})
}

func TestGracePeriodComputation(t *testing.T) {
	r := Setup(t, nil)

	// load the accountability precompiles into the EVM
	accountability.LoadPrecompiles()

	t.Run("test range changes", func(t *testing.T) {
		rangeChange := new(big.Int).SetUint64(20)

		// print out initial values
		config, _, err := r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		initialRange := config.Range
		require.True(t, initialRange.Cmp(rangeChange) > 0)

		initialGracePeriod, _, err := r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)

		t.Logf("range: %d, gracePeriod: %d", initialRange.Uint64(), initialGracePeriod.Uint64())

		// initial grace period should be 0
		require.Equal(t, common.Big0.String(), initialGracePeriod.String())

		// decrease in range should not affect grace period
		lowerRange := new(big.Int).Sub(initialRange, rangeChange)
		t.Logf("setting range to %d", lowerRange.Uint64())
		_, err = r.Accountability.SetRange(r.Operator, lowerRange)
		require.NoError(t, err)

		// nothing should change before finalize
		config, _, err = r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		require.Equal(t, initialRange.Uint64(), config.Range.Uint64())
		gracePeriod, _, err := r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		require.Equal(t, initialGracePeriod.Uint64(), gracePeriod.Uint64())

		// at finalize, range change applies
		_, err = r.Accountability.Finalize(FromAutonity, false)
		require.NoError(t, err)

		// range should change but grace period should still be 0
		config, _, err = r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		require.Equal(t, lowerRange.Uint64(), config.Range.Uint64())
		gracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		require.Equal(t, common.Big0.Uint64(), gracePeriod.Uint64())

		// increase in range should case gracePeriod to increase
		rangeChange = new(big.Int).SetUint64(rangeChange.Uint64() * 2)
		biggerRange := new(big.Int).Add(lowerRange, rangeChange)
		t.Logf("setting range to %d", biggerRange.Uint64())
		_, err = r.Accountability.SetRange(r.Operator, biggerRange)
		require.NoError(t, err)

		// nothing should change before finalize
		config, _, err = r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		require.Equal(t, lowerRange.Uint64(), config.Range.Uint64())
		gracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		require.Equal(t, common.Big0.Uint64(), gracePeriod.Uint64())

		// at finalize, range change applies
		_, err = r.Accountability.Finalize(FromAutonity, false)
		require.NoError(t, err)
		config, _, err = r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		require.Equal(t, biggerRange.Uint64(), config.Range.Uint64())

		// grace period should be the expected one
		gracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		diff := biggerRange.Uint64() - lowerRange.Uint64()
		expectedGracePeriod := diff - diff/4

		t.Logf("range was: %d, changed to %d, expected grace period %d", lowerRange.Uint64(), biggerRange.Uint64(), expectedGracePeriod)
		require.Equal(t, expectedGracePeriod, gracePeriod.Uint64())

		// grace period should decrease at each block
		r.FinalizeBlock()
		newGracePeriod, _, err := r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		t.Logf("mined block, new grace period %d", newGracePeriod.Uint64())
		require.Equal(t, gracePeriod.Uint64()-1, newGracePeriod.Uint64())

		for i := 0; i < 10; i++ {
			r.FinalizeBlock()
		}

		newGracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		t.Logf("mined 10 blocks, new grace period %d", newGracePeriod.Uint64())
		require.Equal(t, gracePeriod.Uint64()-11, newGracePeriod.Uint64())

		for i := 0; i < 20; i++ {
			r.FinalizeBlock()
		}

		newGracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		t.Logf("mined 20 blocks, new grace period %d", newGracePeriod.Uint64())
		require.Equal(t, common.Big0.Uint64(), newGracePeriod.Uint64())

		for i := 0; i < 100; i++ {
			r.FinalizeBlock()
		}

		newGracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		t.Logf("mined 100 blocks, new grace period %d", newGracePeriod.Uint64())
		require.Equal(t, common.Big0.Uint64(), newGracePeriod.Uint64())
	})
	t.Run("range changes edge cases", func(t *testing.T) {
		rangeChange := new(big.Int).SetUint64(4)

		// print out initial values
		config, _, err := r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		initialRange := config.Range
		require.True(t, initialRange.Cmp(rangeChange) > 0)

		initialGracePeriod, _, err := r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)

		t.Logf("range: %d, gracePeriod: %d", initialRange.Uint64(), initialGracePeriod.Uint64())

		// initial grace period should be 0
		require.Equal(t, common.Big0.String(), initialGracePeriod.String())

		// increase in range should case gracePeriod to increase
		rangeChange = new(big.Int).SetUint64(rangeChange.Uint64())
		biggerRange := new(big.Int).Add(initialRange, rangeChange)
		t.Logf("setting range to %d", biggerRange.Uint64())
		_, err = r.Accountability.SetRange(r.Operator, biggerRange)
		require.NoError(t, err)

		// nothing should change before finalize
		config, _, err = r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		require.Equal(t, initialRange.Uint64(), config.Range.Uint64())
		gracePeriod, _, err := r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		require.Equal(t, common.Big0.Uint64(), gracePeriod.Uint64())

		// at finalize, range change applies
		_, err = r.Accountability.Finalize(FromAutonity, false)
		require.NoError(t, err)
		config, _, err = r.Accountability.GetConfig(nil)
		require.NoError(t, err)
		require.Equal(t, biggerRange.Uint64(), config.Range.Uint64())

		// grace period should be the expected one
		gracePeriod, _, err = r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		diff := biggerRange.Uint64() - initialRange.Uint64()
		expectedGracePeriod := diff - diff/4

		t.Logf("range was: %d, changed to %d, expected grace period %d", initialRange.Uint64(), biggerRange.Uint64(), expectedGracePeriod)
		require.Equal(t, expectedGracePeriod, gracePeriod.Uint64())

		// grace period should decrease at each block
		r.FinalizeBlock()
		newGracePeriod, _, err := r.Accountability.GetGracePeriod(nil)
		require.NoError(t, err)
		t.Logf("mined block, new grace period %d", newGracePeriod.Uint64())
		require.Equal(t, gracePeriod.Uint64()-1, newGracePeriod.Uint64())
	})

}
