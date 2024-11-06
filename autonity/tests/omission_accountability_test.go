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
func omissionFinalize(r *Runner, epochEnd bool) {
	_, err := r.OmissionAccountability.Finalize(FromAutonity, epochEnd)
	require.NoError(r.T, err)
	r.T.Logf("Omission accountability, finalized block: %d", r.Evm.Context.BlockNumber)
	// advance the block context as if we mined a block
	r.Evm.Context.BlockNumber = new(big.Int).Add(r.Evm.Context.BlockNumber, common.Big1)
	r.Evm.Context.Time = new(big.Int).Add(r.Evm.Context.Time, common.Big1)
	// clean up activity proof data
	r.Evm.Context.Coinbase = common.Address{}
	r.Evm.Context.ActivityProof = nil
	r.Evm.Context.ActivityProofRound = 0
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

func slashingRatePrecision(r *Runner) *big.Int {
	precision, _, err := r.slasherContract().GetSlashingPrecision(nil)
	require.NoError(r.T, err)
	return precision
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
	r := Setup(t, configOverrideIncreasedStake)

	// set maximum inactivity threshold for this test, we care only about the inactivity counters and not about the jailing
	_, err := r.OmissionAccountability.SetInactivityThreshold(r.Operator, new(big.Int).SetUint64(10000))
	require.NoError(t, err)

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.WaitNBlocks(int(delta.Int64()))

	config, _, err := r.OmissionAccountability.Config(nil)
	require.NoError(r.T, err)
	lookback := int(config.LookbackWindow.Uint64())

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
		targetHeight := r.Evm.Context.BlockNumber.Int64() - delta.Int64()
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
	for i := r.Evm.Context.BlockNumber.Int64(); i < int64(omissionEpochPeriod); i++ {
		setupProofAndAutonityFinalize(r, proposer, nil)
	}
	t.Log("Closing epoch")
	setupProofAndAutonityFinalize(r, proposer, nil)
	r.generateNewCommittee()

	// inactivity counters should be reset
	require.Equal(r.T, 0, inactivityCounter(r, partiallyOffline))
	require.Equal(r.T, 0, inactivityCounter(r, fullyOffline))

	r.WaitNBlocks(int(delta.Int64()))
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
	r := Setup(t, configOverrideIncreasedStake)

	// set maximum inactivity threshold for this test, we care only about the inactivity scores and not about the jailing
	scaleFactorInt := omissionScaleFactor(r)
	_, err := r.OmissionAccountability.SetInactivityThreshold(r.Operator, scaleFactorInt)
	require.NoError(t, err)

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	scaleFactor := newFloat(scaleFactorInt)

	initialCommitteeSize := len(r.Committee.Validators)

	r.WaitNBlocks(int(delta.Int64()))

	config, _, err := r.OmissionAccountability.Config(nil)
	require.NoError(r.T, err)
	lookback := int(config.LookbackWindow.Uint64())
	pastPerformanceWeight := new(big.Float).Quo(newFloat(config.PastPerformanceWeight), scaleFactor)
	currentPerformanceWeight := new(big.Float).Sub(newFloat(common.Big1), pastPerformanceWeight)

	// simulate epoch.
	proposer := r.Committee.Validators[0].NodeAddress
	inactiveBlockStreak := make([]int, len(r.Committee.Validators))
	inactiveCounters := make([]*big.Int, len(r.Committee.Validators))
	for i := range inactiveCounters {
		inactiveCounters[i] = new(big.Int)
	}
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
		absentees := make(map[common.Address]struct{})
		for i := range r.Committee.Validators {
			if r.Committee.Validators[i].NodeAddress == proposer {
				continue // keep proposer always online
			}
			if rand.Intn(30) != 0 {
				absentees[r.Committee.Validators[i].NodeAddress] = struct{}{}
				inactiveBlockStreak[i]++
			} else {
				inactiveBlockStreak[i] = 0
			}
			if inactiveBlockStreak[i] >= lookback {
				inactiveCounters[i] = new(big.Int).Add(inactiveCounters[i], common.Big1)
			}
		}

		t.Logf("number of absentees: %d for height %d", len(absentees), r.Evm.Context.BlockNumber.Uint64())
		setupProofAndAutonityFinalize(r, proposer, absentees)
	}
	r.generateNewCommittee()

	// no one should have been jailed
	require.Equal(t, initialCommitteeSize, len(r.Committee.Validators))

	// check score computation
	pastInactivityScore := make([]*big.Float, len(r.Committee.Validators))
	denominator := new(big.Int).SetUint64(uint64(omissionEpochPeriod) - delta.Uint64() - uint64(lookback) + 1)
	for i, val := range r.Committee.Validators {
		score := new(big.Float).Quo(newFloat(inactiveCounters[i]), newFloat(denominator))

		expectedInactivityScoreFloat := new(big.Float).Mul(score, currentPerformanceWeight) // + 0 * pastPerformanceWeight
		pastInactivityScore[i] = expectedInactivityScoreFloat

		expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, scaleFactor)
		expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Int(nil)

		r.T.Logf("i %d, expectedInactivityScoreFloat %s, expectedInactivityScore %s, inactivityScore %v", i, toString(expectedInactivityScoreFloat), expectedInactivityScore.String(), inactivityScore(r, val.NodeAddress))

		require.Equal(r.T, int(expectedInactivityScore.Uint64()), inactivityScore(r, val.NodeAddress))
	}

	// simulate another epoch
	r.WaitNBlocks(int(delta.Int64()))
	inactiveBlockStreak = make([]int, len(r.Committee.Validators))
	inactiveCounters = make([]*big.Int, len(r.Committee.Validators))
	for i := range inactiveCounters {
		inactiveCounters[i] = new(big.Int)
	}
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
		absentees := make(map[common.Address]struct{})
		for i := range r.Committee.Validators {
			if r.Committee.Validators[i].NodeAddress == proposer {
				continue // keep proposer always online
			}
			if rand.Intn(30) != 0 {
				absentees[r.Committee.Validators[i].NodeAddress] = struct{}{}
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
	require.Equal(t, initialCommitteeSize, len(r.Committee.Validators))

	// check score computation
	for i, val := range r.Committee.Validators {
		score := new(big.Float).Quo(newFloat(inactiveCounters[i]), newFloat(denominator))

		expectedInactivityScoreFloat1 := new(big.Float).Mul(score, currentPerformanceWeight)
		expectedInactivityScoreFloat2 := new(big.Float).Mul(pastInactivityScore[i], pastPerformanceWeight)
		expectedInactivityScoreFloat := new(big.Float).Add(expectedInactivityScoreFloat1, expectedInactivityScoreFloat2)

		expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, scaleFactor)
		expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Int(nil)

		r.T.Logf("i %d, expectedInactivityScoreFloat %s, expectedInactivityScore %s, inactivityScore %v", i, toString(expectedInactivityScoreFloat), expectedInactivityScore.String(), inactivityScore(r, val.NodeAddress))
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
	expectedSlashAmount.Div(expectedSlashAmount, slashingRatePrecision(r))
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
		proposerRewardRatePrecisionBig, _, err := r.Autonity.PROPOSERREWARDRATEPRECISION(nil)
		require.NoError(t, err)
		proposerRewardRatePrecision := newFloat(proposerRewardRatePrecisionBig)
		t.Logf("proposer reward rate: %s, precision: %s", toString(proposerRewardRate), toString(proposerRewardRatePrecision))

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
		ntnBalanceBefore := newFloat(ntnBalance(r, proposerTreasury))
		t.Logf("atn balance before: %s, ntn balance before %s", toString(atnBalanceBefore), toString(ntnBalanceBefore))

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

		atnActualBalance := r.GetBalanceOf(proposerTreasury)
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
		r := Setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			config.EpochPeriod = uint64(omissionEpochPeriod)
			return config
		})

		delta, _, err := r.OmissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.WaitNBlocks(int(delta.Int64()))

		totalEffort := new(big.Int)
		efforts := make([]*big.Int, len(r.Committee.Validators))
		atnBalances := make([]*big.Int, len(r.Committee.Validators))
		ntnBalances := make([]*big.Int, len(r.Committee.Validators))
		totalPower := new(big.Int)
		for i, val := range r.Committee.Validators {
			efforts[i] = new(big.Int)
			atnBalances[i] = r.GetBalanceOf(val.Treasury)
			ntnBalances[i] = ntnBalance(r, val.Treasury)
			totalPower.Add(totalPower, val.BondedStake)
		}
		// simulate epoch
		fullProofEffort := new(big.Int).Sub(totalPower, bft.Quorum(totalPower)) // effort for a full proof
		for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
			proposerIndex := rand.Intn(len(r.Committee.Validators))
			totalEffort.Add(totalEffort, fullProofEffort)
			efforts[proposerIndex].Add(efforts[proposerIndex], fullProofEffort)
			r.setupActivityProofAndCoinbase(r.Committee.Validators[proposerIndex].NodeAddress, nil)
			omissionFinalize(r, h == omissionEpochPeriod)
		}

		simulatedNtnRewards := new(big.Int).SetInt64(5968565)
		simulatedAtnRewards := new(big.Int).SetInt64(4545445)
		r.GiveMeSomeMoney(r.Autonity.address, simulatedAtnRewards)
		_, err = r.Autonity.Mint(r.Operator, r.OmissionAccountability.address, simulatedNtnRewards)
		require.NoError(r.T, err)
		_, err = r.OmissionAccountability.DistributeProposerRewards(&runOptions{origin: r.Autonity.address, value: simulatedAtnRewards}, simulatedNtnRewards)
		require.NoError(t, err)

		for i, val := range r.Committee.Validators {
			atnExpectedIncrement := new(big.Int).Mul(efforts[i], simulatedAtnRewards)
			atnExpectedIncrement.Div(atnExpectedIncrement, totalEffort)
			ntnExpectedIncrement := new(big.Int).Mul(efforts[i], simulatedNtnRewards)
			ntnExpectedIncrement.Div(ntnExpectedIncrement, totalEffort)
			atnExpectedBalance := new(big.Int).Add(atnBalances[i], atnExpectedIncrement)
			ntnExpectedBalance := new(big.Int).Add(ntnBalances[i], ntnExpectedIncrement)
			t.Logf("validator %d, effort %s, total effort %s, expectedBalance atn %s, expectedBalanceNtn %s", i, efforts[i].String(), totalEffort.String(), atnExpectedBalance.String(), ntnExpectedBalance.String())

			atnBalance := r.GetBalanceOf(val.Treasury)
			ntnBalance := ntnBalance(r, val.Treasury)

			require.Equal(t, atnExpectedBalance.String(), atnBalance.String())
			require.Equal(t, ntnExpectedBalance.String(), ntnBalance.String())

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
		return config
	})

	delta, _, err := r.OmissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	// validators over threshold will get all their rewards withheld
	customInactivityThreshold := uint64(6000)
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
				continue // let's keep at least a guy inside the committee
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
	ntnRewards := new(big.Int).SetUint64(8220842843566600000)
	r.GiveMeSomeMoney(r.Autonity.address, atnRewards)

	atnBalancesBefore := make([]*big.Int, len(r.Committee.Validators))
	ntnBalancesBefore := make([]*big.Int, len(r.Committee.Validators))
	totalPower := new(big.Int)
	for i, val := range r.Committee.Validators {
		validatorStruct := validator(r, val.NodeAddress)
		// we assume that all stake is self bonded in this test
		require.Equal(t, validatorStruct.SelfBondedStake.String(), validatorStruct.BondedStake.String())
		atnBalancesBefore[i] = r.GetBalanceOf(val.Treasury)
		ntnBalancesBefore[i] = ntnBalance(r, val.Treasury)
		t.Logf("validator %d, atn balance before: %s, ntn balance before %s", i, atnBalancesBefore[i].String(), ntnBalancesBefore[i].String())
		totalPower.Add(totalPower, validatorStruct.SelfBondedStake)
	}
	atnPoolBefore := r.GetBalanceOf(withheldRewardPool)
	ntnPoolBefore := ntnBalance(r, withheldRewardPool)

	setupProofAndAutonityFinalize(r, proposer, nil)

	atnTotalWithheld := new(big.Int)
	ntnTotalWithheld := new(big.Int)
	for i, val := range r.Committee.Validators {
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
		require.Equal(t, atnExpectedBalance.String(), r.GetBalanceOf(val.Treasury).String())
		require.Equal(t, ntnExpectedBalance.String(), ntnBalance(r, val.Treasury).String())
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

}
