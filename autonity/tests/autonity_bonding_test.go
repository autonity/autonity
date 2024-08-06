package tests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core"
)

func TestBondingRequest(t *testing.T) {
	r := setup(t, nil)

	r.run("Test bond to a valid validator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore)

		rr.keepLogs(true)
		_, err = rr.autonity.Bond(&runOptions{origin: account}, rr.committee.validators[0].NodeAddress, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseNewBondingRequest, &AutonityTestNewBondingRequest{
			Validator:  rr.committee.validators[0].NodeAddress,
			Delegator:  account,
			SelfBonded: false,
			Amount:     amount,
		}), "bond should emit a NewBondingRequest event with the correct params")

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceAfter.Int64())

		bondingRequestId := len(rr.committee.validators)
		bondingRequest, _, err := rr.autonity.GetBondingRequest(nil, big.NewInt(int64(bondingRequestId)))
		require.NoError(t, err)
		require.Equal(t, account, bondingRequest.Delegator)
		require.Equal(t, rr.committee.validators[0].NodeAddress, bondingRequest.Delegatee)
		require.Equal(t, amount, bondingRequest.Amount)

		valLiquid := r.committee.liquidContracts[0]
		accountBalance, _, err := valLiquid.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), accountBalance.Int64())

		rr.waitNextEpoch()

		accountBalance, _, err = valLiquid.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Condition(t, func() bool {
			return accountBalance.Cmp(big.NewInt(0)) > 0
		})
	})

	r.run("Test validator self bonding", func(rr *runner) {
		treasury := rr.committee.validators[0].Treasury
		val := rr.committee.validators[0]
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, treasury, amount)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore)

		rr.keepLogs(true)
		_, err = rr.autonity.Bond(&runOptions{origin: treasury}, val.NodeAddress, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseNewBondingRequest, &AutonityTestNewBondingRequest{
			Validator:  val.NodeAddress,
			Delegator:  treasury,
			SelfBonded: true,
			Amount:     amount,
		}), "bond should emit a NewBondingRequest event with the correct params")

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceAfter.Int64())

		// bondingRequestId is the length of the validators array because each validator has bonded once
		bondingRequestId := len(rr.committee.validators)
		bondingRequest, _, err := rr.autonity.GetBondingRequest(nil, big.NewInt(int64(bondingRequestId)))
		require.NoError(t, err)

		require.Equal(t, treasury, bondingRequest.Delegator)
		require.Equal(t, val.NodeAddress, bondingRequest.Delegatee)
		require.Equal(t, amount, bondingRequest.Amount)

		// wait till end of epoch
		rr.waitNextEpoch()

		valLiquid := r.committee.liquidContracts[0]
		treasuryBalance, _, err := valLiquid.BalanceOf(nil, treasury)
		require.NoError(t, err)
		require.Equal(t, int64(0), treasuryBalance.Int64())
	})

	r.run("Test does not bond on a non registered validator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore)

		_, err = rr.autonity.Bond(&runOptions{origin: account}, rr.randomAccount(), amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: validator not registered")

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfter)
	})
}

func TestBondingAndUnbondingRequests(t *testing.T) {
	// this is all some complex setup to manipulate the inflation rate to 0
	r := setup(t, func(genesis *core.Genesis) *core.Genesis {
		genesis.Config.InflationContractConfig.InflationRateInitial = math.NewHexOrDecimal256(0)
		// set transition period very far away
		transitionPeriod, ok := big.NewInt(0).SetString("1000000000000000000000000", 10)
		require.True(t, ok)
		trans := math.HexOrDecimal256(*transitionPeriod)
		genesis.Config.InflationContractConfig.InflationTransitionPeriod = &trans
		genesis.Config.InflationContractConfig.InflationRateTransition = math.NewHexOrDecimal256(0)
		genesis.Config.InflationContractConfig.InflationCurveConvexity = math.NewHexOrDecimal256(0)

		// set epoch period to 10 blocks
		genesis.Config.AutonityContractConfig.EpochPeriod = uint64(10)
		genesis.Config.AutonityContractConfig.UnbondingPeriod = uint64(20)
		return genesis
	})

	r.run("Test cannot bond to a paused validator", func(rr *runner) {
		_, err := rr.autonity.PauseValidator(
			&runOptions{origin: rr.committee.validators[0].NodeAddress},
			rr.committee.validators[0].NodeAddress,
		)
		require.NoError(t, err)

		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err = rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceBefore)

		_, err = rr.autonity.Bond(&runOptions{origin: account}, rr.committee.validators[0].NodeAddress, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: validator need to be active")

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfter)
	})

	r.run("Test self-unbonding", func(rr *runner) {
		treasury := rr.committee.validators[0].Treasury
		val := rr.committee.validators[0]
		amount := big.NewInt(1e18)

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)

		rr.keepLogs(true)

		_, err = rr.autonity.Unbond(&runOptions{origin: treasury}, val.NodeAddress, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseNewUnbondingRequest, &AutonityTestNewUnbondingRequest{
			Validator:  val.NodeAddress,
			Delegator:  treasury,
			Amount:     amount,
			SelfBonded: true,
		}), "unbond should emit a NewUnbondingRequest event with the correct params")

		unbondingReqId := big.NewInt(int64(0))

		unbondingRequest, _, err := rr.autonity.GetUnbondingRequest(nil, unbondingReqId)
		require.NoError(t, err)
		require.Equal(t, treasury, unbondingRequest.Delegator)
		require.Equal(t, val.NodeAddress, unbondingRequest.Delegatee)
		require.Equal(t, amount, unbondingRequest.Amount)
		require.Equal(t, int64(0), unbondingRequest.UnbondingShare.Int64())
		require.Equal(t, false, unbondingRequest.Unlocked)

		valInfo, _, err := rr.autonity.GetValidator(nil, val.NodeAddress)
		require.NoError(t, err)
		require.Equal(t, amount, valInfo.SelfUnbondingStakeLocked)

		rr.waitNextEpoch()

		unbondingRequest, _, err = rr.autonity.GetUnbondingRequest(nil, unbondingReqId)
		require.NoError(t, err)
		require.Equal(t, amount, unbondingRequest.UnbondingShare)
		require.Equal(t, true, unbondingRequest.Unlocked)

		newValInfo, _, err := rr.autonity.GetValidator(nil, val.NodeAddress)
		verifyValidatorInfoPostUnbonding(t, &newValInfo, &valInfo, amount, amount)

		// before unbonding period, but after epoch end, the balance should be locked
		balanceAfter, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)
		require.Equal(t, balanceBefore, balanceAfter)

		rr.waitNextEpoch()
		rr.waitNextEpoch()

		balanceAfterRelease, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Add(balanceBefore, amount), balanceAfterRelease)

		valInfoAfterRelease, _, err := rr.autonity.GetValidator(nil, val.NodeAddress)
		verifyValidatorInfoPostRelease(t, &valInfoAfterRelease, &newValInfo, amount, amount)
	})

	r.run("Test unbonding from a valid validator", func(rr *runner) {
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		_, err = rr.autonity.Bond(&runOptions{origin: account}, rr.committee.validators[0].NodeAddress, amount)
		require.NoError(t, err)

		rr.waitNextEpoch()

		balanceBefore, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBefore.Int64())

		expectedValInfo, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)

		rr.keepLogs(true)
		_, err = rr.autonity.Unbond(&runOptions{origin: account}, rr.committee.validators[0].NodeAddress, amount)
		require.NoError(t, err)
		emitsEvent(rr.Logs(), rr.autonity.ParseNewUnbondingRequest, &AutonityTestNewUnbondingRequest{
			Validator:  rr.committee.validators[0].NodeAddress,
			Delegator:  account,
			SelfBonded: false,
			Amount:     amount,
		})

		balanceAfter, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceAfter.Int64())

		// this is the first unbonding request
		unbondingRequestId := big.NewInt(int64(0))

		unbondingRequest, _, err := rr.autonity.GetUnbondingRequest(nil, unbondingRequestId)
		// verify unbonding request details
		require.NoError(t, err)
		require.Equal(t, account, unbondingRequest.Delegator)
		require.Equal(t, rr.committee.validators[0].NodeAddress, unbondingRequest.Delegatee)
		require.Equal(t, amount, unbondingRequest.Amount)
		require.Equal(t, int64(0), unbondingRequest.UnbondingShare.Int64())
		require.Equal(t, false, unbondingRequest.Unlocked)

		// check effects of unbond (non-self-bonded):
		// LNTN is locked
		liquidContract := rr.committee.liquidContracts[0]
		require.NoError(t, err)

		lockedBalance, _, err := liquidContract.LockedBalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, lockedBalance)

		rr.waitNextEpoch()

		// verify that the unbonding request is unlocked
		ubRequest, _, err := rr.autonity.GetUnbondingRequest(nil, unbondingRequestId)
		require.NoError(t, err)
		require.Equal(t, amount, ubRequest.UnbondingShare)
		require.Equal(t, true, ubRequest.Unlocked)

		valInfo, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)
		verifyValidatorInfoPostUnbonding(t, &valInfo, &expectedValInfo, big.NewInt(0), amount)

		lockedBalanceAfterEpoch, _, err := liquidContract.LockedBalanceOf(nil, account)
		require.NoError(t, err)

		liquidBalanceAfterEpoch, _, err := liquidContract.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), lockedBalanceAfterEpoch.Int64())
		require.Equal(t, int64(0), liquidBalanceAfterEpoch.Int64())

		balanceBeforeRelease, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, int64(0), balanceBeforeRelease.Int64(), "NTN released before unbonding period")

		// unbonding period = 20 blocks = 2 epochs
		rr.waitNextEpoch()
		rr.waitNextEpoch()

		balanceAfterRelease, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, amount, balanceAfterRelease, "NTN not released after unbonding period")

		valInfoAfterRelease, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)

		verifyValidatorInfoPostRelease(t, &valInfoAfterRelease, &valInfo, big.NewInt(0), amount)
	})
}

func TestBondingAndUnbondingRequestsFailures(t *testing.T) {
	// again we manipulate the config to set inflation to zero, this is to ensure we can correctly
	// test the unbonding balance changes without interference
	r := setup(t, func(genesis *core.Genesis) *core.Genesis {
		genesis.Config.InflationContractConfig.InflationRateInitial = math.NewHexOrDecimal256(0)
		// set transition period very far away
		transitionPeriod, ok := big.NewInt(0).SetString("1000000000000000000000000", 10)
		require.True(t, ok)
		trans := math.HexOrDecimal256(*transitionPeriod)
		genesis.Config.InflationContractConfig.InflationTransitionPeriod = &trans
		genesis.Config.InflationContractConfig.InflationRateTransition = math.NewHexOrDecimal256(0)
		genesis.Config.InflationContractConfig.InflationCurveConvexity = math.NewHexOrDecimal256(0)
		return genesis
	})

	r.run("Test ensure it does not unbond from a non-registered validator", func(rr *runner) {
		amount := big.NewInt(1e18)
		notAValidator := rr.randomAccount()

		_, err := rr.autonity.Unbond(&runOptions{origin: rr.committee.validators[0].Treasury}, notAValidator, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: validator not registered")
	})

	r.run("Test can't unbond from a validator with an amount exeeding balance", func(rr *runner) {
		account := rr.committee.validators[0].Treasury
		val, _, err := rr.autonity.GetValidator(nil, rr.committee.validators[0].NodeAddress)
		require.NoError(t, err)

		tooManyTokens := new(big.Int).Add(val.SelfBondedStake, big.NewInt(1e18))

		_, err = rr.autonity.Unbond(
			&runOptions{origin: account},
			rr.committee.validators[0].NodeAddress,
			tooManyTokens,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: insufficient self bonded newton balance")
	})

	r.run("Test cant unbond non-selfbonded stake of amount 0", func(rr *runner) {
		account := rr.randomAccount()
		validator := rr.committee.validators[0].NodeAddress

		_, err := rr.autonity.Unbond(&runOptions{origin: account}, validator, big.NewInt(0))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: unbonding amount is 0")

		rr.waitNextEpoch()
	})

	r.run("Test self-unbond cannot exceed bonded", func(rr *runner) {
		treasury := rr.committee.validators[0].Treasury
		val := rr.committee.validators[0].NodeAddress

		valInfo, _, err := rr.autonity.GetValidator(nil, val)
		require.NoError(t, err)
		initialBonded := valInfo.SelfBondedStake
		amount := new(big.Int).Add(initialBonded, big.NewInt(1e18))

		_, err = rr.autonity.Mint(rr.operator, treasury, amount)
		require.NoError(t, err)

		initialBalance, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)

		_, err = rr.autonity.Bond(&runOptions{origin: treasury}, val, amount)
		require.NoError(t, err)

		rr.waitNextEpoch()

		// self unbond twice
		rr.keepLogs(true)
		_, err = rr.autonity.Unbond(&runOptions{origin: treasury}, val, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseNewUnbondingRequest, &AutonityTestNewUnbondingRequest{
			Validator:  val,
			Delegator:  treasury,
			SelfBonded: true,
			Amount:     amount,
		}), "bond should emit a NewBondingRequest event with the correct params")

		// should fail the second time
		_, err = rr.autonity.Unbond(&runOptions{origin: treasury}, val, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: insufficient self bonded newton balance")

		rr.waitNextEpoch()

		currentUnbondingPeriod, _, err := rr.autonity.GetUnbondingPeriod(nil)
		require.NoError(t, err)

		rr.waitNBlocks(int(currentUnbondingPeriod.Int64()))

		finalBalance, _, err := rr.autonity.BalanceOf(nil, treasury)
		require.NoError(t, err)
		require.Equal(t, initialBalance, finalBalance)
	})

	r.run("Test non-self unbond cannot exceed bonded", func(rr *runner) {
		account := rr.randomAccount()
		val := rr.committee.validators[0].NodeAddress
		amount := big.NewInt(1e18)

		_, err := rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		initialBalance, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)

		_, err = rr.autonity.Bond(&runOptions{origin: account}, val, amount)
		require.NoError(t, err)

		rr.waitNextEpoch()

		// unbond twice
		rr.keepLogs(true)
		_, err = rr.autonity.Unbond(&runOptions{origin: account}, val, amount)
		require.NoError(t, err)
		require.True(t, emitsEvent(rr.Logs(), rr.autonity.ParseNewUnbondingRequest, &AutonityTestNewUnbondingRequest{
			Validator:  val,
			Delegator:  account,
			SelfBonded: false,
			Amount:     amount,
		}), "bond should emit a NewBondingRequest event with the correct params")

		// should fail the second time
		_, err = rr.autonity.Unbond(&runOptions{origin: account}, val, amount)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted: insufficient unlocked Liquid Newton balance")

		rr.waitNextEpoch()

		unbondingPeriod, _, err := rr.autonity.GetUnbondingPeriod(nil)
		require.NoError(t, err)

		rr.waitNBlocks(int(unbondingPeriod.Int64()))
		finalBalance, _, err := rr.autonity.BalanceOf(nil, account)
		require.NoError(t, err)
		require.Equal(t, initialBalance, finalBalance)
	})
}

func TestBondingAndUnbondingQueues(t *testing.T) {
	r := setup(t, nil)

	r.run("Test bonding queue logic", func(rr *runner) {
		numStakings := big.NewInt(int64(len(rr.committee.validators)))

		tailBondingId, _, err := rr.autonity.GetTailBondingID(nil)
		require.NoError(t, err)

		headBondingReqId, _, err := rr.autonity.GetHeadBondingID(nil)
		require.NoError(t, err)

		// there should be no pending requests
		require.GreaterOrEqual(t,
			tailBondingId.Int64(),
			new(big.Int).Sub(headBondingReqId, big.NewInt(1)).Int64(),
			"pending bonding request found",
		)

		// check ids start from 0
		latestBondingId := new(big.Int).Sub(numStakings, big.NewInt(1))
		require.Equal(t, new(big.Int).Sub(headBondingReqId, big.NewInt(1)), latestBondingId)

		// new bonding request
		account := rr.randomAccount()
		amount := big.NewInt(1e18)

		_, err = rr.autonity.Mint(rr.operator, account, amount)
		require.NoError(t, err)

		_, err = rr.autonity.Bond(&runOptions{origin: account}, rr.committee.validators[0].NodeAddress, amount)
		require.NoError(t, err)
		numStakings = new(big.Int).Add(numStakings, big.NewInt(1))

		// check head and tail are updated
		latestBondingId = new(big.Int).Sub(numStakings, big.NewInt(1))
		headBondingReqId, _, err = rr.autonity.GetHeadBondingID(nil)
		require.NoError(t, err)
		tailBondingId, _, err = rr.autonity.GetTailBondingID(nil)
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Sub(headBondingReqId, big.NewInt(1)), latestBondingId)
		require.Equal(t, tailBondingId, latestBondingId)

		bondingRequest, _, err := rr.autonity.GetBondingRequest(nil, latestBondingId)
		require.NoError(t, err)
		require.Equal(t, account, bondingRequest.Delegator)
		require.Equal(t, rr.committee.validators[0].NodeAddress, bondingRequest.Delegatee)
		require.Equal(t, amount, bondingRequest.Amount)
	})

	r.run("Test unbonding queue logic", func(rr *runner) {
		lastUnlockedUnbonding, _, err := rr.autonity.GetLastUnlockedUnbonding(nil)
		require.NoError(t, err)

		headUnbondingId, _, err := rr.autonity.GetHeadUnbondingID(nil)
		require.NoError(t, err)

		require.GreaterOrEqual(t,
			lastUnlockedUnbonding.Int64(),
			headUnbondingId.Int64(),
			"last unlocked unbonding id should be greater than head unbonding id",
		)

		// there should be no pending requests
		require.Equal(t, headUnbondingId.Int64(), int64(0))

		// new unbonding request
		amount := big.NewInt(1e18)
		account := rr.committee.validators[0].Treasury

		_, err = rr.autonity.Unbond(&runOptions{origin: account}, rr.committee.validators[0].NodeAddress, amount)
		require.NoError(t, err)

		// check head and tail are updated
		latestUnbondingId := big.NewInt(0)

		headUnbondingId, _, err = rr.autonity.GetHeadUnbondingID(nil)
		require.NoError(t, err)
		lastUnlockedUnbonding, _, err = rr.autonity.GetLastUnlockedUnbonding(nil)
		require.NoError(t, err)

		require.Equal(t, new(big.Int).Sub(headUnbondingId, big.NewInt(1)).Int64(), latestUnbondingId.Int64())
		require.Equal(t, lastUnlockedUnbonding.Int64(), latestUnbondingId.Int64())

		unbondingRequest, _, err := rr.autonity.GetUnbondingRequest(nil, latestUnbondingId)
		require.NoError(t, err)

		require.Equal(t, account, unbondingRequest.Delegator)
		require.Equal(t, rr.committee.validators[0].NodeAddress, unbondingRequest.Delegatee)
		require.Equal(t, amount, unbondingRequest.Amount)
		require.Equal(t, false, unbondingRequest.Unlocked)
	})
}

// Unbonding utility functions

func verifyValidatorInfoPostUnbonding(
	t *testing.T,
	valInfo *AutonityValidator,
	expectedValidator *AutonityValidator,
	selfUnbonded,
	totalUnbonded *big.Int,
) {
	nonSelfUnbonded := new(big.Int).Sub(totalUnbonded, selfUnbonded)

	require.Equal(t, new(big.Int).Sub(expectedValidator.BondedStake, totalUnbonded), valInfo.BondedStake)
	require.Equal(t, new(big.Int).Sub(expectedValidator.SelfBondedStake, selfUnbonded), valInfo.SelfBondedStake)
	require.Equal(t, new(big.Int).Add(expectedValidator.UnbondingShares, nonSelfUnbonded).String(), valInfo.UnbondingShares.String())
	require.Equal(t, new(big.Int).Add(expectedValidator.UnbondingStake, nonSelfUnbonded).String(), valInfo.UnbondingStake.String())
	require.Equal(t,
		new(big.Int).Add(expectedValidator.SelfUnbondingStake, selfUnbonded).String(),
		valInfo.SelfUnbondingStake.String(),
	)
	require.Equal(t,
		new(big.Int).Add(expectedValidator.SelfUnbondingShares, selfUnbonded).String(),
		valInfo.SelfUnbondingShares.String(),
	)
	require.Equal(t,
		new(big.Int).Sub(expectedValidator.LiquidSupply, nonSelfUnbonded).String(),
		valInfo.LiquidSupply.String(),
	)
}

func verifyValidatorInfoPostRelease(
	t *testing.T,
	valInfo *AutonityValidator,
	expectedValidator *AutonityValidator,
	selfUnbonded,
	totalUnbonded *big.Int,
) {
	nonSelfUnbonded := new(big.Int).Sub(totalUnbonded, selfUnbonded)
	require.Equal(t, expectedValidator.BondedStake.String(), valInfo.BondedStake.String())
	require.Equal(t, expectedValidator.SelfBondedStake.String(), valInfo.SelfBondedStake.String())
	require.Equal(t,
		new(big.Int).Sub(expectedValidator.UnbondingShares, nonSelfUnbonded).String(),
		valInfo.UnbondingShares.String(),
	)
	require.Equal(t,
		new(big.Int).Sub(expectedValidator.UnbondingStake, nonSelfUnbonded).String(),
		valInfo.UnbondingStake.String(),
	)
	require.Equal(t,
		new(big.Int).Sub(expectedValidator.SelfUnbondingStake, selfUnbonded).String(),
		valInfo.SelfUnbondingStake.String(),
	)
	require.Equal(t,
		new(big.Int).Sub(expectedValidator.SelfUnbondingShares, selfUnbonded).String(),
		valInfo.SelfUnbondingShares.String(),
	)
	require.Equal(t, expectedValidator.LiquidSupply.String(), valInfo.LiquidSupply.String())
}
