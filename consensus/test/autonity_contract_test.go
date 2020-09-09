package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common/acdefault"
	"github.com/clearmatics/autonity/common/graph"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/common/math"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"

	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/ethclient"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
)

const DefaultTestGasPrice = 100000000000

func TestCheckFeeRedirectionAndRedistribution(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	hookGenerator := func() (hook, hook) {
		prevBlockBalance := uint64(0)
		prevSTBalance := new(big.Int)

		fBefore := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			st, _ := validator.service.BlockChain().State()
			if block.NumberU64() == 1 && st.GetBalance(autonity.ContractAddress).Uint64() != 0 {
				return fmt.Errorf("incorrect balance on the first block")
			}
			return nil
		}
		fAfter := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			st, _ := validator.service.BlockChain().State()

			if block.NumberU64() == 1 && prevBlockBalance != 0 {
				return fmt.Errorf("incorrect balance on the first block")
			}
			contractBalance := st.GetBalance(autonity.ContractAddress)
			if block.NumberU64() > 1 && len(block.Transactions()) > 0 && block.NumberU64() <= uint64(tCase.numBlocks) {
				if contractBalance.Uint64() < prevBlockBalance {
					return fmt.Errorf("balance must be increased")
				}
			}
			prevBlockBalance = contractBalance.Uint64()

			if block.NumberU64() > 1 && len(block.Transactions()) > 0 && block.NumberU64() <= uint64(tCase.numBlocks) {
				sh := validator.service.BlockChain().Config().AutonityContractConfig.GetStakeHolderUsers()[0]
				stakeHolderBalance := st.GetBalance(*sh.Address)
				if stakeHolderBalance.Cmp(prevSTBalance) != 1 {
					return fmt.Errorf("balance must be increased")
				}
				prevSTBalance = stakeHolderBalance
			}

			return nil
		}
		return fBefore, fAfter
	}

	case1Before, case1After := hookGenerator()
	case2Before, case2After := hookGenerator()
	case3Before, case3After := hookGenerator()
	cases := []*testCase{
		{
			name:          "no malicious - 1 tx per second",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     1,
			beforeHooks: map[string]hook{
				"VD": case1Before,
			},
			afterHooks: map[string]hook{
				"VD": case1After,
			},
		},
		{
			name:          "no malicious - 10 tx per second",
			numValidators: 6,
			numBlocks:     10,
			txPerPeer:     10,
			beforeHooks: map[string]hook{
				"VF": case2Before,
			},
			afterHooks: map[string]hook{
				"VF": case2After,
			},
		},
		{
			name:          "no malicious - 5 tx per second 4 peers",
			numValidators: 4,
			numBlocks:     5,
			txPerPeer:     5,
			beforeHooks: map[string]hook{
				"VB": case3Before,
			},
			afterHooks: map[string]hook{
				"VB": case3After,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)

		})
	}
}

func TestCheckBlockWithSmallFee(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	hookGenerator := func() (hook, hook) {
		prevBlockBalance := uint64(0)
		fBefore := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			st, _ := validator.service.BlockChain().State()
			if block.NumberU64() == 1 && st.GetBalance(autonity.ContractAddress).Uint64() != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			return nil
		}
		fAfter := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			st, _ := validator.service.BlockChain().State()

			if block.NumberU64() == 1 && prevBlockBalance != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			contractBalance := st.GetBalance(autonity.ContractAddress)

			prevBlockBalance = contractBalance.Uint64()
			return nil
		}
		return fBefore, fAfter
	}

	case1Before, case1After := hookGenerator()
	cases := []*testCase{
		{
			name:          "no malicious - 1 tx per second",
			numValidators: 5,
			numBlocks:     5,
			txPerPeer:     3,
			sendTransactionHooks: map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error){
				"VD": func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error) { //nolint
					nonce := validator.service.TxPool().Nonce(fromAddr)

					tx, err := types.SignTx(
						types.NewTransaction(
							nonce,
							toAddr,
							big.NewInt(1),
							210000000,
							big.NewInt(DefaultTestGasPrice-200),
							nil,
						),
						types.HomesteadSigner{}, validator.privateKey)
					if err != nil {
						return false, nil, err
					}
					err = validator.service.TxPool().AddLocal(tx)
					if err == nil {
						return false, nil, err
					}

					//step 2 valid transaction
					tx, err = types.SignTx(
						types.NewTransaction(
							nonce,
							toAddr,
							big.NewInt(1),
							210000000,
							big.NewInt(DefaultTestGasPrice+200),
							nil,
						),
						types.HomesteadSigner{}, validator.privateKey)
					if err != nil {
						return false, nil, err
					}
					err = validator.service.TxPool().AddLocal(tx)
					if err != nil {
						return false, nil, err
					}

					return false, tx, nil
				},
			},
			beforeHooks: map[string]hook{
				"VD": case1Before,
			},
			afterHooks: map[string]hook{
				"VD": case1After,
			},
			genesisHook: func(g *core.Genesis) *core.Genesis {
				g.Config.AutonityContractConfig.MinGasPrice = DefaultTestGasPrice - 100
				return g
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestRemoveFromValidatorsList(t *testing.T) {
	// to be tracked by https://github.com/clearmatics/autonity/issues/604
	t.Skip("skipping test since the upstream update cause local e2e test framework go routine leak.")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}

	once := sync.Once{}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            10,
		txPerPeer:            1,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(100000000000000000),
			}
			return g
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			validatorUsers := validators["VE"].service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			validatorsListGenesis := []string{}
			for i := range validatorUsers {
				validatorsListGenesis = append(validatorsListGenesis, validatorUsers[i].Address.String())
			}

			stateDB, err := validators["VE"].service.BlockChain().State()
			if err != nil {
				t.Fatal(err)
			}
			validatorList, err := validators["VE"].service.BlockChain().GetAutonityContract().GetCommittee(
				validators["VE"].service.BlockChain().CurrentHeader(),
				stateDB,
			)
			if err != nil {
				t.Fatal(err)
			}
			validatorListStr := []string{}
			for _, v := range validatorList {
				validatorListStr = append(validatorListStr, v.Address.String())
			}

			if len(validatorsListGenesis) <= len(validatorListStr) {
				t.Fatal("Incorrect validator list")
			}
		},
	}
	testCase.sendTransactionHooks["VD"] = func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		skip := true
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
			auth.GasLimit = uint64(300000) // in units
			auth.GasPrice = gasPrice

			instance, err := NewAutonity(autonity.ContractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			validatorsList := validator.service.BlockChain().Config().AutonityContractConfig.GetValidatorUsers()
			_, err = instance.RemoveUser(auth, *validatorsList[0].Address)
			if err != nil {
				t.Fatal(err)
			}
			skip = false
			testCase.removedPeers[*validatorsList[0].Address] = validator.lastBlock
		})

		return skip, nil, nil
	}
	runTest(t, testCase)
}

func TestAddIncorrectStakeholdersToList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	participantKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}

	once := sync.Once{}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            10,
		txPerPeer:            1,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(100000000000000000),
			}
			return g
		},
	}
	testCase.sendTransactionHooks["VD"] = func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		skip := true
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
			auth.GasLimit = uint64(300000) // in units
			auth.GasPrice = gasPrice

			instance, err := NewAutonity(autonity.ContractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			pEnode := enode.NewV4(&participantKey.PublicKey, net.ParseIP("127.0.0.1"), 8527, 8527)
			_, err = instance.AddParticipant(auth, crypto.PubkeyToAddress(participantKey.PublicKey), pEnode.String())
			if err != nil {
				t.Fatal(err)
			}

			skip = false
		})
		return skip, nil, nil
	}
	runTest(t, testCase)
}

func TestAddStakeholderWithCorruptedEnodeToList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	participantKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}

	once := sync.Once{}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            10,
		txPerPeer:            1,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(100000000000000000),
			}
			return g
		},
	}
	testCase.sendTransactionHooks["VD"] = func(validator *testNode, _ common.Address, _ common.Address) (bool, *types.Transaction, error) { //nolint
		if validator.lastBlock <= 3 {
			return true, nil, nil
		}
		skip := true
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
			auth.GasLimit = uint64(300000) // in units
			auth.GasPrice = gasPrice

			instance, err := NewAutonity(autonity.ContractAddress, conn)
			if err != nil {
				t.Fatal(err)
			}
			_, err = instance.AddParticipant(auth, crypto.PubkeyToAddress(participantKey.PublicKey), "enode://some_bad_enode@127.0.0.1:8527")
			if err != nil {
				t.Fatal(err)
			}

			skip = false
		})
		return skip, nil, nil
	}
	runTest(t, testCase)
}

func TestContractUpgrade_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            10,
		txPerPeer:            1,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(math.MaxInt64),
			}
			return g
		},
		afterHooks: map[string]hook{
			"VD": upgradeHook(map[uint64]struct{}{
				5: {},
			},
				operatorAddress,
				operatorKey),
		},
	}
	runTest(t, testCase)
}

func TestContractUpgradeSeveralUpgrades(t *testing.T) {
	t.Skip("test is flaky - https://github.com/clearmatics/autonity/issues/496")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            20,
		txPerPeer:            10,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(math.MaxInt64),
			}
			return g
		},
		afterHooks: map[string]hook{
			"VD": upgradeHook(map[uint64]struct{}{
				5:  {},
				7:  {},
				15: {},
			},
				operatorAddress,
				operatorKey),
		},
	}
	runTest(t, testCase)
}

func TestContractUpgradeSeveralUpgradesOnBusTopology(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	topologyStr := `graph TB
    VA---VB
    VC---VB
    VD---VC
    VE---VD
`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            20,
		txPerPeer:            10,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(math.MaxInt64),
			}
			return g
		},
		afterHooks: map[string]hook{
			"VD": upgradeHook(map[uint64]struct{}{
				5:  {},
				7:  {},
				15: {},
			},
				operatorAddress,
				operatorKey),
		},
		topology: &Topology{
			graph: *topology,
		},
	}
	runTest(t, testCase)
}

func TestContractUpgradeSeveralUpgradesOnStarTopology(t *testing.T) {
	t.Skip("test is flaky - https://github.com/clearmatics/autonity/issues/496")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	operatorKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	topologyStr := `graph TB
    SF---VA
    SF---VB
    SF---VC
    SF---VD
    SF-->VE`

	topology, err := graph.Parse(strings.NewReader(topologyStr))
	if err != nil {
		t.Fatal("parse error")
	}

	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	testCase := &testCase{
		name:                 "no malicious - 1 tx per second",
		numValidators:        5,
		numBlocks:            20,
		txPerPeer:            10,
		removedPeers:         make(map[common.Address]uint64),
		sendTransactionHooks: make(map[string]func(validator *testNode, fromAddr common.Address, toAddr common.Address) (bool, *types.Transaction, error)),
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddress
			g.Alloc[operatorAddress] = core.GenesisAccount{
				Balance: big.NewInt(math.MaxInt64),
			}
			return g
		},
		afterHooks: map[string]hook{
			"VD": upgradeHook(map[uint64]struct{}{
				5:  {},
				7:  {},
				15: {},
			},
				operatorAddress,
				operatorKey),
		},
		topology: &Topology{
			graph: *topology,
		},
	}
	runTest(t, testCase)
}

func upgradeHook(upgradeBlocks map[uint64]struct{}, operatorAddress common.Address, operatorKey *ecdsa.PrivateKey) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		log.Error("Upgrade hook")
		if _, ok := upgradeBlocks[block.Number().Uint64()]; !ok {
			return nil
		}
		conn, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(validator.rpcPort))
		if err != nil {
			return err
		}
		defer conn.Close()

		nonce, err := conn.PendingNonceAt(context.Background(), operatorAddress)
		if err != nil {
			return err
		}

		gasPrice, err := conn.SuggestGasPrice(context.Background())
		if err != nil {
			return err
		}

		auth := bind.NewKeyedTransactor(operatorKey)
		auth.From = operatorAddress
		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasLimit = uint64(30000000) // in units
		auth.GasPrice = gasPrice

		instance, err := NewAutonity(autonity.ContractAddress, conn)
		if err != nil {
			return err
		}

		_, err = instance.UpgradeContract(auth, acdefault.Bytecode(), acdefault.ABI())
		if err != nil {
			return err
		}
		return nil
	}
}
