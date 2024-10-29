package tests

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/eth/tracers"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"

	_ "github.com/autonity/autonity/eth/tracers/native" //nolint
)

var (
	Operator     = &runOptions{origin: defaultAutonityConfig.Protocol.OperatorAccount}
	FromAutonity = &runOptions{origin: params.AutonityContractAddress}
	User         = common.HexToAddress("0x99")

	// todo: replicate truffle tests default config.
	defaultAutonityConfig = AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:             new(big.Int).SetUint64(params.TestAutonityContractConfig.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(params.TestAutonityContractConfig.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(params.TestAutonityContractConfig.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(params.TestAutonityContractConfig.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(params.TestAutonityContractConfig.InitialInflationReserve),
			TreasuryAccount:         params.TestAutonityContractConfig.Operator,
		},
		Contracts: AutonityContracts{
			AccountabilityContract:      params.AccountabilityContractAddress,
			OracleContract:              params.OracleContractAddress,
			AcuContract:                 params.ACUContractAddress,
			SupplyControlContract:       params.SupplyControlContractAddress,
			StabilizationContract:       params.StabilizationContractAddress,
			UpgradeManagerContract:      params.UpgradeManagerContractAddress,
			InflationControllerContract: params.InflationControllerContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount:     params.TestAutonityContractConfig.Operator,
			EpochPeriod:         new(big.Int).SetUint64(params.TestAutonityContractConfig.EpochPeriod),
			BlockPeriod:         new(big.Int).SetUint64(params.TestAutonityContractConfig.BlockPeriod),
			CommitteeSize:       new(big.Int).SetUint64(params.TestAutonityContractConfig.MaxCommitteeSize),
			MaxScheduleDuration: new(big.Int).SetUint64(params.TestAutonityContractConfig.MaxScheduleDuration),
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
	r       *Runner
}

func (c *contract) Address() common.Address {
	return c.address
}

func (c *contract) call(opts *runOptions, method string, params ...any) ([]any, uint64, error) {
	var tracer tracers.Tracer
	if c.r.Tracing {
		tracer, _ = tracers.New("callTracer", new(tracers.Context))
		c.r.Evm.Config = vm.Config{Debug: true, Tracer: tracer}
	}
	input, err := c.abi.Pack(method, params...)
	require.NoError(c.r.T, err)
	out, consumed, err := c.r.call(opts, c.address, input)
	if c.r.Tracing {
		traceResult, err := tracer.GetResult()
		require.NoError(c.r.T, err)
		pretty, _ := json.MarshalIndent(traceResult, "", "    ")
		fmt.Println(string(pretty))
	}
	if err != nil {
		reason, _ := abi.UnpackRevert(out)
		return nil, 0, fmt.Errorf("%w: %s", err, reason)
	}
	res, err := c.abi.Unpack(method, out)
	require.NoError(c.r.T, err)
	return res, consumed, nil
}

// call a method that does not belong to the contract, `c`.
// instead the method can be found in the contract, `methodHouse`.
func (c *contract) CallMethod(methodHouse *contract, opts *runOptions, method string, params ...any) ([]any, uint64, error) {
	var tracer tracers.Tracer
	if c.r.Tracing {
		tracer, _ = tracers.New("callTracer", new(tracers.Context))
		c.r.Evm.Config = vm.Config{Debug: true, Tracer: tracer}
	}
	input, err := methodHouse.abi.Pack(method, params...)
	require.NoError(c.r.T, err)
	out, consumed, err := c.r.call(opts, c.address, input)
	if c.r.Tracing {
		traceResult, err := tracer.GetResult()
		require.NoError(c.r.T, err)
		pretty, _ := json.MarshalIndent(traceResult, "", "    ")
		fmt.Println(string(pretty))
	}
	if err != nil {
		reason, _ := abi.UnpackRevert(out)
		return nil, 0, fmt.Errorf("%w: %s", err, reason)
	}
	res, err := methodHouse.abi.Unpack(method, out)
	require.NoError(c.r.T, err)
	return res, consumed, nil
}

type Committee struct {
	Validators           []AutonityValidator
	LiquidStateContracts []*ILiquid
}

type Runner struct {
	T       *testing.T
	Evm     *vm.EVM
	Origin  common.Address // session's sender, can be overridden via runOptions
	Tracing bool

	// protocol contracts
	// todo: see if genesis deployment flow can be abstracted somehow
	Autonity                *Autonity
	Accountability          *Accountability
	Oracle                  *Oracle
	Acu                     *ACU
	SupplyControl           *SupplyControl
	Stabilization           *Stabilization
	UpgradeManager          *UpgradeManager
	InflationController     *InflationController
	StakeableVestingManager *StakeableVestingManager
	NonStakeableVesting     *NonStakeableVesting

	Committee Committee // genesis validators for easy access
}

func (r *Runner) NoError(gasConsumed uint64, err error) uint64 {
	require.NoError(r.T, err)
	return gasConsumed
}

func (r *Runner) LiquidStateContract(validatorAddress common.Address) *ILiquid {
	validator, _, err := r.Autonity.GetValidator(nil, validatorAddress)
	require.NoError(r.T, err)
	return &ILiquid{r.contractObject(ILiquidMetaData, validator.LiquidStateContract)}
}

func (r *Runner) call(opts *runOptions, addr common.Address, input []byte) ([]byte, uint64, error) {
	r.Evm.Origin = r.Origin
	value := common.Big0
	if opts != nil {
		r.Evm.Origin = opts.origin
		if opts.value != nil {
			value = opts.value
		}
	}
	gas := uint64(math.MaxUint64)
	ret, leftOver, err := r.Evm.Call(vm.AccountRef(r.Evm.Origin), addr, input, gas, value)
	return ret, gas - leftOver, err
}

func (r *Runner) snapshot() int {
	return r.Evm.StateDB.Snapshot()
}

func (r *Runner) revertSnapshot(id int) {
	r.Evm.StateDB.RevertToSnapshot(id)
}

// helpful to run a code snippet without changing the state
func (r *Runner) RunAndRevert(f func(r *Runner)) {
	context := r.Evm.Context
	snap := r.snapshot()
	committee := r.Committee
	f(r)
	r.revertSnapshot(snap)
	r.Evm.Context = context
	r.Committee = committee
}

// run is a convenience wrapper against t.run with automated state snapshot
func (r *Runner) Run(name string, f func(r *Runner)) {
	r.T.Run(name, func(t2 *testing.T) {
		t := r.T
		r.T = t2
		// in the future avoid mutating for supporting parallel testing
		context := r.Evm.Context
		snap := r.snapshot()
		committee := r.Committee
		f(r)
		r.revertSnapshot(snap)
		r.Evm.Context = context
		r.Committee = committee
		r.T = t
	})
}

func RunWithSetup(name string, setup func() *Runner, run func(r *Runner)) {
	r := setup()
	r.T.Run(name, func(t2 *testing.T) {
		run(r)
	})
}

func (r *Runner) GiveMeSomeMoney(account common.Address, amount *big.Int) {
	r.Evm.StateDB.AddBalance(account, amount)
}

func (r *Runner) GetBalanceOf(account common.Address) *big.Int {
	return r.Evm.StateDB.GetBalance(account)
}

func (r *Runner) GetNewtonBalanceOf(account common.Address) *big.Int {
	balance, _, err := r.Autonity.BalanceOf(nil, account)
	require.NoError(r.T, err)
	return balance
}

func (r *Runner) deployContract(opts *runOptions, abi *abi.ABI, bytecode []byte, params ...any) (common.Address, uint64, *contract, error) {
	args, err := abi.Pack("", params...)
	require.NoError(r.T, err)
	data := append(bytecode, args...)
	gas := uint64(math.MaxUint64)
	r.Evm.Origin = r.Origin
	value := common.Big0
	if opts != nil {
		r.Evm.Origin = opts.origin
		if opts.value != nil {
			value = opts.value
		}
	}
	_, contractAddress, leftOverGas, err := r.Evm.Create(vm.AccountRef(r.Evm.Origin), data, gas, value)
	return contractAddress, gas - leftOverGas, &contract{contractAddress, abi, r}, err
}

func (r *Runner) WaitNBlocks(n int) {
	start := r.Evm.Context.BlockNumber
	epochID, _, err := r.Autonity.EpochID(nil)
	require.NoError(r.T, err)
	for i := 0; i < n; i++ {
		// Finalize is not the only block closing operation - fee redistribution is missing and prob
		// other stuff. Left as todo.
		_, err := r.Autonity.Finalize(&runOptions{origin: common.Address{}})
		// consider monitoring gas cost here and fail if it's too much
		require.NoError(r.T, err, "finalize function error in waitNblocks", i)
		r.Evm.Context.BlockNumber = new(big.Int).Add(big.NewInt(int64(i+1)), start)
		r.Evm.Context.Time = new(big.Int).Add(r.Evm.Context.Time, common.Big1)
	}
	newEpochID, _, err := r.Autonity.EpochID(nil)
	require.NoError(r.T, err)
	if newEpochID.Cmp(epochID) != 0 {
		r.generateNewCommittee()
	}
}

func (r *Runner) WaitNextEpoch() {
	info, _, err := r.Autonity.GetEpochInfo(nil)
	require.NoError(r.T, err)

	diff := new(big.Int).Sub(info.NextEpochBlock, r.Evm.Context.BlockNumber)
	r.WaitNBlocks(int(diff.Uint64() + 1))
}

func (r *Runner) contractObject(metadata *bind.MetaData, address common.Address) *contract {
	parsed, err := metadata.GetAbi()
	require.NoError(r.T, err)
	return &contract{address, parsed, r}
}

func (r *Runner) StakeableVestingContractObject(user common.Address, contractID *big.Int) *IStakeableVesting {
	address, _, err := r.StakeableVestingManager.GetContractAccount0(nil, user, contractID)
	require.NoError(r.T, err)
	return &IStakeableVesting{
		r.contractObject(IStakeableVestingMetaData, address),
	}
}

func (r *Runner) generateNewCommittee() {
	committeeMembers, _, err := r.Autonity.GetCommittee(nil)
	require.NoError(r.T, err)
	r.Committee.Validators = make([]AutonityValidator, len(committeeMembers))
	r.Committee.LiquidStateContracts = make([]*ILiquid, len(committeeMembers))
	for i, member := range committeeMembers {
		validator, _, err := r.Autonity.GetValidator(nil, member.Addr)
		require.NoError(r.T, err)
		r.Committee.Validators[i] = validator
		r.Committee.LiquidStateContracts[i] = r.LiquidStateContract(validator.NodeAddress)
	}
}

func (r *Runner) WaitForBlocksUntil(endTime int64) int64 {
	// bcause we have 1 block/s
	r.WaitNBlocks(int(endTime) - int(r.Evm.Context.Time.Int64()))
	return r.Evm.Context.Time.Int64()
}

func (r *Runner) WaitForEpochsUntil(endTime int64) int64 {
	currentTime := r.Evm.Context.Time.Int64()
	for currentTime < endTime {
		r.WaitNextEpoch()
		currentTime = r.Evm.Context.Time.Int64()
	}
	return currentTime
}

func (r *Runner) SendAUT(sender, recipient common.Address, value *big.Int) {
	require.True(r.T, r.Evm.StateDB.GetBalance(sender).Cmp(value) >= 0, "not enough balance to transfer")
	r.Evm.StateDB.SubBalance(sender, value)
	r.Evm.StateDB.AddBalance(recipient, value)
}

type EpochReward struct {
	RewardATN *big.Int
	RewardNTN *big.Int
}

func (r *Runner) RewardsAfterOneEpoch() (rewardsToDistribute EpochReward) {
	// get supply and inflationReserve to calculate inflation reward
	supply, _, err := r.Autonity.CirculatingSupply(nil)
	require.NoError(r.T, err)
	inflationReserve, _, err := r.Autonity.InflationReserve(nil)
	require.NoError(r.T, err)
	info, _, err := r.Autonity.GetEpochInfo(nil)
	require.NoError(r.T, err)
	// get inflation reward
	lastEpochTime, _, err := r.Autonity.LastEpochTime(nil)
	require.NoError(r.T, err)
	currentEpochTime := new(big.Int).Add(lastEpochTime, new(big.Int).Sub(info.NextEpochBlock, info.EpochBlock))
	rewardsToDistribute.RewardNTN, _, err = r.InflationController.CalculateSupplyDelta(nil, supply, inflationReserve, lastEpochTime, currentEpochTime)
	require.NoError(r.T, err)
	// get atn reward
	rewardsToDistribute.RewardATN = r.GetBalanceOf(r.Autonity.Address())
	return rewardsToDistribute
}

func initializeEVM() (*vm.EVM, error) {
	ethDb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(ethDb)
	stateDB, err := state.New(common.Hash{}, db, nil)
	if err != nil {
		return nil, err
	}
	vmBlockContext := vm.BlockContext{
		Transfer: func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
			db.SubBalance(sender, amount)
			db.AddBalance(recipient, amount)
		},
		CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool {
			return db.GetBalance(addr).Cmp(amount) >= 0
		},
		BlockNumber: common.Big0,
		Time:        big.NewInt(time.Now().Unix()),
	}
	txContext := vm.TxContext{
		Origin:   common.Address{},
		GasPrice: common.Big0,
	}
	evm := vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})
	return evm, nil
}

func Setup(t *testing.T, _ *params.ChainConfig) *Runner {
	evm, err := initializeEVM()
	require.NoError(t, err)
	r := &Runner{T: t, Evm: evm}
	/*// todo: left for later..
	var autonityConfig AutonityConfig
	if configOverride != nil && configOverride.AutonityContractConfig != nil {
		// autonityTestConfig prob should use reflection to perform automatic assignments.
		// maybe we could make it generic just like ... Operator in js
		autonityConfig = autonityTestConfig(configOverride.AutonityContractConfig)
	} else {
		autonityConfig = autonityTestConfig(params.TestAutonityContractConfig)
	}
	*/
	//
	// Step 1: Autonity Contract Deployment
	//
	r.Committee.Validators = make([]AutonityValidator, 0, len(params.TestAutonityContractConfig.Validators))
	for _, v := range params.TestAutonityContractConfig.Validators {
		validator := genesisToAutonityVal(v)
		r.Committee.Validators = append(r.Committee.Validators, validator)
	}
	_, _, r.Autonity, err = r.DeployAutonity(nil, r.Committee.Validators, defaultAutonityConfig)
	require.NoError(t, err)
	require.Equal(t, r.Autonity.address, params.AutonityContractAddress)
	_, err = r.Autonity.FinalizeInitialization(nil)
	require.NoError(t, err)
	r.Committee.LiquidStateContracts = make([]*ILiquid, 0, len(params.TestAutonityContractConfig.Validators))
	for _, v := range params.TestAutonityContractConfig.Validators {
		validator, _, err := r.Autonity.GetValidator(nil, *v.NodeAddress)
		require.NoError(r.T, err)
		r.Committee.LiquidStateContracts = append(r.Committee.LiquidStateContracts, r.LiquidStateContract(validator.NodeAddress))
	}
	//
	// Step 2: Accountability Contract Deployment
	//
	_, _, r.Accountability, err = r.DeployAccountability(nil, r.Autonity.address, AccountabilityConfig{
		InnocenceProofSubmissionWindow: big.NewInt(int64(params.DefaultAccountabilityConfig.InnocenceProofSubmissionWindow)),
		BaseSlashingRateLow:            big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateLow)),
		BaseSlashingRateMid:            big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateMid)),
		CollusionFactor:                big.NewInt(int64(params.DefaultAccountabilityConfig.CollusionFactor)),
		HistoryFactor:                  big.NewInt(int64(params.DefaultAccountabilityConfig.HistoryFactor)),
		JailFactor:                     big.NewInt(int64(params.DefaultAccountabilityConfig.JailFactor)),
		SlashingRatePrecision:          big.NewInt(int64(params.DefaultAccountabilityConfig.SlashingRatePrecision)),
	})
	require.NoError(t, err)
	require.Equal(t, r.Accountability.address, params.AccountabilityContractAddress)
	//
	// Step 3: Oracle contract deployment
	//
	voters := make([]common.Address, len(params.TestAutonityContractConfig.Validators))
	for _, val := range params.TestAutonityContractConfig.Validators {
		voters = append(voters, val.OracleAddress)
	}
	_, _, r.Oracle, err = r.DeployOracle(nil,
		voters,
		r.Autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount,
		params.DefaultGenesisOracleConfig.Symbols,
		new(big.Int).SetUint64(params.DefaultGenesisOracleConfig.VotePeriod))
	require.NoError(t, err)
	require.Equal(t, r.Oracle.address, params.OracleContractAddress)
	//
	// Step 4: ACU deployment
	//
	bigQuantities := make([]*big.Int, len(params.DefaultAcuContractGenesis.Quantities))
	for i := range params.DefaultAcuContractGenesis.Quantities {
		bigQuantities[i] = new(big.Int).SetUint64(params.DefaultAcuContractGenesis.Quantities[i])
	}
	_, _, r.Acu, err = r.DeployACU(nil,
		params.DefaultAcuContractGenesis.Symbols,
		bigQuantities,
		new(big.Int).SetUint64(params.DefaultAcuContractGenesis.Scale),
		r.Autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount,
		r.Oracle.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.Oracle.address, params.OracleContractAddress)
	//
	// Step 5: Supply Control Deployment
	//
	r.Evm.StateDB.AddBalance(common.Address{}, (*big.Int)(params.DefaultSupplyControlGenesis.InitialAllocation))
	_, _, r.SupplyControl, err = r.DeploySupplyControl(&runOptions{value: (*big.Int)(params.DefaultSupplyControlGenesis.InitialAllocation)},
		r.Autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount,
		params.StabilizationContractAddress)
	require.NoError(t, err)
	require.Equal(t, r.SupplyControl.address, params.SupplyControlContractAddress)
	//
	// Step 6: Stabilization Control Deployment
	//
	_, _, r.Stabilization, err = r.DeployStabilization(nil,
		StabilizationConfig{
			BorrowInterestRate:        (*big.Int)(params.DefaultStabilizationGenesis.BorrowInterestRate),
			LiquidationRatio:          (*big.Int)(params.DefaultStabilizationGenesis.LiquidationRatio),
			MinCollateralizationRatio: (*big.Int)(params.DefaultStabilizationGenesis.MinCollateralizationRatio),
			MinDebtRequirement:        (*big.Int)(params.DefaultStabilizationGenesis.MinDebtRequirement),
			TargetPrice:               (*big.Int)(params.DefaultStabilizationGenesis.TargetPrice),
		}, params.AutonityContractAddress,
		defaultAutonityConfig.Protocol.OperatorAccount,
		r.Oracle.address,
		r.SupplyControl.address,
		r.Autonity.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.Stabilization.address, params.StabilizationContractAddress)
	//
	// Step 7: Upgrade Manager contract deployment
	//
	_, _, r.UpgradeManager, err = r.DeployUpgradeManager(nil,
		r.Autonity.address,
		defaultAutonityConfig.Protocol.OperatorAccount)
	require.NoError(t, err)
	require.Equal(t, r.UpgradeManager.address, params.UpgradeManagerContractAddress)

	//
	// Step 8: Deploy Inflation Controller
	//
	p := &InflationControllerParams{
		InflationRateInitial:      (*big.Int)(params.DefaultInflationControllerGenesis.InflationRateInitial),
		InflationRateTransition:   (*big.Int)(params.DefaultInflationControllerGenesis.InflationRateTransition),
		InflationCurveConvexity:   (*big.Int)(params.DefaultInflationControllerGenesis.InflationCurveConvexity),
		InflationTransitionPeriod: (*big.Int)(params.DefaultInflationControllerGenesis.InflationTransitionPeriod),
		InflationReserveDecayRate: (*big.Int)(params.DefaultInflationControllerGenesis.InflationReserveDecayRate),
	}
	_, _, r.InflationController, err = r.DeployInflationController(nil, *p)
	require.NoError(r.T, err)
	require.Equal(t, r.InflationController.address, params.InflationControllerContractAddress)

	//
	// Step 9: Stakeable Vesting contract deployment
	//
	_, _, r.StakeableVestingManager, err = r.DeployStakeableVestingManager(
		nil,
		r.Autonity.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.StakeableVestingManager.address, params.StakeableVestingManagerContractAddress)
	r.NoError(
		r.Autonity.Mint(Operator, r.StakeableVestingManager.address, params.DefaultStakeableVestingGenesis.TotalNominal),
	)

	//
	// Step 10: Non-Stakeable Vesting contract deployment
	//
	_, _, r.NonStakeableVesting, err = r.DeployNonStakeableVesting(
		nil,
		r.Autonity.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.NonStakeableVesting.address, params.NonStakeableVestingContractAddress)

	// set protocol contracts
	r.NoError(
		r.Autonity.SetAccountabilityContract(Operator, r.Accountability.address),
	)
	r.NoError(
		r.Autonity.SetAcuContract(Operator, r.Acu.address),
	)
	r.NoError(
		r.Autonity.SetInflationControllerContract(Operator, r.InflationController.address),
	)
	r.NoError(
		r.Autonity.SetOracleContract(Operator, r.Oracle.address),
	)
	r.NoError(
		r.Autonity.SetStabilizationContract(Operator, r.Stabilization.address),
	)
	r.NoError(
		r.Autonity.SetSupplyControlContract(Operator, r.SupplyControl.address),
	)
	r.NoError(
		r.Autonity.SetUpgradeManagerContract(Operator, r.UpgradeManager.address),
	)

	r.Evm.Context.BlockNumber = common.Big1
	r.Evm.Context.Time = new(big.Int).Add(r.Evm.Context.Time, common.Big1)
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
		LiquidStateContract:      *v.LiquidStateContract,
		LiquidSupply:             v.LiquidSupply,
		RegistrationBlock:        v.RegistrationBlock,
		TotalSlashed:             v.TotalSlashed,
		JailReleaseBlock:         v.JailReleaseBlock,
		ProvableFaultCount:       v.ProvableFaultCount,
		ConsensusKey:             v.ConsensusKey,
		State:                    *v.State,
	}
}

func FromSender(sender common.Address, value *big.Int) *runOptions {
	return &runOptions{origin: sender, value: value}
}

func NewAccusationEvent(height uint64, value common.Hash, reporter common.Address) AccountabilityEvent {
	offenderNodeKey, _ := crypto.HexToECDSA(params.TestNodeKeys[0])
	offender := crypto.PubkeyToAddress(offenderNodeKey.PublicKey)
	offenderConsensusKey, _ := blst.SecretKeyFromHex(params.TestConsensusKeys[0])
	cm := types.CommitteeMember{Address: offender, VotingPower: common.Big1, ConsensusKey: offenderConsensusKey.PublicKey(), ConsensusKeyBytes: offenderConsensusKey.PublicKey().Marshal(), Index: 0}
	signer := func(hash common.Hash) blst.Signature {
		return offenderConsensusKey.Sign(hash[:])
	}
	prevote := message.NewPrevote(0, height, value, signer, &cm, 1)

	p := &accountability.Proof{
		Type:    autonity.Accusation,
		Rule:    autonity.PVN,
		Message: prevote,
	}
	rawProof, err := rlp.EncodeToBytes(p)
	if err != nil {
		panic(err)
	}

	return AccountabilityEvent{
		EventType:      uint8(p.Type),
		Rule:           uint8(p.Rule),
		Reporter:       reporter,
		Offender:       offender,
		RawProof:       rawProof,
		Id:             common.Big0,                           // assigned contract-side
		Block:          new(big.Int).SetUint64(p.Message.H()), // assigned contract-side
		ReportingBlock: common.Big0,                           // assigned contract-side
		Epoch:          common.Big0,                           // assigned contract-side
		MessageHash:    common.Big0,                           // assigned contract-side
	}
}
