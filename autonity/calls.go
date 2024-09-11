package autonity

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

var (
	finalizeGas      = metrics.NewRegisteredBufferedGauge("autonity/finalize", nil, nil)
	epochFinalizeGas = metrics.NewRegisteredBufferedGauge("autonity/epoch/finalize", nil, nil)
)

type raw []byte

// GenesisBonds is an intermediary struct used to pass genesis bondings.
// We cannot use autonity/core package here as it would cause import cycle
type GenesisBonds = []GenesisBond
type Delegation = struct {
	Validator common.Address
	Amount    *big.Int
}
type GenesisBond struct {
	Staker        common.Address
	NewtonBalance *big.Int
	Bonds         []Delegation
}

func DeployContracts(genesisConfig *params.ChainConfig, genesisBonds GenesisBonds, evmContracts *GenesisEVMContracts) error {
	if err := DeployAutonityContract(genesisConfig.AutonityContractConfig, genesisBonds, evmContracts, genesisConfig.OmissionAccountabilityConfig.Delta); err != nil {
		return fmt.Errorf("error when deploying the autonity contract: %w", err)
	}
	if err := DeployAccountabilityContract(genesisConfig.AccountabilityConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the accountability contract: %w", err)
	}
	if err := DeployOracleContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the oracle contract: %w", err)
	}
	if err := DeployACUContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the ACU contract: %w", err)
	}
	if err := DeploySupplyControlContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the supply control contract: %w", err)
	}
	if err := DeployStabilizationContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the stabilization contract: %w", err)
	}
	if err := DeployUpgradeManagerContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the upgrade manager contract: %w", err)
	}
	if err := DeployInflationControllerContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the inflation controller contract: %w", err)
	}
	if err := DeployStakeableVestingContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the stakeable vesting contract: %w", err)
	}
	if err := DeployNonStakeableVestingContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the non-stakeable vesting contract: %w", err)
	}
	if err := DeployOmissionAccountabilityContract(genesisConfig, evmContracts); err != nil {
		return fmt.Errorf("error when deploying the omission accountability contract: %w", err)
	}
	return nil
}

func DeployUpgradeManagerContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	err := evmContracts.DeployUpgradeManagerContract(
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		generated.UpgradeManagerBytecode)
	if err != nil {
		log.Error("DeployUpgradeManagerContract failed", "err", err)
		return fmt.Errorf("failed to deploy Upgrade Manager contract: %w", err)
	}
	log.Info("Deployed Upgrade Manager contract", "address", params.UpgradeManagerContractAddress)
	return nil
}

func DeployStabilizationContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	stabilizationConfig := StabilizationConfig{
		BorrowInterestRate:        (*big.Int)(config.ASM.StabilizationContractConfig.BorrowInterestRate),
		LiquidationRatio:          (*big.Int)(config.ASM.StabilizationContractConfig.LiquidationRatio),
		MinCollateralizationRatio: (*big.Int)(config.ASM.StabilizationContractConfig.MinCollateralizationRatio),
		MinDebtRequirement:        (*big.Int)(config.ASM.StabilizationContractConfig.MinDebtRequirement),
		TargetPrice:               (*big.Int)(config.ASM.StabilizationContractConfig.TargetPrice),
	}

	err := evmContracts.DeployStabilizationContract(stabilizationConfig,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.OracleContractAddress,
		params.SupplyControlContractAddress,
		params.AutonityContractAddress,
		generated.StabilizationBytecode)

	if err != nil {
		log.Error("DeployStabilizationContract failed", "err", err)
		return fmt.Errorf("failed to deploy Stabilization contract: %w", err)
	}

	log.Info("Deployed Stabilization contract", "address", params.StabilizationContractAddress)

	return nil
}

func DeploySupplyControlContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	value := (*big.Int)(config.ASM.SupplyControlConfig.InitialAllocation)

	evmContracts.AddBalance(params.DeployerAddress, value)
	err := evmContracts.DeploySupplyControlContract(
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.StabilizationContractAddress,
		generated.SupplyControlBytecode,
		value)

	if err != nil {
		log.Error("DeploySupplyControlContract failed", "err", err)
		return fmt.Errorf("failed to deploy SupplyControl contract: %w", err)
	}

	log.Info("Deployed ASM supply control contract", "address", params.SupplyControlContractAddress)

	return nil
}

func DeployInflationControllerContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	param := InflationControllerParams{
		InflationRateInitial:      (*big.Int)(config.InflationContractConfig.InflationRateInitial),
		InflationRateTransition:   (*big.Int)(config.InflationContractConfig.InflationRateTransition),
		InflationCurveConvexity:   (*big.Int)(config.InflationContractConfig.InflationCurveConvexity),
		InflationTransitionPeriod: (*big.Int)(config.InflationContractConfig.InflationTransitionPeriod),
		InflationReserveDecayRate: (*big.Int)(config.InflationContractConfig.InflationReserveDecayRate),
	}
	if err := evmContracts.DeployInflationControllerContract(generated.InflationControllerBytecode, param); err != nil {
		log.Error("DeployInflationControllerContract failed", "err", err)
		return fmt.Errorf("failed to deploy inflation controller contract: %w", err)
	}
	log.Info("Deployed Inflation Controller contract", "address", params.InflationControllerContractAddress)
	return nil
}

func DeployStakeableVestingContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	if err := evmContracts.DeployStakeableVestingContract(
		generated.StakeableVestingManagerBytecode, params.AutonityContractAddress,
	); err != nil {
		log.Error("DeployStakeableVestingContract failed", "err", err)
		return fmt.Errorf("failed to deploy stakeable vesting contract: %w", err)
	}
	log.Info("Deployed Stakeable Vesting contract", "address", params.StakeableVestingManagerContractAddress)
	if err := evmContracts.Mint(params.StakeableVestingManagerContractAddress, config.StakeableVestingConfig.TotalNominal); err != nil {
		return fmt.Errorf("error while minting total nominal to stakeable vesting contract: %w", err)
	}
	for _, vesting := range config.StakeableVestingConfig.StakeableContracts {
		if err := evmContracts.NewStakeableContract(vesting); err != nil {
			return fmt.Errorf("failed to create new stakeable vesting contract: %w", err)
		}
	}
	return nil
}

func DeployNonStakeableVestingContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	if err := evmContracts.DeployNonStakeableVestingContract(
		generated.NonStakeableVestingBytecode, params.AutonityContractAddress,
	); err != nil {
		log.Error("DeployNonStakeableVestingContract failed", "err", err)
		return fmt.Errorf("failed to deploy non-stakeable vesting contract: %w", err)
	}
	log.Info("Deployed Non-Stakeable Vesting contract", "address", params.NonStakeableVestingContractAddress)
	for _, vesting := range config.NonStakeableVestingConfig.NonStakeableContracts {
		if err := evmContracts.NewNonStakeableContract(vesting); err != nil {
			return fmt.Errorf("failed to create new non-stakeable vesting contract: %w", err)
		}
	}
	return nil
}

func DeployACUContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	bigQuantities := make([]*big.Int, len(config.ASM.ACUContractConfig.Quantities))
	for i := range config.ASM.ACUContractConfig.Quantities {
		bigQuantities[i] = new(big.Int).SetUint64(config.ASM.ACUContractConfig.Quantities[i])
	}

	err := evmContracts.DeployACUContract(config.ASM.ACUContractConfig.Symbols,
		bigQuantities,
		new(big.Int).SetUint64(config.ASM.ACUContractConfig.Scale),
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.OracleContractAddress,
		generated.ACUBytecode)

	if err != nil {
		log.Error("DeployACUContract failed", "err", err)
		return fmt.Errorf("failed to deploy ACU contract: %w", err)
	}

	log.Info("Deployed ACU contract", "address", params.ACUContractAddress)

	return nil
}

func DeployAccountabilityContract(config *params.AccountabilityGenesis, evmContracts *GenesisEVMContracts) error {
	accountabilityConfig := AccountabilityConfig{
		InnocenceProofSubmissionWindow: new(big.Int).SetUint64(config.InnocenceProofSubmissionWindow),
		BaseSlashingRates: AccountabilityBaseSlashingRates{
			Low:  new(big.Int).SetUint64(config.BaseSlashingRateLow),
			Mid:  new(big.Int).SetUint64(config.BaseSlashingRateMid),
			High: new(big.Int).SetUint64(config.BaseSlashingRateHigh),
		},
		Factors: AccountabilityFactors{
			Collusion: new(big.Int).SetUint64(config.CollusionFactor),
			History:   new(big.Int).SetUint64(config.HistoryFactor),
			Jail:      new(big.Int).SetUint64(config.JailFactor),
		},
	}
	err := evmContracts.DeployAccountabilityContract(params.AutonityContractAddress, accountabilityConfig, generated.AccountabilityBytecode)
	if err != nil {
		return fmt.Errorf("failed to deploy accountability contract: %w", err)
	}

	log.Info("Deployed Accountability contract", "address", params.AccountabilityContractAddress)

	return nil
}

func DeployOmissionAccountabilityContract(genesisConfig *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	config := genesisConfig.OmissionAccountabilityConfig

	conf := OmissionAccountabilityConfig{
		InactivityThreshold:    new(big.Int).SetUint64(config.InactivityThreshold),
		LookbackWindow:         new(big.Int).SetUint64(config.LookbackWindow),
		PastPerformanceWeight:  new(big.Int).SetUint64(config.PastPerformanceWeight),
		InitialJailingPeriod:   new(big.Int).SetUint64(config.InitialJailingPeriod),
		InitialProbationPeriod: new(big.Int).SetUint64(config.InitialProbationPeriod),
		InitialSlashingRate:    new(big.Int).SetUint64(config.InitialSlashingRate),
		Delta:                  new(big.Int).SetUint64(config.Delta),
	}

	treasuries := make([]common.Address, len(genesisConfig.AutonityContractConfig.Validators))
	for i, val := range genesisConfig.AutonityContractConfig.Validators {
		treasuries[i] = val.Treasury
	}

	err := evmContracts.DeployOmissionAccountabilityContract(
		params.AutonityContractAddress,
		genesisConfig.AutonityContractConfig.Operator,
		treasuries,
		conf,
		generated.OmissionAccountabilityBytecode,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy omission accountability contract: %w", err)
	}

	log.Info("Deployed Omission Accountability contract", "address", params.OmissionAccountabilityContractAddress)

	return nil
}

func DeployAutonityContract(genesisConfig *params.AutonityContractGenesis, genesisBonds GenesisBonds, evmContracts *GenesisEVMContracts, delta uint64) error {
	contractConfig := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:             new(big.Int).SetUint64(genesisConfig.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(genesisConfig.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(genesisConfig.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(genesisConfig.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(genesisConfig.InitialInflationReserve),
			WithholdingThreshold:    new(big.Int).SetUint64(genesisConfig.WithholdingThreshold),
			ProposerRewardRate:      new(big.Int).SetUint64(genesisConfig.ProposerRewardRate),
			WithheldRewardsPool:     genesisConfig.WithheldRewardsPool,
			TreasuryAccount:         genesisConfig.Treasury,
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
			OperatorAccount:     genesisConfig.Operator,
			EpochPeriod:         new(big.Int).SetUint64(genesisConfig.EpochPeriod),
			BlockPeriod:         new(big.Int).SetUint64(genesisConfig.BlockPeriod),
			CommitteeSize:       new(big.Int).SetUint64(genesisConfig.MaxCommitteeSize),
			MaxScheduleDuration: new(big.Int).SetUint64(genesisConfig.MaxScheduleDuration),
		},
		ContractVersion: big.NewInt(1),
	}
	validators := make([]params.Validator, 0, len(genesisConfig.Validators))
	for _, v := range genesisConfig.Validators {
		validators = append(validators, *v)
	}
	if err := evmContracts.DeployAutonityContract(genesisConfig.Bytecode, validators, contractConfig); err != nil {
		log.Error("DeployAutonityContract failed", "err", err)
		return fmt.Errorf("failed to deploy Autonity contract: %w", err)
	}

	for _, alloc := range genesisBonds {
		balanceToMint := new(big.Int)
		if alloc.NewtonBalance != nil {
			balanceToMint.Add(balanceToMint, alloc.NewtonBalance)
		}
		for _, delegation := range alloc.Bonds {
			balanceToMint.Add(balanceToMint, delegation.Amount)
		}
		if balanceToMint.Cmp(common.Big0) > 0 {
			err := evmContracts.Mint(alloc.Staker, balanceToMint)
			if err != nil {
				return fmt.Errorf("error while minting Newton: %w", err)
			}
			for _, delegation := range alloc.Bonds {
				err = evmContracts.Bond(alloc.Staker, delegation.Validator, delegation.Amount)
				if err != nil {
					return fmt.Errorf("error while bonding: %w", err)
				}
			}
		}
	}

	for _, schedule := range genesisConfig.Schedules {
		if schedule.VaultAddress != params.NonStakeableVestingContractAddress {
			return fmt.Errorf("vault address has no smart contract deployed")
		}
		if err := evmContracts.CreateSchedule(schedule); err != nil {
			return fmt.Errorf("error while creating schedules: %w", err)
		}
	}

	if err := evmContracts.FinalizeInitialization(delta); err != nil {
		return fmt.Errorf("error while calling finalizeInitialization: %w", err)
	}

	log.Info("Deployed Autonity contract", "address", params.AutonityContractAddress)
	return nil
}

func DeployOracleContract(genesisConfig *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	voters := make([]common.Address, len(genesisConfig.AutonityContractConfig.Validators))
	for _, val := range genesisConfig.AutonityContractConfig.Validators {
		voters = append(voters, val.OracleAddress)
	}

	err := evmContracts.DeployOracleContract(
		voters,
		params.AutonityContractAddress,
		genesisConfig.AutonityContractConfig.Operator,
		genesisConfig.OracleContractConfig.Symbols,
		new(big.Int).SetUint64(genesisConfig.OracleContractConfig.VotePeriod),
		genesisConfig.OracleContractConfig.Bytecode,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Oracle contract: %w", err)
	}

	log.Info("Deployed Oracle Contract", "address", params.OracleContractAddress)
	return nil
}

func (c *EVMContract) replaceAutonityBytecode(header *types.Header, statedb vm.StateDB, bytecode []byte) error {
	evm := c.evmProvider(header, params.DeployerAddress, statedb)
	_, _, _, vmerr := evm.Replace(vm.AccountRef(params.DeployerAddress), bytecode, params.AutonityContractAddress)
	if vmerr != nil {
		log.Error("replaceAutonityBytecode evm.Create", "err", vmerr)
		return vmerr
	}
	return nil
}

type errorWithRevertReason struct {
	error
	reason string
}

func (e *errorWithRevertReason) Error() string {
	return e.error.Error() + ": " + e.reason
}

func newErrorWithRevertReason(err error, ret []byte) error {
	if err == nil || !errors.Is(err, vm.ErrExecutionReverted) {
		return err
	}
	reason, unpackErr := abi.UnpackRevert(ret)
	if unpackErr != nil {
		return err // if we cannot unpack the revert reason, return the original error unmodified
	}
	return &errorWithRevertReason{err, reason}
}

// AutonityContractCall calls the specified function of the autonity contract
// with the given args, and returns the output unpacked into the result
// interface.
// It returns the gas used for the call
func (c *AutonityContract) AutonityContractCall(statedb vm.StateDB, header *types.Header, function string, result any, args ...any) (uint64, error) {
	packedArgs, err := c.contractABI.Pack(function, args...)
	if err != nil {
		return 0, err
	}
	ret, usedGas, err := c.CallContractFunc(statedb, header, packedArgs)
	if err != nil {
		return usedGas, newErrorWithRevertReason(err, ret)
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&raw{}) {
		rawPtr := result.(*raw)
		*rawPtr = ret
		return usedGas, nil
	}
	if err := c.contractABI.UnpackIntoInterface(result, function, ret); err != nil {
		log.Error("Could not unpack returned value", "function", function)
		return usedGas, err
	}

	return usedGas, nil
}

func (c *AutonityContract) Mint(header *types.Header, statedb vm.StateDB, address common.Address, amount *big.Int) error {
	packedArgs, err := c.contractABI.Pack("mint", address, amount)
	if err != nil {
		return fmt.Errorf("error while generating call data for mint: %w", err)
	}

	_, err = c.CallContractFuncAs(statedb, header, c.EVMContract.chainConfig.AutonityContractConfig.Operator, packedArgs)
	if err != nil {
		return fmt.Errorf("error while calling mint: %w", err)
	}
	return nil
}

func (c *AutonityContract) Bond(header *types.Header, statedb vm.StateDB, from common.Address, validatorAddress common.Address, amount *big.Int) error {

	packedArgs, err := c.contractABI.Pack("bond", validatorAddress, amount)
	if err != nil {
		return fmt.Errorf("error while generating call data for bond: %w", err)
	}
	_, err = c.CallContractFuncAs(statedb, header, from, packedArgs)

	if err != nil {
		return fmt.Errorf("error while calling bond: %w", err)
	}
	return nil
}

func (c *AutonityContract) FinalizeInitialization(header *types.Header, statedb vm.StateDB, delta uint64) error {
	packedArgs, err := c.contractABI.Pack("finalizeInitialization", new(big.Int).SetUint64(delta))
	if err != nil {
		return fmt.Errorf("error while generating call data for finalizeInitialization: %w", err)
	}

	_, _, err = c.CallContractFunc(statedb, header, packedArgs)
	if err != nil {
		return fmt.Errorf("error while calling finalizeInitialization: %w", err)
	}

	return nil
}

func (c *AutonityContract) CreateSchedule(header *types.Header, statedb vm.StateDB, vault common.Address, schedule params.Schedule) error {
	packedArgs, err := c.contractABI.Pack("createSchedule", vault, schedule.Amount, schedule.Start, schedule.TotalDuration)
	if err != nil {
		return fmt.Errorf("error while generating call data for createSchedule: %w", err)
	}

	_, err = c.CallContractFuncAs(statedb, header, c.chainConfig.AutonityContractConfig.Operator, packedArgs)
	if err != nil {
		return fmt.Errorf("error while calling createSchedule: %w", err)
	}

	return nil
}

func (c *NonStakeableVestingContract) NewContract(header *types.Header, statedb vm.StateDB, contract params.NonStakeableVestingData) error {
	packedArgs, err := c.contractABI.Pack("newContract", contract.Beneficiary, contract.Amount, contract.ScheduleID, contract.CliffDuration)
	if err != nil {
		return fmt.Errorf("error while generating call data for newContract: %w", err)
	}

	ret, err := c.CallContractFuncAs(statedb, header, c.chainConfig.AutonityContractConfig.Operator, packedArgs)
	if err != nil {
		return fmt.Errorf("error while calling newContract: %w, returned %s", err, string(ret))
	}

	return nil
}

func (c *StakeableVestingManagerContract) NewContract(header *types.Header, statedb vm.StateDB, contract params.StakeableVestingData) error {
	packedArgs, err := c.contractABI.Pack("newContract", contract.Beneficiary, contract.Amount, contract.Start, contract.CliffDuration, contract.TotalDuration)
	if err != nil {
		return fmt.Errorf("error while generating call data for newContract: %w", err)
	}

	ret, err := c.CallContractFuncAs(statedb, header, c.chainConfig.AutonityContractConfig.Operator, packedArgs)
	if err != nil {
		return fmt.Errorf("error while calling newContract: %w, %s", err, string(ret))
	}

	return nil
}

// CallContractFunc creates an evm object, uses it to call the
// specified function of the autonity contract with packedArgs and returns the
// packed result. If there is an error making the evm call it will be returned.
// Callers should use the autonity contract ABI to pack and unpack the args and
// result.
// It returns the amount of gas used for the call
func (c *EVMContract) CallContractFunc(statedb vm.StateDB, header *types.Header, contractAddress common.Address, packedArgs []byte) ([]byte, uint64, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, params.DeployerAddress, statedb)
	packedResult, leftOverGas, err := evm.Call(vm.AccountRef(params.DeployerAddress), contractAddress, packedArgs, gas, new(big.Int))
	usedGas := gas - leftOverGas
	return packedResult, usedGas, err
}

func (c *EVMContract) CallContractFuncAs(statedb vm.StateDB, header *types.Header, contractAddress common.Address, origin common.Address, packedArgs []byte) ([]byte, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, origin, statedb)
	packedResult, _, err := evm.Call(vm.AccountRef(origin), contractAddress, packedArgs, gas, new(big.Int))
	return packedResult, err
}

func (c *AutonityContract) callGetCommitteeEnodes(state vm.StateDB, header *types.Header, asACN bool) (*types.Nodes, error) {
	var returnedEnodes []string
	_, err := c.AutonityContractCall(state, header, "getCommitteeEnodes", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes, asACN), nil
}

// callEpochByHeight get the epoch by height.
func (c *AutonityContract) callEpochByHeight(state vm.StateDB, header *types.Header, height *big.Int) (*types.EpochInfo, error) {
	var output raw
	if _, err := c.AutonityContractCall(state, header, "getEpochByHeight", &output, height); err != nil {
		return nil, err
	}

	data, err := c.contractABI.Unpack("getEpochByHeight", output)
	if err != nil {
		return nil, err
	}

	info := *abi.ConvertType(data[0], new(AutonityEpochInfo)).(*AutonityEpochInfo)

	committee := &types.Committee{}
	for _, member := range info.Committee {
		committee.Members = append(committee.Members, types.CommitteeMember{
			Address:           member.Addr,
			VotingPower:       member.VotingPower,
			ConsensusKeyBytes: member.ConsensusKey,
		})
	}
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error()) //nolint
	}

	epochInfo := &types.EpochInfo{
		Epoch: types.Epoch{
			Committee:          committee,
			PreviousEpochBlock: info.PreviousEpochBlock,
			NextEpochBlock:     info.NextEpochBlock,
			Delta:              info.Delta,
		},
		EpochBlock: info.EpochBlock,
	}

	return epochInfo, nil
}

// callGetEpochInfo gets the committee, the epoch boundaries and the omission failure delta based on the input header's state.
func (c *AutonityContract) callGetEpochInfo(state vm.StateDB, header *types.Header) (*types.EpochInfo, error) {
	var output raw
	if _, err := c.AutonityContractCall(state, header, "getEpochInfo", &output); err != nil {
		return nil, err
	}

	data, err := c.contractABI.Unpack("getEpochInfo", output)
	if err != nil {
		return nil, err
	}

	info := *abi.ConvertType(data[0], new(AutonityEpochInfo)).(*AutonityEpochInfo)

	committee := &types.Committee{}
	for _, member := range info.Committee {
		committee.Members = append(committee.Members, types.CommitteeMember{
			Address:           member.Addr,
			VotingPower:       member.VotingPower,
			ConsensusKeyBytes: member.ConsensusKey,
		})
	}
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error()) //nolint
	}

	epoch := &types.EpochInfo{
		EpochBlock: info.EpochBlock,
		Epoch: types.Epoch{
			Committee:          committee,
			NextEpochBlock:     info.NextEpochBlock,
			PreviousEpochBlock: info.PreviousEpochBlock,
			Delta:              info.Delta,
		},
	}
	return epoch, nil
}

func (c *AutonityContract) callGetCommittee(state vm.StateDB, header *types.Header) (*types.Committee, error) {
	var committeeMembers []types.CommitteeMember
	if _, err := c.AutonityContractCall(state, header, "getCommittee", &committeeMembers); err != nil {
		return nil, err
	}
	committee := &types.Committee{}
	committee.Members = committeeMembers
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error()) //nolint
	}
	return committee, nil
}

func (c *AutonityContract) callGetMinimumBaseFee(state vm.StateDB, header *types.Header) (*big.Int, error) {
	minBaseFee := new(big.Int)
	_, err := c.AutonityContractCall(state, header, "getMinimumBaseFee", &minBaseFee)
	if err != nil {
		return nil, err
	}
	return minBaseFee, nil
}

func (c *AutonityContract) callGetEpochPeriod(state vm.StateDB, header *types.Header) (*big.Int, error) {
	epochPeriod := new(big.Int)
	_, err := c.AutonityContractCall(state, header, "getEpochPeriod", &epochPeriod)
	if err != nil {
		return nil, err
	}
	return epochPeriod, nil
}

func recordFinalizeGasUsage(isEpochHeader bool, number uint64, usedGas int64) {
	if isEpochHeader {
		log.Debug("gas used to finalize epoch block", "number", number, "usedGas", usedGas)
		if metrics.Enabled {
			epochFinalizeGas.Add(usedGas)
		}
	} else {
		log.Debug("gas used to finalize block", "number", number, "usedGas", usedGas)
		if metrics.Enabled {
			finalizeGas.Add(usedGas)
		}
	}
}

func (c *AutonityContract) CallConfig(state vm.StateDB, header *types.Header) (*AutonityConfig, error) {
	var config AutonityConfig
	_, err := c.AutonityContractCall(state, header, "config", &config)
	return &config, err
}

func (c *AutonityContract) callEpochID(state vm.StateDB, header *types.Header) (*big.Int, error) {
	epochID := new(big.Int)
	_, err := c.AutonityContractCall(state, header, "epochID", &epochID)
	if err != nil {
		return nil, err
	}
	return epochID, nil
}

func (c *AutonityContract) callFinalize(state vm.StateDB, header *types.Header) (bool, *types.Epoch, error) {
	var updateReady bool
	var epochEnded bool
	var committeeMembers []types.CommitteeMember
	previousEpochBlock := new(big.Int)
	nextEpochBlock := new(big.Int)
	delta := new(big.Int)
	usedGas, err := c.AutonityContractCall(state, header, "finalize", &[]any{&updateReady, &epochEnded, &committeeMembers, &previousEpochBlock, &nextEpochBlock, &delta})
	recordFinalizeGasUsage(epochEnded, header.Number.Uint64(), int64(usedGas))
	if err != nil {
		return false, nil, err
	}

	if !epochEnded {
		return updateReady, nil, nil
	}

	// return with epoch info
	committee := &types.Committee{}
	committee.Members = committeeMembers
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error())
	}

	epoch := &types.Epoch{
		PreviousEpochBlock: previousEpochBlock,
		NextEpochBlock:     nextEpochBlock,
		Committee:          committee,
		Delta:              delta,
	}

	return updateReady, epoch, nil
}

func (c *AutonityContract) callRetrieveContract(state vm.StateDB, header *types.Header) ([]byte, string, error) {
	var bytecode []byte
	var abi string
	if _, err := c.AutonityContractCall(state, header, "getNewContract", &[]any{&bytecode, &abi}); err != nil {
		return nil, "", err
	}
	return bytecode, abi, nil
}
