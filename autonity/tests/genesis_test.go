package tests

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/stretchr/testify/require"
)

// mirrors operations done in `toBlock` function
// when initializing chain
func prepare(g *core.Genesis) {
	// setDefaultHardforks()
	if g.Config.ByzantiumBlock == nil {
		g.Config.ByzantiumBlock = new(big.Int)
	}
	if g.Config.HomesteadBlock == nil {
		g.Config.HomesteadBlock = new(big.Int)
	}
	if g.Config.ConstantinopleBlock == nil {
		g.Config.ConstantinopleBlock = new(big.Int)
	}
	if g.Config.PetersburgBlock == nil {
		g.Config.PetersburgBlock = new(big.Int)
	}
	if g.Config.IstanbulBlock == nil {
		g.Config.IstanbulBlock = new(big.Int)
	}
	if g.Config.MuirGlacierBlock == nil {
		g.Config.MuirGlacierBlock = new(big.Int)
	}
	if g.Config.BerlinBlock == nil {
		g.Config.BerlinBlock = new(big.Int)
	}
	if g.Config.LondonBlock == nil {
		g.Config.LondonBlock = new(big.Int)
	}
	if g.Config.ArrowGlacierBlock == nil {
		g.Config.ArrowGlacierBlock = new(big.Int)
	}
	if g.Config.EIP158Block == nil {
		g.Config.EIP158Block = new(big.Int)
	}
	if g.Config.EIP150Block == nil {
		g.Config.EIP150Block = new(big.Int)
	}
	if g.Config.EIP155Block == nil {
		g.Config.EIP155Block = new(big.Int)
	}
	g.Config.SetDefaults()
	err := g.Config.Prepare()
	if err != nil {
		panic(err)
	}

	if g.Difficulty == nil {
		g.Difficulty = params.GenesisDifficulty
	}
}

// cannot be placed in ../genesis_test.go due to import loops
func TestGenesisSequence(t *testing.T) {
	newEVM := func(timestamp uint64) *vm.EVM {
		stateDB, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		require.NoError(t, err)

		vmBlockContext := vm.BlockContext{
			Transfer: func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
				db.SubBalance(sender, amount)
				db.AddBalance(recipient, amount)
			},
			CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool {
				return db.GetBalance(addr).Cmp(amount) >= 0
			},
			BlockNumber: common.Big0,
			Time:        new(big.Int).SetUint64(timestamp),
		}
		txContext := vm.TxContext{
			Origin:   common.Address{},
			GasPrice: common.Big0,
		}

		return vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})
	}

	// compatibility tests with bakerloo
	g := core.DefaultBakerlooGenesisBlock()
	prepare(g)
	evm := newEVM(g.Timestamp)
	err := autonity.ExecuteGenesisSequence(g.Config, g.Alloc.ToGenesisBonds(), evm)
	require.NoError(t, err)

	// no upgrade should have been deployed
	require.True(t, bytes.Equal(generated.OracleRuntimeBytecode, evm.StateDB.GetCode(params.OracleContractAddress)))
	require.True(t, bytes.Equal(generated.UpgradeManagerRuntimeBytecode, evm.StateDB.GetCode(params.UpgradeManagerContractAddress)))

	// compatibility tests with mainnet
	g = core.DefaultMainnetGenesisBlock()
	prepare(g)
	err = autonity.ExecuteGenesisSequence(g.Config, g.Alloc.ToGenesisBonds(), newEVM(g.Timestamp))
	require.NoError(t, err)

	// no upgrade should have been deployed
	require.True(t, bytes.Equal(generated.OracleRuntimeBytecode, evm.StateDB.GetCode(params.OracleContractAddress)))
	require.True(t, bytes.Equal(generated.UpgradeManagerRuntimeBytecode, evm.StateDB.GetCode(params.UpgradeManagerContractAddress)))
}
