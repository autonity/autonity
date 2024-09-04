package vestingtests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

// 1 million NTN
var reward = new(big.Int).Mul(big.NewInt(1000_000_000_000_000_000), big.NewInt(1000_000))

type StakingRequest struct {
	staker      common.Address
	validator   common.Address
	contractID  *big.Int
	amount      *big.Int
	expectedErr string
	bond        bool
}

type Reward struct {
	rewardATN *big.Int
	rewardNTN *big.Int
}

func TestReleaseFromStakableContract(t *testing.T) {
	r := tests.Setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.Evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	user := tests.User
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0
	// do not modify userBalance
	userBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)

	r.Run("cannot release before cliff", func(r *tests.Runner) {
		r.WaitSomeBlock(cliff)
		require.Equal(r.T, big.NewInt(cliff), r.Evm.Context.Time, "time mismatch")
		_, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())
		_, err = r.StakableVesting.ReleaseFunds(tests.FromSender(user, nil), contractID)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())
		userNewBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, userBalance, userNewBalance, "funds released before cliff period")
	})

	r.Run("release calculation follows epoch based linear function in time", func(r *tests.Runner) {
		currentTime := r.WaitSomeEpoch(cliff + 1)
		require.True(r.T, currentTime <= end+1, "release is not linear after end")
		// contract has the context of last block, so time is 1s less than currentTime
		unlocked := currentTime - 1 - start
		require.True(r.T, contractTotalAmount > unlocked, "cannot test if all funds unlocked")
		epochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		// mine some more blocks, release should be epoch based
		r.WaitNBlocks(10)
		currentTime += 10
		checkReleaseAllNTN(r, user, contractID, big.NewInt(unlocked))

		r.WaitNBlocks(10)
		currentTime += 10
		require.Equal(r.T, big.NewInt(currentTime), r.Evm.Context.Time, "time mismatch, release won't work")
		// no more should be released as epoch did not change
		newEpochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochID, newEpochID, "cannot test if epoch progresses")
		checkReleaseAllNTN(r, user, contractID, common.Big0)
	})

	r.Run("can release in chunks", func(r *tests.Runner) {
		currentTime := r.WaitSomeEpoch(cliff + 1)
		require.True(r.T, currentTime <= end+1, "cannot test, release is not linear after end")
		totalUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
		require.NoError(r.T, err)
		require.True(r.T, totalUnlocked.IsInt64(), "invalid data")
		require.True(r.T, totalUnlocked.Int64() > 1, "cannot test chunks")
		unlockFraction := big.NewInt(totalUnlocked.Int64() / 2)
		// release only a chunk of total unlocked
		r.NoError(
			r.StakableVesting.ReleaseNTN(tests.FromSender(user, nil), contractID, unlockFraction),
		)
		userNewBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(userBalance, unlockFraction), userNewBalance, "balance mismatch")
		data, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.True(r.T, data.IsInt64(), "invalid data")
		epochID := data.Int64()
		r.WaitNBlocks(10)
		data, _, err = r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.True(r.T, data.IsInt64(), "invalid data")
		require.Equal(r.T, epochID, data.Int64(), "epoch progressed, more funds will release")
		require.True(r.T, r.Evm.Context.Time.Cmp(big.NewInt(currentTime)) > 0, "time did not progress")
		checkReleaseAllNTN(r, user, contractID, new(big.Int).Sub(totalUnlocked, unlockFraction))
	})

	r.Run("cannot release more than total", func(r *tests.Runner) {
		r.WaitSomeEpoch(end + 1)
		// progress some more epoch, should not matter after end
		r.WaitNextEpoch()
		currentTime := r.Evm.Context.Time
		checkReleaseAllNTN(r, user, contractID, big.NewInt(contractTotalAmount))
		r.WaitNextEpoch()
		require.True(r.T, r.Evm.Context.Time.Cmp(currentTime) > 0, "time did not progress")
		// cannot release more
		checkReleaseAllNTN(r, user, contractID, common.Big0)
	})
}

func TestBonding(t *testing.T) {
	r := tests.Setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.Evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	contractCount := 2
	users, validators, liquidContracts := setupContracts(r, contractCount, 3, contractTotalAmount, start, cliff, end)

	beneficiary := users[0]
	contractID := common.Big0
	validator := validators[0]
	liquidContract := liquidContracts[0]

	r.Run("can bond all funds before cliff but not before start", func(r *tests.Runner) {
		require.True(r.T, r.Evm.Context.Time.Cmp(big.NewInt(start+1)) < 0, "contract started already")
		bondingAmount := big.NewInt(contractTotalAmount / 2)
		_, err := r.StakableVesting.Bond(tests.FromSender(beneficiary, nil), contractID, validator, bondingAmount)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: contract not started yet", err.Error())
		r.WaitSomeBlock(start + 1)
		require.True(r.T, r.Evm.Context.Time.Cmp(big.NewInt(cliff+1)) < 0, "contract cliff finished already")
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})
	})

	// start contract for bonding for all the tests remaining
	r.WaitSomeBlock(start + 1)

	r.Run("cannot bond more than total", func(r *tests.Runner) {
		bondingAmount := big.NewInt(contractTotalAmount + 10)
		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "execution reverted: not enough tokens", true}

		bondingAmount = big.NewInt(contractTotalAmount / 2)
		requests[1] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "", true}

		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		bondingAmount = new(big.Int).Add(big.NewInt(10), remaining)
		requests[2] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "execution reverted: not enough tokens", true}

		bondAndFinalize(r, requests)

		requests = make([]StakingRequest, 2)
		requests[0] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "execution reverted: not enough tokens", true}
		requests[1] = StakingRequest{beneficiary, validator, contractID, remaining, "", true}

		bondAndFinalize(r, requests)
	})

	r.Run("can release liquid tokens", func(r *tests.Runner) {
		bondingAmount := big.NewInt(contractTotalAmount)
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})
		currentTime := r.WaitSomeEpoch(cliff + 1)
		// contract has context of last block
		unlocked := currentTime - 1 - start
		// mine some more block, release should be epoch based
		r.WaitNBlocks(10)
		r.NoError(
			r.StakableVesting.ReleaseAllLNTN(tests.FromSender(beneficiary, nil), contractID),
		)
		liquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.Equal(
			r.T, big.NewInt(contractTotalAmount-unlocked), liquid,
			"liquid release don't follow epoch based linear function",
		)
		liquid, _, err = liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(contractTotalAmount-unlocked), liquid, "liquid not transferred")
		liquid, _, err = liquidContract.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlocked), liquid, "liquid not received")
		r.WaitSomeEpoch(end + 1)
		// progress more epoch, shouldn't matter
		r.WaitNextEpoch()
		r.NoError(
			r.StakableVesting.ReleaseAllLNTN(tests.FromSender(beneficiary, nil), contractID),
		)
		liquid, _, err = r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.True(r.T, liquid.Cmp(common.Big0) == 0, "all liquid tokens not released")
		liquid, _, err = liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.True(r.T, liquid.Cmp(common.Big0) == 0, "liquid not transferred")
		liquid, _, err = liquidContract.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(contractTotalAmount), liquid, "liquid not received")
	})

	r.Run("track liquids when bonding from multiple contracts to multiple validators", func(r *tests.Runner) {
		bondingIteration := 2
		bondingAmount := big.NewInt(contractTotalAmount / (int64(len(validators) * bondingIteration)))
		require.True(r.T, bondingAmount.Cmp(common.Big0) > 0, "cannot test")

		requests := make([]StakingRequest, 0)
		for _, user := range users {
			for i := 0; i < contractCount; i++ {
				for _, validator := range validators {
					requests = append(requests, StakingRequest{user, validator, big.NewInt(int64(i)), bondingAmount, "", true})
				}
			}
		}

		require.True(
			r.T,
			isInitialBalanceZero(r, requests),
			"initial liquid balances should be 0",
		)

		for i := 0; i < bondingIteration; i++ {
			bondAndFinalize(r, requests)
		}
	})

	r.Run("when bonded, release NTN first", func(r *tests.Runner) {
		liquidBalance, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.True(r.T, contractTotalAmount > 10, "cannot test")
		bondingAmount := big.NewInt(contractTotalAmount / 10)
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})
		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		require.True(r.T, remaining.Cmp(common.Big0) > 0, "no NTN remains")
		r.WaitSomeEpoch(cliff + 1)
		unlocked, _, err := r.StakableVesting.UnlockedFunds(nil, beneficiary, contractID)
		require.NoError(r.T, err)
		require.True(r.T, unlocked.Cmp(remaining) < 0, "don't want to release all NTN in the test")
		balance, _, err := r.Autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		r.NoError(
			r.StakableVesting.ReleaseFunds(tests.FromSender(beneficiary, nil), contractID),
		)
		newLiquidBalance, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(liquidBalance, bondingAmount), newLiquidBalance, "lquid released")
		newBalance, _, err := r.Autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, unlocked), newBalance, "balance not updated")
	})

	r.Run("can release LNTN", func(r *tests.Runner) {
		bondingAmount := big.NewInt(contractTotalAmount)
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})

		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlocked := currentTime - 1 - start
		if unlocked > contractTotalAmount {
			unlocked = contractTotalAmount
		}

		unlockedToken, _, err := r.StakableVesting.UnlockedFunds(nil, beneficiary, contractID)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlocked), unlockedToken)

		require.True(r.T, unlocked > 1, "cannot test")
		r.RunAndRevert(func(r *tests.Runner) {
			checkReleaseAllLNTN(r, beneficiary, contractID, unlockedToken)
		})

		currentUnlocked := big.NewInt(unlocked / 2)
		checkReleaseLNTN(r, beneficiary, validator, contractID, currentUnlocked)
		checkReleaseLNTN(r, beneficiary, validator, contractID, new(big.Int).Sub(unlockedToken, currentUnlocked))
	})

	r.Run("can release LNTN from any validator", func(r *tests.Runner) {
		// bond to all 3 validators, each 300 NTN
		// 100 NTN remains unbonded
		bondingAmount := big.NewInt(300)
		remainingNTN := big.NewInt(contractTotalAmount)
		requests := make([]StakingRequest, 0)
		for _, validator := range validators {
			requests = append(requests, StakingRequest{beneficiary, validator, contractID, bondingAmount, "", true})
			remainingNTN = new(big.Int).Sub(remainingNTN, bondingAmount)
		}
		bondAndFinalize(r, requests)
		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlockAmount := big.NewInt(currentTime - start - 1)
		// got at least 500 unlocked
		releaseAmount := big.NewInt(100)
		for _, validator := range validators {
			checkReleaseLNTN(r, beneficiary, validator, contractID, releaseAmount)
			unlockAmount = new(big.Int).Sub(unlockAmount, releaseAmount)
		}

		if unlockAmount.Cmp(remainingNTN) > 0 {
			unlockAmount = remainingNTN
		}
		checkReleaseAllNTN(r, beneficiary, contractID, unlockAmount)
	})
}

func TestUnbonding(t *testing.T) {
	r := tests.Setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.Evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	validatorCount := 2
	contractCount := 2
	users, validators, _ := setupContracts(r, contractCount, validatorCount, contractTotalAmount, start, cliff, end)

	// bond from all contracts to all validators
	r.WaitSomeBlock(start + 1)
	bondingAmount := big.NewInt(contractTotalAmount / int64(validatorCount))
	require.True(r.T, bondingAmount.Cmp(common.Big0) > 0, "not enough to bond")
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			for _, validator := range validators {
				r.NoError(
					r.StakableVesting.Bond(tests.FromSender(user, nil), big.NewInt(int64(i)), validator, bondingAmount),
				)
			}
		}
	}

	r.WaitNextEpoch()
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			r.NoError(
				r.StakableVesting.UpdateFunds(nil, user, big.NewInt(int64(i))),
			)
			totalLiquid := big.NewInt(0)
			for _, validator := range validators {
				liquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.T, err)
				require.Equal(r.T, bondingAmount, liquid)
				totalLiquid.Add(totalLiquid, liquid)
			}
			require.Equal(r.T, big.NewInt(contractTotalAmount), totalLiquid)
		}
	}

	// for testing single unbonding
	beneficiary := users[0]
	contractID := common.Big0
	validator := validators[0]

	r.Run("can unbond", func(r *tests.Runner) {
		liquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, bondingAmount, liquid, "liquid not minted properly")
		unbondAndRelease(r, []StakingRequest{{beneficiary, validator, contractID, liquid, "", false}})
	})

	r.Run("cannot unbond more than total liquid", func(r *tests.Runner) {
		unbondingAmount := new(big.Int).Add(bondingAmount, big.NewInt(10))
		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "execution reverted: not enough unlocked liquid tokens", false}

		unbondingAmount = big.NewInt(10)
		requests[1] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "", false}

		remaining := new(big.Int).Sub(bondingAmount, unbondingAmount)
		require.True(r.T, remaining.Cmp(common.Big0) > 0, "cannot test if no liquid remains")

		unbondingAmount = new(big.Int).Add(remaining, big.NewInt(10))
		requests[2] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "execution reverted: not enough unlocked liquid tokens", false}
		unbondAndRelease(r, requests)

		requests = make([]StakingRequest, 2)
		requests[0] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "execution reverted: not enough unlocked liquid tokens", false}
		requests[1] = StakingRequest{beneficiary, validator, contractID, remaining, "", false}
		unbondAndRelease(r, requests)
	})

	r.Run("cannot unbond if LNTN withdrawn", func(r *tests.Runner) {
		liquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		validator1 := validators[1]
		liquid1, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator1)
		require.NoError(r.T, err)
		require.True(r.T, liquid1.Cmp(big.NewInt(10)) > 0, "cannot test")

		totalToRelease := liquid.Int64() + 10
		currentTime := r.WaitSomeEpoch(totalToRelease + start + 1)
		totalToRelease = currentTime - 1 - start
		r.NoError(
			r.StakableVesting.ReleaseAllLNTN(tests.FromSender(beneficiary, nil), contractID),
		)

		// LNTN will be released from the first validator in the list
		newLiquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.True(r.T, newLiquid.Cmp(common.Big0) == 0, "liquid remains after releasing")

		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{beneficiary, validator, contractID, liquid, "execution reverted: not enough unlocked liquid tokens", false}

		// if more unlocked funds remain, then LNTN will be released from 2nd validator
		releasedFromValidator1 := totalToRelease - liquid.Int64()
		remainingLiquid := new(big.Int).Sub(liquid1, big.NewInt(releasedFromValidator1))
		requests[1] = StakingRequest{beneficiary, validator1, contractID, liquid1, "execution reverted: not enough unlocked liquid tokens", false}

		liquid1, _, err = r.StakableVesting.LiquidBalanceOf(nil, beneficiary, contractID, validator1)
		require.NoError(r.T, err)
		require.Equal(r.T, remainingLiquid, liquid1, "liquid balance mismatch")

		requests[2] = StakingRequest{beneficiary, validator1, contractID, liquid1, "", false}
		unbondAndRelease(r, requests)
	})

	r.Run("track liquid when unbonding from multiple contracts to multiple validators", func(r *tests.Runner) {
		// unbond few
		unbondingAmount := big.NewInt(100)

		requests := make([]StakingRequest, 0)
		for _, user := range users {
			for _, validator := range validators {
				for i := 0; i < contractCount; i++ {
					requests = append(requests, StakingRequest{user, validator, big.NewInt(int64(i)), unbondingAmount, "", false})
				}
			}
		}
		unbondAndRelease(r, requests)

		// unbond the rest
		unbondingAmount, _, err := r.StakableVesting.LiquidBalanceOf(nil, users[0], common.Big0, validators[0])
		require.NoError(r.T, err)
		requests = make([]StakingRequest, 0)
		for _, user := range users {
			for _, validator := range validators {
				for i := 0; i < contractCount; i++ {
					requests = append(requests, StakingRequest{user, validator, big.NewInt(int64(i)), unbondingAmount, "", false})
				}
			}
		}
		unbondAndRelease(r, requests)
	})
}

func TestRwardTracking(t *testing.T) {
	r := tests.Setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.Evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	contractCount := 2
	users, validators, liquidContracts := setupContracts(r, contractCount, 2, contractTotalAmount, start, cliff, end)

	// start contract to bond
	r.WaitSomeBlock(start + 1)

	r.Run("bond and claim reward", func(r *tests.Runner) {
		beneficiary := users[0]
		contractID := common.Big0
		validator := validators[0]
		liquidContract := liquidContracts[0]
		bondingAmount := big.NewInt(contractTotalAmount)
		r.NoError(
			r.StakableVesting.Bond(
				tests.FromSender(beneficiary, nil), contractID, validator, bondingAmount,
			),
		)
		r.WaitNextEpoch()
		r.GiveMeSomeMoney(r.Autonity.Address(), reward)
		r.WaitNextEpoch()

		rewardOfContract, _, err := liquidContract.UnclaimedRewards(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.True(r.T, rewardOfContract.UnclaimedNTN.Cmp(common.Big0) > 0, "no NTN reward")
		require.True(r.T, rewardOfContract.UnclaimedATN.Cmp(common.Big0) > 0, "no ATN reward")

		rewardOfUser, _, err := r.StakableVesting.UnclaimedRewards(nil, beneficiary, contractID, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, rewardOfContract.UnclaimedATN, rewardOfUser.AtnRewards, "ATN reward mismatch")
		require.Equal(r.T, rewardOfContract.UnclaimedNTN, rewardOfUser.NtnRewards, "NTN reward mismatch")

		rewardOfUser, _, err = r.StakableVesting.UnclaimedRewards1(nil, beneficiary, contractID)
		require.NoError(r.T, err)
		require.Equal(r.T, rewardOfContract.UnclaimedATN, rewardOfUser.AtnRewards, "ATN reward mismatch")
		require.Equal(r.T, rewardOfContract.UnclaimedNTN, rewardOfUser.NtnRewards, "NTN reward mismatch")

		rewardOfUser, _, err = r.StakableVesting.UnclaimedRewards0(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, rewardOfContract.UnclaimedATN, rewardOfUser.AtnRewards, "ATN reward mismatch")
		require.Equal(r.T, rewardOfContract.UnclaimedNTN, rewardOfUser.NtnRewards, "NTN reward mismatch")

		// checking all the variations of claim rewards function
		r.RunAndRevert(func(r *tests.Runner) {
			checkClaimRewardsFunction(
				r, beneficiary, rewardOfUser.AtnRewards, rewardOfUser.NtnRewards,
				func() {
					r.NoError(
						r.StakableVesting.ClaimRewards0(tests.FromSender(beneficiary, nil)),
					)
				},
			)
		})

		r.RunAndRevert(func(r *tests.Runner) {
			checkClaimRewardsFunction(
				r, beneficiary, rewardOfUser.AtnRewards, rewardOfUser.NtnRewards,
				func() {
					r.NoError(
						r.StakableVesting.ClaimRewards(tests.FromSender(beneficiary, nil), contractID),
					)
				},
			)
		})

		r.RunAndRevert(func(r *tests.Runner) {
			checkClaimRewardsFunction(
				r, beneficiary, rewardOfUser.AtnRewards, rewardOfUser.NtnRewards,
				func() {
					r.NoError(
						r.StakableVesting.ClaimRewards1(tests.FromSender(beneficiary, nil), contractID, validator),
					)
				},
			)
		})
	})

	// set commission rate = 0, so all rewards go to delegation
	r.NoError(
		r.Autonity.SetTreasuryFee(operator, common.Big0),
	)
	for _, validator := range r.Committee.Validators {
		r.NoError(
			r.Autonity.ChangeCommissionRate(
				tests.FromSender(validator.Treasury, nil), validator.NodeAddress, common.Big0,
			),
		)
	}
	for r.Committee.Validators[0].CommissionRate.Cmp(common.Big0) > 0 {
		r.WaitNextEpoch()
	}

	// remove all bonding, so we only have bonding from contracts only
	for _, validator := range r.Committee.Validators {
		require.Equal(r.T, validator.SelfBondedStake, validator.BondedStake, "delegation stake should not exist")
		r.NoError(
			r.Autonity.Unbond(
				tests.FromSender(validator.Treasury, nil), validator.NodeAddress, validator.SelfBondedStake,
			),
		)
	}

	// bond from contracts
	bondingAmount := big.NewInt(100)
	totalBonded := new(big.Int)
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			bondedVals, _, err := r.StakableVesting.GetBondedValidators(nil, user, big.NewInt(int64(i)))
			require.NoError(r.T, err)
			require.True(r.T, len(bondedVals) == 0)
			for _, validator := range validators {
				r.NoError(
					r.StakableVesting.Bond(
						tests.FromSender(user, nil), big.NewInt(int64(i)), validator, bondingAmount,
					),
				)
				totalBonded.Add(totalBonded, bondingAmount)
			}
		}
	}

	r.WaitNextEpoch()

	require.Equal(r.T, len(validators), len(r.Committee.Validators), "committee not updated properly")
	eachValidatorDelegation := big.NewInt(int64(len(users) * contractCount))
	eachValidatorStake := new(big.Int).Mul(bondingAmount, eachValidatorDelegation)
	for i, validator := range r.Committee.Validators {
		require.Equal(r.T, eachValidatorStake, validator.BondedStake)
		require.True(r.T, validator.SelfBondedStake.Cmp(common.Big0) == 0)
		balance, _, err := r.Committee.LiquidContracts[i].BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, eachValidatorStake, balance)
	}

	r.Run("bond in differenet epoch and track reward", func(r *tests.Runner) {
		extraBonds := make([]StakingRequest, 0)

		for _, user := range users {
			extraBonds = append(extraBonds, StakingRequest{user, validators[0], common.Big0, bondingAmount, "", true})
			extraBonds = append(extraBonds, StakingRequest{user, validators[0], common.Big1, bondingAmount, "", true})
			extraBonds = append(extraBonds, StakingRequest{user, validators[1], common.Big0, bondingAmount, "", true})
			extraBonds = append(extraBonds, StakingRequest{user, validators[0], common.Big0, bondingAmount, "", true})
		}
		// dummy
		extraBonds = append(extraBonds, StakingRequest{common.Address{}, validators[0], common.Big0, common.Big0, "", true})

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		require.True(
			r.T,
			isAllRewardsZero(r, contractCount, liquidContracts, users, validators),
			"all rewards should be zero initially",
		)

		for _, request := range extraBonds {

			user := request.staker
			if request.amount.Cmp(common.Big0) > 0 {
				r.NoError(
					r.StakableVesting.Bond(
						tests.FromSender(user, nil), request.contractID, request.validator, request.amount,
					),
				)
			}

			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldRewardsFromValidator, oldUserRewards := unclaimedRewards(r, contractCount, liquidContracts, users, validators)
			totalReward := rewardsAfterOneEpoch(r)
			r.WaitNextEpoch()

			// request is not applied yet
			checkRewards(
				r, contractCount, totalStake, totalReward,
				liquidContracts, validators, users, validatorStakes,
				userStakes, oldRewardsFromValidator, oldUserRewards,
			)

			// request is applied, because checkRewards progress 1 epoch
			if request.amount.Cmp(common.Big0) > 0 {
				amount := request.amount
				validator := request.validator
				id := int(request.contractID.Int64())
				validatorStakes[validator].Add(validatorStakes[validator], amount)
				userStakes[user][id][validator].Add(userStakes[user][id][validator], amount)
				totalStake.Add(totalStake, amount)
			}
		}
	})

	r.Run("multiple bonding in same epoch and track rewards", func(r *tests.Runner) {
		extraBondsArray := make([][]StakingRequest, 0)
		extraBonds := make([]StakingRequest, 0)

		for _, user := range users {
			for contractID := 0; contractID < contractCount; contractID++ {
				for _, validator := range validators {
					extraBonds = append(extraBonds, StakingRequest{user, validator, big.NewInt(int64(contractID)), bondingAmount, "", true})
				}
			}
		}

		testCount := 2
		for testCount > 0 {
			extraBondsArray = append(extraBondsArray, extraBonds)
			testCount--
		}
		// dummy
		extraBondsArray = append(extraBondsArray, nil)

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		require.True(
			r.T,
			isAllRewardsZero(r, contractCount, liquidContracts, users, validators),
			"all rewards should be zero initially",
		)

		for _, requests := range extraBondsArray {
			for _, request := range requests {

				user := request.staker
				r.NoError(
					r.StakableVesting.Bond(
						tests.FromSender(user, nil), request.contractID, request.validator, request.amount,
					),
				)
			}

			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldRewardsFromValidator, oldUserRewards := unclaimedRewards(r, contractCount, liquidContracts, users, validators)
			totalReward := rewardsAfterOneEpoch(r)
			r.WaitNextEpoch()

			// request is not applied yet
			checkRewards(
				r, contractCount, totalStake, totalReward,
				liquidContracts, validators, users, validatorStakes,
				userStakes, oldRewardsFromValidator, oldUserRewards,
			)

			for _, request := range requests {
				// request is applied, because checkRewards progress 1 epoch
				user := request.staker
				amount := request.amount
				validator := request.validator
				id := int(request.contractID.Int64())
				validatorStakes[validator].Add(validatorStakes[validator], amount)
				userStakes[user][id][validator].Add(userStakes[user][id][validator], amount)
				totalStake.Add(totalStake, amount)
			}
		}

		validatorNewStakes, userNewStakes, totalNewStake := initialStakes(r, contractCount, liquidContracts, users, validators)
		require.Equal(r.T, totalStake, totalNewStake)
		for validator, newStake := range validatorNewStakes {
			require.Equal(r.T, validatorStakes[validator], newStake)
		}
		for user, newStake := range userNewStakes {
			require.Equal(r.T, userStakes[user], newStake)
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
					r.StakableVesting.Bond(
						tests.FromSender(user, nil), big.NewInt(int64(i)), validator, bondingAmount,
					),
				)
				totalBonded.Add(totalBonded, bondingAmount)
			}
		}
	}
	bondingAmount.Add(bondingAmount, oldBondingAmount)

	r.WaitNextEpoch()

	r.Run("release liquid and track reward", func(r *tests.Runner) {
		r.WaitSomeEpoch(end + 1)
		releaseAmount := big.NewInt(100)
		userLiquidBalance := make(map[common.Address]map[common.Address]*big.Int)
		// unbonding request can be treated as release request
		releaseRequests := make([]StakingRequest, 0)

		for _, user := range users {
			userLiquidBalance[user] = make(map[common.Address]*big.Int)
			for _, validator := range validators {
				userLiquidBalance[user][validator] = new(big.Int)
			}
			releaseRequests = append(releaseRequests, StakingRequest{user, validators[0], common.Big0, releaseAmount, "", false})
			releaseRequests = append(releaseRequests, StakingRequest{user, validators[0], common.Big1, releaseAmount, "", false})
			releaseRequests = append(releaseRequests, StakingRequest{user, validators[1], common.Big0, releaseAmount, "", false})
			releaseRequests = append(releaseRequests, StakingRequest{user, validators[0], common.Big0, releaseAmount, "", false})
		}
		// dummy
		releaseRequests = append(releaseRequests, StakingRequest{common.Address{}, validators[0], common.Big0, common.Big0, "", false})

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		liquidContractsMap := make(map[common.Address]*tests.Liquid)

		for i, liquidContract := range liquidContracts {
			liquidContractsMap[validators[i]] = liquidContract
		}

		for _, request := range releaseRequests {

			// some epoch is passed and we are entitled to some reward,
			// but we don't know about it because we did not get notified
			// or we did not claim them or call unclaimedRewards
			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldRewardsFromValidator, oldUserRewards := unclaimedRewards(r, contractCount, liquidContracts, users, validators)
			totalReward := rewardsAfterOneEpoch(r)
			r.WaitNextEpoch()

			// we release some LNTN and it is applied immediately
			// if unlocked, it is transferred immediately
			// but for reward calculation, it will be applied later
			user := request.staker
			amount := request.amount
			validator := request.validator
			if request.amount.Cmp(common.Big0) > 0 {
				r.NoError(
					r.StakableVesting.ReleaseLNTN(
						tests.FromSender(user, nil),
						request.contractID,
						request.validator,
						request.amount,
					),
				)

				userLiquidBalance[user][validator].Add(userLiquidBalance[user][validator], amount)
				balance, _, err := liquidContractsMap[validator].BalanceOf(nil, user)
				require.NoError(r.T, err)
				require.Equal(r.T, userLiquidBalance[user][validator], balance, "liquid not transferred")
			}

			checkRewards(
				r, contractCount, totalStake, totalReward,
				liquidContracts, validators, users, validatorStakes,
				userStakes, oldRewardsFromValidator, oldUserRewards,
			)

			// for next reward
			if request.amount.Cmp(common.Big0) > 0 {
				id := int(request.contractID.Int64())
				validatorStakes[validator].Sub(validatorStakes[validator], amount)
				userStakes[user][id][validator].Sub(userStakes[user][id][validator], amount)
			}
		}
	})

	r.Run("unbond in different epoch and track reward", func(r *tests.Runner) {
		unbondingAmount := big.NewInt(100)
		extraUnbonds := make([]StakingRequest, 0)
		for _, user := range users {
			extraUnbonds = append(extraUnbonds, StakingRequest{user, validators[0], common.Big0, unbondingAmount, "", false})
			extraUnbonds = append(extraUnbonds, StakingRequest{user, validators[0], common.Big1, unbondingAmount, "", false})
			extraUnbonds = append(extraUnbonds, StakingRequest{user, validators[1], common.Big0, unbondingAmount, "", false})
			extraUnbonds = append(extraUnbonds, StakingRequest{user, validators[0], common.Big0, unbondingAmount, "", false})
		}
		// dummy
		extraUnbonds = append(extraUnbonds, StakingRequest{common.Address{}, validators[0], common.Big0, common.Big0, "", false})

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidContracts, users, validators)

		for _, request := range extraUnbonds {

			user := request.staker
			if request.amount.Cmp(common.Big0) > 0 {
				r.NoError(
					r.StakableVesting.Unbond(
						tests.FromSender(user, nil), request.contractID, request.validator, request.amount,
					),
				)
			}

			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldRewardsFromValidator, oldUserRewards := unclaimedRewards(r, contractCount, liquidContracts, users, validators)
			totalReward := rewardsAfterOneEpoch(r)
			r.WaitNextEpoch()
			// request is not applied yet
			checkRewards(
				r, contractCount, totalStake, totalReward,
				liquidContracts, validators, users, validatorStakes,
				userStakes, oldRewardsFromValidator, oldUserRewards,
			)

			// request is applied, because checkRewards progress 1 epoch
			if request.amount.Cmp(common.Big0) > 0 {
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
	r := tests.Setup(t, nil)
	var contractTotalAmount int64 = 1000
	start := 100 + r.Evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	user := tests.User
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0
	newUser := common.HexToAddress("0x88")
	require.NotEqual(r.T, user, newUser)

	r.Run("beneficiary changes", func(r *tests.Runner) {
		_, _, err := r.StakableVesting.GetContract(nil, user, contractID)
		require.NoError(r.T, err)
		_, _, err = r.StakableVesting.GetContract(nil, newUser, contractID)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: invalid contract id", err.Error())
		r.NoError(
			r.StakableVesting.ChangeContractBeneficiary(operator, user, contractID, newUser),
		)
		_, _, err = r.StakableVesting.GetContract(nil, newUser, contractID)
		require.NoError(r.T, err)
		_, _, err = r.StakableVesting.GetContract(nil, user, contractID)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: invalid contract id", err.Error())
	})

	r.Run("beneficiary does not lose claim to rewards", func(_ *tests.Runner) {
		bondingAmount := big.NewInt(contractTotalAmount)
		validator := r.Committee.Validators[0].NodeAddress
		r.WaitSomeBlock(start + 1)
		r.NoError(
			r.StakableVesting.Bond(
				tests.FromSender(user, nil), contractID, validator, bondingAmount,
			),
		)
		r.WaitNextEpoch()
		r.GiveMeSomeMoney(r.Autonity.Address(), reward)
		r.WaitNextEpoch()
		totalReward, _, err := r.StakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.T, err)
		rewards, _, err := r.StakableVesting.UnclaimedRewards(nil, user, contractID, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, rewards, totalReward)

		// change beneficiary
		r.NoError(
			r.StakableVesting.ChangeContractBeneficiary(operator, user, contractID, newUser),
		)
		totalReward, _, err = r.StakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, rewards, totalReward)
	})
}

func TestSlashingAffect(t *testing.T) {
	r := tests.Setup(t, nil)
	// TODO (tariq): complete tests.Setup

	r.Run("contract total value update when bonded validator slashed", func(r *tests.Runner) {
		// TODO (tariq): complete
	})
}

func TestAccessRestriction(t *testing.T) {
	r := tests.Setup(t, nil)
	user := tests.User

	r.Run("only operator can create contract", func(r *tests.Runner) {
		amount := big.NewInt(1000)
		start := new(big.Int).Add(big.NewInt(100), r.Evm.Context.Time)
		cliff := new(big.Int).Add(start, big.NewInt(100))
		end := new(big.Int).Add(start, amount)
		_, err := r.StakableVesting.NewContract(
			tests.FromSender(user, nil),
			user,
			amount,
			start,
			cliff,
			end,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})

	var contractTotalAmount int64 = 1000
	start := r.Evm.Context.Time.Int64()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	createContract(r, user, contractTotalAmount, start, cliff, end)
	contractID := common.Big0

	r.Run("only operator can change contract beneficiary", func(r *tests.Runner) {
		newUser := common.HexToAddress("0x88")
		require.NotEqual(r.T, user, newUser)
		_, err := r.StakableVesting.ChangeContractBeneficiary(
			tests.FromSender(user, nil),
			user,
			contractID,
			newUser,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())

		_, err = r.StakableVesting.ChangeContractBeneficiary(
			tests.FromSender(newUser, nil),
			user,
			contractID,
			newUser,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())
	})
}

func initialStakes(
	r *tests.Runner,
	contractCount int,
	liquidContracts []*tests.Liquid,
	users, validators []common.Address,
) (
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[int]map[common.Address]*big.Int,
	totalStake *big.Int,
) {

	// need to update funds before querying the initial stakes
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			r.NoError(
				r.StakableVesting.UpdateFunds(nil, user, big.NewInt(int64(i))),
			)
		}
	}

	totalStake = new(big.Int)

	validatorStakes = make(map[common.Address]*big.Int)
	for i, validator := range validators {
		liquidContract := liquidContracts[i]
		balance, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		validatorStakes[validator] = balance
	}

	userStakes = make(map[common.Address]map[int]map[common.Address]*big.Int)
	for _, user := range users {
		userStakes[user] = make(map[int]map[common.Address]*big.Int)
		for i := 0; i < contractCount; i++ {
			userStakes[user][i] = make(map[common.Address]*big.Int)
			for _, validator := range validators {
				balance, _, err := r.StakableVesting.LiquidBalanceOf(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.T, err)
				userStakes[user][i][validator] = balance
				totalStake.Add(totalStake, balance)
			}
		}
	}
	return validatorStakes, userStakes, totalStake
}

func unclaimedRewards(
	r *tests.Runner,
	contractCount int,
	liquidContracts []*tests.Liquid,
	users, validators []common.Address,
) (
	oldRewardsFromValidator map[common.Address]Reward,
	oldUserRewards map[common.Address]map[int]map[common.Address]Reward,
) {

	oldRewardsFromValidator = make(map[common.Address]Reward)
	for i, validator := range validators {
		liquidContract := liquidContracts[i]
		unclaimedReward, _, err := liquidContract.UnclaimedRewards(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		oldRewardsFromValidator[validator] = Reward{unclaimedReward.UnclaimedATN, unclaimedReward.UnclaimedNTN}
	}

	oldUserRewards = make(map[common.Address]map[int]map[common.Address]Reward)
	for _, user := range users {
		oldUserRewards[user] = make(map[int]map[common.Address]Reward)
		for i := 0; i < contractCount; i++ {
			oldUserRewards[user][i] = make(map[common.Address]Reward)
			for _, validator := range validators {
				unclaimedReward, _, err := r.StakableVesting.UnclaimedRewards(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.T, err)
				oldUserRewards[user][i][validator] = Reward{unclaimedReward.AtnRewards, unclaimedReward.NtnRewards}
			}
		}
	}

	return oldRewardsFromValidator, oldUserRewards
}

func rewardsAfterOneEpoch(r *tests.Runner) (rewardsToDistribute Reward) {

	// get supply and inflationReserve to calculate inflation reward
	supply, _, err := r.Autonity.TotalSupply(nil)
	require.NoError(r.T, err)
	inflationReserve, _, err := r.Autonity.InflationReserve(nil)
	require.NoError(r.T, err)
	epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
	require.NoError(r.T, err)

	// get inflation reward
	lastEpochTime, _, err := r.Autonity.LastEpochTime(nil)
	require.NoError(r.T, err)
	currentEpochTime := new(big.Int).Add(lastEpochTime, epochPeriod)
	rewardsToDistribute.rewardNTN, _, err = r.InflationController.CalculateSupplyDelta(nil, supply, inflationReserve, lastEpochTime, currentEpochTime)
	require.NoError(r.T, err)

	// get atn reward
	rewardsToDistribute.rewardATN = r.GetBalanceOf(r.Autonity.Address())
	return rewardsToDistribute
}

func checkClaimRewardsFunction(
	r *tests.Runner,
	account common.Address,
	unclaimedAtnRewards *big.Int,
	unclaimedNtnRewards *big.Int,
	claimRewards func(),
) {
	atnBalance := r.GetBalanceOf(account)
	ntnBalance, _, err := r.Autonity.BalanceOf(nil, account)
	require.NoError(r.T, err)
	claimRewards()

	newAtnBalance := r.GetBalanceOf(account)
	newNtnBalance, _, err := r.Autonity.BalanceOf(nil, account)
	require.NoError(r.T, err)

	atnRewards := new(big.Int).Sub(newAtnBalance, atnBalance)
	ntnRewards := new(big.Int).Sub(newNtnBalance, ntnBalance)

	require.True(
		r.T,
		atnRewards.Cmp(unclaimedAtnRewards) == 0,
		"claimed atn rewards mismatch",
	)

	require.True(
		r.T,
		ntnRewards.Cmp(unclaimedNtnRewards) == 0,
		"claimed ntn rewards mismatch",
	)
}

func checkRewards(
	r *tests.Runner,
	contractCount int,
	totalStake *big.Int,
	totalReward Reward,
	liquidContracts []*tests.Liquid,
	validators, users []common.Address,
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[int]map[common.Address]*big.Int,
	oldRewardsFromValidator map[common.Address]Reward,
	oldUserRewards map[common.Address]map[int]map[common.Address]Reward,
) {

	currentRewards := make(map[common.Address]Reward)
	// check total rewards from each validator
	for i, validator := range validators {
		validatorTotalRewardATN := new(big.Int).Mul(validatorStakes[validator], totalReward.rewardATN)
		validatorTotalRewardNTN := new(big.Int).Mul(validatorStakes[validator], totalReward.rewardNTN)

		if totalStake.Cmp(common.Big0) != 0 {
			validatorTotalRewardATN = validatorTotalRewardATN.Div(validatorTotalRewardATN, totalStake)
			validatorTotalRewardNTN = validatorTotalRewardNTN.Div(validatorTotalRewardNTN, totalStake)
		}

		liquidContract := liquidContracts[i]
		unclaimedReward, _, err := liquidContract.UnclaimedRewards(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)

		diff := new(big.Int).Sub(
			new(big.Int).Add(validatorTotalRewardATN, oldRewardsFromValidator[validator].rewardATN),
			unclaimedReward.UnclaimedATN,
		)
		diff.Abs(diff)
		// difference should be less than or equal to 1 wei
		require.True(
			r.T,
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
			r.T,
			diff.Cmp(common.Big1) <= 0,
			"unclaimed ntn reward not updated in liquid contract",
		)
		currentRewards[validator] = Reward{
			new(big.Int).Sub(unclaimedReward.UnclaimedATN, oldRewardsFromValidator[validator].rewardATN),
			new(big.Int).Sub(unclaimedReward.UnclaimedNTN, oldRewardsFromValidator[validator].rewardNTN),
		}
	}

	// check each user rewards
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

				unclaimedReward, _, err := r.StakableVesting.UnclaimedRewards(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.T, err)

				diff := new(big.Int).Sub(calculatedRewardATN, unclaimedReward.AtnRewards)
				diff.Abs(diff)
				// difference should be less than or equal to 1 wei
				require.True(
					r.T,
					diff.Cmp(common.Big1) <= 0,
					"atn reward calculation mismatch",
				)

				diff = new(big.Int).Sub(calculatedRewardNTN, unclaimedReward.NtnRewards)
				diff.Abs(diff)
				// difference should be less than or equal to 1 wei
				require.True(
					r.T,
					diff.Cmp(common.Big1) <= 0,
					"ntn reward calculation mismatch",
				)
				unclaimedRewardForContractATN.Add(unclaimedRewardForContractATN, unclaimedReward.AtnRewards)
				unclaimedRewardForContractNTN.Add(unclaimedRewardForContractNTN, unclaimedReward.NtnRewards)

				// so that following code snippet reverts
				r.RunAndRevert(func(r *tests.Runner) {
					checkClaimRewardsFunction(
						r, user, unclaimedReward.AtnRewards, unclaimedReward.NtnRewards,
						func() {
							r.NoError(
								r.StakableVesting.ClaimRewards1(tests.FromSender(user, nil), big.NewInt(int64(i)), validator),
							)
						},
					)

				})
			}

			unclaimedReward, _, err := r.StakableVesting.UnclaimedRewards1(nil, user, big.NewInt(int64(i)))
			require.NoError(r.T, err)

			require.Equal(r.T, unclaimedRewardForContractATN, unclaimedReward.AtnRewards)
			require.Equal(r.T, unclaimedRewardForContractNTN, unclaimedReward.NtnRewards)

			userRewardATN.Add(userRewardATN, unclaimedReward.AtnRewards)
			userRewardNTN.Add(userRewardNTN, unclaimedReward.NtnRewards)

			// so that following code snippet reverts
			r.RunAndRevert(func(r *tests.Runner) {
				checkClaimRewardsFunction(
					r, user, unclaimedReward.AtnRewards, unclaimedReward.NtnRewards,
					func() {
						r.NoError(
							r.StakableVesting.ClaimRewards(tests.FromSender(user, nil), big.NewInt(int64(i))),
						)
					},
				)
			})
		}

		unclaimedReward, _, err := r.StakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.T, err)

		require.Equal(
			r.T,
			userRewardATN,
			unclaimedReward.AtnRewards,
			"unclaimed atn reward mismatch",
		)

		require.Equal(
			r.T,
			userRewardNTN,
			unclaimedReward.NtnRewards,
			"unclaimed ntn reward mismatch",
		)

		// so that following code snippet reverts
		r.RunAndRevert(func(r *tests.Runner) {
			checkClaimRewardsFunction(
				r, user, unclaimedReward.AtnRewards, unclaimedReward.NtnRewards,
				func() {
					r.NoError(
						r.StakableVesting.ClaimRewards0(tests.FromSender(user, nil)),
					)
				},
			)
		})
	}
}

func isAllRewardsZero(
	r *tests.Runner,
	contractCount int,
	liquidContracts []*tests.Liquid,
	users, validators []common.Address,
) bool {

	for _, liquidContract := range liquidContracts {
		rewards, _, err := liquidContract.UnclaimedRewards(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)

		if rewards.UnclaimedATN.Cmp(common.Big0) != 0 {
			return false
		}

		if rewards.UnclaimedNTN.Cmp(common.Big0) != 0 {
			return false
		}
	}

	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			for _, validator := range validators {
				rewards, _, err := r.StakableVesting.UnclaimedRewards(nil, user, big.NewInt(int64(i)), validator)
				require.NoError(r.T, err)

				if rewards.AtnRewards.Cmp(common.Big0) != 0 {
					return false
				}

				if rewards.NtnRewards.Cmp(common.Big0) != 0 {
					return false
				}
			}

			rewards, _, err := r.StakableVesting.UnclaimedRewards1(nil, user, big.NewInt(int64(i)))
			require.NoError(r.T, err)

			if rewards.AtnRewards.Cmp(common.Big0) != 0 {
				return false
			}

			if rewards.NtnRewards.Cmp(common.Big0) != 0 {
				return false
			}
		}

		rewards, _, err := r.StakableVesting.UnclaimedRewards0(nil, user)
		require.NoError(r.T, err)

		if rewards.AtnRewards.Cmp(common.Big0) != 0 {
			return false
		}

		if rewards.NtnRewards.Cmp(common.Big0) != 0 {
			return false
		}
	}
	return true
}

func setupContracts(
	r *tests.Runner, contractCount, validatorCount int, contractTotalAmount, start, cliff, end int64,
) (users, validators []common.Address, liquidContracts []*tests.Liquid) {
	users = make([]common.Address, 2)
	users[0] = tests.User
	users[1] = common.HexToAddress("0x88")
	require.NotEqual(r.T, users[0], users[1], "same user")
	for _, user := range users {
		initBalance := new(big.Int).Mul(big.NewInt(1000_000), big.NewInt(1000_000_000_000_000_000))
		r.GiveMeSomeMoney(user, initBalance)
		for i := 0; i < contractCount; i++ {
			createContract(r, user, contractTotalAmount, start, cliff, end)
		}
	}

	// use multiple validators
	validators = make([]common.Address, validatorCount)
	liquidContracts = make([]*tests.Liquid, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validators[i] = r.Committee.Validators[i].NodeAddress
		liquidContracts[i] = r.Committee.LiquidContracts[i]
	}
	return
}

func createContract(r *tests.Runner, beneficiary common.Address, amount, startTime, cliffTime, endTime int64) {
	startBig := big.NewInt(startTime)
	cliffBig := big.NewInt(cliffTime)
	endBig := big.NewInt(endTime)
	r.NoError(
		r.StakableVesting.NewContract(
			operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
			new(big.Int).Sub(cliffBig, startBig), new(big.Int).Sub(endBig, startBig),
		),
	)
}

func checkReleaseAllNTN(r *tests.Runner, user common.Address, contractID, releaseAmount *big.Int) {
	contract, _, err := r.StakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.T, err)
	contractNTN := contract.CurrentNTNAmount
	withdrawn := contract.WithdrawnValue
	initBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)
	totalUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.T, err)
	r.NoError(
		r.StakableVesting.ReleaseAllNTN(tests.FromSender(user, nil), contractID),
	)
	newBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)
	require.Equal(r.T, new(big.Int).Add(initBalance, releaseAmount), newBalance, "balance mismatch")
	contract, _, err = r.StakableVesting.GetContract(nil, user, contractID)
	require.NoError(r.T, err)
	require.True(
		r.T,
		new(big.Int).Sub(contractNTN, releaseAmount).Cmp(contract.CurrentNTNAmount) == 0,
		"contract NTN not updated properly",
	)
	require.True(
		r.T,
		new(big.Int).Add(withdrawn, releaseAmount).Cmp(contract.WithdrawnValue) == 0,
		"contract WithdrawnValue not updated properly",
	)

	remainingUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.T, err)
	require.True(r.T, new(big.Int).Sub(totalUnlocked, releaseAmount).Cmp(remainingUnlocked) == 0)

	// all unlocked NTN withdrawn
	require.True(r.T, contract.CurrentNTNAmount.Cmp(common.Big0) == 0 || remainingUnlocked.Cmp(common.Big0) == 0)
}

func checkReleaseAllLNTN(r *tests.Runner, user common.Address, contractID, releaseAmount *big.Int) {
	totalUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.T, err)

	bondedValidators, _, err := r.StakableVesting.GetBondedValidators(nil, user, contractID)
	require.NoError(r.T, err)

	userLiquidBalances := make([]*big.Int, 0)
	vaultLiquidBalances := make([]*big.Int, 0)
	userLiquidInVesting := make([]*big.Int, 0)
	for _, validator := range bondedValidators {
		liquidContract := r.LiquidContract(validator)
		balance, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		vaultLiquidBalances = append(vaultLiquidBalances, balance)

		balance, _, err = liquidContract.BalanceOf(nil, user)
		require.NoError(r.T, err)
		userLiquidBalances = append(userLiquidBalances, balance)

		balance, _, err = r.StakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
		require.NoError(r.T, err)
		userLiquidInVesting = append(userLiquidInVesting, balance)
	}

	r.NoError(
		r.StakableVesting.ReleaseAllLNTN(tests.FromSender(user, nil), contractID),
	)

	remainingUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.T, err)
	require.True(r.T, new(big.Int).Sub(totalUnlocked, releaseAmount).Cmp(remainingUnlocked) == 0)

	for i, validator := range bondedValidators {
		released := releaseAmount
		if releaseAmount.Cmp(vaultLiquidBalances[i]) > 0 {
			released = vaultLiquidBalances[i]
		}
		liquidContract := r.LiquidContract(validator)
		balance, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Sub(vaultLiquidBalances[i], released).Cmp(balance) == 0)

		balance, _, err = liquidContract.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Add(userLiquidBalances[i], released).Cmp(balance) == 0)

		balance, _, err = r.StakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Sub(userLiquidInVesting[i], released).Cmp(balance) == 0)
	}
}

func checkReleaseLNTN(r *tests.Runner, user, validator common.Address, contractID, releaseAmount *big.Int) {

	var liquidContract *tests.Liquid
	for i, v := range r.Committee.Validators {
		if v.NodeAddress == validator {
			liquidContract = r.Committee.LiquidContracts[i]
			break
		}
	}

	liquidBalance, _, err := liquidContract.BalanceOf(nil, user)
	require.NoError(r.T, err)

	liquidInVesting, _, err := r.StakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
	require.NoError(r.T, err)

	totalUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.T, err)

	r.NoError(
		r.StakableVesting.ReleaseLNTN(tests.FromSender(user, nil), contractID, validator, releaseAmount),
	)

	newLiquidBalance, _, err := liquidContract.BalanceOf(nil, user)
	require.NoError(r.T, err)
	require.Equal(
		r.T,
		new(big.Int).Add(liquidBalance, releaseAmount),
		newLiquidBalance,
	)

	newLiquidInVesting, _, err := r.StakableVesting.LiquidBalanceOf(nil, user, contractID, validator)
	require.NoError(r.T, err)
	require.True(
		r.T,
		newLiquidInVesting.Cmp(new(big.Int).Sub(liquidInVesting, releaseAmount)) == 0,
	)

	remainingUnlocked, _, err := r.StakableVesting.UnlockedFunds(nil, user, contractID)
	require.NoError(r.T, err)
	require.True(r.T, new(big.Int).Sub(totalUnlocked, releaseAmount).Cmp(remainingUnlocked) == 0)
}

func isInitialBalanceZero(
	r *tests.Runner,
	requests []StakingRequest,
) bool {

	updateVestingContractFunds(r, requests)

	for _, request := range requests {
		balance, _, err := r.StakableVesting.LiquidBalanceOf(nil, request.staker, request.contractID, request.validator)
		require.NoError(r.T, err)
		if balance.Cmp(common.Big0) != 0 {
			return false
		}
	}

	for _, liquidContract := range r.Committee.LiquidContracts {
		balance, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		if balance.Cmp(common.Big0) != 0 {
			return false
		}
	}

	return true

}

func initialBalances(
	r *tests.Runner,
	stakingRequests []StakingRequest,
) (
	liquidContracts map[common.Address]*tests.Liquid,
	liquidOfVestingContract map[common.Address]*big.Int,
	liquidOfUser map[common.Address]map[common.Address]map[int64]*big.Int,
	contractNTN map[common.Address]map[int64]*big.Int,
) {
	liquidContracts = make(map[common.Address]*tests.Liquid)
	liquidOfVestingContract = make(map[common.Address]*big.Int)
	liquidOfUser = make(map[common.Address]map[common.Address]map[int64]*big.Int)
	contractNTN = make(map[common.Address]map[int64]*big.Int)

	for _, request := range stakingRequests {
		liquidOfUser[request.staker] = make(map[common.Address]map[int64]*big.Int)
		contractNTN[request.staker] = make(map[int64]*big.Int)
	}

	for _, request := range stakingRequests {
		liquidOfUser[request.staker][request.validator] = make(map[int64]*big.Int)
	}

	for i, validator := range r.Committee.Validators {
		for _, request := range stakingRequests {
			if request.validator == validator.NodeAddress {
				liquidContract := r.Committee.LiquidContracts[i]
				liquidContracts[request.validator] = liquidContract

				balance, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
				require.NoError(r.T, err)
				liquidOfVestingContract[request.validator] = balance

				break
			}
		}
	}

	for _, request := range stakingRequests {
		userLiquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, request.staker, request.contractID, request.validator)
		require.NoError(r.T, err)
		liquidOfUser[request.staker][request.validator][request.contractID.Int64()] = userLiquid

		contract, _, err := r.StakableVesting.GetContract(nil, request.staker, request.contractID)
		require.NoError(r.T, err)
		contractNTN[request.staker][request.contractID.Int64()] = contract.CurrentNTNAmount
	}
	return liquidContracts, liquidOfVestingContract, liquidOfUser, contractNTN
}

func updateVestingContractFunds(r *tests.Runner, stakingRequests []StakingRequest) {
	for _, request := range stakingRequests {
		r.NoError(
			r.StakableVesting.UpdateFunds(nil, request.staker, request.contractID),
		)
	}
}

func bondAndFinalize(
	r *tests.Runner, bondingRequests []StakingRequest,
) {

	liquidContracts, liquidOfVestingContract, liquidOfUser, contractNTNs := initialBalances(r, bondingRequests)

	newtonBalance, _, err := r.Autonity.BalanceOf(nil, r.StakableVesting.Address())
	require.NoError(r.T, err)

	for _, request := range bondingRequests {
		contract, _, err := r.StakableVesting.GetContract(nil, request.staker, request.contractID)
		require.NoError(r.T, err)
		contractNTN := contract.CurrentNTNAmount

		_, bondErr := r.StakableVesting.Bond(
			tests.FromSender(request.staker, nil),
			request.contractID,
			request.validator,
			request.amount,
		)

		contract, _, err = r.StakableVesting.GetContract(nil, request.staker, request.contractID)
		require.NoError(r.T, err)
		remaining := new(big.Int).Sub(contractNTN, common.Big0)

		if request.expectedErr == "" {
			require.NoError(r.T, bondErr)
			validator := request.validator
			id := request.contractID.Int64()
			liquidOfVestingContract[validator].Add(liquidOfVestingContract[validator], request.amount)
			liquidOfUser[request.staker][validator][id].Add(liquidOfUser[request.staker][validator][id], request.amount)

			remaining.Sub(remaining, request.amount)
			contractNTNs[request.staker][id].Sub(contractNTNs[request.staker][id], request.amount)

			newtonBalance.Sub(newtonBalance, request.amount)
		} else {
			require.Error(r.T, bondErr)
			require.Equal(r.T, request.expectedErr, bondErr.Error())
		}
		require.True(r.T, remaining.Cmp(contract.CurrentNTNAmount) == 0, "contract not updated properly")
	}

	// let bonding apply
	r.WaitNextEpoch()

	// need to update funds in vesting contract
	updateVestingContractFunds(r, bondingRequests)

	for _, request := range bondingRequests {
		validator := request.validator
		id := request.contractID.Int64()

		liquidContract := liquidContracts[validator]
		totalLiquid, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.True(
			r.T,
			liquidOfVestingContract[validator].Cmp(totalLiquid) == 0,
			"bonding not applied", // it could happen if Autonity fails to call bondingApplied. Need immediate attention if happens
		)

		userLiquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, request.staker, request.contractID, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			liquidOfUser[request.staker][validator][id].Cmp(userLiquid) == 0,
			"vesting contract cannot track liquid balance",
		)

		contract, _, err := r.StakableVesting.GetContract(nil, request.staker, request.contractID)
		require.NoError(r.T, err)
		require.True(
			r.T,
			contract.CurrentNTNAmount.Cmp(contractNTNs[request.staker][id]) == 0,
			"contract not updated properly",
		)

	}

	newNewtonBalance, _, err := r.Autonity.BalanceOf(nil, r.StakableVesting.Address())
	require.NoError(r.T, err)
	require.True(r.T, newNewtonBalance.Cmp(newtonBalance) == 0, "newton balance not updated")
}

func unbondAndRelease(
	r *tests.Runner, unbondingRequests []StakingRequest,
) {
	liquidContracts, liquidOfVestingContract, liquidOfUser, contractNTN := initialBalances(r, unbondingRequests)

	unbondingRequestBlock := r.Evm.Context.BlockNumber
	newtonBalance, _, err := r.Autonity.BalanceOf(nil, r.StakableVesting.Address())
	require.NoError(r.T, err)

	for _, request := range unbondingRequests {
		lockedLiquid, _, err := r.StakableVesting.LockedLiquidBalanceOf(nil, request.staker, request.contractID, request.validator)
		require.NoError(r.T, err)
		unlockedLiquid, _, err := r.StakableVesting.UnlockedLiquidBalanceOf(nil, request.staker, request.contractID, request.validator)
		require.NoError(r.T, err)
		_, unbondErr := r.StakableVesting.Unbond(
			tests.FromSender(request.staker, nil),
			request.contractID,
			request.validator,
			request.amount,
		)

		if request.expectedErr == "" {
			require.NoError(r.T, unbondErr)
			validator := request.validator
			id := request.contractID.Int64()
			liquidOfVestingContract[validator].Sub(liquidOfVestingContract[validator], request.amount)
			liquidOfUser[request.staker][validator][id].Sub(liquidOfUser[request.staker][validator][id], request.amount)
			contractNTN[request.staker][id].Add(contractNTN[request.staker][id], request.amount)

			newLockedLiquid, _, err := r.StakableVesting.LockedLiquidBalanceOf(nil, request.staker, request.contractID, request.validator)
			require.NoError(r.T, err)
			require.True(
				r.T,
				new(big.Int).Add(lockedLiquid, request.amount).Cmp(newLockedLiquid) == 0,
				"vesting contract cannot track locked liquid",
			)

			newUnlockedLiquid, _, err := r.StakableVesting.UnlockedLiquidBalanceOf(nil, request.staker, request.contractID, request.validator)
			require.NoError(r.T, err)
			require.True(
				r.T,
				new(big.Int).Sub(unlockedLiquid, request.amount).Cmp(newUnlockedLiquid) == 0,
				"vesting contract cannot track unlocked liquid",
			)

			newtonBalance.Add(newtonBalance, request.amount)
		} else {
			require.Error(r.T, unbondErr)
			require.Equal(r.T, request.expectedErr, unbondErr.Error())
		}
	}

	r.WaitNextEpoch()
	// need to update funds in vesting contract
	updateVestingContractFunds(r, unbondingRequests)

	for _, request := range unbondingRequests {
		validator := request.validator
		id := request.contractID.Int64()
		liquidContract := liquidContracts[validator]

		totalLiquid, _, err := liquidContract.BalanceOf(nil, r.StakableVesting.Address())
		require.NoError(r.T, err)
		require.True(
			r.T,
			totalLiquid.Cmp(liquidOfVestingContract[validator]) == 0,
			"unbonding not applied",
		)

		userLiquid, _, err := r.StakableVesting.LiquidBalanceOf(nil, request.staker, request.contractID, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			userLiquid.Cmp(liquidOfUser[request.staker][validator][id]) == 0,
			"vesting contract cannot track liquid",
		)

		lockedLiquid, _, err := r.StakableVesting.LockedLiquidBalanceOf(nil, request.staker, request.contractID, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			lockedLiquid.Cmp(common.Big0) == 0,
			"vesting contract cannot track locked liquid",
		)

		unlockedLiquid, _, err := r.StakableVesting.UnlockedLiquidBalanceOf(nil, request.staker, request.contractID, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			unlockedLiquid.Cmp(userLiquid) == 0,
			"vesting contract cannot track unlocked liquid",
		)
	}

	// release unbonding
	unbondingPeriod, _, err := r.Autonity.GetUnbondingPeriod(nil)
	require.NoError(r.T, err)
	unbondingReleaseBlock := new(big.Int).Add(unbondingRequestBlock, unbondingPeriod)
	for unbondingReleaseBlock.Cmp(r.Evm.Context.BlockNumber) >= 0 {
		r.WaitNextEpoch()
	}

	updateVestingContractFunds(r, unbondingRequests)

	for _, request := range unbondingRequests {
		contract, _, err := r.StakableVesting.GetContract(nil, request.staker, request.contractID)
		require.NoError(r.T, err)

		id := request.contractID.Int64()
		require.True(
			r.T,
			contract.CurrentNTNAmount.Cmp(contractNTN[request.staker][id]) == 0,
			"contract not updated",
		)
	}

	newNewtonBalance, _, err := r.Autonity.BalanceOf(nil, r.StakableVesting.Address())
	require.NoError(r.T, err)
	require.True(r.T, newNewtonBalance.Cmp(newtonBalance) == 0, "vesting contract balance mismatch")
}
