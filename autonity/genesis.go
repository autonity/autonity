package autonity

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"runtime"

	"github.com/autonity/autonity/autonity/bindings"

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
		verifyGenesisSequence,
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

func toContractConfig(acg *params.AutonityContractGenesis) bindings.IAutonityConfig {
	return bindings.IAutonityConfig{
		Policy: bindings.IAutonityPolicy{
			TreasuryFee:              new(big.Int).SetUint64(acg.TreasuryFee),
			MinBaseFee:               new(big.Int).SetUint64(acg.MinBaseFee),
			DelegationRate:           new(big.Int).SetUint64(acg.DelegationRate),
			UnbondingPeriod:          new(big.Int).SetUint64(acg.UnbondingPeriod),
			InitialInflationReserve:  (*big.Int)(acg.InitialInflationReserve),
			WithholdingThreshold:     new(big.Int).SetUint64(acg.WithholdingThreshold),
			ProposerRewardRate:       new(big.Int).SetUint64(acg.ProposerRewardRate),
			OracleRewardRate:         new(big.Int).SetUint64(acg.OracleRewardRate),
			WithheldRewardsPool:      acg.WithheldRewardsPool,
			TreasuryAccount:          acg.Treasury,
			BaseFeeChangeDenominator: new(big.Int).SetUint64(acg.BaseFeeChangeDenominator),
			ElasticityMultiplier:     new(big.Int).SetUint64(acg.ElasticityMultiplier),
		},
		Contracts: bindings.IAutonityContracts{
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
		Protocol: bindings.IAutonityProtocol{
			OperatorAccount:      acg.Operator,
			EpochPeriod:          new(big.Int).SetUint64(acg.EpochPeriod),
			BlockPeriod:          new(big.Int).SetUint64(acg.BlockPeriod),
			CommitteeSize:        new(big.Int).SetUint64(acg.MaxCommitteeSize),
			MaxScheduleDuration:  new(big.Int).SetUint64(acg.MaxScheduleDuration),
			GasLimit:             new(big.Int).SetUint64(acg.GasLimit),
			GasLimitBoundDivisor: new(big.Int).SetUint64(acg.GasLimitBoundDivisor),
		},
		ContractVersion: big.NewInt(1),
	}
}

func deployAutonityContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
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
		toContractConfig(config.AutonityContractConfig),
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
	accountabilityConfig := bindings.IAccountabilityConfig{
		InnocenceProofSubmissionWindow: new(big.Int).SetUint64(config.AccountabilityConfig.InnocenceProofSubmissionWindow),
		Delta:                          new(big.Int).SetUint64(config.AccountabilityConfig.Delta),
		Range:                          new(big.Int).SetUint64(config.AccountabilityConfig.Range),
		BaseSlashingRates: bindings.IAccountabilityBaseSlashingRates{
			Low:  new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateLow),
			Mid:  new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateMid),
			High: new(big.Int).SetUint64(config.AccountabilityConfig.BaseSlashingRateHigh),
		},
		Factors: bindings.IAccountabilityFactors{
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

	conf := bindings.OmissionAccountabilityConfig{
		InactivityThreshold:    new(big.Int).SetUint64(omissionConfig.InactivityThreshold),
		LookbackWindow:         new(big.Int).SetUint64(omissionConfig.LookbackWindow),
		PastPerformanceWeight:  new(big.Int).SetUint64(omissionConfig.PastPerformanceWeight),
		InitialJailingPeriod:   new(big.Int).SetUint64(omissionConfig.InitialJailingPeriod),
		InitialProbationPeriod: new(big.Int).SetUint64(omissionConfig.InitialProbationPeriod),
		InitialSlashingRate:    new(big.Int).SetUint64(omissionConfig.InitialSlashingRate),
		Delta:                  new(big.Int).SetUint64(omissionConfig.Delta),
	}

	err := deploy(
		params.OmissionAccountabilityContractAddress,
		&generated.OmissionAccountabilityAbi,
		generated.OmissionAccountabilityBytecode,
		common.Big0,
		params.AutonityContractAddress,
		config.AutonityContractConfig.Operator,
		conf,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy omission accountability contract: %w", err)
	}
	return nil
}

func deployOracleContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
	voters := make([]common.Address, len(config.AutonityContractConfig.Validators))
	treasuries := make([]common.Address, len(config.AutonityContractConfig.Validators))
	validators := make([]common.Address, len(config.AutonityContractConfig.Validators))
	for i, val := range config.AutonityContractConfig.Validators {
		voters[i] = val.OracleAddress
		treasuries[i] = val.Treasury
		validators[i] = *val.NodeAddress
	}

	oracleConfig := bindings.OracleConfig{
		Autonity:                  params.AutonityContractAddress,
		Operator:                  config.AutonityContractConfig.Operator,
		VotePeriod:                new(big.Int).SetUint64(config.OracleContractConfig.VotePeriod),
		OutlierDetectionThreshold: new(big.Int).SetUint64(config.OracleContractConfig.OutlierDetectionThreshold),
		OutlierSlashingThreshold:  new(big.Int).SetUint64(config.OracleContractConfig.OutlierSlashingThreshold),
		BaseSlashingRate:          new(big.Int).SetUint64(config.OracleContractConfig.BaseSlashingRate),
		NonRevealThreshold:        new(big.Int).SetUint64(config.OracleContractConfig.NonRevealThreshold),
		RevealResetInterval:       new(big.Int).SetUint64(config.OracleContractConfig.RevealResetInterval),
		SlashingRateCap:           new(big.Int).SetUint64(config.OracleContractConfig.SlashingRateCap),
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
	stabilizationConfig := bindings.IStabilizationConfig{
		BorrowInterestRate:        (*big.Int)(config.ASM.StabilizationContractConfig.BorrowInterestRate),
		AnnouncementWindow:        (*big.Int)(config.ASM.StabilizationContractConfig.AnnouncementWindow),
		LiquidationRatio:          (*big.Int)(config.ASM.StabilizationContractConfig.LiquidationRatio),
		MinCollateralizationRatio: (*big.Int)(config.ASM.StabilizationContractConfig.MinCollateralizationRatio),
		MinDebtRequirement:        (*big.Int)(config.ASM.StabilizationContractConfig.MinDebtRequirement),
		TargetPrice:               (*big.Int)(config.ASM.StabilizationContractConfig.TargetPrice),
		DefaultNTNATNPrice:        (*big.Int)(config.ASM.StabilizationContractConfig.DefaultNTNATNPrice),
		DefaultNTNUSDPrice:        (*big.Int)(config.ASM.StabilizationContractConfig.DefaultNTNUSDPrice),
		DefaultACUUSDPrice:        (*big.Int)(config.ASM.StabilizationContractConfig.DefaultACUUSDPrice),
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
	param := bindings.InflationControllerParams{
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

func deployStakableVestingManagerContract(_ *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
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
	auctioneerConfig := bindings.AuctioneerConfig{
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

func verifyGenesisSequence(config *params.ChainConfig, _ GenesisBonds, _ genericDeployer, caller genericCaller) error {
	if !config.AutonityContractConfig.VerifyGenesisSequence {
		return nil
	}
	// verify total allocations
	ret, err := caller(
		common.Address{},
		params.AutonityContractAddress,
		&generated.AutonityAbi,
		"totalSupply",
	)
	if err != nil {
		return fmt.Errorf("error while calling totalSupply: %w", newErrorWithRevertReason(err, ret))
	}

	data, err := generated.AutonityAbi.Unpack("totalSupply", ret)
	if err != nil {
		return fmt.Errorf("error while unpacking totalSupply: %w", err)
	}

	totalSupply := abi.ConvertType(data[0], new(big.Int)).(*big.Int)
	if totalSupply.Cmp((*big.Int)(config.AutonityContractConfig.TokenMint)) != 0 {
		return fmt.Errorf("genesis token allocation mismatch")
	}

	// verify bonded stake
	validatorBondedStake := func(addr common.Address) (*big.Int, error) {
		ret, err := caller(
			common.Address{},
			params.AutonityContractAddress,
			&generated.AutonityAbi,
			"getValidator",
			addr,
		)
		if err != nil {
			return nil, newErrorWithRevertReason(err, ret)
		}

		data, err := generated.AutonityAbi.Unpack("getValidator", ret)
		if err != nil {
			return nil, fmt.Errorf("error while unpacking getValidator: %w", err)
		}

		validator := abi.ConvertType(data[0], new(bindings.IAutonityValidator)).(*bindings.IAutonityValidator)
		return validator.BondedStake, nil
	}

	totalBondedStake := new(big.Int)
	for _, v := range config.AutonityContractConfig.Validators {
		stake, err := validatorBondedStake(*v.NodeAddress)
		if err != nil {
			return fmt.Errorf("error while calling getValidator: %w", err)
		}
		totalBondedStake.Add(totalBondedStake, stake)
	}
	if totalBondedStake.Cmp((*big.Int)(config.AutonityContractConfig.TokenBond)) != 0 {
		return fmt.Errorf("genesis total staking mismatch")
	}

	return nil
}

// *
// Test only functions
// *

func deployAutonityTestContract(config *params.ChainConfig, _ GenesisBonds, deploy genericDeployer, _ genericCaller) error {
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
		toContractConfig(config.AutonityContractConfig),
	)
	if err != nil {
		return fmt.Errorf("failed to deploy AutonityTest contract: %w", err)
	}
	return nil
}
