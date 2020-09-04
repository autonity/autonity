package test

import (
	"context"
	"fmt"
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"strconv"
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

func TestMemberManagementAddNewValidator(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	newValidatorKey, err := keygenerator.Next()
	require.NoError(t, err)
	newValidatorPubKey := newValidatorKey.PublicKey
	eNode := enode.V4DNSUrl(newValidatorPubKey, "VN:8527", 8527, 8527)

	stakeBalance := new(big.Int).SetUint64(300)

	addValidatorHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()
			txOpt, err := contract.transactionOpts(operatorKey)
			if err != nil {
				return true, nil, err
			}
			_, err = contract.AddValidator(txOpt, crypto.PubkeyToAddress(newValidatorPubKey), stakeBalance, eNode)
			if err != nil {
				return true, nil, err
			}
		}
		return false, nil, nil
	}

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	addValidatorCheckerHook := func(t *testing.T, validators map[string]*testNode) error {
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VA"].rpcPort))
		if err != nil {
			return err
		}
		defer conn.Close()

		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			return err
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(validators["VA"].lastBlock),
			Context:     context.Background(),
		}

		// check node presented in white list.
		whiteList, err := instance.GetWhitelist(&auth)
		if err != nil {
			return err
		}

		assert.Contains(t, whiteList, eNode, "eNode is not presented from member list")

		// check node role and its stake balance.
		curNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			return err
		}

		found := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newValidatorPubKey) {
				found = true
				assert.Equal(t, curNetworkMetrics.Stakes[index].Uint64(), stakeBalance.Uint64(), "new validator's stake is not expected")
				assert.Equal(t, int(curNetworkMetrics.Usertypes[index]), 2, "new validator's user type is not expected")
				break
			}
		}

		// compare the total stake supply before and after new node added.
		auth.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			return err
		}

		b := curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64()
		assert.Equal(t, b, stakeBalance.Uint64(), "stake total supply is not expected")
		assert.True(t, found, "new validator is not presented")
		return nil
	}

	testCase := &testCase{
		name:          "member management test add validator",
		numValidators: 6,
		numBlocks:     20,
		txPerPeer:     1,
		sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
			"VA": addValidatorHook,
		},
		genesisHook: genesisHook,
		finalAssert: addValidatorCheckerHook,
	}
	runTest(t, testCase)
}

func TestMemberManagementAddNewStakeHolder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newStakeHolderKey, err := keygenerator.Next()
	require.NoError(t, err)

	eNode := enode.V4DNSUrl(newStakeHolderKey.PublicKey, "SN:8527", 8527, 8527)
	stakeBalance := new(big.Int).SetUint64(100)

	addStakeHolderHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()

			txOpt, err := contract.transactionOpts(operatorKey)
			if err != nil {
				return true, nil, err
			}

			_, err = contract.AddStakeholder(txOpt, crypto.PubkeyToAddress(newStakeHolderKey.PublicKey), eNode, stakeBalance)
			if err != nil {
				return true, nil, err
			}
		}
		return false, nil, nil
	}

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	addStakeHolderCheckerHook := func(t *testing.T, validators map[string]*testNode) error {
		contract, err := autonityInstance(validators["VA"].rpcPort)
		if err != nil {
			return err
		}
		defer contract.Close()
		callOpt := contract.callOpts(validators["VA"].lastBlock)

		// check node presented in white list.
		whiteList, err := contract.GetWhitelist(callOpt)
		if err != nil {
			return err
		}

		assert.Contains(t, whiteList, eNode)

		// check node role and its stake balance.
		curNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		found := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newStakeHolderKey.PublicKey) {
				found = true
				assert.Equal(t, curNetworkMetrics.Stakes[index].Uint64(), stakeBalance.Uint64(), "new stakeholder's stake is not expected")
				assert.Equal(t, int(curNetworkMetrics.Usertypes[index]), 1, "new stakeholder's user type is not expected")
				break
			}
		}

		// compare the total stake supply before and after new node added.
		callOpt.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		b := curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64()
		assert.Equal(t, b, stakeBalance.Uint64(), "stake total supply is not expected")
		assert.True(t, found, "new stakeholder is not presented")
		return nil
	}

	testCase := &testCase{
		name:          "member management test add stakeholder",
		numValidators: 6,
		numBlocks:     20,
		txPerPeer:     1,
		sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
			"VA": addStakeHolderHook,
		},
		genesisHook: genesisHook,
		finalAssert: addStakeHolderCheckerHook,
	}
	runTest(t, testCase)
}

func TestMemberManagementAddNewParticipant(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newParticipantKey, err := keygenerator.Next()
	require.NoError(t, err)
	eNode := enode.V4DNSUrl(newParticipantKey.PublicKey, "PN:8527", 8527, 8527)

	addParticipantHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()
			txOpt, err := contract.transactionOpts(operatorKey)
			if err != nil {
				return true, nil, err
			}

			_, err = contract.AddParticipant(txOpt, crypto.PubkeyToAddress(newParticipantKey.PublicKey), eNode)
			if err != nil {
				return true, nil, err
			}
		}
		return false, nil, nil
	}

	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	addParticipantCheckerHook := func(t *testing.T, validators map[string]*testNode) error {
		contract, err := autonityInstance(validators["VA"].rpcPort)
		if err != nil {
			return err
		}

		callOpt := contract.callOpts(validators["VA"].lastBlock)

		whiteList, err := contract.GetWhitelist(callOpt)
		if err != nil {
			return err
		}

		assert.Contains(t, whiteList, eNode)

		// check node role and its stake balance.
		curNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		found := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newParticipantKey.PublicKey) {
				found = true
				assert.Equal(t, curNetworkMetrics.Stakes[index].Uint64(), uint64(0), "new participant's stake is not expected")
				assert.Equal(t, int(curNetworkMetrics.Usertypes[index]), 0, "new participant's user type is not expected")
				break
			}
		}

		// compare the total stake supply before and after new node added.
		callOpt.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := contract.DumpEconomicsMetricData(callOpt)
		if err != nil {
			return err
		}

		assert.Zero(t, curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64(), "stake total supply is not expected")
		assert.True(t, found, "new participant is not presented")
		return nil
	}

	testCase := &testCase{
		name:          "member management test add participant",
		numValidators: 6,
		numBlocks:     20,
		txPerPeer:     1,
		sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
			"VA": addParticipantHook,
		},
		genesisHook: genesisHook,
		finalAssert: addParticipantCheckerHook,
	}
	runTest(t, testCase)
}

func TestMemberManagementRemoveUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	// genesis hook
	genesisHook := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		return g
	}

	removeUserCheckerHook := func(t *testing.T, validators map[string]*testNode) error {
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VD"].rpcPort))
		if err != nil {
			return err
		}
		defer conn.Close()

		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			return err
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(validators["VD"].lastBlock),
			Context:     context.Background(),
		}
		validatorsList := validators["VD"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		isMember, err := instance.CheckMember(&auth, *validatorsList[0].Address)
		if err != nil {
			return err
		}
		if isMember {
			return fmt.Errorf("wrong membership for removed user")
		}
		return nil
	}

	testCase := &testCase{
		name:                 "member management test remove user",
		numValidators:        6,
		numBlocks:            20,
		txPerPeer:            1,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook:          genesisHook,
		finalAssert:          removeUserCheckerHook,
	}
	testCase.sendTransactionHooks["VD"] = func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock == 4 {
			contract, err := autonityInstance(validator.rpcPort)
			if err != nil {
				return true, nil, err
			}
			defer contract.Close()
			txOpt, err := contract.transactionOpts(operatorKey)
			if err != nil {
				return true, nil, err
			}
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = contract.RemoveUser(txOpt, *validatorsList[0].Address)
			if err != nil {
				return true, nil, err
			}
			testCase.removedPeers[*validatorsList[0].Address] = validator.lastBlock
		}
		return false, nil, nil
	}
	runTest(t, testCase)
}
