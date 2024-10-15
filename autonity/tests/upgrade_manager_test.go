package tests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"

	"github.com/stretchr/testify/require"
)

func TestExample1(t *testing.T) {
	r := Setup(t, nil)
	validators, consumed, err := r.Autonity.GetValidators(nil)
	require.NoError(t, err)
	require.Equal(t, *params.TestAutonityContractConfig.Validators[0].NodeAddress, validators[0])
	require.LessOrEqual(t, consumed, uint64(2000))
}

func TestExample2(t *testing.T) {
	// Setup phase ....
	r := Setup(t, nil)
	_, err := r.Autonity.Mint(Operator, User, common.Big2)
	require.NoError(t, err)
	// End setup - state snapshot here
	r.Run("sub-test1", func(r *Runner) {
		balance, _, _ := r.Autonity.BalanceOf(nil, User)
		require.Equal(r.T, common.Big2, balance)
		_, _ = r.Autonity.Mint(Operator, User, common.Big1)
		balance, _, _ = r.Autonity.BalanceOf(nil, User)
		require.Equal(r.T, common.Big3, balance)
	})
	r.Run("sub-test2", func(r *Runner) {
		balance, _, _ := r.Autonity.BalanceOf(nil, User)
		require.Equal(r.T, common.Big2, balance)
	})
}

func TestSetOperatorAccount(t *testing.T) {
	r := Setup(t, nil)
	r.Run("setOperatorAccount is restricted to Autonity contract", func(r *Runner) {
		_, err := r.UpgradeManager.SetOperator(&runOptions{origin: User}, User)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(Operator, User)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(&runOptions{origin: r.Autonity.address}, User)
		require.NoError(r.T, err)
	})
}

func TestUpgrade(t *testing.T) {
	makeCalldata := func(newContract *bind.MetaData, args ...any) (calldata []byte) {
		// craft upgrade transaction's calldata
		// append the upgrade contract's deployment bytecode
		calldata = append(calldata, common.FromHex(newContract.Bin)...)
		// then finally append the arguments
		parsed, _ := newContract.GetAbi()
		packedArgs, _ := parsed.Pack("", args...)
		calldata = append(calldata, packedArgs...)
		return
	}
	r := Setup(t, nil)
	r.Run("restricted to the Operator", func(r *Runner) {
		_, err := r.UpgradeManager.Upgrade(&runOptions{origin: User}, r.Autonity.address, "0x1111")
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted) // maybe check revert reason
	})
	r.Run("upgrade target contract", func(r *Runner) {
		// deploy first dummy contract
		_, _, base, err := r.DeployTestBase(Operator, "v1")
		require.NoError(r.T, err, base)
		v1string, _, _ := base.Foo(nil)
		require.Equal(r.T, v1string, "v1")
		calldata := makeCalldata(TestUpgradedMetaData, "hello", "v2")
		// call the replace function
		gas, err := r.UpgradeManager.Upgrade(Operator, base.address, string(calldata))
		require.NoError(r.T, err)
		r.T.Log("gas consumed:", gas)
		// check if base has been updated
		v2string, _, _ := base.Foo(nil)
		require.Equal(r.T, v2string, "v2")
		// todo: attach TestUpgraded to this address and check if new functions are exposed
	})
	r.Run("upgrade autonity contract", func(r *Runner) {
		calldata := makeCalldata(AutonityUpgradeTestMetaData)
		r.T.Log("upgrade autonity: calldata size:", len(calldata)/1000, "kB")
		gas, err := r.UpgradeManager.Upgrade(Operator, r.Autonity.address, string(calldata))
		require.NoError(r.T, err)
		r.T.Log("upgrade autonity: gas consumed:", gas)
		cfg, _, err := r.Autonity.Config(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, cfg.ContractVersion.Uint64(), common.Big2.Uint64())
		// test the hot patched _transfer operation, see AutonityUpgradeTest.sol
		r.Autonity.Transfer(Operator, User, big.NewInt(50))
		balance, _, _ := r.Autonity.BalanceOf(nil, User)
		require.Equal(r.T, balance.Uint64(), big.NewInt(100).Uint64())
	})
}
