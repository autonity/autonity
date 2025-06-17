package autonity

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common/math"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/internal/testrand"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

func TestGenesisSteps(t *testing.T) {
	newEVM := func() *vm.EVM {
		stateDB, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		require.NoError(t, err)

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

		return vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})
	}

	t.Run("Test autonity deploy step", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(params.TestChainConfig, []GenesisBond{}, evm, []genesisStep{deployAutonityContract})
		require.NoError(t, err)

		// Check that the autonity contract was deployed
		code := evm.StateDB.GetCode(params.AutonityContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test execute genesis delegations", func(t *testing.T) {
		evm := newEVM()
		validator1 := *params.TestChainConfig.AutonityContractConfig.Validators[0].NodeAddress
		validator2 := *params.TestChainConfig.AutonityContractConfig.Validators[1].NodeAddress

		ntnBalance1 := big.NewInt(100)
		bondedNtnBalance1 := big.NewInt(50)

		account2 := testrand.Address()
		ntnBalance2 := big.NewInt(200)
		bondedNtnBalance2 := big.NewInt(100)

		err := executeGenesisSequence(params.TestChainConfig, []GenesisBond{
			// validator self bonded
			{
				Staker:        validator1,
				NewtonBalance: ntnBalance1,
				Bonds: []Delegation{
					{
						Validator: validator1,
						Amount:    bondedNtnBalance1,
					},
				},
			},
			// delegated to another account
			{
				Staker:        account2,
				NewtonBalance: ntnBalance2,
				Bonds: []Delegation{
					{
						Validator: validator2,
						Amount:    bondedNtnBalance2,
					},
				},
			},
		}, evm, []genesisStep{deployAutonityContract, executeGenesisDelegations})
		require.NoError(t, err)
		balanceOf := func(addr common.Address) *big.Int {
			result := new(big.Int)
			_, err := AutonityContractCall(&generated.AutonityAbi, evm, "balanceOf", &result, addr)
			require.NoError(t, err)
			return result
		}
		require.Equal(t, ntnBalance1, balanceOf(validator1))
		require.Equal(t, ntnBalance2, balanceOf(account2))
	})

	t.Run("Test create autonity schedules", func(t *testing.T) {
		evm := newEVM()
		config := params.TestChainConfig
		config.AutonityContractConfig.Schedules = []params.Schedule{
			{
				Start:         big.NewInt(time.Now().Unix() + 10),
				TotalDuration: big.NewInt(100),
				Amount:        big.NewInt(100),
				VaultAddress:  params.NonStakeableVestingContractAddress,
			},
		}
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{deployAutonityContract, executeGenesisDelegations, createAutonitySchedules},
		)
		require.NoError(t, err)

		schedule := new(bindings.IScheduleControllerSchedule)
		_, err = AutonityContractCall(
			&generated.AutonityAbi,
			evm,
			"getSchedule",
			&schedule,
			config.AutonityContractConfig.Schedules[0].VaultAddress,
			big.NewInt(int64(0)),
		)
		require.NoError(t, err)

		require.Equal(t, config.AutonityContractConfig.Schedules[0].Start, schedule.Start)
		require.Equal(t, config.AutonityContractConfig.Schedules[0].TotalDuration, schedule.TotalDuration)
		require.Equal(t, config.AutonityContractConfig.Schedules[0].Amount, schedule.TotalAmount)
	})

	t.Run("Test finalize autonity initialization", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{deployAutonityContract, executeGenesisDelegations, createAutonitySchedules},
		)
		require.NoError(t, err)
		getCommitteeEnodes := func() []string {
			var result []string
			_, err := AutonityContractCall(&generated.AutonityAbi, evm, "getCommitteeEnodes", &result)
			require.NoError(t, err)
			return result
		}
		require.Empty(t, getCommitteeEnodes())

		err = executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{finalizeAutonityInitialization},
		)
		require.NoError(t, err)
		require.NotEmpty(t, getCommitteeEnodes())
	})

	t.Run("Test deploy accountability contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				executeGenesisDelegations,
				createAutonitySchedules,
				finalizeAutonityInitialization,
				deployAccountabilityContract,
			},
		)
		require.NoError(t, err)

		// Check that the accountability contract was deployed
		code := evm.StateDB.GetCode(params.AccountabilityContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy oracle contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				executeGenesisDelegations,
				createAutonitySchedules,
				finalizeAutonityInitialization,
				deployAccountabilityContract,
				deployOracleContract,
			},
		)
		require.NoError(t, err)

		// Check that the oracle contract was deployed
		code := evm.StateDB.GetCode(params.OracleContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy ACU contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				executeGenesisDelegations,
				createAutonitySchedules,
				finalizeAutonityInitialization,
				deployAccountabilityContract,
				deployOracleContract,
				deployACUContract,
			},
		)
		require.NoError(t, err)

		// Check that the ACU contract was deployed
		code := evm.StateDB.GetCode(params.ACUContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy supply control contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				executeGenesisDelegations,
				createAutonitySchedules,
				finalizeAutonityInitialization,
				deployAccountabilityContract,
				deployOracleContract,
				deployACUContract,
				deploySupplyControlContract,
			},
		)
		require.NoError(t, err)

		// Check that the supply control contract was deployed
		code := evm.StateDB.GetCode(params.SupplyControlContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy stabilization contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				executeGenesisDelegations,
				createAutonitySchedules,
				finalizeAutonityInitialization,
				deployAccountabilityContract,
				deployOracleContract,
				deployACUContract,
				deploySupplyControlContract,
				deployStabilizationContract,
			},
		)
		require.NoError(t, err)

		// Check that the stabilization contract was deployed
		code := evm.StateDB.GetCode(params.StabilizationContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy upgrade manager contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				executeGenesisDelegations,
				createAutonitySchedules,
				finalizeAutonityInitialization,
				deployAccountabilityContract,
				deployOracleContract,
				deployACUContract,
				deploySupplyControlContract,
				deployStabilizationContract,
				deployUpgradeManagerContract,
			},
		)
		require.NoError(t, err)

		// Check that the upgrade manager contract was deployed
		code := evm.StateDB.GetCode(params.UpgradeManagerContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy inflation control contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
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
			},
		)
		require.NoError(t, err)

		// Check that the inflation controller contract was deployed
		code := evm.StateDB.GetCode(params.InflationControllerContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test deploy stakable vesting manager contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
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
			},
		)
		require.NoError(t, err)

		// Check that the stakable vesting manager contract was deployed
		code := evm.StateDB.GetCode(params.StakeableVestingManagerContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test create stakable vesting contracts", func(t *testing.T) {
		evm := newEVM()
		owner1 := testrand.Address()
		owner2 := testrand.Address()
		config := &params.ChainConfig{
			AutonityContractConfig: params.TestAutonityContractConfig,
			OracleContractConfig:   params.TestOracleConfig,
			StakeableVestingConfig: &params.StakeableVestingGenesis{
				StakeableContracts: []params.StakeableVestingData{
					{
						Beneficiary:   owner1,
						Amount:        new(big.Int).SetUint64(100),
						Start:         big.NewInt(time.Now().Unix() + 10),
						CliffDuration: new(big.Int).SetUint64(10),
						TotalDuration: new(big.Int).SetUint64(100),
					},
					{
						Beneficiary:   owner2,
						Amount:        new(big.Int).SetUint64(200),
						Start:         big.NewInt(time.Now().Unix() + 10),
						CliffDuration: new(big.Int).SetUint64(10),
						TotalDuration: new(big.Int).SetUint64(100),
					},
				},
				TotalNominal: big.NewInt(300),
			},
		}

		config.SetDefaults()

		err := executeGenesisSequence(
			config,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
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
			},
		)
		require.NoError(t, err)

		// Check that the stakable vesting contracts were created
		totalContracts := func(address common.Address) *big.Int {
			result := new(big.Int)
			callContractFunc(
				t,
				evm,
				params.StakeableVestingManagerContractAddress,
				&generated.StakeableVestingManagerAbi,
				"totalContracts",
				&result,
				address,
			)
			return result
		}
		require.Equal(t, big.NewInt(1), totalContracts(owner1))
		require.Equal(t, big.NewInt(1), totalContracts(owner2))

	})

	t.Run("Test deploy non-stakable vesting manager contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
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
			},
		)
		require.NoError(t, err)

		// Check that the non-stakable vesting manager contract was deployed
		code := evm.StateDB.GetCode(params.NonStakeableVestingContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test create non-stakable vesting contracts", func(t *testing.T) {
		owner1 := testrand.Address()
		owner2 := testrand.Address()

		evm := newEVM()
		config := &params.ChainConfig{
			AutonityContractConfig: params.TestAutonityContractConfig,
			OracleContractConfig:   params.TestOracleConfig,
			NonStakeableVestingConfig: &params.NonStakeableVestingGenesis{
				NonStakeableContracts: []params.NonStakeableVestingData{
					{
						Beneficiary:   owner1,
						Amount:        big.NewInt(100),
						ScheduleID:    big.NewInt(0),
						CliffDuration: big.NewInt(10),
					},
					{
						Beneficiary:   owner2,
						Amount:        big.NewInt(200),
						ScheduleID:    big.NewInt(0),
						CliffDuration: big.NewInt(10),
					},
				},
			},
		}
		config.AutonityContractConfig.Schedules = []params.Schedule{
			{
				Start:         big.NewInt(time.Now().Unix() + 10),
				TotalDuration: big.NewInt(100),
				Amount:        big.NewInt(500),
				VaultAddress:  params.NonStakeableVestingContractAddress,
			},
		}
		config.SetDefaults()

		err := executeGenesisSequence(
			config,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
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
			},
		)
		require.NoError(t, err)

		// Check that the nonstakable vesting contracts were created
		totalContracts := func(address common.Address) *big.Int {
			result := new(big.Int)
			callContractFunc(
				t,
				evm,
				params.NonStakeableVestingContractAddress,
				&generated.StakeableVestingManagerAbi,
				"totalContracts",
				&result,
				address,
			)
			return result
		}
		require.Equal(t, big.NewInt(1), totalContracts(owner1))
		require.Equal(t, big.NewInt(1), totalContracts(owner2))

	})

	t.Run("Test deploy omission accountability contract", func(t *testing.T) {
		evm := newEVM()
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
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
			},
		)
		require.NoError(t, err)

		// Check that the omission accountability contract was deployed
		code := evm.StateDB.GetCode(params.OmissionAccountabilityContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test genesis sequence verifier", func(t *testing.T) {
		stake := validatorTotalStake(params.TestChainConfig)
		newTestConfig, err := chainConfig(params.TestChainConfig, stake, stake)
		require.NoError(t, err)

		evm := newEVM()
		err = executeGenesisSequence(
			newTestConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{
				deployAutonityContract,
				finalizeAutonityInitialization,
				verifyGenesisSequence,
			},
		)
		require.NoError(t, err)
	})
}

func validatorTotalStake(config *params.ChainConfig) *big.Int {
	stake := new(big.Int)
	for _, v := range config.AutonityContractConfig.Validators {
		stake.Add(stake, v.BondedStake)
	}
	return stake
}

func chainConfig(config *params.ChainConfig, genesisMint, genesisBond *big.Int) (*params.ChainConfig, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	copy := new(params.ChainConfig)
	if err = json.Unmarshal(data, &copy); err != nil {
		return nil, err
	}

	copy.AutonityContractConfig.VerifyGenesisSequence = true
	copy.AutonityContractConfig.TokenBond = (*math.HexOrDecimal256)(genesisBond)
	copy.AutonityContractConfig.TokenMint = (*math.HexOrDecimal256)(genesisMint)
	return copy, nil
}

func callContractFunc(
	t *testing.T,
	evm *vm.EVM,
	contractAddress common.Address,
	contractAbi *abi.ABI,
	function string,
	result interface{},
	args ...interface{},
) {
	packedArgs, err := contractAbi.Pack(function, args...)
	require.NoError(t, err)
	ret, _, err := evm.Call(
		vm.AccountRef(params.DeployerAddress),
		contractAddress,
		packedArgs,
		math.MaxUint64,
		common.Big0,
	)
	require.NoError(t, err)

	err = contractAbi.UnpackIntoInterface(result, function, ret)
	require.NoError(t, err)
}
