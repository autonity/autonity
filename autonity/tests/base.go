package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/eth/tracers"
	"github.com/autonity/autonity/params"

	_ "github.com/autonity/autonity/eth/tracers/native" //nolint
)

var fromAutonity = &runOptions{origin: params.AutonityContractAddress}

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
	var tracer tracers.Tracer
	if c.r.tracing {
		tracer, _ = tracers.New("callTracer", new(tracers.Context))
		c.r.evm.Config = vm.Config{Debug: true, Tracer: tracer}
	}
	input, err := c.abi.Pack(method, params...)
	require.NoError(c.r.t, err)
	out, consumed, err := c.r.call(opts, c.address, input)
	if c.r.tracing {
		traceResult, err := tracer.GetResult()
		require.NoError(c.r.t, err)
		pretty, _ := json.MarshalIndent(traceResult, "", "    ")
		fmt.Println(string(pretty))
	}
	if err != nil {
		reason, _ := abi.UnpackRevert(out)
		return nil, 0, fmt.Errorf("%w: %s", err, reason)
	}
	res, err := c.abi.Unpack(method, out)
	require.NoError(c.r.t, err)
	return res, consumed, nil
}

// call a contract function and then revert. helpful to get output of the function without changing state.
// similar to making a method.call() in truffle
func (c *contract) SimulateCall(methodHouse *contract, opts *runOptions, method string, params ...any) ([]any, uint64, error) {
	snap := c.r.snapshot()
	out, consumed, err := c.CallMethod(methodHouse, opts, method, params...)
	c.r.revertSnapshot(snap)
	return out, consumed, err
}

// call a method that does not belong to the contract, `c`.
// instead the method can be found in the contract, `methodHouse`.
func (c *contract) CallMethod(methodHouse *contract, opts *runOptions, method string, params ...any) ([]any, uint64, error) {
	var tracer tracers.Tracer
	if c.r.tracing {
		tracer, _ = tracers.New("callTracer", new(tracers.Context))
		c.r.evm.Config = vm.Config{Debug: true, Tracer: tracer}
	}
	input, err := methodHouse.abi.Pack(method, params...)
	require.NoError(c.r.t, err)
	out, consumed, err := c.r.call(opts, c.address, input)
	if c.r.tracing {
		traceResult, err := tracer.GetResult()
		require.NoError(c.r.t, err)
		pretty, _ := json.MarshalIndent(traceResult, "", "    ")
		fmt.Println(string(pretty))
	}
	if err != nil {
		reason, _ := abi.UnpackRevert(out)
		return nil, 0, fmt.Errorf("%w: %s", err, reason)
	}
	res, err := methodHouse.abi.Unpack(method, out)
	require.NoError(c.r.t, err)
	return res, consumed, nil
}

type Committee struct {
	validators           []AutonityValidator
	liquidStateContracts []*ILiquidLogic
}

type runner struct {
	t       *testing.T
	evm     *vm.EVM
	origin  common.Address // session's sender, can be overridden via runOptions
	tracing bool

	// protocol contracts
	// todo: see if genesis deployment flow can be abstracted somehow
	autonity               *Autonity
	accountability         *Accountability
	oracle                 *Oracle
	acu                    *ACU
	supplyControl          *SupplyControl
	stabilization          *Stabilization
	upgradeManager         *UpgradeManager
	inflationController    *InflationController
	stakableVesting        *StakableVesting
	nonStakableVesting     *NonStakableVesting
	omissionAccountability *OmissionAccountability

	committee Committee   // genesis validators for easy access
	operator  *runOptions // operator runOptions for easy access
}

func (r *runner) CallNoError(output []any, gasConsumed uint64, err error) ([]any, uint64) {
	require.NoError(r.t, err)
	return output, gasConsumed
}

func (r *runner) NoError(gasConsumed uint64, err error) uint64 {
	require.NoError(r.t, err)
	return gasConsumed
}

// returns an object of LiquidLogic contract with address set to 0
func (r *runner) LiquidLogicContractObject() *LiquidLogic {
	parsed, err := LiquidLogicMetaData.GetAbi()
	require.NoError(r.t, err)
	require.NotEqual(r.t, nil, parsed)
	return &LiquidLogic{
		contract: &contract{
			common.Address{},
			parsed,
			r,
		},
	}
}

func (r *runner) liquidStateContract(v AutonityValidator) *ILiquidLogic {
	abi, err := ILiquidLogicMetaData.GetAbi()
	require.NoError(r.t, err)
	return &ILiquidLogic{&contract{v.LiquidStateContract, abi, r}}
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

// generates an activity proof signed by all committee members, `absentees` excluded
// NOTE: if additional validators whose key is not in params.TestConsensusKey are registered in the tests,
// then this func needs to be modified to add their signatures as well.
func activityProof(committee []AutonityValidator, headerSeal common.Hash, absentees map[common.Address]struct{}) *types.AggregateSignature {
	var signatures []blst.Signature //nolint
	signers := types.NewSigners(len(committee))
	numSigners := 0
	for _, keyHex := range params.TestConsensusKeys {
		// deserialize key
		key, err := blst.SecretKeyFromHex(keyHex)
		if err != nil {
			panic(err)
		}

		// find index in committee, and skip if not in it
		index := -1
		for j, member := range committee {
			if bytes.Equal(member.ConsensusKey, key.PublicKey().Marshal()) {
				index = j
			}
		}
		if index == -1 { // not in the committee
			continue
		}

		// skip if absent
		if _, isAbsent := absentees[committee[index].NodeAddress]; isAbsent {
			continue
		}
		signatures = append(signatures, key.Sign(headerSeal[:]))

		signers.Bits.Set(index, 1)
		numSigners++
	}
	// if there are no signers, return an empty proof
	if numSigners == 0 {
		return nil
	}

	aggregateSig := blst.AggregateSignatures(signatures)

	return types.NewAggregateSignature(aggregateSig.(*blst.BlsSignature), signers)
}

// sets up an activity proof, `absentees` are excluded from it
func (r *runner) setupActivityProofAndCoinbase(proposer common.Address, absentees map[common.Address]struct{}) {
	// initialize empty map if absentees was left nil
	if absentees == nil {
		absentees = make(map[common.Address]struct{})
	}
	epochInfo, _, err := r.autonity.GetEpochInfo(nil)
	require.NoError(r.t, err)

	mustBeEmpty := r.evm.Context.BlockNumber.Uint64() <= epochInfo.EpochBlock.Uint64()+epochInfo.Delta.Uint64()
	if !mustBeEmpty {
		r.evm.Context.Coinbase = proposer
		targetHeight := r.evm.Context.BlockNumber.Uint64() - epochInfo.Delta.Uint64()

		r.evm.Context.ActivityProofRound = 0
		r.evm.Context.ActivityProof = activityProof(r.committee.validators, sealFaker(targetHeight, r.evm.Context.ActivityProofRound), absentees)
	}
}

func (r *runner) waitNBlocks(n int) { //nolint
	start := r.evm.Context.BlockNumber
	epochID, _, err := r.autonity.EpochID(nil)
	require.NoError(r.t, err)
	for i := 0; i < n; i++ {
		// set validator 0 as proposer always
		r.setupActivityProofAndCoinbase(r.committee.validators[0].NodeAddress, nil)
		// Finalize is not the only block closing operation - fee redistribution is missing and prob
		// other stuff. Left as todo.
		_, err = r.autonity.Finalize(&runOptions{origin: common.Address{}})
		// consider monitoring gas cost here and fail if it's too much
		require.NoError(r.t, err, "finalize function error in waitNblocks - ", i, r.evm.Context.BlockNumber.Uint64())
		r.evm.Context.BlockNumber = new(big.Int).Add(big.NewInt(int64(i+1)), start)
		r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
		// clean up activity proof related data
		r.evm.Context.ActivityProof = nil
		r.evm.Context.ActivityProofRound = 0
		r.evm.Context.Coinbase = common.Address{}
	}
	newEpochID, _, err := r.autonity.EpochID(nil)
	require.NoError(r.t, err)
	if newEpochID.Cmp(epochID) != 0 {
		r.generateNewCommittee()
	}
}

func (r *runner) waitNextEpoch() { //nolint
	epochInfo, _, err := r.autonity.GetEpochInfo(nil)
	require.NoError(r.t, err)

	diff := new(big.Int).Sub(epochInfo.NextEpochBlock, r.evm.Context.BlockNumber)
	r.waitNBlocks(int(diff.Uint64() + 1))
}

func (r *runner) generateNewCommittee() {
	committeeMembers, _, err := r.autonity.GetCommittee(nil)
	require.NoError(r.t, err)
	r.committee.validators = make([]AutonityValidator, len(committeeMembers))
	r.committee.liquidStateContracts = make([]*ILiquidLogic, len(committeeMembers))
	for i, member := range committeeMembers {
		validator, _, err := r.autonity.GetValidator(nil, member.Addr)
		require.NoError(r.t, err)
		r.committee.validators[i] = validator
		r.committee.liquidStateContracts[i] = r.liquidStateContract(validator)
	}
}

func (r *runner) waitSomeBlock(endTime int64) int64 { //nolint
	// because we have 1 block/s
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

func hashFaker(h uint64) common.Hash {
	return common.BytesToHash(new(big.Int).SetUint64(h).Bytes())
}

func sealFaker(targetHeight uint64, round uint64) common.Hash {
	return message.PrepareCommittedSeal(hashFaker(targetHeight), int64(round), new(big.Int).SetUint64(targetHeight))
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
		// used by the AbsenteeComputer precompile to verify activity proofs
		GetHash: hashFaker,
	}
	txContext := vm.TxContext{
		Origin:   common.Address{},
		GasPrice: common.Big0,
	}
	evm := vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})
	return evm, nil
}

func copyConfig(original *params.AutonityContractGenesis) *params.AutonityContractGenesis {
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		panic("cannot marshal autonity genesis config: " + err.Error())
	}
	genesisCopy := &params.AutonityContractGenesis{}
	err = json.Unmarshal(jsonBytes, genesisCopy)
	if err != nil {
		panic("cannot unmarshal autonity genesis config: " + err.Error())
	}
	return genesisCopy
}

func setup(t *testing.T, configOverride func(*params.AutonityContractGenesis) *params.AutonityContractGenesis) *runner {
	evm, err := initializeEVM()
	require.NoError(t, err)
	r := &runner{t: t, evm: evm}

	var autonityGenesis *params.AutonityContractGenesis
	if configOverride != nil {
		autonityGenesis = configOverride(copyConfig(params.TestAutonityContractConfig))
	} else {
		autonityGenesis = params.TestAutonityContractConfig
	}
	//TODO: implement override also for the other contracts

	//
	// Step 1: Autonity Contract Deployment
	//

	// TODO: replicate truffle tests default config.
	autonityConfig := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:             new(big.Int).SetUint64(autonityGenesis.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(autonityGenesis.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(autonityGenesis.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(autonityGenesis.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(autonityGenesis.InitialInflationReserve),
			WithholdingThreshold:    new(big.Int).SetUint64(autonityGenesis.WithholdingThreshold),
			ProposerRewardRate:      new(big.Int).SetUint64(autonityGenesis.ProposerRewardRate),
			WithheldRewardsPool:     autonityGenesis.Operator,
			TreasuryAccount:         autonityGenesis.Operator,
		},
		Contracts: AutonityContracts{
			AccountabilityContract:         params.AccountabilityContractAddress,
			OracleContract:                 params.OracleContractAddress,
			AcuContract:                    params.ACUContractAddress,
			SupplyControlContract:          params.SupplyControlContractAddress,
			StabilizationContract:          params.StabilizationContractAddress,
			UpgradeManagerContract:         params.UpgradeManagerContractAddress,
			InflationControllerContract:    params.InflationControllerContractAddress,
			OmissionAccountabilityContract: params.OmissionAccountabilityContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount: autonityGenesis.Operator,
			EpochPeriod:     new(big.Int).SetUint64(autonityGenesis.EpochPeriod),
			BlockPeriod:     new(big.Int).SetUint64(autonityGenesis.BlockPeriod),
			CommitteeSize:   new(big.Int).SetUint64(autonityGenesis.MaxCommitteeSize),
		},
		ContractVersion: big.NewInt(1),
	}
	r.operator = &runOptions{origin: autonityConfig.Protocol.OperatorAccount}
	r.committee.validators = make([]AutonityValidator, 0, len(autonityGenesis.Validators))
	for _, v := range autonityGenesis.Validators {
		validator := genesisToAutonityVal(v)
		r.committee.validators = append(r.committee.validators, validator)
	}
	_, _, r.autonity, err = r.deployAutonity(nil, r.committee.validators, autonityConfig)
	require.NoError(t, err)
	require.Equal(t, r.autonity.address, params.AutonityContractAddress)
	_, err = r.autonity.FinalizeInitialization(nil, new(big.Int).SetUint64(params.DefaultOmissionAccountabilityConfig.Delta))
	require.NoError(t, err)
	r.committee.liquidStateContracts = make([]*ILiquidLogic, 0, len(autonityGenesis.Validators))
	for _, v := range autonityGenesis.Validators {
		validator, _, err := r.autonity.GetValidator(nil, *v.NodeAddress)
		require.NoError(r.t, err)
		r.committee.liquidStateContracts = append(r.committee.liquidStateContracts, r.liquidStateContract(validator))
	}
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
	})
	require.NoError(t, err)
	require.Equal(t, r.accountability.address, params.AccountabilityContractAddress)
	//
	// Step 3: Oracle contract deployment
	//
	voters := make([]common.Address, len(autonityGenesis.Validators))
	for _, val := range autonityGenesis.Validators {
		voters = append(voters, val.OracleAddress)
	}
	_, _, r.oracle, err = r.deployOracle(nil,
		voters,
		r.autonity.address,
		autonityConfig.Protocol.OperatorAccount,
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
		autonityConfig.Protocol.OperatorAccount,
		r.oracle.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.oracle.address, params.OracleContractAddress)
	//
	// Step 5: Supply Control Deployment
	//
	r.evm.StateDB.AddBalance(common.Address{}, (*big.Int)(params.DefaultSupplyControlGenesis.InitialAllocation))
	_, _, r.supplyControl, err = r.deploySupplyControl(&runOptions{value: (*big.Int)(params.DefaultSupplyControlGenesis.InitialAllocation)},
		r.autonity.address,
		autonityConfig.Protocol.OperatorAccount,
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
		autonityConfig.Protocol.OperatorAccount,
		r.oracle.address,
		r.supplyControl.address,
		r.autonity.address,
	)
	require.NoError(t, err)
	require.Equal(t, r.stabilization.address, params.StabilizationContractAddress)
	//
	// Step 7: Upgrade Manager contract deployment
	//
	_, _, r.upgradeManager, err = r.deployUpgradeManager(nil,
		r.autonity.address,
		autonityConfig.Protocol.OperatorAccount)
	require.NoError(t, err)
	require.Equal(t, r.upgradeManager.address, params.UpgradeManagerContractAddress)

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
	_, _, r.inflationController, err = r.deployInflationController(nil, *p)
	require.NoError(r.t, err)
	require.Equal(t, r.inflationController.address, params.InflationControllerContractAddress)

	//
	// Step 9: Stakable Vesting contract deployment
	//
	_, _, r.stakableVesting, err = r.deployStakableVesting(
		nil,
		r.autonity.address,
		autonityConfig.Protocol.OperatorAccount,
	)
	require.NoError(t, err)
	require.Equal(t, r.stakableVesting.address, params.StakableVestingContractAddress)
	r.NoError(
		r.autonity.Mint(r.operator, r.stakableVesting.address, params.DefaultStakableVestingGenesis.TotalNominal),
	)
	r.NoError(
		r.stakableVesting.SetTotalNominal(r.operator, params.DefaultStakableVestingGenesis.TotalNominal),
	)

	//
	// Step 10: Non-Stakable Vesting contract deployment
	//
	_, _, r.nonStakableVesting, err = r.deployNonStakableVesting(
		nil,
		r.autonity.address,
		autonityConfig.Protocol.OperatorAccount,
	)
	require.NoError(t, err)
	require.Equal(t, r.nonStakableVesting.address, params.NonStakableVestingContractAddress)
	r.NoError(
		r.nonStakableVesting.SetTotalNominal(r.operator, params.DefaultNonStakableVestingGenesis.TotalNominal),
	)
	r.NoError(
		r.nonStakableVesting.SetMaxAllowedDuration(r.operator, params.DefaultNonStakableVestingGenesis.MaxAllowedDuration),
	)

	//
	// Step 11: Omission Accountability Contract Deployment
	//
	treasuries := make([]common.Address, len(autonityGenesis.Validators))
	for i, val := range autonityGenesis.Validators {
		treasuries[i] = val.Treasury
	}
	_, _, r.omissionAccountability, err = r.deployOmissionAccountability(nil, r.autonity.address, autonityConfig.Protocol.OperatorAccount, treasuries, OmissionAccountabilityConfig{
		InactivityThreshold:    big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InactivityThreshold)),
		LookbackWindow:         big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.LookbackWindow)),
		PastPerformanceWeight:  big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.PastPerformanceWeight)),
		InitialJailingPeriod:   big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InitialJailingPeriod)),
		InitialProbationPeriod: big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InitialProbationPeriod)),
		InitialSlashingRate:    big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InitialSlashingRate)),
		Delta:                  big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.Delta)),
	})
	require.NoError(t, err)
	require.Equal(t, r.omissionAccountability.address, params.OmissionAccountabilityContractAddress)

	// set protocol contracts
	r.NoError(
		r.autonity.SetAccountabilityContract(r.operator, r.accountability.address),
	)
	r.NoError(
		r.autonity.SetAcuContract(r.operator, r.acu.address),
	)
	r.NoError(
		r.autonity.SetInflationControllerContract(r.operator, r.inflationController.address),
	)
	r.NoError(
		r.autonity.SetOracleContract(r.operator, r.oracle.address),
	)
	r.NoError(
		r.autonity.SetStabilizationContract(r.operator, r.stabilization.address),
	)
	r.NoError(
		r.autonity.SetSupplyControlContract(r.operator, r.supplyControl.address),
	)
	r.NoError(
		r.autonity.SetUpgradeManagerContract(r.operator, r.upgradeManager.address),
	)
	r.NoError(
		r.autonity.SetNonStakableVestingContract(r.operator, r.nonStakableVesting.address),
	)
	r.NoError(
		r.autonity.SetOmissionAccountabilityContract(r.operator, r.omissionAccountability.address),
	)

	r.evm.Context.BlockNumber = common.Big1
	r.evm.Context.Time = new(big.Int).Add(r.evm.Context.Time, common.Big1)
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
