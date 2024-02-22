package autonity

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

type raw []byte

// GenesisBonds is an intermediary struct used to pass genesis bondings.
// We cannot use autonity/core package here as it would cause import cycle
type GenesisBonds = map[common.Address]GenesisBond

type GenesisBond struct {
	NewtonBalance *big.Int
	Bonds         map[common.Address]*big.Int
}

func DeployContracts(genesisConfig *params.ChainConfig, genesisBonds GenesisBonds, evmContracts *GenesisEVMContracts) error {
	if err := DeployAutonityContract(genesisConfig.AutonityContractConfig, genesisBonds, evmContracts); err != nil {
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
	if config.ASM.SupplyControlConfig == nil {
		log.Info("Config missing, using default parameters for the Supply Control contract")
		config.ASM.SupplyControlConfig = params.DefaultSupplyControlGenesis
	} else {
		config.ASM.SupplyControlConfig.SetDefaults()
	}

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

func DeployACUContract(config *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
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
	if config == nil {
		config = params.DefaultAccountabilityConfig
	}
	accountabilityConfig := AccountabilityConfig{
		InnocenceProofSubmissionWindow: new(big.Int).SetUint64(config.InnocenceProofSubmissionWindow),
		BaseSlashingRateLow:            new(big.Int).SetUint64(config.BaseSlashingRateLow),
		BaseSlashingRateMid:            new(big.Int).SetUint64(config.BaseSlashingRateMid),
		CollusionFactor:                new(big.Int).SetUint64(config.CollusionFactor),
		HistoryFactor:                  new(big.Int).SetUint64(config.HistoryFactor),
		JailFactor:                     new(big.Int).SetUint64(config.JailFactor),
		SlashingRatePrecision:          new(big.Int).SetUint64(config.SlashingRatePrecision),
	}
	err := evmContracts.DeployAccountabilityContract(params.AutonityContractAddress, accountabilityConfig, generated.AccountabilityBytecode)
	if err != nil {
		return fmt.Errorf("failed to deploy accountability contract: %w", err)
	}

	log.Info("Deployed Accountability contract", "address", params.AccountabilityContractAddress)

	return nil
}

func DeployAutonityContract(genesisConfig *params.AutonityContractGenesis, genesisBonds GenesisBonds, evmContracts *GenesisEVMContracts) error {
	contractConfig := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:     new(big.Int).SetUint64(genesisConfig.TreasuryFee),
			MinBaseFee:      new(big.Int).SetUint64(genesisConfig.MinBaseFee),
			DelegationRate:  new(big.Int).SetUint64(genesisConfig.DelegationRate),
			UnbondingPeriod: new(big.Int).SetUint64(genesisConfig.UnbondingPeriod),
			TreasuryAccount: genesisConfig.Treasury,
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
			OperatorAccount: genesisConfig.Operator,
			EpochPeriod:     new(big.Int).SetUint64(genesisConfig.EpochPeriod),
			BlockPeriod:     new(big.Int).SetUint64(genesisConfig.BlockPeriod),
			CommitteeSize:   new(big.Int).SetUint64(genesisConfig.MaxCommitteeSize),
		},
		ContractVersion: big.NewInt(1),
	}

	validators := make([]params.Validator, 0, len(genesisConfig.Validators))
	for _, v := range genesisConfig.Validators {
		validators = append(validators, *v)
	}

	err := evmContracts.DeployAutonityContract(genesisConfig.Bytecode, validators, contractConfig)
	if err != nil {
		log.Error("DeployAutonityContract failed", "err", err)
		return fmt.Errorf("failed to deploy Autonity contract: %w", err)
	}

	for addr, account := range genesisBonds {

		balanceToMint := new(big.Int)

		if account.NewtonBalance != nil {
			balanceToMint.Add(balanceToMint, account.NewtonBalance)
		}

		for _, amount := range account.Bonds {
			balanceToMint.Add(balanceToMint, amount)
		}

		if balanceToMint.Cmp(common.Big0) > 0 {
			err := evmContracts.Mint(addr, balanceToMint)
			if err != nil {
				return fmt.Errorf("error while minting Newton: %w", err)
			}

			for validatorAddress, amount := range account.Bonds {
				err = evmContracts.Bond(addr, validatorAddress, amount)
				if err != nil {
					return fmt.Errorf("error while bonding: %w", err)
				}
			}
		}

	}

	err = evmContracts.FinalizeInitialization()
	if err != nil {
		return fmt.Errorf("error while calling finalizeInitialization: %w", err)
	}

	log.Info("Deployed Autonity contract", "address", params.AutonityContractAddress)

	return nil
}

func DeployOracleContract(genesisConfig *params.ChainConfig, evmContracts *GenesisEVMContracts) error {
	if genesisConfig.OracleContractConfig == nil {
		log.Info("Using default genesis parameters for the Oracle Contract")
		genesisConfig.OracleContractConfig = params.DefaultGenesisOracleConfig
	}
	if err := genesisConfig.OracleContractConfig.SetDefaults(); err != nil {
		log.Crit("Error with Oracle Contract configuration", "err", err)
	}

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

// AutonityContractCall calls the specified function of the autonity contract
// with the given args, and returns the output unpacked into the result
// interface.
func (c *AutonityContract) AutonityContractCall(statedb vm.StateDB, header *types.Header, function string, result any, args ...any) error {
	packedArgs, err := c.contractABI.Pack(function, args...)
	if err != nil {
		return err
	}
	ret, err := c.CallContractFunc(statedb, header, packedArgs)
	if err != nil {
		return err
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&raw{}) {
		rawPtr := result.(*raw)
		*rawPtr = ret
		return nil
	}
	if err := c.contractABI.UnpackIntoInterface(result, function, ret); err != nil {
		log.Error("Could not unpack returned value", "function", function)
		return err
	}

	return nil
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

func (c *AutonityContract) FinalizeInitialization(header *types.Header, statedb vm.StateDB) error {
	packedArgs, err := c.contractABI.Pack("finalizeInitialization")
	if err != nil {
		return fmt.Errorf("error while generating call data for finalizeInitialization: %w", err)
	}

	_, err = c.CallContractFunc(statedb, header, packedArgs)
	if err != nil {
		return fmt.Errorf("error while calling finalizeInitialization: %w", err)
	}

	return nil
}

// CallContractFunc creates an evm object, uses it to call the
// specified function of the autonity contract with packedArgs and returns the
// packed result. If there is an error making the evm call it will be returned.
// Callers should use the autonity contract ABI to pack and unpack the args and
// result.
func (c *EVMContract) CallContractFunc(statedb vm.StateDB, header *types.Header, contractAddress common.Address, packedArgs []byte) ([]byte, uint64, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, params.DeployerAddress, statedb)
	return evm.Call(vm.AccountRef(params.DeployerAddress), contractAddress, packedArgs, gas, new(big.Int))
}

func (c *EVMContract) CallContractFuncAs(statedb vm.StateDB, header *types.Header, contractAddress common.Address, origin common.Address, packedArgs []byte) ([]byte, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, origin, statedb)
	packedResult, _, err := evm.Call(vm.AccountRef(origin), contractAddress, packedArgs, gas, new(big.Int))
	return packedResult, err
}

func (c *AutonityContract) callGetCommitteeEnodes(state vm.StateDB, header *types.Header, asACN bool) (*types.Nodes, error) {
	var returnedEnodes []string
	err := c.AutonityContractCall(state, header, "getCommitteeEnodes", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes, asACN), nil
}

func (c *AutonityContract) callGetCommittee(state vm.StateDB, header *types.Header) ([]types.CommitteeMember, error) {
	var committee []types.CommitteeMember
	err := c.AutonityContractCall(state, header, "getCommittee", &committee)
	return committee, err
}

func (c *AutonityContract) callGetMinimumBaseFee(state vm.StateDB, header *types.Header) (*big.Int, error) {
	minBaseFee := new(big.Int)
	err := c.AutonityContractCall(state, header, "getMinimumBaseFee", &minBaseFee)
	if err != nil {
		return nil, err
	}
	return minBaseFee, nil
}

func (c *AutonityContract) callFinalize(state vm.StateDB, header *types.Header) (bool, types.Committee, error) {
	var updateReady bool
	var committee types.Committee
	if err := c.AutonityContractCall(state, header, "finalize", &[]any{&updateReady, &committee}); err != nil {
		return false, nil, err
	}
	return updateReady, committee, nil
}

func (c *AutonityContract) callRetrieveContract(state vm.StateDB, header *types.Header) ([]byte, string, error) {
	var bytecode []byte
	var abi string
	if err := c.AutonityContractCall(state, header, "getNewContract", &[]any{&bytecode, &abi}); err != nil {
		return nil, "", err
	}
	return bytecode, abi, nil
}
