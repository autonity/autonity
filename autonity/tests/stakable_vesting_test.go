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

type StakingRequest struct {
	amount      *big.Int
	contractID  *big.Int
	validator   common.Address
	expectedErr string
	bond        bool
}

type Reward struct {
	rewardATN *big.Int
	rewardNTN *big.Int
}

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
	maxBondGas, _, err := r.autonity.MaxBondAppliedGas(nil)
	require.NoError(r.t, err)
	maxRewardsDistributionGas, _, err := r.autonity.MaxRewardsDistributionGas(nil)
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
		bondingAmount = big.NewInt(contractTotalAmount / iteration)
		for ; iteration > 0; iteration-- {
			checkGasForBonding(r, []StakingRequest{{bondingAmount, common.Big0, validator, "", true}}, bondingID)
			bondingID++
		}
	})

	r.run("single bond after a successful bonding", func(r *runner) {
		oldLiquidBalance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		bondingID := len(r.committee.validators) + 1
		bondingAmount = big.NewInt(contractTotalAmount)
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
		checkGasForBonding(r, []StakingRequest{{bondingAmount, common.Big0, validator, "", true}}, bondingID)
	})

	r.run("bonding rejected", func(r *runner) {
		bondingID := len(r.committee.validators) + 1
		bondingAmount = big.NewInt(contractTotalAmount)
		gasUsedDistribute, gasUsedBond := bondAndApply(
			r, user, []StakingRequest{{bondingAmount, common.Big0, validator, "", true}}, bondingID, bondingGas, true,
		)
		totalGasUsed := big.NewInt(int64(gasUsedDistribute + gasUsedBond[0]))
		fmt.Printf("total gas used %v\n", totalGasUsed)
		fmt.Printf("gas to notify bond %v\n", gasUsedBond[0])
		fmt.Printf("gas to notify rewards distribution %v\n", gasUsedDistribute)
		require.True(
			r.t,
			bondingGas.Cmp(new(big.Int).Mul(totalGasUsed, stakingGas)) >= 0,
			"need more gas to notify bonding operations",
		)
		require.True(
			r.t,
			maxBondGas.Cmp(big.NewInt(int64(gasUsedBond[0]))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxRewardsDistributionGas.Cmp(big.NewInt(int64(gasUsedDistribute))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
	})

	for _, validator := range r.committee.validators {
		r.NoError(
			r.autonity.Mint(operator, user, bondingAmount),
		)
		r.NoError(
			r.autonity.Bond(fromSender(user, nil), validator.NodeAddress, bondingAmount),
		)
	}
	r.waitNextEpoch()

	r.run("bond to multiple validators from single contract", func(r *runner) {
		contractID := common.Big0
		bondingAmount = new(big.Int).Div(big.NewInt(contractTotalAmount), big.NewInt(int64(len(r.committee.validators))))
		require.Equal(r.t, new(big.Int).Mul(bondingAmount, big.NewInt(int64(len(r.committee.validators)))), big.NewInt(contractTotalAmount))
		bondingID := len(r.committee.validators)*2 + 1
		requests := make([]StakingRequest, 0)

		for _, validator := range r.committee.validators {
			requests = append(requests, StakingRequest{bondingAmount, contractID, validator.NodeAddress, "", true})
		}

		checkGasForBonding(r, requests, bondingID)

	})

	r.run("bond to multiple validators from multiple contract", func(r *runner) {
		bondingAmount = new(big.Int).Div(big.NewInt(contractTotalAmount), big.NewInt(int64(len(r.committee.validators))))
		require.Equal(r.t, new(big.Int).Mul(bondingAmount, big.NewInt(int64(len(r.committee.validators)))), big.NewInt(contractTotalAmount))
		bondingID := len(r.committee.validators)*2 + 1
		requests := make([]StakingRequest, 0)

		for i := 0; i < contractCount; i++ {
			for _, validator := range r.committee.validators {
				requests = append(requests, StakingRequest{bondingAmount, big.NewInt(int64(i)), validator.NodeAddress, "", true})
			}
		}

		checkGasForBonding(r, requests, bondingID)
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
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
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
			checkGasForUnbonding(
				r, []StakingRequest{{unbondingAmount, common.Big0, validator, "", false}}, unbondingID,
			)
			unbondingID++
		}
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
		bondAndFinalize(r, beneficiary, []StakingRequest{{bondingAmount, contractID, validator, "", true}}, bondingGas)
	})

	// start contract for bonding for all the tests remaining
	r.waitSomeBlock(start + 1)

	r.run("cannot bond more than total", func(r *runner) {
		bondingAmount := big.NewInt(contractTotalAmount + 10)
		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{bondingAmount, contractID, validator, "execution reverted: not enough tokens", true}

		bondingAmount = big.NewInt(contractTotalAmount / 2)
		requests[1] = StakingRequest{bondingAmount, contractID, validator, "", true}

		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		bondingAmount = new(big.Int).Add(big.NewInt(10), remaining)
		requests[2] = StakingRequest{bondingAmount, contractID, validator, "execution reverted: not enough tokens", true}

		bondAndFinalize(r, beneficiary, requests, bondingGas)

		requests = make([]StakingRequest, 2)
		requests[0] = StakingRequest{bondingAmount, contractID, validator, "execution reverted: not enough tokens", true}
		requests[1] = StakingRequest{remaining, contractID, validator, "", true}

		bondAndFinalize(r, beneficiary, requests, bondingGas)
	})

	r.run("can release liquid tokens", func(r *runner) {
		bondingAmount := big.NewInt(contractTotalAmount)
		bondAndFinalize(r, user, []StakingRequest{{bondingAmount, contractID, validator, "", true}}, bondingGas)
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
		// TODO (tariq): complete
	})

	r.run("when bonded, release NTN first", func(r *runner) {
		liquidBalance, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.True(r.t, contractTotalAmount > 10, "cannot test")
		bondingAmount := big.NewInt(contractTotalAmount / 10)
		bondAndFinalize(r, user, []StakingRequest{{bondingAmount, contractID, validator, "", true}}, bondingGas)
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

	r.run("can release LNTN", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("can release LNTN from any validator", func(r *runner) {
		// TODO (tariq): complete
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
	users, validators, _ := setupContracts(r, contractCount, validatorCount, contractTotalAmount, start, cliff, end)

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

	r.run("can unbond", func(r *runner) {
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, bondingAmount, liquid, "liquid not minted properly")
		unbondAndRelease(r, beneficiary, []StakingRequest{{liquid, contractID, validator, "", false}}, unbondingGas)
	})

	r.run("cannot unbond more than total liquid", func(r *runner) {
		unbondingAmount := new(big.Int).Add(bondingAmount, big.NewInt(10))
		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{unbondingAmount, contractID, validator, "execution reverted: not enough unlocked liquid tokens", false}

		unbondingAmount = big.NewInt(10)
		requests[1] = StakingRequest{unbondingAmount, contractID, validator, "", false}

		remaining := new(big.Int).Sub(bondingAmount, unbondingAmount)
		require.True(r.t, remaining.Cmp(common.Big0) > 0, "cannot test if no liquid remains")

		unbondingAmount = new(big.Int).Add(remaining, big.NewInt(10))
		requests[2] = StakingRequest{unbondingAmount, contractID, validator, "execution reverted: not enough unlocked liquid tokens", false}
		unbondAndRelease(r, user, requests, unbondingGas)

		requests = make([]StakingRequest, 2)
		requests[0] = StakingRequest{unbondingAmount, contractID, validator, "execution reverted: not enough unlocked liquid tokens", false}
		requests[1] = StakingRequest{remaining, contractID, validator, "", false}
		unbondAndRelease(r, user, requests, unbondingGas)
	})

	r.run("cannot unbond if LNTN withdrawn", func(r *runner) {
		liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		validator1 := validators[1]
		liquid1, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator1)
		require.NoError(r.t, err)
		require.True(r.t, liquid1.Cmp(big.NewInt(10)) > 0, "cannot test")

		totalToRelease := liquid.Int64() + 10
		currentTime := r.waitSomeEpoch(totalToRelease + start + 1)
		totalToRelease = currentTime - 1 - start
		r.NoError(
			r.stakableVesting.ReleaseAllLNTN(fromSender(beneficiary, nil), contractID),
		)

		// LNTN will be released from the first validator in the list
		newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.t, err)
		require.True(r.t, newLiquid.Cmp(common.Big0) == 0, "liquid remains after releasing")

		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{liquid, contractID, validator, "execution reverted: not enough unlocked liquid tokens", false}

		// if more unlocked funds remain, then LNTN will be released from 2nd validator
		releasedFromValidator1 := totalToRelease - liquid.Int64()
		remainingLiquid := new(big.Int).Sub(liquid1, big.NewInt(releasedFromValidator1))
		requests[1] = StakingRequest{liquid1, contractID, validator1, "execution reverted: not enough unlocked liquid tokens", false}

		liquid1, _, err = r.stakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator1)
		require.NoError(r.t, err)
		require.Equal(r.t, remainingLiquid, liquid1, "liquid balance mismatch")

		requests[2] = StakingRequest{liquid1, contractID, validator1, "", false}
		unbondAndRelease(r, beneficiary, requests, unbondingGas)
	})

	r.run("track liquid when unbonding from multiple contracts to multiple validators", func(r *runner) {
		// TODO (tariq): complete
	})
}

// TODO (tariq): low priority
func TestStakingRevert(t *testing.T) {
	r := setup(t, nil)
	// TODO (tariq): complete setup

	r.run("fails to notify reward distribution", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("reject bonding request and notify rejection", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("reject bonding request but fails to notify", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("revert applied bonding", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("reject unbonding request and notify rejection", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("reject unbonding request but fails to notify", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("revert applied unbonding", func(r *runner) {
		// TODO (tariq): complete
	})

	r.run("revert released unbonding", func(r *runner) {
		// TODO (tariq): complete
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
	unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	require.NoError(r.t, err)

	// start contract to bond
	r.waitSomeBlock(start + 1)

	r.run("bond and claim reward", func(r *runner) {
		beneficiary := users[0]
		contractID := common.Big0
		validator := validators[0]
		liquidContract := liquidContracts[0]
		bondingAmount := big.NewInt(contractTotalAmount)
		r.NoError(
			r.stakableVesting.Bond(
				fromSender(beneficiary, bondingGas), contractID, validator, bondingAmount,
			),
		)
		r.waitNextEpoch()

		r.NoError(
			r.autonity.ReceiveATN(
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
		extraBonds := make([]StakingRequest, 5)
		extraBonds[0] = StakingRequest{bondingAmount, common.Big0, validators[0], "", true}
		extraBonds[1] = StakingRequest{bondingAmount, common.Big1, validators[0], "", true}
		extraBonds[2] = StakingRequest{bondingAmount, common.Big0, validators[1], "", true}
		extraBonds[3] = StakingRequest{bondingAmount, common.Big0, validators[0], "", true}
		// dummy
		extraBonds[4] = StakingRequest{common.Big0, common.Big0, validators[0], "", true}

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		for _, user := range users {
			for _, request := range extraBonds {

				if request.amount.Cmp(common.Big0) > 0 {
					r.NoError(
						r.stakableVesting.Bond(
							fromSender(user, bondingGas), request.contractID, request.validator, request.amount,
						),
					)
				}

				r.giveMeSomeMoney(r.autonity.address, reward)
				totalReward, oldRewardsFromValidator, oldUserRewards := getRewardsAfterOneEpoch(r, contractCount, liquidContracts, users, validators)
				r.waitNextEpoch()
				// request is not applied yet
				checkRewards(
					r, contractCount, totalStake, totalReward,
					liquidContracts, validators, users, validatorStakes,
					userStakes, oldRewardsFromValidator, oldUserRewards,
				)

				// request is applied, because checkRewards progress 1 epoch
				amount := request.amount
				validator := request.validator
				id := int(request.contractID.Int64())
				validatorStakes[validator].Add(validatorStakes[validator], amount)
				userStakes[user][id][validator].Add(userStakes[user][id][validator], amount)
				totalStake.Add(totalStake, amount)
			}
		}
	})

	// bond everything
	oldBondingAmount := bondingAmount
	bondingPerContract := new(big.Int).Mul(oldBondingAmount, big.NewInt(int64(len(validators))))
	remainingNTN := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingPerContract)
	bondingAmount = new(big.Int).Div(remainingNTN, big.NewInt(int64(len(validators))))
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
	bondingAmount.Add(bondingAmount, oldBondingAmount)

	r.waitNextEpoch()

	r.run("release liquid and track reward", func(r *runner) {
		r.waitSomeEpoch(end + 1)
		releaseAmount := big.NewInt(100)
		// unbonding request can be treated as release request
		releaseRequests := make([]StakingRequest, 5)
		releaseRequests[0] = StakingRequest{releaseAmount, common.Big0, validators[0], "", false}
		releaseRequests[1] = StakingRequest{releaseAmount, common.Big1, validators[0], "", false}
		releaseRequests[2] = StakingRequest{releaseAmount, common.Big0, validators[1], "", false}
		releaseRequests[3] = StakingRequest{releaseAmount, common.Big0, validators[0], "", false}
		// dummy
		releaseRequests[4] = StakingRequest{common.Big0, common.Big0, validators[0], "", false}

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		liquidContractsMap := make(map[common.Address]*Liquid)

		for i, liquidContract := range liquidContracts {
			liquidContractsMap[validators[i]] = liquidContract
		}

		for _, user := range users {

			userLiquidBalance := make(map[common.Address]*big.Int)
			for _, validator := range validators {
				userLiquidBalance[validator] = new(big.Int)
			}

			for _, request := range releaseRequests {

				// some epoch is passed and we are entitled to some reward,
				// but we don't know about it because we did not get notified
				// or we did not claim them or call unclaimedRewards
				r.giveMeSomeMoney(r.autonity.address, reward)
				totalReward, oldRewardsFromValidator, oldUserRewards := getRewardsAfterOneEpoch(r, contractCount, liquidContracts, users, validators)
				r.waitNextEpoch()

				// we release some LNTN and it is applied immediately
				// if unlocked, it is transferred immediately
				// but for reward calculation, it will be applied later
				if request.amount.Cmp(common.Big0) > 0 {
					r.NoError(
						r.stakableVesting.ReleaseLNTN(
							fromSender(user, nil),
							request.contractID,
							request.validator,
							request.amount,
						),
					)
				}

				amount := request.amount
				validator := request.validator
				userLiquidBalance[validator].Add(userLiquidBalance[validator], amount)
				balance, _, err := liquidContractsMap[validator].BalanceOf(nil, user)
				require.NoError(r.t, err)
				require.Equal(r.t, userLiquidBalance[validator], balance, "liquid not transferred")

				checkRewards(
					r, contractCount, totalStake, totalReward,
					liquidContracts, validators, users, validatorStakes,
					userStakes, oldRewardsFromValidator, oldUserRewards,
				)

				// for next reward
				id := int(request.contractID.Int64())
				validatorStakes[validator].Sub(validatorStakes[validator], amount)
				userStakes[user][id][validator].Sub(userStakes[user][id][validator], amount)
			}
		}
	})

	r.run("unbond in different epoch and track reward", func(r *runner) {
		unbondingAmount := big.NewInt(100)
		extraUnbonds := make([]StakingRequest, 5)
		extraUnbonds[0] = StakingRequest{unbondingAmount, common.Big0, validators[0], "", false}
		extraUnbonds[1] = StakingRequest{unbondingAmount, common.Big1, validators[0], "", false}
		extraUnbonds[2] = StakingRequest{unbondingAmount, common.Big0, validators[1], "", false}
		extraUnbonds[3] = StakingRequest{unbondingAmount, common.Big0, validators[0], "", false}
		// dummy
		extraUnbonds[4] = StakingRequest{common.Big0, common.Big0, validators[0], "", false}

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		for _, user := range users {
			for _, request := range extraUnbonds {

				if request.amount.Cmp(common.Big0) > 0 {
					r.NoError(
						r.stakableVesting.Unbond(
							fromSender(user, unbondingGas), request.contractID, request.validator, request.amount,
						),
					)
				}

				r.giveMeSomeMoney(r.autonity.address, reward)
				totalReward, oldRewardsFromValidator, oldUserRewards := getRewardsAfterOneEpoch(r, contractCount, liquidContracts, users, validators)
				r.waitNextEpoch()
				// request is not applied yet
				checkRewards(
					r, contractCount, totalStake, totalReward,
					liquidContracts, validators, users, validatorStakes,
					userStakes, oldRewardsFromValidator, oldUserRewards,
				)

				// request is applied, because checkRewards progress 1 epoch
				amount := request.amount
				validator := request.validator
				id := int(request.contractID.Int64())
				validatorStakes[validator].Sub(validatorStakes[validator], amount)
				userStakes[user][id][validator].Sub(userStakes[user][id][validator], amount)
				totalStake.Sub(totalStake, amount)
			}
		}
	})
}

func TestChangeContractBeneficiary(t *testing.T) {
	r := setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0

	r.run("beneficiary changes", func(r *runner) {
		_, _, err := r.stakableVesting.GetContract(nil, user, contractID)
		require.NoError(r.t, err)
		newUser := common.HexToAddress("0x88")
		_, _, err = r.stakableVesting.GetContract(nil, newUser, contractID)
		require.Equal(r.t, "execution reverted: invalid contract id", err.Error())
		r.stakableVesting.ChangeContractBeneficiary(operator, user, contractID, newUser)
		_, _, err = r.stakableVesting.GetContract(nil, newUser, contractID)
		require.NoError(r.t, err)
		_, _, err = r.stakableVesting.GetContract(nil, user, contractID)
		require.Equal(r.t, "execution reverted: invalid contract id", err.Error())
	})
}

func TestContractUpdateWhenSlashed(t *testing.T) {
	r := setup(t, nil)
	// TODO (tariq): complete setup

	r.run("contract total value update when bonded validator slashed", func(r *runner) {
		// TODO (tariq): complete
	})
}

func TestAccessRestriction(t *testing.T) {
	r := setup(t, nil)

	r.run("only operator can create contract", func(r *runner) {
		amount := big.NewInt(1000)
		start := new(big.Int).Add(big.NewInt(100), r.evm.Context.Time)
		cliff := new(big.Int).Add(start, big.NewInt(100))
		end := new(big.Int).Add(start, amount)
		_, err := r.stakableVesting.NewContract(
			fromSender(user, nil),
			user,
			amount,
			start,
			cliff,
			end,
		)
		require.Equal(r.t, "execution reverted: caller is not the operator", err.Error())
	})

	r.run("only operator can set gas cost", func(r *runner) {
		_, err := r.stakableVesting.SetRequiredGasBond(
			fromSender(user, nil),
			big.NewInt(100),
		)
		require.Equal(r.t, "execution reverted: caller is not the operator", err.Error())

		_, err = r.stakableVesting.SetRequiredGasUnbond(
			fromSender(user, nil),
			big.NewInt(100),
		)
		require.Equal(r.t, "execution reverted: caller is not the operator", err.Error())
	})

	var contractTotalAmount int64 = 1000
	start := r.evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0
	validator := r.committee.validators[0].NodeAddress

	r.run("cannot request bonding or unbonding without enough gas", func(r *runner) {
		bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
		require.NoError(r.t, err)
		balance := big.NewInt(1000_000_000_000_000_000)
		r.giveMeSomeMoney(user, balance)
		bondingAmount := big.NewInt(100)
		_, err = r.stakableVesting.Bond(
			fromSender(user, new(big.Int).Sub(bondingGas, common.Big1)),
			contractID,
			validator,
			bondingAmount,
		)
		require.Equal(r.t, "execution reverted: not enough gas given for notification on bonding", err.Error())

		r.NoError(
			r.stakableVesting.Bond(
				fromSender(user, bondingGas),
				contractID,
				validator,
				bondingAmount,
			),
		)
		r.waitNextEpoch()

		unbondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
		require.NoError(r.t, err)
		_, err = r.stakableVesting.Unbond(
			fromSender(user, new(big.Int).Sub(unbondingGas, common.Big1)),
			contractID,
			validator,
			bondingAmount,
		)

		require.Equal(r.t, "execution reverted: not enough gas given for notification on unbonding", err.Error())

		r.NoError(
			r.stakableVesting.Unbond(
				fromSender(user, unbondingGas),
				contractID,
				validator,
				bondingAmount,
			),
		)

	})

	r.run("only operator can change contract beneficiary", func(r *runner) {
		newUser := common.HexToAddress("0x88")
		require.NotEqual(r.t, user, newUser)
		_, err := r.stakableVesting.ChangeContractBeneficiary(
			fromSender(user, nil),
			user,
			contractID,
			newUser,
		)
		require.Equal(r.t, "execution reverted: caller is not the operator", err.Error())

		_, err = r.stakableVesting.ChangeContractBeneficiary(
			fromSender(newUser, nil),
			user,
			contractID,
			newUser,
		)
		require.Equal(r.t, "execution reverted: caller is not the operator", err.Error())
	})

	r.run("only autonity can notify staking operations", func(r *runner) {

		_, err := r.stakableVesting.RewardsDistributed(
			fromSender(user, nil),
			[]common.Address{},
		)
		require.Equal(r.t, "execution reverted: function restricted to Autonity contract", err.Error())

		_, err = r.stakableVesting.BondingApplied(
			fromSender(user, nil),
			common.Big0,
			validator,
			common.Big1,
			true,
			true,
		)
		require.Equal(r.t, "execution reverted: function restricted to Autonity contract", err.Error())

		_, err = r.stakableVesting.UnbondingApplied(
			fromSender(user, nil),
			common.Big0,
			validator,
			true,
		)
		require.Equal(r.t, "execution reverted: function restricted to Autonity contract", err.Error())

		_, err = r.stakableVesting.UnbondingReleased(
			fromSender(user, nil),
			common.Big0,
			common.Big1,
			true,
		)
		require.Equal(r.t, "execution reverted: function restricted to Autonity contract", err.Error())
	})
}

func initialStakes(
	r *runner,
	contractCount int,
	liquidContracts []*Liquid,
	users, validators []common.Address,
) (
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[int]map[common.Address]*big.Int,
	totalStake *big.Int,
) {

	totalStake = new(big.Int)

	validatorStakes = make(map[common.Address]*big.Int)
	for i, validator := range validators {
		liquidContract := liquidContracts[i]
		balance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		validatorStakes[validator] = balance
	}

	userStakes = make(map[common.Address]map[int]map[common.Address]*big.Int)
	for _, user := range users {
		userStakes[user] = make(map[int]map[common.Address]*big.Int)
		for i := 0; i < contractCount; i++ {
			userStakes[user][i] = make(map[common.Address]*big.Int)
			for _, validator := range validators {
				balance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.t, err)
				userStakes[user][i][validator] = balance
				totalStake.Add(totalStake, balance)
			}
		}
	}
	return validatorStakes, userStakes, totalStake
}

func getRewardsAfterOneEpoch(
	r *runner,
	contractCount int,
	liquidContracts []*Liquid,
	users, validators []common.Address,
) (
	currentReward Reward,
	oldRewardsFromValidator map[common.Address]Reward,
	oldUserRewards map[common.Address]map[int]map[common.Address]Reward,
) {
	oldRewardsFromValidator = make(map[common.Address]Reward)
	for i, validator := range validators {
		liquidContract := liquidContracts[i]
		unclaimedReward, _, err := liquidContract.UnclaimedRewards(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		oldRewardsFromValidator[validator] = Reward{unclaimedReward.UnclaimedATN, unclaimedReward.UnclaimedNTN}
	}

	oldUserRewards = make(map[common.Address]map[int]map[common.Address]Reward)
	for _, user := range users {
		oldUserRewards[user] = make(map[int]map[common.Address]Reward)
		for i := 0; i < contractCount; i++ {
			oldUserRewards[user][i] = make(map[common.Address]Reward)
			for _, validator := range validators {
				unclaimedReward, _, err := r.stakableVesting.UnclaimedRewards(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.t, err)
				oldUserRewards[user][i][validator] = Reward{unclaimedReward.AtnFee, unclaimedReward.NtnFee}
			}
		}
	}

	// get supply and inflationReserve to calculate inflation reward
	supply, _, err := r.autonity.TotalSupply(nil)
	require.NoError(r.t, err)
	inflationReserve, _, err := r.autonity.InflationReserve(nil)
	require.NoError(r.t, err)
	epochPeriod, _, err := r.autonity.GetEpochPeriod(nil)
	require.NoError(r.t, err)

	// get inflation reward
	lastEpochTime, _, err := r.autonity.LastEpochTime(nil)
	require.NoError(r.t, err)
	currentEpochTime := new(big.Int).Add(lastEpochTime, epochPeriod)
	currentReward.rewardNTN, _, err = r.inflationController.CalculateSupplyDelta(nil, supply, inflationReserve, lastEpochTime, currentEpochTime)
	require.NoError(r.t, err)

	// get atn reward
	currentReward.rewardATN = r.getBalanceOf(r.autonity.address)
	return currentReward, oldRewardsFromValidator, oldUserRewards
}

func checkRewards(
	r *runner,
	contractCount int,
	totalStake *big.Int,
	totalReward Reward,
	liquidContracts []*Liquid,
	validators, users []common.Address,
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[int]map[common.Address]*big.Int,
	oldRewardsFromValidator map[common.Address]Reward,
	oldUserRewards map[common.Address]map[int]map[common.Address]Reward,
) {

	currentRewards := make(map[common.Address]Reward)
	for i, validator := range validators {
		validatorTotalRewardATN := new(big.Int).Mul(validatorStakes[validator], totalReward.rewardATN)
		validatorTotalRewardNTN := new(big.Int).Mul(validatorStakes[validator], totalReward.rewardNTN)

		if totalStake.Cmp(common.Big0) != 0 {
			validatorTotalRewardATN = validatorTotalRewardATN.Div(validatorTotalRewardATN, totalStake)
			validatorTotalRewardNTN = validatorTotalRewardNTN.Div(validatorTotalRewardNTN, totalStake)
		}

		liquidContract := liquidContracts[i]
		unclaimedReward, _, err := liquidContract.UnclaimedRewards(nil, r.stakableVesting.address)
		require.NoError(r.t, err)

		diff := new(big.Int).Sub(
			new(big.Int).Add(validatorTotalRewardATN, oldRewardsFromValidator[validator].rewardATN),
			unclaimedReward.UnclaimedATN,
		)
		diff.Abs(diff)
		// difference should be less than or equal to 1 wei
		require.True(
			r.t,
			diff.Cmp(common.Big1) <= 0,
			"unclaimed atn reward not updated in liquid contract",
		)

		diff = new(big.Int).Sub(
			new(big.Int).Add(validatorTotalRewardNTN, oldRewardsFromValidator[validator].rewardNTN),
			unclaimedReward.UnclaimedNTN,
		)
		diff.Abs(diff)
		// difference should be less than or equal to 1 wei
		require.True(
			r.t,
			diff.Cmp(common.Big1) <= 0,
			"unclaimed ntn reward not updated in liquid contract",
		)
		currentRewards[validator] = Reward{
			new(big.Int).Sub(unclaimedReward.UnclaimedATN, oldRewardsFromValidator[validator].rewardATN),
			new(big.Int).Sub(unclaimedReward.UnclaimedNTN, oldRewardsFromValidator[validator].rewardNTN),
		}
	}

	for _, user := range users {
		userRewardATN := new(big.Int)
		userRewardNTN := new(big.Int)
		// The following loops is equivalent to: (user_all_stake_to_all_validator) * totalReward / totalStake
		// But StakableVesting contract handles reward for each validator separately, so there can be some funds lost due to
		// integer division in solidity. So we simulate the calculation with the for loop instead
		for i := 0; i < contractCount; i++ {
			unclaimedRewardForContractATN := new(big.Int)
			unclaimedRewardForContractNTN := new(big.Int)
			for _, validator := range validators {
				calculatedRewardATN := new(big.Int).Mul(userStakes[user][i][validator], currentRewards[validator].rewardATN)
				calculatedRewardNTN := new(big.Int).Mul(userStakes[user][i][validator], currentRewards[validator].rewardNTN)

				if validatorStakes[validator].Cmp(common.Big0) != 0 {
					calculatedRewardATN.Div(calculatedRewardATN, validatorStakes[validator])
					calculatedRewardNTN.Div(calculatedRewardNTN, validatorStakes[validator])
				}
				calculatedRewardATN.Add(calculatedRewardATN, oldUserRewards[user][i][validator].rewardATN)

				calculatedRewardNTN.Add(calculatedRewardNTN, oldUserRewards[user][i][validator].rewardNTN)

				unclaimedReward, _, err := r.stakableVesting.UnclaimedRewards(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.t, err)

				diff := new(big.Int).Sub(calculatedRewardATN, unclaimedReward.AtnFee)
				diff.Abs(diff)
				// difference should be less than or equal to 1 wei
				require.True(
					r.t,
					diff.Cmp(common.Big1) <= 0,
					"atn reward calculation mismatch",
				)

				diff = new(big.Int).Sub(calculatedRewardNTN, unclaimedReward.NtnFee)
				diff.Abs(diff)
				// difference should be less than or equal to 1 wei
				require.True(
					r.t,
					diff.Cmp(common.Big1) <= 0,
					"ntn reward calculation mismatch",
				)
				unclaimedRewardForContractATN.Add(unclaimedRewardForContractATN, unclaimedReward.AtnFee)
				unclaimedRewardForContractNTN.Add(unclaimedRewardForContractNTN, unclaimedReward.NtnFee)
			}

			unclaimedReward, _, err := r.stakableVesting.UnclaimedRewards1(nil, user, big.NewInt(int64(i)))
			require.NoError(r.t, err)
			require.Equal(r.t, unclaimedRewardForContractATN, unclaimedReward.AtnFee)
			require.Equal(r.t, unclaimedRewardForContractNTN, unclaimedReward.NtnFee)

			userRewardATN.Add(userRewardATN, unclaimedReward.AtnFee)
			userRewardNTN.Add(userRewardNTN, unclaimedReward.NtnFee)
		}

		unclaimedReward, _, err := r.stakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.t, err)

		require.Equal(
			r.t,
			userRewardATN,
			unclaimedReward.AtnTotalFee,
			"unclaimed atn reward mismatch",
		)

		require.Equal(
			r.t,
			userRewardNTN,
			unclaimedReward.NtnTotalFee,
			"unclaimed ntn reward mismatch",
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
	startBig := big.NewInt(startTime)
	cliffBig := big.NewInt(cliffTime)
	endBig := big.NewInt(endTime)
	r.NoError(
		r.stakableVesting.NewContract(
			operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
			new(big.Int).Sub(cliffBig, startBig), new(big.Int).Sub(endBig, startBig),
		),
	)
}

func fromSender(sender common.Address, value *big.Int) *runOptions {
	return &runOptions{origin: sender, value: value}
}

func checkGasForUnbonding(r *runner, requests []StakingRequest, unbondingID int) {
	unbondingGas, _, err := r.stakableVesting.RequiredUnbondingGasCost(nil)
	require.NoError(r.t, err)
	stakingGas, _, err := r.autonity.StakingGasPrice(nil)
	require.NoError(r.t, err)
	maxUnbondGas, _, err := r.autonity.MaxUnbondAppliedGas(nil)
	require.NoError(r.t, err)
	maxUnbondReleaseGas, _, err := r.autonity.MaxUnbondReleasedGas(nil)
	require.NoError(r.t, err)
	maxRewardsDistributionGas, _, err := r.autonity.MaxRewardsDistributionGas(nil)
	require.NoError(r.t, err)
	gasUsedDistribute, gasUsedUnbond, gasUsedRelease := unbondAndApply(
		r, user, requests, unbondingID, unbondingGas, true,
	)

	require.Equal(r.t, len(requests), len(gasUsedUnbond))
	require.Equal(r.t, len(requests), len(gasUsedRelease))

	totalGasUsed := big.NewInt(int64(gasUsedDistribute))
	avgGasUsedDistribute := gasUsedDistribute / uint64(len(requests))
	for i, gasUsed := range gasUsedUnbond {
		gasUsedForRelease := gasUsedRelease[i]
		fmt.Printf("gas to notify unbond %v\n", gasUsed)
		fmt.Printf("gas to notify unbond release %v\n", gasUsedForRelease)
		fmt.Printf("gas to notify rewards distribution %v\n", avgGasUsedDistribute)
		avgGasUsed := big.NewInt(int64(gasUsed + gasUsedForRelease + avgGasUsedDistribute))
		require.True(
			r.t,
			unbondingGas.Cmp(new(big.Int).Mul(avgGasUsed, stakingGas)) >= 0,
			"need more avg gas to notify unbonding",
		)
		require.True(
			r.t,
			maxUnbondGas.Cmp(big.NewInt(int64(gasUsed))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxUnbondReleaseGas.Cmp(big.NewInt(int64(gasUsedForRelease))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxRewardsDistributionGas.Cmp(big.NewInt(int64(avgGasUsedDistribute))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		totalGasUsed.Add(totalGasUsed, big.NewInt(int64(gasUsed)))
		totalGasUsed.Add(totalGasUsed, big.NewInt(int64(gasUsedForRelease)))
	}
	fmt.Printf("total gas used %v\n", totalGasUsed)
	require.True(
		r.t,
		new(big.Int).Mul(unbondingGas, big.NewInt(int64(len(requests)))).Cmp(new(big.Int).Mul(totalGasUsed, stakingGas)) >= 0,
		"need more gas to notify bonding operations",
	)
	require.True(
		r.t,
		new(big.Int).Mul(maxRewardsDistributionGas, big.NewInt(int64(len(requests)))).Cmp(big.NewInt(int64(gasUsedDistribute))) >= 0,
		"gas usage exceeds autonity allowed gas",
	)

	gasUsedDistribute, gasUsedUnbond, gasUsedRelease = unbondAndApply(
		r, user, requests, unbondingID, unbondingGas, false,
	)

	require.Equal(r.t, len(requests), len(gasUsedUnbond))
	require.Equal(r.t, len(requests), len(gasUsedRelease))

	totalGasUsed = big.NewInt(int64(gasUsedDistribute))
	avgGasUsedDistribute = gasUsedDistribute / uint64(len(requests))
	for i, gasUsed := range gasUsedUnbond {
		gasUsedForRelease := gasUsedRelease[i]
		fmt.Printf("gas to notify unbond %v\n", gasUsed)
		fmt.Printf("gas to notify unbond release %v\n", gasUsedForRelease)
		fmt.Printf("gas to notify rewards distribution %v\n", avgGasUsedDistribute)
		avgGasUsed := big.NewInt(int64(gasUsed + gasUsedForRelease + avgGasUsedDistribute))
		require.True(
			r.t,
			unbondingGas.Cmp(new(big.Int).Mul(avgGasUsed, stakingGas)) >= 0,
			"need more avg gas to notify unbonding",
		)
		require.True(
			r.t,
			maxUnbondGas.Cmp(big.NewInt(int64(gasUsed))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxUnbondReleaseGas.Cmp(big.NewInt(int64(gasUsedForRelease))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxRewardsDistributionGas.Cmp(big.NewInt(int64(avgGasUsedDistribute))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		totalGasUsed.Add(totalGasUsed, big.NewInt(int64(gasUsed)))
		totalGasUsed.Add(totalGasUsed, big.NewInt(int64(gasUsedForRelease)))
	}
	fmt.Printf("total gas used %v\n", totalGasUsed)
	require.True(
		r.t,
		new(big.Int).Mul(unbondingGas, big.NewInt(int64(len(requests)))).Cmp(new(big.Int).Mul(totalGasUsed, stakingGas)) >= 0,
		"need more gas to notify bonding operations",
	)
	require.True(
		r.t,
		new(big.Int).Mul(maxRewardsDistributionGas, big.NewInt(int64(len(requests)))).Cmp(big.NewInt(int64(gasUsedDistribute))) >= 0,
		"gas usage exceeds autonity allowed gas",
	)
}

func checkGasForBonding(r *runner, requests []StakingRequest, bondingID int) {
	bondingGas, _, err := r.stakableVesting.RequiredBondingGasCost(nil)
	require.NoError(r.t, err)
	stakingGas, _, err := r.autonity.StakingGasPrice(nil)
	require.NoError(r.t, err)
	maxBondGas, _, err := r.autonity.MaxBondAppliedGas(nil)
	require.NoError(r.t, err)
	maxRewardsDistributionGas, _, err := r.autonity.MaxRewardsDistributionGas(nil)
	require.NoError(r.t, err)
	gasUsedDistribute, gasUsedBond := bondAndApply(
		r, user, requests, bondingID, bondingGas, true,
	)
	require.Equal(r.t, len(requests), len(gasUsedBond))

	totalGasUsed := big.NewInt(int64(gasUsedDistribute))
	avgGasUsedDistribute := gasUsedDistribute / uint64(len(requests))
	for _, gasUsed := range gasUsedBond {
		fmt.Printf("gas to notify bond %v\n", gasUsed)
		fmt.Printf("gas to notify rewards distribution %v\n", avgGasUsedDistribute)
		avgGasUsed := big.NewInt(int64(gasUsed + avgGasUsedDistribute))
		require.True(
			r.t,
			bondingGas.Cmp(new(big.Int).Mul(avgGasUsed, stakingGas)) >= 0,
			"need more avg gas to notify bonding",
		)
		require.True(
			r.t,
			maxBondGas.Cmp(big.NewInt(int64(gasUsed))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxRewardsDistributionGas.Cmp(big.NewInt(int64(avgGasUsedDistribute))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		totalGasUsed.Add(totalGasUsed, big.NewInt(int64(gasUsed)))
	}
	fmt.Printf("total gas used %v\n", totalGasUsed)
	require.True(
		r.t,
		new(big.Int).Mul(bondingGas, big.NewInt(int64(len(requests)))).Cmp(new(big.Int).Mul(totalGasUsed, stakingGas)) >= 0,
		"need more gas to notify bonding operations",
	)
	require.True(
		r.t,
		new(big.Int).Mul(maxRewardsDistributionGas, big.NewInt(int64(len(requests)))).Cmp(big.NewInt(int64(gasUsedDistribute))) >= 0,
		"gas usage exceeds autonity allowed gas",
	)

	gasUsedDistribute, gasUsedBond = bondAndApply(
		r, user, requests, bondingID, bondingGas, false,
	)
	require.Equal(r.t, len(requests), len(gasUsedBond))

	totalGasUsed = big.NewInt(int64(gasUsedDistribute))
	avgGasUsedDistribute = gasUsedDistribute / uint64(len(requests))
	for _, gasUsed := range gasUsedBond {
		fmt.Printf("gas to notify bond %v\n", gasUsed)
		fmt.Printf("gas to notify rewards distribution %v\n", avgGasUsedDistribute)
		avgGasUsed := big.NewInt(int64(gasUsed + avgGasUsedDistribute))
		require.True(
			r.t,
			bondingGas.Cmp(new(big.Int).Mul(avgGasUsed, stakingGas)) >= 0,
			"need more avg gas to notify bonding",
		)
		require.True(
			r.t,
			maxBondGas.Cmp(big.NewInt(int64(gasUsed))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		require.True(
			r.t,
			maxRewardsDistributionGas.Cmp(big.NewInt(int64(avgGasUsedDistribute))) >= 0,
			"gas usage exceeds autonity allowed gas",
		)
		totalGasUsed.Add(totalGasUsed, big.NewInt(int64(gasUsed)))
	}
	fmt.Printf("total gas used %v\n", totalGasUsed)
	require.True(
		r.t,
		new(big.Int).Mul(bondingGas, big.NewInt(int64(len(requests)))).Cmp(new(big.Int).Mul(totalGasUsed, stakingGas)) >= 0,
		"need more gas to notify bonding operations",
	)
	require.True(
		r.t,
		new(big.Int).Mul(maxRewardsDistributionGas, big.NewInt(int64(len(requests)))).Cmp(big.NewInt(int64(gasUsedDistribute))) >= 0,
		"gas usage exceeds autonity allowed gas",
	)
}

func bondAndApply(
	r *runner, user common.Address, bondingRequests []StakingRequest, bondingID int, bondingGas *big.Int, rejected bool,
) (uint64, []uint64) {

	liquidContracts := make(map[common.Address]*Liquid)

	for i, validator := range r.committee.validators {
		for _, request := range bondingRequests {
			if request.validator == validator.NodeAddress {
				liquidContracts[request.validator] = r.committee.liquidContracts[i]
				break
			}
		}
	}

	liquid := make(map[int64]map[common.Address]*big.Int)
	for _, request := range bondingRequests {
		liquid[request.contractID.Int64()] = make(map[common.Address]*big.Int)
	}
	validatorExist := make(map[common.Address]bool)

	for _, request := range bondingRequests {
		validator := request.validator
		id := request.contractID.Int64()
		balance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		liquid[id][validator] = balance
		r.NoError(
			r.stakableVesting.Bond(fromSender(user, bondingGas), request.contractID, validator, request.amount),
		)
		validatorExist[validator] = true
	}

	bondedValidators := make([]common.Address, 0)
	for key := range validatorExist {
		liquidContract := liquidContracts[key]
		bondedValidators = append(bondedValidators, key)
		r.giveMeSomeMoney(r.autonity.address, reward)
		r.autonity.Mint(operator, liquidContract.address, reward)
		r.NoError(
			liquidContract.Redistribute(fromSender(r.autonity.address, reward), reward),
		)
	}

	gasUsedDistribute := r.NoError(
		r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators),
	)

	if rejected == false {
		for _, request := range bondingRequests {
			validator := request.validator
			id := request.contractID.Int64()
			liquidContract := liquidContracts[validator]
			r.NoError(
				liquidContract.Mint(fromAutonity, r.stakableVesting.address, request.amount),
			)
			liquid[id][validator].Add(liquid[id][validator], request.amount)
		}
	}

	gasUsedBond := make([]uint64, 0)
	for i, request := range bondingRequests {
		validator := request.validator
		id := request.contractID.Int64()
		curBondingID := new(big.Int).Add(big.NewInt(int64(bondingID)), big.NewInt(int64(i)))
		gasUsed := r.NoError(
			r.stakableVesting.BondingApplied(
				fromAutonity, curBondingID, validator, request.amount, false, rejected,
			),
		)
		newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, liquid[id][validator], newLiquid)
		gasUsedBond = append(gasUsedBond, gasUsed)
	}

	// for _, request := range bondingRequests {
	// 	validator := request.validator
	// 	id := request.contractID.Int64()
	// 	newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
	// 	require.NoError(r.t, err)
	// 	require.Equal(r.t, liquid[id][validator], newLiquid)
	// }
	return gasUsedDistribute, gasUsedBond
}

func unbondAndApply(
	r *runner, user common.Address, unbondingRequests []StakingRequest, unbondingID int, unbondingGas *big.Int, rejected bool,
) (uint64, []uint64, []uint64) {

	liquidContracts := make(map[common.Address]*Liquid)

	for i, validator := range r.committee.validators {
		for _, request := range unbondingRequests {
			if request.validator == validator.NodeAddress {
				liquidContracts[request.validator] = r.committee.liquidContracts[i]
				break
			}
		}
	}

	liquid := make(map[int64]map[common.Address]*big.Int)

	for _, request := range unbondingRequests {
		liquid[request.contractID.Int64()] = make(map[common.Address]*big.Int)
	}

	validatorExist := make(map[common.Address]bool)
	for _, request := range unbondingRequests {
		id := request.contractID.Int64()
		validator := request.validator
		balance, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		liquid[id][validator] = balance
		r.NoError(
			r.stakableVesting.Unbond(fromSender(user, unbondingGas), request.contractID, validator, request.amount),
		)
		validatorExist[validator] = true
	}
	// liquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validatorAddress)
	// require.NoError(r.t, err)
	// r.NoError(
	// 	r.stakableVesting.Unbond(fromSender(user, unbondingGas), contractID, validatorAddress, unbondingAmount),
	// )

	bondedValidators := make([]common.Address, 0)
	for key := range validatorExist {
		liquidContract := liquidContracts[key]
		r.giveMeSomeMoney(r.autonity.address, reward)
		r.NoError(
			r.autonity.Mint(operator, liquidContract.address, reward),
		)
		r.NoError(
			liquidContract.Redistribute(fromSender(r.autonity.address, reward), reward),
		)
		bondedValidators = append(bondedValidators, key)
	}

	// bondedValidators[0] = validatorAddress
	gasUsedDistribute := r.NoError(
		r.stakableVesting.RewardsDistributed(fromAutonity, bondedValidators),
	)

	if rejected == false {
		// r.NoError(
		// 	liquidContract.Unlock(fromAutonity, r.stakableVesting.address, unbondingAmount),
		// )
		// r.NoError(
		// 	liquidContract.Burn(fromAutonity, r.stakableVesting.address, unbondingAmount),
		// )
		// liquid = liquid.Sub(liquid, unbondingAmount)

		for _, request := range unbondingRequests {
			id := request.contractID.Int64()
			validator := request.validator
			liquidContract := liquidContracts[validator]
			r.NoError(
				liquidContract.Unlock(fromAutonity, r.stakableVesting.address, request.amount),
			)
			r.NoError(
				liquidContract.Burn(fromAutonity, r.stakableVesting.address, request.amount),
			)
			liquid[id][validator].Sub(liquid[id][validator], request.amount)
		}
	}

	gasUsedUnbond := make([]uint64, 0)
	for i, request := range unbondingRequests {
		// gasUsedUnbond := r.NoError(
		// 	r.stakableVesting.UnbondingApplied(fromAutonity, big.NewInt(int64(unbondingID)), validatorAddress, rejected),
		// )
		validator := request.validator
		id := request.contractID.Int64()
		curBondingID := new(big.Int).Add(big.NewInt(int64(unbondingID)), big.NewInt(int64(i)))
		gasUsed := r.NoError(
			r.stakableVesting.UnbondingApplied(fromAutonity, curBondingID, validator, rejected),
		)

		newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		require.Equal(r.t, liquid[id][validator], newLiquid)
		gasUsedUnbond = append(gasUsedUnbond, gasUsed)
	}

	// newLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, contractID, validatorAddress)
	// require.NoError(r.t, err)
	// require.Equal(r.t, liquid, newLiquid)

	gasUsedRelease := make([]uint64, 0)
	for i, request := range unbondingRequests {
		curBondingID := new(big.Int).Add(big.NewInt(int64(unbondingID)), big.NewInt(int64(i)))
		gasUsed := r.NoError(
			r.stakableVesting.UnbondingReleased(fromAutonity, curBondingID, request.amount, rejected),
		)
		gasUsedRelease = append(gasUsedRelease, gasUsed)
	}
	return gasUsedDistribute, gasUsedUnbond, gasUsedRelease
}

func checkReleaseAllNTN(r *runner, user common.Address, contractID, unlockAmount *big.Int) {
	contract, _, err := r.stakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.t, err)
	contractNTN := contract.CurrentNTNAmount
	withdrawn := contract.WithdrawnValue
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
	require.True(
		r.t,
		new(big.Int).Sub(contractNTN, unlockAmount).Cmp(contract.CurrentNTNAmount) == 0,
		"contract NTN not updated properly",
	)
	require.True(
		r.t,
		new(big.Int).Add(withdrawn, unlockAmount).Cmp(contract.WithdrawnValue) == 0,
		"contract WithdrawnValue not updated properly",
	)
}

func bondAndFinalize(
	r *runner, user common.Address, bondingRequests []StakingRequest, bondingGas *big.Int,
) {
	liquidContracts := make(map[common.Address]*Liquid)
	liquidOfVestingContract := make(map[common.Address]*big.Int)
	liquidOfUser := make(map[common.Address]map[int64]*big.Int)

	for i, validator := range r.committee.validators {
		for _, request := range bondingRequests {
			if request.validator == validator.NodeAddress {
				liquidContract := r.committee.liquidContracts[i]
				liquidContracts[request.validator] = liquidContract

				balance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
				require.NoError(r.t, err)
				liquidOfVestingContract[request.validator] = balance

				liquidOfUser[request.validator] = make(map[int64]*big.Int)
				break
			}
		}
	}

	for _, request := range bondingRequests {
		userLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, request.validator)
		require.NoError(r.t, err)
		liquidOfUser[request.validator][request.contractID.Int64()] = userLiquid
	}

	newtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)

	for _, request := range bondingRequests {
		contract, _, err := r.stakableVesting.GetContract(nil, user, request.contractID)
		require.NoError(r.t, err)
		contractNTN := contract.CurrentNTNAmount

		_, err = r.stakableVesting.Bond(
			fromSender(user, bondingGas),
			request.contractID,
			request.validator,
			request.amount,
		)

		if request.expectedErr == "" {
			require.NoError(r.t, err)
			validator := request.validator
			id := request.contractID.Int64()
			liquidOfVestingContract[validator].Add(liquidOfVestingContract[validator], request.amount)
			liquidOfUser[validator][id].Add(liquidOfUser[validator][id], request.amount)

			contract, _, err = r.stakableVesting.GetContract(nil, user, request.contractID)
			require.NoError(r.t, err)
			remaining := new(big.Int).Sub(contractNTN, request.amount)
			require.True(r.t, remaining.Cmp(contract.CurrentNTNAmount) == 0, "contract not updated properly")

			newtonBalance.Sub(newtonBalance, request.amount)
		} else {
			require.Equal(r.t, request.expectedErr, err.Error())
		}
	}

	// let bonding apply
	r.waitNextEpoch()

	for _, request := range bondingRequests {
		validator := request.validator
		id := request.contractID.Int64()

		liquidContract := liquidContracts[validator]
		totalLiquid, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.True(
			r.t,
			liquidOfVestingContract[validator].Cmp(totalLiquid) == 0,
			"bonding not applied", // it could happen if Autonity fails to call bondingApplied. Need immediate attention if happens
		)

		userLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		require.True(
			r.t,
			liquidOfUser[validator][id].Cmp(userLiquid) == 0,
			"vesting contract cannot track liquid balance",
		)

	}

	newNewtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.True(r.t, newNewtonBalance.Cmp(newtonBalance) == 0, "newton balance not updated")
}

func unbondAndRelease(
	r *runner, user common.Address, unbondingRequests []StakingRequest, unbondingGas *big.Int,
) {
	liquidContracts := make(map[common.Address]*Liquid)
	liquidOfUser := make(map[common.Address]map[int64]*big.Int)
	liquidOfVestingContract := make(map[common.Address]*big.Int)

	for i, validator := range r.committee.validators {
		for _, request := range unbondingRequests {
			if request.validator == validator.NodeAddress {
				liquidContract := r.committee.liquidContracts[i]
				liquidContracts[request.validator] = liquidContract

				balance, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
				require.NoError(r.t, err)
				liquidOfVestingContract[request.validator] = balance

				liquidOfUser[request.validator] = make(map[int64]*big.Int)
				break
			}
		}
	}

	contractNTN := make(map[int64]*big.Int)
	for _, request := range unbondingRequests {
		userLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, request.validator)
		require.NoError(r.t, err)
		liquidOfUser[request.validator][request.contractID.Int64()] = userLiquid

		contract, _, err := r.stakableVesting.GetContract(nil, user, request.contractID)
		require.NoError(r.t, err)
		contractNTN[request.contractID.Int64()] = contract.CurrentNTNAmount
	}

	unbondingRequestBlock := r.evm.Context.BlockNumber
	newtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)

	for _, request := range unbondingRequests {
		lockedLiquid, _, err := r.stakableVesting.LockedLiquidBalanceOf(nil, user, request.contractID, request.validator)
		require.NoError(r.t, err)
		unlockedLiquid, _, err := r.stakableVesting.UnlockedLiquidBalanceOf(nil, user, request.contractID, request.validator)
		require.NoError(r.t, err)
		_, err = r.stakableVesting.Unbond(
			fromSender(user, unbondingGas),
			request.contractID,
			request.validator,
			request.amount,
		)

		if request.expectedErr == "" {
			require.NoError(r.t, err)
			validator := request.validator
			id := request.contractID.Int64()
			liquidOfVestingContract[validator].Sub(liquidOfVestingContract[validator], request.amount)
			liquidOfUser[validator][id].Sub(liquidOfUser[validator][id], request.amount)
			contractNTN[id].Add(contractNTN[id], request.amount)

			newLockedLiquid, _, err := r.stakableVesting.LockedLiquidBalanceOf(nil, user, request.contractID, request.validator)
			require.NoError(r.t, err)
			require.True(
				r.t,
				new(big.Int).Add(lockedLiquid, request.amount).Cmp(newLockedLiquid) == 0,
				"vesting contract cannot track locked liquid",
			)

			newUnlockedLiquid, _, err := r.stakableVesting.UnlockedLiquidBalanceOf(nil, user, request.contractID, request.validator)
			require.NoError(r.t, err)
			require.True(
				r.t,
				new(big.Int).Sub(unlockedLiquid, request.amount).Cmp(newUnlockedLiquid) == 0,
				"vesting contract cannot track unlocked liquid",
			)

			newtonBalance.Add(newtonBalance, request.amount)
		} else {
			require.Equal(r.t, request.expectedErr, err.Error())
		}
	}

	r.waitNextEpoch()

	for _, request := range unbondingRequests {
		validator := request.validator
		id := request.contractID.Int64()
		liquidContract := liquidContracts[validator]

		totalLiquid, _, err := liquidContract.BalanceOf(nil, r.stakableVesting.address)
		require.NoError(r.t, err)
		require.True(
			r.t,
			totalLiquid.Cmp(liquidOfVestingContract[validator]) == 0,
			"unbonding not applied",
		)

		userLiquid, _, err := r.stakableVesting.LiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		require.True(
			r.t,
			userLiquid.Cmp(liquidOfUser[validator][id]) == 0,
			"vesting contract cannot track liquid",
		)

		lockedLiquid, _, err := r.stakableVesting.LockedLiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		require.True(
			r.t,
			lockedLiquid.Cmp(common.Big0) == 0,
			"vesting contract cannot track locked liquid",
		)

		unlockedLiquid, _, err := r.stakableVesting.UnlockedLiquidBalanceOf(nil, user, request.contractID, validator)
		require.NoError(r.t, err)
		require.True(
			r.t,
			unlockedLiquid.Cmp(userLiquid) == 0,
			"vesting contract cannot track unlocked liquid",
		)
	}

	// release unbonding
	unbondingPeriod, _, err := r.autonity.GetUnbondingPeriod(nil)
	require.NoError(r.t, err)
	unbondingReleaseBlock := new(big.Int).Add(unbondingRequestBlock, unbondingPeriod)
	for unbondingReleaseBlock.Cmp(r.evm.Context.BlockNumber) >= 0 {
		r.waitNextEpoch()
	}

	for _, request := range unbondingRequests {
		contract, _, err := r.stakableVesting.GetContract(nil, user, request.contractID)
		require.NoError(r.t, err)

		id := request.contractID.Int64()
		require.True(
			r.t,
			contract.CurrentNTNAmount.Cmp(contractNTN[id]) == 0,
			"contract not updated",
		)
	}

	newNewtonBalance, _, err := r.autonity.BalanceOf(nil, r.stakableVesting.address)
	require.NoError(r.t, err)
	require.True(r.t, newNewtonBalance.Cmp(newtonBalance) == 0, "vesting contract balance mismatch")
}
