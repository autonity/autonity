package autonity

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
)

func TestGenesisSteps(t *testing.T) {
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

	evm := vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})

	t.Run("Test autonity deploy step", func(t *testing.T) {
		err := executeGenesisSequence(params.TestChainConfig, []GenesisBond{}, evm, []genesisStep{deployAutonityContract})
		require.NoError(t, err)

		// Check that the autonity contract was deployed
		code := evm.StateDB.GetCode(params.AutonityContractAddress)
		require.NotEmpty(t, code)
	})

	t.Run("Test execute genesis delegations", func(t *testing.T) {
		validator1 := *params.TestChainConfig.AutonityContractConfig.Validators[0].NodeAddress
		validator2 := *params.TestChainConfig.AutonityContractConfig.Validators[1].NodeAddress

		ntnBalance1 := big.NewInt(100)
		bondedNtnBalance1 := big.NewInt(50)

		account2 := randomAddress(t)
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
			_, err := AutonityContractCall(evm, "balanceOf", &result, addr)
			require.NoError(t, err)
			return result
		}
		require.Equal(t, ntnBalance1, balanceOf(validator1))
		require.Equal(t, ntnBalance2, balanceOf(account2))
	})

	t.Run("Test create autonity schedules", func(t *testing.T) {
		err := executeGenesisSequence(
			params.TestChainConfig,
			[]GenesisBond{},
			evm,
			[]genesisStep{deployAutonityContract, executeGenesisDelegations, createAutonitySchedules},
		)
		require.NoError(t, err)

	})
	// TODO(scott) test the remaining genesis steps
}

func randomAddress(t *testing.T) common.Address {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return crypto.PubkeyToAddress(key.PublicKey)
}
