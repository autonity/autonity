package vestingtests

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
)

// 1 million NTN
var reward = new(big.Int).Mul(big.NewInt(1000_000_000_000_000_000), big.NewInt(1000_000))

type Reward tests.EpochReward

type StakingRequest struct {
	staker      common.Address
	validator   common.Address
	contractID  *big.Int
	amount      *big.Int
	expectedErr string
	bond        bool
}

func TestSeparateAccountForStakableVestingContract(t *testing.T) {
	user := tests.User
	start := time.Now().Unix()
	var contractAmount int64 = 100
	cliff := start + 10
	end := start + contractAmount
	contractID := common.Big0
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("contract creation transfer funds", setup, func(r *tests.Runner) {
		managerBalance, _, err := r.Autonity.BalanceOf(nil, r.StakableVestingManager.Address())
		require.NoError(r.T, err)
		var stakableContractAddress common.Address

		r.RunAndRevert(func(r *tests.Runner) {
			createContract(r, user, contractAmount, start, cliff, end)
			var err error
			stakableContractAddress, _, err = r.StakableVestingManager.GetContractAccount(nil, user, contractID)
			require.NoError(r.T, err)
		})
		contractBalance, _, err := r.Autonity.BalanceOf(nil, stakableContractAddress)
		require.NoError(r.T, err)
		require.True(r.T, contractBalance.Cmp(common.Big0) == 0)

		createContract(r, user, contractAmount, start, cliff, end)
		managerBalanceNew, _, err := r.Autonity.BalanceOf(nil, r.StakableVestingManager.Address())
		require.NoError(r.T, err)
		require.True(
			r.T,
			managerBalanceNew.Cmp(new(big.Int).Sub(managerBalance, big.NewInt(contractAmount))) == 0,
			"manager balance not updated",
		)

		stakableContract := r.StakableVestingContractObject(user, contractID)
		require.Equal(r.T, stakableContractAddress, stakableContract.Address())
		contractBalance, _, err = r.Autonity.BalanceOf(nil, stakableContractAddress)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(contractAmount), contractBalance, "ntn not transferred")
	})

	tests.RunWithSetup("stakable vesting contract NTN amount is equal to contract balance", setup, func(r *tests.Runner) {
		createContract(r, user, contractAmount, start, cliff, end)
		stakableContract := r.StakableVestingContractObject(user, contractID)
		contract, _, err := stakableContract.GetContract(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(contractAmount), contract.CurrentNTNAmount)
		balance, _, err := r.Autonity.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, contract.CurrentNTNAmount, balance)
	})

	tests.RunWithSetup("separate stakable vesting contracts are in separate smart contract", setup, func(r *tests.Runner) {
		createContract(r, user, contractAmount, start, cliff, end)
		createContract(r, user, contractAmount, start, cliff, end)
		newUser := common.HexToAddress("0x123132")
		createContract(r, newUser, contractAmount, start, cliff, end)
		stakableContractAddress1 := r.StakableVestingContractObject(user, common.Big0).Address()
		stakableContractAddress2 := r.StakableVestingContractObject(user, common.Big1).Address()
		stakableContractAddress3 := r.StakableVestingContractObject(newUser, common.Big0).Address()
		require.NotEqual(r.T, stakableContractAddress1, stakableContractAddress2)
		require.NotEqual(r.T, stakableContractAddress1, stakableContractAddress3)
		require.NotEqual(r.T, stakableContractAddress2, stakableContractAddress3)
	})
}

func TestReleaseFromStakableContract(t *testing.T) {
	contractID := common.Big0
	var contractTotalAmount int64 = 1000
	start := 100 + time.Now().Unix()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	user := tests.User
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}
	initiate := func(r *tests.Runner) (
		userBalance *big.Int,
		stakableContract *tests.StakableVesting,
	) {
		createContract(r, user, contractTotalAmount, start, cliff, end)
		stakableContract = r.StakableVestingContractObject(user, contractID)
		// do not modify userBalance
		userBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		return
	}

	tests.RunWithSetup("cannot release before cliff", setup, func(r *tests.Runner) {
		userBalance, stakableContract := initiate(r)
		r.WaitSomeBlock(cliff)
		require.Equal(r.T, big.NewInt(cliff), r.Evm.Context.Time, "time mismatch")
		_, _, err := stakableContract.UnlockedFunds(nil)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())
		_, err = stakableContract.ReleaseFunds(tests.FromSender(user, nil))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: cliff period not reached yet", err.Error())
		userNewBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, userBalance, userNewBalance, "funds released before cliff period")
	})

	tests.RunWithSetup("unlocking mechanism follows epoch based linear function in time", setup, func(r *tests.Runner) {
		_, stakableContract := initiate(r)
		currentTime := r.WaitSomeEpoch(cliff + 1)
		require.True(r.T, currentTime <= end+1, "release is not linear after end")
		// contract has the context of last block, so time is 1s less than currentTime
		unlocked := currentTime - 1 - start
		require.True(r.T, contractTotalAmount > unlocked, "cannot test if all funds unlocked")
		epochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		// mine some more blocks, release should be epoch based
		r.WaitNBlocks(10)
		newEpochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochID, newEpochID, "cannot test if epoch progresses")
		unlockedFunds, _, err := stakableContract.UnlockedFunds(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlocked), unlockedFunds)
	})

	tests.RunWithSetup("release calculation follows epoch based linear function in time", setup, func(r *tests.Runner) {
		initiate(r)
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

	tests.RunWithSetup("can release in chunks", setup, func(r *tests.Runner) {
		userBalance, stakableContract := initiate(r)
		currentTime := r.WaitSomeEpoch(cliff + 1)
		require.True(r.T, currentTime <= end+1, "cannot test, release is not linear after end")
		totalUnlocked, _, err := stakableContract.UnlockedFunds(nil)
		require.NoError(r.T, err)
		require.True(r.T, totalUnlocked.IsInt64(), "invalid data")
		require.True(r.T, totalUnlocked.Int64() > 1, "cannot test chunks")
		unlockFraction := big.NewInt(totalUnlocked.Int64() / 2)
		// release only a chunk of total unlocked
		r.NoError(
			stakableContract.ReleaseNTN(tests.FromSender(user, nil), unlockFraction),
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

	tests.RunWithSetup("cannot release more than total", setup, func(r *tests.Runner) {
		initiate(r)
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
	contractCount := 2
	contractID := common.Big0
	var contractTotalAmount int64 = 1000
	start := 10 + time.Now().Unix()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	setup := func() *tests.Runner {
		r := tests.Setup(t, nil)
		return r
	}
	initiate := func(r *tests.Runner) (
		users, validators []common.Address,
		liquidStateContract *tests.ILiquidLogic,
		beneficiary, validator common.Address,
		stakableContract *tests.StakableVesting,
	) {
		users, validators, liquidStateContracts := setupContracts(r, contractCount, 3, contractTotalAmount, start, cliff, end)
		beneficiary = users[0]
		validator = validators[0]
		liquidStateContract = liquidStateContracts[0]
		stakableContract = r.StakableVestingContractObject(beneficiary, contractID)
		return
	}

	tests.RunWithSetup("can bond all funds before cliff but not before start", setup, func(r *tests.Runner) {
		_, _, _, beneficiary, validator, stakableContract := initiate(r)
		require.True(r.T, r.Evm.Context.Time.Cmp(big.NewInt(start+1)) < 0, "contract started already")
		bondingAmount := big.NewInt(contractTotalAmount / 2)
		_, err := stakableContract.Bond(tests.FromSender(beneficiary, nil), validator, bondingAmount)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: contract not started yet", err.Error())
		r.WaitSomeBlock(start + 1)
		require.True(r.T, r.Evm.Context.Time.Cmp(big.NewInt(cliff+1)) < 0, "contract cliff finished already")
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})
	})

	initiate2 := func(r *tests.Runner) (
		users, validators []common.Address,
		liquidStateContract *tests.ILiquidLogic,
		beneficiary, validator common.Address,
		stakableContract *tests.StakableVesting,
	) {
		users, validators, liquidStateContract, beneficiary, validator, stakableContract = initiate(r)
		r.WaitSomeBlock(start + 1)
		return
	}
	// start contract for bonding for all the tests remaining

	tests.RunWithSetup("cannot bond more than total", setup, func(r *tests.Runner) {
		_, _, _, beneficiary, validator, _ := initiate2(r)
		bondingAmount := big.NewInt(contractTotalAmount + 10)
		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "execution reverted: insufficient Newton balance", true}

		bondingAmount = big.NewInt(contractTotalAmount / 2)
		requests[1] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "", true}

		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		bondingAmount = new(big.Int).Add(big.NewInt(10), remaining)
		requests[2] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "execution reverted: insufficient Newton balance", true}

		bondAndFinalize(r, requests)

		requests = make([]StakingRequest, 2)
		requests[0] = StakingRequest{beneficiary, validator, contractID, bondingAmount, "execution reverted: insufficient Newton balance", true}
		requests[1] = StakingRequest{beneficiary, validator, contractID, remaining, "", true}

		bondAndFinalize(r, requests)
	})

	tests.RunWithSetup("can release liquid tokens", setup, func(r *tests.Runner) {
		_, _, liquidStateContract, beneficiary, validator, stakableContract := initiate2(r)
		bondingAmount := big.NewInt(contractTotalAmount)
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})
		currentTime := r.WaitSomeEpoch(cliff + 1)
		// contract has context of last block
		unlocked := currentTime - 1 - start
		// mine some more block, release should be epoch based
		r.WaitNBlocks(10)
		r.NoError(
			stakableContract.ReleaseAllLNTN(tests.FromSender(beneficiary, nil)),
		)
		liquid, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.Equal(
			r.T, big.NewInt(contractTotalAmount-unlocked), liquid,
			"liquid release don't follow epoch based linear function",
		)
		liquid, _, err = liquidStateContract.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(contractTotalAmount-unlocked), liquid, "liquid not transferred")
		liquid, _, err = liquidStateContract.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(unlocked), liquid, "liquid not received")
		r.WaitSomeEpoch(end + 1)
		// progress more epoch, shouldn't matter
		r.WaitNextEpoch()
		r.NoError(
			stakableContract.ReleaseAllLNTN(tests.FromSender(beneficiary, nil)),
		)
		liquid, _, err = stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(r.T, liquid.Cmp(common.Big0) == 0, "all liquid tokens not released")
		liquid, _, err = liquidStateContract.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.True(r.T, liquid.Cmp(common.Big0) == 0, "liquid not transferred")
		liquid, _, err = liquidStateContract.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, big.NewInt(contractTotalAmount), liquid, "liquid not received")
	})

	tests.RunWithSetup("track liquids when bonding from multiple contracts to multiple validators", setup, func(r *tests.Runner) {
		users, validators, _, _, _, _ := initiate2(r)
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

	tests.RunWithSetup("when bonded, release NTN first", setup, func(r *tests.Runner) {
		_, _, _, beneficiary, validator, stakableContract := initiate2(r)
		liquidBalance, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(r.T, contractTotalAmount > 10, "cannot test")
		bondingAmount := big.NewInt(contractTotalAmount / 10)
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})
		remaining := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingAmount)
		require.True(r.T, remaining.Cmp(common.Big0) > 0, "no NTN remains")
		r.WaitSomeEpoch(cliff + 1)
		unlocked, _, err := stakableContract.UnlockedFunds(nil)
		require.NoError(r.T, err)
		require.True(r.T, unlocked.Cmp(remaining) < 0, "don't want to release all NTN in the test")
		balance, _, err := r.Autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		r.NoError(
			stakableContract.ReleaseFunds(tests.FromSender(beneficiary, nil)),
		)
		newLiquidBalance, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(liquidBalance, bondingAmount), newLiquidBalance, "lquid released")
		newBalance, _, err := r.Autonity.BalanceOf(nil, beneficiary)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(balance, unlocked), newBalance, "balance not updated")
	})

	tests.RunWithSetup("can release LNTN", setup, func(r *tests.Runner) {
		_, _, _, beneficiary, validator, stakableContract := initiate2(r)
		bondingAmount := big.NewInt(contractTotalAmount)
		bondAndFinalize(r, []StakingRequest{{beneficiary, validator, contractID, bondingAmount, "", true}})

		currentTime := r.WaitSomeEpoch(cliff + 1)
		unlocked := currentTime - 1 - start
		if unlocked > contractTotalAmount {
			unlocked = contractTotalAmount
		}

		unlockedToken, _, err := stakableContract.UnlockedFunds(nil)
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

	tests.RunWithSetup("can release LNTN from any validator", setup, func(r *tests.Runner) {
		_, validators, _, beneficiary, _, _ := initiate2(r)
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
	beneficiary := tests.User
	contractID := common.Big0
	var contractTotalAmount int64 = 1000
	start := 100 + time.Now().Unix()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	validatorCount := 2
	contractCount := 2
	bondingAmount := big.NewInt(contractTotalAmount / int64(validatorCount))

	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	initiate := func(r *tests.Runner) (
		users, validators []common.Address,
		validator common.Address,
		stakableContract *tests.StakableVesting,
	) {
		users, validators, _ = setupContracts(r, contractCount, validatorCount, contractTotalAmount, start, cliff, end)

		// bond from all contracts to all validators
		r.WaitSomeBlock(start + 1)
		require.True(r.T, bondingAmount.Cmp(common.Big0) > 0, "not enough to bond")
		for _, user := range users {
			for i := 0; i < contractCount; i++ {
				vestingContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
				for _, validator := range validators {
					r.NoError(
						vestingContract.Bond(tests.FromSender(user, nil), validator, bondingAmount),
					)
				}
			}
		}

		r.WaitNextEpoch()
		for _, user := range users {
			for i := 0; i < contractCount; i++ {
				vestingContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
				// r.NoError(
				// 	r.StakableVestingManager.UpdateFunds(nil, user, big.NewInt(int64(i))),
				// )
				totalLiquid := big.NewInt(0)
				for _, validator := range validators {
					liquid, _, err := vestingContract.LiquidBalance(nil, validator)
					require.NoError(r.T, err)
					require.Equal(r.T, bondingAmount, liquid)
					totalLiquid.Add(totalLiquid, liquid)
				}
				require.Equal(r.T, big.NewInt(contractTotalAmount), totalLiquid)
			}
		}

		// for testing single unbonding
		validator = validators[0]
		stakableContract = r.StakableVestingContractObject(beneficiary, contractID)
		return users, validators, validator, stakableContract
	}

	tests.RunWithSetup("can unbond", setup, func(r *tests.Runner) {
		_, _, validator, stakableContract := initiate(r)
		liquid, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, bondingAmount, liquid, "liquid not minted properly")
		unbondAndRelease(r, []StakingRequest{{beneficiary, validator, contractID, liquid, "", false}})
	})

	tests.RunWithSetup("cannot unbond more than total liquid", setup, func(r *tests.Runner) {
		_, _, validator, _ := initiate(r)
		unbondingAmount := new(big.Int).Add(bondingAmount, big.NewInt(10))
		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "execution reverted: insufficient unlocked Liquid Newton balance", false}

		unbondingAmount = big.NewInt(10)
		requests[1] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "", false}

		remaining := new(big.Int).Sub(bondingAmount, unbondingAmount)
		require.True(r.T, remaining.Cmp(common.Big0) > 0, "cannot test if no liquid remains")

		unbondingAmount = new(big.Int).Add(remaining, big.NewInt(10))
		requests[2] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "execution reverted: insufficient unlocked Liquid Newton balance", false}
		unbondAndRelease(r, requests)

		requests = make([]StakingRequest, 2)
		requests[0] = StakingRequest{beneficiary, validator, contractID, unbondingAmount, "execution reverted: insufficient unlocked Liquid Newton balance", false}
		requests[1] = StakingRequest{beneficiary, validator, contractID, remaining, "", false}
		unbondAndRelease(r, requests)
	})

	tests.RunWithSetup("cannot unbond if LNTN withdrawn", setup, func(r *tests.Runner) {
		_, validators, validator, stakableContract := initiate(r)
		liquid, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		validator1 := validators[1]
		liquid1, _, err := stakableContract.LiquidBalance(nil, validator1)
		require.NoError(r.T, err)
		require.True(r.T, liquid1.Cmp(big.NewInt(10)) > 0, "cannot test")

		totalToRelease := liquid.Int64() + 10
		currentTime := r.WaitSomeEpoch(totalToRelease + start + 1)
		totalToRelease = currentTime - 1 - start
		r.NoError(
			stakableContract.ReleaseAllLNTN(tests.FromSender(beneficiary, nil)),
		)

		// LNTN will be released from the first validator in the list
		newLiquid, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(r.T, newLiquid.Cmp(common.Big0) == 0, "liquid remains after releasing")

		requests := make([]StakingRequest, 3)
		requests[0] = StakingRequest{beneficiary, validator, contractID, liquid, "execution reverted: insufficient unlocked Liquid Newton balance", false}

		// if more unlocked funds remain, then LNTN will be released from 2nd validator
		releasedFromValidator1 := totalToRelease - liquid.Int64()
		remainingLiquid := new(big.Int).Sub(liquid1, big.NewInt(releasedFromValidator1))
		requests[1] = StakingRequest{beneficiary, validator1, contractID, liquid1, "execution reverted: insufficient unlocked Liquid Newton balance", false}

		liquid1, _, err = stakableContract.LiquidBalance(nil, validator1)
		require.NoError(r.T, err)
		require.Equal(r.T, remainingLiquid, liquid1, "liquid balance mismatch")

		requests[2] = StakingRequest{beneficiary, validator1, contractID, liquid1, "", false}
		unbondAndRelease(r, requests)
	})

	tests.RunWithSetup("track liquid when unbonding from multiple contracts to multiple validators", setup, func(r *tests.Runner) {
		users, validators, _, stakableContract := initiate(r)
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
		unbondingAmount, _, err := stakableContract.LiquidBalance(nil, validators[0])
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

func TestRewardTracking(t *testing.T) {
	var contractTotalAmount int64 = 1000
	start := 100 + time.Now().Unix()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	contractCount := 2

	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	initiate := func(r *tests.Runner) (
		users, validators []common.Address,
		liquidStateContracts []*tests.ILiquidLogic,
	) {
		users, validators, liquidStateContracts = setupContracts(r, contractCount, 2, contractTotalAmount, start, cliff, end)
		// start contract to bond
		r.WaitSomeBlock(start + 1)
		return
	}

	tests.RunWithSetup("bond and claim reward", setup, func(r *tests.Runner) {
		users, validators, liquidStateContracts := initiate(r)
		beneficiary := users[0]
		contractID := common.Big0
		stakableContract := r.StakableVestingContractObject(beneficiary, contractID)
		validator := validators[0]
		liquidStateContract := liquidStateContracts[0]
		bondingAmount := big.NewInt(contractTotalAmount)
		r.NoError(
			stakableContract.Bond(
				tests.FromSender(beneficiary, nil), validator, bondingAmount,
			),
		)
		r.WaitNextEpoch()
		r.GiveMeSomeMoney(r.Autonity.Address(), reward)
		r.WaitNextEpoch()

		rewardOfContract, _, err := liquidStateContract.UnclaimedRewards(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.True(r.T, rewardOfContract.UnclaimedNTN.Cmp(common.Big0) > 0, "no NTN reward")
		require.True(r.T, rewardOfContract.UnclaimedATN.Cmp(common.Big0) > 0, "no ATN reward")

		rewardOfUser, _, err := stakableContract.UnclaimedRewards(nil, validator)
		require.NoError(r.T, err)
		require.Equal(r.T, rewardOfContract.UnclaimedATN, rewardOfUser.AtnRewards, "ATN reward mismatch")
		require.Equal(r.T, rewardOfContract.UnclaimedNTN, rewardOfUser.NtnRewards, "NTN reward mismatch")

		rewardOfUser, _, err = stakableContract.UnclaimedRewards0(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, rewardOfContract.UnclaimedATN, rewardOfUser.AtnRewards, "ATN reward mismatch")
		require.Equal(r.T, rewardOfContract.UnclaimedNTN, rewardOfUser.NtnRewards, "NTN reward mismatch")

		// checking all the variations of claim rewards function
		r.RunAndRevert(func(r *tests.Runner) {
			checkClaimRewardsFunction(
				r, beneficiary, rewardOfUser.AtnRewards, rewardOfUser.NtnRewards,
				func() {
					r.NoError(
						stakableContract.ClaimRewards(tests.FromSender(beneficiary, nil)),
					)
				},
			)
		})

		r.RunAndRevert(func(r *tests.Runner) {
			checkClaimRewardsFunction(
				r, beneficiary, rewardOfUser.AtnRewards, rewardOfUser.NtnRewards,
				func() {
					r.NoError(
						stakableContract.ClaimRewards0(tests.FromSender(beneficiary, nil), validator),
					)
				},
			)
		})
	})

	bondingAmount := big.NewInt(100)
	initiate2 := func(r *tests.Runner) (
		users, validators []common.Address,
		liquidStateContracts []*tests.ILiquidLogic,
	) {
		users, validators, liquidStateContracts = initiate(r)
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
		for _, user := range users {
			for i := 0; i < contractCount; i++ {
				stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
				bondedVals, _, err := stakableContract.GetLinkedValidators(nil)
				require.NoError(r.T, err)
				require.True(r.T, len(bondedVals) == 0)
				for _, validator := range validators {
					r.NoError(
						stakableContract.Bond(
							tests.FromSender(user, nil), validator, bondingAmount,
						),
					)
					// totalBonded.Add(totalBonded, bondingAmount)
				}
			}
		}

		r.WaitNextEpoch()

		require.Equal(r.T, len(validators), len(r.Committee.Validators), "committee not updated properly")
		eachValidatorDelegationCount := big.NewInt(int64(len(users) * contractCount))
		eachValidatorStake := new(big.Int).Mul(bondingAmount, eachValidatorDelegationCount)
		for i, validator := range r.Committee.Validators {
			require.Equal(r.T, eachValidatorStake, validator.BondedStake)
			require.True(r.T, validator.SelfBondedStake.Cmp(common.Big0) == 0)
			for _, user := range users {
				for j := 0; j < contractCount; j++ {
					stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(j)))
					balance, _, err := r.Committee.LiquidStateContracts[i].BalanceOf(nil, stakableContract.Address())
					require.NoError(r.T, err)
					require.Equal(r.T, bondingAmount, balance)
				}
			}
		}
		return users, validators, liquidStateContracts
	}

	tests.RunWithSetup("bond in different epoch and track reward", setup, func(r *tests.Runner) {
		users, validators, liquidStateContracts := initiate2(r)
		extraBonds := make([]StakingRequest, 0)

		for _, user := range users {
			extraBonds = append(extraBonds, StakingRequest{user, validators[0], common.Big0, bondingAmount, "", true})
			extraBonds = append(extraBonds, StakingRequest{user, validators[0], common.Big1, bondingAmount, "", true})
			extraBonds = append(extraBonds, StakingRequest{user, validators[1], common.Big0, bondingAmount, "", true})
			extraBonds = append(extraBonds, StakingRequest{user, validators[0], common.Big0, bondingAmount, "", true})
		}
		// dummy
		extraBonds = append(extraBonds, StakingRequest{common.Address{}, validators[0], common.Big0, common.Big0, "", true})

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidStateContracts, users, validators)

		require.True(
			r.T,
			isAllRewardsZero(r, contractCount, liquidStateContracts, users, validators),
			"all rewards should be zero initially",
		)

		for _, request := range extraBonds {

			user := request.staker
			if request.amount.Cmp(common.Big0) > 0 {
				stakableContract := r.StakableVestingContractObject(user, request.contractID)
				r.NoError(
					stakableContract.Bond(
						tests.FromSender(user, nil), request.validator, request.amount,
					),
				)
			}

			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldUserRewards := unclaimedRewards(r, contractCount, users, validators)
			totalReward := r.RewardsAfterOneEpoch()
			r.WaitNextEpoch()

			// request is not applied yet
			checkRewards(
				r, contractCount, totalStake, totalReward,
				validators, users, validatorStakes,
				userStakes, oldUserRewards,
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

	tests.RunWithSetup("multiple bonding in same epoch and track rewards", setup, func(r *tests.Runner) {
		users, validators, liquidStateContracts := initiate2(r)
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

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidStateContracts, users, validators)

		require.True(
			r.T,
			isAllRewardsZero(r, contractCount, liquidStateContracts, users, validators),
			"all rewards should be zero initially",
		)

		for _, requests := range extraBondsArray {
			for _, request := range requests {

				user := request.staker
				stakableContract := r.StakableVestingContractObject(user, request.contractID)
				r.NoError(
					stakableContract.Bond(
						tests.FromSender(user, nil), request.validator, request.amount,
					),
				)
			}

			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldUserRewards := unclaimedRewards(r, contractCount, users, validators)
			totalReward := r.RewardsAfterOneEpoch()
			r.WaitNextEpoch()

			// request is not applied yet
			checkRewards(
				r, contractCount, totalStake, totalReward,
				validators, users, validatorStakes,
				userStakes, oldUserRewards,
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

		validatorNewStakes, userNewStakes, totalNewStake := initialStakes(r, contractCount, liquidStateContracts, users, validators)
		require.Equal(r.T, totalStake, totalNewStake)
		for validator, newStake := range validatorNewStakes {
			require.Equal(r.T, validatorStakes[validator], newStake)
		}
		for user, newStake := range userNewStakes {
			require.Equal(r.T, userStakes[user], newStake)
		}
	})

	// bond everything
	initiate3 := func(r *tests.Runner) (
		users, validators []common.Address,
		liquidStateContracts []*tests.ILiquidLogic,
	) {
		users, validators, liquidStateContracts = initiate2(r)
		bondingPerContract := new(big.Int).Mul(bondingAmount, big.NewInt(int64(len(validators))))
		remainingNTN := new(big.Int).Sub(big.NewInt(contractTotalAmount), bondingPerContract)
		newBondingAmount := new(big.Int).Div(remainingNTN, big.NewInt(int64(len(validators))))
		for _, user := range users {
			for i := 0; i < contractCount; i++ {
				stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
				for _, validator := range validators {
					r.NoError(
						stakableContract.Bond(
							tests.FromSender(user, nil), validator, newBondingAmount,
						),
					)
				}
			}
		}

		r.WaitNextEpoch()
		return
	}

	tests.RunWithSetup("release liquid and track reward", setup, func(r *tests.Runner) {
		users, validators, liquidStateContracts := initiate3(r)
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

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidStateContracts, users, validators)

		liquidStateContractsMap := make(map[common.Address]*tests.ILiquidLogic)

		for i, liquidStateContract := range liquidStateContracts {
			liquidStateContractsMap[validators[i]] = liquidStateContract
		}

		for _, request := range releaseRequests {

			// some epoch is passed and we are entitled to some reward,
			// but we don't know about it because we did not get notified
			// or we did not claim them or call unclaimedRewards
			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldUserRewards := unclaimedRewards(r, contractCount, users, validators)
			totalReward := r.RewardsAfterOneEpoch()
			r.WaitNextEpoch()

			// we release some LNTN and it is applied immediately
			// if unlocked, it is transferred immediately
			// but for reward calculation, it will be applied later
			user := request.staker
			amount := request.amount
			validator := request.validator
			if request.amount.Cmp(common.Big0) > 0 {
				stakableContract := r.StakableVestingContractObject(user, request.contractID)
				r.NoError(
					stakableContract.ReleaseLNTN(
						tests.FromSender(user, nil),
						request.validator,
						request.amount,
					),
				)

				userLiquidBalance[user][validator].Add(userLiquidBalance[user][validator], amount)
				balance, _, err := liquidStateContractsMap[validator].BalanceOf(nil, user)
				require.NoError(r.T, err)
				require.Equal(r.T, userLiquidBalance[user][validator], balance, "liquid not transferred")
			}

			checkRewards(
				r, contractCount, totalStake, totalReward,
				validators, users, validatorStakes,
				userStakes, oldUserRewards,
			)

			// for next reward
			if request.amount.Cmp(common.Big0) > 0 {
				id := int(request.contractID.Int64())
				validatorStakes[validator].Sub(validatorStakes[validator], amount)
				userStakes[user][id][validator].Sub(userStakes[user][id][validator], amount)
			}
		}
	})

	tests.RunWithSetup("unbond in different epoch and track reward", setup, func(r *tests.Runner) {
		users, validators, liquidStateContracts := initiate3(r)
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

		validatorStakes, userStakes, totalStake := initialStakes(r, contractCount, liquidStateContracts, users, validators)

		for _, request := range extraUnbonds {

			user := request.staker
			if request.amount.Cmp(common.Big0) > 0 {
				stakableContract := r.StakableVestingContractObject(user, request.contractID)
				r.NoError(
					stakableContract.Unbond(
						tests.FromSender(user, nil), request.validator, request.amount,
					),
				)
			}

			r.GiveMeSomeMoney(r.Autonity.Address(), reward)
			oldUserRewards := unclaimedRewards(r, contractCount, users, validators)
			totalReward := r.RewardsAfterOneEpoch()
			r.WaitNextEpoch()
			// request is not applied yet
			checkRewards(
				r, contractCount, totalStake, totalReward,
				validators, users, validatorStakes,
				userStakes, oldUserRewards,
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
	var contractTotalAmount int64 = 1000
	start := 100 + time.Now().Unix()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	user := tests.User
	contractID := common.Big0
	newUser := common.HexToAddress("0x88")
	require.NotEqual(t, user, newUser)

	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	initiate := func(r *tests.Runner) *tests.StakableVesting {
		createContract(r, user, contractTotalAmount, start, cliff, end)
		return r.StakableVestingContractObject(user, contractID)
	}

	tests.RunWithSetup("beneficiary changes", setup, func(r *tests.Runner) {
		stakableContract := initiate(r)
		_, _, err := stakableContract.GetContract(nil)
		require.NoError(r.T, err)
		beneficiary, _, err := stakableContract.Beneficiary(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, user, beneficiary)
		_, _, err = r.StakableVestingManager.GetContractAccount(nil, newUser, contractID)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: invalid contract id", err.Error())
		r.NoError(
			r.StakableVestingManager.ChangeContractBeneficiary(operator, user, contractID, newUser),
		)
		beneficiary, _, err = stakableContract.Beneficiary(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, newUser, beneficiary)
		_, _, err = r.StakableVestingManager.GetContractAccount(nil, user, contractID)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: invalid contract id", err.Error())
	})

	tests.RunWithSetup("changing beneficiary transfers rewards to the old beneficiary", setup, func(r *tests.Runner) {
		stakableContract := initiate(r)
		bondingAmount := big.NewInt(contractTotalAmount)
		validator := r.Committee.Validators[0].NodeAddress
		r.WaitSomeBlock(start + 1)
		r.NoError(
			stakableContract.Bond(
				tests.FromSender(user, nil), validator, bondingAmount,
			),
		)
		r.WaitNextEpoch()
		r.GiveMeSomeMoney(r.Autonity.Address(), reward)
		r.WaitNextEpoch()
		rewards, _, err := stakableContract.UnclaimedRewards0(nil)
		require.NoError(r.T, err)
		atnBalance := r.GetBalanceOf(user)
		ntnBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)

		// change beneficiary
		r.NoError(
			r.StakableVestingManager.ChangeContractBeneficiary(operator, user, contractID, newUser),
		)
		newAtnBalance := r.GetBalanceOf(user)
		require.Equal(r.T, new(big.Int).Add(rewards.AtnRewards, atnBalance), newAtnBalance)
		newNtnBalance, _, err := r.Autonity.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(rewards.NtnRewards, ntnBalance), newNtnBalance)
	})
}

func TestSlashingAffect(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}
	// TODO (tariq): complete tests.Setup

	tests.RunWithSetup("contract total value update when bonded validator slashed", setup, func(r *tests.Runner) {
		// TODO (tariq): complete
	})
}

func TestAccessRestriction(t *testing.T) {
	user := tests.User

	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("only operator can create contract", setup, func(r *tests.Runner) {
		amount := big.NewInt(1000)
		start := new(big.Int).Add(big.NewInt(100), r.Evm.Context.Time)
		cliff := new(big.Int).Add(start, big.NewInt(100))
		end := new(big.Int).Add(start, amount)
		_, err := r.StakableVestingManager.NewContract(
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
	start := time.Now().Unix()
	cliff := 500 + start
	// by making (end - start == contractTotalAmount) we have (totalUnlocked = currentTime - start)
	end := contractTotalAmount + start
	contractID := common.Big0

	newSetup := func() *tests.Runner {
		r := setup()
		createContract(r, user, contractTotalAmount, start, cliff, end)
		return r
	}

	tests.RunWithSetup("only operator can change contract beneficiary", newSetup, func(r *tests.Runner) {
		newUser := common.HexToAddress("0x88")
		require.NotEqual(r.T, user, newUser)
		_, err := r.StakableVestingManager.ChangeContractBeneficiary(
			tests.FromSender(user, nil),
			user,
			contractID,
			newUser,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: caller is not the operator", err.Error())

		_, err = r.StakableVestingManager.ChangeContractBeneficiary(
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
	liquidStateContracts []*tests.ILiquidLogic,
	users, validators []common.Address,
) (
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[int]map[common.Address]*big.Int,
	totalStake *big.Int,
) {

	// need to update funds before querying the initial stakes
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
			r.NoError(
				stakableContract.UpdateFunds(nil),
			)
		}
	}

	totalStake = new(big.Int)

	validatorStakes = make(map[common.Address]*big.Int)
	for _, validator := range validators {
		validatorStakes[validator] = big.NewInt(0)
	}
	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
			for i, validator := range validators {
				liquidStateContract := liquidStateContracts[i]
				balance, _, err := liquidStateContract.BalanceOf(nil, stakableContract.Address())
				require.NoError(r.T, err)
				validatorStakes[validator].Add(validatorStakes[validator], balance)
			}
		}
	}

	userStakes = make(map[common.Address]map[int]map[common.Address]*big.Int)
	for _, user := range users {
		userStakes[user] = make(map[int]map[common.Address]*big.Int)
		for i := 0; i < contractCount; i++ {
			stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
			userStakes[user][i] = make(map[common.Address]*big.Int)
			for _, validator := range validators {
				balance, _, err := stakableContract.LiquidBalance(nil, validator)
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
	users, validators []common.Address,
) (
	oldUserRewards map[common.Address]map[int]map[common.Address]Reward,
) {

	oldUserRewards = make(map[common.Address]map[int]map[common.Address]Reward)
	for _, user := range users {
		oldUserRewards[user] = make(map[int]map[common.Address]Reward)
		for i := 0; i < contractCount; i++ {
			oldUserRewards[user][i] = make(map[common.Address]Reward)
			stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
			for _, validator := range validators {
				unclaimedReward, _, err := stakableContract.UnclaimedRewards(nil, validator)
				require.NoError(r.T, err)
				oldUserRewards[user][i][validator] = Reward{unclaimedReward.AtnRewards, unclaimedReward.NtnRewards}
			}
		}
	}

	return oldUserRewards
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
	totalReward tests.EpochReward,
	validators, users []common.Address,
	validatorStakes map[common.Address]*big.Int,
	userStakes map[common.Address]map[int]map[common.Address]*big.Int,
	oldUserRewards map[common.Address]map[int]map[common.Address]Reward,
) {

	currentRewards := make(map[common.Address]Reward)
	// check total rewards from each validator
	for _, validator := range validators {
		validatorTotalRewardATN := new(big.Int).Mul(validatorStakes[validator], totalReward.RewardATN)
		validatorTotalRewardNTN := new(big.Int).Mul(validatorStakes[validator], totalReward.RewardNTN)

		if totalStake.Cmp(common.Big0) != 0 {
			validatorTotalRewardATN = validatorTotalRewardATN.Div(validatorTotalRewardATN, totalStake)
			validatorTotalRewardNTN = validatorTotalRewardNTN.Div(validatorTotalRewardNTN, totalStake)
		}

		currentRewards[validator] = Reward{
			validatorTotalRewardATN,
			validatorTotalRewardNTN,
		}
	}

	// check each user rewards
	for _, user := range users {

		for i := 0; i < contractCount; i++ {
			unclaimedRewardForContractATN := new(big.Int)
			unclaimedRewardForContractNTN := new(big.Int)
			stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
			for _, validator := range validators {
				calculatedRewardATN := new(big.Int).Mul(userStakes[user][i][validator], currentRewards[validator].RewardATN)
				calculatedRewardNTN := new(big.Int).Mul(userStakes[user][i][validator], currentRewards[validator].RewardNTN)

				if validatorStakes[validator].Cmp(common.Big0) != 0 {
					calculatedRewardATN.Div(calculatedRewardATN, validatorStakes[validator])
					calculatedRewardNTN.Div(calculatedRewardNTN, validatorStakes[validator])
				}
				calculatedRewardATN.Add(calculatedRewardATN, oldUserRewards[user][i][validator].RewardATN)

				calculatedRewardNTN.Add(calculatedRewardNTN, oldUserRewards[user][i][validator].RewardNTN)

				unclaimedReward, _, err := stakableContract.UnclaimedRewards(nil, validator)
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
								stakableContract.ClaimRewards0(tests.FromSender(user, nil), validator),
							)
						},
					)

				})
			}

			unclaimedReward, _, err := stakableContract.UnclaimedRewards0(nil)
			require.NoError(r.T, err)

			require.Equal(r.T, unclaimedRewardForContractATN, unclaimedReward.AtnRewards)
			require.Equal(r.T, unclaimedRewardForContractNTN, unclaimedReward.NtnRewards)

			// so that following code snippet reverts
			r.RunAndRevert(func(r *tests.Runner) {
				checkClaimRewardsFunction(
					r, user, unclaimedReward.AtnRewards, unclaimedReward.NtnRewards,
					func() {
						r.NoError(
							stakableContract.ClaimRewards(tests.FromSender(user, nil)),
						)
					},
				)
			})
		}
	}
}

func isAllRewardsZero(
	r *tests.Runner,
	contractCount int,
	liquidStateContracts []*tests.ILiquidLogic,
	users, validators []common.Address,
) bool {

	for _, user := range users {
		for i := 0; i < contractCount; i++ {
			stakableContract := r.StakableVestingContractObject(user, big.NewInt(int64(i)))
			for _, validator := range validators {
				rewards, _, err := stakableContract.UnclaimedRewards(nil, validator)
				require.NoError(r.T, err)

				if rewards.AtnRewards.Cmp(common.Big0) != 0 {
					return false
				}

				if rewards.NtnRewards.Cmp(common.Big0) != 0 {
					return false
				}
			}

			for _, liquidStateContract := range liquidStateContracts {
				rewards, _, err := liquidStateContract.UnclaimedRewards(nil, stakableContract.Address())
				require.NoError(r.T, err)

				if rewards.UnclaimedATN.Cmp(common.Big0) != 0 {
					return false
				}

				if rewards.UnclaimedNTN.Cmp(common.Big0) != 0 {
					return false
				}
			}

			rewards, _, err := stakableContract.UnclaimedRewards0(nil)
			require.NoError(r.T, err)

			if rewards.AtnRewards.Cmp(common.Big0) != 0 {
				return false
			}

			if rewards.NtnRewards.Cmp(common.Big0) != 0 {
				return false
			}
		}
	}
	return true
}

func setupContracts(
	r *tests.Runner, contractCount, validatorCount int, contractTotalAmount, start, cliff, end int64,
) (users, validators []common.Address, liquidStateContracts []*tests.ILiquidLogic) {
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
	liquidStateContracts = make([]*tests.ILiquidLogic, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validators[i] = r.Committee.Validators[i].NodeAddress
		liquidStateContracts[i] = r.Committee.LiquidStateContracts[i]
	}
	return
}

func createContract(r *tests.Runner, beneficiary common.Address, amount, startTime, cliffTime, endTime int64) {
	startBig := big.NewInt(startTime)
	cliffBig := big.NewInt(cliffTime)
	endBig := big.NewInt(endTime)
	r.NoError(
		r.StakableVestingManager.NewContract(
			operator, beneficiary, big.NewInt(amount), big.NewInt(startTime),
			new(big.Int).Sub(cliffBig, startBig), new(big.Int).Sub(endBig, startBig),
		),
	)
}

func checkReleaseAllNTN(r *tests.Runner, user common.Address, contractID, releaseAmount *big.Int) {

	stakableContract := r.StakableVestingContractObject(user, contractID)
	contract, _, err := stakableContract.GetContract(nil)
	require.NoError(r.T, err)
	contractNTN := contract.CurrentNTNAmount
	withdrawn := contract.WithdrawnValue
	initBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)
	totalUnlocked, _, err := stakableContract.UnlockedFunds(nil)
	require.NoError(r.T, err)
	r.NoError(
		stakableContract.ReleaseAllNTN(tests.FromSender(user, nil)),
	)
	newBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)
	require.Equal(r.T, new(big.Int).Add(initBalance, releaseAmount), newBalance, "balance mismatch")
	contract, _, err = stakableContract.GetContract(nil)
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

	remainingUnlocked, _, err := stakableContract.UnlockedFunds(nil)
	require.NoError(r.T, err)
	require.True(r.T, new(big.Int).Sub(totalUnlocked, releaseAmount).Cmp(remainingUnlocked) == 0)

	// all unlocked NTN withdrawn
	require.True(r.T, contract.CurrentNTNAmount.Cmp(common.Big0) == 0 || remainingUnlocked.Cmp(common.Big0) == 0)
}

func checkReleaseAllLNTN(r *tests.Runner, user common.Address, contractID, releaseAmount *big.Int) {

	stakableContract := r.StakableVestingContractObject(user, contractID)
	totalUnlocked, _, err := stakableContract.UnlockedFunds(nil)
	require.NoError(r.T, err)

	bondedValidators, _, err := stakableContract.GetLinkedValidators(nil)
	require.NoError(r.T, err)

	userLiquidBalances := make([]*big.Int, 0)
	vaultLiquidBalances := make([]*big.Int, 0)
	userLiquidInVesting := make([]*big.Int, 0)
	for _, validator := range bondedValidators {
		liquidStateContract := r.LiquidStateContract(validator)
		balance, _, err := liquidStateContract.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		vaultLiquidBalances = append(vaultLiquidBalances, balance)

		balance, _, err = liquidStateContract.BalanceOf(nil, user)
		require.NoError(r.T, err)
		userLiquidBalances = append(userLiquidBalances, balance)

		balance, _, err = stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		userLiquidInVesting = append(userLiquidInVesting, balance)
	}

	r.NoError(
		stakableContract.ReleaseAllLNTN(tests.FromSender(user, nil)),
	)

	remainingUnlocked, _, err := stakableContract.UnlockedFunds(nil)
	require.NoError(r.T, err)
	require.True(r.T, new(big.Int).Sub(totalUnlocked, releaseAmount).Cmp(remainingUnlocked) == 0)

	for i, validator := range bondedValidators {
		released := releaseAmount
		if releaseAmount.Cmp(vaultLiquidBalances[i]) > 0 {
			released = vaultLiquidBalances[i]
		}
		liquidStateContract := r.LiquidStateContract(validator)
		balance, _, err := liquidStateContract.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Sub(vaultLiquidBalances[i], released).Cmp(balance) == 0)

		balance, _, err = liquidStateContract.BalanceOf(nil, user)
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Add(userLiquidBalances[i], released).Cmp(balance) == 0)

		balance, _, err = stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(r.T, new(big.Int).Sub(userLiquidInVesting[i], released).Cmp(balance) == 0)
	}
}

func checkReleaseLNTN(r *tests.Runner, user, validator common.Address, contractID, releaseAmount *big.Int) {

	stakableContract := r.StakableVestingContractObject(user, contractID)
	var liquidStateContract *tests.ILiquidLogic
	for i, v := range r.Committee.Validators {
		if v.NodeAddress == validator {
			liquidStateContract = r.Committee.LiquidStateContracts[i]
			break
		}
	}

	liquidBalance, _, err := liquidStateContract.BalanceOf(nil, user)
	require.NoError(r.T, err)

	liquidInVesting, _, err := stakableContract.LiquidBalance(nil, validator)
	require.NoError(r.T, err)

	totalUnlocked, _, err := stakableContract.UnlockedFunds(nil)
	require.NoError(r.T, err)

	r.NoError(
		stakableContract.ReleaseLNTN(tests.FromSender(user, nil), validator, releaseAmount),
	)

	newLiquidBalance, _, err := liquidStateContract.BalanceOf(nil, user)
	require.NoError(r.T, err)
	require.Equal(
		r.T,
		new(big.Int).Add(liquidBalance, releaseAmount),
		newLiquidBalance,
	)

	newLiquidInVesting, _, err := stakableContract.LiquidBalance(nil, validator)
	require.NoError(r.T, err)
	require.True(
		r.T,
		newLiquidInVesting.Cmp(new(big.Int).Sub(liquidInVesting, releaseAmount)) == 0,
	)

	remainingUnlocked, _, err := stakableContract.UnlockedFunds(nil)
	require.NoError(r.T, err)
	require.True(r.T, new(big.Int).Sub(totalUnlocked, releaseAmount).Cmp(remainingUnlocked) == 0)
}

func isInitialBalanceZero(
	r *tests.Runner,
	requests []StakingRequest,
) bool {

	updateVestingContractFunds(r, requests)

	for _, request := range requests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		balance, _, err := stakableContract.LiquidBalance(nil, request.validator)
		require.NoError(r.T, err)
		if balance.Cmp(common.Big0) != 0 {
			return false
		}
	}

	for _, request := range requests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		for _, liquidStateContract := range r.Committee.LiquidStateContracts {
			balance, _, err := liquidStateContract.BalanceOf(nil, stakableContract.Address())
			require.NoError(r.T, err)
			if balance.Cmp(common.Big0) != 0 {
				return false
			}
		}
	}

	return true

}

func initialBalances(
	r *tests.Runner,
	stakingRequests []StakingRequest,
) (
	liquidStateContracts map[common.Address]*tests.ILiquidLogic,
	// liquidOfVestingContract map[common.Address]*big.Int,
	ntnBalanceOfVestingContract map[common.Address]map[int64]*big.Int,
	liquidOfUser map[common.Address]map[common.Address]map[int64]*big.Int,
) {
	liquidStateContracts = make(map[common.Address]*tests.ILiquidLogic)
	// liquidOfVestingContract = make(map[common.Address]*big.Int)
	ntnBalanceOfVestingContract = make(map[common.Address]map[int64]*big.Int)
	liquidOfUser = make(map[common.Address]map[common.Address]map[int64]*big.Int)

	for _, request := range stakingRequests {
		liquidOfUser[request.staker] = make(map[common.Address]map[int64]*big.Int)
		ntnBalanceOfVestingContract[request.staker] = make(map[int64]*big.Int)
	}

	for _, request := range stakingRequests {
		liquidOfUser[request.staker][request.validator] = make(map[int64]*big.Int)
	}

	for i, validator := range r.Committee.Validators {
		for _, request := range stakingRequests {
			if request.validator == validator.NodeAddress {
				liquidStateContract := r.Committee.LiquidStateContracts[i]
				liquidStateContracts[request.validator] = liquidStateContract
				break
			}
		}
	}

	for _, request := range stakingRequests {
		stakableVesting := r.StakableVestingContractObject(request.staker, request.contractID)
		userLiquid, _, err := stakableVesting.LiquidBalance(nil, request.validator)
		require.NoError(r.T, err)
		liquidOfUser[request.staker][request.validator][request.contractID.Int64()] = userLiquid

		balance, _, err := r.Autonity.BalanceOf(nil, stakableVesting.Address())
		require.NoError(r.T, err)
		ntnBalanceOfVestingContract[request.staker][request.contractID.Int64()] = balance
	}
	return liquidStateContracts, ntnBalanceOfVestingContract, liquidOfUser
}

func updateVestingContractFunds(r *tests.Runner, stakingRequests []StakingRequest) {
	for _, request := range stakingRequests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		r.NoError(
			stakableContract.UpdateFunds(nil),
		)
	}
}

func bondAndFinalize(
	r *tests.Runner, bondingRequests []StakingRequest,
) {

	liquidStateContracts, ntnBalanceOfContract, liquidOfUser := initialBalances(r, bondingRequests)

	for _, request := range bondingRequests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		contract, _, err := stakableContract.GetContract(nil)
		require.NoError(r.T, err)
		contractNTN := contract.CurrentNTNAmount

		_, bondErr := stakableContract.Bond(
			tests.FromSender(request.staker, nil),
			request.validator,
			request.amount,
		)

		contract, _, err = stakableContract.GetContract(nil)
		require.NoError(r.T, err)
		remaining := new(big.Int).Sub(contractNTN, common.Big0)

		if request.expectedErr == "" {
			require.NoError(r.T, bondErr)
			validator := request.validator
			id := request.contractID.Int64()
			liquidOfUser[request.staker][validator][id].Add(liquidOfUser[request.staker][validator][id], request.amount)

			remaining.Sub(remaining, request.amount)

			currentBalance := ntnBalanceOfContract[request.staker][id]
			ntnBalanceOfContract[request.staker][id] = new(big.Int).Sub(currentBalance, request.amount)
		} else {
			require.Error(r.T, bondErr)
			require.Equal(r.T, request.expectedErr, bondErr.Error())
		}
	}

	// let bonding apply
	r.WaitNextEpoch()

	// need to update funds in vesting contract
	updateVestingContractFunds(r, bondingRequests)

	for _, request := range bondingRequests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		validator := request.validator
		id := request.contractID.Int64()

		userLiquid, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			liquidOfUser[request.staker][validator][id].Cmp(userLiquid) == 0,
			"vesting contract cannot track liquid balance",
		)

		liquidStateContract := liquidStateContracts[validator]
		contractLiquid, _, err := liquidStateContract.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.True(
			r.T,
			userLiquid.Cmp(contractLiquid) == 0,
			"liquid from Liquid State contract does not match",
		)

		contract, _, err := stakableContract.GetContract(nil)
		require.NoError(r.T, err)

		newNewtonBalance, _, err := r.Autonity.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.True(
			r.T,
			contract.CurrentNTNAmount.Cmp(newNewtonBalance) == 0,
			"contract NTN should match contract balance",
		)
		require.True(
			r.T,
			newNewtonBalance.Cmp(ntnBalanceOfContract[request.staker][id]) == 0,
			"newton balance not updated",
		)

	}
}

func unbondAndRelease(
	r *tests.Runner, unbondingRequests []StakingRequest,
) {
	liquidStateContracts, ntnBalanceOfContract, liquidOfUser := initialBalances(r, unbondingRequests)

	unbondingRequestBlock := r.Evm.Context.BlockNumber

	for _, request := range unbondingRequests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		lockedLiquid, _, err := stakableContract.LockedLiquidBalance(nil, request.validator)
		require.NoError(r.T, err)
		unlockedLiquid, _, err := stakableContract.UnlockedLiquidBalance(nil, request.validator)
		require.NoError(r.T, err)
		_, unbondErr := stakableContract.Unbond(
			tests.FromSender(request.staker, nil),
			request.validator,
			request.amount,
		)

		if request.expectedErr == "" {
			require.NoError(r.T, unbondErr)
			validator := request.validator
			id := request.contractID.Int64()
			liquidOfUser[request.staker][validator][id].Sub(liquidOfUser[request.staker][validator][id], request.amount)

			newLockedLiquid, _, err := stakableContract.LockedLiquidBalance(nil, request.validator)
			require.NoError(r.T, err)
			require.True(
				r.T,
				new(big.Int).Add(lockedLiquid, request.amount).Cmp(newLockedLiquid) == 0,
				"vesting contract cannot track locked liquid",
			)

			newUnlockedLiquid, _, err := stakableContract.UnlockedLiquidBalance(nil, request.validator)
			require.NoError(r.T, err)
			require.True(
				r.T,
				new(big.Int).Sub(unlockedLiquid, request.amount).Cmp(newUnlockedLiquid) == 0,
				"vesting contract cannot track unlocked liquid",
			)

			currentBalance := ntnBalanceOfContract[request.staker][id]
			ntnBalanceOfContract[request.staker][id] = new(big.Int).Add(currentBalance, request.amount)
		} else {
			require.Error(r.T, unbondErr)
			require.Equal(r.T, request.expectedErr, unbondErr.Error())
		}
	}

	r.WaitNextEpoch()
	// need to update funds in vesting contract
	updateVestingContractFunds(r, unbondingRequests)

	for _, request := range unbondingRequests {
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		validator := request.validator
		id := request.contractID.Int64()

		userLiquid, _, err := stakableContract.LiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			userLiquid.Cmp(liquidOfUser[request.staker][validator][id]) == 0,
			"vesting contract cannot track liquid",
		)

		liquidStateContract := liquidStateContracts[validator]
		contractLiquid, _, err := liquidStateContract.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.True(
			r.T,
			userLiquid.Cmp(contractLiquid) == 0,
			"liquid from Liquid State does not match with stakable contract",
		)

		lockedLiquid, _, err := stakableContract.LockedLiquidBalance(nil, validator)
		require.NoError(r.T, err)
		require.True(
			r.T,
			lockedLiquid.Cmp(common.Big0) == 0,
			"vesting contract cannot track locked liquid",
		)

		unlockedLiquid, _, err := stakableContract.UnlockedLiquidBalance(nil, validator)
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
		stakableContract := r.StakableVestingContractObject(request.staker, request.contractID)
		contract, _, err := stakableContract.GetContract(nil)
		require.NoError(r.T, err)

		newNewtonBalance, _, err := r.Autonity.BalanceOf(nil, stakableContract.Address())
		require.NoError(r.T, err)
		require.Equal(
			r.T,
			contract.CurrentNTNAmount, newNewtonBalance,
			"contract NTN should match contract ntn balance",
		)
		id := request.contractID.Int64()
		require.True(
			r.T,
			newNewtonBalance.Cmp(ntnBalanceOfContract[request.staker][id]) == 0,
			"vesting contract balance mismatch",
		)
	}
}
