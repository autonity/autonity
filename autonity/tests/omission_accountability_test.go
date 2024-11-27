package tests

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/params"
)

const BigFloatPrecision = 256 // will allow full representation of the solidity uint256 range

func newFloat(value *big.Int) *big.Float {
	return new(big.Float).SetPrec(BigFloatPrecision).SetInt(value)
}

func toString(value *big.Float) string {
	return value.Text('f', -1)
}

var omissionEpochPeriod = 130

const active = uint8(0)
const jailed = uint8(2)
const jailbound = uint8(3)
const jailedForInactivity = uint8(4)
const jailboundForInactivity = uint8(5)

// need a longer epoch for omission accountability tests
var configOverride = func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
	config.EpochPeriod = uint64(omissionEpochPeriod)
	for _, val := range config.Validators {
		val.BondedStake = new(big.Int).SetUint64(1)
	}
	return config
}

var configOverrideIncreasedStake = func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
	// apply default override and additionally increase voting power of validator 0 to reach quorum in proofs easily
	defaultOverrideGenesis := configOverride(config)
	defaultOverrideGenesis.Validators[0].BondedStake = new(big.Int).SetUint64(10)
	return defaultOverrideGenesis
}

func autonityFinalize(r *Runner) { //nolint
	_, err := r.Autonity.Finalize(nil)
	require.NoError(r.T, err)
	r.T.Logf("Autonity, finalized block: %d", r.Evm.Context.BlockNumber)
	// advance the block context as if we mined a block
	r.Evm.Context.BlockNumber = new(big.Int).Add(r.Evm.Context.BlockNumber, common.Big1)
	r.Evm.Context.Time = new(big.Int).Add(r.Evm.Context.Time, common.Big1)
	// clean up activity proof data
	r.Evm.Context.Coinbase = common.Address{}
	r.Evm.Context.ActivityProof = nil
	r.Evm.Context.ActivityProofRound = 0
}

func setupProofAndAutonityFinalize(r *Runner, proposer common.Address, absentees map[common.Address]struct{}) {
	r.setupActivityProofAndCoinbase(proposer, absentees)
	autonityFinalize(r)
}

func faultyProposers(r *Runner) uint64 {
	n, _, err := r.OmissionAccountability.FaultyProposersInWindow(nil)
	require.NoError(r.T, err)
	return n.Uint64()
}

func inactivityCounter(r *Runner, validator common.Address) int {
	counter, _, err := r.OmissionAccountability.InactivityCounter(nil, validator)
	require.NoError(r.T, err)
	return int(counter.Uint64())
}

func probation(r *Runner, validator common.Address) int {
	probation, _, err := r.OmissionAccountability.ProbationPeriods(nil, validator)
	require.NoError(r.T, err)
	return int(probation.Uint64())
}

func offences(r *Runner, validator common.Address) int {
	offences, _, err := r.OmissionAccountability.RepeatedOffences(nil, validator)
	require.NoError(r.T, err)
	return int(offences.Uint64())
}

func inactivityScore(r *Runner, validator common.Address) int {
	score, _, err := r.OmissionAccountability.InactivityScores(nil, validator)
	require.NoError(r.T, err)
	return int(score.Uint64())
}

func omissionScaleFactor(r *Runner) *big.Int {
	factor, _, err := r.OmissionAccountability.GetScaleFactor(nil)
	require.NoError(r.T, err)
	return factor
}

func slashingRateScaleFactor(r *Runner) *big.Int {
	scaleFactor, _, err := r.slasherContract().GetSlashingScaleFactor(nil)
	require.NoError(r.T, err)
	return scaleFactor
}

func proposerEffort(r *Runner, validator common.Address) *big.Int {
	effort, _, err := r.OmissionAccountability.ProposerEffort(nil, validator)
	require.NoError(r.T, err)
	return effort
}

func totalProposerEffort(r *Runner) *big.Int {
	effort, _, err := r.OmissionAccountability.TotalEffort(nil)
	require.NoError(r.T, err)
	return effort
}

func faultyProposer(r *Runner, targetHeight int64) bool {
	faulty, _, err := r.OmissionAccountability.FaultyProposers(nil, new(big.Int).SetInt64(targetHeight))
	require.NoError(r.T, err)
	return faulty
}

func absenteesLastHeight(r *Runner) []common.Address {
	absentees, _, err := r.OmissionAccountability.GetAbsenteesLastHeight(nil)
	require.NoError(r.T, err)

	return absentees
}

func lastActive(r *Runner, addr common.Address) int64 {
	lastActive, _, err := r.OmissionAccountability.LastActive(nil, addr)
	require.NoError(r.T, err)
	return lastActive.Int64()
}

func isValidatorInactive(r *Runner, targetHeight int64, validator common.Address) bool {
	inactive, _, err := r.OmissionAccountability.InactiveValidators(nil, new(big.Int).SetInt64(targetHeight), validator)
	require.NoError(r.T, err)
	return inactive
}

func validator(r *Runner, addr common.Address) AutonityValidator {
	val, _, err := r.Autonity.GetValidator(nil, addr)
	require.NoError(r.T, err)
	return val
}

func ntnBalance(r *Runner, addr common.Address) *big.Int {
	balance, _, err := r.Autonity.BalanceOf(nil, addr)
	require.NoError(r.T, err)
	return balance
}

func expectNearlyEqual(t *testing.T, expected, actual *big.Int, tolerance int64) {
	diff := new(big.Int).Abs(new(big.Int).Sub(expected, actual))
	require.True(t, diff.Cmp(big.NewInt(tolerance)) <= 0, "expected %v, got %v", expected, actual)
}

func applyCustomParams(r *Runner, customDelta uint64, customLookback uint64) {
	_, err := r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(customDelta))
	require.NoError(r.T, err)
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(customLookback))
	require.NoError(r.T, err)
	// close the epoch to apply new omission params
	r.WaitNextEpoch()
}

func TestAccessControl(t *testing.T) {
	r := Setup(t, nil)

	_, err := r.OmissionAccountability.Finalize(r.Operator, false)
	require.Error(r.T, err)
	_, err = r.OmissionAccountability.Finalize(FromAutonity, false)
	require.NoError(r.T, err)

	_, err = r.OmissionAccountability.DistributeProposerRewards(r.Operator, common.Big256)
	require.Error(r.T, err)
	_, err = r.OmissionAccountability.DistributeProposerRewards(FromAutonity, common.Big256)
	require.NoError(r.T, err)

	_, err = r.OmissionAccountability.SetCommittee(r.Operator, []AutonityCommitteeMember{}, []common.Address{})
	require.Error(r.T, err)
	_, err = r.OmissionAccountability.SetCommittee(FromAutonity, []AutonityCommitteeMember{}, []common.Address{})
	require.NoError(r.T, err)

	_, err = r.OmissionAccountability.SetEpochBlock(r.Operator, common.Big256)
	require.Error(r.T, err)
	_, err = r.OmissionAccountability.SetEpochBlock(FromAutonity, common.Big256)
	require.NoError(r.T, err)
}

func TestProposerLogic(t *testing.T) {
	t.Run("Faulty proposer inactive score increases and height is marked as invalid", func(t *testing.T) {
		r := Setup(t, configOverride)

		delta, _, err := r.OmissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.WaitNBlocks(int(delta.Int64()))

		targetHeight := r.Evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer := r.Committee.Validators[0].NodeAddress
		r.Evm.Context.Coinbase = proposer
		r.Evm.Context.ActivityProof = nil
		autonityFinalize(r)

		require.True(r.T, faultyProposer(r, targetHeight))
		require.Equal(r.T, 1, inactivityCounter(r, proposer))

		r.Evm.Context.Coinbase = proposer
		r.Evm.Context.ActivityProof = nil
		autonityFinalize(r)
		require.True(r.T, faultyProposer(r, targetHeight+1))
		require.Equal(r.T, 2, inactivityCounter(r, proposer))

		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.T, faultyProposer(r, targetHeight+2))
		require.Equal(r.T, 2, inactivityCounter(r, proposer))
	})
	t.Run("Proposer effort is correctly computed", func(t *testing.T) {
		r := Setup(t, configOverride)

		delta, _, err := r.OmissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.WaitNBlocks(int(delta.Int64()))

		totalVotingPower := new(big.Int)
		for _, val := range r.Committee.Validators {
			totalVotingPower.Add(totalVotingPower, val.BondedStake)
		}
		quorum := bft.Quorum(totalVotingPower)
		fullProofEffort := new(big.Int).Sub(totalVotingPower, quorum) // proposer effort when a full activity proof is provided

		targetHeight := r.Evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer := r.Committee.Validators[0].NodeAddress
		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.T, faultyProposer(r, targetHeight))
		require.Equal(r.T, 0, inactivityCounter(r, proposer))
		require.Equal(r.T, fullProofEffort.String(), proposerEffort(r, proposer).String())
		require.Equal(r.T, fullProofEffort.String(), totalProposerEffort(r).String())

		// finalize 3 more times with full proof
		targetHeight = r.Evm.Context.BlockNumber.Int64() - delta.Int64()
		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.T, faultyProposer(r, targetHeight))
		targetHeight++
		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.T, faultyProposer(r, targetHeight))
		targetHeight++
		setupProofAndAutonityFinalize(r, proposer, nil)
		expectedEffort := new(big.Int).Mul(fullProofEffort, common.Big4) // we finalized 4 times up to now
		require.False(r.T, faultyProposer(r, targetHeight))
		require.Equal(r.T, 0, inactivityCounter(r, proposer))
		require.Equal(r.T, expectedEffort.String(), proposerEffort(r, proposer).String())
		require.Equal(r.T, expectedEffort.String(), totalProposerEffort(r).String())

		targetHeight = r.Evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer = r.Committee.Validators[1].NodeAddress
		setupProofAndAutonityFinalize(r, proposer, nil)
		expectedTotalEffort := new(big.Int).Add(expectedEffort, fullProofEffort) // validators[0] effort + validator[1] effort
		require.False(r.T, faultyProposer(r, targetHeight))
		require.Equal(r.T, 0, inactivityCounter(r, proposer))
		require.Equal(r.T, fullProofEffort.String(), proposerEffort(r, proposer).String())
		require.Equal(r.T, expectedTotalEffort.String(), totalProposerEffort(r).String())
	})
}

// checks that the inactivity counters are correctly updated according to the lookback window (_recordAbsentees function)
func TestInactivityCounter(t *testing.T) {
	t.Run("delta=5, lookback=40", func(t *testing.T) {
		runTestInactivityCounter(t, 5, 40, 130)
	})
	t.Run("delta=2, lookback=6", func(t *testing.T) {
		runTestInactivityCounter(t, 2, 6, 60)
	})
	t.Run("delta=20, lookback=6", func(t *testing.T) {
		runTestInactivityCounter(t, 20, 6, 130)
	})
	t.Run("delta=20, lookback=68", func(t *testing.T) {
		runTestInactivityCounter(t, 20, 68, 200)
	})
	t.Run("delta=34, lookback=37", func(t *testing.T) {
		runTestInactivityCounter(t, 34, 37, 200)
	})
}

func runTestInactivityCounter(t *testing.T, delta int, lookback int, epochPeriod uint64) {
	r := Setup(t, func(genesis *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		customGenesis := configOverrideIncreasedStake(genesis)
		customGenesis.EpochPeriod = epochPeriod
		return customGenesis
	})

	// set maximum inactivity threshold for this test, we care only about the inactivity counters and not about the jailing
	_, err := r.OmissionAccountability.SetInactivityThreshold(r.Operator, new(big.Int).SetUint64(10000))
	require.NoError(t, err)

	applyCustomParams(r, uint64(delta), uint64(lookback))
	t.Logf("last mined block: %d, delta: %d", r.lastMinedHeight(), delta)

	r.WaitNBlocks(delta)

	proposer := r.Committee.Validators[0].NodeAddress
	fullyOffline := r.Committee.Validators[1].NodeAddress
	partiallyOffline := r.Committee.Validators[2].NodeAddress

	absentees := make(map[common.Address]struct{})
	absentees[fullyOffline] = struct{}{}
	partiallyOfflineCounter := 0
	for i := 0; i < lookback-2; i++ {
		if i == lookback/2 {
			absentees[partiallyOffline] = struct{}{}
		}
		targetHeight := r.Evm.Context.BlockNumber.Int64() - int64(delta)
		setupProofAndAutonityFinalize(r, proposer, absentees)
		for absentee := range absentees {
			require.True(r.T, isValidatorInactive(r, targetHeight, absentee))
			if absentee == partiallyOffline {
				partiallyOfflineCounter++
			}
		}
	}

	// we still need two height to have a full lookback window
	// insert a proposer faulty height (no activity proof), it should be ignored
	r.T.Logf("current block number in evm: %d", r.Evm.Context.BlockNumber.Uint64())
	r.Evm.Context.Coinbase = proposer
	r.Evm.Context.ActivityProof = nil
	autonityFinalize(r)
	require.Equal(r.T, 0, inactivityCounter(r, fullyOffline))

	// here we should update the inactivity counter to 1, but since there was a faulty proposer we extend the lookback period
	r.T.Logf("current block number in evm: %d", r.Evm.Context.BlockNumber.Uint64())
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.T, 0, inactivityCounter(r, fullyOffline))

	// now we have a full lookback period
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.T, 1, inactivityCounter(r, fullyOffline))
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.T, 2, inactivityCounter(r, fullyOffline))
	require.Equal(r.T, 0, inactivityCounter(r, partiallyOffline))
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.T, 3, inactivityCounter(r, fullyOffline))
	require.Equal(r.T, 0, inactivityCounter(r, partiallyOffline))

	// fill up enough blocks for partiallyOffline as well
	for i := partiallyOfflineCounter; i < lookback-1; i++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}

	require.Equal(r.T, 0, inactivityCounter(r, partiallyOffline))
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.T, 1, inactivityCounter(r, partiallyOffline))

	fullyOfflineIC := inactivityCounter(r, fullyOffline)
	partiallyOfflineIC := inactivityCounter(r, partiallyOffline)
	// every two block, one has faulty proposer
	n := 20
	for i := 0; i < (n * 2); i++ {
		proposerFaulty := i%2 == 0
		if !proposerFaulty {
			r.setupActivityProofAndCoinbase(proposer, absentees)
		} else {
			r.Evm.Context.Coinbase = proposer
			r.Evm.Context.ActivityProof = nil
		}
		autonityFinalize(r)
	}

	// inactivity counter should still have increased by n due to lookback period extension
	require.Equal(r.T, fullyOfflineIC+n, inactivityCounter(r, fullyOffline))
	require.Equal(r.T, partiallyOfflineIC+n, inactivityCounter(r, partiallyOffline))

	//  close the epoch
	for i := r.Evm.Context.BlockNumber.Int64(); i < int64(epochPeriod)*2; i++ {
		setupProofAndAutonityFinalize(r, proposer, nil)
	}
	t.Log("Closing epoch")
	setupProofAndAutonityFinalize(r, proposer, nil)
	r.generateNewCommittee()

	// inactivity counters should be reset
	require.Equal(r.T, 0, inactivityCounter(r, partiallyOffline))
	require.Equal(r.T, 0, inactivityCounter(r, fullyOffline))

	r.WaitNBlocks(delta)
	t.Logf("current consensus instance for height %d", r.Evm.Context.BlockNumber.Uint64())
	otherValidator := r.Committee.Validators[3].NodeAddress
	newAbsentees := make(map[common.Address]struct{})
	newAbsentees[otherValidator] = struct{}{}

	for i := 0; i < lookback/2; i++ {
		setupProofAndAutonityFinalize(r, proposer, newAbsentees)
	}

	t.Log("online at following block")
	setupProofAndAutonityFinalize(r, proposer, nil)

	// one block online is going to "save" the validator for the next lookback window
	for i := 0; i < lookback-1; i++ {
		setupProofAndAutonityFinalize(r, proposer, newAbsentees)
		require.Equal(r.T, 0, inactivityCounter(r, otherValidator))
	}

	// proposer faulty
	r.Evm.Context.Coinbase = proposer
	r.Evm.Context.ActivityProof = nil
	autonityFinalize(r)
	require.Equal(r.T, 0, inactivityCounter(r, otherValidator))

	setupProofAndAutonityFinalize(r, proposer, newAbsentees)
	require.Equal(r.T, 1, inactivityCounter(r, otherValidator))
}

// checks that the inactivity scores are computed correctly
func TestInactivityScore(t *testing.T) {
	t.Run("delta=5, lookback=40, past=1000", func(t *testing.T) {
		runTestInactivityScore(t, 5, 40, 130, 1000)
	})
	t.Run("delta=2, lookback=6, past=1000", func(t *testing.T) {
		runTestInactivityScore(t, 2, 6, 60, 1000)
	})
	t.Run("delta=20, lookback=6, past=1000", func(t *testing.T) {
		runTestInactivityScore(t, 20, 6, 130, 1000)
	})
	t.Run("delta=20, lookback=68, past=500", func(t *testing.T) {
		runTestInactivityScore(t, 20, 68, 200, 500)
	})
	t.Run("delta=34, lookback=37, past=2500", func(t *testing.T) {
		runTestInactivityScore(t, 34, 37, 200, 2500)
	})
	t.Run("delta=2, lookback=1, past=1000", func(t *testing.T) {
		runTestInactivityScore(t, 2, 1, 60, 1000)
	})
	t.Run("delta=29, lookback=1, past=1000", func(t *testing.T) {
		runTestInactivityScore(t, 2, 1, 60, 1000)
	})
	t.Run("delta=5, lookback=40, past=0", func(t *testing.T) {
		runTestInactivityScore(t, 5, 40, 130, 0)
	})
	t.Run("delta=5, lookback=40, past=7500", func(t *testing.T) {
		runTestInactivityScore(t, 5, 40, 130, 7500)
	})
	t.Run("delta=5, lookback=40, past=10000", func(t *testing.T) {
		runTestInactivityScore(t, 5, 40, 130, 10000)
	})
}

func runTestInactivityScore(t *testing.T, delta int, lookback int, epochPeriod uint64, pastPerformanceWeight uint64) {
	r := Setup(t, func(genesis *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		customGenesis := configOverrideIncreasedStake(genesis)
		customGenesis.EpochPeriod = epochPeriod
		return customGenesis
	})

	// set maximum inactivity threshold for this test, we care only about the inactivity scores and not about the jailing
	scaleFactorInt := omissionScaleFactor(r)
	_, err := r.OmissionAccountability.SetInactivityThreshold(r.Operator, scaleFactorInt)
	require.NoError(t, err)
	scaleFactor := newFloat(scaleFactorInt)

	// set past performance weight
	pastPerformanceWeightBig := new(big.Int).SetUint64(pastPerformanceWeight)
	_, err = r.OmissionAccountability.SetPastPerformanceWeight(r.Operator, pastPerformanceWeightBig)
	require.NoError(t, err)

	applyCustomParams(r, uint64(delta), uint64(lookback))
	t.Logf("last mined block: %d, delta: %d", r.lastMinedHeight(), delta)

	initialCommitteeSize := len(r.Committee.Validators)

	r.WaitNBlocks(delta)

	pastPerformanceWeightFloat := new(big.Float).Quo(newFloat(pastPerformanceWeightBig), scaleFactor)
	currentPerformanceWeight := new(big.Float).Sub(newFloat(common.Big1), pastPerformanceWeightFloat)

	// simulate epoch.
	proposer := r.Committee.Validators[0].NodeAddress
	inactiveBlockStreak := make(map[common.Address]int)
	inactiveCounters := make(map[common.Address]*big.Int)
	for _, v := range r.Committee.Validators {
		inactiveCounters[v.NodeAddress] = new(big.Int)
	}
	for h := delta + 1; h < int(epochPeriod)+1; h++ {
		absentees := make(map[common.Address]struct{})
		for _, v := range r.Committee.Validators {
			if v.NodeAddress == proposer {
				continue // keep proposer always online
			}
			if rand.Intn(30) != 0 {
				absentees[v.NodeAddress] = struct{}{}
				inactiveBlockStreak[v.NodeAddress]++
			} else {
				inactiveBlockStreak[v.NodeAddress] = 0
			}
			if inactiveBlockStreak[v.NodeAddress] >= lookback {
				inactiveCounters[v.NodeAddress] = new(big.Int).Add(inactiveCounters[v.NodeAddress], common.Big1)
			}
		}

		t.Logf("number of absentees: %d for height %d", len(absentees), r.Evm.Context.BlockNumber.Uint64())
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	// no one should have been jailed
	require.Equal(t, initialCommitteeSize, len(r.Committee.Validators))

	// check score computation
	pastInactivityScore := make(map[common.Address]*big.Float)
	denominator := new(big.Int).SetUint64(epochPeriod - uint64(delta) - uint64(lookback) + 1)
	for _, val := range r.Committee.Validators {
		score := new(big.Float).Quo(newFloat(inactiveCounters[val.NodeAddress]), newFloat(denominator))

		expectedInactivityScoreFloat := new(big.Float).Mul(score, currentPerformanceWeight) // + 0 * pastPerformanceWeight
		pastInactivityScore[val.NodeAddress] = expectedInactivityScoreFloat

		expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, scaleFactor)
		expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Int(nil)

		r.T.Logf("validator %s, expectedInactivityScoreFloat %s, expectedInactivityScore %s, inactivityScore %v", common.Bytes2Hex(val.NodeAddress[:]), toString(expectedInactivityScoreFloat), expectedInactivityScore.String(), inactivityScore(r, val.NodeAddress))

		require.Equal(r.T, int(expectedInactivityScore.Uint64()), inactivityScore(r, val.NodeAddress))
	}

	// simulate another epoch
	r.WaitNBlocks(delta)
	inactiveBlockStreak = make(map[common.Address]int)
	inactiveCounters = make(map[common.Address]*big.Int)
	for _, v := range r.Committee.Validators {
		inactiveCounters[v.NodeAddress] = new(big.Int)
	}
	for h := delta + 1; h < int(epochPeriod)+1; h++ {
		absentees := make(map[common.Address]struct{})
		for _, v := range r.Committee.Validators {
			if v.NodeAddress == proposer {
				continue // keep proposer always online
			}
			if rand.Intn(30) != 0 {
				absentees[v.NodeAddress] = struct{}{}
				inactiveBlockStreak[v.NodeAddress]++
			} else {
				inactiveBlockStreak[v.NodeAddress] = 0
			}
			if inactiveBlockStreak[v.NodeAddress] >= lookback {
				inactiveCounters[v.NodeAddress] = new(big.Int).Add(inactiveCounters[v.NodeAddress], common.Big1)
			}
		}

		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	// no one should have been jailed
	require.Equal(t, initialCommitteeSize, len(r.Committee.Validators))

	// check score computation
	for _, val := range r.Committee.Validators {
		score := new(big.Float).Quo(newFloat(inactiveCounters[val.NodeAddress]), newFloat(denominator))

		expectedInactivityScoreFloat1 := new(big.Float).Mul(score, currentPerformanceWeight)
		expectedInactivityScoreFloat2 := new(big.Float).Mul(pastInactivityScore[val.NodeAddress], pastPerformanceWeightFloat)
		expectedInactivityScoreFloat := new(big.Float).Add(expectedInactivityScoreFloat1, expectedInactivityScoreFloat2)

		expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, scaleFactor)
		expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Int(nil)

		r.T.Logf("validator %s, expectedInactivityScoreFloat %s, expectedInactivityScore %s, inactivityScore %v", common.Bytes2Hex(val.NodeAddress[:]), toString(expectedInactivityScoreFloat), expectedInactivityScore.String(), inactivityScore(r, val.NodeAddress))
		// precision loss cumulates across subsequent epochs due to the weighted sum, thus as there can be some instability
		// in this test, we allow a tolerance of 1
		expectNearlyEqual(t, expectedInactivityScore, big.NewInt(int64(inactivityScore(r, val.NodeAddress))), 1)
	}
}

func TestOmissionPunishments(t *testing.T) {
	r := Setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		config.EpochPeriod = uint64(omissionEpochPeriod)
		// increase voting power of validator 0 to reach quorum in proofs easily
		config.Validators[0].BondedStake = new(big.Int).Mul(config.Validators[1].BondedStake, big.NewInt(6))
		return config
	})

	// deploy and set the AccountabilityTest contract. We will need write access to the beneficiaries map later
	_, _, accountabilityTest, err := r.DeployAccountabilityTest(nil, r.Autonity.address, AccountabilityConfig{
		InnocenceProofSubmissionWindow: big.NewInt(int64(params.DefaultAccountabilityConfig.InnocenceProofSubmissionWindow)),
		BaseSlashingRates: AccountabilityBaseSlashingRates{
			Low:  big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateLow)),
			Mid:  big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateMid)),
			High: big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateHigh)),
		},
		Factors: AccountabilityFactors{
			Collusion: big.NewInt(int64(params.DefaultAccountabilityConfig.CollusionFactor)),
			History:   big.NewInt(int64(params.DefaultAccountabilityConfig.HistoryFactor)),
			Jail:      big.NewInt(int64(params.DefaultAccountabilityConfig.JailFactor)),
		},
	})
	require.NoError(t, err)
	r.NoError(
		r.Autonity.SetAccountabilityContract(r.Operator, accountabilityTest.address),
	)

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.WaitNBlocks(int(delta.Int64()))

	config, _, err := r.OmissionAccountability.Config(nil)
	require.NoError(r.T, err)
	initialJailingPeriod := int(config.InitialJailingPeriod.Uint64())
	initialProbationPeriod := int(config.InitialProbationPeriod.Uint64())
	pastPerformanceWeight := int(config.PastPerformanceWeight.Uint64())
	initialSlashingRate := int(config.InitialSlashingRate.Uint64())
	scaleFactor := int(omissionScaleFactor(r).Uint64())

	proposer := r.Committee.Validators[0].NodeAddress
	absentees := make(map[common.Address]struct{})
	val1Address := r.Committee.Validators[1].NodeAddress // will be handy to have those later
	val2Address := r.Committee.Validators[2].NodeAddress
	absentees[val1Address] = struct{}{}
	absentees[val2Address] = struct{}{}
	val1Treasury := r.Committee.Validators[1].Treasury
	val2Treasury := r.Committee.Validators[2].Treasury

	// simulate epoch with two validator at 100% inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	// close the epoch
	setupProofAndAutonityFinalize(r, proposer, absentees)
	r.generateNewCommittee()

	// the two validators should have been jailed and be under probation + offence counter should have been incremented
	expectedFullOfflineScore := scaleFactor - pastPerformanceWeight
	for absentee := range absentees {
		require.Equal(r.T, expectedFullOfflineScore, inactivityScore(r, absentee))
		val := validator(r, absentee)
		require.Equal(r.T, jailedForInactivity, val.State)
		require.Equal(r.T, uint64(omissionEpochPeriod+initialJailingPeriod), val.JailReleaseBlock.Uint64())
		require.Equal(r.T, initialProbationPeriod, probation(r, absentee))
		require.Equal(t, 1, offences(r, absentee))
	}

	// wait that the jailing finishes and reactivate validators
	r.WaitNBlocks(initialJailingPeriod)
	_, err = r.Autonity.ActivateValidator(&runOptions{origin: val1Treasury}, val1Address)
	require.NoError(r.T, err)
	_, err = r.Autonity.ActivateValidator(&runOptions{origin: val2Treasury}, val2Address)
	require.NoError(r.T, err)

	// pass some epochs, probation period should decrease
	r.WaitNextEpoch() // re-activation epoch, val not part of committee
	// inactivity score should still be the same as before
	for absentee := range absentees {
		require.Equal(r.T, expectedFullOfflineScore, inactivityScore(r, absentee))
	}
	r.WaitNextEpoch()
	// should be decreased  now
	for absentee := range absentees {
		require.Equal(r.T, (expectedFullOfflineScore*pastPerformanceWeight)/scaleFactor, inactivityScore(r, absentee))
	}
	r.WaitNextEpoch()
	r.WaitNextEpoch()

	// probation periods should have decreased of once for every epoch that passed with the validator as part of the committee
	passedEpochs := 3

	for absentee := range absentees {
		require.Equal(t, initialProbationPeriod-passedEpochs, probation(r, absentee))
	}

	// simulate another epoch where:
	// - val 1 gets slashed by accountability and therefore doesn't get punished by omission accountability
	// - val 2 gets punished again for omission while in the probation period, therefore he gets slashed
	r.WaitNBlocks(int(delta.Int64()))
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	val1 := validator(r, val1Address)
	totalSlashedVal1 := val1.TotalSlashed
	_, err = r.Autonity.Jail(&runOptions{origin: accountabilityTest.address}, val1Address, new(big.Int).SetUint64(uint64(omissionEpochPeriod*10)), jailed)
	require.NoError(t, err)
	_, err = accountabilityTest.AddBeneficiary(nil, val1Address, proposer)
	require.NoError(t, err)

	val2BeforeSlash := validator(r, val2Address)
	// close epoch
	setupProofAndAutonityFinalize(r, proposer, absentees)
	r.generateNewCommittee()

	// val1, punished by accountability, shouldn't have been slashed by omission even if 100% offline and still under probation
	val1 = validator(r, val1Address)
	require.Equal(r.T, jailed, val1.State)
	require.True(r.T, probation(r, val1.NodeAddress) > 0)
	require.Equal(r.T, totalSlashedVal1.String(), val1.TotalSlashed.String())
	require.Equal(r.T, 1, offences(r, val1Address))

	// val2 offline while on probation, should have been slashed by omission
	val2 := validator(r, val2Address)
	require.Equal(r.T, jailedForInactivity, val2.State)
	require.True(r.T, probation(r, val2.NodeAddress) > 0)
	require.True(r.T, val2.TotalSlashed.Cmp(val2BeforeSlash.TotalSlashed) > 0)
	require.Equal(r.T, 2, offences(r, val2Address))
	expectedSlashRate := new(big.Int).SetInt64(int64(initialSlashingRate * 4 * 2)) // rate * offence^2 * collusion
	availableFunds := new(big.Int).Add(val2BeforeSlash.BondedStake, val2.UnbondingStake)
	availableFunds.Add(availableFunds, val2.SelfUnbondingStake)
	expectedSlashAmount := new(big.Int).Mul(expectedSlashRate, availableFunds)
	expectedSlashAmount.Div(expectedSlashAmount, slashingRateScaleFactor(r))
	t.Logf("expected slash rate: %s, available funds: %s, expected slash amount: %s", expectedSlashRate.String(), availableFunds.String(), expectedSlashAmount.String())
	require.Equal(r.T, expectedSlashAmount.String(), new(big.Int).Sub(val2.TotalSlashed, val2BeforeSlash.TotalSlashed).String())
}

func TestProposerRewardDistribution(t *testing.T) {
	t.Run("Rewards are correctly allocated based on config", func(t *testing.T) {
		r := Setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			config.EpochPeriod = uint64(omissionEpochPeriod)
			config.TreasuryFee = 0
			return config
		})

		// deploy and set the AccountabilityTest contract. We will need write access to the beneficiaries map later
		_, _, accountabilityTest, err := r.DeployAccountabilityTest(nil, r.Autonity.address, AccountabilityConfig{
			InnocenceProofSubmissionWindow: big.NewInt(int64(params.DefaultAccountabilityConfig.InnocenceProofSubmissionWindow)),
			BaseSlashingRates: AccountabilityBaseSlashingRates{
				Low:  big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateLow)),
				Mid:  big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateMid)),
				High: big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateHigh)),
			},
			Factors: AccountabilityFactors{
				Collusion: big.NewInt(int64(params.DefaultAccountabilityConfig.CollusionFactor)),
				History:   big.NewInt(int64(params.DefaultAccountabilityConfig.HistoryFactor)),
				Jail:      big.NewInt(int64(params.DefaultAccountabilityConfig.JailFactor)),
			},
		})
		require.NoError(t, err)
		r.NoError(
			r.Autonity.SetAccountabilityContract(r.Operator, accountabilityTest.address),
		)

		maxCommitteeSizeBig, _, err := r.Autonity.GetMaxCommitteeSize(nil)
		require.NoError(r.T, err)
		maxCommitteeSize := newFloat(maxCommitteeSizeBig)
		t.Logf("max committee size: %s", toString(maxCommitteeSize))

		config, _, err := r.Autonity.Config(nil)
		require.NoError(r.T, err)
		proposerRewardRateBig := config.Policy.ProposerRewardRate
		proposerRewardRate := newFloat(proposerRewardRateBig)
		standardScaleFactorBig, _, err := r.Autonity.STANDARDSCALEFACTOR(nil)
		require.NoError(t, err)
		standardScaleFactor := newFloat(standardScaleFactorBig)
		t.Logf("proposer reward rate: %s, scale factor: %s", toString(proposerRewardRate), toString(standardScaleFactor))

		autonityAtnsBig := new(big.Int).SetUint64(54644455456467) // random amount
		t.Logf("atn rewards: %s", autonityAtnsBig.String())
		// this has to match the ntn inflation unlocked NTNs.
		// Can be retrieved by adding in solidity a revert(Helpers.toString(accounts[address(this)])); in Finalize
		ntnRewardsBig := new(big.Int).SetUint64(8205384319979600000)
		t.Logf("ntn rewards: %s", ntnRewardsBig.String())
		r.GiveMeSomeMoney(r.Autonity.address, autonityAtnsBig)

		autonityAtns := newFloat(autonityAtnsBig)
		ntnRewards := newFloat(ntnRewardsBig)

		// all rewards should go to val 0
		proposer := r.Committee.Validators[0].NodeAddress
		proposerTreasury := r.Committee.Validators[0].Treasury
		atnBalanceBefore := newFloat(r.GetBalanceOf(proposerTreasury))
		stakesBefore := newFloat(r.Committee.Validators[0].BondedStake)
		t.Logf("atn balance before: %s, ntn balance before %s", toString(atnBalanceBefore), toString(stakesBefore))

		// set validator state to jailed so that he will not receive any reward other the proposer one
		_, err = r.Autonity.Jail(&runOptions{origin: accountabilityTest.address}, proposer, new(big.Int).SetUint64(uint64(omissionEpochPeriod*10)), jailed)
		require.NoError(t, err)
		_, err = accountabilityTest.AddBeneficiary(nil, proposer, r.Committee.Validators[1].NodeAddress)
		require.NoError(t, err)

		r.Evm.Context.BlockNumber = new(big.Int).SetInt64(int64(omissionEpochPeriod))
		r.Evm.Context.Time.Add(r.Evm.Context.Time, new(big.Int).SetInt64(int64(omissionEpochPeriod-1)))
		setupProofAndAutonityFinalize(r, proposer, nil)

		committeeSize := newFloat(new(big.Int).SetUint64(uint64(len(r.Committee.Validators))))
		t.Logf("proposer reward rate: %s, committee size: %s", toString(proposerRewardRate), toString(committeeSize))
		numeratorFactor := new(big.Float).Mul(proposerRewardRate, committeeSize)
		t.Logf("numeratorFactor: %s", toString(numeratorFactor))

		t.Logf("proposer reward rate scale factor: %s, max committee size: %s", toString(standardScaleFactor), toString(maxCommitteeSize))
		denominator := new(big.Float).Mul(standardScaleFactor, maxCommitteeSize)
		t.Logf("denominator: %s", toString(denominator))

		numeratorAtn := new(big.Float).Mul(autonityAtns, numeratorFactor)
		t.Logf("numerator atn: %s", toString(numeratorAtn))
		atnExpectedReward := new(big.Float).Quo(numeratorAtn, denominator)

		numeratorNtn := new(big.Float).Mul(ntnRewards, numeratorFactor)
		t.Logf("numerator ntn: %s", toString(numeratorNtn))
		ntnExpectedReward := new(big.Float).Quo(numeratorNtn, denominator)

		t.Logf("atn expected reward: %s, ntn expected reward: %s", toString(atnExpectedReward), toString(ntnExpectedReward))

		atnExpectedBalance := new(big.Float).Add(atnBalanceBefore, atnExpectedReward)
		stakesExpected := new(big.Float).Add(stakesBefore, ntnExpectedReward)

		t.Logf("atn expected balance: %s, ntn expected balance: %s", toString(atnExpectedBalance), toString(stakesExpected))

		atnActualBalance := r.GetBalanceOf(proposerTreasury)
		validatorInfo, _, err := r.Autonity.GetValidator(nil, proposer)
		require.NoError(t, err)
		stakesActual := validatorInfo.BondedStake
		t.Logf("atn actual balance: %s, ntn actual balance: %s", atnActualBalance.String(), stakesActual.String())

		atnExpectedBalanceInt, _ := atnExpectedBalance.Int(nil)
		stakesExpectedInt, _ := stakesExpected.Int(nil)

		t.Logf("balance diff atn: %s", new(big.Int).Sub(atnExpectedBalanceInt, atnActualBalance).String())
		t.Logf("balance diff ntn: %s", new(big.Int).Sub(stakesExpectedInt, stakesActual).String())

		require.Equal(t, atnExpectedBalanceInt.String(), atnActualBalance.String())
		require.Equal(t, stakesExpectedInt.String(), stakesActual.String())
	})

	t.Run("Rewards are correctly distributed among proposers", func(t *testing.T) {
		r := Setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			config.EpochPeriod = uint64(omissionEpochPeriod)
			config.TreasuryFee = 0
			// set max committee and ProposerRewardRate so that all rewards go to proposer
			config.ProposerRewardRate = 10_000
			config.OracleRewardRate = 0
			config.MaxCommitteeSize = 4
			return config
		})

		delta, _, err := r.OmissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.WaitNBlocks(int(delta.Int64()))

		totalEffort := new(big.Int)
		efforts := make(map[common.Address]*big.Int)
		atnBalances := make(map[common.Address]*big.Int)
		stakes := make(map[common.Address]*big.Int)
		totalPower := new(big.Int)
		for _, val := range r.Committee.Validators {
			efforts[val.NodeAddress] = new(big.Int)
			atnBalances[val.NodeAddress] = r.GetBalanceOf(val.Treasury)
			stakes[val.NodeAddress] = val.BondedStake
			totalPower.Add(totalPower, val.BondedStake)
		}
		simulatedAtnRewards := new(big.Int).SetInt64(4545445)
		r.GiveMeSomeMoney(r.Autonity.address, simulatedAtnRewards)
		simulatedNtnRewards := r.RewardsAfterOneEpoch().RewardNTN
		// simulate epoch
		fullProofEffort := new(big.Int).Sub(totalPower, bft.Quorum(totalPower)) // effort for a full proof
		for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
			proposerIndex := rand.Intn(len(r.Committee.Validators))
			totalEffort.Add(totalEffort, fullProofEffort)
			nodeAddress := r.Committee.Validators[proposerIndex].NodeAddress
			efforts[nodeAddress].Add(efforts[nodeAddress], fullProofEffort)
			r.setupActivityProofAndCoinbase(nodeAddress, nil)
			r.FinalizeBlock()
			// omissionFinalize(r, h == omissionEpochPeriod)
		}
		r.generateNewCommittee()

		// _, err = r.Autonity.Mint(r.Operator, r.OmissionAccountability.address, simulatedNtnRewards)
		// require.NoError(r.T, err)
		// _, err = r.OmissionAccountability.DistributeProposerRewards(&runOptions{origin: r.Autonity.address, value: simulatedAtnRewards}, simulatedNtnRewards)
		// require.NoError(t, err)

		for _, val := range r.Committee.Validators {
			atnExpectedIncrement := new(big.Int).Mul(efforts[val.NodeAddress], simulatedAtnRewards)
			atnExpectedIncrement.Div(atnExpectedIncrement, totalEffort)
			ntnExpectedIncrement := new(big.Int).Mul(efforts[val.NodeAddress], simulatedNtnRewards)
			ntnExpectedIncrement.Div(ntnExpectedIncrement, totalEffort)
			atnExpectedBalance := new(big.Int).Add(atnBalances[val.NodeAddress], atnExpectedIncrement)
			stakesExpected := new(big.Int).Add(stakes[val.NodeAddress], ntnExpectedIncrement)
			t.Logf("validator %s, effort %s, total effort %s, expectedBalance atn %s, expectedStakes %s", common.Bytes2Hex(val.NodeAddress[:]), efforts[val.NodeAddress].String(), totalEffort.String(), atnExpectedBalance.String(), stakesExpected.String())

			atnBalance := r.GetBalanceOf(val.Treasury)

			require.Equal(t, atnExpectedBalance.String(), atnBalance.String())
			require.Equal(t, stakesExpected.String(), val.BondedStake.String())

			// effort counters should be zeroed out
			require.Equal(r.T, common.Big0.String(), proposerEffort(r, val.NodeAddress).String())
		}

		require.Equal(r.T, common.Big0.String(), totalProposerEffort(r).String())
	})
}

// past performance weight and inactivity threshold should be set low enough that if:
// - a validator gets 100% inactivity in epoch x
// - then he gets 0% inactivity in epoch x+n (after he reactivated)
// he shouldn't get slashed in epoch x+n
func TestConfigSanity(t *testing.T) {
	r := Setup(t, configOverrideIncreasedStake)

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.WaitNBlocks(int(delta.Int64()))

	config, _, err := r.OmissionAccountability.Config(nil)
	require.NoError(r.T, err)
	initialJailingPeriod := int(config.InitialJailingPeriod.Uint64())

	proposer := r.Committee.Validators[0].NodeAddress
	absentees := make(map[common.Address]struct{})
	val1Address := r.Committee.Validators[1].NodeAddress // will be handy later
	val2Address := r.Committee.Validators[2].NodeAddress // will be handy later
	absentees[val1Address] = struct{}{}
	absentees[val2Address] = struct{}{}
	val1Treasury := r.Committee.Validators[1].Treasury
	val2Treasury := r.Committee.Validators[2].Treasury

	// simulate epoch with two validator at 100% inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	for absentee := range absentees {
		val := validator(r, absentee)
		require.Equal(r.T, jailedForInactivity, val.State)
		require.Equal(t, 1, offences(r, absentee))
	}

	// wait that the jailing finishes and reactivate validators
	r.WaitNBlocks(initialJailingPeriod)
	_, err = r.Autonity.ActivateValidator(&runOptions{origin: val1Treasury}, val1Address)
	require.NoError(r.T, err)
	_, err = r.Autonity.ActivateValidator(&runOptions{origin: val2Treasury}, val2Address)
	require.NoError(r.T, err)

	r.WaitNextEpoch() // re-activation epoch, val not part of committee
	r.WaitNextEpoch()

	// validator should not have been punished since he did 0% offline
	for absentee := range absentees {
		val := validator(r, absentee)
		require.Equal(r.T, active, val.State)
		require.Equal(t, 1, offences(r, absentee))
	}

}

func TestRewardWithholding(t *testing.T) {
	r := Setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		config.EpochPeriod = uint64(omissionEpochPeriod)
		config.ProposerRewardRate = 0 // no rewards to proposers to make computation simpler
		config.TreasuryFee = 0        // same
		// increase voting power of validator 0 to reach quorum in proofs easily
		config.Validators[0].BondedStake = new(big.Int).Mul(config.Validators[1].BondedStake, big.NewInt(6))
		config.OracleRewardRate = 0
		return config
	})

	extractNodeAddresses := func(vals []AutonityValidator) []common.Address {
		addresses := make([]common.Address, 0)
		for _, val := range vals {
			addresses = append(addresses, val.NodeAddress)
		}
		return addresses
	}

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	// validators over threshold will get all their rewards withheld
	customInactivityThreshold := uint64(4000)
	_, err = r.OmissionAccountability.SetInactivityThreshold(r.Operator, new(big.Int).SetUint64(customInactivityThreshold))
	require.NoError(t, err)

	r.WaitNBlocks(int(delta.Int64()))

	config, _, err := r.Autonity.Config(nil)
	require.NoError(t, err)
	withheldRewardPool := config.Policy.WithheldRewardsPool

	proposer := r.Committee.Validators[0].NodeAddress

	// simulate epoch with random levels of inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		absentees := make(map[common.Address]struct{})
		for i := range r.Committee.Validators {
			if i == 0 {
				continue // let's keep at least one guy inside the committee
			}
			if rand.Intn(30) != 0 {
				absentees[r.Committee.Validators[i].NodeAddress] = struct{}{}
			}
		}
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}

	atnRewards := new(big.Int).SetUint64(5467879877987) // random amount
	// this has to match the ntn inflation unlocked NTNs.
	// Can be retrieved by adding in solidity a revert(Helpers.toString(accounts[address(this)])); in Finalize
	// Also can be used `r.RewardsAfterOneEpoch()` function to calculate rewards
	ntnRewards := new(big.Int).SetUint64(8220842843566600000)
	r.GiveMeSomeMoney(r.Autonity.address, atnRewards)

	atnBalancesBefore := make(map[common.Address]*big.Int)
	stakesBefore := make(map[common.Address]*big.Int)
	totalPower := new(big.Int)
	for _, val := range r.Committee.Validators {
		validatorStruct := validator(r, val.NodeAddress)
		// we assume that all stake is self bonded in this test
		require.Equal(t, validatorStruct.SelfBondedStake.String(), validatorStruct.BondedStake.String())
		atnBalancesBefore[val.NodeAddress] = r.GetBalanceOf(val.Treasury)
		stakesBefore[val.NodeAddress] = validatorStruct.BondedStake
		t.Logf("validator %s, atn balance before: %s, bonded stake before balance before %s", common.Bytes2Hex(val.NodeAddress[:]), atnBalancesBefore[val.NodeAddress].String(), stakesBefore[val.NodeAddress].String())
		totalPower.Add(totalPower, validatorStruct.SelfBondedStake)
	}
	atnPoolBefore := r.GetBalanceOf(withheldRewardPool)
	ntnPoolBefore := ntnBalance(r, withheldRewardPool)

	nodeAddresses := extractNodeAddresses(r.Committee.Validators) // save node addresses as some vals might get jailed if over inactivity threshold
	setupProofAndAutonityFinalize(r, proposer, nil)
	r.generateNewCommittee()

	atnTotalWithheld := new(big.Int)
	ntnTotalWithheld := new(big.Int)
	for _, nodeAddress := range nodeAddresses {
		val := validator(r, nodeAddress)
		power := stakesBefore[val.NodeAddress]

		// compute reward without withholding
		atnFullReward := new(big.Int).Mul(power, atnRewards)
		atnFullReward.Div(atnFullReward, totalPower)
		ntnFullReward := new(big.Int).Mul(power, ntnRewards)
		ntnFullReward.Div(ntnFullReward, totalPower)

		// compute withheld amount
		score := new(big.Int).SetInt64(int64(inactivityScore(r, val.NodeAddress)))
		var ntnWithheld *big.Int
		var atnWithheld *big.Int
		t.Logf("validator %s, score: %d", common.Bytes2Hex(val.NodeAddress[:]), score.Uint64())
		if score.Uint64() <= customInactivityThreshold {
			atnWithheld = new(big.Int).Mul(atnFullReward, score)
			atnWithheld.Div(atnWithheld, omissionScaleFactor(r))
			ntnWithheld = new(big.Int).Mul(ntnFullReward, score)
			ntnWithheld.Div(ntnWithheld, omissionScaleFactor(r))
		} else {
			// all rewards are withheld
			atnWithheld = new(big.Int).Set(atnFullReward)
			ntnWithheld = new(big.Int).Set(ntnFullReward)
		}
		t.Logf("validator %s, atn withheld %d, ntn withheld %d", common.Bytes2Hex(val.NodeAddress[:]), atnWithheld.Uint64(), ntnWithheld.Uint64())
		atnTotalWithheld.Add(atnTotalWithheld, atnWithheld)
		ntnTotalWithheld.Add(ntnTotalWithheld, ntnWithheld)

		// check validator balance
		atnExpectedBalance := new(big.Int).Add(atnFullReward, atnBalancesBefore[val.NodeAddress])
		atnExpectedBalance.Sub(atnExpectedBalance, atnWithheld)
		stakesExpected := new(big.Int).Add(ntnFullReward, stakesBefore[val.NodeAddress])
		stakesExpected.Sub(stakesExpected, ntnWithheld)
		require.Equal(t, atnExpectedBalance.String(), r.GetBalanceOf(val.Treasury).String())
		require.Equal(t, stakesExpected.String(), val.BondedStake.String())
	}
	atnExpectedPoolBalance := atnPoolBefore.Add(atnPoolBefore, atnTotalWithheld)
	ntnExpectedPoolBalance := ntnPoolBefore.Add(ntnPoolBefore, ntnTotalWithheld)
	require.Equal(t, atnExpectedPoolBalance.String(), r.GetBalanceOf(withheldRewardPool).String())
	require.Equal(t, ntnExpectedPoolBalance.String(), ntnBalance(r, withheldRewardPool).String())
}

// operator can disable omission punishments by increasing the inactivity threshold
// there are also other ways to disable accountability, by lowering the initial punishments values
// proposer rewards can be disabled by lowering the proposer reward rate in the autonity contract
func TestOmissionDisabling(t *testing.T) {
	r := Setup(t, configOverride)

	_, err := r.OmissionAccountability.SetInactivityThreshold(r.Operator, omissionScaleFactor(r))
	require.NoError(t, err)

	// validator 1 absent for entire epoch
	absentees := make(map[common.Address]struct{})
	absentees[r.Committee.Validators[1].NodeAddress] = struct{}{}

	csize := len(r.Committee.Validators)

	for i := 0; i < omissionEpochPeriod; i++ {
		setupProofAndAutonityFinalize(r, r.Committee.Validators[0].NodeAddress, absentees)
	}

	epochID, _, err := r.Autonity.EpochID(nil)
	require.NoError(t, err)
	require.Equal(t, common.Big1.String(), epochID.String())

	// validator 1 should still be in the committee and not jailed
	require.Equal(t, csize, len(r.Committee.Validators))
	val := validator(r, r.Committee.Validators[1].NodeAddress)
	require.Equal(t, uint8(0), val.State)
	require.Equal(t, 0, offences(r, r.Committee.Validators[1].NodeAddress))
}

func TestProtocolParameterChange(t *testing.T) {
	r := Setup(t, nil)
	/* default config:
	- epoch period: 50
	- lookback window: 40
	- delta: 5
	*/
	lookback, _, err := r.OmissionAccountability.GetLookbackWindow(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(40), lookback.Uint64())
	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(5), delta.Uint64())
	epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(50), epochPeriod.Uint64())

	// past performance weight cannot be greater than inactivity threshold
	_, err = r.OmissionAccountability.SetPastPerformanceWeight(r.Operator, omissionScaleFactor(r))
	t.Log(err)
	require.Error(t, err)

	// equation epochPeriod > delta+lookback-1 needs to be respected

	/*
		current params:
		- epoch period: 50
		- delta: 5
		- lookback window: 40
		updated attempted equation:
		30 > 5+40-1 --> false --> err
	*/
	_, err = r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).SetUint64(30))
	t.Log(err)
	require.Error(t, err)

	/*
		current params:
		- epoch period: 50
		- delta: 5
		- lookback window: 40
		updated attempted equation:
		100 > 5+40-1 --> true --> no error
	*/
	_, err = r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).SetUint64(100))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 100
		- delta: 5
		- lookback window: 40
		updated attempted equation:
		100 > 5+100-1 --> false --> err
	*/
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(100))
	t.Log(err)
	require.Error(t, err)

	/*
		current params:
		- epoch period: 100
		- delta: 5
		- lookback window: 40
		updated attempted equation:
		100 > 5+60-1 --> true --> no err
		this will pass because we already increase the period for next epoch, although equation would not be respected for the current epoch
	*/
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(60))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 100
		- delta: 5
		- lookback window: 60
		updated attempted equation:
		30 > 5+20-1 --> true --> no error
	*/
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(20))
	require.NoError(t, err)
	_, err = r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).SetUint64(30))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 30
		- delta: 5
		- lookback window: 20
		updated attempted equation:
		30 > 20+20-1 --> false --> err
	*/
	_, err = r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(20))
	t.Log(err)
	require.Error(t, err)

	/*
		current params:
		- epoch period: 30
		- delta: 5
		- lookback window: 20
		updated attempted equation:
		30 > 10+20-1 --> true --> no err
	*/
	_, err = r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(10))
	require.NoError(t, err)

	// params should still be unchanged in current epoch
	config, _, err := r.OmissionAccountability.Config(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(40), config.LookbackWindow.Uint64())
	require.Equal(t, uint64(5), config.Delta.Uint64())
	epochPeriod, _, err = r.Autonity.GetCurrentEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(50), epochPeriod.Uint64())

	// getters should already return new values
	lookback, _, err = r.OmissionAccountability.GetLookbackWindow(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(20), lookback.Uint64())
	delta, _, err = r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), delta.Uint64())
	epochPeriod, _, err = r.Autonity.GetEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(30), epochPeriod.Uint64())

	// end epoch, new protocol params should be applied
	// both getters and config should return new values
	r.WaitNextEpoch()

	config, _, err = r.OmissionAccountability.Config(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(20), config.LookbackWindow.Uint64())
	require.Equal(t, uint64(10), config.Delta.Uint64())
	epochPeriod, _, err = r.Autonity.GetCurrentEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(30), epochPeriod.Uint64())

	lookback, _, err = r.OmissionAccountability.GetLookbackWindow(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(20), lookback.Uint64())
	delta, _, err = r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), delta.Uint64())
	epochPeriod, _, err = r.Autonity.GetEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(30), epochPeriod.Uint64())

	/*
		current params:
		- epoch period: 30
		- delta: 10
		- lookback window: 20
		updated attempted equation:
		40 > 20+20-1 --> true --> no err
	*/
	_, err = r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).SetUint64(40))
	require.NoError(t, err)
	_, err = r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(20))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 40
		- delta: 20
		- lookback window: 20
		updated attempted equation:
		400 > 20+20-1 --> true --> no err
	*/
	_, err = r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).SetUint64(400))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 400
		- delta: 20
		- lookback window: 20
		updated attempted equation:
		400 > 20+200-1 --> true --> no err
	*/
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(200))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 400
		- delta: 20
		- lookback window: 200
		updated attempted equation:
		400 > 2+200-1 --> true --> no err
	*/
	_, err = r.OmissionAccountability.SetDelta(r.Operator, common.Big2)
	require.NoError(t, err)

	// cannot set lookback and delta to 0
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(0))
	t.Log(err)
	require.Error(t, err)
	_, err = r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(0))
	t.Log(err)
	require.Error(t, err)

	// cannot set delta to 1, but can set lookback window to 1
	_, err = r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(1))
	t.Log(err)
	require.Error(t, err)
	_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(1))
	require.NoError(t, err)

}

func TestFaultyProposerCount(t *testing.T) {
	t.Run("standard scenario", func(t *testing.T) {
		r := Setup(t, configOverride)

		customDelta := uint64(5)
		customLookback := uint64(5)
		_, err := r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(customDelta))
		require.NoError(t, err)
		_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(customLookback))
		require.NoError(t, err)

		// close the epoch to apply new omission params
		r.WaitNextEpoch()
		t.Logf("Starting new epoch at block: %d", r.Evm.Context.BlockNumber.Int64())

		r.WaitNBlocks(int(customDelta) * 2)

		// NOTE: blocks numbers mentioned below are relative to the epoch block of the past epoch
		// leave activity proof empty at block 11, proposer is faulty
		autonityFinalize(r)
		require.Equal(t, uint64(1), faultyProposers(r))
		require.True(t, faultyProposer(r, int64(r.lastTargetHeight())))
		t.Logf("proposer faulty for height: %d, target height: %d", r.lastMinedHeight(), r.lastTargetHeight())

		// mine three blocks with non-empty activity proof
		r.WaitNBlocks(3)

		// another empty proof
		autonityFinalize(r)
		require.Equal(t, uint64(2), faultyProposers(r))
		require.True(t, faultyProposer(r, int64(r.lastTargetHeight())))
		t.Logf("proposer faulty for height: %d, target height: %d", r.lastMinedHeight(), r.lastTargetHeight())

		// one more non-empty proof
		r.WaitNBlocks(1)

		/* so in terms of targetHeight, we are now in this situation:
		*  block 1-5: ok
		*  block 6: faulty
		*  block 7: ok
		*  block 8: ok
		*  block 9: ok
		*  block 10: faulty
		*  block 11: ok
		*
		* For block 11, if we don't count the faulty proposer extension, the lookback window should be (6,11]
		* However since block 10 is faulty, it will get extended to (5,11]. However also block 6 is faulty,
		* So it should get further extended to (4,11]. Therefore the number of faulty proposers should still be 2
		 */

		require.Equal(t, uint64(2), faultyProposers(r))

		// now it should decrease to 1, as the window will be (7,12]. Then gets extended to (6,12] due to block 10
		r.WaitNBlocks(1)
		require.Equal(t, uint64(1), faultyProposers(r))

	})
	t.Run("multiple faulty proposer at the tail of the window", func(t *testing.T) {
		r := Setup(t, configOverride)

		customDelta := uint64(5)
		customLookback := uint64(5)
		_, err := r.OmissionAccountability.SetDelta(r.Operator, new(big.Int).SetUint64(customDelta))
		require.NoError(t, err)
		_, err = r.OmissionAccountability.SetLookbackWindow(r.Operator, new(big.Int).SetUint64(customLookback))
		require.NoError(t, err)

		// close the epoch to apply new omission params
		r.WaitNextEpoch()
		t.Logf("Starting new epoch at block: %d", r.Evm.Context.BlockNumber.Int64())

		r.WaitNBlocks(int(customDelta) * 2)

		// NOTE: blocks numbers mentioned below are relative to the epoch block of the past epoch
		// leave activity proof empty at block 11, proposer is faulty
		autonityFinalize(r)
		require.Equal(t, uint64(1), faultyProposers(r))
		require.True(t, faultyProposer(r, int64(r.lastTargetHeight())))
		t.Logf("proposer faulty for height: %d, target height: %d", r.lastMinedHeight(), r.lastTargetHeight())

		// leave activity proof empty at block 12, proposer is faulty
		autonityFinalize(r)
		require.Equal(t, uint64(2), faultyProposers(r))
		require.True(t, faultyProposer(r, int64(r.lastTargetHeight())))
		t.Logf("proposer faulty for height: %d, target height: %d", r.lastMinedHeight(), r.lastTargetHeight())

		// mine four blocks with non-empty activity proof
		r.WaitNBlocks(4)

		// faulty proposers should still be two
		require.Equal(t, uint64(2), faultyProposers(r))

		// one more block with non-empty activity proof, now faulty proposer should go directly to 0
		r.WaitNBlocks(1)
		require.Equal(t, uint64(0), faultyProposers(r))

	})
	t.Run("random scenario - delta=5, lookback=5, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 5, 5, 4)
	})
	t.Run("random scenario - delta=5, lookback=5, faulty prob=1/2", func(t *testing.T) {
		runRandomScenario(t, 5, 5, 2)
	})
	t.Run("random scenario - delta=5, lookback=5, faulty prob=1", func(t *testing.T) {
		runRandomScenario(t, 5, 5, 1)
	})
	t.Run("random scenario - delta=5, lookback=5, faulty prob=1/10", func(t *testing.T) {
		runRandomScenario(t, 5, 5, 10)
	})
	t.Run("random scenario - delta=5, lookback=5, faulty prob=1/20", func(t *testing.T) {
		runRandomScenario(t, 5, 5, 20)
	})
	t.Run("random scenario - delta=5, lookback=5, faulty prob=1/999999999", func(t *testing.T) {
		runRandomScenario(t, 5, 5, 999999999)
	})
	t.Run("random scenario - delta=2, lookback=5, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 2, 5, 4)
	})
	t.Run("random scenario - delta=5, lookback=1, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 5, 1, 4)
	})
	t.Run("random scenario - delta=2, lookback=1, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 2, 1, 4)
	})
	t.Run("random scenario - delta=45, lookback=5, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 45, 5, 4)
	})
	t.Run("random scenario - delta=5, lookback=55, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 5, 55, 4)
	})
	t.Run("random scenario - delta=23, lookback=31, faulty prob=1/4", func(t *testing.T) {
		runRandomScenario(t, 23, 31, 4)
	})

}

// e.g. nonFaultyProbability 4 --> 1/4 of probability that each height is faulty
func runRandomScenario(t *testing.T, customDelta uint64, customLookback uint64, nonFaultyProbability int) {
	r := Setup(t, configOverride)
	applyCustomParams(r, customDelta, customLookback)

	epochBlock := r.Evm.Context.BlockNumber.Uint64() - 1
	t.Logf("Starting new epoch at block: %d, epoch block: %d", r.Evm.Context.BlockNumber.Int64(), epochBlock)

	t.Logf("Mining %d blocks", customDelta)
	r.WaitNBlocks(int(customDelta))

	faults := make(map[int]struct{})

	countFaults := func(currentHeight uint64) uint64 {
		if currentHeight <= epochBlock+customDelta {
			return 0
		}

		windowEnd := int(currentHeight - customDelta)                      // included
		windowStart := max(int(epochBlock), windowEnd-int(customLookback)) // excluded

		counter := 0
		lastFaulty := -1
		for h := windowEnd; h > windowStart-counter; h-- {
			if h == int(epochBlock) {
				break // finished available blocks
			}
			_, wasFaulty := faults[h]
			if wasFaulty {
				counter++
				lastFaulty = h
			}
		}

		// no faulty proposers in the window
		if lastFaulty == -1 {
			return uint64(counter)
		}

		// if after the last faulty (the oldest) there are lookback - 1 non faulty blocks, we can remove the faulty proposers at the tail of the window
		// as soon as a non-faulty block will come, we will be able to complete the window without needing to extend to block oldest than `lastFaulty`
		nonFaulty := 0
		// exclude window end because in the contract this cleanup is done before we know whether windowEnd is faulty or not
		// here we are instead mirroring it after the call is already executed
		for h := lastFaulty; h < windowEnd; h++ {
			_, wasFaulty := faults[h]
			if !wasFaulty {
				nonFaulty++
			}
		}
		if nonFaulty >= int(customLookback)-1 {
			// exclude window end because in the contract this cleanup is done before we know whether windowEnd is faulty or not
			// here we are instead mirroring it after the call is already executed
			for h := lastFaulty; h < windowEnd; h++ {
				_, wasFaulty := faults[h]
				if !wasFaulty {
					break
				}
				counter--
			}
		}
		return uint64(counter)
	}

	// randomly assign faulty proposers
	for i := int(r.Evm.Context.BlockNumber.Uint64()); i < omissionEpochPeriod*2; i++ {
		t.Logf("Mining block %d, targetHeight %d", i, i-int(customDelta))
		if rand.Intn(nonFaultyProbability) != 0 {
			// all good
			r.WaitNBlocks(1)
			t.Log("Proposer not faulty")
		} else {
			// proposer is faulty
			autonityFinalize(r)
			faults[i-int(customDelta)] = struct{}{}
			t.Log("Proposer IS faulty")
		}
		faultCount := countFaults(uint64(i))
		t.Logf("faulty proposers in window: %d", faultCount)
		require.Equal(t, faultCount, faultyProposers(r))
	}

	// close epoch, faulty proposer number should drop to 0 regardless of fault or not
	if rand.Intn(4) != 0 {
		// all good
		r.WaitNBlocks(1)
	} else {
		// proposer is faulty
		autonityFinalize(r)
	}
	require.Equal(t, uint64(0), faultyProposers(r))
}

func TestLastActive(t *testing.T) {
	r := Setup(t, configOverrideIncreasedStake)

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.WaitNBlocks(int(delta.Int64()))

	proposer := r.Committee.Validators[0].NodeAddress

	// no one should be recorded in the absentees of last height + last active should be -1 for everyone
	require.Equal(t, 0, len(absenteesLastHeight(r)))
	for _, val := range r.Committee.Validators {
		require.Equal(t, int64(-1), lastActive(r, val.NodeAddress))
	}

	// absent should be recorded into the last height absentees and his last active value should be updated
	offline := r.Committee.Validators[1].NodeAddress
	offlineVals := make(map[common.Address]struct{})
	offlineVals[offline] = struct{}{}
	setupProofAndAutonityFinalize(r, proposer, offlineVals)

	absentees := absenteesLastHeight(r)
	t.Log(absentees)
	require.Equal(t, 1, len(absentees))
	require.Equal(t, offline, absentees[0])

	expectedLastActiveHeight := int64(r.lastTargetHeight()) - 1
	t.Logf("Expected last active height: %d", expectedLastActiveHeight)

	for _, val := range r.Committee.Validators {
		if val.NodeAddress != offline {
			require.Equal(t, int64(-1), lastActive(r, val.NodeAddress))
		} else {
			require.Equal(t, expectedLastActiveHeight, lastActive(r, val.NodeAddress))
		}
	}

	// same guy is offline, his lastActive block should remain the same
	// the other offline validators instead should have a more recent lastActive height
	offline2 := r.Committee.Validators[2].NodeAddress
	offlineVals[offline2] = struct{}{}

	setupProofAndAutonityFinalize(r, proposer, offlineVals)
	absentees = absenteesLastHeight(r)
	t.Log(absentees)
	require.Equal(t, 2, len(absentees))
	require.Equal(t, offline, absentees[0])
	require.Equal(t, offline2, absentees[1])
	for _, val := range r.Committee.Validators {
		switch val.NodeAddress {
		case offline:
			require.Equal(t, expectedLastActiveHeight, lastActive(r, val.NodeAddress))
		case offline2:
			require.Equal(t, expectedLastActiveHeight+1, lastActive(r, val.NodeAddress))
		default:
			require.Equal(t, int64(-1), lastActive(r, val.NodeAddress))
		}
	}

	// everything should be cleared like in the initial state
	r.WaitNBlocks(1)
	require.Equal(t, 0, len(absenteesLastHeight(r)))
	for _, val := range r.Committee.Validators {
		require.Equal(t, int64(-1), lastActive(r, val.NodeAddress))
	}

	// reach last epoch block - 1
	currentBlock := r.Evm.Context.BlockNumber.Uint64() - 1
	r.WaitNBlocks(omissionEpochPeriod - int(currentBlock) - 2)

	// have some absents, their last active block should change
	expectedLastActiveHeight = int64(r.lastTargetHeight())
	setupProofAndAutonityFinalize(r, proposer, offlineVals)
	for _, val := range r.Committee.Validators {
		switch val.NodeAddress {
		case offline:
			fallthrough
		case offline2:
			require.Equal(t, expectedLastActiveHeight, lastActive(r, val.NodeAddress))
		default:
			require.Equal(t, int64(-1), lastActive(r, val.NodeAddress))
		}
	}

	// same validators still offline, however their lastActive counter should get reset by the epoch end (`setCommittee`)
	setupProofAndAutonityFinalize(r, proposer, offlineVals)
	for _, val := range r.Committee.Validators {
		require.Equal(t, int64(-1), lastActive(r, val.NodeAddress))
	}
}
