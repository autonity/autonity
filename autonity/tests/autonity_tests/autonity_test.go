package autonity_tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

const (
	FailAlways = iota
	FailAfterRewardsDistribution
	FailAfterUnboningApplied
)

type StakingRequest struct {
	staker    common.Address
	amount    *big.Int
	validator common.Address
	bond      bool
}

func TestBondingOperationReverted(t *testing.T) {
	r := tests.Setup(t, nil)
	stakingContract := deployDummyStakingContract(r)
	validator := r.Committee.Validators[0].NodeAddress
	bondingAmount := big.NewInt(100)

	stakingContracts := make(map[common.Address]*tests.DummyStakintgContract)
	stakingContracts[stakingContract.Address()] = stakingContract

	r.Run("reverts bonding, single request", func(r *tests.Runner) {
		testBondingReverted(r, stakingContracts, []StakingRequest{{stakingContract.Address(), bondingAmount, validator, true}})
	})

	extraContract := 2
	for extraContract > 0 {
		extraContract--
		contract := deployDummyStakingContract(r)
		_, ok := stakingContracts[contract.Address()]
		require.False(r.T, ok, "contract already exists")
		stakingContracts[contract.Address()] = contract
	}

	r.Run("reverts bonding, multiple request from single contract", func(r *tests.Runner) {
		// TODO (tariq): complete
	})

	r.Run("reverts bonding, multiple request from multiple contract", func(r *tests.Runner) {
		// TODO (tariq): complete
	})
}

func TestUnbondingOperationReverted(t *testing.T) {
	r := tests.Setup(t, nil)
	_, _, stakingContract, err := r.DeployDummyStakintgContract(nil, r.Autonity.Address())
	require.NoError(r.T, err)
	r.GiveMeSomeMoney(stakingContract.Address(), big.NewInt(1000_000_000_000_000_000))
	validator := r.Committee.Validators[0].NodeAddress
	liquidContract := r.Committee.LiquidContracts[0]
	bondingAmount := big.NewInt(100)
	r.NoError(
		r.Autonity.Mint(tests.Operator, stakingContract.Address(), bondingAmount),
	)
	// bond to test unbonding
	balance, _, err := r.Autonity.BalanceOf(nil, stakingContract.Address())
	require.NoError(r.T, err)
	liquidBalance, _, err := liquidContract.BalanceOf(nil, stakingContract.Address())
	require.NoError(r.T, err)
	r.NoError(
		stakingContract.Bond(nil, validator, bondingAmount),
	)
	r.WaitNextEpoch()

	newBalance, _, err := r.Autonity.BalanceOf(nil, stakingContract.Address())
	require.NoError(r.T, err)
	require.True(r.T, newBalance.Cmp(new(big.Int).Sub(balance, bondingAmount)) == 0)
	newLiquidBalance, _, err := liquidContract.BalanceOf(nil, stakingContract.Address())
	require.NoError(r.T, err)
	require.True(r.T, newLiquidBalance.Cmp(new(big.Int).Add(liquidBalance, bondingAmount)) == 0)

	stakingContracts := make(map[common.Address]*tests.DummyStakintgContract)
	stakingContracts[stakingContract.Address()] = stakingContract

	r.Run("revert unbonding, single request", func(r *tests.Runner) {
		testUnbondingReverted(r, stakingContracts, []StakingRequest{{stakingContract.Address(), bondingAmount, validator, false}})
	})

	extraContract := 2
	for extraContract > 0 {
		extraContract--
		contract := deployDummyStakingContract(r)
		_, ok := stakingContracts[contract.Address()]
		require.False(r.T, ok, "contract already exists")
		stakingContracts[contract.Address()] = contract
	}

	r.Run("reverts unbonding, multiple request from single contract", func(r *tests.Runner) {
		// TODO (tariq): complete
	})

	r.Run("reverts unbonding, multiple request from multiple contract", func(r *tests.Runner) {
		// TODO (tariq): complete
	})

}

func deployDummyStakingContract(r *tests.Runner) *tests.DummyStakintgContract {
	_, _, contract, err := r.DeployDummyStakintgContract(nil, r.Autonity.Address())
	require.NoError(r.T, err)
	r.GiveMeSomeMoney(contract.Address(), big.NewInt(1000_000_000_000_000_000))
	r.NoError(
		r.Autonity.Mint(tests.Operator, contract.Address(), big.NewInt(1000_000_000_000_000_000)),
	)
	return contract
}

func checkValidatorStake(r *tests.Runner, validator common.Address, oldValidatorInfo tests.AutonityValidator) {

	validatorInfo, _, err := r.Autonity.GetValidator(nil, validator)
	require.NoError(r.T, err)

	require.True(
		r.T,
		oldValidatorInfo.BondedStake.Cmp(validatorInfo.BondedStake) == 0,
		"bonded stake mismatch",
	)

	require.True(
		r.T,
		oldValidatorInfo.SelfBondedStake.Cmp(validatorInfo.SelfBondedStake) == 0,
		"self bonded stake mismatch",
	)

	require.True(
		r.T,
		oldValidatorInfo.LiquidSupply.Cmp(validatorInfo.LiquidSupply) == 0,
		"liquid supply mismatch",
	)
}

func checkValidatorAndLiquid(
	r *tests.Runner,
	liquidContracts map[common.Address]*tests.Liquid,
	stakingContracts map[common.Address]*tests.DummyStakintgContract,
	liquidBalances map[common.Address]map[common.Address]*big.Int,
	validatorInfos map[common.Address]tests.AutonityValidator,
) {
	for validator, liquidContract := range liquidContracts {
		for staker := range stakingContracts {
			newLiquidBalance, _, err := liquidContract.BalanceOf(nil, staker)
			require.NoError(r.T, err)
			require.True(
				r.T,
				liquidBalances[staker][validator].Cmp(newLiquidBalance) == 0,
				"liquid balance mismatch",
			)
		}

		checkValidatorStake(r, validator, validatorInfos[validator])
	}
}

func checkRewardNotification(
	r *tests.Runner,
	epochID *big.Int,
	stakingContracts map[common.Address]*tests.DummyStakintgContract,
	validatorStaked map[common.Address][]common.Address,
	failStep int,
) {
	for staker, validators := range validatorStaked {
		stakingContract := stakingContracts[staker]
		for _, validator := range validators {
			state, _, err := stakingContract.ValidatorsState(nil, epochID, validator)
			require.NoError(r.T, err)
			constant := new(big.Int)

			if failStep == FailAlways {
				constant, _, err = stakingContract.ValidatorStaked(nil)
				require.NoError(r.T, err)
			} else {
				constant, _, err = stakingContract.ValidatorRewarded(nil)
				require.NoError(r.T, err)
			}
			require.Equal(r.T, constant, state)
		}
	}
}

func testBondingReverted(r *tests.Runner, stakingContracts map[common.Address]*tests.DummyStakintgContract, requests []StakingRequest) {
	for _, stakingContract := range stakingContracts {
		r.NoError(
			stakingContract.RevertStakingOperations(nil),
		)
	}

	bondFromContractAndRevert(r, stakingContracts, requests, FailAlways)

	for _, stakingContract := range stakingContracts {
		stakingContract.RemoveRequestedBondingIDs(nil)
		r.NoError(
			stakingContract.ProcessStakingOperations(nil),
		)
		r.NoError(
			// the notification fails after rewards distribution
			stakingContract.FailAfterRewardsDistribution(nil),
		)
	}

	bondFromContractAndRevert(r, stakingContracts, requests, FailAfterRewardsDistribution)
}

func testUnbondingReverted(r *tests.Runner, stakingContracts map[common.Address]*tests.DummyStakintgContract, requests []StakingRequest) {
	for _, stakingContract := range stakingContracts {
		r.NoError(
			stakingContract.RevertStakingOperations(nil),
		)
	}

	unbondFromContractAndRevert(r, stakingContracts, requests, FailAlways)

	for _, stakingContract := range stakingContracts {
		stakingContract.RemoveRequestedUnbondingIDs(nil)
		r.NoError(
			stakingContract.ProcessStakingOperations(nil),
		)
		r.NoError(
			stakingContract.FailAfterRewardsDistribution(nil),
		)
	}

	unbondFromContractAndRevert(r, stakingContracts, requests, FailAfterRewardsDistribution)

	for _, stakingContract := range stakingContracts {
		stakingContract.RemoveRequestedUnbondingIDs(nil)
		r.NoError(
			stakingContract.ProcessStakingOperations(nil),
		)
		r.NoError(
			stakingContract.FailAfterUnbondingApplied(nil),
		)
	}

	unbondFromContractAndRevert(r, stakingContracts, requests, FailAfterUnboningApplied)
}

func bondFromContractAndRevert(r *tests.Runner, stakingContracts map[common.Address]*tests.DummyStakintgContract, requests []StakingRequest, failStep int) {

	balances, liquidBalances, liquidContracts, validatorInfos := initialState(r, stakingContracts, requests)

	validatorStaked := make(map[common.Address][]common.Address)
	for staker := range stakingContracts {
		validatorStaked[staker] = make([]common.Address, 0)
	}
	// make requests
	bondingIDs := make([]*big.Int, 0)
	for i, request := range requests {
		// get the bonding id
		staker := request.staker
		stakingContract := stakingContracts[staker]
		out, _, err := r.SimulateCall(stakingContract.Contract(), nil, "bond", request.validator, request.amount)
		require.NoError(r.T, err)
		bondingID := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
		bondingIDs = append(bondingIDs, bondingID)
		// checking if the above call reverted
		_, _, err = stakingContract.RequestedBondings(nil, big.NewInt(int64(i)))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: ", err.Error())

		r.NoError(
			stakingContract.Bond(nil, request.validator, request.amount),
		)
		storedBondingID, _, err := stakingContract.RequestedBondings(nil, big.NewInt(int64(i)))
		require.NoError(r.T, err)
		require.True(
			r.T,
			bondingID.Cmp(storedBondingID) == 0,
			"bonding id mismatch",
		)
		balances[staker].Sub(balances[staker], request.amount)
		validatorStaked[staker] = append(validatorStaked[staker], request.validator)
	}

	for staker := range stakingContracts {
		newBalance, _, err := r.Autonity.BalanceOf(nil, staker)
		require.NoError(r.T, err)
		require.True(
			r.T,
			balances[staker].Cmp(newBalance) == 0,
			"bonding not working",
		)
	}

	oldEpochID, _, err := r.Autonity.EpochID(nil)
	require.NoError(r.T, err)
	r.WaitNextEpoch()

	checkRewardNotification(r, oldEpochID, stakingContracts, validatorStaked, failStep)

	// check if requests are reverted
	for i, bondingID := range bondingIDs {
		request := requests[i]
		stakingContract := stakingContracts[request.staker]
		info, _, err := stakingContract.NotifiedBondings(nil, bondingID)
		require.NoError(r.T, err)

		require.Equal(r.T, false, info.Applied, "bonding applied")
		balances[request.staker].Add(balances[request.staker], request.amount)
	}

	for staker := range stakingContracts {
		newBalance, _, err := r.Autonity.BalanceOf(nil, staker)
		require.NoError(r.T, err)
		require.True(
			r.T,
			balances[staker].Cmp(newBalance) == 0,
			"balance not updating",
		)
	}

	checkValidatorAndLiquid(r, liquidContracts, stakingContracts, liquidBalances, validatorInfos)

}

func unbondFromContractAndRevert(r *tests.Runner, stakingContracts map[common.Address]*tests.DummyStakintgContract, requests []StakingRequest, failStep int) {

	balances, liquidBalances, liquidContracts, validatorInfos := initialState(r, stakingContracts, requests)
	validatorStaked := make(map[common.Address][]common.Address)
	for staker := range stakingContracts {
		validatorStaked[staker] = make([]common.Address, 0)
	}

	// make requests
	unbondingBlock := r.Evm.Context.BlockNumber
	unbondingIDs := make([]*big.Int, 0)
	for i, request := range requests {
		// get the unbonding id
		validator := request.validator
		staker := request.staker
		stakingContract := stakingContracts[staker]
		liquidContract := liquidContracts[validator]
		out, _, err := r.SimulateCall(stakingContract.Contract(), nil, "unbond", validator, request.amount)
		require.NoError(r.T, err)
		unbondingID := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
		unbondingIDs = append(unbondingIDs, unbondingID)
		// checking if the above call reverted
		_, _, err = stakingContract.RequestedUnbondings(nil, big.NewInt(int64(i)))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: ", err.Error())

		unlocked, _, err := liquidContract.UnlockedBalanceOf(nil, request.staker)
		require.NoError(r.T, err)
		r.NoError(
			stakingContract.Unbond(nil, validator, request.amount),
		)
		newUnlocked, _, err := liquidContract.UnlockedBalanceOf(nil, request.staker)
		require.NoError(r.T, err)
		require.True(
			r.T,
			newUnlocked.Cmp(new(big.Int).Sub(unlocked, request.amount)) == 0,
			"unbonding not working",
		)

		storedUnbondingID, _, err := stakingContract.RequestedUnbondings(nil, big.NewInt(int64(i)))
		require.NoError(r.T, err)
		require.True(
			r.T,
			unbondingID.Cmp(storedUnbondingID) == 0,
			"unbonding id mismatch",
		)

		if failStep == FailAfterUnboningApplied {
			// change in liquid balance will be seen after unbonding is applied
			liquidBalances[staker][validator].Sub(liquidBalances[staker][validator], request.amount)
		}

		validatorStaked[staker] = append(validatorStaked[staker], validator)
	}

	oldEpochID, _, err := r.Autonity.EpochID(nil)
	require.NoError(r.T, err)
	r.WaitNextEpoch()

	checkRewardNotification(r, oldEpochID, stakingContracts, validatorStaked, failStep)

	// check if unbonding is applied or reverted
	for i, unbondingID := range unbondingIDs {
		request := requests[i]
		staker := request.staker
		stakingContract := stakingContracts[staker]
		info, _, err := stakingContract.NotifiedUnbonding(nil, unbondingID)
		require.NoError(r.T, err)

		if failStep == FailAfterUnboningApplied {
			require.Equal(r.T, true, info.Applied, "unbonding not applied")
		} else {
			require.Equal(r.T, false, info.Applied, "unbonding applied")
			require.Equal(r.T, request.validator, info.Validator)
			require.Equal(r.T, false, info.Rejected)
		}

		validator := request.validator
		liquidContract := liquidContracts[validator]
		balance, _, err := liquidContract.BalanceOf(nil, request.staker)
		require.NoError(r.T, err)
		require.True(
			r.T,
			liquidBalances[staker][validator].Cmp(balance) == 0,
			"liquid balance not updated after unbonding applied",
		)

		unlocked, _, err := liquidContract.UnlockedBalanceOf(nil, request.staker)
		require.NoError(r.T, err)
		require.True(
			r.T,
			unlocked.Cmp(balance) == 0,
			"liquid not unlocked",
		)
	}

	if failStep == FailAfterUnboningApplied {
		// liquid balance was changed, but now it will revert
		for _, request := range requests {
			validator := request.validator
			staker := request.staker
			liquidBalances[staker][validator].Add(liquidBalances[staker][validator], request.amount)
		}
	}

	// release unbonding
	unbondingPeriod, _, err := r.Autonity.GetUnbondingPeriod(nil)
	require.NoError(r.T, err)
	targetBlock := new(big.Int).Add(unbondingBlock, unbondingPeriod)
	for r.Evm.Context.BlockNumber.Cmp(targetBlock) <= 0 {
		r.WaitNextEpoch()
	}

	for staker := range stakingContracts {
		newBalance, _, err := r.Autonity.BalanceOf(nil, staker)
		require.NoError(r.T, err)
		require.True(
			r.T,
			balances[staker].Cmp(newBalance) == 0,
			"balance should not update",
		)
	}

	checkValidatorAndLiquid(r, liquidContracts, stakingContracts, liquidBalances, validatorInfos)
}

func initialState(r *tests.Runner, stakingContracts map[common.Address]*tests.DummyStakintgContract, requests []StakingRequest) (
	balances map[common.Address]*big.Int,
	liquidBalances map[common.Address]map[common.Address]*big.Int,
	liquidContracts map[common.Address]*tests.Liquid,
	validatorInfos map[common.Address]tests.AutonityValidator,
) {

	balances = make(map[common.Address]*big.Int)
	for address := range stakingContracts {
		balance, _, err := r.Autonity.BalanceOf(nil, address)
		require.NoError(r.T, err)
		balances[address] = balance
	}

	liquidContracts = make(map[common.Address]*tests.Liquid)
	validatorInfos = make(map[common.Address]tests.AutonityValidator)
	liquidBalances = make(map[common.Address]map[common.Address]*big.Int)

	for address := range stakingContracts {
		liquidBalances[address] = make(map[common.Address]*big.Int)
	}

	for _, request := range requests {
		validator := request.validator
		for i, v := range r.Committee.Validators {
			if v.NodeAddress == validator {
				liquidContract := r.Committee.LiquidContracts[i]
				liquidContracts[validator] = liquidContract

				info, _, err := r.Autonity.GetValidator(nil, validator)
				require.NoError(r.T, err)
				validatorInfos[validator] = info

				break
			}
		}
	}

	for validator, liquidContract := range liquidContracts {
		for address := range stakingContracts {
			balance, _, err := liquidContract.BalanceOf(nil, address)
			require.NoError(r.T, err)
			liquidBalances[address][validator] = balance
		}
	}
	return
}
