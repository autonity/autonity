package tests

import (
	"math/big"
	"testing"

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
	r := Setup(t, nil)
	// Mint Newton to some few accounts
	r.Autonity.Mint(Operator, staker1, params.Ntn10000)
	r.Autonity.Mint(Operator, staker2, params.Ntn10000)
	r.Autonity.Mint(Operator, staker3, params.Ntn10000)
	r.Autonity.Bond(&runOptions{origin: staker1}, r.Committee.Validators[0].NodeAddress, params.Ntn10000)
	r.Autonity.Bond(&runOptions{origin: staker2}, r.Committee.Validators[1].NodeAddress, params.Ntn10000)
	r.Autonity.Bond(&runOptions{origin: staker3}, r.Committee.Validators[1].NodeAddress, new(big.Int).Mul(common.Big2, params.Ntn10000))

	// create liquid staking contract per validator
	r.WaitNextEpoch()
	// .. test here claiming rewards, checking if NTN/ATN reward is coherent and accurate.
	// transactions fees can be simulated be sending atns directly to the autonity contract account.
	// todo: Think about in base.go to assign at each epoch the current list of validators / committee
	// in r.validators with the liquid stake contract bindings already prepared so that's easy to manipulate
	// or maybe just create some helpers for it.
}

func TestAccess(t *testing.T) {
	r := Setup(t, nil)
	validator := r.Committee.Validators[0].NodeAddress
	treasury := r.Committee.Validators[0].Treasury
	liquidState := deployLiquid(r, validator, treasury)

	r.Run("only autonity can mint", func(r *Runner) {
		_, err := liquidState.Mint(nil, validator, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
	})

	r.Run("only autonity can burn", func(r *Runner) {
		_, err := liquidState.Burn(nil, validator, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
	})

	r.Run("only autonity can lock", func(r *Runner) {
		_, err := liquidState.Lock(
			nil, validator, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
	})

	r.Run("only autonity can unlock", func(r *Runner) {
		_, err := liquidState.Unlock(
			nil, validator, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
	})

	r.Run("only autonity can redistribute", func(r *Runner) {
		_, err := liquidState.Redistribute(
			nil, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
	})

	r.Run("only autonity can setCommissionRate", func(r *Runner) {
		_, err := liquidState.SetCommissionRate(
			nil, common.Big1,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
	})

}

func TestLogicOperation(t *testing.T) {
	r := Setup(t, nil)

	validator := r.Committee.Validators[0].NodeAddress
	treasury := r.Committee.Validators[0].Treasury

	r.Run("liquid logic can be updated", func(r *Runner) {
		stateContract := deployLiquidTest(r, validator, treasury)
		liquidLogicFromAutonity, _, err := r.Autonity.LiquidLogicContract(nil)
		require.NoError(r.T, err)
		liquidLogicFromState, _, err := stateContract.LiquidLogicContract(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, liquidLogicFromAutonity, liquidLogicFromState)

		newLiquidLogic, _, _, err := r.DeployLiquidLogic(nil)
		require.NoError(r.T, err)
		require.NotEqual(r.T, liquidLogicFromAutonity, newLiquidLogic)

		r.NoError(
			r.Autonity.SetLiquidLogicContract(Operator, newLiquidLogic),
		)

		liquidLogicFromAutonity, _, err = r.Autonity.LiquidLogicContract(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, newLiquidLogic, liquidLogicFromAutonity)
		liquidLogicFromState, _, err = stateContract.LiquidLogicContract(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, liquidLogicFromAutonity, liquidLogicFromState)
	})

	liquidState := deployLiquid(r, validator, treasury)

	r.Run("updating liquid logic does not update state", func(r *Runner) {
		checkLiquidBalance(r, liquidState, validator, common.Big0)
		r.NoError(
			liquidState.Mint(FromAutonity, validator, common.Big1),
		)
		checkLiquidBalance(r, liquidState, validator, common.Big1)

		newLiquidLogic, _, _, err := r.DeployLiquidLogic(nil)
		require.NoError(r.T, err)

		r.NoError(
			r.Autonity.SetLiquidLogicContract(Operator, newLiquidLogic),
		)
		checkLiquidBalance(r, liquidState, validator, common.Big1)

		_, err = liquidState.Mint(nil, validator, common.Big1)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: Call restricted to the Autonity Contract", err.Error())
		checkLiquidBalance(r, liquidState, validator, common.Big1)

		r.NoError(
			liquidState.Mint(FromAutonity, validator, common.Big1),
		)
		checkLiquidBalance(r, liquidState, validator, common.Big2)
	})

	r.Run("liquid logic storage is separate than liquid state storage", func(r *Runner) {
		_, _, newLiquidLogic, err := r.DeployLiquidLogic(FromAutonity)
		require.NoError(r.T, err)

		r.NoError(
			r.Autonity.SetLiquidLogicContract(Operator, newLiquidLogic.address),
		)

		r.NoError(
			newLiquidLogic.Mint(FromAutonity, validator, common.Big1),
		)
		checkLiquidBalance(r, liquidState, validator, common.Big0)
	})

	r.Run("non-implemented method reverts", func(r *Runner) {
		_, _, err := liquidState.CallMethod(r.Autonity.contract, nil, "finalize")
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: fallback not implemented for LiquidLogic", err.Error())
	})
}

func TestFunctions(t *testing.T) {
	r := Setup(t, nil)

	validator := r.Committee.Validators[0].NodeAddress
	treasury := r.Committee.Validators[0].Treasury
	delegatorA := r.Committee.Validators[1].NodeAddress
	delegatorB := r.Committee.Validators[2].NodeAddress
	delegatorC := r.Committee.Validators[3].NodeAddress
	liquidState := deployLiquid(r, validator, treasury)

	validatorMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
	r.NoError(
		liquidState.Mint(
			FromAutonity, validator, validatorMint,
		),
	)

	supply, _, err := liquidState.TotalSupply(nil)
	require.NoError(r.T, err)
	require.Equal(r.T, validatorMint, supply)
	checkLiquidBalance(r, liquidState, validator, validatorMint)

	r.GiveMeSomeMoney(r.Autonity.address, new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor))

	r.Run("check name and symbol", func(r *Runner) {
		name, _, err := liquidState.Name(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, "LNTN-27", name)
		symbol, _, err := liquidState.Symbol(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, "LNTN-27", symbol)
	})

	r.Run("reward single validator", func(r *Runner) {
		// Initial state
		checkLiquidBalance(r, liquidState, delegatorA, common.Big0)
		checkLiquidBalance(r, liquidState, delegatorB, common.Big0)

		checkReward(r, liquidState, validator, common.Big0, common.Big0)
		checkReward(r, liquidState, delegatorA, common.Big0, common.Big0)
		checkReward(r, liquidState, delegatorB, common.Big0, common.Big0)

		// Send 10 ATN as a reward.  Perform a call first (not a tx)
		// in order to check the returned value.
		liquidReward := new(big.Int).Mul(big.NewInt(10), params.DecimalFactor)
		atnDistributed, ntnDistributed, _, err := liquidState.CallRedistribute(r, FromSender(r.Autonity.address, liquidReward), liquidReward)
		require.NoError(r.T, err)
		// out, _ := r.CallNoError(
		// 	liquidState.SimulateCall(
		// 		liquidState.contract,
		// 		FromSender(r.Autonity.address, liquidReward),
		// 		"redistribute", liquidReward,
		// 	),
		// )
		require.True(r.T, liquidReward.Cmp(atnDistributed) >= 0)
		require.True(r.T, liquidReward.Cmp(ntnDistributed) >= 0)
		precision := new(big.Int).Div(params.DecimalFactor, big.NewInt(10000))
		require.True(r.T, new(big.Int).Sub(liquidReward, precision).Cmp(ntnDistributed) <= 0)
		require.True(r.T, new(big.Int).Sub(liquidReward, precision).Cmp(ntnDistributed) <= 0)

		redistributeLiquidReward(r, liquidState, liquidReward)

		// Check distribution (only validator should hold this)
		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, validatorMint, supply)
		checkLiquidBalance(r, liquidState, validator, validatorMint)
		checkLiquidBalance(r, liquidState, delegatorA, common.Big0)
		checkLiquidBalance(r, liquidState, delegatorB, common.Big0)

		checkReward(r, liquidState, validator, liquidReward, liquidReward)
		checkReward(r, liquidState, delegatorA, common.Big0, common.Big0)
		checkReward(r, liquidState, delegatorB, common.Big0, common.Big0)
	})

	r.Run("reward multiple validators", func(r *Runner) {
		// delegatorA bonds 8000 NEW
		// delegatorB bonds 2000 NEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		delegatorBMint := new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorB, delegatorBMint),
		)

		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Mul(big.NewInt(20000), params.DecimalFactor), supply)

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

	r.Run("transfer LNEW", func(r *Runner) {
		// delegatorA bonds 8000 NEW
		// delegatorB bonds 2000 NEW
		// 20 AUT reward
		delegatorAMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		delegatorBMint := new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorB, delegatorBMint),
		)

		// Send 20 AUT as a reward and check distribution
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// delegatorA gives delegatorC 3000 LNEW
		transfer := new(big.Int).Mul(big.NewInt(3000), params.DecimalFactor)
		r.NoError(
			liquidState.Transfer(FromSender(delegatorA, nil), delegatorC, transfer),
		)

		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Mul(big.NewInt(20000), params.DecimalFactor), supply)

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

	r.Run("burn LNEW", func(r *Runner) {
		// delegatorA bonds 8000 NEW and burns 3000 LNEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		delegatorABurn := new(big.Int).Mul(big.NewInt(3000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)
		r.NoError(
			liquidState.Burn(FromAutonity, delegatorA, delegatorABurn),
		)

		supply, _, err := liquidState.TotalSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Mul(big.NewInt(15000), params.DecimalFactor), supply)

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

	r.Run("claiming rewards", func(r *Runner) {
		// delegatorA bonds 10000 NEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)

		// Send 20 AUT as a reward (validator and delegatorA each earns 10). Withdraw and check balance.
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)
		expectedReward := new(big.Int).Mul(big.NewInt(10), params.DecimalFactor)
		withdrawAndCheck(r, liquidState, delegatorA, expectedReward, expectedReward)

		// Send 40 AUT as a reward (validator and delegatorA each earns 20). Withdraw and check balance.
		liquidReward = new(big.Int).Mul(big.NewInt(40), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)
		expectedReward = new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		withdrawAndCheck(r, liquidState, delegatorA, expectedReward, expectedReward)
	})

	r.Run("accumulating rewards", func(r *Runner) {
		// delegatorA bonds 10000 NEW (total 20000 delegated)
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)

		// Send 20 AUT as a reward (delegatorA earns 10)
		liquidReward := new(big.Int).Mul(big.NewInt(20), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Other delegators bond 20000 NEW (total of 40000 NEW bonded)
		delegatorBMint := new(big.Int).Mul(big.NewInt(12000), params.DecimalFactor)
		delegatorCMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorB, delegatorBMint),
		)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorC, delegatorCMint),
		)

		// Send 20 AUT as a reward (delegatorA earns 5)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Other delegators bond 10000 NEW (total of 50000 NEW bonded)
		r.NoError(
			liquidState.Mint(
				FromAutonity, validator, new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor),
			),
		)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorC, delegatorCMint),
		)

		// Send 50 AUT as a reward (delegatorA earns 10)
		liquidReward = new(big.Int).Mul(big.NewInt(50), params.DecimalFactor)
		redistributeLiquidReward(r, liquidState, liquidReward)

		// Check delegatorA's total fees were 10 + 5 + 10 = 25
		expectedReward := new(big.Int).Mul(big.NewInt(25), params.DecimalFactor)
		checkReward(r, liquidState, delegatorA, expectedReward, expectedReward)
	})

	r.Run("commission", func(r *Runner) {
		// use 50% commission for simplcity
		newLiquidState := deployLiquid(r, validator, treasury, 50)
		validatorMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.NoError(
			newLiquidState.Mint(
				FromAutonity, validator, validatorMint,
			),
		)
		treasuryBalance := r.GetBalanceOf(treasury)

		// delegatorA bonds 10000 NEW (total 20000 delegated)
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.NoError(
			newLiquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)

		// Send 40 AUT as a reward (treasury earns 20, delegatorA earns 10)
		liquidReward := new(big.Int).Mul(big.NewInt(40), params.DecimalFactor)
		redistributeLiquidReward(r, newLiquidState, liquidReward)

		// Other delegators bond 20000 NEW (total of 40000 NEW bonded)
		delegatorBMint := new(big.Int).Mul(big.NewInt(12000), params.DecimalFactor)
		delegatorCMint := new(big.Int).Mul(big.NewInt(8000), params.DecimalFactor)
		r.NoError(
			newLiquidState.Mint(FromAutonity, delegatorB, delegatorBMint),
		)
		r.NoError(
			newLiquidState.Mint(FromAutonity, delegatorC, delegatorCMint),
		)

		// Send 40 AUT as a reward (treasury earns 20 delegatorA earns 5)
		redistributeLiquidReward(r, newLiquidState, liquidReward)
		// Other delegators bond 10000 NEW (total of 50000 NEW bonded)
		r.NoError(
			newLiquidState.Mint(
				FromAutonity, validator, new(big.Int).Mul(big.NewInt(2000), params.DecimalFactor),
			),
		)
		r.NoError(
			newLiquidState.Mint(FromAutonity, delegatorC, delegatorCMint),
		)

		// Send 100 AUT as a reward (treasury earns 50, delegatorA earns 10)
		liquidReward = new(big.Int).Mul(big.NewInt(100), params.DecimalFactor)
		redistributeLiquidReward(r, newLiquidState, liquidReward)

		// Check treasury balance increased by: 20 + 20 + 50 = 90
		require.Equal(
			r.T,
			new(big.Int).Add(treasuryBalance, new(big.Int).Mul(big.NewInt(90), params.DecimalFactor)),
			r.GetBalanceOf(treasury),
		)

		// Check delegatorA's total fees: 10 + 5 + 10 = 25
		expectedReward := new(big.Int).Mul(big.NewInt(25), params.DecimalFactor)
		checkReward(r, newLiquidState, delegatorA, expectedReward, expectedReward)
	})

	r.Run("allowances", func(r *Runner) {
		// delegatorA bonds 10000 NEW
		delegatorAMint := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		r.NoError(
			liquidState.Mint(FromAutonity, delegatorA, delegatorAMint),
		)

		// delegatorC should not be able to transfer on A's behalf
		allowance, _, err := liquidState.Allowance(nil, delegatorA, delegatorC)
		require.NoError(r.T, err)
		require.True(r.T, allowance.Cmp(common.Big0) == 0)
		transfer := new(big.Int).Mul(big.NewInt(1000), params.DecimalFactor)
		_, err = liquidState.TransferFrom(
			FromSender(delegatorC, nil), delegatorA, delegatorB, transfer,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: ERC20: transfer amount exceeds allowance", err.Error())

		// A grants C permission to spend 5000.
		approval := new(big.Int).Mul(big.NewInt(5000), params.DecimalFactor)
		r.NoError(
			liquidState.Approve(FromSender(delegatorA, nil), delegatorC, approval),
		)
		allowance, _, err = liquidState.Allowance(nil, delegatorA, delegatorC)
		require.NoError(r.T, err)
		require.Equal(r.T, approval, allowance)

		// C sends 1000 of A's LNEW to B
		r.NoError(
			liquidState.TransferFrom(
				FromSender(delegatorC, nil), delegatorA, delegatorB, transfer,
			),
		)

		// Check balances and allowances
		checkLiquidBalance(r, liquidState, delegatorA, new(big.Int).Sub(delegatorAMint, transfer))
		checkLiquidBalance(r, liquidState, delegatorB, transfer)
		checkLiquidBalance(r, liquidState, delegatorC, common.Big0)

		allowance, _, err = liquidState.Allowance(nil, delegatorA, delegatorC)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Sub(approval, transfer), allowance)

		// Sending 4001 should fail.
		transfer = new(big.Int).Mul(big.NewInt(4001), params.DecimalFactor)
		_, err = liquidState.TransferFrom(
			FromSender(delegatorC, nil), delegatorA, delegatorB, transfer,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: ERC20: transfer amount exceeds allowance", err.Error())
	})

	r.Run("locking", func(r *Runner) {
		balance := new(big.Int).Mul(big.NewInt(10000), params.DecimalFactor)
		balanceToLock := new(big.Int).Mul(big.NewInt(1000), params.DecimalFactor)
		// increment := new(big.Int).Mul(big.NewInt(100), params.DecimalFactor)

		// mint
		r.NoError(
			liquidState.Mint(
				FromAutonity, delegatorA, balance,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, balance)
		checkLockedLiquidBalance(r, liquidState, delegatorA, common.Big0)

		// lock more than balance
		_, err = liquidState.Lock(
			FromAutonity, delegatorA, new(big.Int).Add(balance, common.Big1),
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: can't lock more funds than available", err.Error())

		// lock less than balance
		r.NoError(
			liquidState.Lock(
				FromAutonity, delegatorA, balanceToLock,
			),
		)
		checkLockedLiquidBalance(r, liquidState, delegatorA, balanceToLock)
		checkLiquidBalance(r, liquidState, delegatorA, balance)

		maxTransferable := new(big.Int).Sub(balance, balanceToLock)
		// transfer more than unlocked
		_, err = liquidState.Transfer(
			FromSender(delegatorA, nil), delegatorB, new(big.Int).Add(maxTransferable, common.Big1),
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: insufficient unlocked funds", err.Error())

		// burn more than unlocked
		_, err = liquidState.Burn(
			FromAutonity, delegatorA, new(big.Int).Add(maxTransferable, common.Big1),
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: insufficient unlocked funds", err.Error())

		// cannot unlock more than locked
		_, err = liquidState.Unlock(
			FromAutonity, delegatorA, new(big.Int).Add(maxTransferable, common.Big1),
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: can't unlock more funds than locked", err.Error())

		// unlock
		r.NoError(
			liquidState.Unlock(
				FromAutonity, delegatorA, balanceToLock,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, balance)
		checkLockedLiquidBalance(r, liquidState, delegatorA, common.Big0)

		// transfer and burn whole amount
		transferAmount := new(big.Int).Add(maxTransferable, common.Big1)
		burnAmount := new(big.Int).Sub(balance, transferAmount)
		r.NoError(
			liquidState.Transfer(
				FromSender(delegatorA, nil), delegatorB, transferAmount,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, burnAmount)

		r.NoError(
			liquidState.Burn(
				FromAutonity, delegatorA, burnAmount,
			),
		)
		checkLiquidBalance(r, liquidState, delegatorA, common.Big0)
	})
}

func checkLiquidBalance(r *Runner, liquidState *ILiquid, user common.Address, expecedBalance *big.Int) {
	balance, _, err := liquidState.BalanceOf(nil, user)
	require.NoError(r.T, err)
	require.True(r.T, balance.Cmp(expecedBalance) == 0)
}

func checkLockedLiquidBalance(r *Runner, liquidState *ILiquid, user common.Address, expecedLockedBalance *big.Int) {
	lockedBalance, _, err := liquidState.LockedBalanceOf(nil, user)
	require.NoError(r.T, err)
	require.True(r.T, lockedBalance.Cmp(expecedLockedBalance) == 0)
}

func checkReward(r *Runner, liquidState *ILiquid, user common.Address, atnReward, ntnReward *big.Int) {
	abi, err := ILiquidMetaData.GetAbi()
	require.NoError(r.T, err)
	liquidLogicInterface := ILiquid{
		&contract{liquidState.address, abi, r},
	}
	unclaimedRewards, _, err := liquidLogicInterface.UnclaimedRewards(nil, user)
	require.NoError(r.T, err)

	require.True(r.T, unclaimedRewards.UnclaimedATN.Cmp(atnReward) == 0)
	require.True(r.T, unclaimedRewards.UnclaimedNTN.Cmp(ntnReward) == 0)
}

func withdrawAndCheck(
	r *Runner, liquidState *ILiquid, user common.Address, atnReward, ntnReward *big.Int,
) {
	ntnBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)
	atnBalance := r.GetBalanceOf(user)

	checkReward(r, liquidState, user, atnReward, ntnReward)

	r.NoError(
		liquidState.ClaimRewards(FromSender(user, nil)),
	)

	checkReward(r, liquidState, user, common.Big0, common.Big0)

	ntnNewBalance, _, err := r.Autonity.BalanceOf(nil, user)
	require.NoError(r.T, err)
	atnNewBalance := r.GetBalanceOf(user)

	require.Equal(r.T, new(big.Int).Add(ntnBalance, ntnReward), ntnNewBalance)
	require.Equal(r.T, new(big.Int).Add(atnBalance, atnReward), atnNewBalance)
}

func redistributeLiquidReward(r *Runner, liquidState *ILiquid, reward *big.Int) {
	r.NoError(
		r.Autonity.Mint(Operator, liquidState.address, reward),
	)
	r.NoError(
		liquidState.Redistribute(
			FromSender(r.Autonity.address, reward),
			reward,
		),
	)
}

func deployLiquid(
	r *Runner, validator, treasury common.Address, commissionRatePercent ...int64,
) *ILiquid {

	liquidLogic, _, err := r.Autonity.LiquidLogicContract(nil)
	require.NoError(r.T, err)

	var commissionRate int64
	if len(commissionRatePercent) > 0 {
		commissionRate = commissionRatePercent[0]
	}

	_, _, liquidState, err := r.DeployLiquidState(
		FromAutonity, validator, treasury,
		big.NewInt(commissionRate*100), "27", liquidLogic,
	)
	require.NoError(r.T, err)

	abi, err := ILiquidMetaData.GetAbi()
	require.NoError(r.T, err)
	return &ILiquid{
		&contract{liquidState.address, abi, r},
	}
}

func deployLiquidTest(
	r *Runner, validator, treasury common.Address, commissionRatePercent ...int64,
) *LiquidStateTest {

	liquidLogic, _, err := r.Autonity.LiquidLogicContract(nil)
	require.NoError(r.T, err)

	var commissionRate int64
	if len(commissionRatePercent) > 0 {
		commissionRate = commissionRatePercent[0]
	}

	_, _, liquidState, err := r.DeployLiquidStateTest(
		FromAutonity, validator, treasury,
		big.NewInt(commissionRate*100), "27", liquidLogic,
	)
	require.NoError(r.T, err)
	return liquidState
}
