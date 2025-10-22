package tests

import (
	"bytes"
	"math"
	"math/big"
	"testing"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	generated0 "github.com/autonity/autonity/params/upgrades/generated/0"
	generated1 "github.com/autonity/autonity/params/upgrades/generated/1"
	"github.com/stretchr/testify/require"
)

func TestExample1(t *testing.T) {
	r := Setup(t, nil)
	validators, consumed, err := r.Autonity.GetValidators(nil)
	require.NoError(t, err)
	require.Equal(t, *params.TestAutonityContractConfig.Validators[0].NodeAddress, validators[0])
	require.LessOrEqual(t, consumed, uint64(2200))
}

func TestExample2(t *testing.T) {
	// Setup phase ....
	r := Setup(t, nil)
	_, err := r.Autonity.Mint(r.Operator, User, common.Big2)
	require.NoError(t, err)
	// End setup - state snapshot here
	r.Run("sub-test1", func(r *Runner) {
		balance, _, _ := r.Autonity.BalanceOf(nil, User)
		require.Equal(r.T, common.Big2, balance)
		_, _ = r.Autonity.Mint(r.Operator, User, common.Big1)
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
		_, err := r.UpgradeManager.SetOperator(&RunOptions{origin: User}, User)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(r.Operator, User)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(&RunOptions{origin: r.Autonity.address}, User)
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
		_, err := r.UpgradeManager.Upgrade(&RunOptions{origin: User}, r.Autonity.address, "0x1111")
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted) // maybe check revert reason
	})
	r.Run("upgrade target contract", func(r *Runner) {
		// reset deployment params to fake a whitelisted protocol contract
		r.Evm.StateDB.SetNonce(common.Address{}, 0)
		r.Evm.StateDB.SetCode(params.AutonityContractAddress, []byte{})
		r.Evm.StateDB.SetNonce(params.AutonityContractAddress, 0)

		// deploy first dummy contract
		_, _, base, err := r.DeployTestBase(&RunOptions{origin: common.Address{}, value: new(big.Int)}, "v1")

		require.NoError(r.T, err, base)
		v1string, _, _ := base.Foo(nil)
		require.Equal(r.T, v1string, "v1")
		calldata := makeCalldata(TestUpgradedMetaData, "hello", "v2")
		// call the replace function
		gas, err := r.UpgradeManager.Upgrade(r.Operator, base.address, string(calldata))
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
		gas, err := r.UpgradeManager.Upgrade(r.Operator, r.Autonity.address, string(calldata))
		require.NoError(r.T, err)
		r.T.Log("upgrade autonity: gas consumed:", gas)
		cfg, _, err := r.Autonity.GetConfig(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, cfg.ContractVersion.Uint64(), common.Big2.Uint64())
		// test the hot patched _transfer operation, see AutonityUpgradeTest.sol
		r.Autonity.Transfer(r.Operator, User, big.NewInt(50))
		balance, _, _ := r.Autonity.BalanceOf(nil, User)
		require.Equal(r.T, balance.Uint64(), big.NewInt(100).Uint64())
	})

	r.Run("upgrade the upgrade manager itself", func(r *Runner) {
		_, err := r.UpgradeManager.Upgrade(r.Operator, r.UpgradeManager.address, string(generated1.UpgradeManagerBytecode))
		require.NoError(r.T, err)

		// TODO: set and fetch versions

	})
}

// tests that manually patched oracle deployment bytecode returns
// the correct runtime bytecode ( == to the one deployed on mainnet)
func TestOracleUpgradePatchedBytecode(t *testing.T) {
	r := Setup(t, nil)

	evmContract := vm.NewContract(vm.AccountRef(params.DeployerAddress), vm.AccountRef(params.OracleContractAddress), common.Big0, math.MaxUint64)
	evmContract.Code = generated0.OracleBytecode
	evmContract.CodeAddr = &params.OracleContractAddress

	// run deployment bytecode to get runtime bytecode
	ret, err := r.Evm.Interpreter().Run(evmContract, nil, false)
	require.NoError(t, err)
	// should be equal to the patched one
	require.True(t, bytes.Equal(ret, generated0.OracleRuntimeBytecode))
}

// test if upgrading all ASM contract at once would fit in a single tx
func TestASMAtomicUpdate(t *testing.T) {
	cumulativeBytecodeSize := 0
	cumulativeBytecodeSize += len(generated.ACUBytecode)
	cumulativeBytecodeSize += len(generated.SupplyControlBytecode)
	cumulativeBytecodeSize += len(generated.StabilizationBytecode)
	cumulativeBytecodeSize += len(generated.InflationControllerBytecode)
	cumulativeBytecodeSize += len(generated.AuctioneerBytecode)
	t.Logf("cumulativeBytecodeSize: %d bytes (~ %.2f kb)", cumulativeBytecodeSize, float64(cumulativeBytecodeSize)/float64(1024))
	require.True(t, cumulativeBytecodeSize < core.TxMaxSize)

}
