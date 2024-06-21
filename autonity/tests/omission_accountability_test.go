package tests

//var omissionEpochPeriod = 100
//var inflationAfter100Blocks = 6311834092292000000 //TODO(lorenzo) remove

const ScaleFactor = 10_000

/*
// need a longer epoch for omission accountability tests
var configOverride = func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
	config.EpochPeriod = uint64(omissionEpochPeriod)
	return config
}

// helpers
// TODO(lorenzo) refactor
// func omissionFinalize(r *runner, absents []common.Address, proposer common.Address, effort *big.Int, proposerFaulty bool, epochEnd bool) {
func omissionFinalize(r *runner, epochEnd bool) {
	//_, err := r.omissionAccountability.Finalize(fromAutonity, absents, proposer, effort, proposerFaulty, epochEnd)
	_, err := r.omissionAccountability.Finalize(fromAutonity, epochEnd)
	require.NoError(r.t, err)
	r.t.Logf("Omission accountability, finalized block: %d", r.evm.Context.BlockNumber)
	// advance the block context as if we mined a block
	r.evm.Context.BlockNumber = new(big.Int).Add(r.evm.Context.BlockNumber, common.Big1)
	r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
}

// TODO(lorenzo) refactor
// func autonityFinalize(r *runner, absents []common.Address, proposer common.Address, effort *big.Int, proposerFaulty bool) { //nolint
func autonityFinalize(r *runner) { //nolint
	_, err := r.autonity.Finalize(nil)
	require.NoError(r.t, err)
	r.t.Logf("Autonity, finalized block: %d", r.evm.Context.BlockNumber)
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

func validator(r *runner, addr common.Address) AutonityValidator {
	val, _, err := r.autonity.GetValidator(nil, addr)
	require.NoError(r.t, err)
	return val
}

func ntnBalance(r *runner, addr common.Address) *big.Int {
	balance, _, err := r.autonity.BalanceOf(nil, addr)
	require.NoError(r.t, err)
	return balance
}*/

/*
TODO(lorenzo) refactor this tests on the new precompile implementation
func TestAccessControl(t *testing.T) {
	r := setup(t, nil)

	delta, _, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Uint64()))

	_, err = r.omissionAccountability.Finalize(r.operator, []common.Address{}, common.Address{}, common.Big0, true, false)
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

		delta, _, _, err := r.omissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.waitNBlocks(int(delta.Int64()))

		targetHeight := r.evm.Context.BlockNumber.Int64() - delta.Int64()
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

		delta, _, _, err := r.omissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.waitNBlocks(int(delta.Int64()))

		targetHeight := r.evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer := r.committee.validators[0].NodeAddress
		omissionFinalize(r, []common.Address{}, proposer, common.Big1, false, false)
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, 1, proposerEffort(r, proposer))
		require.Equal(r.t, 1, totalProposerEffort(r))

		targetHeight = r.evm.Context.BlockNumber.Int64() - delta.Int64()
		proposer = r.committee.validators[0].NodeAddress
		omissionFinalize(r, []common.Address{}, proposer, common.Big3, false, false)
		require.False(r.t, faultyProposer(r, targetHeight))
		require.Equal(r.t, 0, inactivityCounter(r, proposer))
		require.Equal(r.t, 4, proposerEffort(r, proposer))
		require.Equal(r.t, 4, totalProposerEffort(r))

		targetHeight = r.evm.Context.BlockNumber.Int64() - delta.Int64()
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

	delta, _, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

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
		targetHeight := r.evm.Context.BlockNumber.Int64() - delta.Int64()
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
	autonityFinalize(r, []common.Address{}, proposer, common.Big1, false)

	// inactivity counters should be reset
	require.Equal(r.t, 0, inactivityCounter(r, partiallyOffline))
	require.Equal(r.t, 0, inactivityCounter(r, fullyOffline))

	r.waitNBlocks(int(delta.Int64()))
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

	delta, _, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	lookback := int(config.LookbackWindow.Uint64())
	pastPerformanceWeight := float64(config.PastPerformanceWeight.Uint64()) / ScaleFactor

	// simulate epoch
	inactiveBlockStreak := make([]int, len(r.committee.validators))
	inactiveCounters := make([]int, len(r.committee.validators))
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
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

	// NOTE: theoretically the validators over the threshold would get jailed and therefore excluded from the next committee
	// however since we care only about scores in this test I just do not call the autonity finalize, so that they remain in the committee
	// even if they are jailed. We just set the last epoch block.
	_, err = r.omissionAccountability.SetLastEpochBlock(fromAutonity, new(big.Int).SetUint64(uint64(omissionEpochPeriod)))
	require.NoError(t, err)

	// check score computation
	pastInactivityScore := make([]float64, len(r.committee.validators))
	for i, val := range r.committee.validators {
		score := float64(inactiveCounters[i]) / float64(omissionEpochPeriod-int(delta.Int64())-lookback+1)
		score = math.Floor(score*ScaleFactor) / ScaleFactor // mimic precision loss due to fixed point arithmetic used in solidity
		expectedInactivityScoreFloat := score*(1-pastPerformanceWeight) + 0*pastPerformanceWeight
		expectedInactivityScoreFloat = math.Floor(expectedInactivityScoreFloat*ScaleFactor) / ScaleFactor // mimic precision loss due to fixed point arithmetic used in solidity
		pastInactivityScore[i] = expectedInactivityScoreFloat
		expectedInactivityScore := int(math.Round(expectedInactivityScoreFloat * ScaleFactor)) // using round to mitigate precision loss due to floating point arithmetic
		r.t.Logf("expectedInactivityScore %v, inactivityScore %v", expectedInactivityScore, inactivityScore(r, val.NodeAddress))
		require.Equal(r.t, expectedInactivityScore, inactivityScore(r, val.NodeAddress))
	}

	// simulate another epoch
	r.waitNBlocks(int(delta.Int64()))
	inactiveBlockStreak = make([]int, len(r.committee.validators))
	inactiveCounters = make([]int, len(r.committee.validators))
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
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
		score := float64(inactiveCounters[i]) / float64(omissionEpochPeriod-int(delta.Int64())-lookback+1)
		score = math.Floor(score*ScaleFactor) / ScaleFactor // mimic precision loss due to fixed point arithmetic used in solidity
		expectedInactivityScoreFloat := score*(1-pastPerformanceWeight) + pastInactivityScore[i]*pastPerformanceWeight

		expectedInactivityScoreFloatScaled := expectedInactivityScoreFloat * ScaleFactor
		// detect and address floating point precision loss. A bit hackish but it works
		// this is to address where floating point representation makes us end up with number like 3533.999999999999 instead of 3534
		if math.Floor(expectedInactivityScoreFloatScaled+0.0000000001) > math.Floor(expectedInactivityScoreFloatScaled) {
			t.Log("Detected and corrected floating point precision loss")
			expectedInactivityScoreFloatScaled = math.Floor(expectedInactivityScoreFloatScaled) + 1
		}
		expectedInactivityScoreFraction := math.Floor(expectedInactivityScoreFloatScaled) / ScaleFactor // mimic precision loss due to fixed point arithmetic used in solidity
		expectedInactivityScore := int(math.Round(expectedInactivityScoreFraction * ScaleFactor))       // round to mitigate precision loss due to floating point
		r.t.Logf("expectedInactivityScore %v, inactivityScore %v", expectedInactivityScore, inactivityScore(r, val.NodeAddress))
		require.Equal(r.t, expectedInactivityScore, inactivityScore(r, val.NodeAddress))
	}
}

func TestOmissionPunishments(t *testing.T) {
	r := setup(t, configOverride)

	delta, _, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	initialJailingPeriod := int(config.InitialJailingPeriod.Uint64())
	initialProbationPeriod := int(config.InitialProbationPeriod.Uint64())
	pastPerformanceWeight := int(config.PastPerformanceWeight.Uint64())
	initialSlashingRate := int(config.InitialSlashingRate.Uint64())
	slashingRatePrecision := int(config.SlashingRatePrecision.Uint64())

	proposer := r.committee.validators[0].NodeAddress
	absents := []common.Address{r.committee.validators[1].NodeAddress, r.committee.validators[2].NodeAddress}
	treasuries := []common.Address{r.committee.validators[1].Treasury, r.committee.validators[2].Treasury}

	// simulate epoch with two validator at 100% inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		omissionFinalize(r, absents, proposer, common.Big1, false, false)
	}
	autonityFinalize(r, absents, proposer, common.Big1, false)

	// the two validators should have been jailed and be under probation + offence counter should have been incremented
	expectedFullOfflineScore := ScaleFactor - pastPerformanceWeight
	for _, absent := range absents {
		require.Equal(r.t, expectedFullOfflineScore, inactivityScore(r, absent))
		val := validator(r, absent)
		require.Equal(r.t, uint8(2), val.State)
		require.Equal(r.t, uint64(omissionEpochPeriod+initialJailingPeriod), val.JailReleaseBlock.Uint64())
		require.Equal(r.t, initialProbationPeriod, probation(r, absent))
		require.Equal(t, 1, offences(r, absent))
	}

	// wait that the jailing finishes and reactivate validators
	r.waitNBlocks(initialJailingPeriod)
	for i, absent := range absents {
		_, err = r.autonity.ActivateValidator(&runOptions{origin: treasuries[i]}, absent)
		require.NoError(r.t, err)
	}

	// pass some epochs, probation period should decrease
	r.waitNextEpoch() // re-activation epoch, val not part of committee
	// inactivity score should still be the same as before
	for _, absent := range absents {
		require.Equal(r.t, expectedFullOfflineScore, inactivityScore(r, absent))
	}
	r.waitNextEpoch()
	// should be decreased  now
	for _, absent := range absents {
		require.Equal(r.t, (expectedFullOfflineScore*pastPerformanceWeight)/ScaleFactor, inactivityScore(r, absent))
	}
	r.waitNextEpoch()
	r.waitNextEpoch()

	// probation periods should have decreased of once for every epoch that passed with the validator as part of the committee
	passedEpochs := 3

	for _, absent := range absents {
		require.Equal(t, initialProbationPeriod-passedEpochs, probation(r, absent))
	}

	// simulate another epoch where:
	// - val 1 gets slashed by accountability and therefore doesn't get punished by omission accountability
	// - val 2 gets punished again for omission while in the probation period, therefore he gets slashed
	r.waitNBlocks(int(delta.Int64()))
	val1Address := absents[0]
	val2Address := absents[1]
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		omissionFinalize(r, absents, proposer, common.Big1, false, false)
	}
	val1 := validator(r, val1Address)
	val1.State = uint8(2) // jailed
	val1.JailReleaseBlock = new(big.Int).SetInt64(r.evm.Context.BlockNumber.Int64() + int64(omissionEpochPeriod*10))
	totalSlashedVal1 := val1.TotalSlashed
	_, err = r.autonity.UpdateValidatorAndTransferSlashedFunds(&runOptions{origin: r.accountability.address}, val1)
	require.NoError(t, err)

	val2BeforeSlash := validator(r, val2Address)
	autonityFinalize(r, absents, proposer, common.Big1, false)

	// val1, punished by accountability, shouldn't have been slashed by omission even if 100% offline and still under probation
	val1 = validator(r, val1Address)
	require.Equal(r.t, uint8(2), val1.State)
	require.True(r.t, probation(r, val1.NodeAddress) > 0)
	require.Equal(r.t, totalSlashedVal1.String(), val1.TotalSlashed.String())
	require.Equal(r.t, 1, offences(r, val1Address))

	// val2 offline while on probation, should have been slashed by omission
	val2 := validator(r, val2Address)
	require.Equal(r.t, uint8(2), val2.State)
	require.True(r.t, probation(r, val2.NodeAddress) > 0)
	require.True(r.t, val2.TotalSlashed.Cmp(val2BeforeSlash.TotalSlashed) > 0)
	require.Equal(r.t, 2, offences(r, val2Address))
	expectedSlashRate := new(big.Int).SetInt64(int64(initialSlashingRate * 4 * 2)) // rate * offence^2 * collusion
	availableFunds := new(big.Int).Add(val2BeforeSlash.BondedStake, val2.UnbondingStake)
	availableFunds.Add(availableFunds, val2.SelfUnbondingStake)
	expectedSlashAmount := new(big.Int).Mul(expectedSlashRate, availableFunds)
	expectedSlashAmount.Div(expectedSlashAmount, new(big.Int).SetInt64(int64(slashingRatePrecision)))
	require.Equal(r.t, expectedSlashAmount.String(), new(big.Int).Sub(val2.TotalSlashed, val2BeforeSlash.TotalSlashed).String())
}

func TestProposerRewardDistribution(t *testing.T) {
	t.Run("Rewards are correctly allocated based on config", func(t *testing.T) {
		r := setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			config.EpochPeriod = uint64(omissionEpochPeriod)
			return config
		})

		maxCommitteeSize, _, err := r.autonity.GetMaxCommitteeSize(nil)
		require.NoError(r.t, err)
		config, _, err := r.autonity.Config(nil)
		require.NoError(r.t, err)
		proposerRewardRate := config.Policy.ProposerRewardRate.Uint64()
		treasuryRate := config.Policy.TreasuryFee
		proposerRewardRatePrecisionBig, _, err := r.autonity.PROPOSERREWARDRATEPRECISION(nil)
		require.NoError(t, err)
		proposerRewardRatePrecision := float64(proposerRewardRatePrecisionBig.Uint64())
		committeeFactorPrecisionBig, _, err := r.autonity.COMMITTEEFRACTIONPRECISION(nil)
		require.NoError(t, err)
		committeeFactorPrecision := float64(committeeFactorPrecisionBig.Uint64())

		autonityAtns := new(big.Int).SetInt64(54644455456465)               // random amount
		ntnRewards := new(big.Int).SetInt64(int64(inflationAfter100Blocks)) // this has to match the ntn inflation unlocked NTNs
		r.giveMeSomeMoney(r.autonity.address, autonityAtns)

		// compute actual rewards for validator (subtract treasury fee)
		treasuryFee := new(big.Int).Mul(treasuryRate, autonityAtns)
		ten := new(big.Int).SetInt64(10)
		eighteen := new(big.Int).SetInt64(18)
		treasuryFee.Div(treasuryFee, new(big.Int).Exp(ten, eighteen, nil))
		atnRewards := new(big.Int).Sub(autonityAtns, treasuryFee)

		// all rewards should go to val 0
		proposer := r.committee.validators[0].NodeAddress
		atnBalanceBefore := float64(r.getBalanceOf(proposer).Uint64())
		ntnBalanceBefore := float64(ntnBalance(r, proposer).Uint64())

		// set validator state to jailed so that he will not receive any reward other the proposer one
		val := validator(r, proposer)
		val.State = uint8(2) // jailed
		_, err = r.autonity.UpdateValidatorAndTransferSlashedFunds(&runOptions{origin: r.accountability.address}, val)
		require.NoError(t, err)

		r.evm.Context.BlockNumber = new(big.Int).SetInt64(int64(omissionEpochPeriod))
		r.evm.Context.Time.Add(r.evm.Context.Time, new(big.Int).SetInt64(int64(omissionEpochPeriod-1)))
		autonityFinalize(r, []common.Address{}, proposer, common.Big1, false)

		committeeFactor := float64(len(r.committee.validators)) / float64(maxCommitteeSize.Int64())
		committeeFactor = math.Floor(committeeFactor*committeeFactorPrecision) / committeeFactorPrecision // simulate loss of precision due to fixed point arithmetic
		atnExpectedReward := (float64(atnRewards.Uint64()) * committeeFactor * float64(proposerRewardRate)) / proposerRewardRatePrecision
		ntnExpectedReward := (float64(ntnRewards.Uint64()) * committeeFactor * float64(proposerRewardRate)) / proposerRewardRatePrecision

		atnExpectedBalance := int64(math.Floor(atnBalanceBefore + atnExpectedReward))
		ntnExpectedBalance := int64(math.Floor(ntnBalanceBefore + ntnExpectedReward))
		require.Equal(t, atnExpectedBalance, r.getBalanceOf(proposer).Int64())
		require.Equal(t, ntnExpectedBalance, ntnBalance(r, proposer).Int64())
	})
	t.Run("Rewards are correctly distributed among proposers", func(t *testing.T) {
		r := setup(t, configOverride)

		delta, _, _, err := r.omissionAccountability.GetDelta(nil)
		require.NoError(t, err)

		r.waitNBlocks(int(delta.Int64()))

		totalEffort := new(big.Int)
		efforts := make([]*big.Int, len(r.committee.validators))
		atnBalances := make([]*big.Int, len(r.committee.validators))
		ntnBalances := make([]*big.Int, len(r.committee.validators))
		for i, val := range r.committee.validators {
			efforts[i] = new(big.Int)
			atnBalances[i] = r.getBalanceOf(val.Treasury)
			ntnBalances[i] = ntnBalance(r, val.Treasury)
		}
		// simulate epoch
		for h := int(delta.Int64()) + 1; h < omissionEpochPeriod+1; h++ {
			proposerIndex := rand.Intn(len(r.committee.validators))
			totalEffort.Add(totalEffort, common.Big1)
			efforts[proposerIndex].Add(efforts[proposerIndex], common.Big1)
			omissionFinalize(r, []common.Address{}, r.committee.validators[proposerIndex].NodeAddress, common.Big1, false, h == omissionEpochPeriod)
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

			atnBalance := r.getBalanceOf(val.Treasury)
			ntnBalance := ntnBalance(r, val.Treasury)

			require.Equal(t, atnExpectedBalance.String(), atnBalance.String())
			require.Equal(t, ntnExpectedBalance.String(), ntnBalance.String())

			// effort counters should be zeroed out
			require.Equal(r.t, 0, proposerEffort(r, val.NodeAddress))
		}

		require.Equal(r.t, 0, totalProposerEffort(r))
	})
}

// past performance weight and inactivity threshold should be set low enough that if:
// - a validator gets 100% inactivity in epoch x
// - then he gets 0% inactivity in epoch x+n (after he reactivated)
// he shouldn't get slashed in epoch x+n
func TestConfigSanity(t *testing.T) {
	r := setup(t, configOverride)

	delta, _, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.omissionAccountability.Config(nil)
	require.NoError(r.t, err)
	initialJailingPeriod := int(config.InitialJailingPeriod.Uint64())

	proposer := r.committee.validators[0].NodeAddress
	absents := []common.Address{r.committee.validators[1].NodeAddress, r.committee.validators[2].NodeAddress}
	treasuries := []common.Address{r.committee.validators[1].Treasury, r.committee.validators[2].Treasury}

	// simulate epoch with two validator at 100% inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		omissionFinalize(r, absents, proposer, common.Big1, false, false)
	}
	autonityFinalize(r, absents, proposer, common.Big1, false)

	for _, absent := range absents {
		val := validator(r, absent)
		require.Equal(r.t, uint8(2), val.State)
		require.Equal(t, 1, offences(r, absent))
	}

	// wait that the jailing finishes and reactivate validators
	r.waitNBlocks(initialJailingPeriod)
	for i, absent := range absents {
		_, err = r.autonity.ActivateValidator(&runOptions{origin: treasuries[i]}, absent)
		require.NoError(r.t, err)
	}

	r.waitNextEpoch() // re-activation epoch, val not part of committee
	r.waitNextEpoch()

	// validator should not have been punished since he did 0% offline
	for _, absent := range absents {
		val := validator(r, absent)
		require.Equal(r.t, uint8(0), val.State)
		require.Equal(t, 1, offences(r, absent))
	}

}

func TestRewardWithholding(t *testing.T) {
	r := setup(t, func(config *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		config.EpochPeriod = uint64(omissionEpochPeriod)
		config.ProposerRewardRate = 0 // no rewards to proposers to make computation simpler
		config.TreasuryFee = 0        // same
		return config
	})

	delta, _, _, err := r.omissionAccountability.GetDelta(nil)
	require.NoError(t, err)

	r.waitNBlocks(int(delta.Int64()))

	config, _, err := r.autonity.Config(nil)
	require.NoError(t, err)
	withheldRewardPool := config.Policy.WithheldRewardsPool

	proposer := r.committee.validators[0].NodeAddress

	// simulate epoch with random levels of inactivity
	for h := int(delta.Int64()) + 1; h < omissionEpochPeriod; h++ {
		var absents []common.Address
		for i := range r.committee.validators {
			if i == 0 {
				continue // let's keep at least a guy inside the committee
			}
			if rand.Intn(30) != 0 {
				absents = append(absents, r.committee.validators[i].NodeAddress)
			}
		}
		omissionFinalize(r, absents, proposer, common.Big1, false, false)
	}

	atnRewards := new(big.Int).SetInt64(5467879877987)                  // random amount
	ntnRewards := new(big.Int).SetInt64(int64(inflationAfter100Blocks)) // this has to match the ntn inflation unlocked NTNs
	r.giveMeSomeMoney(r.autonity.address, atnRewards)

	atnBalancesBefore := make([]*big.Int, len(r.committee.validators))
	ntnBalancesBefore := make([]*big.Int, len(r.committee.validators))
	totalPower := new(big.Int)
	for i, val := range r.committee.validators {
		validatorStruct := validator(r, val.NodeAddress)
		// we assume that all stake is self bonded in this test
		require.Equal(t, validatorStruct.SelfBondedStake.String(), validatorStruct.BondedStake.String())
		atnBalancesBefore[i] = r.getBalanceOf(val.NodeAddress)
		ntnBalancesBefore[i] = ntnBalance(r, val.NodeAddress)
		totalPower.Add(totalPower, validatorStruct.SelfBondedStake)
	}
	atnPoolBefore := r.getBalanceOf(withheldRewardPool)
	ntnPoolBefore := ntnBalance(r, withheldRewardPool)
	autonityFinalize(r, []common.Address{}, proposer, common.Big1, false)

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
		t.Logf("validator index %d, score: %d", i, score.Uint64())
		atnWithheld := new(big.Int).Mul(atnFullReward, score)
		atnWithheld.Div(atnWithheld, omissionScaleFactor(r))
		ntnWithheld := new(big.Int).Mul(ntnFullReward, score)
		ntnWithheld.Div(ntnWithheld, omissionScaleFactor(r))
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
*/
