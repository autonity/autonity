package tests

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/params/generated"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"

	_ "github.com/autonity/autonity/eth/tracers/native" //nolint
)

var (
	operator     = &runOptions{origin: params.TestAutonityContractConfig.Operator}
	fromAutonity = &runOptions{origin: params.AutonityContractAddress}
)

type runOptions struct {
	origin common.Address
	value  *big.Int
}

type committee struct {
	validators      []AutonityValidator
	liquidContracts []*Liquid
}

type runner struct {
	t       *testing.T
	evm     *vm.EVM
	stateDB *state.StateDB
	origin  common.Address // session's sender, can be overridden via runOptions
	tracing bool

	// protocol contracts
	// todo: see if genesis deployment flow can be abstracted somehow
	autonity            *AutonityTest
	accountability      *Accountability
	oracle              *Oracle
	acu                 *ACU
	supplyControl       *SupplyControl
	stabilization       *Stabilization
	upgradeManager      *UpgradeManager
	inflationController *InflationController
	stakableVesting     *StakableVesting
	nonStakableVesting  *NonStakableVesting

	committee committee     // genesis validators for easy access
	operator  *runOptions   // operator runOptions for easy access
	genesis   *core.Genesis // genesis config for easy access

	// temporary state
	prepared    bool
	persistLogs bool
	logs        []*types.Log
}

func (r *runner) NoError(gasConsumed uint64, err error) uint64 {
	require.NoError(r.t, err)
	return gasConsumed
}

func (r *runner) liquidContract(v AutonityValidator) *Liquid {
	abi, err := LiquidMetaData.GetAbi()
	require.NoError(r.t, err)
	return &Liquid{&contract{v.LiquidContract, abi, r}}
}

func (r *runner) call(opts *runOptions, addr common.Address, input []byte) ([]byte, uint64, error) {
	txHash := common.HexToHash("0x2f1085bc13c8a5a7d8e0544216e110fcc22b4c571657f313c9716a3bdc7453f3")

	r.evm.Origin = r.origin
	value := common.Big0
	if opts != nil {
		r.evm.Origin = opts.origin
		if opts.value != nil {
			value = opts.value
		}
	}

	sender := r.stateDB.GetOrNewStateObject(r.evm.Origin)

	// by only preparing once, we avoid the errors we encounter in the access list on revert
	if !r.prepared {
		r.stateDB.Prepare(txHash, 0)
		r.prepared = true
	}

	gas := uint64(math.MaxUint64)
	ret, leftOver, err := r.evm.Call(sender, addr, input, gas, value)

	if r.persistLogs {
		logs := r.stateDB.GetLogs(txHash, common.Hash{})
		r.logs = append(r.logs, logs...)
	}

	return ret, gas - leftOver, err
}

func (r *runner) snapshot() int {
	return r.evm.StateDB.Snapshot()
}

func (r *runner) revertSnapshot(id int) {
	r.prepared = false
	r.clearLogs()
	r.persistLogs = false
	r.evm.StateDB.RevertToSnapshot(id)
}

// run is a convenience wrapper against t.run with automated state snapshot
func (r *runner) run(name string, f func(r *runner)) {
	r.t.Run(name, func(t2 *testing.T) {
		t := r.t
		r.t = t2
		// in the future avoid mutating for supporting parallel testing
		context := r.evm.Context
		snap := r.snapshot()
		committee := r.committee
		f(r)
		r.revertSnapshot(snap)
		r.evm.Context = context
		r.committee = committee
		r.t = t
	})
}

func (r *runner) giveMeSomeMoney(user common.Address, amount *big.Int) { //nolint
	r.evm.StateDB.AddBalance(user, amount)
}

func (r *runner) randomAccount() common.Address {
	key, err := crypto.GenerateKey()
	require.NoError(r.t, err)
	address := crypto.PubkeyToAddress(key.PublicKey)
	r.giveMeSomeMoney(address, big.NewInt(1e18))
	return address
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
	epochID, _, err := r.autonity.EpochID(nil)
	require.NoError(r.t, err)
	for i := 0; i < n; i++ {
		// Finalize is not the only block closing operation - fee redistribution is missing and prob
		// other stuff. Left as todo.
		_, err := r.autonity.Finalize(&runOptions{origin: common.Address{}})
		// consider monitoring gas cost here and fail if it's too much
		require.NoError(r.t, err, "finalize function error in waitNblocks - ", i, r.evm.Context.BlockNumber.Uint64())
		r.evm.Context.BlockNumber = new(big.Int).Add(big.NewInt(int64(i+1)), start)
		r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
	}
	newEpochID, _, err := r.autonity.EpochID(nil)
	require.NoError(r.t, err)
	if newEpochID.Cmp(epochID) != 0 {
		r.generateNewCommittee()
	}
}

func (r *runner) waitNextEpoch() { //nolint
	epochPeriod, _, err := r.autonity.GetEpochPeriod(nil)
	require.NoError(r.t, err)
	lastEpochBlock, _, err := r.autonity.LastEpochBlock(nil)
	require.NoError(r.t, err)
	nextEpochBlock := new(big.Int).Add(epochPeriod, lastEpochBlock)
	diff := new(big.Int).Sub(nextEpochBlock, r.evm.Context.BlockNumber)
	r.waitNBlocks(int(diff.Uint64() + 1))
	r.generateNewCommittee()
}

func (r *runner) generateNewCommittee() {
	committeeMembers, _, err := r.autonity.GetCommittee(nil)
	require.NoError(r.t, err)
	r.committee.validators = make([]AutonityValidator, len(committeeMembers))
	r.committee.liquidContracts = make([]*Liquid, len(committeeMembers))
	for i, member := range committeeMembers {
		validator, _, err := r.autonity.GetValidator(nil, member.Addr)
		require.NoError(r.t, err)
		r.committee.validators[i] = validator
		r.committee.liquidContracts[i] = r.liquidContract(validator)
	}
}

func (r *runner) waitSomeBlock(endTime int64) int64 { //nolint
	// bcause we have 1 block/s
	r.waitNBlocks(int(endTime) - int(r.evm.Context.Time.Int64()))
	return r.evm.Context.Time.Int64()
}

func (r *runner) waitSomeEpoch(endTime int64) int64 {
	currentTime := r.evm.Context.Time.Int64()
	for currentTime < endTime {
		r.waitNextEpoch()
		currentTime = r.evm.Context.Time.Int64()
	}
	return currentTime
}

func (r *runner) sendAUT(sender, recipient common.Address, value *big.Int) { //nolint
	require.True(r.t, r.evm.StateDB.GetBalance(sender).Cmp(value) >= 0, "not enough balance to transfer")
	r.evm.StateDB.SubBalance(sender, value)
	r.evm.StateDB.AddBalance(recipient, value)
}

func initializeEVM(chainConfig *params.ChainConfig) (*vm.EVM, *state.StateDB, error) {
	ethDb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(ethDb)
	stateDB, err := state.New(common.Hash{}, db, nil)
	if err != nil {
		return nil, nil, err
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
	evm := vm.NewEVM(vmBlockContext, txContext, stateDB, chainConfig, vm.Config{})
	return evm, stateDB, nil
}

func (r *runner) keepLogs(persist bool) {
	r.persistLogs = persist
}

func (r *runner) Logs() []*types.Log {
	return r.logs
}

func (r *runner) clearLogs() {
	r.logs = nil
}

func setup(t *testing.T, configOverride func(genesis *core.Genesis) *core.Genesis) *runner {
	genesisConfig := &core.Genesis{
		Config:  params.TestChainConfig,
		BaseFee: big.NewInt(params.InitialBaseFee),
	}

	if configOverride != nil {
		genesisConfig = configOverride(genesisConfig)
	}

	evm, stateDb, err := initializeEVM(genesisConfig.Config)
	require.NoError(t, err)
	r := &runner{t: t, evm: evm, stateDB: stateDb}

	//
	// Step 1: Autonity Contract Deployment
	//

	if err := autonity.ExecuteTestGenesisSequence(genesisConfig.Config, genesisConfig.Alloc.ToGenesisBonds(), evm); err != nil {
		require.NoError(t, err)
		return nil
	}

	// TODO: replicate truffle tests default config.
	r.genesis = genesisConfig
	r.operator = operator

	// set up internal bindings
	r.autonity = &AutonityTest{&contract{params.AutonityContractAddress, &generated.AutonityTestAbi, r}}
	r.accountability = &Accountability{&contract{params.AccountabilityContractAddress, &generated.AccountabilityAbi, r}}
	r.oracle = &Oracle{&contract{params.OracleContractAddress, &generated.OracleAbi, r}}
	r.acu = &ACU{&contract{params.ACUContractAddress, &generated.ACUAbi, r}}
	r.supplyControl = &SupplyControl{&contract{params.SupplyControlContractAddress, &generated.SupplyControlAbi, r}}
	r.stabilization = &Stabilization{&contract{params.StabilizationContractAddress, &generated.StabilizationAbi, r}}
	r.upgradeManager = &UpgradeManager{&contract{params.UpgradeManagerContractAddress, &generated.UpgradeManagerAbi, r}}
	r.inflationController = &InflationController{&contract{params.InflationControllerContractAddress, &generated.InflationControllerAbi, r}}
	r.stakableVesting = &StakableVesting{&contract{params.StakableVestingContractAddress, &generated.StakableVestingAbi, r}}
	r.nonStakableVesting = &NonStakableVesting{&contract{params.NonStakableVestingContractAddress, &generated.NonStakableVestingAbi, r}}

	r.committee.liquidContracts = make([]*Liquid, 0, len(genesisConfig.Config.AutonityContractConfig.Validators))
	r.committee.validators = make([]AutonityValidator, 0, len(genesisConfig.Config.AutonityContractConfig.Validators))
	for _, v := range genesisConfig.Config.AutonityContractConfig.Validators {
		validator, _, err := r.autonity.GetValidator(nil, *v.NodeAddress)
		require.NoError(r.t, err)
		r.committee.liquidContracts = append(r.committee.liquidContracts, r.liquidContract(validator))
		r.committee.validators = append(r.committee.validators, validator)
	}

	r.evm.Context.BlockNumber = common.Big1
	r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
	return r
}
