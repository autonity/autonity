package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

var (
	staker1 = common.HexToAddress("0x1000000000000000000000000000000000000000")
	staker2 = common.HexToAddress("0x2000000000000000000000000000000000000000")
	staker3 = common.HexToAddress("0x3000000000000000000000000000000000000000")
)

func TestClaimRewards(t *testing.T) {
	// Test 1 validator 1 staker
	r := setup(t, nil)
	// Mint Newton to some few accounts
	r.autonity.Mint(operator, staker1, params.Ntn10000)
	r.autonity.Mint(operator, staker2, params.Ntn10000)
	r.autonity.Mint(operator, staker3, params.Ntn10000)
	r.autonity.Bond(&runOptions{origin: staker1}, r.committee.validators[0].NodeAddress, params.Ntn10000)
	r.autonity.Bond(&runOptions{origin: staker2}, r.committee.validators[1].NodeAddress, params.Ntn10000)
	r.autonity.Bond(&runOptions{origin: staker3}, r.committee.validators[1].NodeAddress, new(big.Int).Mul(common.Big2, params.Ntn10000))

	// create liquid staking contract per validator
	r.waitNextEpoch()
	// .. test here claiming rewards, checking if NTN/ATN reward is coherent and accurate.
	// transactions fees can be simulated be sending atns directly to the autonity contract account.
	// todo: Think about in base.go to assign at each epoch the current list of validators / committee
	// in r.validators with the liquid stake contract bindings already prepared so that's easy to manipulate
	// or maybe just create some helpers for it.
}

func TestFunctions(t *testing.T) {
	r := setup(t, nil)

	validator := r.committee.validators[0].NodeAddress
	treasury := r.committee.validators[0].Treasury
	delegatorA := r.committee.validators[1].NodeAddress
	delegatorB := r.committee.validators[2].NodeAddress
	delegatorC := r.committee.validators[3].NodeAddress
	liquidLogic := r.LiquidLogicContractObject()
	liquidState := deployLiquid(r, validator, treasury)

	validatorMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
	r.CallNoError(
		liquidState.CallMethod(
			liquidLogic.contract, fromAutonity, "mint", validator, validatorMint,
		),
	)

	supply, _, err := liquidState.TotalSupply(nil)
	require.NoError(r.t, err)
	require.Equal(r.t, validatorMint, supply)
	checkLiquidBalance(r, liquidState, validator, validatorMint)

	r.giveMeSomeMoney(r.autonity.address, new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor))

	r.run("check name and symbol", func(r *runner) {
		name, _, err := liquidState.Name(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, "LNTN-27", name)
		symbol, _, err := liquidState.Symbol(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, "LNTN-27", symbol)
	})

	r.run("reward single validator", func(r *runner) {
		// Initial state
		checkLiquidBalance(r, liquidState, delegatorA, common.Big0)
		checkLiquidBalance(r, liquidState, delegatorB, common.Big0)

		checkReward(r, liquidState, validator, common.Big0, common.Big0)
		checkReward(r, liquidState, delegatorA, common.Big0, common.Big0)
		checkReward(r, liquidState, delegatorB, common.Big0, common.Big0)

		// Send 10 ATN as a reward.  Perform a call first (not a tx)
		// in order to check the returned value.
		liquidReward := new(big.Int).Mul(big.NewInt(10), params.DecimalFactor)
		out, _ := r.CallNoError(
			liquidState.SimulateCall(
				liquidLogic.contract,
				fromSender(r.autonity.address, liquidReward),
				"redistribute", liquidReward,
			),
		)
		atnDistributed := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
		ntnDistributed := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
		require.True(r.t, liquidReward.Cmp(atnDistributed) >= 0)
		require.True(r.t, liquidReward.Cmp(ntnDistributed) >= 0)
		precision := new(big.Int).Div(params.DecimalFactor, big.NewInt(10000))
		require.True(r.t, new(big.Int).Sub(liquidReward, precision).Cmp(ntnDistributed) <= 0)
		require.True(r.t, new(big.Int).Sub(liquidReward, precision).Cmp(ntnDistributed) <= 0)

		redistributeLiquidReward(r, liquidState, liquidReward)

		// Check distribution (only validator should hold this)
		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, validatorMint, supply)
		checkLiquidBalance(r, liquidState, validator, validatorMint)
		checkLiquidBalance(r, liquidState, delegatorA, common.Big0)
		checkLiquidBalance(r, liquidState, delegatorB, common.Big0)

		checkReward(r, liquidState, validator, liquidReward, liquidReward)
		checkReward(r, liquidState, delegatorA, common.Big0, common.Big0)
		checkReward(r, liquidState, delegatorB, common.Big0, common.Big0)
	})

	r.run("reward multiple validators", func(r *runner) {
		// delegatorA bonds 8000 NEW
		// delegatorB bonds 2000 NEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		delegatorBMint := new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorB, delegatorBMint),
		)

		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Mul(big.NewInt(20000), params.DecimalFactor), supply)

		checkLiquidBalance(r, liquidState, validator, validatorMint)
		checkLiquidBalance(r, liquidState, delegatorA, delegatorAMint)
		checkLiquidBalance(r, liquidState, delegatorB, delegatorBMint)

		// Send 20 AUT as a reward and check distribution
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		expectedReward := new(big.Int).Mul(big.NewInt(10), params.DecimalFactor)
		checkReward(r, liquidState, validator, expectedReward, expectedReward)

		expectedReward = new(big.Int).Mul(big.NewInt(8), params.DecimalFactor)
		checkReward(r, liquidState, delegatorA, expectedReward, expectedReward)

		expectedReward = new(big.Int).Mul(big.NewInt(2), params.DecimalFactor)
		checkReward(r, liquidState, delegatorB, expectedReward, expectedReward)
	})

	r.run("transfer LNEW", func(r *runner) {
		// delegatorA bonds 8000 NEW
		// delegatorB bonds 2000 NEW
		// 20 AUT reward
		delegatorAMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		delegatorBMint := new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorB, delegatorBMint),
		)

		// Send 20 AUT as a reward and check distribution
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// delegatorA gives delegatorC 3000 LNEW
		transfer := new(big.Int).Mul(big.NewInt(3000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromSender(delegatorA, nil), "transfer", delegatorC, transfer),
		)

		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Mul(big.NewInt(20000), params.DecimalFactor), supply)

		checkLiquidBalance(r, liquidState, validator, validatorMint)
		checkLiquidBalance(r, liquidState, delegatorA, new(big.Int).Sub(delegatorAMint, transfer))
		checkLiquidBalance(r, liquidState, delegatorB, delegatorBMint)
		checkLiquidBalance(r, liquidState, delegatorC, transfer)

		// Another 20 AUT reward.  Check distribution.
		redistributeLiquidReward(r, liquidState, liquidReward)
		// validator has 10 + 10
		expectedReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		checkReward(r, liquidState, validator, expectedReward, expectedReward)
		// delegatorA has 8 + 5
		expectedReward = new(big.Int).Mul(big.NewInt(13), params.DecimalFactor)
		checkReward(r, liquidState, delegatorA, expectedReward, expectedReward)
		// delegatorB has 2 + 2
		expectedReward = new(big.Int).Mul(big.NewInt(4), params.DecimalFactor)
		checkReward(r, liquidState, delegatorB, expectedReward, expectedReward)
		// delegatorC has 3
		expectedReward = new(big.Int).Mul(big.NewInt(3), params.DecimalFactor)
		checkReward(r, liquidState, delegatorC, expectedReward, expectedReward)
	})

	r.run("burn LNEW", func(r *runner) {
		// delegatorA bonds 8000 NEW and burns 3000 LNEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		delegatorABurn := new(big.Int).Mul(big.NewInt(3000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "burn", delegatorA, delegatorABurn),
		)

		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Mul(big.NewInt(15000), params.DecimalFactor), supply)

		checkLiquidBalance(r, liquidState, validator, validatorMint)
		checkLiquidBalance(r, liquidState, delegatorA, new(big.Int).Sub(delegatorAMint, delegatorABurn))

		// Send 15 AUT as a reward and check distribution
		liquidReward := new(big.Int).Mul(big.NewInt(15), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		expectedReward := new(big.Int).Mul(big.NewInt(10), params.DecimalFactor)
		checkReward(r, liquidState, validator, expectedReward, expectedReward)

		expectedReward = new(big.Int).Mul(big.NewInt(5), params.DecimalFactor)
		checkReward(r, liquidState, delegatorA, expectedReward, expectedReward)
	})

	r.run("claiming rewards", func(r *runner) {
		// delegatorA bonds 10000 NEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)

		// Send 20 AUT as a reward (validator and delegatorA each earns 10). Withdraw and check balance.
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)
		expectedReward := new(big.Int).Mul(big.NewInt(10), params.DecimalFactor)
		withdrawAndCheck(r, liquidState, liquidLogic, delegatorA, expectedReward, expectedReward)

		// Send 40 AUT as a reward (validator and delegatorA each earns 20). Withdraw and check balance.
		liquidReward = new(big.Int).Mul(big.NewInt(40), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)
		expectedReward = new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		withdrawAndCheck(r, liquidState, liquidLogic, delegatorA, expectedReward, expectedReward)
	})

	r.run("accumulating rewards", func(r *runner) {
		// delegatorA bonds 10000 NEW (total 20000 delegated)
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)

		// Send 20 AUT as a reward (delegatorA earns 10)
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Other delegators bond 20000 NEW (total of 40000 NEW bonded)
		delegatorBMint := new(big.Int).Mul(big.NewInt(12000), params.DecimalFactor)
		delegatorCMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorB, delegatorBMint),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorC, delegatorCMint),
		)

		// Send 20 AUT as a reward (delegatorA earns 5)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Other delegators bond 10000 NEW (total of 50000 NEW bonded)
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity,
				"mint", validator, new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor),
			),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorC, delegatorCMint),
		)

		// Send 50 AUT as a reward (delegatorA earns 10)
		liquidReward = new(big.Int).Mul(big.NewInt(50), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Check delegatorA's total fees were 10 + 5 + 10 = 25
		expectedReward := new(big.Int).Mul(big.NewInt(25), params.DecimalFactor)
		checkReward(r, liquidState, delegatorA, expectedReward, expectedReward)
	})

	r.run("commission", func(r *runner) {
		// use 50% commission for simplcity
		liquidState := deployLiquid(r, validator, treasury, 50)
		validatorMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity, "mint", validator, validatorMint,
			),
		)
		treasuryBalance := r.getBalanceOf(treasury)

		// delegatorA bonds 10000 NEW (total 20000 delegated)
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)

		// Send 40 AUT as a reward (treasury earns 20, delegatorA earns 10)
		liquidReward := new(big.Int).Mul(big.NewInt(40), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Other delegators bond 20000 NEW (total of 40000 NEW bonded)
		delegatorBMint := new(big.Int).Mul(big.NewInt(12000), params.DecimalFactor)
		delegatorCMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorB, delegatorBMint),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorC, delegatorCMint),
		)

		// Send 40 AUT as a reward (treasury earns 20 delegatorA earns 5)
		redistributeLiquidReward(r, liquidState, liquidReward)
		// Other delegators bond 10000 NEW (total of 50000 NEW bonded)
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity,
				"mint", validator, new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor),
			),
		)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorC, delegatorCMint),
		)

		// Send 100 AUT as a reward (treasury earns 50, delegatorA earns 10)
		liquidReward = new(big.Int).Mul(big.NewInt(100), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Check treasury balance increased by: 20 + 20 + 50 = 90
		require.Equal(
			r.t,
			new(big.Int).Add(treasuryBalance, new(big.Int).Mul(big.NewInt(90), params.DecimalFactor)),
			r.getBalanceOf(treasury),
		)

		// Check delegatorA's total fees: 10 + 5 + 10 = 25
		expectedReward := new(big.Int).Mul(big.NewInt(25), params.DecimalFactor)
		checkReward(r, liquidState, delegatorA, expectedReward, expectedReward)
	})

	r.run("allowances", func(r *runner) {
		// delegatorA bonds 10000 NEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromAutonity, "mint", delegatorA, delegatorAMint),
		)

		// delegatorC should not be able to transfer on A's behalf
		allowance, _, err := liquidState.Allowance(nil, delegatorA, delegatorC)
		require.NoError(r.t, err)
		require.True(r.t, allowance.Cmp(common.Big0) == 0)
		transfer := new(big.Int).Mul(big.NewInt(1000), params.DecimalFactor)
		_, _, err = liquidState.CallMethod(
			liquidLogic.contract, fromSender(delegatorC, nil), "transferFrom", delegatorA, delegatorB, transfer,
		)
		require.Error(r.t, err)
		require.Equal(r.t, "execution reverted: ERC20: transfer amount exceeds allowance", err.Error())

		// A grants C permission to spend 5000.
		approval := new(big.Int).Mul(big.NewInt(5000), params.DecimalFactor)
		r.CallNoError(
			liquidState.CallMethod(liquidLogic.contract, fromSender(delegatorA, nil), "approve", delegatorC, approval),
		)
		allowance, _, err = liquidState.Allowance(nil, delegatorA, delegatorC)
		require.NoError(r.t, err)
		require.Equal(r.t, approval, allowance)

		// C sends 1000 of A's LNEW to B
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromSender(delegatorC, nil), "transferFrom", delegatorA, delegatorB, transfer,
			),
		)

		// Check balances and allowances
		checkLiquidBalance(r, liquidState, delegatorA, new(big.Int).Sub(delegatorAMint, transfer))
		checkLiquidBalance(r, liquidState, delegatorB, transfer)
		checkLiquidBalance(r, liquidState, delegatorC, common.Big0)

		allowance, _, err = liquidState.Allowance(nil, delegatorA, delegatorC)
		require.NoError(r.t, err)
		require.Equal(r.t, new(big.Int).Sub(approval, transfer), allowance)

		// Sending 4001 should fail.
		transfer = new(big.Int).Mul(big.NewInt(4001), params.DecimalFactor)
		_, _, err = liquidState.CallMethod(
			liquidLogic.contract, fromSender(delegatorC, nil), "transferFrom", delegatorA, delegatorB, transfer,
		)
		require.Error(r.t, err)
		require.Equal(r.t, "execution reverted: ERC20: transfer amount exceeds allowance", err.Error())
	})

	r.run("locking", func(r *runner) {
		balance := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		balanceToLock := new(big.Int).Mul(big.NewInt(1000), params.DecimalFactor)
		// increment := new(big.Int).Mul(big.NewInt(100), params.DecimalFactor)

		// mint
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity, "mint", delegatorA, balance,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, balance)
		checkLockedLiquidBalance(r, liquidState, delegatorA, common.Big0)

		// lock more than balance
		_, _, err = liquidState.CallMethod(
			liquidLogic.contract, fromAutonity, "lock", delegatorA, new(big.Int).Add(balance, common.Big1),
		)
		require.Error(r.t, err)
		require.Equal(r.t, "execution reverted: can't lock more funds than available", err.Error())

		// lock less than balance
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity, "lock", delegatorA, balanceToLock,
			),
		)
		checkLockedLiquidBalance(r, liquidState, delegatorA, balanceToLock)
		checkLiquidBalance(r, liquidState, delegatorA, balance)

		maxTransferable := new(big.Int).Sub(balance, balanceToLock)
		// transfer more than unlocked
		_, _, err = liquidState.CallMethod(
			liquidLogic.contract, fromSender(delegatorA, nil), "transfer", delegatorB, new(big.Int).Add(maxTransferable, common.Big1),
		)
		require.Error(r.t, err)
		require.Equal(r.t, "execution reverted: insufficient unlocked funds", err.Error())

		// burn more than unlocked
		_, _, err = liquidState.CallMethod(
			liquidLogic.contract, fromAutonity, "burn", delegatorA, new(big.Int).Add(maxTransferable, common.Big1),
		)
		require.Error(r.t, err)
		require.Equal(r.t, "execution reverted: insufficient unlocked funds", err.Error())

		// cannot unlock more than locked
		_, _, err = liquidState.CallMethod(
			liquidLogic.contract, fromAutonity, "unlock", delegatorA, new(big.Int).Add(maxTransferable, common.Big1),
		)
		require.Error(r.t, err)
		require.Equal(r.t, "execution reverted: can't unlock more funds than locked", err.Error())

		// unlock
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity, "unlock", delegatorA, balanceToLock,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, balance)
		checkLockedLiquidBalance(r, liquidState, delegatorA, common.Big0)

		// transfer and burn whole amount
		transferAmount := new(big.Int).Add(maxTransferable, common.Big1)
		burnAmount := new(big.Int).Sub(balance, transferAmount)
		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromSender(delegatorA, nil), "transfer", delegatorB, transferAmount,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, burnAmount)

		r.CallNoError(
			liquidState.CallMethod(
				liquidLogic.contract, fromAutonity, "burn", delegatorA, burnAmount,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, common.Big0)
	})
}

func checkLiquidBalance(r *runner, liquidState *LiquidState, user common.Address, expecedBalance *big.Int) {
	balance, _, err := liquidState.BalanceOf(nil, user)
	require.NoError(r.t, err)
	require.True(r.t, balance.Cmp(expecedBalance) == 0)
}

func checkLockedLiquidBalance(r *runner, liquidState *LiquidState, user common.Address, expecedLockedBalance *big.Int) {
	lockedBalance, _, err := liquidState.LockedBalanceOf(nil, user)
	require.NoError(r.t, err)
	require.True(r.t, lockedBalance.Cmp(expecedLockedBalance) == 0)
}

func checkReward(r *runner, liquidState *LiquidState, user common.Address, atnReward, ntnReward *big.Int) {
	unclaimedRewards, _, err := liquidState.UnclaimedRewards(nil, user)
	require.NoError(r.t, err)
	fmt.Printf("unclaimedRewards.UnclaimedATN %v\n", unclaimedRewards.UnclaimedATN)
	fmt.Printf("atnReward %v\n", atnReward)
	fmt.Printf("unclaimedRewards.UnclaimedNTN %v\n", unclaimedRewards.UnclaimedNTN)
	fmt.Printf("ntnReward %v\n", ntnReward)

	require.True(r.t, unclaimedRewards.UnclaimedATN.Cmp(atnReward) == 0)
	require.True(r.t, unclaimedRewards.UnclaimedNTN.Cmp(ntnReward) == 0)
}

func withdrawAndCheck(
	r *runner, liquidState *LiquidState, liquidLogic *LiquidLogic, user common.Address, atnReward, ntnReward *big.Int,
) {
	ntnBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	atnBalance := r.getBalanceOf(user)

	checkReward(r, liquidState, user, atnReward, ntnReward)

	r.CallNoError(
		liquidState.CallMethod(liquidLogic.contract, fromSender(user, nil), "claimRewards"),
	)

	checkReward(r, liquidState, user, common.Big0, common.Big0)

	ntnNewBalance, _, err := r.autonity.BalanceOf(nil, user)
	require.NoError(r.t, err)
	atnNewBalance := r.getBalanceOf(user)

	require.Equal(r.t, new(big.Int).Add(ntnBalance, ntnReward), ntnNewBalance)
	require.Equal(r.t, new(big.Int).Add(atnBalance, atnReward), atnNewBalance)
}

func redistributeLiquidReward(r *runner, liquidState *LiquidState, reward *big.Int) {
	r.NoError(
		r.autonity.Mint(operator, liquidState.address, reward),
	)
	r.CallNoError(
		liquidState.CallMethod(
			r.LiquidLogicContractObject().contract,
			fromSender(r.autonity.address, reward),
			"redistribute", reward,
		),
	)
}

func deployLiquid(
	r *runner, validator, treasury common.Address, commissionRatePercent ...int64,
) *LiquidState {

	liquidLogic, _, err := r.autonity.LiquidLogicContract(nil)
	require.NoError(r.t, err)

	var commissionRate int64 = 0
	if len(commissionRatePercent) > 0 {
		commissionRate = commissionRatePercent[0]
	}

	_, _, liquidState, err := r.deployLiquidState(
		fromAutonity, validator, treasury,
		big.NewInt(commissionRate*100), "27", liquidLogic,
	)
	require.NoError(r.t, err)

	return liquidState
}
