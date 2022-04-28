package test

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

const (
	treasuryFee      = 100
	minBaseFee       = 10
	delegationRate   = 1
	epochPeriod      = 5
	unBondingPeriod  = 5
	maxCommitteeSize = 7
	blockPeriod      = 1
	defaultVersion   = "v0.0.0"
)

var (
	// new settings to be submitted for protocol in testcase.
	newMinBaseFee      = new(big.Int).SetUint64(20)
	newCommitteeSize   = new(big.Int).SetUint64(10)
	newUnBondingPeriod = new(big.Int).SetUint64(30)
	newEpochPeriod     = new(big.Int).SetUint64(45)
	newTreasuryFee     = new(big.Int).SetUint64(10)
	mintAmount         = new(big.Int).SetUint64(100)
	burnAmount         = new(big.Int).SetUint64(50)
)

// test registerValidator, unRegisterValidators, bond and unbond operations.
func TestACPublicWritters(t *testing.T) {
	numOfValidators := 2

	operator, err := makeAccount()
	require.NoError(t, err)
	operatorAddr := crypto.PubkeyToAddress(operator.PublicKey)

	newValidator, err := makeAccount()
	require.NoError(t, err)
	newValidatorAddr := crypto.PubkeyToAddress(newValidator.PublicKey)
	enodeUrl := enode.V4DNSUrl(newValidator.PublicKey, "127.0.0.1", 30303, 30303) + ":30303"

	// newton to be mint
	amount := new(big.Int).SetUint64(10)

	cases := []*testCase{
		{
			name:          "Test register new validator",
			numValidators: numOfValidators,
			numBlocks:     10,
			txPerPeer:     1,
			// register validator right after block #5 is committed from client V0.
			afterHooks: map[string]hook{
				"V0": registerValidatorHook(map[uint64]struct{}{
					5: {},
				},
					enodeUrl,
				),
			},
			finalAssert: func(t *testing.T, validators map[string]*testNode) {
				client := interact(validators["V0"].rpcPort)
				defer client.close()
				val, err := client.call(validators["V0"].lastBlock).getValidator(newValidatorAddr)
				require.NoError(t, err)
				require.Equal(t, newValidatorAddr, val.Addr)
			},
		},
		{
			name:          "bond stake to validator",
			numValidators: numOfValidators,
			numBlocks:     10,
			txPerPeer:     1,
			genesisHook: func(g *core.Genesis) *core.Genesis {
				g.Config.AutonityContractConfig.Operator = operatorAddr
				// pre-mine Auton for system operator and new validator.
				g.Alloc[operatorAddr] = core.GenesisAccount{
					Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
				}
				g.Alloc[newValidatorAddr] = core.GenesisAccount{
					Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
				}
				return g
			},
			// mint newton for new validator right after block #3 via client V0.
			// bond newton to validator right after block #7 via client V1.
			afterHooks: map[string]hook{
				"V0": mintStakeHook(map[uint64]struct{}{
					2: {},
				},
					operator,
					newValidator,
					amount,
				),
				"V1": bondStakeHook(map[uint64]struct{}{
					7: {},
				},
					newValidator,
					amount,
				),
			},
			finalAssert: func(t *testing.T, validators map[string]*testNode) {
				node := validators["V1"]
				client := interact(node.rpcPort)
				defer client.close()
				reqs, err := client.call(node.lastBlock).getBondingReq(new(big.Int).SetUint64(0),
					new(big.Int).SetUint64(3))
				require.NoError(t, err)
				require.Equal(t, node.EthAddress(), reqs[2].Delegatee)
				require.Equal(t, newValidatorAddr, reqs[2].Delegator)
				require.Equal(t, amount, reqs[2].Amount)
			},
		},
		{
			name:          "unbond stake from validator",
			numValidators: numOfValidators,
			numBlocks:     10,
			txPerPeer:     1,
			afterHooks: map[string]hook{
				"V0": unBondStakeHook(map[uint64]struct{}{
					2: {},
				},
					amount,
				),
			},
			finalAssert: func(t *testing.T, validators map[string]*testNode) {
				node := validators["V0"]
				client := interact(node.rpcPort)
				defer client.close()
				reqs, err := client.call(node.lastBlock).getUnBondingReq(new(big.Int).SetUint64(0),
					new(big.Int).SetUint64(1))
				require.NoError(t, err)
				require.Equal(t, node.EthAddress(), reqs[0].Delegatee)
				require.Equal(t, node.EthAddress(), reqs[0].Delegator)
				require.Equal(t, amount, reqs[0].Amount)
			},
		},
	}

	for _, testcase := range cases {
		runTest(t, testcase)
	}
}

// test system settings management by operator account.
func TestACSystemOperatorOPs(t *testing.T) {
	numOfValidators := 3
	initialOperator, err := makeAccount()
	require.NoError(t, err)
	initialOperatorAddr := crypto.PubkeyToAddress(initialOperator.PublicKey)

	newOperator, err := makeAccount()
	require.NoError(t, err)
	newOperatorAddr := crypto.PubkeyToAddress(newOperator.PublicKey)

	testCase := &testCase{
		name:          "Test AC system operator change settings",
		numValidators: numOfValidators,
		numBlocks:     15,
		txPerPeer:     1,
		// set AC configs in genesis hook.
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = initialOperatorAddr
			g.Alloc[initialOperatorAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			g.Alloc[newOperatorAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			setACConfig(g.Config.AutonityContractConfig)
			return g
		},
		afterHooks: map[string]hook{
			// change settings right after block #1 is committed from client V0.
			"V0": changeSettingHook(
				map[uint64]struct{}{
					1: {},
				},
				initialOperator,
				newOperatorAddr,
			),
			// burn stake right after block #5 from client V1.
			"V1": burnStakeHook(
				map[uint64]struct{}{
					5: {},
				},
				initialOperator,
				newOperatorAddr,
			),
			// change operator after block #10 from client V3.
			"V2": setOperatorHook(
				map[uint64]struct{}{
					10: {},
				},
				initialOperator,
				newOperatorAddr,
			),
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			node := validators["V0"]
			client := interact(node.rpcPort)
			defer client.close()
			mBaseFee, err := client.call(node.lastBlock).getMinBaseFee()
			require.NoError(t, err)
			require.Equal(t, newMinBaseFee.Uint64(), mBaseFee.Uint64())
			comSize, err := client.call(node.lastBlock).getMaxCommitteeSize()
			require.NoError(t, err)
			require.Equal(t, newCommitteeSize.Uint64(), comSize.Uint64())
			newOP, err := client.call(node.lastBlock).getOperator()
			require.NoError(t, err)
			require.Equal(t, newOperatorAddr, newOP)
			balance, err := client.call(node.lastBlock).balanceOf(newOperatorAddr)
			require.NoError(t, err)
			require.Equal(t, mintAmount.Sub(mintAmount, burnAmount).Uint64(), balance.Uint64())

			// new operator can change settings after the account update without revert error.
			_, err = client.tx(newOperator).setEpochPeriod(new(big.Int).SetUint64(45))
			require.NoError(t, err)
		},
	}
	runTest(t, testCase)
}

// test system settings / state getters
func TestACStateGetters(t *testing.T) {
	operatorKey, err := makeAccount()
	require.NoError(t, err)
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	numOfValidators := 2
	testCase := &testCase{
		name:          "Test AC state getters",
		numValidators: numOfValidators,
		numBlocks:     10,
		txPerPeer:     1,
		// set AC configs in genesis hook.
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Config.AutonityContractConfig.Treasury = operatorAddress
			setACConfig(g.Config.AutonityContractConfig)
			return g
		},
		// start AC state getter verifications right after block #5 is committed from client V0.
		afterHooks: map[string]hook{
			"V0": acStateGettersHook(map[uint64]struct{}{
				5: {},
			},
				operatorAddress,
				numOfValidators,
			),
		},
	}
	runTest(t, testCase)
}

func burnStakeHook(upgradeBlocks map[uint64]struct{}, op *ecdsa.PrivateKey, ac common.Address) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(op).burn(ac, burnAmount); err != nil {
			return err
		}
		return nil
	}
}

func setOperatorHook(upgradeBlocks map[uint64]struct{}, operator *ecdsa.PrivateKey, opAddr common.Address) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(operator).setOperator(opAddr); err != nil {
			return err
		}
		return nil
	}
}

func changeSettingHook(upgradeBlocks map[uint64]struct{}, opKey *ecdsa.PrivateKey, newOpAddr common.Address) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(opKey).setMinBaseFee(newMinBaseFee); err != nil {
			return err
		}
		if _, err := interaction.tx(opKey).setCommitteeSize(newCommitteeSize); err != nil {
			return err
		}
		if _, err := interaction.tx(opKey).setUnBondingPeriod(newUnBondingPeriod); err != nil {
			return err
		}
		if _, err := interaction.tx(opKey).setEpochPeriod(newEpochPeriod); err != nil {
			return err
		}
		if _, err := interaction.tx(opKey).setTreasuryAccount(newOpAddr); err != nil {
			return err
		}
		if _, err := interaction.tx(opKey).setTreasuryFee(newTreasuryFee); err != nil {
			return err
		}
		if _, err := interaction.tx(opKey).mint(newOpAddr, mintAmount); err != nil {
			return err
		}
		return nil
	}
}

func unBondStakeHook(upgradeBlocks map[uint64]struct{}, amount *big.Int) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(validator.privateKey).unbond(validator.EthAddress(), amount); err != nil {
			return err
		}
		return nil
	}
}

func bondStakeHook(upgradeBlocks map[uint64]struct{}, newVal *ecdsa.PrivateKey, amount *big.Int) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(newVal).bond(validator.EthAddress(), amount); err != nil {
			return err
		}
		return nil
	}
}

func mintStakeHook(upgradeBlocks map[uint64]struct{}, operator *ecdsa.PrivateKey, newVal *ecdsa.PrivateKey, amount *big.Int) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		newValAddr := crypto.PubkeyToAddress(newVal.PublicKey)
		if _, err := interaction.tx(operator).mint(newValAddr, amount); err != nil {
			return err
		}
		return nil
	}
}

func registerValidatorHook(upgradeBlocks map[uint64]struct{}, enode string) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(validator.privateKey).registerValidator(enode); err != nil {
			return err
		}
		return nil
	}
}

func acStateGettersHook(upgradeBlocks map[uint64]struct{}, operator common.Address, numVals int) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := upgradeBlocks[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()

		if err := checkVersion(interaction, blockNum, defaultVersion); err != nil {
			return err
		}

		if err := checkCommittee(interaction, blockNum, numVals); err != nil {
			return err
		}

		if err := checkValidators(interaction, blockNum, numVals); err != nil {
			return err
		}

		if err := checkValidator(interaction, blockNum, validator.EthAddress(), validator.netNode.url); err != nil {
			return err
		}

		if err := checkMaxCommitteeSize(interaction, blockNum, maxCommitteeSize); err != nil {
			return err
		}

		if err := checkCommitteeEnodes(interaction, blockNum, numVals); err != nil {
			return err
		}

		if err := checkMinBaseFee(interaction, blockNum, minBaseFee); err != nil {
			return err
		}

		if err := checkOperatorAddress(interaction, blockNum, operator); err != nil {
			return err
		}

		if err := checkNewContract(interaction, blockNum, "", ""); err != nil {
			return err
		}

		start := new(big.Int).SetUint64(0)
		end := new(big.Int).SetUint64(uint64(numVals))
		if err := checkBondingReqs(interaction, blockNum, start, end, numVals); err != nil {
			return err
		}
		if err := checkUnBondingReqs(interaction, blockNum, start, end, numVals); err != nil {
			return err
		}
		return nil
	}
}

func checkVersion(client *interactor, height uint64, expected string) error {
	version, err := client.call(height).getVersion()
	if err != nil {
		return err
	}
	if version != expected {
		return fmt.Errorf("unexpected version")
	}
	return nil
}

func checkMaxCommitteeSize(client *interactor, height uint64, size int) error {
	maxSize, err := client.call(height).getMaxCommitteeSize()
	if err != nil {
		return err
	}
	if maxSize.Uint64() != uint64(size) {
		return fmt.Errorf("unexpected max committee size")
	}
	return nil
}

// todo: check each enode urls
func checkCommitteeEnodes(client *interactor, height uint64, numEnodes int) error {
	enodes, err := client.call(height).getCommitteeEnodes()
	if err != nil {
		return err
	}
	if len(enodes) != numEnodes {
		return fmt.Errorf("unexpected committee enodes")
	}
	return nil
}

// todo: get committtee from genesis file, and check them on each seat
func checkCommittee(client *interactor, height uint64, lenCommittee int) error {
	committee, err := client.call(height).getCommittee()
	if err != nil {
		return err
	}
	if len(committee) != lenCommittee {
		return fmt.Errorf("unexpected committee")
	}
	return nil
}

// todo: check each validator
func checkValidators(client *interactor, height uint64, numVals int) error {
	vals, err := client.call(height).getValidators()
	if err != nil {
		return err
	}
	if len(vals) != numVals {
		return fmt.Errorf("unexpected validators")
	}
	return nil
}

// todo: check each property of validator
func checkValidator(client *interactor, height uint64, address common.Address, enode string) error {
	val, err := client.call(height).getValidator(address)
	if err != nil {
		return err
	}
	if enode != val.Enode {
		return fmt.Errorf("unexpected validator")
	}
	return nil
}

func checkMinBaseFee(client *interactor, height uint64, minBaseFee int) error {
	fee, err := client.call(height).getMinBaseFee()
	if err != nil {
		return err
	}
	if fee.Uint64() != uint64(minBaseFee) {
		return fmt.Errorf("unexpected min base fee")
	}
	return nil
}

func checkOperatorAddress(client *interactor, height uint64, opAddr common.Address) error {
	op, err := client.call(height).getOperator()
	if err != nil {
		return err
	}
	if op != opAddr {
		return fmt.Errorf("unexpected operator account")
	}
	return nil
}

func checkNewContract(client *interactor, height uint64, byteCode string, abi string) error {
	byteC, a, err := client.call(height).getNewContract()
	if err != nil {
		return err
	}
	if byteC != byteCode && a != abi {
		return fmt.Errorf("unexpected new contract")
	}
	return nil
}

// todo check each bonding requests
func checkBondingReqs(client *interactor, height uint64, s *big.Int, e *big.Int, num int) error {
	reqs, err := client.call(height).getBondingReq(s, e)
	if err != nil {
		return err
	}
	if len(reqs) != num {
		return fmt.Errorf("unexpected bonding reqs")
	}
	return nil
}

// todo check each unbonding requests.
func checkUnBondingReqs(client *interactor, height uint64, s *big.Int, e *big.Int, num int) error {
	reqs, err := client.call(height).getUnBondingReq(s, e)
	if err != nil {
		return err
	}
	if len(reqs) != num {
		return fmt.Errorf("unexpected unbonding reqs")
	}
	return nil
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
