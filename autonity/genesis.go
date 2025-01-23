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
		Staker        common.Address
		NewtonBalance *big.Int
		Bonds         []Delegation
	}
	genericDeployer func(address common.Address, abi *abi.ABI, bytecode []byte, value *big.Int, args ...interface{}) error
	genericCaller   func(caller common.Address, contractAddress common.Address, abi *abi.ABI, method string, args ...interface{}) ([]byte, error)
	genesisStep     func(chainConfig *params.ChainConfig, genesisBonds GenesisBonds, deployer genericDeployer, caller genericCaller) error
)

var (
	commonSequence = []genesisStep{
		executeGenesisDelegations,
		createAutonitySchedules,
		finalizeAutonityInitialization,
		deployAccountabilityContract,
		deployOracleContract,
		deployACUContract,
		deploySupplyControlContract,
		deployStabilizationContract,
		deployUpgradeManagerContract,
		deployInflationControllerContract,
		deployStakableVestingManagerContract,
		createDefaultStakableVestingContracts,
		deployNonStakableVestingContract,
		createDefaultNonStakableVestingContracts,
		deployOmissionAccountabilityContract,
		deployAuctioneerContract,
	}
	genesisSequence = append(
		[]genesisStep{
			deployAutonityContract,
		},
		commonSequence...,
	)
	testGenesisSequence = append(
		[]genesisStep{
			deployAutonityTestContract,
		},
		commonSequence...,
	)
	errBadDeploymentAddress = errors.New("mismatch with params deployment address")
)

// *
// Genesis sequence execution
// *

func ExecuteGenesisSequence(genesisConfig *params.ChainConfig, genesisBonds GenesisBonds, evm *vm.EVM) error {
	return executeGenesisSequence(genesisConfig, genesisBonds, evm, genesisSequence)
}

func ExecuteTestGenesisSequence(genesisConfig *params.ChainConfig, genesisBonds GenesisBonds, evm *vm.EVM) error {
	return executeGenesisSequence(genesisConfig, genesisBonds, evm, testGenesisSequence)
}

func executeGenesisSequence(genesisConfig *params.ChainConfig, genesisBonds GenesisBonds, evm *vm.EVM, genesisSeq []genesisStep) error {
	contractDeployer := func(
		address common.Address,
		abi *abi.ABI,
		bytecode []byte,
		value *big.Int,
		args ...interface{},
	) error {
		constructorParams, err := abi.Pack("", args...)
		if err != nil {
			return fmt.Errorf("failed to pack parameters: %w", err)
		}
		if value.BitLen() != 0 && evm.StateDB.GetBalance(params.DeployerAddress).Cmp(value) < 0 {
			evm.StateDB.AddBalance(params.DeployerAddress, value)
		}
		data := append(bytecode, constructorParams...)
		gas := uint64(math.MaxUint64)
		_, addr, _, err := evm.Create(vm.AccountRef(params.DeployerAddress), data, gas, value)
		if err != nil {
			return err
		}
		if addr != address {
			return errBadDeploymentAddress
		}
		return nil
	}

	contractCaller := func(
		origin common.Address,
		contractAddress common.Address,
		abi *abi.ABI,
		method string,
		args ...interface{},
	) ([]byte, error) {
		packedArgs, err := abi.Pack(method, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to pack parameters for method: %s %w", method, err)
		}
		gas := uint64(math.MaxUint64)
		packedResult, _, err := evm.Call(vm.AccountRef(origin), contractAddress, packedArgs, gas, common.Big0)
		return packedResult, err
	}

	for i, fn := range genesisSeq {
		if err := fn(genesisConfig, genesisBonds, contractDeployer, contractCaller); err != nil {
			log.Error(
				"Failed to execute genesis step", "i", i, "err", err, "fn", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
			)
			return err
		}
	}
	return nil
}

// *
// Main protocol steps
// *

func deployAutonityContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	contractConfig := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:             new(big.Int).SetUint64(config.AutonityContractConfig.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(config.AutonityContractConfig.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(config.AutonityContractConfig.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(config.AutonityContractConfig.InitialInflationReserve),
			WithholdingThreshold:    new(big.Int).SetUint64(config.AutonityContractConfig.WithholdingThreshold),
			ProposerRewardRate:      new(big.Int).SetUint64(config.AutonityContractConfig.ProposerRewardRate),
			OracleRewardRate:        new(big.Int).SetUint64(config.AutonityContractConfig.OracleRewardRate),
			WithheldRewardsPool:     config.AutonityContractConfig.WithheldRewardsPool,
			TreasuryAccount:         config.AutonityContractConfig.Treasury,
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
			AuctioneerContract:             params.AuctioneerContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount:     config.AutonityContractConfig.Operator,
			EpochPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.EpochPeriod),
			BlockPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.BlockPeriod),
			CommitteeSize:       new(big.Int).SetUint64(config.AutonityContractConfig.MaxCommitteeSize),
			MaxScheduleDuration: new(big.Int).SetUint64(config.AutonityContractConfig.MaxScheduleDuration),
		},
		ContractVersion: big.NewInt(1),
	}
	validators := make([]params.Validator, 0, len(config.AutonityContractConfig.Validators))
	for _, v := range config.AutonityContractConfig.Validators {
		validators = append(validators, *v)
	}
	err := deploy(
		params.AutonityContractAddress,
		&generated.AutonityAbi,
		generated.AutonityBytecode,
		common.Big0,
		validators,
		contractConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Autonity contract: %w", err)
	}
	return nil
}

func executeGenesisDelegations(config *params.ChainConfig, genesisBonds GenesisBonds, _ genericDeployer, caller genericCaller) error {
	mint := func(address common.Address, amount *big.Int) error {
		ret, err := caller(config.AutonityContractConfig.Operator, params.AutonityContractAddress, &generated.AutonityAbi, "mint", address, amount)
		return newErrorWithRevertReason(err, ret)
	}
	bond := func(staker, validator common.Address, amount *big.Int) error {
		ret, err := caller(staker, params.AutonityContractAddress, &generated.AutonityAbi, "bond", validator, amount)
		return newErrorWithRevertReason(err, ret)
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
			if err := mint(alloc.Staker, balanceToMint); err != nil {
				return fmt.Errorf("error while minting Newton: %w", err)
			}
			for _, delegation := range alloc.Bonds {
				if err := bond(alloc.Staker, delegation.Validator, delegation.Amount); err != nil {
					return fmt.Errorf("error while bonding: %w", err)
				}
			}
		}
	}
	return nil
}

func createAutonitySchedules(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	createSchedule := func(schedule params.Schedule) error {
		if schedule.VaultAddress != params.NonStakeableVestingContractAddress {
			return fmt.Errorf(
				"invalid Schedule configuration, should match non stakable vesting contract address: %s",
				schedule.VaultAddress,
			)
		}
		ret, err := caller(
			config.AutonityContractConfig.Operator,
			params.AutonityContractAddress,
			&generated.AutonityAbi,
			"createSchedule",
			schedule.VaultAddress,
			schedule.Amount,
			schedule.Start,
			schedule.TotalDuration,
		)
		return newErrorWithRevertReason(err, ret)
	}
	for _, schedule := range config.AutonityContractConfig.Schedules {
		if err := createSchedule(schedule); err != nil {
			return fmt.Errorf("error while creating schedule: %w", err)
		}
	}
	return nil
}

func finalizeAutonityInitialization(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	ret, err := caller(
		params.DeployerAddress,
		params.AutonityContractAddress,
		&generated.AutonityAbi,
		"finalizeInitialization",
		new(big.Int).SetUint64(config.OmissionAccountabilityConfig.Delta),
	)
	return newErrorWithRevertReason(err, ret)
}

func deployAccountabilityContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.AccountabilityConfig == nil {
		config.AccountabilityConfig = params.DefaultAccountabilityConfig
	}
	accountabilityConfig := AccountabilityConfig{
		InnocenceProofSubmissionWindow: new(big.Int).SetUint64(config.AccountabilityConfig.InnocenceProofSubmissionWindow),
		BaseSlashingRates: AccountabilityBaseSlashingRates{
			Low:  new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateLow),
			Mid:  new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateMid),
			High: new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateHigh),
		},
		Factors: AccountabilityFactors{
			Collusion: new(big.Int).SetUint64(config.AccountabilityConfig.CollusionFactor),
			History:   new(big.Int).SetUint64(config.AccountabilityConfig.HistoryFactor),
			Jail:      new(big.Int).SetUint64(config.AccountabilityConfig.JailFactor),
		},
	}
	err := deploy(
		params.AccountabilityContractAddress,
		&generated.AccountabilityAbi,
		generated.AccountabilityBytecode,
		common.Big0,
		params.AutonityContractAddress,
		accountabilityConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy accountability contract: %w", err)
	}
	return nil
}

func deployOmissionAccountabilityContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	omissionConfig := config.OmissionAccountabilityConfig

	conf := OmissionAccountabilityConfig{
		InactivityThreshold:    new(big.Int).SetUint64(omissionConfig.InactivityThreshold),
		LookbackWindow:         new(big.Int).SetUint64(omissionConfig.LookbackWindow),
		PastPerformanceWeight:  new(big.Int).SetUint64(omissionConfig.PastPerformanceWeight),
		InitialJailingPeriod:   new(big.Int).SetUint64(omissionConfig.InitialJailingPeriod),
		InitialProbationPeriod: new(big.Int).SetUint64(omissionConfig.InitialProbationPeriod),
		InitialSlashingRate:    new(big.Int).SetUint64(omissionConfig.InitialSlashingRate),
		Delta:                  new(big.Int).SetUint64(omissionConfig.Delta),
	}

	treasuries := make([]common.Address, len(config.AutonityContractConfig.Validators))
	for i, val := range config.AutonityContractConfig.Validators {
		treasuries[i] = val.Treasury
	}
	err := deploy(
		params.OmissionAccountabilityContractAddress,
		&generated.OmissionAccountabilityAbi,
		generated.OmissionAccountabilityBytecode,
		common.Big0,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		treasuries,
		conf,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy omission accountability contract: %w", err)
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
	treasuries := make([]common.Address, len(config.AutonityContractConfig.Validators))
	validators := make([]common.Address, len(config.AutonityContractConfig.Validators))
	for i, val := range config.AutonityContractConfig.Validators {
		voters[i] = val.OracleAddress
		treasuries[i] = val.Treasury
		validators[i] = *val.NodeAddress
	}

	oracleConfig := OracleConfig{
		Autonity:                  params.AutonityContractAddress,
		Operator:                  config.AutonityContractConfig.Operator,
		VotePeriod:                new(big.Int).SetUint64(config.OracleContractConfig.VotePeriod),
		OutlierDetectionThreshold: new(big.Int).SetUint64(config.OracleContractConfig.OutlierDetectionThreshold),
		OutlierSlashingThreshold:  new(big.Int).SetUint64(config.OracleContractConfig.OutlierSlashingThreshold),
		BaseSlashingRate:          new(big.Int).SetUint64(config.OracleContractConfig.BaseSlashingRate),
	}

	err := deploy(
		params.OracleContractAddress,
		&generated.OracleAbi,
		generated.OracleBytecode,
		common.Big0,
		voters,
		validators,
		treasuries,
		config.OracleContractConfig.Symbols,
		oracleConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Oracle contract: %w", err)
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

	err := deploy(
		params.ACUContractAddress,
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
	err := deploy(
		params.SupplyControlContractAddress,
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
	return nil
}

func deployUpgradeManagerContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	err := deploy(
		params.UpgradeManagerContractAddress,
		&generated.UpgradeManagerAbi,
		generated.UpgradeManagerBytecode,
		common.Big0,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Upgrade Manager contract: %w", err)
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

	stabilizationConfig := IStabilizationConfig{
		BorrowInterestRate:        (*big.Int)(config.ASM.StabilizationContractConfig.BorrowInterestRate),
		LiquidationRatio:          (*big.Int)(config.ASM.StabilizationContractConfig.LiquidationRatio),
		MinCollateralizationRatio: (*big.Int)(config.ASM.StabilizationContractConfig.MinCollateralizationRatio),
		MinDebtRequirement:        (*big.Int)(config.ASM.StabilizationContractConfig.MinDebtRequirement),
		TargetPrice:               (*big.Int)(config.ASM.StabilizationContractConfig.TargetPrice),
	}

	err := deploy(
		params.StabilizationContractAddress,
		&generated.StabilizationAbi,
		generated.StabilizationBytecode,
		common.Big0,
		stabilizationConfig,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		params.OracleContractAddress,
		params.SupplyControlContractAddress,
		params.AuctioneerContractAddress,
		params.ACUContractAddress,
		params.AutonityContractAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Stabilization contract: %w", err)
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
	err := deploy(
		params.InflationControllerContractAddress,
		&generated.InflationControllerAbi,
		generated.InflationControllerBytecode,
		common.Big0,
		param,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy inflation controller contract: %w", err)
	}
	return nil
}

func deployStakableVestingManagerContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.StakeableVestingConfig == nil {
		log.Info("Config missing, using default parameters for the Stakable Vesting contract")
		config.StakeableVestingConfig = params.DefaultStakeableVestingGenesis
	} else {
		config.StakeableVestingConfig.SetDefaults()
	}
	err := deploy(
		params.StakeableVestingManagerContractAddress,
		&generated.StakeableVestingManagerAbi,
		generated.StakeableVestingManagerBytecode,
		common.Big0,
		params.AutonityContractAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Stakable vesting contract: %w", err)
	}
	return nil
}

func createDefaultStakableVestingContracts(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	if ret, err := caller(
		config.AutonityContractConfig.Operator,
		params.AutonityContractAddress,
		&generated.AutonityAbi,
		"mint",
		params.StakeableVestingManagerContractAddress,
		config.StakeableVestingConfig.TotalNominal,
	); err != nil {
		return fmt.Errorf(
			"error while minting total nominal to Stakable vesting contract: %w",
			newErrorWithRevertReason(err, ret),
		)
	}

	callNewStakableContract := func(contract params.StakeableVestingData) error {
		ret, err := caller(
			config.AutonityContractConfig.Operator,
			params.StakeableVestingManagerContractAddress,
			&generated.StakeableVestingManagerAbi,
			"newContract",
			contract.Beneficiary,
			contract.Amount,
			contract.Start,
			contract.CliffDuration,
			contract.TotalDuration,
		)
		return newErrorWithRevertReason(err, ret)
	}

	for i, data := range config.StakeableVestingConfig.StakeableContracts {
		if err := callNewStakableContract(data); err != nil {
			return fmt.Errorf("failed to create new Stakable vesting contract (i=%d): %w", i, err)
		}
	}
	return nil
}

func deployNonStakableVestingContract(_ *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	err := deploy(
		params.NonStakeableVestingContractAddress,
		&generated.NonStakeableVestingAbi,
		generated.NonStakeableVestingBytecode,
		common.Big0,
		params.AutonityContractAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy non-Stakable vesting contract: %w", err)
	}
	return nil
}

func createDefaultNonStakableVestingContracts(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	createNonStakableVestingContract := func(contract params.NonStakeableVestingData) error {
		ret, err := caller(
			config.AutonityContractConfig.Operator,
			params.NonStakeableVestingContractAddress,
			&generated.NonStakeableVestingAbi,
			"newContract",
			contract.Beneficiary,
			contract.Amount,
			contract.ScheduleID,
			contract.CliffDuration,
		)
		return newErrorWithRevertReason(err, ret)
	}
	for _, schedule := range config.NonStakeableVestingConfig.NonStakeableContracts {
		if err := createNonStakableVestingContract(schedule); err != nil {
			return fmt.Errorf("error while creating new non-stakable schedule: %w", err)
		}
	}
	return nil
}

func deployAuctioneerContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	if config.ASM.AuctioneerContractConfig == nil {
		log.Info("Config missing, using default parameters for the Auctioneer contract")
		config.ASM.AuctioneerContractConfig = params.DefaultAuctioneerGenesis
	} else {
		config.ASM.AuctioneerContractConfig.SetDefaults()
	}
	auctioneerConfig := AuctioneerConfig{
		LiquidationAuctionDuration: config.ASM.AuctioneerContractConfig.LiquidationAuctionDuration,
		InterestAuctionDuration:    config.ASM.AuctioneerContractConfig.InterestAuctionDuration,
		InterestAuctionDiscount:    config.ASM.AuctioneerContractConfig.InterestAuctionDiscount,
		InterestAuctionThreshold:   config.ASM.AuctioneerContractConfig.InterestAuctionThreshold,
	}
	err := deploy(
		params.AuctioneerContractAddress,
		&generated.AuctioneerAbi,
		generated.AuctioneerBytecode,
		common.Big0,
		auctioneerConfig,
		params.StabilizationContractAddress,
		params.OracleContractAddress,
		params.AutonityContractAddress,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy Auctioneer contract: %w", err)
	}

	return nil
}

// *
// Test only functions
// *

func deployAutonityTestContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	contractConfig := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:             new(big.Int).SetUint64(config.AutonityContractConfig.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(config.AutonityContractConfig.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(config.AutonityContractConfig.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(config.AutonityContractConfig.InitialInflationReserve),
			WithholdingThreshold:    new(big.Int).SetUint64(config.AutonityContractConfig.WithholdingThreshold),
			ProposerRewardRate:      new(big.Int).SetUint64(config.AutonityContractConfig.ProposerRewardRate),
			OracleRewardRate:        new(big.Int).SetUint64(config.AutonityContractConfig.OracleRewardRate),
			WithheldRewardsPool:     config.AutonityContractConfig.WithheldRewardsPool,
			TreasuryAccount:         config.AutonityContractConfig.Treasury,
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
			AuctioneerContract:             params.AuctioneerContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount:     config.AutonityContractConfig.Operator,
			EpochPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.EpochPeriod),
			BlockPeriod:         new(big.Int).SetUint64(config.AutonityContractConfig.BlockPeriod),
			CommitteeSize:       new(big.Int).SetUint64(config.AutonityContractConfig.MaxCommitteeSize),
			MaxScheduleDuration: new(big.Int).SetUint64(config.AutonityContractConfig.MaxScheduleDuration),
		},
		ContractVersion: big.NewInt(1),
	}
	validators := make([]params.Validator, 0, len(config.AutonityContractConfig.Validators))
	for _, v := range config.AutonityContractConfig.Validators {
		validators = append(validators, *v)
	}
	err := deploy(
		params.AutonityContractAddress,
		&generated.AutonityTestAbi,
		generated.AutonityTestBytecode,
		common.Big0,
		validators,
		contractConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy AutonityTest contract: %w", err)
	}
	return nil
}
