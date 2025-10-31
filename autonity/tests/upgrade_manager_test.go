package tests

import (
	"bytes"
	"math"
	"math/big"
	"testing"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/stretchr/testify/require"
)

// meta-test: tests that the testing framework works correctly
// doesn't strictly belong in this file but whatever
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
		_, err := r.UpgradeManager.SetOperator(&runOptions{origin: User}, User)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(r.Operator, User)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.SetOperator(&runOptions{origin: r.Autonity.address}, User)
		require.NoError(r.T, err)
	})
}

func TestUpgrade(t *testing.T) {
	upgradeBytecode := func(abi *abi.ABI, bytecode []byte, args ...any) ([]byte, error) {
		// if no args upgradeBytecode == bytecode
		if len(args) == 0 {
			return bytecode, nil
		}
		// otherwise append the packed constructor args
		packedArgs, err := abi.Pack("", args...)
		if err != nil {
			return nil, err
		}
		return append(bytecode, packedArgs...), nil
	}
	customOperator := common.Address{0xca, 0xfe}
	r := Setup(t, func(genesis *params.AutonityContractGenesis) *params.AutonityContractGenesis {
		genesis.Operator = customOperator
		return genesis
	})
	// sanity checks on operator
	originalOperator := r.Operator.origin
	fetchedOperator, _, err := r.UpgradeManager.GetOperator(nil)
	require.NoError(t, err)
	require.Equal(t, customOperator, originalOperator)
	require.Equal(t, customOperator, fetchedOperator)
	r.Run("upgrade functionality is restricted to the Operator", func(r *Runner) {
		_, err := r.UpgradeManager.Upgrade(&runOptions{origin: User},
			r.Acu.address,
			string(generated.ACUTestUpgradeBytecode),
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.Upgrade0(&runOptions{origin: User},
			r.Acu.address,
			string(generated.ACUTestUpgradeBytecode),
			"1.4.5",
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.UpgradeMultiple(&runOptions{origin: User},
			[]common.Address{r.Acu.address},
			[]string{string(generated.ACUTestUpgradeBytecode)},
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.UpgradeMultiple0(&runOptions{origin: User},
			[]common.Address{r.Acu.address},
			[]string{string(generated.ACUTestUpgradeBytecode)},
			[]string{"1.4.5"},
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		// if calling with operator account, should work fine
		_, err = r.UpgradeManager.Upgrade(r.Operator,
			r.Acu.address,
			string(generated.ACUTestUpgradeBytecode),
		)
		require.NoError(t, err)
		_, err = r.UpgradeManager.Upgrade0(r.Operator,
			r.Acu.address,
			string(generated.ACUTestUpgradeBytecode),
			"1.4.5",
		)
		require.NoError(t, err)
		_, err = r.UpgradeManager.UpgradeMultiple(r.Operator,
			[]common.Address{r.Acu.address},
			[]string{string(generated.ACUTestUpgradeBytecode)},
		)
		require.NoError(t, err)
		_, err = r.UpgradeManager.UpgradeMultiple0(r.Operator,
			[]common.Address{r.Acu.address},
			[]string{string(generated.ACUTestUpgradeBytecode)},
			[]string{"1.4.5"},
		)
		require.NoError(t, err)
	})
	r.Run("upgrade target contract", func(r *Runner) {
		// reset deployment params to fake a whitelisted protocol contract
		r.Evm.StateDB.SetNonce(common.Address{}, 0)
		r.Evm.StateDB.SetCode(params.AutonityContractAddress, []byte{})
		r.Evm.StateDB.SetNonce(params.AutonityContractAddress, 0)

		// deploy first dummy contract
		_, _, base, err := r.DeployTestBase(&runOptions{origin: common.Address{}, value: new(big.Int)}, "v1")
		require.NoError(r.T, err, base)
		v1string, _, err := base.Foo(nil)
		require.NoError(t, err)
		require.Equal(r.T, v1string, "v1")

		// upgrade to v2
		abiV2, err := TestUpgradedMetaData.GetAbi()
		require.NoError(t, err)
		upgradePayload, err := upgradeBytecode(abiV2, common.FromHex(TestUpgradedMetaData.Bin), "hello", "v2")
		require.NoError(t, err)
		// call the replace function
		gas, err := r.UpgradeManager.Upgrade(r.Operator, base.address, string(upgradePayload))
		require.NoError(r.T, err)
		r.T.Log("gas consumed:", gas)
		// check if base has been updated
		v2string, _, err := base.Foo(nil)
		require.NoError(t, err)
		require.Equal(r.T, v2string, "v2")

		upgraded := &TestUpgraded{
			&contract{
				address: base.address,
				abi:     abiV2,
				r:       r,
			},
		}
		v2bar, _, err := upgraded.Bar(nil)
		require.NoError(t, err)
		require.Equal(r.T, v2bar, "hello")

		_, err = upgraded.FooBar(nil, "modified")
		require.NoError(t, err)

		modFoo, _, err := upgraded.Foo(nil)
		require.NoError(t, err)
		require.Equal(r.T, modFoo, "modified")
	})
	r.Run("upgrade autonity contract", func(r *Runner) {
		autonityUpgradeAbi, err := AutonityUpgradeTestMetaData.GetAbi()
		require.NoError(t, err)
		upgradePayload, err := upgradeBytecode(autonityUpgradeAbi, common.FromHex(AutonityUpgradeTestMetaData.Bin))
		require.NoError(t, err)

		r.T.Log("upgrade autonity: payload size:", len(upgradePayload)/1000, "kB")
		gas, err := r.UpgradeManager.Upgrade(r.Operator, r.Autonity.address, string(upgradePayload))
		require.NoError(r.T, err)
		r.T.Log("upgrade autonity: gas consumed:", gas)
		cfg, _, err := r.Autonity.GetConfig(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, cfg.ContractVersion.Uint64(), common.Big2.Uint64())
		// test the hot patched _transfer operation, see AutonityUpgradeTest.sol
		_, err = r.Autonity.Transfer(r.Operator, User, big.NewInt(50))
		require.NoError(r.T, err)
		balance, _, err := r.Autonity.BalanceOf(nil, User)
		require.NoError(r.T, err)
		require.Equal(r.T, balance.Uint64(), big.NewInt(100).Uint64())
	})

	r.Run("upgrade the upgrade manager itself", func(r *Runner) {
		upgradePayload, err := upgradeBytecode(&generated.UpgradeManagerTestUpgradeAbi, generated.UpgradeManagerTestUpgradeBytecode)
		require.NoError(t, err)

		_, err = r.UpgradeManager.Upgrade(r.Operator, r.UpgradeManager.address, string(upgradePayload))
		require.NoError(r.T, err)

		// get version should be bricked
		_, _, err = r.UpgradeManager.GetVersion(nil, generated.AutonityCodeHash)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)

		// re-upgrade to original version

		upgradePayload, err = upgradeBytecode(&generated.UpgradeManager1Abi, generated.UpgradeManager1Bytecode, []common.Hash{}, []bindings.UpgradeManager1version{})
		require.NoError(t, err)

		_, err = r.UpgradeManager.Upgrade(r.Operator, r.UpgradeManager.address, string(upgradePayload))
		require.NoError(r.T, err)

		// get version should work again
		ver, _, err := r.UpgradeManager.GetVersion(nil, generated.AutonityCodeHash)
		t.Log(ver)
		require.NoError(t, err)

		// autonity and operator should still be the original values
		autonityAddress, _, err := r.UpgradeManager.GetAutonity(nil)
		require.NoError(t, err)
		require.Equal(r.T, params.AutonityContractAddress, autonityAddress)
		operatorAddress, _, err := r.UpgradeManager.GetOperator(nil)
		require.NoError(t, err)
		require.Equal(r.T, originalOperator, operatorAddress)
	})
	r.Run("upgrade with version tag", func(r *Runner) {
		auctioneerCodeHash := r.Evm.StateDB.GetCodeHash(params.AuctioneerContractAddress)
		initialVersion, _, err := r.UpgradeManager.GetVersion(nil, auctioneerCodeHash)
		require.NoError(t, err)
		t.Logf("initial version for auctioneer: %v (hash: %s)", initialVersion, auctioneerCodeHash)

		expectedVersionString := "12.47.11"
		_, err = r.UpgradeManager.Upgrade0(r.Operator,
			params.AuctioneerContractAddress,
			string(generated.AuctioneerTestUpgradeBytecode),
			expectedVersionString,
		)
		require.NoError(t, err)

		// old hash should still be tagged with initial version
		oldVersion, _, err := r.UpgradeManager.GetVersion(nil, auctioneerCodeHash)
		t.Logf("old version %v (hash: %s)", oldVersion, auctioneerCodeHash)
		require.NoError(t, err)
		require.Equal(t, oldVersion.Number, initialVersion.Number)
		require.Equal(t, oldVersion.Block.String(), initialVersion.Block.String())

		// code hash should have changed in statedb and new version applied
		upgradedAuctioneerCodeHash := r.Evm.StateDB.GetCodeHash(params.AuctioneerContractAddress)
		require.NotEqual(t, auctioneerCodeHash, upgradedAuctioneerCodeHash)
		upgradedVersion, _, err := r.UpgradeManager.GetVersion(nil, upgradedAuctioneerCodeHash)
		require.NoError(t, err)
		t.Logf("upgraded version for auctioneer: %v (hash: %s)", upgradedVersion, upgradedAuctioneerCodeHash)
		require.Equal(t, expectedVersionString, upgradedVersion.Number)
	})
	r.Run("upgrade multiple fails if args are not in same number", func(r *Runner) {
		_, err := r.UpgradeManager.UpgradeMultiple(r.Operator,
			[]common.Address{
				params.ACUContractAddress,
				//params.SupplyControlContractAddress,
			},
			[]string{
				string(generated.ACUTestUpgradeBytecode),
				string(generated.SupplyControlTestUpgradeBytecode),
			},
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.UpgradeMultiple(r.Operator,
			[]common.Address{
				params.ACUContractAddress,
				params.SupplyControlContractAddress,
			},
			[]string{
				string(generated.ACUTestUpgradeBytecode),
				//string(generated.SupplyControlTestUpgradeBytecode),
			},
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
		_, err = r.UpgradeManager.UpgradeMultiple0(r.Operator,
			[]common.Address{
				params.ACUContractAddress,
				params.SupplyControlContractAddress,
			},
			[]string{
				string(generated.ACUTestUpgradeBytecode),
				string(generated.SupplyControlTestUpgradeBytecode),
			},
			[]string{
				"1.6.0",
				//"1.6.0",
			},
		)
		t.Log(err)
		require.ErrorIs(r.T, err, vm.ErrExecutionReverted)
	})
	r.Run("upgrade multiple works as expected", func(r *Runner) {
		acuCodeHash := r.Evm.StateDB.GetCodeHash(params.ACUContractAddress)
		supplyControlCodeHash := r.Evm.StateDB.GetCodeHash(params.SupplyControlContractAddress)
		_, err := r.UpgradeManager.UpgradeMultiple(r.Operator,
			[]common.Address{
				params.ACUContractAddress,
				params.SupplyControlContractAddress,
			},
			[]string{
				string(generated.ACUTestUpgradeBytecode),
				string(generated.SupplyControlTestUpgradeBytecode),
			},
		)
		require.NoError(t, err)

		// code hash for both should have been changed
		require.NotEqual(t, acuCodeHash, r.Evm.StateDB.GetCodeHash(params.ACUContractAddress))
		require.NotEqual(t, supplyControlCodeHash, r.Evm.StateDB.GetCodeHash(params.SupplyControlContractAddress))
	})
	r.Run("upgrade multiple with versions works as expected", func(r *Runner) {
		acuCodeHash := r.Evm.StateDB.GetCodeHash(params.ACUContractAddress)
		supplyControlCodeHash := r.Evm.StateDB.GetCodeHash(params.SupplyControlContractAddress)

		acuInitialVersion, _, err := r.UpgradeManager.GetVersion(nil, acuCodeHash)
		require.NoError(t, err)
		t.Logf("initial version for acu: %v (hash: %s)", acuInitialVersion, acuCodeHash)

		supplyControlInitialVersion, _, err := r.UpgradeManager.GetVersion(nil, supplyControlCodeHash)
		require.NoError(t, err)
		t.Logf("initial version for supply control: %v (hash: %s)", supplyControlInitialVersion, supplyControlCodeHash)

		expectedAcuVersion := "11.6.0"
		expectedSupplyControlVersion := "13.7.0"
		_, err = r.UpgradeManager.UpgradeMultiple0(r.Operator,
			[]common.Address{
				params.ACUContractAddress,
				params.SupplyControlContractAddress,
			},
			[]string{
				string(generated.ACUTestUpgradeBytecode),
				string(generated.SupplyControlTestUpgradeBytecode),
			},
			[]string{
				expectedAcuVersion,
				expectedSupplyControlVersion,
			},
		)
		require.NoError(t, err)

		// old hashes should still be tagged with initial version
		oldAcuVersion, _, err := r.UpgradeManager.GetVersion(nil, acuCodeHash)
		require.NoError(t, err)
		t.Logf("old acu version %v (hash: %s)", oldAcuVersion, acuCodeHash)
		require.Equal(t, oldAcuVersion.Number, acuInitialVersion.Number)
		require.Equal(t, oldAcuVersion.Block.String(), acuInitialVersion.Block.String())

		oldSupplyControlVersion, _, err := r.UpgradeManager.GetVersion(nil, supplyControlCodeHash)
		require.NoError(t, err)
		t.Logf("old supply control version %v (hash: %s)", oldSupplyControlVersion, supplyControlCodeHash)
		require.Equal(t, oldSupplyControlVersion.Number, supplyControlInitialVersion.Number)
		require.Equal(t, oldSupplyControlVersion.Block.String(), supplyControlInitialVersion.Block.String())

		// code hash should have changed in statedb and new version applied
		upgradedAcuCodeHash := r.Evm.StateDB.GetCodeHash(params.ACUContractAddress)
		require.NotEqual(t, acuCodeHash, upgradedAcuCodeHash)
		upgradedAcuVersion, _, err := r.UpgradeManager.GetVersion(nil, upgradedAcuCodeHash)
		require.NoError(t, err)
		t.Logf("upgraded version for acu: %v (hash: %s)", upgradedAcuVersion, upgradedAcuCodeHash)
		require.Equal(t, expectedAcuVersion, upgradedAcuVersion.Number)

		upgradedSupplyControlCodeHash := r.Evm.StateDB.GetCodeHash(params.SupplyControlContractAddress)
		require.NotEqual(t, supplyControlCodeHash, upgradedSupplyControlCodeHash)
		upgradedSupplyControlVersion, _, err := r.UpgradeManager.GetVersion(nil, upgradedSupplyControlCodeHash)
		require.NoError(t, err)
		t.Logf("upgraded version for supplyControl: %v (hash: %s)", upgradedSupplyControlVersion, upgradedSupplyControlCodeHash)
		require.Equal(t, expectedSupplyControlVersion, upgradedSupplyControlVersion.Number)
	})
}

func TestSetAndGetVersion(t *testing.T) {
	r := Setup(t, nil)

	// restricted to operator
	_, err := r.UpgradeManager.SetVersion(&runOptions{origin: User}, common.Hash{}, "", new(big.Int))
	t.Log(err)
	require.ErrorIs(t, err, vm.ErrExecutionReverted)

	// getter is open
	version, _, err := r.UpgradeManager.GetVersion(nil, generated.AutonityCodeHash)
	require.NoError(t, err)
	t.Logf("version: %v", version)
	require.Equal(t, "1.0.0", version.Number)
	require.Equal(t, "0", version.Block.String())

	// empty hash doesn't have a version
	version, _, err = r.UpgradeManager.GetVersion(nil, common.Hash{})
	require.NoError(t, err)
	t.Logf("version: %v", version)
	require.Equal(t, "", version.Number)
	require.Equal(t, "0", version.Block.String())

	// set a version for empty code hash
	_, err = r.UpgradeManager.SetVersion(r.Operator, common.Hash{}, "145.5.6", new(big.Int).SetUint64(123))
	require.NoError(t, err)
	version, _, err = r.UpgradeManager.GetVersion(nil, common.Hash{})
	require.NoError(t, err)
	t.Logf("version: %v", version)
	require.Equal(t, "145.5.6", version.Number)
	require.Equal(t, "123", version.Block.String())

	// overwrite autonity contract hash
	_, err = r.UpgradeManager.SetVersion(r.Operator, generated.AutonityCodeHash, "a_weird_version_number", new(big.Int).SetUint64(456))
	require.NoError(t, err)
	version, _, err = r.UpgradeManager.GetVersion(nil, generated.AutonityCodeHash)
	require.NoError(t, err)
	t.Logf("version: %v", version)
	require.Equal(t, "a_weird_version_number", version.Number)
	require.Equal(t, "456", version.Block.String())

}

// tests that manually patched oracle deployment bytecode returns
// the correct runtime bytecode ( == to the one deployed on mainnet)
func TestOracleUpgradePatchedBytecode(t *testing.T) {
	r := Setup(t, nil)

	evmContract := vm.NewContract(vm.AccountRef(params.DeployerAddress), vm.AccountRef(params.OracleContractAddress), common.Big0, math.MaxUint64)
	evmContract.Code = generated.Oracle0Bytecode
	evmContract.CodeAddr = &params.OracleContractAddress

	// run deployment bytecode to get runtime bytecode
	ret, err := r.Evm.Interpreter().Run(evmContract, nil, false)
	require.NoError(t, err)
	// should be equal to the patched one
	require.True(t, bytes.Equal(ret, generated.Oracle0RuntimeBytecode))
}
