package tests

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"testing"
)

const BigFloatPrecision = 256 // will allow full representation of the solidity uint256 range

func newFloat(value *big.Int) *big.Float {
	return new(big.Float).SetPrec(BigFloatPrecision).SetInt(value)
}

func toString(value *big.Float) string {
	return value.Text('f', -1)
}

var omissionEpochPeriod = 130

const SlashingRatePrecision = 10_000 // must match the precision used in Slasher.sol

const active = uint8(0)
const jailed = uint8(2)
const jailedForInactivity = uint8(4)

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

// helpers
func omissionFinalize(r *runner, epochEnd bool) {
	_, err := r.omissionAccountability.Finalize(fromAutonity, epochEnd)
	require.NoError(r.t, err)
	r.t.Logf("Omission accountability, finalized block: %d", r.evm.Context.BlockNumber)
	// advance the block context as if we mined a block
	r.evm.Context.BlockNumber = new(big.Int).Add(r.evm.Context.BlockNumber, common.Big1)
	r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
	// clean up activity proof data
	r.evm.Context.Coinbase = common.Address{}
	r.evm.Context.ActivityProof = nil
	r.evm.Context.ActivityProofRound = 0
}

func autonityFinalize(r *runner) { //nolint
	_, err := r.autonity.Finalize(nil)
	require.NoError(r.t, err)
	r.t.Logf("Autonity, finalized block: %d", r.evm.Context.BlockNumber)
	// advance the block context as if we mined a block
	r.evm.Context.BlockNumber = new(big.Int).Add(r.evm.Context.BlockNumber, common.Big1)
	r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
	// clean up activity proof data
	r.evm.Context.Coinbase = common.Address{}
	r.evm.Context.ActivityProof = nil
	r.evm.Context.ActivityProofRound = 0
}

func setupProofAndAutonityFinalize(r *runner, proposer common.Address, absentees map[common.Address]struct{}) {
	r.setupActivityProofAndCoinbase(proposer, absentees)
	autonityFinalize(r)
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

func offences(r *runner, validator common.Address) int {
	offences, _, err := r.omissionAccountability.RepeatedOffences(nil, validator)
	require.NoError(r.t, err)
	return int(offences.Uint64())
}

func inactivityScore(r *runner, validator common.Address) int {
	score, _, err := r.omissionAccountability.InactivityScores(nil, validator)
	require.NoError(r.t, err)
	return int(score.Uint64())
}

func omissionScaleFactor(r *runner) *big.Int {
	factor, _, err := r.omissionAccountability.GetScaleFactor(nil)
	require.NoError(r.t, err)
	return factor
}

func proposerEffort(r *runner, validator common.Address) *big.Int {
	effort, _, err := r.omissionAccountability.ProposerEffort(nil, validator)
	require.NoError(r.t, err)
	return effort
}

func totalProposerEffort(r *runner) *big.Int {
	effort, _, err := r.omissionAccountability.TotalEffort(nil)
	require.NoError(r.t, err)
	return effort
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

func validator(r *runner, addr common.Address) AutonityValidator {
	val, _, err := r.autonity.GetValidator(nil, addr)
	require.NoError(r.t, err)
	return val
}

func ntnBalance(r *runner, addr common.Address) *big.Int {
	balance, _, err := r.autonity.BalanceOf(nil, addr)
	require.NoError(r.t, err)
	return balance
}

func TestAccessControl(t *testing.T) {
	r := setup(t, nil)

	_, err := r.omissionAccountability.Finalize(r.operator, false)
	require.Error(r.t, err)
	_, err = r.omissionAccountability.Finalize(fromAutonity, false)
	require.NoError(r.t, err)

	_, err = r.omissionAccountability.DistributeProposerRewards(r.operator, common.Big256)
	require.Error(r.t, err)
	_, err = r.omissionAccountability.DistributeProposerRewards(fromAutonity, common.Big256)
	require.NoError(r.t, err)

	_, err = r.omissionAccountability.SetCommittee(r.operator, []AutonityCommitteeMember{}, []common.Address{})
	require.Error(r.t, err)
	_, err = r.omissionAccountability.SetCommittee(fromAutonity, []AutonityCommitteeMember{}, []common.Address{})
	require.NoError(r.t, err)

	_, err = r.omissionAccountability.SetEpochBlock(r.operator, common.Big256)
	require.Error(r.t, err)
	_, err = r.omissionAccountability.SetEpochBlock(fromAutonity, common.Big256)
	require.NoError(r.t, err)
}

func TestProposerLogic(t *testing.T) {
	t.Run("Faulty proposer inactive score increases and height is marked as invalid", func(t *testing.T) {
		r := setup(t, configOverride)

		delta, _, err := r.omissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.waitNBlocks(int(delta.Int64()))

		targetHeight := r.evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer := r.committee.validators[0].NodeAddress
		r.evm.Context.Coinbase = proposer
		r.evm.Context.ActivityProof = nil
		autonityFinalize(r)

		require.True(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 1, inactivityCounter(r, proposer))

		r.evm.Context.Coinbase = proposer
		r.evm.Context.ActivityProof = nil
		autonityFinalize(r)
		require.True(r.t, faultyProposer(r, targetHeight+1))
		require.Equal(r.t, 2, inactivityCounter(r, proposer))

		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.t, faultyProposer(r, targetHeight+2))
		require.Equal(r.t, 2, inactivityCounter(r, proposer))
	})
	t.Run("Proposer effort is correctly computed", func(t *testing.T) {
		r := setup(t, configOverride)

		delta, _, err := r.omissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.waitNBlocks(int(delta.Int64()))

		totalVotingPower := new(big.Int)
		for _, val := range r.committee.validators {
			totalVotingPower.Add(totalVotingPower, val.BondedStake)
		}
		quorum := bft.Quorum(totalVotingPower)
		fullProofEffort := new(big.Int).Sub(totalVotingPower, quorum) // proposer effort when a full activity proof is provided

		targetHeight := r.evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer := r.committee.validators[0].NodeAddress
		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, fullProofEffort.String(), proposerEffort(r, proposer).String())
		require.Equal(r.t, fullProofEffort.String(), totalProposerEffort(r).String())

		// finalize 3 more times with full proof
		targetHeight = r.evm.Context.BlockNumber.Int64() - delta.Int64()
		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.t, faultyProposer(r, targetHeight))
		targetHeight++
		setupProofAndAutonityFinalize(r, proposer, nil)
		require.False(r.t, faultyProposer(r, targetHeight))
		targetHeight++
		setupProofAndAutonityFinalize(r, proposer, nil)
		expectedEffort := new(big.Int).Mul(fullProofEffort, common.Big4) // we finalized 4 times up to now
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, expectedEffort.String(), proposerEffort(r, proposer).String())
		require.Equal(r.t, expectedEffort.String(), totalProposerEffort(r).String())

		targetHeight = r.evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer = r.committee.validators[1].NodeAddress
		setupProofAndAutonityFinalize(r, proposer, nil)
		expectedTotalEffort := new(big.Int).Add(expectedEffort, fullProofEffort) // validators[0] effort + validator[1] effort
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, fullProofEffort.String(), proposerEffort(r, proposer).String())
		require.Equal(r.t, expectedTotalEffort.String(), totalProposerEffort(r).String())
	})
}

// checks that the inactivity counters are correctly updated according to the lookback window (_recordAbsentees function)
func TestInactivityCounter(t *testing.T) {
	r := setup(t, configOverrideIncreasedStake)

	// set maximum inactivity threshold for this test, we care only about the inactivity counters and not about the jailing
	_, err := r.omissionAccountability.SetInactivityThreshold(r.operator, new(big.Int).SetUint64(10000))
	require.NoError(t, err)

	delta, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	lookback := int(config.LookbackWindow.Uint64())

	proposer := r.committee.validators[0].NodeAddress
	fullyOffline := r.committee.validators[1].NodeAddress
	partiallyOffline := r.committee.validators[2].NodeAddress

	absentees := make(map[common.Address]struct{})
	absentees[fullyOffline] = struct{}{}
	partiallyOfflineCounter := 0
	for i := 0; i < lookback-2; i++ {
		if i == lookback/2 {
			absentees[partiallyOffline] = struct{}{}
		}
		targetHeight := r.evm.Context.BlockNumber.Int64() - delta.Int64()
		setupProofAndAutonityFinalize(r, proposer, absentees)
		for absentee := range absentees {
			require.True(r.t, isValidatorInactive(r, targetHeight, absentee))
			if absentee == partiallyOffline {
				partiallyOfflineCounter++
			}
		}
	}

	// we still need two height to have a full lookback window
	// insert a proposer faulty height (no activity proof), it should be ignored
	r.t.Logf("current block number in evm: %d", r.evm.Context.BlockNumber.Uint64())
	r.evm.Context.Coinbase = proposer
	r.evm.Context.ActivityProof = nil
	autonityFinalize(r)
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	// here we should update the inactivity counter to 1, but since there was a faulty proposer we extend the lookback period
	r.t.Logf("current block number in evm: %d", r.evm.Context.BlockNumber.Uint64())
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	// now we have a full lookback period
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.t, 1, inactivityCounter(r, fullyOffline))
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.t, 2, inactivityCounter(r, fullyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	partiallyOfflineCounter++
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.t, 3, inactivityCounter(r, fullyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))

	// fill up enough blocks for partiallyOffline as well
	for i := partiallyOfflineCounter; i < lookback-1; i++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}

	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	setupProofAndAutonityFinalize(r, proposer, absentees)
	require.Equal(r.t, 1, inactivityCounter(r, partiallyOffline))

	fullyOfflineIC := inactivityCounter(r, fullyOffline)
	partiallyOfflineIC := inactivityCounter(r, partiallyOffline)
	// every two block, one has faulty proposer
	n := 20
	for i := 0; i < (n * 2); i++ {
		proposerFaulty := i%2 == 0
		if !proposerFaulty {
			r.setupActivityProofAndCoinbase(proposer, absentees)
		} else {
			r.evm.Context.Coinbase = proposer
			r.evm.Context.ActivityProof = nil
		}
		autonityFinalize(r)
	}

	// inactivity counter should still have increased by n due to lookback period extension
	require.Equal(r.t, fullyOfflineIC+n, inactivityCounter(r, fullyOffline))
	require.Equal(r.t, partiallyOfflineIC+n, inactivityCounter(r, partiallyOffline))

	//  close the epoch
	for i := r.evm.Context.BlockNumber.Int64(); i < int64(omissionEpochPeriod); i++ {
		setupProofAndAutonityFinalize(r, proposer, nil)
	}
	t.Log("Closing epoch")
	setupProofAndAutonityFinalize(r, proposer, nil)
	r.generateNewCommittee()

	// inactivity counters should be reset
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	r.waitNBlocks(int(delta.Int64()))
	t.Logf("current consensus instance for height %d", r.evm.Context.BlockNumber.Uint64())
	otherValidator := r.committee.validators[3].NodeAddress
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
		require.Equal(r.t, 0, inactivityCounter(r, otherValidator))
	}

	// proposer faulty
	r.evm.Context.Coinbase = proposer
	r.evm.Context.ActivityProof = nil
	autonityFinalize(r)
	require.Equal(r.t, 0, inactivityCounter(r, otherValidator))

	setupProofAndAutonityFinalize(r, proposer, newAbsentees)
	require.Equal(r.t, 1, inactivityCounter(r, otherValidator))
}

// checks that the inactivity scores are computed correctly
func TestInactivityScore(t *testing.T) {
	r := setup(t, configOverrideIncreasedStake)

	// set maximum inactivity threshold for this test, we care only about the inactivity scores and not about the jailing
	scaleFactorInt := omissionScaleFactor(r)
	_, err := r.omissionAccountability.SetInactivityThreshold(r.operator, scaleFactorInt)
	require.NoError(t, err)

	delta, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	scaleFactor := newFloat(scaleFactorInt)

	initialCommitteeSize := len(r.committee.validators)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	lookback := int(config.LookbackWindow.Uint64())
	pastPerformanceWeight := new(big.Float).Quo(newFloat(config.PastPerformanceWeight), scaleFactor)
	currentPerformanceWeight := new(big.Float).Sub(newFloat(common.Big1), pastPerformanceWeight)

	// simulate epoch.
	proposer := r.committee.validators[0].NodeAddress
	inactiveBlockStreak := make([]int, len(r.committee.validators))
	inactiveCounters := make([]*big.Int, len(r.committee.validators))
	for i := range inactiveCounters {
		inactiveCounters[i] = new(big.Int)
	}
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
		absentees := make(map[common.Address]struct{})
		for i := range r.committee.validators {
			if r.committee.validators[i].NodeAddress == proposer {
				continue // keep proposer always online
			}
			if rand.Intn(30) != 0 {
				absentees[r.committee.validators[i].NodeAddress] = struct{}{}
				inactiveBlockStreak[i]++
			} else {
				inactiveBlockStreak[i] = 0
			}
			if inactiveBlockStreak[i] >= lookback {
				inactiveCounters[i] = new(big.Int).Add(inactiveCounters[i], common.Big1)
			}
		}

		t.Logf("number of absentees: %d for height %d", len(absentees), r.evm.Context.BlockNumber.Uint64())
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	// no one should have been jailed
	require.Equal(t, initialCommitteeSize, len(r.committee.validators))

	// check score computation
	pastInactivityScore := make([]*big.Float, len(r.committee.validators))
	denominator := new(big.Int).SetUint64(uint64(omissionEpochPeriod) - delta.Uint64() - uint64(lookback) + 1)
	for i, val := range r.committee.validators {
		score := new(big.Float).Quo(newFloat(inactiveCounters[i]), newFloat(denominator))

		expectedInactivityScoreFloat := new(big.Float).Mul(score, currentPerformanceWeight) // + 0 * pastPerformanceWeight
		pastInactivityScore[i] = expectedInactivityScoreFloat

		expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, scaleFactor)
		expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Int(nil)

		r.t.Logf("i %d, expectedInactivityScoreFloat %s, expectedInactivityScore %s, inactivityScore %v", i, toString(expectedInactivityScoreFloat), expectedInactivityScore.String(), inactivityScore(r, val.NodeAddress))

		require.Equal(r.t, int(expectedInactivityScore.Uint64()), inactivityScore(r, val.NodeAddress))
	}

	// simulate another epoch
	r.waitNBlocks(int(delta.Int64()))
	inactiveBlockStreak = make([]int, len(r.committee.validators))
	inactiveCounters = make([]*big.Int, len(r.committee.validators))
	for i := range inactiveCounters {
		inactiveCounters[i] = new(big.Int)
	}
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
		absentees := make(map[common.Address]struct{})
		for i := range r.committee.validators {
			if r.committee.validators[i].NodeAddress == proposer {
				continue // keep proposer always online
			}
			if rand.Intn(30) != 0 {
				absentees[r.committee.validators[i].NodeAddress] = struct{}{}
				inactiveBlockStreak[i]++
			} else {
				inactiveBlockStreak[i] = 0
			}
			if inactiveBlockStreak[i] >= lookback {
				inactiveCounters[i] = new(big.Int).Add(inactiveCounters[i], common.Big1)
			}
		}

		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	// no one should have been jailed
	require.Equal(t, initialCommitteeSize, len(r.committee.validators))

	// check score computation
	for i, val := range r.committee.validators {
		score := new(big.Float).Quo(newFloat(inactiveCounters[i]), newFloat(denominator))

		expectedInactivityScoreFloat1 := new(big.Float).Mul(score, currentPerformanceWeight)
		expectedInactivityScoreFloat2 := new(big.Float).Mul(pastInactivityScore[i], pastPerformanceWeight)
		expectedInactivityScoreFloat := new(big.Float).Add(expectedInactivityScoreFloat1, expectedInactivityScoreFloat2)

		expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, scaleFactor)
		expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Int(nil)

		r.t.Logf("i %d, expectedInactivityScoreFloat %s, expectedInactivityScore %s, inactivityScore %v", i, toString(expectedInactivityScoreFloat), expectedInactivityScore.String(), inactivityScore(r, val.NodeAddress))
		require.Equal(r.t, int(expectedInactivityScore.Uint64()), inactivityScore(r, val.NodeAddress))
	}
}

func TestOmissionPunishments(t *testing.T) {
	r := setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		config.EpochPeriod = uint64(omissionEpochPeriod)
		// increase voting power of validator 0 to reach quorum in proofs easily
		config.Validators[0].BondedStake = new(big.Int).Mul(config.Validators[1].BondedStake, big.NewInt(6))
		return config
	})

	delta, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	initialJailingPeriod := int(config.InitialJailingPeriod.Uint64())
	initialProbationPeriod := int(config.InitialProbationPeriod.Uint64())
	pastPerformanceWeight := int(config.PastPerformanceWeight.Uint64())
	initialSlashingRate := int(config.InitialSlashingRate.Uint64())
	scaleFactor := int(omissionScaleFactor(r).Uint64())

	proposer := r.committee.validators[0].NodeAddress
	absentees := make(map[common.Address]struct{})
	val1Address := r.committee.validators[1].NodeAddress // will be handy to have those later
	val2Address := r.committee.validators[2].NodeAddress
	absentees[val1Address] = struct{}{}
	absentees[val2Address] = struct{}{}
	val1Treasury := r.committee.validators[1].Treasury
	val2Treasury := r.committee.validators[2].Treasury

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
		require.Equal(r.t, expectedFullOfflineScore, inactivityScore(r, absentee))
		val := validator(r, absentee)
		require.Equal(r.t, jailedForInactivity, val.State)
		require.Equal(r.t, uint64(omissionEpochPeriod+initialJailingPeriod), val.JailReleaseBlock.Uint64())
		require.Equal(r.t, initialProbationPeriod, probation(r, absentee))
		require.Equal(t, 1, offences(r, absentee))
	}

	// wait that the jailing finishes and reactivate validators
	r.waitNBlocks(initialJailingPeriod)
	_, err = r.autonity.ActivateValidator(&runOptions{origin: val1Treasury}, val1Address)
	require.NoError(r.t, err)
	_, err = r.autonity.ActivateValidator(&runOptions{origin: val2Treasury}, val2Address)
	require.NoError(r.t, err)

	// pass some epochs, probation period should decrease
	r.waitNextEpoch() // re-activation epoch, val not part of committee
	// inactivity score should still be the same as before
	for absentee := range absentees {
		require.Equal(r.t, expectedFullOfflineScore, inactivityScore(r, absentee))
	}
	r.waitNextEpoch()
	// should be decreased  now
	for absentee := range absentees {
		require.Equal(r.t, (expectedFullOfflineScore*pastPerformanceWeight)/scaleFactor, inactivityScore(r, absentee))
	}
	r.waitNextEpoch()
	r.waitNextEpoch()

	// probation periods should have decreased of once for every epoch that passed with the validator as part of the committee
	passedEpochs := 3

	for absentee := range absentees {
		require.Equal(t, initialProbationPeriod-passedEpochs, probation(r, absentee))
	}

	// simulate another epoch where:
	// - val 1 gets slashed by accountability and therefore doesn't get punished by omission accountability
	// - val 2 gets punished again for omission while in the probation period, therefore he gets slashed
	r.waitNBlocks(int(delta.Int64()))
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	val1 := validator(r, val1Address)
	val1.State = jailed
	val1.JailReleaseBlock = new(big.Int).SetInt64(r.evm.Context.BlockNumber.Int64() + int64(omissionEpochPeriod*10))
	totalSlashedVal1 := val1.TotalSlashed
	_, err = r.autonity.UpdateValidatorAndTransferSlashedFunds(&runOptions{origin: r.accountability.address}, val1)
	require.NoError(t, err)

	val2BeforeSlash := validator(r, val2Address)
	// close epoch
	setupProofAndAutonityFinalize(r, proposer, absentees)
	r.generateNewCommittee()

	// val1, punished by accountability, shouldn't have been slashed by omission even if 100% offline and still under probation
	val1 = validator(r, val1Address)
	require.Equal(r.t, jailed, val1.State)
	require.True(r.t, probation(r, val1.NodeAddress) > 0)
	require.Equal(r.t, totalSlashedVal1.String(), val1.TotalSlashed.String())
	require.Equal(r.t, 1, offences(r, val1Address))

	// val2 offline while on probation, should have been slashed by omission
	val2 := validator(r, val2Address)
	require.Equal(r.t, jailedForInactivity, val2.State)
	require.True(r.t, probation(r, val2.NodeAddress) > 0)
	require.True(r.t, val2.TotalSlashed.Cmp(val2BeforeSlash.TotalSlashed) > 0)
	require.Equal(r.t, 2, offences(r, val2Address))
	expectedSlashRate := new(big.Int).SetInt64(int64(initialSlashingRate * 4 * 2)) // rate * offence^2 * collusion
	availableFunds := new(big.Int).Add(val2BeforeSlash.BondedStake, val2.UnbondingStake)
	availableFunds.Add(availableFunds, val2.SelfUnbondingStake)
	expectedSlashAmount := new(big.Int).Mul(expectedSlashRate, availableFunds)
	expectedSlashAmount.Div(expectedSlashAmount, new(big.Int).SetInt64(SlashingRatePrecision))
	t.Logf("expected slash rate: %s, available funds: %s, expected slash amount: %s", expectedSlashRate.String(), availableFunds.String(), expectedSlashAmount.String())
	require.Equal(r.t, expectedSlashAmount.String(), new(big.Int).Sub(val2.TotalSlashed, val2BeforeSlash.TotalSlashed).String())
}

func TestProposerRewardDistribution(t *testing.T) {
	t.Run("Rewards are correctly allocated based on config", func(t *testing.T) {
		r := setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			config.EpochPeriod = uint64(omissionEpochPeriod)
			config.TreasuryFee = 0
			return config
		})

		maxCommitteeSizeBig, _, err := r.autonity.GetMaxCommitteeSize(nil)
		require.NoError(r.t, err)
		maxCommitteeSize := newFloat(maxCommitteeSizeBig)
		t.Logf("max committee size: %s", toString(maxCommitteeSize))

		config, _, err := r.autonity.Config(nil)
		require.NoError(r.t, err)
		proposerRewardRateBig := config.Policy.ProposerRewardRate
		proposerRewardRate := newFloat(proposerRewardRateBig)
		proposerRewardRatePrecisionBig, _, err := r.autonity.PROPOSERREWARDRATEPRECISION(nil)
		require.NoError(t, err)
		proposerRewardRatePrecision := newFloat(proposerRewardRatePrecisionBig)
		t.Logf("proposer reward rate: %s, precision: %s", toString(proposerRewardRate), toString(proposerRewardRatePrecision))

		autonityAtnsBig := new(big.Int).SetUint64(54644455456467) // random amount
		t.Logf("atn rewards: %s", autonityAtnsBig.String())
		// this has to match the ntn inflation unlocked NTNs.
		// Can be retrieved by adding in solidity a revert(Helpers.toString(accounts[address(this)])); in Finalize
		ntnRewardsBig := new(big.Int).SetUint64(8205384319979600000)
		t.Logf("ntn rewards: %s", ntnRewardsBig.String())
		r.giveMeSomeMoney(r.autonity.address, autonityAtnsBig)

		autonityAtns := newFloat(autonityAtnsBig)
		ntnRewards := newFloat(ntnRewardsBig)

		// all rewards should go to val 0
		proposer := r.committee.validators[0].NodeAddress
		proposerTreasury := r.committee.validators[0].Treasury
		atnBalanceBefore := newFloat(r.getBalanceOf(proposerTreasury))
		ntnBalanceBefore := newFloat(ntnBalance(r, proposerTreasury))
		t.Logf("atn balance before: %s, ntn balance before %s", toString(atnBalanceBefore), toString(ntnBalanceBefore))

		// set validator state to jailed so that he will not receive any reward other the proposer one
		val := validator(r, proposer)
		val.State = jailed
		_, err = r.autonity.UpdateValidatorAndTransferSlashedFunds(&runOptions{origin: r.accountability.address}, val)
		require.NoError(t, err)

		r.evm.Context.BlockNumber = new(big.Int).SetInt64(int64(omissionEpochPeriod))
		r.evm.Context.Time.Add(r.evm.Context.Time, new(big.Int).SetInt64(int64(omissionEpochPeriod-1)))
		setupProofAndAutonityFinalize(r, proposer, nil)

		committeeSize := newFloat(new(big.Int).SetUint64(uint64(len(r.committee.validators))))
		t.Logf("proposer reward rate: %s, committee size: %s", toString(proposerRewardRate), toString(committeeSize))
		numeratorFactor := new(big.Float).Mul(proposerRewardRate, committeeSize)
		t.Logf("numeratorFactor: %s", toString(numeratorFactor))

		t.Logf("proposer reward rate precision: %s, max committee size: %s", toString(proposerRewardRatePrecision), toString(maxCommitteeSize))
		denominator := new(big.Float).Mul(proposerRewardRatePrecision, maxCommitteeSize)
		t.Logf("denominator: %s", toString(denominator))

		numeratorAtn := new(big.Float).Mul(autonityAtns, numeratorFactor)
		t.Logf("numerator atn: %s", toString(numeratorAtn))
		atnExpectedReward := new(big.Float).Quo(numeratorAtn, denominator)

		numeratorNtn := new(big.Float).Mul(ntnRewards, numeratorFactor)
		t.Logf("numerator ntn: %s", toString(numeratorNtn))
		ntnExpectedReward := new(big.Float).Quo(numeratorNtn, denominator)

		t.Logf("atn expected reward: %s, ntn expected reward: %s", toString(atnExpectedReward), toString(ntnExpectedReward))

		atnExpectedBalance := new(big.Float).Add(atnBalanceBefore, atnExpectedReward)
		ntnExpectedBalance := new(big.Float).Add(ntnBalanceBefore, ntnExpectedReward)

		t.Logf("atn expected balance: %s, ntn expected balance: %s", toString(atnExpectedBalance), toString(ntnExpectedBalance))

		atnActualBalance := r.getBalanceOf(proposerTreasury)
		ntnActualBalance := ntnBalance(r, proposerTreasury)
		t.Logf("atn actual balance: %s, ntn actual balance: %s", atnActualBalance.String(), ntnActualBalance.String())

		atnExpectedBalanceInt, _ := atnExpectedBalance.Int(nil)
		ntnExpectedBalanceInt, _ := ntnExpectedBalance.Int(nil)

		t.Logf("balance diff atn: %s", new(big.Int).Sub(atnExpectedBalanceInt, atnActualBalance).String())
		t.Logf("balance diff ntn: %s", new(big.Int).Sub(ntnExpectedBalanceInt, ntnActualBalance).String())

		require.Equal(t, atnExpectedBalanceInt.String(), atnActualBalance.String())
		require.Equal(t, ntnExpectedBalanceInt.String(), ntnActualBalance.String())
	})
	t.Run("Rewards are correctly distributed among proposers", func(t *testing.T) {
		r := setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			config.EpochPeriod = uint64(omissionEpochPeriod)
			return config
		})

		delta, _, err := r.omissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.waitNBlocks(int(delta.Int64()))

		totalEffort := new(big.Int)
		efforts := make([]*big.Int, len(r.committee.validators))
		atnBalances := make([]*big.Int, len(r.committee.validators))
		ntnBalances := make([]*big.Int, len(r.committee.validators))
		totalPower := new(big.Int)
		for i, val := range r.committee.validators {
			efforts[i] = new(big.Int)
			atnBalances[i] = r.getBalanceOf(val.Treasury)
			ntnBalances[i] = ntnBalance(r, val.Treasury)
			totalPower.Add(totalPower, val.BondedStake)
		}
		// simulate epoch
		fullProofEffort := new(big.Int).Sub(totalPower, bft.Quorum(totalPower)) // effort for a full proof
		for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
			proposerIndex := rand.Intn(len(r.committee.validators))
			totalEffort.Add(totalEffort, fullProofEffort)
			efforts[proposerIndex].Add(efforts[proposerIndex], fullProofEffort)
			r.setupActivityProofAndCoinbase(r.committee.validators[proposerIndex].NodeAddress, nil)
			omissionFinalize(r, h == omissionEpochPeriod)
		}

		simulatedNtnRewards := new(big.Int).SetInt64(5968565)
		simulatedAtnRewards := new(big.Int).SetInt64(4545445)
		r.giveMeSomeMoney(r.autonity.address, simulatedAtnRewards)
		_, err = r.autonity.Mint(r.operator, r.omissionAccountability.address, simulatedNtnRewards)
		require.NoError(r.t, err)
		_, err = r.omissionAccountability.DistributeProposerRewards(&runOptions{origin: r.autonity.address, value: simulatedAtnRewards}, simulatedNtnRewards)
		require.NoError(t, err)

		for i, val := range r.committee.validators {
			atnExpectedIncrement := new(big.Int).Mul(efforts[i], simulatedAtnRewards)
			atnExpectedIncrement.Div(atnExpectedIncrement, totalEffort)
			ntnExpectedIncrement := new(big.Int).Mul(efforts[i], simulatedNtnRewards)
			ntnExpectedIncrement.Div(ntnExpectedIncrement, totalEffort)
			atnExpectedBalance := new(big.Int).Add(atnBalances[i], atnExpectedIncrement)
			ntnExpectedBalance := new(big.Int).Add(ntnBalances[i], ntnExpectedIncrement)
			t.Logf("validator %d, effort %s, total effort %s, expectedBalance atn %s, expectedBalanceNtn %s", i, efforts[i].String(), totalEffort.String(), atnExpectedBalance.String(), ntnExpectedBalance.String())

			atnBalance := r.getBalanceOf(val.Treasury)
			ntnBalance := ntnBalance(r, val.Treasury)

			require.Equal(t, atnExpectedBalance.String(), atnBalance.String())
			require.Equal(t, ntnExpectedBalance.String(), ntnBalance.String())

			// effort counters should be zeroed out
			require.Equal(r.t, common.Big0.String(), proposerEffort(r, val.NodeAddress).String())
		}

		require.Equal(r.t, common.Big0.String(), totalProposerEffort(r).String())
	})
}

// past performance weight and inactivity threshold should be set low enough that if:
// - a validator gets 100% inactivity in epoch x
// - then he gets 0% inactivity in epoch x+n (after he reactivated)
// he shouldn't get slashed in epoch x+n
func TestConfigSanity(t *testing.T) {
	r := setup(t, configOverrideIncreasedStake)

	delta, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	initialJailingPeriod := int(config.InitialJailingPeriod.Uint64())

	proposer := r.committee.validators[0].NodeAddress
	absentees := make(map[common.Address]struct{})
	val1Address := r.committee.validators[1].NodeAddress // will be handy later
	val2Address := r.committee.validators[2].NodeAddress // will be handy later
	absentees[val1Address] = struct{}{}
	absentees[val2Address] = struct{}{}
	val1Treasury := r.committee.validators[1].Treasury
	val2Treasury := r.committee.validators[2].Treasury

	// simulate epoch with two validator at 100% inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	for absentee := range absentees {
		val := validator(r, absentee)
		require.Equal(r.t, jailedForInactivity, val.State)
		require.Equal(t, 1, offences(r, absentee))
	}

	// wait that the jailing finishes and reactivate validators
	r.waitNBlocks(initialJailingPeriod)
	_, err = r.autonity.ActivateValidator(&runOptions{origin: val1Treasury}, val1Address)
	require.NoError(r.t, err)
	_, err = r.autonity.ActivateValidator(&runOptions{origin: val2Treasury}, val2Address)
	require.NoError(r.t, err)

	r.waitNextEpoch() // re-activation epoch, val not part of committee
	r.waitNextEpoch()

	// validator should not have been punished since he did 0% offline
	for absentee := range absentees {
		val := validator(r, absentee)
		require.Equal(r.t, active, val.State)
		require.Equal(t, 1, offences(r, absentee))
	}

}

func TestRewardWithholding(t *testing.T) {
	r := setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		config.EpochPeriod = uint64(omissionEpochPeriod)
		config.ProposerRewardRate = 0 // no rewards to proposers to make computation simpler
		config.TreasuryFee = 0        // same
		// increase voting power of validator 0 to reach quorum in proofs easily
		config.Validators[0].BondedStake = new(big.Int).Mul(config.Validators[1].BondedStake, big.NewInt(6))
		return config
	})

	delta, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	// validators over threshold will get all their rewards withheld
	customInactivityThreshold := uint64(6000)
	_, err = r.omissionAccountability.SetInactivityThreshold(r.operator, new(big.Int).SetUint64(customInactivityThreshold))
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.autonity.Config(nil)
	require.NoError(t, err)
	withheldRewardPool := config.Policy.WithheldRewardsPool

	proposer := r.committee.validators[0].NodeAddress

	// simulate epoch with random levels of inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		absentees := make(map[common.Address]struct{})
		for i := range r.committee.validators {
			if i == 0 {
				continue // let's keep at least a guy inside the committee
			}
			if rand.Intn(30) != 0 {
				absentees[r.committee.validators[i].NodeAddress] = struct{}{}
			}
		}
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}

	atnRewards := new(big.Int).SetUint64(5467879877987) // random amount
	// this has to match the ntn inflation unlocked NTNs.
	// Can be retrieved by adding in solidity a revert(Helpers.toString(accounts[address(this)])); in Finalize
	ntnRewards := new(big.Int).SetUint64(8220842843566600000)
	r.giveMeSomeMoney(r.autonity.address, atnRewards)

	atnBalancesBefore := make([]*big.Int, len(r.committee.validators))
	ntnBalancesBefore := make([]*big.Int, len(r.committee.validators))
	totalPower := new(big.Int)
	for i, val := range r.committee.validators {
		validatorStruct := validator(r, val.NodeAddress)
		// we assume that all stake is self bonded in this test
		require.Equal(t, validatorStruct.SelfBondedStake.String(), validatorStruct.BondedStake.String())
		atnBalancesBefore[i] = r.getBalanceOf(val.Treasury)
		ntnBalancesBefore[i] = ntnBalance(r, val.Treasury)
		t.Logf("validator %d, atn balance before: %s, ntn balance before %s", i, atnBalancesBefore[i].String(), ntnBalancesBefore[i].String())
		totalPower.Add(totalPower, validatorStruct.SelfBondedStake)
	}
	atnPoolBefore := r.getBalanceOf(withheldRewardPool)
	ntnPoolBefore := ntnBalance(r, withheldRewardPool)

	setupProofAndAutonityFinalize(r, proposer, nil)

	atnTotalWithheld := new(big.Int)
	ntnTotalWithheld := new(big.Int)
	for i, val := range r.committee.validators {
		validatorStruct := validator(r, val.NodeAddress)
		power := validatorStruct.SelfBondedStake

		// compute reward without withholding
		atnFullReward := new(big.Int).Mul(power, atnRewards)
		atnFullReward.Div(atnFullReward, totalPower)
		ntnFullReward := new(big.Int).Mul(power, ntnRewards)
		ntnFullReward.Div(ntnFullReward, totalPower)

		// compute withheld amount
		score := new(big.Int).SetInt64(int64(inactivityScore(r, val.NodeAddress)))
		var ntnWithheld *big.Int
		var atnWithheld *big.Int
		t.Logf("validator index %d, score: %d", i, score.Uint64())
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
		atnTotalWithheld.Add(atnTotalWithheld, atnWithheld)
		ntnTotalWithheld.Add(ntnTotalWithheld, ntnWithheld)

		// check validator balance
		atnExpectedBalance := new(big.Int).Add(atnFullReward, atnBalancesBefore[i])
		atnExpectedBalance.Sub(atnExpectedBalance, atnWithheld)
		ntnExpectedBalance := new(big.Int).Add(ntnFullReward, ntnBalancesBefore[i])
		ntnExpectedBalance.Sub(ntnExpectedBalance, ntnWithheld)
		require.Equal(t, atnExpectedBalance.String(), r.getBalanceOf(val.Treasury).String())
		require.Equal(t, ntnExpectedBalance.String(), ntnBalance(r, val.Treasury).String())
	}
	atnExpectedPoolBalance := atnPoolBefore.Add(atnPoolBefore, atnTotalWithheld)
	ntnExpectedPoolBalance := ntnPoolBefore.Add(ntnPoolBefore, ntnTotalWithheld)
	require.Equal(t, atnExpectedPoolBalance.String(), r.getBalanceOf(withheldRewardPool).String())
	require.Equal(t, ntnExpectedPoolBalance.String(), ntnBalance(r, withheldRewardPool).String())
}

// operator can disable omission punishments by increasing the inactivity threshold
// there are also other ways to disable accountability, by lowering the initial punishments values
// proposer rewards can be disabled by lowering the proposer reward rate in the autonity contract
func TestOmissionDisabling(t *testing.T) {
	r := setup(t, configOverride)

	_, err := r.omissionAccountability.SetInactivityThreshold(r.operator, omissionScaleFactor(r))
	require.NoError(t, err)

	// validator 1 absent for entire epoch
	absentees := make(map[common.Address]struct{})
	absentees[r.committee.validators[1].NodeAddress] = struct{}{}

	csize := len(r.committee.validators)

	for i := 0; i < omissionEpochPeriod; i++ {
		setupProofAndAutonityFinalize(r, r.committee.validators[0].NodeAddress, absentees)
	}

	epochID, _, err := r.autonity.EpochID(nil)
	require.NoError(t, err)
	require.Equal(t, common.Big1.String(), epochID.String())

	// validator 1 should still be in the committee and not jailed
	require.Equal(t, csize, len(r.committee.validators))
	val := validator(r, r.committee.validators[1].NodeAddress)
	require.Equal(t, uint8(0), val.State)
	require.Equal(t, 0, offences(r, r.committee.validators[1].NodeAddress))
}

func TestProtocolParameterChange(t *testing.T) {
	r := setup(t, nil)
	/* default config:
	- epoch period: 50
	- lookback window: 40
	- delta: 5
	*/
	lookback, _, err := r.omissionAccountability.GetLookbackWindow(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(40), lookback.Uint64())
	delta, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(5), delta.Uint64())
	epochPeriod, _, err := r.autonity.GetEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(50), epochPeriod.Uint64())

	// past performance weight cannot be greater than inactivity threshold
	_, err = r.omissionAccountability.SetPastPerformanceWeight(r.operator, omissionScaleFactor(r))
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
	_, err = r.autonity.SetEpochPeriod(r.operator, new(big.Int).SetUint64(30))
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
	_, err = r.autonity.SetEpochPeriod(r.operator, new(big.Int).SetUint64(100))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 100
		- delta: 5
		- lookback window: 40
		updated attempted equation:
		100 > 5+100-1 --> false --> err
	*/
	_, err = r.omissionAccountability.SetLookbackWindow(r.operator, new(big.Int).SetUint64(100))
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
	_, err = r.omissionAccountability.SetLookbackWindow(r.operator, new(big.Int).SetUint64(60))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 100
		- delta: 5
		- lookback window: 60
		updated attempted equation:
		30 > 5+20-1 --> true --> no error
	*/
	_, err = r.omissionAccountability.SetLookbackWindow(r.operator, new(big.Int).SetUint64(20))
	require.NoError(t, err)
	_, err = r.autonity.SetEpochPeriod(r.operator, new(big.Int).SetUint64(30))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 30
		- delta: 5
		- lookback window: 20
		updated attempted equation:
		30 > 20+20-1 --> false --> err
	*/
	_, err = r.omissionAccountability.SetDelta(r.operator, new(big.Int).SetUint64(20))
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
	_, err = r.omissionAccountability.SetDelta(r.operator, new(big.Int).SetUint64(10))
	require.NoError(t, err)

	// params should still be unchanged in current epoch
	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(40), config.LookbackWindow.Uint64())
	require.Equal(t, uint64(5), config.Delta.Uint64())
	epochPeriod, _, err = r.autonity.GetCurrentEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(50), epochPeriod.Uint64())

	// getters should already return new values
	lookback, _, err = r.omissionAccountability.GetLookbackWindow(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(20), lookback.Uint64())
	delta, _, err = r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), delta.Uint64())
	epochPeriod, _, err = r.autonity.GetEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(30), epochPeriod.Uint64())

	// end epoch, new protocol params should be applied
	// both getters and config should return new values
	r.waitNextEpoch()

	config, _, err = r.omissionAccountability.Config(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(20), config.LookbackWindow.Uint64())
	require.Equal(t, uint64(10), config.Delta.Uint64())
	epochPeriod, _, err = r.autonity.GetCurrentEpochPeriod(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(30), epochPeriod.Uint64())

	lookback, _, err = r.omissionAccountability.GetLookbackWindow(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(20), lookback.Uint64())
	delta, _, err = r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)
	require.Equal(t, uint64(10), delta.Uint64())
	epochPeriod, _, err = r.autonity.GetEpochPeriod(nil)
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
	_, err = r.autonity.SetEpochPeriod(r.operator, new(big.Int).SetUint64(40))
	require.NoError(t, err)
	_, err = r.omissionAccountability.SetDelta(r.operator, new(big.Int).SetUint64(20))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 40
		- delta: 20
		- lookback window: 20
		updated attempted equation:
		400 > 20+20-1 --> true --> no err
	*/
	_, err = r.autonity.SetEpochPeriod(r.operator, new(big.Int).SetUint64(400))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 400
		- delta: 20
		- lookback window: 20
		updated attempted equation:
		400 > 20+200-1 --> true --> no err
	*/
	_, err = r.omissionAccountability.SetLookbackWindow(r.operator, new(big.Int).SetUint64(200))
	require.NoError(t, err)

	/*
		current params:
		- epoch period: 400
		- delta: 20
		- lookback window: 20
		updated attempted equation:
		400 > 1+200-1 --> true --> no err
	*/
	_, err = r.omissionAccountability.SetDelta(r.operator, common.Big1)
	require.NoError(t, err)

	// cannot set lookback and delta to 0
	_, err = r.omissionAccountability.SetLookbackWindow(r.operator, new(big.Int).SetUint64(0))
	t.Log(err)
	require.Error(t, err)
	_, err = r.omissionAccountability.SetDelta(r.operator, new(big.Int).SetUint64(0))
	t.Log(err)
	require.Error(t, err)

	// can set to 1
	_, err = r.omissionAccountability.SetDelta(r.operator, new(big.Int).SetUint64(1))
	require.NoError(t, err)
	_, err = r.omissionAccountability.SetLookbackWindow(r.operator, new(big.Int).SetUint64(1))
	require.NoError(t, err)

}
