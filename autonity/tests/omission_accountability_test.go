package tests

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/params"
)

var omissionEpochPeriod = 100

const SCALE_FACTOR = 10_000

// need a longer epoch for omission accountability tests
var configOverride = func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
	config.EpochPeriod = uint64(omissionEpochPeriod)
	return config
}

// helpers
func omissionFinalize(r *runner, absents []common.Address, proposer common.Address, effort *big.Int, proposerFaulty bool, epochEnd bool) {
	_, err := r.omissionAccountability.Finalize(fromAutonity, absents, proposer, effort, proposerFaulty, epochEnd)
	require.NoError(r.t, err)
	r.t.Logf("Omission accountability, finalized block: %d", r.evm.Context.BlockNumber)
	// advance the block context as if we mined a block
	r.evm.Context.BlockNumber = new(big.Int).Add(r.evm.Context.BlockNumber, common.Big1)
	r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
}

func inactivityCounter(r *runner, validator common.Address) int {
	counter, _, err := r.omissionAccountability.InactivityCounter(nil, validator)
	require.NoError(r.t, err)
	return int(counter.Uint64())
}

func probation(r *runner, validator common.Address) int {
	probation, _, err := r.omissionAccountability.ProbationPeriods(nil, validator)
	require.NoError(r.t, err)
	return int(probation.Uint64())
}

func inactivityScore(r *runner, validator common.Address) int {
	score, _, err := r.omissionAccountability.InactivityScores(nil, validator)
	require.NoError(r.t, err)
	return int(score.Uint64())
}

func proposerEffort(r *runner, validator common.Address) int {
	effort, _, err := r.omissionAccountability.ProposerEffort(nil, validator)
	require.NoError(r.t, err)
	return int(effort.Uint64())
}

func totalProposerEffort(r *runner) int {
	effort, _, err := r.omissionAccountability.TotalEffort(nil)
	require.NoError(r.t, err)
	return int(effort.Uint64())
}

func faultyProposer(r *runner, targetHeight int64) bool {
	faulty, _, err := r.omissionAccountability.FaultyProposers(nil, new(big.Int).SetInt64(targetHeight))
	require.NoError(r.t, err)
	return faulty
}

func isValidatorInactive(r *runner, targetHeight int64, validator common.Address) bool {
	inactive, _, err := r.omissionAccountability.InactiveValidators(nil, new(big.Int).SetInt64(targetHeight), validator)
	require.NoError(r.t, err)
	return inactive
}

func TestAccessControl(t *testing.T) {
	r := setup(t, nil)
	r.waitNBlocks(tendermint.DeltaBlocks)

	_, err := r.omissionAccountability.Finalize(r.operator, []common.Address{}, common.Address{}, common.Big0, true, false)
	require.Error(r.t, err)
	_, err = r.omissionAccountability.Finalize(fromAutonity, []common.Address{}, common.Address{}, common.Big0, true, false)
	require.NoError(r.t, err)

	_, err = r.omissionAccountability.DistributeProposerRewards(r.operator, common.Big256)
	require.Error(r.t, err)
	_, err = r.omissionAccountability.DistributeProposerRewards(fromAutonity, common.Big256)
	require.NoError(r.t, err)

	_, err = r.omissionAccountability.SetCommittee(r.operator, []common.Address{}, []common.Address{})
	require.Error(r.t, err)
	_, err = r.omissionAccountability.SetCommittee(fromAutonity, []common.Address{}, []common.Address{})
	require.NoError(r.t, err)

	_, err = r.omissionAccountability.SetLastEpochBlock(r.operator, common.Big256)
	require.Error(r.t, err)
	_, err = r.omissionAccountability.SetLastEpochBlock(fromAutonity, common.Big256)
	require.NoError(r.t, err)

}

func TestProposerLogic(t *testing.T) {
	t.Run("Faulty proposer inactive score increases and height is marked as invalid", func(t *testing.T) {
		r := setup(t, configOverride)
		r.waitNBlocks(tendermint.DeltaBlocks)

		targetHeight := r.evm.Context.BlockNumber.Int64() - tendermint.DeltaBlocks
		proposer := r.committee.validators[0].NodeAddress
		omissionFinalize(r, []common.Address{}, proposer, common.Big0, true, false)

		require.True(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 1, inactivityCounter(r, proposer))

		omissionFinalize(r, []common.Address{}, proposer, common.Big0, true, false)
		require.True(r.t, faultyProposer(r, targetHeight+1))
		require.Equal(r.t, 2, inactivityCounter(r, proposer))

		omissionFinalize(r, []common.Address{}, proposer, common.Big0, false, false)
		require.False(r.t, faultyProposer(r, targetHeight+2))
		require.Equal(r.t, 2, inactivityCounter(r, proposer))

	})
	t.Run("Proposer effort is correctly computed", func(t *testing.T) {
		r := setup(t, configOverride)
		r.waitNBlocks(tendermint.DeltaBlocks)

		targetHeight := r.evm.Context.BlockNumber.Int64() - tendermint.DeltaBlocks
		proposer := r.committee.validators[0].NodeAddress
		omissionFinalize(r, []common.Address{}, proposer, common.Big1, false, false)
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, 1, proposerEffort(r, proposer))
		require.Equal(r.t, 1, totalProposerEffort(r))

		targetHeight = r.evm.Context.BlockNumber.Int64() - tendermint.DeltaBlocks
		proposer = r.committee.validators[0].NodeAddress
		omissionFinalize(r, []common.Address{}, proposer, common.Big3, false, false)
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, 4, proposerEffort(r, proposer))
		require.Equal(r.t, 4, totalProposerEffort(r))

		targetHeight = r.evm.Context.BlockNumber.Int64() - tendermint.DeltaBlocks
		proposer = r.committee.validators[1].NodeAddress
		omissionFinalize(r, []common.Address{}, proposer, common.Big3, false, false)
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, 3, proposerEffort(r, proposer))
		require.Equal(r.t, 7, totalProposerEffort(r))
	})
}

// checks that the inactivity counters are correctly updated according to the lookback window (_recordAbsentees function)
func TestInactivityCounter(t *testing.T) {
	r := setup(t, configOverride)
	r.waitNBlocks(tendermint.DeltaBlocks)

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	lookback := int(config.LookbackWindow.Uint64())

	proposer := r.committee.validators[0].NodeAddress
	fullyOffline := r.committee.validators[1].NodeAddress
	partiallyOffline := r.committee.validators[2].NodeAddress

	absents := []common.Address{fullyOffline}
	partiallyOfflineCounter := 0
	for i := 0; i < lookback-2; i++ {
		if i == lookback/2 {
			absents = append(absents, partiallyOffline)
		}
		targetHeight := r.evm.Context.BlockNumber.Int64() - tendermint.DeltaBlocks
		omissionFinalize(r, absents, proposer, common.Big1, false, false)
		for _, absent := range absents {
			require.True(r.t, isValidatorInactive(r, targetHeight, absent))
			if absent == partiallyOffline {
				partiallyOfflineCounter++
			}
		}
	}

	// we still need two height to have a full lookback window
	// insert a validator faulty height, it should be ignored
	r.t.Logf("current block number in evm: %d", r.evm.Context.BlockNumber.Uint64())
	omissionFinalize(r, absents, proposer, common.Big1, true, false)
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	// here we should update the inactivity counter to 1, but since there was a faulty proposer we extend the lookback period
	r.t.Logf("current block number in evm: %d", r.evm.Context.BlockNumber.Uint64())
	partiallyOfflineCounter++
	omissionFinalize(r, absents, proposer, common.Big1, false, false)
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	// now we have a full lookback period
	partiallyOfflineCounter++
	omissionFinalize(r, absents, proposer, common.Big1, false, false)
	require.Equal(r.t, 1, inactivityCounter(r, fullyOffline))
	partiallyOfflineCounter++
	omissionFinalize(r, absents, proposer, common.Big1, false, false)
	require.Equal(r.t, 2, inactivityCounter(r, fullyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	partiallyOfflineCounter++
	omissionFinalize(r, absents, proposer, common.Big1, false, false)
	require.Equal(r.t, 3, inactivityCounter(r, fullyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))

	// fill up enough blocks for partiallyOffline as well
	for i := partiallyOfflineCounter; i < lookback-1; i++ {
		omissionFinalize(r, absents, proposer, common.Big1, false, false)
	}

	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	omissionFinalize(r, absents, proposer, common.Big1, false, false)
	require.Equal(r.t, 1, inactivityCounter(r, partiallyOffline))

	fullyOfflineIC := inactivityCounter(r, fullyOffline)
	partiallyOfflineIC := inactivityCounter(r, partiallyOffline)

	// every two block, one has faulty proposer
	n := 20
	for i := 0; i < (n * 2); i++ {
		proposerFaulty := i%2 == 0
		omissionFinalize(r, absents, proposer, common.Big1, proposerFaulty, false)
	}

	// inactivity counter should still have increased by n due to lookback period extension
	require.Equal(r.t, fullyOfflineIC+n, inactivityCounter(r, fullyOffline))
	require.Equal(r.t, partiallyOfflineIC+n, inactivityCounter(r, partiallyOffline))

	// reach block 100 and close the epoch
	for i := r.evm.Context.BlockNumber.Int64(); i < 100; i++ {
		omissionFinalize(r, []common.Address{}, proposer, common.Big1, false, false)
	}
	// close the epoch
	t.Log("Closing epoch")
	omissionFinalize(r, []common.Address{}, proposer, common.Big1, false, true)

	// inactivity counters should be reset
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	r.waitNBlocks(tendermint.DeltaBlocks)
	otherValidator := r.committee.validators[3].NodeAddress

	for i := 0; i < lookback/2; i++ {
		omissionFinalize(r, []common.Address{otherValidator}, proposer, common.Big1, false, false)
	}

	t.Log("online at following block")
	omissionFinalize(r, []common.Address{}, proposer, common.Big1, false, false)

	// one block online is going to "save" the validator for the next lookback window
	for i := 0; i < lookback-1; i++ {
		omissionFinalize(r, []common.Address{otherValidator}, proposer, common.Big1, false, false)
		require.Equal(r.t, 0, inactivityCounter(r, otherValidator))
	}

	// proposer faulty
	omissionFinalize(r, []common.Address{otherValidator}, proposer, common.Big1, true, false)
	require.Equal(r.t, 0, inactivityCounter(r, otherValidator))

	omissionFinalize(r, []common.Address{otherValidator}, proposer, common.Big1, false, false)
	require.Equal(r.t, 1, inactivityCounter(r, otherValidator))

}

// checks that the inactivity are computed correctly
func TestInactivityScore(t *testing.T) {
	r := setup(t, configOverride)
	r.waitNBlocks(tendermint.DeltaBlocks)

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	lookback := int(config.LookbackWindow.Uint64())
	pastPerformanceWeight := float64(config.PastPerformanceWeight.Uint64()) / SCALE_FACTOR

	// simulate epoch
	inactiveBlockStreak := make([]int, len(r.committee.validators))
	inactiveCounters := make([]int, len(r.committee.validators))
	for h := tendermint.DeltaBlocks + 1; h < omissionEpochPeriod+1; h++ {
		var absents []common.Address
		for i := range r.committee.validators {
			if rand.Intn(30) != 0 {
				absents = append(absents, r.committee.validators[i].NodeAddress)
				inactiveBlockStreak[i]++
			} else {
				inactiveBlockStreak[i] = 0
			}
			if inactiveBlockStreak[i] >= lookback {
				inactiveCounters[i]++
			}
		}

		epochEnded := h == omissionEpochPeriod
		omissionFinalize(r, absents, r.committee.validators[0].NodeAddress, common.Big1, false, epochEnded)
	}

	// check score computation
	pastInactivityScore := make([]float64, len(r.committee.validators))
	for i, val := range r.committee.validators {
		score := float64(inactiveCounters[i]) / float64(omissionEpochPeriod-tendermint.DeltaBlocks-lookback+1)
		expectedInactivityScoreFloat := score*(1-pastPerformanceWeight) + 0*pastPerformanceWeight
		pastInactivityScore[i] = expectedInactivityScoreFloat
		expectedInactivityScore := int(math.Floor(expectedInactivityScoreFloat * SCALE_FACTOR))
		r.t.Logf("expectedInactivityScore %v, inactivityScore %v", expectedInactivityScore, inactivityScore(r, val.NodeAddress))
		require.Equal(r.t, expectedInactivityScore, inactivityScore(r, val.NodeAddress))
	}

	// simulate another epoch
	r.waitNBlocks(tendermint.DeltaBlocks)
	inactiveBlockStreak = make([]int, len(r.committee.validators))
	inactiveCounters = make([]int, len(r.committee.validators))
	for h := tendermint.DeltaBlocks + 1; h < omissionEpochPeriod+1; h++ {
		var absents []common.Address
		for i := range r.committee.validators {
			if rand.Intn(30) != 0 {
				absents = append(absents, r.committee.validators[i].NodeAddress)
				inactiveBlockStreak[i]++
			} else {
				inactiveBlockStreak[i] = 0
			}
			if inactiveBlockStreak[i] >= lookback {
				inactiveCounters[i]++
			}
		}

		epochEnded := h == omissionEpochPeriod
		omissionFinalize(r, absents, r.committee.validators[0].NodeAddress, common.Big1, false, epochEnded)
	}

	// check score computation
	for i, val := range r.committee.validators {
		score := float64(inactiveCounters[i]) / float64(omissionEpochPeriod-tendermint.DeltaBlocks-lookback+1)
		expectedInactivityScoreFloat := score*(1-pastPerformanceWeight) + pastInactivityScore[i]*pastPerformanceWeight
		expectedInactivityScore := int(math.Floor(expectedInactivityScoreFloat * SCALE_FACTOR))
		r.t.Logf("expectedInactivityScore %v, inactivityScore %v", expectedInactivityScore, inactivityScore(r, val.NodeAddress))
		//TODO(lorenzo) most of the time it passes, but some times it fail by 1 of difference. Probably there is a precision issue.
		require.Equal(r.t, expectedInactivityScore, inactivityScore(r, val.NodeAddress))
	}
}

func TestOmissionPunishments(t *testing.T) {
	r := setup(t, configOverride)
	r.waitNBlocks(tendermint.DeltaBlocks)
	//TODO: implement
}

func TestProposerRewardDistribution(t *testing.T) {
	r := setup(t, configOverride)
	r.waitNBlocks(tendermint.DeltaBlocks)
	//TODO: implement
}
