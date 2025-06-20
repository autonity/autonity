package autonitytests

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

func TestCirculatingSupply(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("total supply and circulating supply are same", setup, func(r *tests.Runner) {
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, totalSupply, circulatingSupply)
	})

	tests.RunWithSetup("minting increases both circulating and total supply", setup, func(r *tests.Runner) {
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)

		amount := big.NewInt(100)
		r.NoError(
			r.Autonity.Mint(r.Operator, tests.User, amount),
		)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(totalSupply, amount), newTotalSupply)
		require.Equal(r.T, new(big.Int).Add(circulatingSupply, amount), newCirculatingSupply)
	})

	tests.RunWithSetup("burning decreases both circulating and total supply", setup, func(r *tests.Runner) {
		mintAmount := big.NewInt(100)
		r.NoError(
			r.Autonity.Mint(r.Operator, tests.User, mintAmount),
		)
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		burnAmount := big.NewInt(60)
		r.NoError(
			r.Autonity.Burn(r.Operator, tests.User, burnAmount),
		)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Sub(totalSupply, burnAmount), newTotalSupply)
		require.Equal(r.T, new(big.Int).Sub(circulatingSupply, burnAmount), newCirculatingSupply)
	})

}

func TestDuplicateOracleAddress(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("no duplicate oracle address in autonity deployment", setup, func(r *tests.Runner) {
		config, _, err := r.Autonity.GetConfig(nil)
		require.NoError(r.T, err)
		validators := r.Committee.Validators
		for i := 0; i < len(validators); i++ {
			// there is a check for empty address in Prepare(), so there is no check for this on the contract side
			validators[i].OracleAddress = common.Address{}
		}
		config.ContractVersion = common.Big0
		_, _, _, err = r.DeployAutonity(
			nil,
			validators,
			config,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: oracle server exists", err.Error())
	})

	tests.RunWithSetup("no duplicate oracle address in validator registration", setup, func(r *tests.Runner) {
		// register validator with genesis validator oracle address
		validator, signature, nodeKey, blsKey, err := tests.RandomValidator()
		require.NoError(r.T, err)
		duplicateOracleKey, err := crypto.HexToECDSA(params.TestNodeKeys[0])
		require.NoError(r.T, err)
		pop, err := crypto.AutonityPOPProof(nodeKey, duplicateOracleKey, validator.Treasury.Hex(), blsKey)
		require.NoError(r.T, err)
		validator.OracleAddress = crypto.PubkeyToAddress(duplicateOracleKey.PublicKey)
		_, err = r.Autonity.RegisterValidator(
			tests.FromSender(validator.Treasury, nil),
			validator.Enode,
			validator.OracleAddress,
			validator.ConsensusKey,
			pop,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: oracle server exists", err.Error())

		validator.OracleAddress = validator.Treasury
		r.NoError(
			r.Autonity.RegisterValidator(
				tests.FromSender(validator.Treasury, nil),
				validator.Enode,
				validator.OracleAddress,
				validator.ConsensusKey,
				signature,
			),
		)

		// register validator with existing validator oracle address
		newValidator, signature, newNodeKey, blsKey, err := tests.RandomValidator()
		require.NoError(r.T, err)
		require.NotEqual(r.T, nodeKey, newNodeKey)
		pop, err = crypto.AutonityPOPProof(newNodeKey, nodeKey, newValidator.Treasury.Hex(), blsKey)
		require.NoError(r.T, err)
		newValidator.OracleAddress = validator.OracleAddress
		_, err = r.Autonity.RegisterValidator(
			tests.FromSender(newValidator.Treasury, nil),
			newValidator.Enode,
			newValidator.OracleAddress,
			newValidator.ConsensusKey,
			pop,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: oracle server exists", err.Error())

		newValidator.OracleAddress = newValidator.Treasury
		r.NoError(
			r.Autonity.RegisterValidator(
				tests.FromSender(newValidator.Treasury, nil),
				newValidator.Enode,
				newValidator.OracleAddress,
				newValidator.ConsensusKey,
				signature,
			),
		)
	})
}

func TestAutonityBalance(t *testing.T) {
	r := tests.Setup(t, nil)
	user := tests.User
	delegation := new(big.Int).Mul(
		big.NewInt(100_000_000),
		params.DecimalFactor,
	)

	isBalanceZero := func() {
		autonityBalance := r.GetNewtonBalanceOf(r.Autonity.Address())
		t.Logf("newton balance in autonity %v", autonityBalance)
		threshold := big.NewInt(int64(len(r.Committee.Validators))) // due to integer division in rewards distribution
		require.True(r.T, autonityBalance.Cmp(threshold) <= 0)
	}
	isBalanceZero()

	for _, v := range r.Committee.Validators {
		r.NoError(
			r.Autonity.Mint(
				r.Operator,
				user,
				delegation,
			),
		)

		r.NoError(
			r.Autonity.Bond(
				tests.FromSender(user, nil),
				v.NodeAddress,
				delegation,
			),
		)
	}
	isBalanceZero()

	r.WaitNextEpoch()
	isBalanceZero()
}
func TestConversionRatio(t *testing.T) {
	r := tests.Setup(t, nil)

	// initial conversion ratio should be 10 ^ 18 for everyone
	for _, member := range r.Committee.Validators {
		val, _, err := r.Autonity.GetValidator(nil, member.NodeAddress)
		require.NoError(t, err)
		require.Equal(t, uint64(1e18), val.ConversionRatio.Uint64())
	}

	var (
		staker1    = common.HexToAddress("0x1000000000000000000000000000000000000000")
		staker2    = common.HexToAddress("0x2000000000000000000000000000000000000000")
		staker3    = common.HexToAddress("0x3000000000000000000000000000000000000000")
		validator0 = r.Committee.Validators[0].NodeAddress
		validator1 = r.Committee.Validators[1].NodeAddress
	)

	// bond some delegated stake
	_, err := r.Autonity.Mint(r.Operator, staker1, params.Ntn10000)
	require.NoError(t, err)
	_, err = r.Autonity.Mint(r.Operator, staker2, params.Ntn10000)
	require.NoError(t, err)
	_, err = r.Autonity.Mint(r.Operator, staker3, params.Ntn40000)
	require.NoError(t, err)
	_, err = r.Autonity.Bond(tests.FromSender(staker1, nil), validator0, params.Ntn10000)
	require.NoError(t, err)
	_, err = r.Autonity.Bond(tests.FromSender(staker2, nil), validator1, params.Ntn10000)
	require.NoError(t, err)
	_, err = r.Autonity.Bond(tests.FromSender(staker3, nil), validator1, new(big.Int).Mul(common.Big2, params.Ntn10000))
	require.NoError(t, err)

	// let delegations apply
	r.WaitNextEpoch()

	// autobond at epoch end should increase the ratio
	r.WaitNextEpoch()

	for _, member := range r.Committee.Validators {
		// check only vals with delegated stake
		if bytes.Equal(member.NodeAddress.Bytes(), validator0.Bytes()) || bytes.Equal(member.NodeAddress.Bytes(), validator1.Bytes()) {
			val, _, err := r.Autonity.GetValidator(nil, member.NodeAddress)
			require.NoError(t, err)
			t.Logf("conversion ratio: %s", val.ConversionRatio.String())
			require.Greater(t, val.ConversionRatio.Uint64(), uint64(1e18))
		}
	}

	// slashing lowers conversion ratio
	_, err = r.Autonity.Slash(tests.FromSender(r.Accountability.Address(), nil), validator1, big.NewInt(9000))
	require.NoError(t, err)

	val, _, err := r.Autonity.GetValidator(nil, validator1)
	require.NoError(t, err)
	t.Logf("conversion ratio: %s", val.ConversionRatio.String())
	t.Logf("liquid supply: %s", val.LiquidSupply.String())
	require.Less(t, val.ConversionRatio.Uint64(), uint64(1e18))

	// unbond everything
	liquidContract := r.LiquidStateContract(validator1)
	staker2LiquidBalance, _, err := liquidContract.BalanceOf(nil, staker2)
	t.Log(staker2LiquidBalance)
	require.NoError(t, err)
	staker3LiquidBalance, _, err := liquidContract.BalanceOf(nil, staker3)
	t.Log(staker3LiquidBalance)
	require.NoError(t, err)

	_, err = r.Autonity.Unbond(tests.FromSender(staker2, nil), validator1, staker2LiquidBalance)
	require.NoError(t, err)
	_, err = r.Autonity.Unbond(tests.FromSender(staker3, nil), validator1, staker3LiquidBalance)
	require.NoError(t, err)

	// ratio will still increase a bit due to the epoch rewards
	r.WaitNextEpoch()

	// from here on ratio should stay constant

	val, _, err = r.Autonity.GetValidator(nil, validator1)
	require.NoError(t, err)
	t.Logf("conversion ratio: %s", val.ConversionRatio.String())
	ratioAfterUnbond := val.ConversionRatio.Uint64()

	r.WaitNextEpoch()
	r.WaitNextEpoch()
	r.WaitNextEpoch()
	r.WaitNextEpoch()

	val, _, err = r.Autonity.GetValidator(nil, validator1)
	require.NoError(t, err)
	t.Logf("conversion ratio: %s", val.ConversionRatio.String())
	require.Equal(t, ratioAfterUnbond, val.ConversionRatio.Uint64())

	// bonding restores the ratio to 1:1
	_, err = r.Autonity.Bond(tests.FromSender(staker3, nil), validator1, params.Ntn10000)
	require.NoError(t, err)

	r.WaitNextEpoch()

	val, _, err = r.Autonity.GetValidator(nil, validator1)
	require.NoError(t, err)
	t.Logf("conversion ratio: %s", val.ConversionRatio.String())
	require.Equal(t, uint64(1e18), val.ConversionRatio.Uint64())

}
