package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
)

var fromAutonity = &runOptions{origin: params.AutonityContractAddress}

var reward = big.NewInt(1000_000_000)

func TestBondingGasConsumption(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	contractCount := 10
	start := r.evm.Context.Time.Int64()
	cliff := start
	end := 1000 + start
	for i := 0; i < contractCount; i++ {
		createContract(r, user, contractTotalAmount, start, cliff, end)
	}
	validator := r.committee.validators[0].NodeAddress
	liquidContract := r.committee.liquidContracts[0]
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	stakingGas, _, err := r.autonity.StakingGasPrice(nil)
	require.NoError(r.t, err)
	bondingAmount := big.NewInt(contractTotalAmount)
	r.NoError(
		r.autonity.Mint(operator, user, bondingAmount),
	)
	r.NoError(
		r.autonity.Bond(fromSender(user, nil), validator, bondingAmount),
	)
	initBalance := new(big.Int).Mul(big.NewInt(1000_000), big.NewInt(1000_000_000_000_000_000))
	r.giveMeSomeMoney(user, initBalance)
	r.waitNextEpoch()

	r.run("single bond", func(r *runner) {
		bondingID := len(r.committee.validators) + 1
		var iteration int64 = 10
		bondingAmount := big.NewInt(contractTotalAmount / iteration)
		for ; iteration > 0; iteration-- {
			gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, user, liquidContract, bondingID, common.Big0, bondingAmount, bondingGas, false)
			totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedBond)), stakingGas)
			require.True(
				r.t,
				bondingGas.Cmp(totalGasUsed) >= 0,
				"need more gas to notify bonding operations",
			)
			bondingID++
		}
	})

	r.run("multiple bonding", func(r *runner) {
		oldLiquidBalance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		bondingID := len(r.committee.validators) + 1
		bondingAmount := big.NewInt(contractTotalAmount)
		for i := 1; i < contractCount; i++ {
			r.NoError(
				r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount),
			)
			bondingID++
		}
		r.waitNextEpoch()
		delegatedStake := big.NewInt(contractTotalAmount * int64(contractCount-1))
		liquidBalance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.Equal(r.t, oldLiquidBalance.Add(oldLiquidBalance, delegatedStake), liquidBalance)
		for i := 1; i < contractCount; i++ {
			liquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
			require.NoError(r.t, err)
			require.Equal(r.t, bondingAmount, liquidBalance)
		}
		gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, user, liquidContract, bondingID, common.Big0, bondingAmount, bondingGas, false)
		totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedBond)), stakingGas)
		require.True(
			r.t,
			bondingGas.Cmp(totalGasUsed) >= 0,
			"need more gas to notify bonding operations",
		)
	})

	r.run("bonding rejected", func(r *runner) {
		bondingID := len(r.committee.validators) + 1
		bondingAmount := big.NewInt(contractTotalAmount)
		gasUsedDistribute, gasUsedBond := bondAndApply(r, validator, user, liquidContract, bondingID, common.Big0, bondingAmount, bondingGas, true)
		totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedBond)), stakingGas)
		require.True(
			r.t,
			bondingGas.Cmp(totalGasUsed) >= 0,
			"need more gas to notify bonding operations",
		)
	})
}

func TestUnbondingGasConsumption(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	contractCount := 10
	start := r.evm.Context.Time.Int64()
	cliff := start
	end := 1000 + start
	for i := 0; i < contractCount; i++ {
		createContract(r, user, contractTotalAmount, start, cliff, end)
	}
	validator := r.committee.validators[0].NodeAddress
	liquidContract := r.committee.liquidContracts[0]
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	require.NoError(r.t, err)
	stakingGas, _, err := r.autonity.StakingGasPrice(nil)
	require.NoError(r.t, err)
	initBalance := new(big.Int).Mul(big.NewInt(1000_000), big.NewInt(1000_000_000_000_000_000))
	r.giveMeSomeMoney(user, initBalance)

	bondingAmount := big.NewInt(contractTotalAmount)
	for i := 0; i < contractCount; i++ {
		r.NoError(
			r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount),
		)
	}
	r.waitNextEpoch()

	r.run("single unbond", func(r *runner) {
		var iteration int64 = 10
		unbondingAmount := big.NewInt(contractTotalAmount / iteration)
		unbondingID := 0
		for ; iteration > 0; iteration-- {
			gasUsedDistribute, gasUsedUnbond, gasUsedRelease := unbondAndApply(r, validator, user, liquidContract, unbondingID, common.Big0, unbondingAmount, unbondingGas, false)
			totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedUnbond+gasUsedRelease)), stakingGas)
			require.True(
				r.t,
				unbondingGas.Cmp(totalGasUsed) >= 0,
				"need more gas to notify unbonding operations",
			)
			unbondingID++
		}
	})

	r.run("unbond rejected", func(r *runner) {
		unbondingID := 0
		unbondingAmount := big.NewInt(contractTotalAmount)
		gasUsedDistribute, gasUsedUnbond, gasUsedRelease := unbondAndApply(r, validator, user, liquidContract, unbondingID, common.Big0, unbondingAmount, unbondingGas, true)
		totalGasUsed := new(big.Int).Mul(big.NewInt(int64(gasUsedDistribute+gasUsedUnbond+gasUsedRelease)), stakingGas)
		require.True(
			r.t,
			unbondingGas.Cmp(totalGasUsed) >= 0,
			"need more gas to notify unbonding operations",
		)
	})
}

func TestRelease(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0
	// do not modify userBalance
	userBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)

	r.run("cannot release before cliff", func(r *runner) {
		r.waitSomeBlock(cliff)
		require.Equal(r.t, big.NewInt(cliff), r.evm.Context.Time, "time mismatch")
		_, _, err := r.stakableVesting.UnlockedFunds(nil, user, contractID)
		require.Equal(r.t, "execution reverted: cliff period not reached yet", err.Error())
		_, err = r.stakableVesting.ReleaseFunds(fromSender(user, nil), contractID)
		require.Equal(r.t, "execution reverted: cliff period not reached yet", err.Error())
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, userBalance, userNewBalance, "funds released before cliff period")
	})

	r.run("release calculation follows epoch based linear function in time", func(r *runner) {
		currentTime := r.waitSomeEpoch(cliff + 1)
		require.True(r.t, currentTime <= end, "release is not linear after end")
		// contract has the context of last block, so time is 1s less than currentTime
		unlocked := currentTime - 1 - start
		require.True(r.t, contractTotalAmount > unlocked, "cannot test if all funds unlocked")
		epochID, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		// mine some more blocks, release should be epoch based
		r.waitNBlocks(10)
		currentTime += 10
		checkReleaseAllNTN(r, user, contractID, big.NewInt(unlocked))

		r.waitNBlocks(10)
		currentTime += 10
		require.Equal(r.t, big.NewInt(currentTime), r.evm.Context.Time, "time mismatch, release won't work")
		// no more should be released as epoch did not change
		newEpochID, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, epochID, newEpochID, "cannot test if epoch progresses")
		checkReleaseAllNTN(r, user, contractID, common.Big0)
	})

	r.run("can release in chunks", func(r *runner) {
		currentTime := r.waitSomeEpoch(cliff + 1)
		require.True(r.t, currentTime <= end, "cannot test, release is not linear after end")
		totalUnlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, contractID)
		require.NoError(r.t, err)
		require.True(r.t, totalUnlocked.IsInt64(), "invalid data")
		require.True(r.t, totalUnlocked.Int64() > 1, "cannot test chunks")
		unlockFraction := big.NewInt(totalUnlocked.Int64() / 2)
		// release only a chunk of total unlocked
		r.NoError(
			r.stakableVesting.ReleaseNTN(fromSender(user, nil), contractID, unlockFraction),
		)
		userNewBalance, _, err := r.autonity.BalanceOf(nil, user)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(userBalance, unlockFraction), userNewBalance, "balance mismatch")
		data, _, err := r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.True(r.t, data.IsInt64(), "invalid data")
		epochID := data.Int64()
		r.waitNBlocks(10)
		data, _, err = r.autonity.EpochID(nil)
		require.NoError(r.t, err)
		require.True(r.t, data.IsInt64(), "invalid data")
		require.Equal(r.t, epochID, data.Int64(), "epoch progressed, more funds will release")
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(currentTime)) > 0, "time did not progress")
		checkReleaseAllNTN(r, user, contractID, new(big.Int).Sub(totalUnlocked, unlockFraction))
	})

	r.run("cannot release more than total", func(r *runner) {
		r.waitSomeEpoch(end + 1)
		// progress some more epoch, should not matter after end
		r.waitNextEpoch()
		currentTime := r.evm.Context.Time
		checkReleaseAllNTN(r, user, contractID, big.NewInt(contractTotalAmount))
		r.waitNextEpoch()
		require.True(r.t, r.evm.Context.Time.Cmp(currentTime) > 0, "time did not progress")
		// cannot release more
		checkReleaseAllNTN(r, user, contractID, common.Big0)
	})
}

func TestBonding(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	users, validators, liquidContracts := setupContracts(r, 2, 2, contractTotalAmount, start, cliff, end)

	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)

	beneficiary := users[0]
	contractID := common.Big0
	validator := validators[0]
	liquidContract := liquidContracts[0]

	r.run("can bond all funds before cliff but not before start", func(r *runner) {
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(start+1)) < 0, "contract started already")
		bondingAmount := big.NewInt(contractTotalAmount / 2)
		_, err := r.stakableVesting.Bond(fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: contract not started yet", err.Error())
		r.waitSomeBlock(start + 1)
		require.True(r.t, r.evm.Context.Time.Cmp(big.NewInt(cliff+1)) < 0, "contract cliff finished already")
		bondAndFinalize(r, beneficiary, validator, liquidContract, contractID, bondingAmount, bondingGas)
	})

	// start contract for bonding for all the tests remaining
	r.waitSomeBlock(start + 1)

	r.run("cannot bond more than total", func(r *runner) {
		bondingAmount := big.NewInt(contractTotalAmount + 10)
		_, err := r.stakableVesting.Bond(fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: not enough tokens", err.Error())
		bondingAmount = big.NewInt(contractTotalAmount / 2)
		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		r.NoError(
			r.stakableVesting.Bond(fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount),
		)
		contract, _, err := r.stakableVesting.GetContract(nil, beneficiary, contractID)
		require.NoError(r.t, err)
		require.Equal(r.t, remaining, contract.CurrentNTNAmount, "contract not updated properly")
		bondingAmount = new(big.Int).Add(big.NewInt(10), remaining)
		_, err = r.stakableVesting.Bond(fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: not enough tokens", err.Error())
		// let bonding apply
		r.waitNextEpoch()
		_, err = r.stakableVesting.Bond(fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount)
		require.Equal(r.t, "execution reverted: not enough tokens", err.Error())
		bondAndFinalize(r, beneficiary, validator, liquidContract, contractID, remaining, bondingGas)
	})

	r.run("can release liquid tokens", func(r *runner) {
		bondingAmount := big.NewInt(contractTotalAmount)
		r.NoError(
			r.stakableVesting.Bond(fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount),
		)
		// let bonding apply
		r.waitNextEpoch()
		currentTime := r.waitSomeEpoch(cliff + 1)
		// contract has context of last block
		unlocked := currentTime - 1 - start
		// mine some more block, release should be epoch based
		r.waitNBlocks(10)
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(beneficiary, nil), contractID),
		)
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.Equal(
			r.t, big.NewInt(contractTotalAmount-unlocked), liquid,
			"liquid release don't follow epoch based linear function",
		)
		liquid, _, err = liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.Equal(r.t, big.NewInt(contractTotalAmount-unlocked), liquid, "liquid not transferred")
		liquid, _, err = liquidContract.BalanceOf(nil, beneficiary)
		require.NoError(r.t, err)
		require.Equal(r.t, big.NewInt(unlocked), liquid, "liquid not received")
		r.waitSomeEpoch(end + 1)
		// progress more epoch, shouldn't matter
		r.waitNextEpoch()
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(beneficiary, nil), contractID),
		)
		liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.True(r.t, liquid.Cmp(common.Big0) == 0, "all liquid tokens not released")
		liquid, _, err = liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.True(r.t, liquid.Cmp(common.Big0) == 0, "liquid not transferred")
		liquid, _, err = liquidContract.BalanceOf(nil, beneficiary)
		require.NoError(r.t, err)
		require.Equal(r.t, big.NewInt(contractTotalAmount), liquid, "liquid not received")
	})

	r.run("track liquids when bonding from multiple contracts to multiple validators", func(r *runner) {
		// TODO: complete
	})

	r.run("when bonded, release NTN first", func(r *runner) {
		liquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.True(r.t, contractTotalAmount > 10, "cannot test")
		bondingAmount := big.NewInt(contractTotalAmount / 10)
		bondAndFinalize(r, beneficiary, validator, liquidContract, contractID, bondingAmount, bondingGas)
		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		require.True(r.t, remaining.Cmp(common.Big0) > 0, "no NTN remains")
		r.waitSomeEpoch(cliff + 1)
		unlocked, _, err := r.stakableVesting.UnlockedFunds(nil, beneficiary, contractID)
		require.NoError(r.t, err)
		require.True(r.t, unlocked.Cmp(remaining) < 0, "don't want to release all NTN in the test")
		balance, _, err := r.autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.t, err)
		r.NoError(
			r.stakableVesting.ReleaseFunds(fromSender(beneficiary, nil), contractID),
		)
		newLiquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(liquidBalance, bondingAmount), newLiquidBalance, "lquid released")
		newBalance, _, err := r.autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(balance, unlocked), newBalance, "balance not updated")
	})

	r.run("test release when bonding to multiple validator", func(r *runner) {
		// TODO: complete
	})
}

func TestUnbonding(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	validatorCount := 2
	contractCount := 2
	users, validators, liquidContracts := setupContracts(r, contractCount, validatorCount, contractTotalAmount, start, cliff, end)

	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	require.NoError(r.t, err)

	// bond from all contracts to all validators
	r.waitSomeBlock(start + 1)
	bondingAmount := big.NewInt(contractTotalAmount / int64(validatorCount))
	require.True(r.t, bondingAmount.Cmp(common.Big0) > 0, "not enough to bond")
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			for _, validator := range validators {
				r.NoError(
					r.stakableVesting.Bond(fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount),
				)
			}
		}
	}

	r.waitNextEpoch()
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			totalLiquid := big.NewInt(0)
			for _, validator := range validators {
				liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.t, err)
				require.Equal(r.t, bondingAmount, liquid)
				totalLiquid.Add(totalLiquid, liquid)
			}
			require.Equal(r.t, big.NewInt(contractTotalAmount), totalLiquid)
		}
	}

	// for testing single unbonding
	beneficiary := users[0]
	contractID := common.Big0
	validator := validators[0]
	liquidContract := liquidContracts[0]

	r.run("can unbond", func(r *runner) {
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, bondingAmount, liquid, "liquid not minted properly")
		unbondAndRelease(r, beneficiary, validator, liquidContract, contractID, liquid, unbondingGas)
	})

	r.run("cannot unbond more than total liquid", func(r *runner) {
		unbondingAmount := new(big.Int).Add(bondingAmount, big.NewInt(10))
		_, err = r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator, unbondingAmount)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		unbondingAmount = big.NewInt(10)
		remaining := new(big.Int).Sub(bondingAmount, unbondingAmount)
		require.True(r.t, remaining.Cmp(common.Big0) > 0, "cannot test if no liquid remains")
		r.NoError(
			r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator, unbondingAmount),
		)
		lockedLiquid, _, err := r.stakableVesting.LockedLiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, unbondingAmount, lockedLiquid)
		unbondingAmount = new(big.Int).Add(remaining, big.NewInt(10))
		_, err = r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator, unbondingAmount)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		r.waitNextEpoch()
		_, err = r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator, unbondingAmount)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		r.NoError(
			r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator, remaining),
		)
	})

	r.run("cannot unbond if released", func(r *runner) {
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		totalToRelease := liquid.Int64() + 10
		currentTime := r.waitSomeEpoch(totalToRelease + start + 1)
		totalToRelease = currentTime - 1 - start
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(beneficiary, nil), contractID),
		)
		_, err = r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator, liquid)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		// LNTN will be released from then first validator in the list
		newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.True(r.t, newLiquid.Cmp(common.Big0) == 0, "liquid remains after unbonding")
		// if more unlocked funds remain, then LNTN will be released from 2nd validator
		validator1 := validators[1]
		_, err = r.stakableVesting.Unbond(fromSender(beneficiary, unbondingGas), contractID, validator1, liquid)
		require.Equal(r.t, "execution reverted: not enough unlocked liquid tokens", err.Error())
		releasedFromAnotherValidator1 := totalToRelease - liquid.Int64()
		remainingLiquid := new(big.Int).Sub(liquid, big.NewInt(releasedFromAnotherValidator1))
		liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator1)
		require.NoError(r.t, err)
		require.Equal(r.t, remainingLiquid, liquid, "liquid balance mismatch")
		unbondAndRelease(r, beneficiary, validator1, liquidContracts[1], contractID, remainingLiquid, unbondingGas)
	})

	r.run("track liquid when unbonding from multiple contracts to multiple validators", func(r *runner) {
		// TODO: complete
	})
}

func TestStakingRevert(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("fails to notify reward distribution", func(r *runner) {
		// TODO: complete
	})

	r.run("reject bonding request and notify rejection", func(r *runner) {
		// TODO: complete
	})

	r.run("reject bonding request but fails to notify", func(r *runner) {
		// TODO: complete
	})

	r.run("revert applied bonding", func(r *runner) {
		// TODO: complete
	})

	r.run("reject unbonding request and notify rejection", func(r *runner) {
		// TODO: complete
	})

	r.run("reject unbonding request but fails to notify", func(r *runner) {
		// TODO: complete
	})

	r.run("revert applied unbonding", func(r *runner) {
		// TODO: complete
	})

	r.run("revert released unbonding", func(r *runner) {
		// TODO: complete
	})
}

func TestRwardTracking(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	contractCount := 2
	users, validators, liquidContracts := setupContracts(r, contractCount, 2, contractTotalAmount, start, cliff, end)

	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	// unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	// require.NoError(r.t, err)

	// for testing single unbonding
	beneficiary := users[0]
	contractID := common.Big0
	validator := validators[0]
	liquidContract := liquidContracts[0]

	// start contract to bond
	r.waitSomeBlock(start + 1)

	r.run("bond and claim reward", func(r *runner) {
		bondingAmount := big.NewInt(contractTotalAmount)
		r.NoError(
			r.stakableVesting.Bond(
				fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount,
			),
		)
		r.waitNextEpoch()

		r.NoError(
			r.autonity.ReceiveAut(
				fromSender(user, reward),
			),
		)
		r.waitNextEpoch()
		rewardOfContract, _, err := liquidContract.UnclaimedRewards(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.True(r.t, rewardOfContract.UnclaimedNTN.Cmp(common.Big0) > 0, "no NTN reward")
		require.True(r.t, rewardOfContract.UnclaimedATN.Cmp(common.Big0) > 0, "no ATN reward")
		rewardOfUser, _, err := r.stakableVesting.UnclaimedRewards0(nil, beneficiary)
		require.NoError(r.t, err)
		require.Equal(r.t, rewardOfContract.UnclaimedATN, rewardOfUser.AtnTotalFee, "ATN reward mismatch")
		require.Equal(r.t, rewardOfContract.UnclaimedNTN, rewardOfUser.NtnTotalFee, "NTN reward mismatch")
		balanceNTN, _, err := r.autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.t, err)
		balanceATN := r.getBalanceOf(beneficiary)
		r.tracing = true
		r.NoError(
			r.stakableVesting.ClaimRewards0(
				fromSender(beneficiary, nil),
			),
		)
		r.tracing = false
		newBalanceNTN, _, err := r.autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Add(balanceNTN, rewardOfUser.NtnTotalFee), newBalanceNTN, "NTN reward not claimed")
		newBalanceATN := r.getBalanceOf(beneficiary)
		require.Equal(r.t, new(big.Int).Add(balanceATN, rewardOfUser.AtnTotalFee), newBalanceATN, "ATN reward not claimed")
	})

	// set commission rate = 0, so all rewards go to delegation
	r.NoError(
		r.autonity.SetTreasuryFee(operator, common.Big0),
	)
	// remove all bonding, so we only have bonding from contracts only
	for _, validator := range r.committee.validators {
		require.Equal(r.t, validator.SelfBondedStake, validator.BondedStake, "delegation stake should not exist")
		r.NoError(
			r.autonity.Unbond(
				fromSender(validator.Treasury, nil), validator.NodeAddress, validator.SelfBondedStake,
			),
		)
		r.NoError(
			r.autonity.ChangeCommissionRate(
				fromSender(validator.Treasury, nil), validator.NodeAddress, common.Big0,
			),
		)
	}

	// bond from contracts
	bondingAmount := big.NewInt(100)
	totalBonded := new(big.Int)
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			for _, validator := range validators {
				r.NoError(
					r.stakableVesting.Bond(
						fromSender(user, bondingGas), big.NewInt(int64(i)), validator, bondingAmount,
					),
				)
				totalBonded.Add(totalBonded, bondingAmount)
			}
		}
	}
	r.waitNextEpoch()
	require.Equal(r.t, len(validators), len(r.committee.validators), "committee not updated properly")
	eachValidatorDelegation := big.NewInt(int64(len(users) * contractCount))
	eachValidatorStake := new(big.Int).Mul(bondingAmount, eachValidatorDelegation)
	for i, validator := range r.committee.validators {
		require.Equal(r.t, eachValidatorStake, validator.BondedStake)
		require.True(r.t, validator.SelfBondedStake.Cmp(common.Big0) == 0)
		balance, _, err := r.committee.liquidContracts[i].BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.Equal(r.t, eachValidatorStake, balance)
	}
	for r.committee.validators[0].CommissionRate.Cmp(common.Big0) > 0 {
		r.waitNextEpoch()
	}

	r.run("bond in differenet epoch and track reward", func(r *runner) {
		extraBond1 := bondingAmount
		r.NoError(
			r.stakableVesting.Bond(
				fromSender(beneficiary, bondingGas), contractID, validator, extraBond1,
			),
		)

		// check reward ratio, extraBond1 not applied yet
		r.giveMeSomeMoney(r.autonity.address, reward)
		totalReward := r.getBalanceOf(r.autonity.address)

		// every contract gets same reward
		validatorStakes := make(map[common.Address]*big.Int)
		userStakes := make(map[common.Address]map[common.Address]*big.Int)
		for _, user := range users {
			userStakes[user] = make(map[common.Address]*big.Int)
		}
		for _, validator := range validators {
			validatorStakes[validator] = new(big.Int).Add(common.Big0, eachValidatorStake)
			for _, user := range users {
				userStakes[user][validator] = new(big.Int).Mul(bondingAmount, big.NewInt(int64(contractCount)))
			}
		}
		totalStake := totalBonded
		checkRewards(r, totalStake, totalReward, validatorStakes, userStakes, liquidContracts, validators, users)
		for validator, stake := range validatorStakes {
			fmt.Printf("\nvalidator %v total stake %v\n", validator, stake)
		}
		for user, stake := range validatorStakes {
			fmt.Printf("\nuser %v total stake %v\n", user, stake)
		}

		extraBond2 := bondingAmount
		anotherContractID := common.Big1
		require.True(r.t, contractID.Cmp(anotherContractID) != 0, "cannot test with same contract")
		r.NoError(
			r.stakableVesting.Bond(
				fromSender(beneficiary, bondingGas), anotherContractID, validator, extraBond2,
			),
		)
		// check reward ratio after extraBond1 applied
		r.giveMeSomeMoney(r.autonity.address, reward)
		totalReward = r.getBalanceOf(r.autonity.address)

		// add extraBond1
		for validator, stake := range validatorStakes {
			fmt.Printf("\nvalidator %v total stake %v\n", validator, stake)
		}
		for user, stake := range validatorStakes {
			fmt.Printf("\nuser %v total stake %v\n", user, stake)
		}
		validatorStakes[validator].Add(validatorStakes[validator], extraBond1)
		userStakes[beneficiary][validator].Add(userStakes[beneficiary][validator], extraBond1)
		totalStake.Add(totalStake, extraBond1)
		for validator, stake := range validatorStakes {
			fmt.Printf("\nvalidator %v total stake %v\n", validator, stake)
		}
		for user, stake := range validatorStakes {
			fmt.Printf("\nuser %v total stake %v\n", user, stake)
		}
		checkRewards(r, totalStake, totalReward, validatorStakes, userStakes, liquidContracts, validators, users)
	})

	r.run("release liquid and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("unbond and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("unbond in different epoch and track reward", func(r *runner) {
		// TODO: complete
	})

	r.run("bond from multiple contracts to multiple validators and track reward", func(r *runner) {
		// TODO: complete
	})
}

func TestContractCancel(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0

	r.run("beneficiary changes when contract is canceled", func(r *runner) {
		newUser := common.HexToAddress("0x88")
		_, _, err := r.stakableVesting.GetContract(nil, user, contractID)
		require.NoError(r.t, err)
		_, _, err = r.stakableVesting.GetContract(nil, newUser, contractID)
		require.Equal(r.t, "execution reverted: invalid contract id", err.Error())
		r.stakableVesting.CancelContract(operator, user, contractID, newUser)
		_, _, err = r.stakableVesting.GetContract(nil, newUser, contractID)
		require.NoError(r.t, err)
		_, _, err = r.stakableVesting.GetContract(nil, user, contractID)
		require.Equal(r.t, "execution reverted: invalid contract id", err.Error())
	})
}

func TestContractUpdateWhenSlashed(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("contract total value update when bonded validator slashed", func(r *runner) {
		// TODO: complete
	})
}

func TestAccessRestriction(t *testing.T) {
	r := setup(t, nil)
	// TODO: complete setup

	r.run("only operator can create contract", func(r *runner) {
		// TODO: complete
	})

	r.run("only operator can cancel contract", func(r *runner) {
		// TODO: complete
	})

	r.run("only autonity can notify staking operations", func(r *runner) {
		// TODO: complete
		/*
			1. rewards distribution
			2. bonding applied
			3. unbonding applied
			4. unbonding released
		*/
	})

	r.run("only operator can set gas cost", func(r *runner) {
		// TODO: complete
	})

	r.run("cannot request bonding or unbonding without enough gas", func(r *runner) {
		// TODO: complete
	})
}

func checkRewards(
	r *runner,
	totalStake, totalReward *big.Int,
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[common.Address]*big.Int,
	liquidContracts []*Liquid,
	validators, users []common.Address,
) {

	fmt.Printf("\n\ntotal reward %v\ntotal stake %v\n", totalReward, totalStake)
	oldRewardsFromValidator := make(map[common.Address]*big.Int)
	for i, validator := range validators {
		liquidContract := liquidContracts[i]
		unclaimedReward, _, err := liquidContract.UnclaimedRewards(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		oldRewardsFromValidator[validator] = unclaimedReward.UnclaimedATN
		fmt.Printf("validator %v old reward %v\n", validator, oldRewardsFromValidator[validator])
	}

	oldUserRewards := make(map[common.Address]*big.Int)
	for _, user := range users {
		unclaimedReward, _, err := r.stakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.t, err)
		oldUserRewards[user] = unclaimedReward.AtnTotalFee
		fmt.Printf("user %v old reward %v\n", user, oldUserRewards[user])
	}

	r.waitNextEpoch()
	validatorCurrentReward := make(map[common.Address]*big.Int)
	for i, validator := range validators {
		liquidContract := liquidContracts[i]
		unclaimedReward, _, err := liquidContract.UnclaimedRewards(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		validatorTotalReward := new(big.Int).Mul(validatorStakes[validator], totalReward)
		validatorTotalReward = validatorTotalReward.Div(validatorTotalReward, totalStake)
		fmt.Printf("validator %v total stake %v\n", validator, validatorStakes[validator])
		fmt.Printf("validator %v current reward %v\n", validator, validatorTotalReward)
		fmt.Printf("validator %v total reward %v\n", validator, unclaimedReward.UnclaimedATN)
		require.Equal(
			r.t,
			new(big.Int).Add(validatorTotalReward, oldRewardsFromValidator[validator]),
			unclaimedReward.UnclaimedATN,
			"unclaimed atn reward not updated in liquid contract",
		)
		validatorCurrentReward[validator] = validatorTotalReward
	}

	for _, user := range users {
		userReward := new(big.Int)
		// The following loop is equivalent to: (userStakes[user][V0] + .. + userStakes[user][Vn]) * totalReward / totalStake
		// But StakableVesting contract handles reward for each validator separately, so there can be some funds lost due to
		// integer division in solidity. So we simulate the calculation with the for loop instead
		for _, validator := range validators {
			rewardFromValidator := new(big.Int).Mul(userStakes[user][validator], validatorCurrentReward[validator])
			rewardFromValidator.Div(rewardFromValidator, validatorStakes[validator])
			userReward.Add(userReward, rewardFromValidator)
		}
		unclaimedReward, _, err := r.stakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.t, err)
		// fmt.Printf("user %v total stake %v\n", user, userStakes[user])
		fmt.Printf("user %v current reward %v\n", user, userReward)
		fmt.Printf("user %v total reward %v\n", user, unclaimedReward.AtnTotalFee)
		require.Equal(
			r.t,
			new(big.Int).Add(userReward, oldUserRewards[user]),
			unclaimedReward.AtnTotalFee,
			"unclaimed atn reward mismatch",
		)
	}
}

func setupContracts(
	r *runner, contractCount, validatorCount int, contractTotalAmount, start, cliff, end int64,
) (users, validators []common.Address, liquidContracts []*Liquid) {
	users = make([]common.Address, 2)
	users[0] = user
	users[1] = common.HexToAddress("0x88")
	require.NotEqual(r.t, users[0], users[1], "same user")
	for _, user := range users {
		initBalance := new(big.Int).Mul(big.NewInt(1000_000), big.NewInt(1000_000_000_000_000_000))
		r.giveMeSomeMoney(user, initBalance)
		for i := 0; i < contractCount; i++ {
			createContract(r, user, contractTotalAmount, start, cliff, end)
		}
	}

	// use multiple validators
	validators = make([]common.Address, validatorCount)
	liquidContracts = make([]*Liquid, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validators[i] = r.committee.validators[i].NodeAddress
		liquidContracts[i] = r.committee.liquidContracts[i]
	}
	return
}

func createContract(r *runner, beneficiary common.Address, amount, startTime, cliffTime, endTime int64) {
	amountBigInt := big.NewInt(amount)
	r.NoError(
		r.autonity.Mint(operator, r.stakableVesting.address, amountBigInt),
	)
	r.NoError(
		r.stakableVesting.NewContract(
			operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
			big.NewInt(cliffTime), big.NewInt(endTime),
		),
	)
}

func fromSender(sender common.Address, value *big.Int) *runOptions {
	return &runOptions{origin: sender, value: value}
}

func bondAndApply(
	r *runner, validatorAddress, user common.Address, liquidContract *Liquid, bondingID int, contractID, bondingAmount, bondingGas *big.Int, rejected bool,
) (uint64, uint64) {
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validatorAddress)
	require.NoError(r.t, err)
	r.NoError(
		r.stakableVesting.Bond(fromSender(user, bondingGas), contractID, validatorAddress, bondingAmount),
	)
	r.giveMeSomeMoney(r.autonity.address, reward)
	r.NoError(
		liquidContract.Redistribute(fromSender(r.autonity.address, reward), common.Big0),
	)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(
		r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators),
	)
	if rejected == false {
		r.NoError(
			liquidContract.Mint(fromAutonity, r.stakableVesting.address, bondingAmount),
		)
		liquid = liquid.Add(liquid, bondingAmount)
	}
	gasUsedBond := r.NoError(
		r.stakableVesting.BondingApplied(
			fromAutonity, big.NewInt(int64(bondingID)), validatorAddress, bondingAmount, false, rejected,
		),
	)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	return gasUsedDistribute, gasUsedBond
}

func unbondAndApply(
	r *runner, validatorAddress, user common.Address, liquidContract *Liquid, unbondingID int, contractID, unbondingAmount, unbondingGas *big.Int, rejected bool,
) (uint64, uint64, uint64) {
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validatorAddress)
	require.NoError(r.t, err)
	r.NoError(
		r.stakableVesting.Unbond(fromSender(user, unbondingGas), contractID, validatorAddress, unbondingAmount),
	)
	r.giveMeSomeMoney(r.autonity.address, reward)
	r.NoError(
		liquidContract.Redistribute(fromSender(r.autonity.address, reward), common.Big0),
	)
	bondedValidators := make([]common.Address, 1)
	bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(
		r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators),
	)
	if rejected == false {
		r.NoError(
			liquidContract.Unlock(fromAutonity, r.stakableVesting.address, unbondingAmount),
		)
		r.NoError(
			liquidContract.Burn(fromAutonity, r.stakableVesting.address, unbondingAmount),
		)
		liquid = liquid.Sub(liquid, unbondingAmount)
	}
	gasUsedUnbond := r.NoError(
		r.stakableVesting.UnbondingApplied(fromAutonity, big.NewInt(int64(unbondingID)), validatorAddress, rejected),
	)
	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validatorAddress)
	require.NoError(r.t, err)
	require.Equal(r.t, liquid, newLiquid)
	gasUsedRelease := r.NoError(
		r.stakableVesting.UnbondingReleased(fromAutonity, big.NewInt(int64(unbondingID)), unbondingAmount, rejected),
	)
	return gasUsedDistribute, gasUsedUnbond, gasUsedRelease
}

func checkReleaseAllNTN(r *runner, user common.Address, contractID, unlockAmount *big.Int) {
	contract, _, err := r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	contractNTN := contract.CurrentNTNAmount
	initBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	totalUnlocked, _, err := r.stakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.t, err)
	require.True(r.t, unlockAmount.Cmp(totalUnlocked) == 0, "unlocked amount mismatch")
	r.NoError(
		r.stakableVesting.ReleaseAllNTN(fromSender(user, nil), contractID),
	)
	newBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	require.Equal(r.t, new(big.Int).Add(initBalance, totalUnlocked), newBalance, "balance mismatch")
	contract, _, err = r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	require.True(r.t, new(big.Int).Sub(contractNTN, unlockAmount).Cmp(contract.CurrentNTNAmount) == 0, "contract not updated properly")
}

func bondAndFinalize(
	r *runner, user, validator common.Address, liquidContract *Liquid, contractID, bondingAmount, bondingGas *big.Int,
) {
	newtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	liquidOfVestingContract, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	liquidOfUser, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
	require.NoError(r.t, err)
	contract, _, err := r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	remaining := new(big.Int).Sub(contract.CurrentNTNAmount, bondingAmount)
	r.NoError(
		r.stakableVesting.Bond(fromSender(user, bondingGas), contractID, validator, bondingAmount),
	)
	contract, _, err = r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	require.Equal(r.t, remaining, contract.CurrentNTNAmount, "contract not updated properly")
	// let bonding apply
	r.waitNextEpoch()
	liquidOfVestingContract = new(big.Int).Add(bondingAmount, liquidOfVestingContract)
	liquidOfUser = new(big.Int).Add(bondingAmount, liquidOfUser)
	liquid, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.Equal(r.t, liquidOfVestingContract, liquid, "liquid not minted")
	liquid, _, err = r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
	require.NoError(r.t, err)
	require.Equal(r.t, liquidOfUser, liquid, "vesting contract cannot track liquid balance")
	newNewtonBalance := new(big.Int).Sub(newtonBalance, bondingAmount)
	newtonBalance, _, err = r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.True(r.t, newNewtonBalance.Cmp(newtonBalance) == 0, "newton balance not updated")
}

func unbondAndRelease(
	r *runner, user, validator common.Address, liquidContract *Liquid, contractID, unbondingAmount, unbondingGas *big.Int,
) {
	newtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	liquidOfUser, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
	require.NoError(r.t, err)
	liquidOfVestingContract, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	contract, _, err := r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	unbondingRequestBlock := r.evm.Context.BlockNumber
	r.NoError(
		r.stakableVesting.Unbond(fromSender(user, unbondingGas), contractID, validator, unbondingAmount),
	)
	r.waitNextEpoch()
	liquidOfUser = new(big.Int).Sub(liquidOfUser, unbondingAmount)
	liquidOfVestingContract = new(big.Int).Sub(liquidOfVestingContract, unbondingAmount)
	liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
	require.NoError(r.t, err)
	require.True(r.t, liquid.Cmp(liquidOfUser) == 0, "liquid balance mismatch after unbonding")
	liquid, _, err = liquidContract.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.True(r.t, liquid.Cmp(liquidOfVestingContract) == 0, "liquid balance mismatch after unbonding")

	// release unbonding
	unbondingPeriod, _, err := r.autonity.GetUnbondingPeriod(nil)
	require.NoError(r.t, err)
	unbondingReleaseBlock := new(big.Int).Add(unbondingRequestBlock, unbondingPeriod)
	for unbondingReleaseBlock.Cmp(r.evm.Context.BlockNumber) >= 0 {
		r.waitNextEpoch()
	}
	newNewtonAmount := new(big.Int).Add(contract.CurrentNTNAmount, unbondingAmount)
	contract, _, err = r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	require.Equal(r.t, newNewtonAmount, contract.CurrentNTNAmount, "contract not updated")
	newNewtonBalance := new(big.Int).Add(newtonBalance, unbondingAmount)
	newtonBalance, _, err = r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.Equal(r.t, newNewtonBalance, newtonBalance, "vesting contract balance mismatch")
}
