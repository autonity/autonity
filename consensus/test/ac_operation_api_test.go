package test

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

const (
	treasuryFee      = 100
	minBaseFee       = 10
	delegationRate   = 1
	epochPeriod      = 5
	unBondingPeriod  = 5
	maxCommitteeSize = 7
	blockPeriod      = 1
)

// Test the contract upgrade mechanism.
// The new contract is ./autonity/solidity/contracts/Upgrade_test.sol
func TestUpgradeMechanism(t *testing.T) {
	numOfValidators := 3
	operator, err := makeAccount()
	require.NoError(t, err)
	initialOperatorAddr := crypto.PubkeyToAddress(operator.PublicKey)
	//var upgradeTxs []*types.Transaction
	bytecode := generated.AutonityUpgradeTestBytecode
	var interactor *interactor
	var transactor *transactor

	testCase := &testCase{
		name:          "Test AC system operator change settings",
		numValidators: numOfValidators,
		numBlocks:     30,
		// set AC configs in genesis hook.
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = initialOperatorAddr
			g.Alloc[initialOperatorAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			setACConfig(g.Config.AutonityContractConfig)
			return g
		},
		beforeHooks: map[string]hook{"V0": func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			if block.Number().Uint64() == 2 {
				interactor = interact(validator.rpcPort)
				transactor = interactor.tx(operator)
			}
			return nil
		}},

		afterHooks: map[string]hook{
			// upgrade is triggered at block #5
			"V0": func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
				var err error
				switch block.Number().Uint64() {
				case 5:
					_, err = transactor.upgradeContract(bytecode[0:len(bytecode)/2], "")
				case 10:
					res, _ := json.Marshal(&generated.AutonityUpgradeTestAbi)
					_, err = transactor.upgradeContract(nil, string(res))
				case 15:
					_, err = transactor.upgradeContract(bytecode[len(bytecode)/2:], "")
				default:
					return nil
				}
				return err
			},
			"V2": func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
				if block.Number().Uint64() == 20 {
					if _, err := transactor.completeContractUpgrade(); err != nil {
						return err
					}
				}
				return nil
			},
			"V1": func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
				if block.Number().Uint64() == 25 {
					if _, err = transactor.transfer(common.HexToAddress("0x111"), big.NewInt(1000)); err != nil {
						return err
					}
				}
				return nil
			},
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			transactor.analyze(t)
			node := validators["V0"]
			client := interact(node.rpcPort)
			defer client.close()
			//upgradeReceipt, err := client.client.TransactionReceipt(context.Background(), upgradeTx.Hash())
			//require.NoError(t, err)
			//t.Log("upgrade gas used:", upgradeReceipt.GasUsed)
			//check that validator 1 has 50 bonded stake. (initial stake is 100)
			val1, err := client.call(node.lastBlock).getValidator(validators["V1"].address)
			val2, err := client.call(node.lastBlock).getValidator(validators["V2"].address)
			val0, err := client.call(node.lastBlock).getValidator(validators["V0"].address)
			t.Log("val0 bs", val0.BondedStake.Uint64())
			t.Log("val1 bs", val1.BondedStake.Uint64())
			t.Log("val2 bs", val2.BondedStake.Uint64())

			require.NoError(t, err)
			stake := new(big.Int).SetUint64(defaultStake)
			stake.Div(stake, big.NewInt(2))
			require.Equal(t, stake.Bytes(), val1.BondedStake.Bytes())
			//check the contract version is now 2.0.0
			version, err := client.call(node.lastBlock).getVersion()
			require.NoError(t, err)
			require.Equal(t, uint64(2), version)
			//check the new transfer operation
			balance, err := client.call(node.lastBlock).balanceOf(common.HexToAddress("0x111"))
			require.NoError(t, err)
			require.Equal(t, big.NewInt(2000), balance)
		},
	}
	runTest(t, testCase)
}

func makeAccount() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func setACConfig(contractConf *params.AutonityContractGenesis) {
	contractConf.BlockPeriod = uint64(blockPeriod)
	contractConf.DelegationRate = uint64(delegationRate)
	contractConf.MaxCommitteeSize = uint64(maxCommitteeSize)
	contractConf.EpochPeriod = uint64(epochPeriod)
	contractConf.MinBaseFee = uint64(minBaseFee)
	contractConf.UnbondingPeriod = uint64(unBondingPeriod)
	contractConf.TreasuryFee = uint64(treasuryFee)
}
