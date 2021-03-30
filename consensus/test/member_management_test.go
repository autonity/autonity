package test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	initHeight := uint64(0)
	startHeight := uint64(1)

	// prepare chain operator
	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	validatorRole := uint8(2)
	stakeHolderRole := uint8(1)
	participantRole := uint8(0)

	newValidatorNodeKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	newValidatorPubKey := newValidatorNodeKey.PublicKey
	newValidatorENode := enode.V4DNSUrl(newValidatorPubKey, "VN:8527", 8527, 8527)

	newStakeholderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	newStakeholderPubKey := newStakeholderKey.PublicKey
	newStakeholderEnode := enode.V4DNSUrl(newStakeholderPubKey, "SN:8528", 8528, 8528)

	newParticipantKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	newParticipantPubKey := newParticipantKey.PublicKey
	newParticipantEnode := enode.V4DNSUrl(newParticipantPubKey, "PN:8529", 8529, 8529)

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

	addUsersHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == startHeight {
			port := validator.rpcPort

			err := interact(port).tx(operatorKey).addUser(crypto.PubkeyToAddress(newValidatorPubKey), validatorStake, newValidatorENode, validatorRole)
			if err != nil {
				return false, nil, err
			}

			err = interact(port).tx(operatorKey).addUser(crypto.PubkeyToAddress(newStakeholderPubKey), stakeHolderStake, newStakeholderEnode, stakeHolderRole)
			if err != nil {
				return false, nil, err
			}

			err = interact(port).tx(operatorKey).addUser(crypto.PubkeyToAddress(newParticipantPubKey), participantStake, newParticipantEnode, participantRole)
			if err != nil {
				return false, nil, err
			}
		}
		return false, nil, nil
	}

	isNetworkParticipant := func(port int, height uint64, address common.Address, eNode string) bool {
		whiteList, err := interact(port).call(height).getWhitelist()
		require.NoError(t, err)
		var inWhitelist bool
		for _, en := range whiteList {
			if en == eNode {
				inWhitelist = true
				break
			}
		}
		user, err := interact(port).call(height).getUser(address)
		require.NoError(t, err)
		isMember := user.Addr == address
		assert.Equal(t, isMember, inWhitelist)
		return isMember
	}

	// to check user membership, user type, stake balance.
	validateAddedUser := func(t *testing.T, port int, height uint64, address common.Address, eNode string, role uint8, stake uint64, economicMetric AutonityEconomicMetrics) {
		assert.True(t, isNetworkParticipant(port, height, address, eNode), "wrong membership for added user")
		// check validator and stakeholder's stake balance
		userBalance, err := interact(port).call(height).getAccountStake(address)
		require.NoError(t, err)
		user, err := interact(port).call(height).getUser(address)
		require.NoError(t, err)
		require.Equal(t, user.Stake.Uint64(), userBalance.Uint64())
		require.Equal(t, user.UserType, role)
		require.Equal(t, user.Enode, eNode)
		require.Equal(t, user.Stake.Uint64(), stake)

		for index, v := range economicMetric.Accounts {
			if v == address {
				assert.Equal(t, role, economicMetric.Usertypes[index], "user type is not expected")
				break
			}
		}
	}

	addUsersCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		port := validators["VA"].rpcPort
		lastHeight := validators["VA"].lastBlock
		curNetworkMetrics, err := interact(port).call(lastHeight).dumpEconomicsMetricData()
		require.NoError(t, err)

		validateAddedUser(t, port, lastHeight, crypto.PubkeyToAddress(newValidatorPubKey), newValidatorENode, validatorRole, validatorStake.Uint64(), curNetworkMetrics)
		validateAddedUser(t, port, lastHeight, crypto.PubkeyToAddress(newStakeholderPubKey), newStakeholderEnode, stakeHolderRole, stakeHolderStake.Uint64(), curNetworkMetrics)
		validateAddedUser(t, port, lastHeight, crypto.PubkeyToAddress(newParticipantPubKey), newParticipantEnode, participantRole, participantStake.Uint64(), curNetworkMetrics)

		// compare the total stake supply before and after new node added.
		initNetworkMetrics, err := interact(validators["VA"].rpcPort).call(initHeight).dumpEconomicsMetricData()
		require.NoError(t, err)

		// new_total_stake - init_total_stake == new added (validatorStake + stakeHolderStake + participantStake)
		b := curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64()
		assert.Equal(t, b, validatorStake.Uint64()+stakeHolderStake.Uint64()+participantStake.Uint64(), "stake total supply is not expected")
	}

	removeUserHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == startHeight {
			return true, nil, interact(validator.rpcPort).tx(operatorKey).removeUser(addressToRemove)
		}
		return false, nil, nil
	}

	removeUserCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		port := validators["VA"].rpcPort
		lastHeight := validators["VA"].lastBlock
		assert.False(t, isNetworkParticipant(port, lastHeight, addressToRemove, eNodeToRemove), "wrong membership for removed user")
	}

	// numBlocks are used to stop the test on current test framework, to let user management TX to be mined before the test end,
	// bigger numBlocks in below test cases are set.
	cases := []*testCase{
		{
			name:          "add users",
			numValidators: 6,
			numBlocks:     15,
			txPerPeer:     1,
			sendTransactionHooks: map[string]sendTransactionHook{
				"VA": addUsersHook,
			},
			genesisHook: genesisHook,
			finalAssert: addUsersCheckerHook,
		},
		{
			name:          "remove user",
			numValidators: 6,
			numBlocks:     15,
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
