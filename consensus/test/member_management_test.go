package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
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
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newValidatorNodeKey, err := keygenerator.Next()
	require.NoError(t, err)
	newValidatorPubKey := newValidatorNodeKey.PublicKey
	newValidatorENode := enode.V4DNSUrl(newValidatorPubKey, "VN:8527", 8527, 8527)

	newStakeholderKey, err := keygenerator.Next()
	require.NoError(t, err)
	newStakeholderPubKey := newStakeholderKey.PublicKey
	newStakeholderEnode := enode.V4DNSUrl(newStakeholderPubKey, "VN:8528", 8528, 8528)

	newParticipantKey, err := keygenerator.Next()
	require.NoError(t, err)
	newParticipantPubKey := newParticipantKey.PublicKey
	newParticipantEnode := enode.V4DNSUrl(newParticipantPubKey, "VN:8529", 8529, 8529)

	removeNodeKey, err := keygenerator.Next()
	require.NoError(t, err)
	removeNodePubKey := removeNodeKey.PublicKey
	eNodeToRemove := enode.V4DNSUrl(removeNodePubKey, "VM:8527", 8527, 8527)
	addressToRemove := crypto.PubkeyToAddress(removeNodePubKey)

	stakeBalance := new(big.Int).SetUint64(300)

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

	addUserHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			err := interact(validator.rpcPort).tx(operatorKey).addUser(crypto.PubkeyToAddress(newValidatorPubKey), stakeBalance, newValidatorENode, uint8(2))
			if err != nil {
				return false, nil, err
			}
			err = interact(validator.rpcPort).tx(operatorKey).addUser(crypto.PubkeyToAddress(newStakeholderPubKey), stakeBalance, newStakeholderEnode, uint8(1))
			if err != nil {
				return false, nil, err
			}
			err = interact(validator.rpcPort).tx(operatorKey).addUser(crypto.PubkeyToAddress(newParticipantPubKey), stakeBalance, newParticipantEnode, uint8(0))
			if err != nil {
				return false, nil, err
			}
		}
		return false, nil, nil
	}

	addUserCheckerHook := func(t *testing.T, validators map[string]*testNode) {

		whiteList, err := interact(validators["VA"].rpcPort).call(validators["VA"].lastBlock).getWhitelist()
		require.NoError(t, err)

		assert.Contains(t, whiteList, newValidatorENode, "new validator eNode is not presented from member list")
		assert.Contains(t, whiteList, newStakeholderEnode, "new stake holder eNode is not presented from member list")
		assert.Contains(t, whiteList, newParticipantEnode, "new participant eNode is not presented from member list")

		// check node role and its stake balance.
		curNetworkMetrics, err := interact(validators["VA"].rpcPort).call(validators["VA"].lastBlock).dumpEconomicsMetricData()
		require.NoError(t, err)

		newValidatorFind := false
		newStakeHolderFind := false
		newParticipantFind := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newValidatorPubKey) {
				newValidatorFind = true
				assert.Equal(t, curNetworkMetrics.Stakes[index].Uint64(), stakeBalance.Uint64(), "new validator's stake is not expected")
				assert.Equal(t, 2, int(curNetworkMetrics.Usertypes[index]), "new validator's user type is not expected")
			}
			if v == crypto.PubkeyToAddress(newStakeholderPubKey) {
				newStakeHolderFind = true
				assert.Equal(t, curNetworkMetrics.Stakes[index].Uint64(), stakeBalance.Uint64(), "new stakeholder's stake is not expected")
				assert.Equal(t, 1, int(curNetworkMetrics.Usertypes[index]), "new stakeholder's user type is not expected")
			}
			if v == crypto.PubkeyToAddress(newParticipantPubKey) {
				newParticipantFind = true
				assert.Equal(t, curNetworkMetrics.Stakes[index].Uint64(), uint64(0), "new participant's stake is not expected")
				assert.Equal(t, 0, int(curNetworkMetrics.Usertypes[index]), "new participant's user type is not expected")
			}
		}

		// compare the total stake supply before and after new node added.
		initNetworkMetrics, err := interact(validators["VA"].rpcPort).call(3).dumpEconomicsMetricData()
		require.NoError(t, err)

		b := curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64()
		assert.Equal(t, b, stakeBalance.Uint64()*2, "stake total supply is not expected")
		assert.True(t, newValidatorFind, "new validator is not presented")
		assert.True(t, newStakeHolderFind, "new stakeholder is not presented")
		assert.True(t, newParticipantFind, "new participant is not presented")
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

	// set up of hooks, which should be refactored since they share virtually the same code
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
