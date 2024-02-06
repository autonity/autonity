package test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/ethclient"
)

// test committee members should all start the mining workers.
func TestAutoMiningForCommitteeMembers(t *testing.T) {
	numOfValidators := 2
	testCase := &testCase{
		name:          "Auto mining test, committee members should start mining",
		numValidators: numOfValidators,
		numBlocks:     15,
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			miningChecker(t, validators["V0"], true)
			miningChecker(t, validators["V1"], true)
		},
	}

	runTest(t, testCase)
}

// test non-committee member validators should not start the mining workers.
func TestAutoMiningForNonCommitteeValidator(t *testing.T) {
	numOfValidators := 2
	testCase := &testCase{
		name:          "Auto mining test, none committee member should not start mining",
		numValidators: numOfValidators,
		numBlocks:     15,
		genesisHook: func(g *core.Genesis) *core.Genesis {
			// set committee seat as 1, set bonded stake from 100 to 1 to keep validator 0 out of committee.
			g.Config.AutonityContractConfig.MaxCommitteeSize = 1
			g.Config.AutonityContractConfig.Validators[0].BondedStake.SetUint64(1)
			return g
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			miningChecker(t, validators["V0"], false)
			miningChecker(t, validators["V1"], true)
		},
	}

	runTest(t, testCase)
}

// test on committee seats rotation, new selected validators should start mining workers.
func TestAutoMiningForNewSelectedValidator(t *testing.T) {
	numOfValidators := 2
	operator, err := makeAccount()
	require.NoError(t, err)
	operatorAddr := crypto.PubkeyToAddress(operator.PublicKey)

	testCase := &testCase{
		name:          "Auto mining test, new committee member should start mining",
		numValidators: numOfValidators,
		numBlocks:     25,
		genesisHook: func(g *core.Genesis) *core.Genesis {
			// set committee seat as 1, set bonded stake from 100 to 1 to keep validator 0 out of committee.
			g.Config.AutonityContractConfig.Operator = operatorAddr
			g.Config.AutonityContractConfig.MaxCommitteeSize = 1
			g.Config.AutonityContractConfig.EpochPeriod = 10
			g.Config.AutonityContractConfig.Validators[0].BondedStake.SetUint64(1)
			g.Alloc[operatorAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			return g
		},
		afterHooks: map[string]hook{
			// increase committee size right after block #2, in epoch1, new validator should be mining.
			"V1": setCommitteeSize(
				map[uint64]struct{}{
					2: {},
					3: {},
					4: {},
				},
				operator,
				2,
			),
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			miningChecker(t, validators["V0"], true)
			miningChecker(t, validators["V1"], true)
		},
	}

	runTest(t, testCase)
}

// test on committee seats rotation, removed validators should stop mining workers.
func TestAutoMiningForRemovedValidator(t *testing.T) {
	//t.Skip("todo: this case is unstable in CI, fix it")
	numOfValidators := 2
	operator, err := makeAccount()
	require.NoError(t, err)
	operatorAddr := crypto.PubkeyToAddress(operator.PublicKey)

	testCase := &testCase{
		name:          "Auto mining test, removed committee member should not start mining",
		numValidators: numOfValidators,
		numBlocks:     25,
		genesisHook: func(g *core.Genesis) *core.Genesis {
			// set committee seat as 2
			g.Config.AutonityContractConfig.Operator = operatorAddr
			g.Config.AutonityContractConfig.MaxCommitteeSize = 2
			g.Config.AutonityContractConfig.EpochPeriod = 10
			g.Config.AutonityContractConfig.Validators[0].BondedStake.SetUint64(1)
			g.Alloc[operatorAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			return g
		},
		afterHooks: map[string]hook{
			// decrease committee size right after block #2, in epoch1, validator0 should not be mining.
			"V1": setCommitteeSize(
				map[uint64]struct{}{
					2: {},
					3: {},
					4: {},
				},
				operator,
				1,
			),
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			miningChecker(t, validators["V0"], false)
			miningChecker(t, validators["V1"], true)
		},
	}

	runTest(t, testCase)
}

func miningChecker(t *testing.T, autonity *testNode, isMining bool) {
	c, err := autonity.node.Attach()
	require.NoError(t, err)
	client := ethclient.NewClient(c)
	defer client.Close()
	mining, err := client.IsMining(context.Background())
	require.NoError(t, err)
	require.Equal(t, isMining, mining)
}

func setCommitteeSize(upgradeBlocks map[uint64]struct{}, opKey *ecdsa.PrivateKey, size uint64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(opKey).setCommitteeSize(new(big.Int).SetUint64(size)); err != nil {
			return err
		}
		return nil
	}
}
