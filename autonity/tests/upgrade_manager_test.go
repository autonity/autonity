package tests

import (
	"fmt"
	"testing"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"

	"github.com/stretchr/testify/require"
)

var (
	user = common.HexToAddress("0x99")
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
		_, err := r.upgradeManager.SetOperator(&runOptions{origin: user}, user)
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		_, err = r.upgradeManager.SetOperator(operator, user)
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		_, err = r.upgradeManager.SetOperator(&runOptions{origin: r.autonity.address}, user)
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
	r.run("upgrade accountability contract", func(r *runner) {
		calldata := makeCalldata(Accountability3MetaData)
		r.t.Log("upgrade autonity: calldata size:", len(calldata)/1000, "kB")
		gas, err := r.upgradeManager.Upgrade(operator, r.accountability.address, string(calldata))
		fmt.Println(common.Bytes2Hex(calldata))
		require.NoError(r.t, err)
		r.t.Log("upgrade autonity: gas consumed:", gas)
	})
}
