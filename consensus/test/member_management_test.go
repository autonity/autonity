package test

import (
	"context"
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/p2p/enode"
	"math/big"
	"strconv"
	"sync"
	"testing"
)

func TestMemberManagementAddNewValidator(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	once := sync.Once{}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newValidatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}

	stakeBalance := new(big.Int).SetUint64(300)

	addValidatorHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		once.Do(func() {
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(operatorKey)
			auth.From = operatorAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(30000000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			eNode := enode.V4DNSUrl(newValidatorKey.PublicKey, "VN:8527", 8527, 8527)
			_, err = instance.AddValidator(auth, crypto.PubkeyToAddress(newValidatorKey.PublicKey), stakeBalance, eNode)
			if err != nil {
				t.Fatal(err)
			}
		})
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

	addValidatorCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VA"].rpcPort))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			t.Fatal(err)
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(validators["VA"].lastBlock),
			Context:     context.Background(),
		}

		// check node presented in white list.
		eNode := enode.V4DNSUrl(newValidatorKey.PublicKey, "VN:8527", 8527, 8527)
		whiteList, err := instance.GetWhitelist(&auth)
		whiteListed := false
		for _, node := range whiteList {
			if node == eNode {
				whiteListed = true
				break
			}
		}

		// check node role and its stake balance.
		curNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			t.Fatal(err)
		}

		founded := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newValidatorKey.PublicKey) {
				founded = true
				if curNetworkMetrics.Stakes[index].Uint64() != stakeBalance.Uint64() {
					t.Fatal("new validator's stake is not expected")
				}

				if curNetworkMetrics.Usertypes[index] != 2 {
					t.Fatal("new validator's user type is not expected")
				}

				break
			}
		}

		// compare the total stake supply before and after new node added.
		auth.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64() != stakeBalance.Uint64() {
			t.Fatal("stake total supply is not expected")
		}

		if !whiteListed || !founded {
			t.Fatal("new validator is not presented")
		}
	}

	testCase := &testCase{
		name:          "member management test add validator",
		numValidators: 6,
		numBlocks:     10,
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
	once := sync.Once{}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newStakeHolderKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}

	stakeBalance := new(big.Int).SetUint64(100)

	addStakeHolderHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		once.Do(func() {
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(operatorKey)
			auth.From = operatorAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(30000000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			eNode := enode.V4DNSUrl(newStakeHolderKey.PublicKey, "SN:8527", 8527, 8527)
			_, err = instance.AddStakeholder(auth, crypto.PubkeyToAddress(newStakeHolderKey.PublicKey), eNode, stakeBalance)
			if err != nil {
				t.Fatal(err)
			}
		})
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

	addStakeHolderCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VA"].rpcPort))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			t.Fatal(err)
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(validators["VA"].lastBlock),
			Context:     context.Background(),
		}

		// check node presented in white list.
		eNode := enode.V4DNSUrl(newStakeHolderKey.PublicKey, "SN:8527", 8527, 8527)
		whiteList, err := instance.GetWhitelist(&auth)
		whiteListed := false
		for _, node := range whiteList {
			if node == eNode {
				whiteListed = true
				break
			}
		}

		// check node role and its stake balance.
		curNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			t.Fatal(err)
		}

		founded := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newStakeHolderKey.PublicKey) {
				founded = true
				if curNetworkMetrics.Stakes[index].Uint64() != stakeBalance.Uint64() {
					t.Fatal("new stakeholder's stake is not expected")
				}

				if curNetworkMetrics.Usertypes[index] != 1 {
					t.Fatal("new stakeholder's user type is not expected")
				}

				break
			}
		}

		// compare the total stake supply before and after new node added.
		auth.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64() != stakeBalance.Uint64() {
			t.Fatal("stake total supply is not expected")
		}

		if !whiteListed || !founded {
			t.Fatal("new stakeholder is not presented")
		}
	}

	testCase := &testCase{
		name:          "member management test add stakeholder",
		numValidators: 6,
		numBlocks:     10,
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
	once := sync.Once{}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)

	newParticipantKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}

	addParticipantHook := func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		once.Do(func() {
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(operatorKey)
			auth.From = operatorAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(30000000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			eNode := enode.V4DNSUrl(newParticipantKey.PublicKey, "PN:8527", 8527, 8527)
			_, err = instance.AddParticipant(auth, crypto.PubkeyToAddress(newParticipantKey.PublicKey), eNode)
			if err != nil {
				t.Fatal(err)
			}
		})
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

	addParticipantCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VA"].rpcPort))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			t.Fatal(err)
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(validators["VA"].lastBlock),
			Context:     context.Background(),
		}

		// check node presented in white list.
		eNode := enode.V4DNSUrl(newParticipantKey.PublicKey, "PN:8527", 8527, 8527)
		whiteList, err := instance.GetWhitelist(&auth)
		whiteListed := false
		for _, node := range whiteList {
			if node == eNode {
				whiteListed = true
				break
			}
		}

		// check node role and its stake balance.
		curNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if err != nil {
			t.Fatal(err)
		}

		founded := false
		for index, v := range curNetworkMetrics.Accounts {
			if v == crypto.PubkeyToAddress(newParticipantKey.PublicKey) {
				founded = true
				if curNetworkMetrics.Stakes[index].Uint64() != 0 {
					t.Fatal("new participant's stake is not expected")
				}

				if curNetworkMetrics.Usertypes[index] != 0 {
					t.Fatal("new participant's user type is not expected")
				}

				break
			}
		}

		// compare the total stake supply before and after new node added.
		auth.BlockNumber.SetUint64(3)
		initNetworkMetrics, err := instance.DumpEconomicsMetricData(&auth)
		if curNetworkMetrics.Stakesupply.Sub(curNetworkMetrics.Stakesupply, initNetworkMetrics.Stakesupply).Uint64() != 0 {
			t.Fatal("stake total supply is not expected")
		}

		if !whiteListed || !founded {
			t.Fatal("new participant is not presented")
		}
	}

	testCase := &testCase{
		name:          "member management test add participant",
		numValidators: 6,
		numBlocks:     10,
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
	once := sync.Once{}
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

	removeUserCheckerHook := func(t *testing.T, validators map[string]*testNode) {
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validators["VD"].rpcPort))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		contractAddress := autonity.ContractAddress
		instance, err := NewAutonity(contractAddress, conn)
		if err != nil {
			t.Fatal(err)
		}

		auth := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: new(big.Int).SetUint64(validators["VD"].lastBlock),
			Context:     context.Background(),
		}
		validatorsList := validators["VD"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
		isMember, err := instance.CheckMember(&auth, *validatorsList[0].Address)
		if isMember {
			t.Fatal("Wrong membership for removed user")
		}
	}

	testCase := &testCase{
		name:                 "member management test remove user",
		numValidators:        6,
		numBlocks:            10,
		txPerPeer:            1,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook:          genesisHook,
		finalAssert:          removeUserCheckerHook,
	}
	testCase.sendTransactionHooks["VD"] = func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		once.Do(func() {
			conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
			if err != nil {
				t.Fatal(err)
			}

			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			auth := bind.NewKeyedTransactor(operatorKey)
			auth.From = operatorAddress
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasLimit = uint64(30000000) // in units
			auth.GasPrice = gasPrice

			contractAddress := autonity.ContractAddress
			instance, err := NewAutonity(contractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = instance.RemoveUser(auth, *validatorsList[0].Address)
			if err != nil {
				t.Fatal(err)
			}
			testCase.removedPeers[*validatorsList[0].Address] = validator.lastBlock
		})
		return false, nil, nil
	}
	runTest(t, testCase)
}
