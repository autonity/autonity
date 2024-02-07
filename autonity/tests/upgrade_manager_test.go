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

var (
	operator = &runOptions{origin: defaultAutonityConfig.Protocol.OperatorAccount}
	user     = common.HexToAddress("0x99")
)

func TestExample1(t *testing.T) {
	r := setup(t, nil)
	validators, consumed, err := r.autonity.GetValidators(nil)
	require.NoError(t, err)
	require.Equal(t, *params.TestAutonityContractConfig.Validators[0].NodeAddress, validators[0])
	require.LessOrEqual(t, consumed, uint64(2000))
}

func TestExample2(t *testing.T) {
	// Setup phase ....
	r := setup(t, nil)
	_, err := r.autonity.Mint(operator, user, common.Big2)
	require.NoError(t, err)
	// End setup - state snapshot here
	r.run("sub-test1", func(r *runner) {
		balance, _, _ := r.autonity.BalanceOf(nil, user)
		require.Equal(r.t, common.Big2, balance)
		_, _ = r.autonity.Mint(operator, user, common.Big1)
		balance, _, _ = r.autonity.BalanceOf(nil, user)
		require.Equal(r.t, common.Big3, balance)
	})
	r.run("sub-test2", func(r *runner) {
		balance, _, _ := r.autonity.BalanceOf(nil, user)
		require.Equal(r.t, common.Big2, balance)
	})
}

func TestSetOperatorAccount(t *testing.T) {
	r := setup(t, nil)
	r.run("setOperatorAccount is restricted to Autonity contract", func(r *runner) {
		_, err := r.UpgradeManager.SetOperator(&runOptions{origin: user}, user)
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(operator, user)
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(&runOptions{origin: r.autonity.address}, user)
		require.NoError(r.t, err)
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
	r := setup(t, nil)
	r.run("restricted to the operator", func(r *runner) {
		_, err := r.UpgradeManager.Upgrade(&runOptions{origin: user}, r.autonity.address, "0x1111")
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted) // maybe check revert reason
	})
	r.run("upgrade target contract", func(r *runner) {
		// deploy first dummy contract
		_, _, base, err := r.deployTestBase(operator, "v1")
		require.NoError(r.t, err, base)
		v1string, _, _ := base.Foo(nil)
		require.Equal(r.t, v1string, "v1")
		calldata := makeCalldata(TestUpgradedMetaData, "hello", "v2")
		// call the replace function
		gas, err := r.UpgradeManager.Upgrade(operator, base.address, string(calldata))
		require.NoError(r.t, err)
		r.t.Log("gas consumed:", gas)
		// check if base has been updated
		v2string, _, _ := base.Foo(nil)
		require.Equal(r.t, v2string, "v2")
		// todo: attach TestUpgraded to this address and check if new functions are exposed
	})
	r.run("upgrade autonity contract", func(r *runner) {
		calldata := makeCalldata(AutonityUpgradeTestMetaData)
		r.t.Log("upgrade autonity: calldata size:", len(calldata)/1000, "kB")
		gas, err := r.UpgradeManager.Upgrade(operator, r.autonity.address, string(calldata))
		require.NoError(r.t, err)
		r.t.Log("upgrade autonity: gas consumed:", gas)
		cfg, _, err := r.autonity.Config(nil)
		require.NoError(r.t, err)
		require.Equal(r.t, cfg.ContractVersion.Uint64(), common.Big2.Uint64())
		// test the hot patched _transfer operation, see AutonityUpgradeTest.sol
		r.autonity.Transfer(operator, user, big.NewInt(50))
		balance, _, _ := r.autonity.BalanceOf(nil, user)
		require.Equal(r.t, balance.Uint64(), big.NewInt(100).Uint64())
	})
}
