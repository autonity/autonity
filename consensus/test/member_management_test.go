package test

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

/*
  In this file, it create 4 test cases which have similar work flow base on the local e2e test framework's main flow.

  First it setup an autontiy network by according to the genesis hook function, then from the specific chain height, it
  start to issue transactions via the transaction hook function specified for the target node, for example in the
  addValidatorHook, it issues transaction to call autonity contract via operator account to add a new validator.

  Then the test case verify the output from its finalAssert hook function on the specified height of the blockchain, for
  example the addValidatorCheckerHook checks if the new validator is presented in the white list, and its stake balance
  checked too, and finally it checks the total stake supply after the membership updates.

  for the other cases in the file: add stake_holder, participants, or remove user, they follow the same work flow and
  some rules to check the outputs.
*/

func TestMemberManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newValidatorNodeKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	newValidatorPubKey := newValidatorNodeKey.PublicKey
	newValidatorENode := enode.V4DNSUrl(newValidatorPubKey, "VN:8527", 8527, 8527)

	newStakeholderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	newStakeholderPubKey := newStakeholderKey.PublicKey
	newStakeholderEnode := enode.V4DNSUrl(newStakeholderPubKey, "VN:8528", 8528, 8528)

	newParticipantKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	newParticipantPubKey := newParticipantKey.PublicKey
	newParticipantEnode := enode.V4DNSUrl(newParticipantPubKey, "VN:8529", 8529, 8529)

	removeNodeKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	removeNodePubKey := removeNodeKey.PublicKey
	eNodeToRemove := enode.V4DNSUrl(removeNodePubKey, "VM:8527", 8527, 8527)
	addressToRemove := crypto.PubkeyToAddress(removeNodePubKey)

	validatorStake := new(big.Int).SetUint64(200)
	stakeHolderStake := new(big.Int).SetUint64(100)
	participantStake := new(big.Int).SetUint64(0)

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}

		// the user to be removed.
		user := &params.User{
			Address: &addressToRemove,
			Enode:   eNodeToRemove,
			Type:    "participant",
			Stake:   0,
		}
		g.Config.AutonityContractConfig.Users = append(g.Config.AutonityContractConfig.Users, *user)
		return g
	}

	addUser := func(operator *ecdsa.PrivateKey, port int, userPubKey ecdsa.PublicKey, stake *big.Int, enode string, userType uint8) error {
		return interact(port).tx(operator).addUser(crypto.PubkeyToAddress(userPubKey), stake, enode, userType)
	}

	addUserHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			err := addUser(operatorKey, validator.rpcPort, newValidatorPubKey, validatorStake, newValidatorENode, uint8(2))
			if err != nil {
				return false, nil, err
			}

			err = addUser(operatorKey, validator.rpcPort, newStakeholderPubKey, stakeHolderStake, newStakeholderEnode, uint8(1))
			if err != nil {
				return false, nil, err
			}

			err = addUser(operatorKey, validator.rpcPort, newParticipantPubKey, participantStake, newParticipantEnode, uint8(0))
			if err != nil {
				return false, nil, err
			}
		}
		return false, nil, nil
	}

	// to check user membership, user type, stake balance.
	validateAddedUser := func(t *testing.T, port int, height uint64, address common.Address, eNode string, role uint8, stake uint64, economicMetric Struct1) {
		whiteList, err := interact(port).call(height).getWhitelist()
		require.NoError(t, err)
		assert.Contains(t, whiteList, eNode, "new user ENode is not presented from ENode list")
		isMember, err := interact(port).call(height).checkMember(address)
		require.NoError(t, err)
		assert.True(t, isMember, "wrong membership for added user")

		// check validator and stakeholder's stake balance
		actualStake, err := interact(port).call(height).getAccountStake(address)
		if role != uint8(0) {
			require.NoError(t, err)
			require.Equal(t, stake, actualStake.Uint64())
		} else {
			// for participants, it is not allow to have stake, getAccountStake is limited only for stakeholder and validator.
			require.EqualError(t, err, "execution reverted: address not allowed to use stake")
		}

		for index, v := range economicMetric.Accounts {
			if v == address {
				assert.Equal(t, role, economicMetric.Usertypes[index], "user type is not expected")
				break
			}
		}
	}

	addUserCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		port := validators["VA"].rpcPort
		lastHeight := validators["VA"].lastBlock
		curNetworkMetrics, err := interact(port).call(lastHeight).dumpEconomicsMetricData()
		require.NoError(t, err)

		validateAddedUser(t, port, lastHeight, crypto.PubkeyToAddress(newValidatorPubKey), newValidatorENode, uint8(2), validatorStake.Uint64(), curNetworkMetrics)
		validateAddedUser(t, port, lastHeight, crypto.PubkeyToAddress(newStakeholderPubKey), newStakeholderEnode, uint8(1), stakeHolderStake.Uint64(), curNetworkMetrics)
		validateAddedUser(t, port, lastHeight, crypto.PubkeyToAddress(newParticipantPubKey), newParticipantEnode, uint8(0), participantStake.Uint64(), curNetworkMetrics)

		// compare the total stake supply before and after new node added.
		initNetworkMetrics, err := interact(validators["VA"].rpcPort).call(3).dumpEconomicsMetricData()
		require.NoError(t, err)

		// new_total_stake - init_total_stake == new added (validatorStake + stakeHolderStake + participantStake)
		b := curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64()
		assert.Equal(t, b, validatorStake.Uint64()+stakeHolderStake.Uint64()+participantStake.Uint64(), "stake total supply is not expected")
	}

	removeUserHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			return true, nil, interact(validator.rpcPort).tx(operatorKey).removeUser(addressToRemove)
		}
		return false, nil, nil
	}

	removeUserCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		isMember, err := interact(validators["VA"].rpcPort).call(validators["VA"].lastBlock).checkMember(addressToRemove)
		require.NoError(t, err)
		assert.False(t, isMember, "wrong membership for removed user")
	}

	cases := []*testCase{
		{
			name:          "add users",
			numValidators: 6,
			numBlocks:     30,
			txPerPeer:     5,
			sendTransactionHooks: map[string]sendTransactionHook{
				"VA": addUserHook,
			},
			genesisHook: genesisHook,
			finalAssert: addUserCheckerHook,
		},
		{
			name:          "remove user",
			numValidators: 6,
			numBlocks:     20,
			txPerPeer:     1,
			sendTransactionHooks: map[string]sendTransactionHook{
				"VD": removeUserHook,
			},
			genesisHook: genesisHook,
			finalAssert: removeUserCheckerHook,
		},
	}
	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
