package autonity

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"runtime"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

// This file should contain all the logic necessary for setting up state at genesis with nothing more and nothing less.

type (
	// GenesisBonds is an intermediary struct used to pass genesis delegations.
	// We cannot use autonity/core package here as it would cause import cycle
	GenesisBonds = []GenesisBond
	Delegation   = struct {
		Validator common.Address
		Amount    *big.Int
	}
	GenesisBond struct {
		Staker                common.Address
		UnbondedNewtonBalance *big.Int
		Delegations           []Delegation
	}
	genericDeployer func(abi *abi.ABI, bytecode []byte, value *big.Int, args ...interface{}) (common.Address, error)
	genericCaller   func(caller common.Address, contractAddress common.Address, abi *abi.ABI, method string, args ...interface{}) ([]byte, error)
	genesisStep     func(chainConfig *params.ChainConfig, genesisBonds GenesisBonds, deployer genericDeployer, caller genericCaller) error
)

var (
	genesisSequence = []genesisStep{
		deployAutonityContract,
		executeGenesisDelegations,
		deployAccountabilityContract,
		deployOracleContract,
		deployACUContract,
		deploySupplyControlContract,
		deployStabilizationContract,
		deployUpgradeManagerContract,
		deployInflationControllerContract,
		deployStakeableVestingContract,
		createDefaultStakingVestingContracts,
		deployNonStakeableVestingContract,
		createDefaultNonStakeableVestingContracts,
	}
	errBadDeploymentAddress = errors.New("mismatch with params deployment address")
)

func ExecuteGenesisSequence(genesisConfig *params.ChainConfig, genesisBonds GenesisBonds, evm *vm.EVM) error {

	contractDeployer := func(abi *abi.ABI, bytecode []byte, value *big.Int, args ...interface{}) (common.Address, error) {
		constructorParams, err := abi.Pack("", args...)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to pack parameters: %w", err)
		}
		if value.BitLen() != 0 {
			evm.StateDB.AddBalance(params.DeployerAddress, value)
		}
		data := append(bytecode, constructorParams...)
		gas := uint64(math.MaxUint64)
		_, addr, _, err := evm.Create(vm.AccountRef(params.DeployerAddress), data, gas, value)
		return addr, err
	}

	contractCaller := func(origin common.Address, contractAddress common.Address, abi *abi.ABI, method string, args ...interface{}) ([]byte, error) {
		packedArgs, err := abi.Pack(method, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to pack parameters: %w", err)
		}
		gas := uint64(math.MaxUint64)
		packedResult, _, err := evm.Call(vm.AccountRef(origin), contractAddress, packedArgs, gas, common.Big0)
		return packedResult, err
	}

	for i, fn := range genesisSequence {
		if err := fn(genesisConfig, genesisBonds, contractDeployer, contractCaller); err != nil {
			log.Error("Failed to execute genesis step", "i", i, "err", err, "fn", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
			return err
		}
	}
	return nil
}

func deployAutonityContract(config *params.ChainConfig, genesisBonds GenesisBonds, deploy genericDeployer, caller genericCaller) error {
	contractConfig := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:             new(big.Int).SetUint64(config.AutonityContractConfig.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(config.AutonityContractConfig.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(config.AutonityContractConfig.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(config.AutonityContractConfig.InitialInflationReserve),
			TreasuryAccount:         config.AutonityContractConfig.Treasury,
		},
		Contracts: AutonityContracts{
			AccountabilityContract:      params.AccountabilityContractAddress,
			OracleContract:              params.OracleContractAddress,
			AcuContract:                 params.ACUContractAddress,
			SupplyControlContract:       params.SupplyControlContractAddress,
			StabilizationContract:       params.StabilizationContractAddress,
			UpgradeManagerContract:      params.UpgradeManagerContractAddress,
			InflationControllerContract: params.InflationControllerContractAddress,
			NonStakableVestingContract:  params.NonStakeableVestingContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount: config.AutonityContractConfig.Operator,
			EpochPeriod:     new(big.Int).SetUint64(config.AutonityContractConfig.EpochPeriod),
			BlockPeriod:     new(big.Int).SetUint64(config.AutonityContractConfig.BlockPeriod),
			CommitteeSize:   new(big.Int).SetUint64(config.AutonityContractConfig.MaxCommitteeSize),
		},
		ContractVersion: big.NewInt(1),
	}
	validators := make([]params.Validator, 0, len(config.AutonityContractConfig.Validators))
	for _, v := range config.AutonityContractConfig.Validators {
		validators = append(validators, *v)
	}
	addr, err := deploy(&generated.AutonityAbi, generated.AutonityBytecode, common.Big0, validators, contractConfig)
	if err != nil {
		return fmt.Errorf("failed to deploy Autonity contract: %w", err)
	}
	if addr != params.AutonityContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func executeGenesisDelegations(config *params.ChainConfig, genesisBonds GenesisBonds, deploy genericDeployer, caller genericCaller) error {
	mint := func(address common.Address, amount *big.Int) error {
		_, err := caller(config.AutonityContractConfig.Operator, params.AutonityContractAddress, &generated.AutonityAbi, "mint", address, amount)
		return err
	}
	bond := func(staker, validator common.Address, amount *big.Int) error {
		_, err := caller(staker, params.AutonityContractAddress, &generated.AutonityAbi, "bond", validator, amount)
		return err
	}
	finalizeInitialization := func() error {
		_, err := caller(params.DeployerAddress, params.AutonityContractAddress, &generated.AutonityAbi, "finalizeInitialization")
		return err
	}
	for _, alloc := range genesisBonds {
		balanceToMint := new(big.Int)
		if alloc.UnbondedNewtonBalance != nil {
			balanceToMint.Add(balanceToMint, alloc.UnbondedNewtonBalance)
		}
		for _, delegation := range alloc.Delegations {
			balanceToMint.Add(balanceToMint, delegation.Amount)
		}
		if balanceToMint.Cmp(common.Big0) > 0 {
			if err := mint(alloc.Staker, balanceToMint); err != nil {
				return fmt.Errorf("error while minting Newton: %w", err)
			}
			for _, delegation := range alloc.Delegations {
				if err := bond(alloc.Staker, delegation.Validator, delegation.Amount); err != nil {
					return fmt.Errorf("error while bonding: %w", err)
				}
			}
		}
	}
	if err := finalizeInitialization(); err != nil {
		return fmt.Errorf("error while calling finalizeInitialization: %w", err)
	}
	return nil
}

func deployAccountabilityContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.AccountabilityConfig == nil {
		config.AccountabilityConfig = params.DefaultAccountabilityConfig
	}
	accountabilityConfig := AccountabilityConfig{
		InnocenceProofSubmissionWindow: new(big.Int).SetUint64(config.AccountabilityConfig.InnocenceProofSubmissionWindow),
		BaseSlashingRateLow:            new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateLow),
		BaseSlashingRateMid:            new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateMid),
		CollusionFactor:                new(big.Int).SetUint64(config.AccountabilityConfig.CollusionFactor),
		HistoryFactor:                  new(big.Int).SetUint64(config.AccountabilityConfig.HistoryFactor),
		JailFactor:                     new(big.Int).SetUint64(config.AccountabilityConfig.JailFactor),
		SlashingRatePrecision:          new(big.Int).SetUint64(config.AccountabilityConfig.SlashingRatePrecision),
	}
	addr, err := deploy(&generated.AccountabilityAbi, generated.AccountabilityBytecode, common.Big0, params.AutonityContractAddress, accountabilityConfig)
	if err != nil {
		return fmt.Errorf("failed to deploy accountability contract: %w", err)
	}
	if addr != params.AccountabilityContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deployOracleContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.OracleContractConfig == nil {
		log.Info("Using default genesis parameters for the Oracle Contract")
		config.OracleContractConfig = params.DefaultGenesisOracleConfig
	}
	if err := config.OracleContractConfig.SetDefaults(); err != nil {
		log.Crit("Error with Oracle Contract configuration", "err", err)
	}
	voters := make([]common.Address, len(config.AutonityContractConfig.Validators))
	for _, val := range config.AutonityContractConfig.Validators {
		voters = append(voters, val.OracleAddress)
	}

	addr, err := deploy(
		&generated.OracleAbi,
		generated.OracleBytecode,
		common.Big0,
		voters,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		config.OracleContractConfig.Symbols,
		new(big.Int).SetUint64(config.OracleContractConfig.VotePeriod),
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Oracle contract: %w", err)
	}
	if addr != params.OracleContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deployACUContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.ASM.ACUContractConfig == nil {
		log.Info("Config missing, using default parameters for the ACU contract")
		config.ASM.ACUContractConfig = params.DefaultAcuContractGenesis
	} else {
		config.ASM.ACUContractConfig.SetDefaults()
	}

	bigQuantities := make([]*big.Int, len(config.ASM.ACUContractConfig.Quantities))
	for i := range config.ASM.ACUContractConfig.Quantities {
		bigQuantities[i] = new(big.Int).SetUint64(config.ASM.ACUContractConfig.Quantities[i])
	}

	addr, err := deploy(
		&generated.ACUAbi,
		generated.ACUBytecode,
		common.Big0,
		config.ASM.ACUContractConfig.Symbols,
		bigQuantities,
		new(big.Int).SetUint64(config.ASM.ACUContractConfig.Scale),
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.OracleContractAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy ACU contract: %w", err)
	}
	if addr != params.ACUContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deploySupplyControlContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.ASM.SupplyControlConfig == nil {
		log.Info("Config missing, using default parameters for the Supply Control contract")
		config.ASM.SupplyControlConfig = params.DefaultSupplyControlGenesis
	} else {
		config.ASM.SupplyControlConfig.SetDefaults()
	}

	value := (*big.Int)(config.ASM.SupplyControlConfig.InitialAllocation)
	addr, err := deploy(
		&generated.SupplyControlAbi,
		generated.SupplyControlBytecode,
		value,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.StabilizationContractAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy SupplyControl contract: %w", err)
	}
	if addr != params.SupplyControlContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deployUpgradeManagerContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	addr, err := deploy(
		&generated.UpgradeManagerAbi,
		generated.UpgradeManagerBytecode,
		common.Big0,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Upgrade Manager contract: %w", err)
	}
	if addr != params.UpgradeManagerContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deployStabilizationContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.ASM.StabilizationContractConfig == nil {
		log.Info("Config missing, using default parameters for the Stabilization contract")
		config.ASM.StabilizationContractConfig = params.DefaultStabilizationGenesis
	} else {
		config.ASM.StabilizationContractConfig.SetDefaults()
	}

	stabilizationConfig := StabilizationConfig{
		BorrowInterestRate:        (*big.Int)(config.ASM.StabilizationContractConfig.BorrowInterestRate),
		LiquidationRatio:          (*big.Int)(config.ASM.StabilizationContractConfig.LiquidationRatio),
		MinCollateralizationRatio: (*big.Int)(config.ASM.StabilizationContractConfig.MinCollateralizationRatio),
		MinDebtRequirement:        (*big.Int)(config.ASM.StabilizationContractConfig.MinDebtRequirement),
		TargetPrice:               (*big.Int)(config.ASM.StabilizationContractConfig.TargetPrice),
	}

	addr, err := deploy(
		&generated.StabilizationAbi,
		generated.StabilizationBytecode,
		common.Big0,
		stabilizationConfig,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.OracleContractAddress,
		params.SupplyControlContractAddress,
		params.AutonityContractAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Stabilization contract: %w", err)
	}
	if addr != params.StabilizationContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deployInflationControllerContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.InflationContractConfig == nil {
		log.Info("Config missing, using default parameters for the Inflation Controller contract")
		config.InflationContractConfig = params.DefaultInflationControllerGenesis
	} else {
		config.InflationContractConfig.SetDefaults()
	}
	param := InflationControllerParams{
		InflationRateInitial:      (*big.Int)(config.InflationContractConfig.InflationRateInitial),
		InflationRateTransition:   (*big.Int)(config.InflationContractConfig.InflationRateTransition),
		InflationCurveConvexity:   (*big.Int)(config.InflationContractConfig.InflationCurveConvexity),
		InflationTransitionPeriod: (*big.Int)(config.InflationContractConfig.InflationTransitionPeriod),
		InflationReserveDecayRate: (*big.Int)(config.InflationContractConfig.InflationReserveDecayRate),
	}
	addr, err := deploy(&generated.InflationControllerAbi, generated.InflationControllerBytecode, common.Big0, param)
	if err != nil {
		return fmt.Errorf("failed to deploy inflation controller contract: %w", err)
	}
	if addr != params.InflationControllerContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func deployStakeableVestingContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.StakeableVestingConfig == nil {
		log.Info("Config missing, using default parameters for the Stakeable Vesting contract")
		config.StakeableVestingConfig = params.DefaultStakeableVestingGenesis
	} else {
		config.StakeableVestingConfig.SetDefaults()
	}
	addr, err := deploy(
		&generated.StakableVestingAbi,
		generated.StakableVestingBytecode,
		common.Big0,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy stakeable vesting contract: %w", err)
	}
	if addr != params.StakeableVestingContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func createDefaultStakingVestingContracts(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	if _, err := caller(
		config.AutonityContractConfig.Operator,
		params.AutonityContractAddress,
		&generated.AutonityAbi,
		"mint",
		params.StakeableVestingContractAddress,
		config.StakeableVestingConfig.TotalNominal,
	); err != nil {
		return fmt.Errorf("error while minting total nominal to stakeable vesting contract: %w", err)
	}

	if _, err := caller(
		config.AutonityContractConfig.Operator,
		params.StakeableVestingContractAddress,
		&generated.StakableVestingAbi,
		"setTotalNominal",
		config.StakeableVestingConfig.TotalNominal,
	); err != nil {
		return fmt.Errorf("error while setting total nominal in stakeable vesting contract: %w", err)
	}

	callNewStakeableContract := func(data params.StakeableVestingData) error {
		_, err := caller(
			config.AutonityContractConfig.Operator,
			params.StakeableVestingContractAddress,
			&generated.StakableVestingAbi,
			"newContract",
			data.Beneficiary,
			data.Amount,
			data.Start,
			data.CliffDuration,
			data.TotalDuration,
		)
		return err
	}

	for i, data := range config.StakeableVestingConfig.StakeableContracts {
		if err := callNewStakeableContract(data); err != nil {
			return fmt.Errorf("failed to create new stakeable vesting contract (i=%d): %w", i, err)
		}
	}
	return nil
}

func deployNonStakeableVestingContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.NonStakeableVestingConfig == nil {
		log.Info("Config missing, using default parameters for the Non-Stakable Vesting contract")
		config.NonStakeableVestingConfig = params.DefaultNonStakeableVestingGenesis
	} else {
		config.NonStakeableVestingConfig.SetDefaults()
	}
	addr, err := deploy(
		&generated.NonStakableVestingAbi,
		generated.NonStakableVestingBytecode,
		common.Big0,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy non-stakeable vesting contract: %w", err)
	}
	if addr != params.NonStakeableVestingContractAddress {
		return errBadDeploymentAddress
	}
	return nil
}

func createDefaultNonStakeableVestingContracts(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	if _, err := caller(
		config.AutonityContractConfig.Operator,
		params.NonStakeableVestingContractAddress,
		&generated.NonStakableVestingAbi,
		"setTotalNominal",
		config.NonStakeableVestingConfig.TotalNominal,
	); err != nil {
		return fmt.Errorf("error while seting total nominal in non-stakable vesting contract: %w", err)
	}

	if _, err := caller(
		config.AutonityContractConfig.Operator,
		params.NonStakeableVestingContractAddress,
		&generated.NonStakableVestingAbi,
		"setMaxAllowedDuration",
		config.NonStakeableVestingConfig.MaxAllowedDuration,
	); err != nil {
		return fmt.Errorf("error while seting max allowed duration in non-stakable vesting contract: %w", err)
	}

	createNonStakeableSchedule := func(data params.NonStakeableSchedule) error {
		// 	packedArgs, err := c.contractABI.Pack("createSchedule", schedule.Amount, schedule.Start, schedule.CliffDuration, schedule.TotalDuration)
		_, err := caller(
			config.AutonityContractConfig.Operator,
			params.NonStakeableVestingContractAddress,
			&generated.NonStakableVestingAbi,
			"createSchedule",
			data.Amount,
			data.Start,
			data.CliffDuration,
			data.TotalDuration,
		)
		return err
	}
	for _, schedule := range config.NonStakeableVestingConfig.NonStakeableSchedules {
		if err := createNonStakeableSchedule(schedule); err != nil {
			return fmt.Errorf("error while creating new non-stakable schedule: %w", err)
		}
	}

	newNonStakeableContract := func(data params.NonStakeableVestingData) error {
		//	packedArgs, err := c.contractABI.Pack("newContract", contract.Beneficiary, contract.Amount, contract.ScheduleID)
		_, err := caller(
			config.AutonityContractConfig.Operator,
			params.NonStakeableVestingContractAddress,
			&generated.NonStakableVestingAbi,
			"newContract",
			data.Beneficiary,
			data.Amount,
			data.ScheduleID,
		)
		return err
	}
	for _, data := range config.NonStakeableVestingConfig.NonStakeableContracts {
		if err := newNonStakeableContract(data); err != nil {
			return fmt.Errorf("failed to create new non-stakable vesting contract: %w", err)
		}
	}
	return nil
}
