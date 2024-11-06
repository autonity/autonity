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
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/eth/tracers"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"

	_ "github.com/autonity/autonity/eth/tracers/native" //nolint
)

var (
	FromAutonity = &runOptions{origin: params.AutonityContractAddress}
	User         = common.HexToAddress("0x99")
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
	OmissionAccountability  *OmissionAccountability

	Committee Committee   // genesis validators for easy access
	Operator  *runOptions // operator runOptions for easy access
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

func (r *Runner) slasherContract() *Slasher {
	slasherAddr, _, err := r.Autonity.Slasher(nil)
	require.NoError(r.T, err)
	abi, err := SlasherMetaData.GetAbi()
	require.NoError(r.T, err)
	return &Slasher{&contract{slasherAddr, abi, r}}
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
func (r *Runner) setupActivityProofAndCoinbase(proposer common.Address, absentees map[common.Address]struct{}) {
	// if we are in TestMode, no need to set the activity proof, it will not be verified anyway
	if r.Evm.ChainConfig().TestMode {
		return
	}
	// initialize empty map if absentees was left nil
	if absentees == nil {
		absentees = make(map[common.Address]struct{})
	}
	epochInfo, _, err := r.Autonity.GetEpochInfo(nil)
	require.NoError(r.T, err)

	mustBeEmpty := r.Evm.Context.BlockNumber.Uint64() <= epochInfo.EpochBlock.Uint64()+epochInfo.Delta.Uint64()
	if !mustBeEmpty {
		r.Evm.Context.Coinbase = proposer
		targetHeight := r.Evm.Context.BlockNumber.Uint64() - epochInfo.Delta.Uint64()

		r.Evm.Context.ActivityProofRound = 0
		r.Evm.Context.ActivityProof = activityProof(r.Committee.Validators, sealFaker(targetHeight, r.Evm.Context.ActivityProofRound), absentees)
	}
}

func (r *Runner) lastMinedHeight() uint64 {
	return r.Evm.Context.BlockNumber.Uint64() - 1
}

// for omission
func (r *Runner) lastTargetHeight() uint64 {
	config, _, err := r.OmissionAccountability.Config(nil)
	require.NoError(r.T, err)
	delta := config.Delta.Uint64()

	epochBlockBig, _, err := r.Autonity.GetLastEpochBlock(nil)
	require.NoError(r.T, err)
	epochBlock := epochBlockBig.Uint64()

	lastHeight := r.lastMinedHeight()

	require.Greater(r.T, lastHeight, epochBlock+delta)
	return lastHeight - delta
}

func (r *Runner) WaitNBlocks(n int) {
	start := r.Evm.Context.BlockNumber
	epochID, _, err := r.Autonity.EpochID(nil)
	require.NoError(r.T, err)
	for i := 0; i < n; i++ {
		// set validator 0 as proposer always
		r.setupActivityProofAndCoinbase(r.Committee.Validators[0].NodeAddress, nil)
		// Finalize is not the only block closing operation - fee redistribution is missing and prob
		// other stuff. Left as todo.
		_, err := r.Autonity.Finalize(&runOptions{origin: common.Address{}})
		// consider monitoring gas cost here and fail if it's too much
		require.NoError(r.T, err, "finalize function error in waitNblocks", i)
		r.Evm.Context.BlockNumber = new(big.Int).Add(big.NewInt(int64(i+1)), start)
		r.Evm.Context.Time = new(big.Int).Add(r.Evm.Context.Time, common.Big1)
		// clean up activity proof related data
		r.Evm.Context.ActivityProof = nil
		r.Evm.Context.ActivityProofRound = 0
		r.Evm.Context.Coinbase = common.Address{}
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

func Setup(t *testing.T, configOverride func(*params.AutonityContractGenesis) *params.AutonityContractGenesis) *Runner {
	evm, err := initializeEVM()
	require.NoError(t, err)
	r := &Runner{T: t, Evm: evm}

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
			OperatorAccount:     autonityGenesis.Operator,
			EpochPeriod:         new(big.Int).SetUint64(autonityGenesis.EpochPeriod),
			BlockPeriod:         new(big.Int).SetUint64(autonityGenesis.BlockPeriod),
			CommitteeSize:       new(big.Int).SetUint64(autonityGenesis.MaxCommitteeSize),
			MaxScheduleDuration: new(big.Int).SetUint64(autonityGenesis.MaxScheduleDuration),
		},
		ContractVersion: big.NewInt(1),
	}
	r.Operator = &runOptions{origin: autonityConfig.Protocol.OperatorAccount}

	r.Committee.Validators = make([]AutonityValidator, 0, len(autonityGenesis.Validators))
	for _, v := range autonityGenesis.Validators {
		validator := genesisToAutonityVal(v)
		r.Committee.Validators = append(r.Committee.Validators, validator)
	}
	_, _, r.Autonity, err = r.DeployAutonity(nil, r.Committee.Validators, autonityConfig)
	require.NoError(t, err)
	require.Equal(t, r.Autonity.address, params.AutonityContractAddress)
	_, err = r.Autonity.FinalizeInitialization(nil, new(big.Int).SetUint64(params.DefaultOmissionAccountabilityConfig.Delta))
	require.NoError(t, err)
	r.Committee.LiquidStateContracts = make([]*ILiquid, 0, len(autonityGenesis.Validators))
	for _, v := range autonityGenesis.Validators {
		validator, _, err := r.Autonity.GetValidator(nil, *v.NodeAddress)
		require.NoError(r.T, err)
		r.Committee.LiquidStateContracts = append(r.Committee.LiquidStateContracts, r.LiquidStateContract(validator.NodeAddress))
	}
	//
	// Step 2: Accountability Contract Deployment
	//
	_, _, r.Accountability, err = r.DeployAccountability(nil, r.Autonity.address, AccountabilityConfig{
		InnocenceProofSubmissionWindow: big.NewInt(int64(params.DefaultAccountabilityConfig.InnocenceProofSubmissionWindow)),
		BaseSlashingRates: AccountabilityBaseSlashingRates{
			Low:  big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateLow)),
			Mid:  big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateMid)),
			High: big.NewInt(int64(params.DefaultAccountabilityConfig.BaseSlashingRateHigh)),
		},
		Factors: AccountabilityFactors{
			Collusion: big.NewInt(int64(params.DefaultAccountabilityConfig.CollusionFactor)),
			History:   big.NewInt(int64(params.DefaultAccountabilityConfig.HistoryFactor)),
			Jail:      big.NewInt(int64(params.DefaultAccountabilityConfig.JailFactor)),
		},
	})
	require.NoError(t, err)
	require.Equal(t, r.Accountability.address, params.AccountabilityContractAddress)
	//
	// Step 3: Oracle contract deployment
	//
	voters := make([]common.Address, len(autonityGenesis.Validators))
	for _, val := range autonityGenesis.Validators {
		voters = append(voters, val.OracleAddress)
	}
	_, _, r.Oracle, err = r.DeployOracle(nil,
		voters,
		r.Autonity.address,
		autonityConfig.Protocol.OperatorAccount,
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
		autonityConfig.Protocol.OperatorAccount,
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
		autonityConfig.Protocol.OperatorAccount,
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
		autonityConfig.Protocol.OperatorAccount,
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
		autonityConfig.Protocol.OperatorAccount)
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
		r.Autonity.Mint(r.Operator, r.StakeableVestingManager.address, params.DefaultStakeableVestingGenesis.TotalNominal),
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

	//
	// Step 11: Omission Accountability Contract Deployment
	//
	treasuries := make([]common.Address, len(autonityGenesis.Validators))
	for i, val := range autonityGenesis.Validators {
		treasuries[i] = val.Treasury
	}
	_, _, r.OmissionAccountability, err = r.DeployOmissionAccountability(nil, r.Autonity.address, autonityConfig.Protocol.OperatorAccount, treasuries, OmissionAccountabilityConfig{
		InactivityThreshold:    big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InactivityThreshold)),
		LookbackWindow:         big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.LookbackWindow)),
		PastPerformanceWeight:  big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.PastPerformanceWeight)),
		InitialJailingPeriod:   big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InitialJailingPeriod)),
		InitialProbationPeriod: big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InitialProbationPeriod)),
		InitialSlashingRate:    big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.InitialSlashingRate)),
		Delta:                  big.NewInt(int64(params.DefaultOmissionAccountabilityConfig.Delta)),
	})
	require.NoError(t, err)
	require.Equal(t, r.OmissionAccountability.address, params.OmissionAccountabilityContractAddress)

	// set protocol contracts
	r.NoError(
		r.Autonity.SetAccountabilityContract(r.Operator, r.Accountability.address),
	)
	r.NoError(
		r.Autonity.SetAcuContract(r.Operator, r.Acu.address),
	)
	r.NoError(
		r.Autonity.SetInflationControllerContract(r.Operator, r.InflationController.address),
	)
	r.NoError(
		r.Autonity.SetOracleContract(r.Operator, r.Oracle.address),
	)
	r.NoError(
		r.Autonity.SetStabilizationContract(r.Operator, r.Stabilization.address),
	)
	r.NoError(
		r.Autonity.SetSupplyControlContract(r.Operator, r.SupplyControl.address),
	)
	r.NoError(
		r.Autonity.SetUpgradeManagerContract(r.Operator, r.UpgradeManager.address),
	)
	r.NoError(
		r.Autonity.SetOmissionAccountabilityContract(r.Operator, r.OmissionAccountability.address),
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
