package tests

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/params"
)

var (
	// todo: replicate truffle tests default config.
	defaultAutonityConfig = AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:     new(big.Int).SetUint64(params.TestAutonityContractConfig.TreasuryFee),
			MinBaseFee:      new(big.Int).SetUint64(params.TestAutonityContractConfig.MinBaseFee),
			DelegationRate:  new(big.Int).SetUint64(params.TestAutonityContractConfig.DelegationRate),
			UnbondingPeriod: new(big.Int).SetUint64(params.TestAutonityContractConfig.UnbondingPeriod),
			TreasuryAccount: params.TestAutonityContractConfig.Operator,
		},
		Contracts: AutonityContracts{
			AccountabilityContract: params.AccountabilityContractAddress,
			OracleContract:         params.OracleContractAddress,
			AcuContract:            params.ACUContractAddress,
			SupplyControlContract:  params.SupplyControlContractAddress,
			StabilizationContract:  params.StabilizationContractAddress,
			UpgradeManagerContract: params.UpgradeManagerContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount: params.TestAutonityContractConfig.Operator,
			EpochPeriod:     new(big.Int).SetUint64(params.TestAutonityContractConfig.EpochPeriod),
			BlockPeriod:     new(big.Int).SetUint64(params.TestAutonityContractConfig.BlockPeriod),
			CommitteeSize:   new(big.Int).SetUint64(params.TestAutonityContractConfig.MaxCommitteeSize),
		},
		ContractVersion: big.NewInt(1),
	}
)

type runOptions struct {
	origin common.Address
	value  *big.Int
}

type contract struct {
	address common.Address
	abi     *abi.ABI
	r       *runner
}

func (c *contract) call(opts *runOptions, method string, params ...any) ([]any, uint64, error) {
	input, err := c.abi.Pack(method, params...)
	require.NoError(c.r.t, err)
	out, consumed, err := c.r.call(opts, c.address, input)
	if err != nil {
		reason, _ := abi.UnpackRevert(out)
		return nil, 0, fmt.Errorf("%w: %s", err, reason)
	}
	res, err := c.abi.Unpack(method, out)
	require.NoError(c.r.t, err)
	return res, consumed, nil
}

type runner struct {
	t      *testing.T
	evm    *vm.EVM
	origin common.Address // session's sender, can be overridden via runOptions

	// protocol contracts
	// todo: see if genesis deployment flow can be abstracted somehow
	autonity       *Autonity
	accountability *Accountability
	oracle         *Oracle
	acu            *ACU
	supplyControl  *SupplyControl
	stabilization  *Stabilization
	UpgradeManager *UpgradeManager
}

func (r *runner) call(opts *runOptions, addr common.Address, input []byte) ([]byte, uint64, error) {
	r.evm.Origin = r.origin
	value := common.Big0
	if opts != nil {
		r.evm.Origin = opts.origin
		if opts.value != nil {
			value = opts.value
		}
	}
	gas := uint64(math.MaxUint64)
	ret, leftOver, err := r.evm.Call(vm.AccountRef(r.evm.Origin), addr, input, gas, value)
	return ret, gas - leftOver, err
}

func (r *runner) snapshot() int {
	return r.evm.StateDB.Snapshot()
}

func (r *runner) revertSnapshot(id int) {
	r.evm.StateDB.RevertToSnapshot(id)
}

// run is a convenience wrapper against t.run with automated state snapshot
func (r *runner) run(name string, f func(r *runner)) {
	r.t.Run(name, func(t2 *testing.T) {
		t := r.t
		r.t = t2
		// in the future avoid mutating for supporting parallel testing
		snap := r.snapshot()
		f(r)
		r.revertSnapshot(snap)
		r.t = t
	})
}

func (r *runner) getBalanceOf(account common.Address) *big.Int { //nolint
	return r.evm.StateDB.GetBalance(account)
}

func (r *runner) deployContract(opts *runOptions, abi *abi.ABI, bytecode []byte, params ...any) (common.Address, uint64, *contract, error) {
	args, err := abi.Pack("", params...)
	require.NoError(r.t, err)
	data := append(bytecode, args...)
	gas := uint64(math.MaxUint64)
	r.evm.Origin = r.origin
	value := common.Big0
	if opts != nil {
		r.evm.Origin = opts.origin
		if opts.value != nil {
			value = opts.value
		}
	}
	_, contractAddress, leftOverGas, err := r.evm.Create(vm.AccountRef(r.evm.Origin), data, gas, value)
	return contractAddress, gas - leftOverGas, &contract{contractAddress, abi, r}, err
}

func (r *runner) waitNBlocks(n int) { //nolint
	start := r.evm.Context.BlockNumber
	for i := 0; i < n; i++ {
		// Finalize is not the only block closing operation - fee redistribution is missing and prob
		// other stuff. Left as todo.
		_, err := r.autonity.Finalize(&runOptions{origin: common.Address{}})
		// consider monitoring gas cost here and fail if it's too much
		require.NoError(r.t, err, "finalize function error in waitNblocks", i)
		r.evm.Context.BlockNumber = new(big.Int).Add(big.NewInt(int64(i+1)), start)
	}
}

func (r *runner) waitNextEpoch() { //nolint
	epochPeriod, _, err := r.autonity.GetEpochPeriod(nil)
	require.NoError(r.t, err)
	lastEpochBlock, _, err := r.autonity.LastEpochBlock(nil)
	require.NoError(r.t, err)
	nextEpochBlock := new(big.Int).Add(epochPeriod, lastEpochBlock)
	diff := new(big.Int).Sub(nextEpochBlock, r.evm.Context.BlockNumber)
	r.waitNBlocks(int(diff.Uint64()))
}

func (r *runner) sendAUT(sender, recipient common.Address, value *big.Int) { //nolint
	//...
}

func initalizeEVM() (*vm.EVM, error) {
	ethDb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(ethDb)
	stateDB, err := state.New(common.Hash{}, db, nil)
	if err != nil {
		return nil, err
	}
	vmBlockContext := vm.BlockContext{
		Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
		CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
		BlockNumber: common.Big0,
	}
	txContext := vm.TxContext{
		Origin:   common.Address{},
		GasPrice: common.Big0,
	}
	evm := vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})
	return evm, nil
}

func setup(t *testing.T, _ *params.ChainConfig) *runner {
	evm, err := initalizeEVM()
	require.NoError(t, err)
	r := &runner{t: t, evm: evm}
	/*// todo: left for later..
	var autonityConfig AutonityConfig
	if configOverride != nil && configOverride.AutonityContractConfig != nil {
		// autonityTestConfig prob should use reflection to perform automatic assignments.
		// maybe we could make it generic just like ... operator in js
		autonityConfig = autonityTestConfig(configOverride.AutonityContractConfig)
	} else {
		autonityConfig = autonityTestConfig(params.TestAutonityContractConfig)
	}
	*/
	//
	// Step 1: Autonity Contract Deployment
	//
	validators := make([]AutonityValidator, 0, len(params.TestAutonityContractConfig.Validators))
	for _, v := range params.TestAutonityContractConfig.Validators {
		validators = append(validators, genesisToAutonityVal(v))
	}
	_, _, r.autonity, err = r.deployAutonity(nil, validators, defaultAutonityConfig)
	require.NoError(t, err)
	require.Equal(t, r.autonity.address, params.AutonityContractAddress)
	_, err = r.autonity.FinalizeInitialization(nil)
	require.NoError(t, err)
	//
	// Step 2: Accountability Contract Deployment
	//
	_, _, r.accountability, err = r.deployAccountability(nil, r.autonity.address, AccountabilityConfig{
		InnocenceProofSubmissionWindow: big.NewInt(int64(params.DefaultAccountabilityConfig.InnocenceProofSubmissionWindow)),
		BaseSlashingRateLow:            big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateLow)),
		BaseSlashingRateMid:            big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateMid)),
		CollusionFactor:                big.NewInt(int64(params.DefaultAccountabilityConfig.CollusionFactor)),
		HistoryFactor:                  big.NewInt(int64(params.DefaultAccountabilityConfig.HistoryFactor)),
		JailFactor:                     big.NewInt(int64(params.DefaultAccountabilityConfig.JailFactor)),
		SlashingRatePrecision:          big.NewInt(int64(params.DefaultAccountabilityConfig.SlashingRatePrecision)),
	})
	require.NoError(t, err)
	require.Equal(t, r.accountability.address, params.AccountabilityContractAddress)
	//
	// Step 3: Oracle contract deployment
	//
	voters := make([]common.Address, len(params.TestAutonityContractConfig.Validators))
	for _, val := range params.TestAutonityContractConfig.Validators {
		voters = append(voters, val.OracleAddress)
	}
	_, _, r.oracle, err = r.deployOracle(nil,
		voters,
		r.autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount,
		params.DefaultGenesisOracleConfig.Symbols,
		new(big.Int).SetUint64(params.DefaultGenesisOracleConfig.VotePeriod))
	require.NoError(t, err)
	require.Equal(t, r.oracle.address, params.OracleContractAddress)
	//
	// Step 4: ACU deployment
	//
	bigQuantities := make([]*big.Int, len(params.DefaultAcuContractGenesis.Quantities))
	for i := range params.DefaultAcuContractGenesis.Quantities {
		bigQuantities[i] = new(big.Int).SetUint64(params.DefaultAcuContractGenesis.Quantities[i])
	}
	_, _, r.acu, err = r.deployACU(nil,
		params.DefaultAcuContractGenesis.Symbols,
		bigQuantities,
		new(big.Int).SetUint64(params.DefaultAcuContractGenesis.Scale),
		r.autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount,
		r.oracle.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.oracle.address, params.OracleContractAddress)
	//
	// Step 5: Supply Control Deployment
	//
	_, _, r.supplyControl, err = r.deploySupplyControl(&runOptions{value: (*big.Int)(params.DefaultSupplyControlGenesis.InitialAllocation)},
		r.autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount,
		params.StabilizationContractAddress)
	require.NoError(t, err)
	require.Equal(t, r.supplyControl.address, params.SupplyControlContractAddress)
	//
	// Step 6: Stabilization Control Deployment
	//
	_, _, r.stabilization, err = r.deployStabilization(nil,
		StabilizationConfig{
			BorrowInterestRate:        (*big.Int)(params.DefaultStabilizationGenesis.BorrowInterestRate),
			LiquidationRatio:          (*big.Int)(params.DefaultStabilizationGenesis.LiquidationRatio),
			MinCollateralizationRatio: (*big.Int)(params.DefaultStabilizationGenesis.MinCollateralizationRatio),
			MinDebtRequirement:        (*big.Int)(params.DefaultStabilizationGenesis.MinDebtRequirement),
			TargetPrice:               (*big.Int)(params.DefaultStabilizationGenesis.TargetPrice),
		}, params.AutonityContractAddress,
		defaultAutonityConfig.Protocol.OperatorAccount,
		r.oracle.address,
		r.supplyControl.address,
		r.autonity.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.stabilization.address, params.StabilizationContractAddress)
	//
	// Step 7: Upgrade Manager contract deployment
	//
	_, _, r.UpgradeManager, err = r.deployUpgradeManager(nil,
		r.autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount)
	require.NoError(t, err)
	require.Equal(t, r.UpgradeManager.address, params.UpgradeManagerContractAddress)
	return r
}

// temporary until we find a better solution
func genesisToAutonityVal(v *params.Validator) AutonityValidator {
	return AutonityValidator{
		Treasury:                 v.Treasury,
		NodeAddress:              *v.NodeAddress,
		OracleAddress:            v.OracleAddress,
		Enode:                    v.Enode,
		CommissionRate:           v.CommissionRate,
		BondedStake:              v.BondedStake,
		UnbondingStake:           v.UnbondingStake,
		UnbondingShares:          v.UnbondingShares,
		SelfBondedStake:          v.SelfBondedStake,
		SelfUnbondingStake:       v.SelfUnbondingStake,
		SelfUnbondingShares:      v.SelfUnbondingShares,
		SelfUnbondingStakeLocked: v.SelfUnbondingStakeLocked,
		LiquidContract:           *v.LiquidContract,
		LiquidSupply:             v.LiquidSupply,
		RegistrationBlock:        v.RegistrationBlock,
		TotalSlashed:             v.TotalSlashed,
		JailReleaseBlock:         v.JailReleaseBlock,
		ProvableFaultCount:       v.ProvableFaultCount,
		State:                    0,
	}
}
