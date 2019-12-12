package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/clearmatics/autonity/accounts"
	"github.com/clearmatics/autonity/accounts/keystore"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/fdlimit"
	"github.com/clearmatics/autonity/consensus"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"
)

const DefaultTestGasPrice = 100000000000

func TestTendermintSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintOneMalicious(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	addedValidatorsBlocks := make(map[common.Hash]uint64)
	removedValidatorsBlocks := make(map[common.Hash]uint64)
	changedValidators := tendermintCore.NewChanges()

	cases := []*testCase{
		{
			name:      "replace a valid validator with invalid one",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			maliciousPeers: map[int]injectors{
				4: {
					cons: func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewReplaceValidatorCore(basic, changedValidators)
					},
					backs: func(basic tendermintCore.Backend) tendermintCore.Backend {
						return tendermintCore.NewChangeValidatorBackend(t, basic, changedValidators, removedValidatorsBlocks, addedValidatorsBlocks)
					},
				},
			},
			addedValidatorsBlocks:   addedValidatorsBlocks,
			removedValidatorsBlocks: removedValidatorsBlocks,
			changedValidators:       changedValidators,
		},
		{
			name:      "add a validator",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			maliciousPeers: map[int]injectors{
				4: {
					cons: func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewAddValidatorCore(basic, changedValidators)
					},
					backs: func(basic tendermintCore.Backend) tendermintCore.Backend {
						return tendermintCore.NewChangeValidatorBackend(t, basic, changedValidators, removedValidatorsBlocks, addedValidatorsBlocks)
					},
				},
			},
			addedValidatorsBlocks:   addedValidatorsBlocks,
			removedValidatorsBlocks: removedValidatorsBlocks,
			changedValidators:       changedValidators,
		},
		{
			name:      "remove a validator",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			maliciousPeers: map[int]injectors{
				4: {
					cons: func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewRemoveValidatorCore(basic, changedValidators)
					},
					backs: func(basic tendermintCore.Backend) tendermintCore.Backend {
						return tendermintCore.NewChangeValidatorBackend(t, basic, changedValidators, removedValidatorsBlocks, addedValidatorsBlocks)
					},
				},
			},
			addedValidatorsBlocks:   addedValidatorsBlocks,
			removedValidatorsBlocks: removedValidatorsBlocks,
			changedValidators:       changedValidators,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintSlowConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious, one slow node",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			networkRates: map[int]networkRate{
				4: {50 * 1024, 50 * 1024},
			},
		},
		{
			name:      "no malicious, all nodes are slow",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			networkRates: map[int]networkRate{
				0: {50 * 1024, 50 * 1024},
				1: {50 * 1024, 50 * 1024},
				2: {50 * 1024, 50 * 1024},
				3: {50 * 1024, 50 * 1024},
				4: {50 * 1024, 50 * 1024},
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


func TestTendermintMemoryLeak(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious - 30 tx per second",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 30,
		},
		{
			name:      "no malicious - 100 blocks",
			numPeers:  7,
			numBlocks: 100,
			txPerPeer: 10,
		},
	}

	const repeats = 10
	leaks := make([][]float64, len(cases))

	for i, testCase := range cases {
		i := i
		testCase := testCase

		leaks[i] = make([]float64, repeats)
		for n:=0;n<repeats; n++ {
			n:=n
			t.Run(fmt.Sprintf("test case %s, try %d", testCase.name, n), func(t *testing.T) {
				m := leak.MarkMemory()
				runTest(t, testCase)
				leaks[i][n] = float64(m.Release())
			})
		}
	}

	type stats struct {
		mean float64
		std float64
		stdErr float64
		n int
	}

	leaksStats := make([]stats, len(leaks))
	for i, cases := range leaks {
		mean, std := stat.MeanStdDev(cases, nil)
		stdErr := stat.StdErr(std, float64(len(cases)))

		leaksStats[i] = stats{
			mean:   mean,
			std:    std,
			stdErr: stdErr,
			n: len(cases),
		}
	}

	sed := gmath.Sqrt(gmath.Pow(leaksStats[0].stdErr, 2) + gmath.Pow(leaksStats[1].stdErr, 2))
	tstat := (leaksStats[0].mean - leaksStats[1].mean) / sed

	// degrees of freedom
	df := leaksStats[0].n + leaksStats[1].n - 2

	// calculate the critical value
	alpha := 0.05
	cv := stat.PP.ppf(1.0 - alpha, df)

	// calculate the p-value
	p := (1 - stat.CDF(gmath.Abs(tstat), stat.Empirical, leaks[0], nil)) * 2

	meansEqual := gmath.Abs(tstat) <= cv && p > alpha

	fmt.Println("!!!", meansEqual, gmath.Abs(tstat) <= cv, p > alpha)
}

func runTestWithLeakChecks(t *testing.T, test *testCase, blocksMap ...map[common.Hash]uint64) {
	m := leak.MarkMemory()
	defer func() {
		leaks := m.Release()

		if leaks > 0 {
			log.Error("some code is leaking", "size", leaks)
		}
	}()

	runTest(t, test, blocksMap...)
}

func TestTendermintLongRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "no malicious - 30 tx per second",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 30,
		},
		{
			name:      "no malicious - 100 blocks",
			numPeers:  5,
			numBlocks: 100,
			txPerPeer: 5,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStopUpToFNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "one node stops at block 1",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 1),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				4: {},
			},
		},
		{
			name:      "one node stops at block 5",
			numPeers:  5,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				4: {},
			},
		},
		{
			name:      "F nodes stop at block 1",
			numPeers:  7,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 1),
				4: hookStopNode(4, 1),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				3: {},
				4: {},
			},
		},
		{
			name:      "F nodes stop at block 5",
			numPeers:  7,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				3: {},
				4: {},
			},
		},
		{
			name:      "F nodes stop at blocks 4,5",
			numPeers:  7,
			numBlocks: 10,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 4),
				4: hookStopNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
			maliciousPeers: map[int]injectors{
				3: {},
				4: {},
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

func TestCheckFeeRedirectionAndRedistribution(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	hookGenerator := func() (hook, hook) {
		prevBlockBalance := uint64(0)
		prevSTBalance := new(big.Int)

		fBefore := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			addr, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				return err
			}
			st, _ := validator.service.BlockChain().State()
			if block.NumberU64() == 1 && st.GetBalance(addr).Uint64() != 0 {
				return fmt.Errorf("incorrect balance on the first block")
			}
			return nil
		}
		fAfter := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			autonityContractAddress, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				return err
			}
			st, _ := validator.service.BlockChain().State()

			if block.NumberU64() == 1 && prevBlockBalance != 0 {
				return fmt.Errorf("incorrect balance on the first block")
			}
			contractBalance := st.GetBalance(autonityContractAddress)
			if block.NumberU64() > 1 && len(block.Transactions()) > 0 && block.NumberU64() <= uint64(tCase.numBlocks) {
				if contractBalance.Uint64() < prevBlockBalance {
					return fmt.Errorf("balance must be increased")
				}
			}
			prevBlockBalance = contractBalance.Uint64()

			if block.NumberU64() > 1 && len(block.Transactions()) > 0 && block.NumberU64() <= uint64(tCase.numBlocks) {
				sh := validator.service.BlockChain().Config().AutonityContractConfig.GetStakeHolderUsers()[0]
				stakeHolderBalance := st.GetBalance(sh.Address)
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
			name:      "no malicious - 1 tx per second",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: case1Before,
			},
			afterHooks: map[int]hook{
				3: case1After,
			},
		},
		{
			name:      "no malicious - 10 tx per second",
			numPeers:  6,
			numBlocks: 10,
			txPerPeer: 10,
			beforeHooks: map[int]hook{
				5: case2Before,
			},
			afterHooks: map[int]hook{
				5: case2After,
			},
		},
		{
			name:      "no malicious - 5 tx per second 4 peers",
			numPeers:  4,
			numBlocks: 5,
			txPerPeer: 5,
			beforeHooks: map[int]hook{
				1: case3Before,
			},
			afterHooks: map[int]hook{
				1: case3After,
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
			addr, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				t.Fatal(err)
			}
			st, _ := validator.service.BlockChain().State()
			if block.NumberU64() == 1 && st.GetBalance(addr).Uint64() != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			return nil
		}
		fAfter := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			autonityContractAddress, err := validator.service.BlockChain().Config().AutonityContractConfig.GetContractAddress()
			if err != nil {
				t.Fatal(err)
			}
			st, _ := validator.service.BlockChain().State()

			if block.NumberU64() == 1 && prevBlockBalance != 0 {
				t.Fatal("incorrect balance on the first block")
			}
			contractBalance := st.GetBalance(autonityContractAddress)

			prevBlockBalance = contractBalance.Uint64()
			return nil
		}
		return fBefore, fAfter
	}

	case1Before, case1After := hookGenerator()
	cases := []*testCase{
		{
			name:      "no malicious - 1 tx per second",
			numPeers:  5,
			numBlocks: 5,
			txPerPeer: 3,
			sendTransactionHooks: map[int]func(service *eth.Ethereum, key *ecdsa.PrivateKey, fromAddr common.Address, toAddr common.Address) (*types.Transaction, error){
				3: func(service *eth.Ethereum, key *ecdsa.PrivateKey, fromAddr common.Address, toAddr common.Address) (*types.Transaction, error) {
					nonce := service.TxPool().Nonce(fromAddr)

					//step 1 invalid transaction. It must return error.
					tx, err := types.SignTx(
						types.NewTransaction(
							nonce,
							toAddr,
							big.NewInt(1),
							210000000,
							big.NewInt(DefaultTestGasPrice-200),
							nil,
						),
						types.HomesteadSigner{}, key)
					if err != nil {
						return nil, err
					}
					err = service.TxPool().AddLocal(tx)
					if err == nil {
						return nil, err
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
						types.HomesteadSigner{}, key)
					if err != nil {
						return nil, err
					}
					err = service.TxPool().AddLocal(tx)
					if err != nil {
						return nil, err
					}

					return tx, nil
				},
			},
			beforeHooks: map[int]hook{
				3: case1Before,
			},
			afterHooks: map[int]hook{
				3: case1After,
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

func TestTendermintStartStopSingleNode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "one node stops for 5 seconds",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "one node stops for 10 seconds",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "one node stops for 20 seconds",
			numPeers:  5,
			numBlocks: 30,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopFNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "f nodes stop for 5 seconds at the same block",
			numPeers:  7,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 5 seconds at different blocks",
			numPeers:  7,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at the same block",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at different blocks",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at the same block",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f nodes stop for 10 seconds at different blocks",
			numPeers:  7,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopFPlusOneNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "f+1 nodes stop for 5 seconds at the same block",
			numPeers:  5,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 5 seconds at different blocks",
			numPeers:  5,
			numBlocks: 15,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 10 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 10 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 20 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+1 nodes stop for 20 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 6),
			},
			afterHooks: map[int]hook{
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopFPlusTwoNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "f+2 nodes stop for 5 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 5),
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 5 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 5),
				3: hookStartNode(3, 5),
				4: hookStartNode(4, 5),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 10 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 10),
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 10 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 10),
				3: hookStartNode(3, 10),
				4: hookStartNode(4, 10),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 20 seconds at the same block",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 20),
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "f+2 nodes stop for 20 seconds at different blocks",
			numPeers:  5,
			numBlocks: 20,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				2: hookStopNode(2, 4),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				2: hookStartNode(2, 20),
				3: hookStartNode(3, 20),
				4: hookStartNode(4, 20),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintStartStopAllNodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:      "all nodes stop for 60 seconds at different blocks(2+2+1)",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 7),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 60),
				1: hookStartNode(1, 60),
				2: hookStartNode(2, 60),
				3: hookStartNode(3, 60),
				4: hookStartNode(4, 60),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "all nodes stop for 60 seconds at different blocks (2+3)",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 5),
				3: hookStopNode(3, 5),
				4: hookStopNode(4, 5),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 60),
				1: hookStartNode(1, 60),
				2: hookStartNode(2, 60),
				3: hookStartNode(3, 60),
				4: hookStartNode(4, 60),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "all nodes stop for 30 seconds at the same block",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 3),
				3: hookStopNode(3, 3),
				4: hookStopNode(4, 3),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 30),
				1: hookStartNode(1, 30),
				2: hookStartNode(2, 30),
				3: hookStartNode(3, 30),
				4: hookStartNode(4, 30),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:      "all nodes stop for 60 seconds at the same block",
			numPeers:  5,
			numBlocks: 50,
			txPerPeer: 1,
			beforeHooks: map[int]hook{
				0: hookStopNode(0, 3),
				1: hookStopNode(1, 3),
				2: hookStopNode(2, 3),
				3: hookStopNode(3, 3),
				4: hookStopNode(4, 3),
			},
			afterHooks: map[int]hook{
				0: hookStartNode(0, 60),
				1: hookStartNode(1, 60),
				2: hookStartNode(2, 60),
				3: hookStartNode(3, 60),
				4: hookStartNode(4, 60),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

func TestTendermintTC7(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	test := &testCase{
		name:      "3 nodes stop, 1 recover and sync blocks and state",
		numPeers:  6,
		numBlocks: 40,
		txPerPeer: 1,
		beforeHooks: map[int]hook{
			3: hookStopNode(3, 10),
			4: hookStopNode(4, 15),
			5: hookStopNode(5, 20),
		},
		afterHooks: map[int]hook{
			3: hookStartNode(3, 40),
		},
		maliciousPeers: map[int]injectors{
			4: {},
			5: {},
		},
		stopTime: make(map[int]time.Time),
	}

	for i := 0; i < 20; i++ {
		t.Run(fmt.Sprintf("test case %s - %d", test.name, i), func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func TestTendermintNoQuorum(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []*testCase{
		{
			name:               "2 validators, one goes down after block 3",
			numPeers:           2,
			numBlocks:          5,
			txPerPeer:          1,
			noQuorumAfterBlock: 3,
			beforeHooks: map[int]hook{
				1: hookForceStopNode(1, 3),
			},
			stopTime: make(map[int]time.Time),
		},
		{
			name:               "3 validators, two go down after block 3",
			numPeers:           3,
			numBlocks:          5,
			txPerPeer:          1,
			noQuorumAfterBlock: 3,
			noQuorumTimeout:    time.Second * 3,
			beforeHooks: map[int]hook{
				1: hookForceStopNode(1, 3),
				2: hookForceStopNode(2, 3),
			},
			stopTime: make(map[int]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

type testCase struct {
	name                   string
	isSkipped              bool
	numPeers               int
	numBlocks              int
	txPerPeer              int
	validatorsCanBeStopped *int64

	maliciousPeers          map[int]injectors
	addedValidatorsBlocks   map[common.Hash]uint64
	removedValidatorsBlocks map[common.Hash]uint64 //nolint: unused, structcheck
	changedValidators       tendermintCore.Changes //nolint: unused,structcheck

	networkRates         map[int]networkRate //map[validatorIndex]networkRate
	beforeHooks          map[int]hook        //map[validatorIndex]beforeHook
	afterHooks           map[int]hook        //map[validatorIndex]afterHook
	sendTransactionHooks map[int]func(service *eth.Ethereum, key *ecdsa.PrivateKey, fromAddr common.Address, toAddr common.Address) (*types.Transaction, error)
	stopTime             map[int]time.Time
	genesisHook          func(g *core.Genesis) *core.Genesis
	mu                   sync.RWMutex
	noQuorumAfterBlock   uint64
	noQuorumTimeout      time.Duration
}

type injectors struct {
	cons  func(basic consensus.Engine) consensus.Engine
	backs func(basic tendermintCore.Backend) tendermintCore.Backend
}

func (test *testCase) getBeforeHook(index int) hook {
	test.mu.Lock()
	defer test.mu.Unlock()

	if test.beforeHooks == nil {
		return nil
	}

	validatorHook, ok := test.beforeHooks[index]
	if !ok || validatorHook == nil {
		return nil
	}

	return validatorHook
}

func (test *testCase) getAfterHook(index int) hook {
	test.mu.Lock()
	defer test.mu.Unlock()

	if test.afterHooks == nil {
		return nil
	}

	validatorHook, ok := test.afterHooks[index]
	if !ok || validatorHook == nil {
		return nil
	}

	return validatorHook
}

func (test *testCase) setStopTime(index int, stopTime time.Time) {
	test.mu.Lock()
	test.stopTime[index] = stopTime
	test.mu.Unlock()
}

func (test *testCase) getStopTime(index int) time.Time {
	test.mu.RLock()
	currentTime := test.stopTime[index]
	test.mu.RUnlock()

	return currentTime
}

type hook func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error

func runTest(t *testing.T, test *testCase) {
	if test.isSkipped {
		t.SkipNow()
	}

	defer goleak.VerifyNone(t)

	log.Root().SetHandler(log.LvlFilterHandler(log.LvlError, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	_, err := fdlimit.Raise(512 * uint64(test.numPeers))
	if err != nil {
		t.Log("can't rise file description limit. errors are possible")
	}

	// Generate a batch of accounts to seal and fund with
	validators := make([]*testNode, test.numPeers)

	for i := range validators {
		validators[i] = new(testNode)
		validators[i].privateKey, err = crypto.GenerateKey()
		if err != nil {
			t.Fatal("cant make pk", err)
		}
	}

	for i := range validators {
		validators[i].listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}
	}

	for _, validator := range validators {
		listener := validator.listener
		validator.address = listener.Addr().String()
		port := strings.Split(listener.Addr().String(), ":")[1]
		validator.port, _ = strconv.Atoi(port)

		validator.url = enode.V4URL(
			validator.privateKey.PublicKey,
			net.IPv4(127, 0, 0, 1),
			validator.port,
			validator.port,
		)
	}

	genesis := makeGenesis(validators)
	if test.genesisHook != nil {
		genesis = test.genesisHook(genesis)
	}
	for i, validator := range validators {
		var engineConstructor func(basic consensus.Engine) consensus.Engine
		var backendConstructor func(basic tendermintCore.Backend) tendermintCore.Backend
		if test.maliciousPeers != nil {
			engineConstructor = test.maliciousPeers[i].cons
			backendConstructor = test.maliciousPeers[i].backs
		}

		validator.listener.Close()

		rates := test.networkRates[i]

		validator.node, err = makeValidator(genesis, validator.privateKey, validator.address, rates.in, rates.out, engineConstructor, backendConstructor) //делаем валидатор. Он просит переменную, которая типа функция engineConstructor
		if err != nil {
			t.Fatal("cant make a validator", i, err)
		}
	}

	wg := &errgroup.Group{}
	for _, validator := range validators {
		validator := validator

		wg.Go(func() error {
			return validator.startNode()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		wgClose := &errgroup.Group{}
		for _, validator := range validators {
			validatorInner := validator
			wgClose.Go(func() error {
				if !validatorInner.isRunning {
					return nil
				}

				errInner := validatorInner.node.Close()
				if errInner != nil {
					return fmt.Errorf("error on node close %v", err)
				}

				validatorInner.node.Wait()

				return nil
			})
		}

		err = wgClose.Wait()
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second) //level DB needs a second to close
	}()

	wg = &errgroup.Group{}
	for _, validator := range validators {
		validator := validator

		wg.Go(func() error {
			return validator.startService()
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	wg = &errgroup.Group{}
	for i, validator := range validators {
		validator := validator
		i := i

		wg.Go(func() error {
			log.Debug("peers", "i", i,
				"peers", len(validator.node.Server().Peers()),
				"nodes", len(validators))
			return nil
		})
	}
	err = wg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, validator := range validators {
			validator.subscription.Unsubscribe()
		}
	}()

	// each peer sends one tx per block
	sendTransactions(t, test, validators, test.txPerPeer, true)

	if len(test.maliciousPeers) != 0 {
		maliciousTest(t, test, validators)
	}
}

func maliciousTest(t *testing.T, test *testCase, validators []*testNode) {
	for index, validator := range validators {
		for number, block := range validator.blocks {
			if test.addedValidatorsBlocks != nil {
				if maliciousBlock, ok := test.addedValidatorsBlocks[block.hash]; ok {
					t.Errorf("a malicious block %d(%v)\nwas added to %d(%v)", number, maliciousBlock, index, validator)
				}
			}
		}
	}
}

func (validator *testNode) startNode() error {
	// Inject the signer key and start sealing with it
	store := validator.node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	var (
		err    error
		signer accounts.Account
	)

	if !validator.isInited {
		signer, err = store.ImportECDSA(validator.privateKey, "")
		if err != nil {
			return fmt.Errorf("import pk: %s", err)
		}

		for {
			// wait until the private key is imported
			_, err = validator.node.AccountManager().Find(signer)
			if err == nil {
				break
			}
			time.Sleep(50 * time.Microsecond)
		}

		validator.isInited = true
	} else {
		signer = store.Accounts()[0]
	}

	if err = store.Unlock(signer, ""); err != nil {
		return fmt.Errorf("cant unlock: %s", err)
	}

	validator.node.ResetEventMux()

	if err := validator.node.Start(); err != nil {
		return fmt.Errorf("cannot start a node %s", err)
	}

	// Start tracking the node and it's enode
	validator.enode = validator.node.Server().Self()
	return nil
}

func (validator *testNode) stopNode() error {
	//remove pending transactions
	addr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		pendingTxsMap, queuedTxsMap := validator.service.TxPool().Content()
		if len(pendingTxsMap) == 0 && len(queuedTxsMap) == 0 {
			break
		}

		canBreak := true
		for txAddr, txs := range pendingTxsMap {
			if addr != txAddr {
				continue
			}
			if len(txs) != 0 {
				canBreak = false
			}
		}
		for txAddr, txs := range queuedTxsMap {
			if addr != txAddr {
				continue
			}
			if len(txs) != 0 {
				canBreak = false
			}
		}
		if canBreak {
			break
		}
	}

	return validator.forceStopNode()
}

func (validator *testNode) forceStopNode() error {
	if err := validator.node.Stop(); err != nil {
		return fmt.Errorf("cannot stop a node on block %d: %q", validator.lastBlock, err)
	}
	validator.node.Wait()
	validator.isRunning = false
	validator.wasStopped = true

	return nil
}

func (validator *testNode) startService() error {
	var ethereum *eth.Ethereum
	if err := validator.node.Service(&ethereum); err != nil {
		return fmt.Errorf("cant start a node %s", err)
	}

	time.Sleep(100 * time.Millisecond)

	validator.service = ethereum

	if validator.eventChan == nil {
		validator.eventChan = make(chan core.ChainEvent, 1024)
		validator.transactions = make(map[common.Hash]struct{})
		validator.blocks = make(map[uint64]block)
		validator.txsSendCount = new(int64)
		validator.txsChainCount = make(map[uint64]int64)
	} else {
		// validator is restarting
		// we need to retrieve missed block events since last stop as we're not subscribing fast enough
		curBlock := validator.service.BlockChain().CurrentBlock().Number().Uint64()
		for blockNum := validator.lastBlock + 1; blockNum <= curBlock; blockNum++ {
			block := validator.service.BlockChain().GetBlockByNumber(blockNum)
			event := core.ChainEvent{
				Block: block,
				Hash:  block.Hash(),
				Logs:  nil,
			}
			validator.eventChan <- event
		}
	}

	validator.subscription = validator.service.BlockChain().SubscribeChainEvent(validator.eventChan)

	if err := ethereum.StartMining(1); err != nil {
		return fmt.Errorf("cant start mining %s", err)
	}

	for !ethereum.IsMining() {
		time.Sleep(50 * time.Millisecond)
	}

	validator.isRunning = true

	return nil
}

func sendTransactions(t *testing.T, test *testCase, validators []*testNode, txPerPeer int, errorOnTx bool) {
	const blocksToWait = 15

	txs := make(map[uint64]int) // blockNumber to count
	txsMu := sync.Mutex{}

	test.validatorsCanBeStopped = new(int64)
	wg, ctx := errgroup.WithContext(context.Background())

	for index, validator := range validators {
		index := index
		validator := validator

		logger := log.New("addr", crypto.PubkeyToAddress(validator.privateKey.PublicKey).String(), "idx", index)

		// skip malicious nodes
		if len(test.maliciousPeers) != 0 {
			if _, ok := test.maliciousPeers[index]; ok {
				atomic.AddInt64(test.validatorsCanBeStopped, 1)
				continue
			}
		}

		wg.Go(func() error {
			var err error
			testCanBeStopped := new(uint32)
			fromAddr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)

		wgLoop:
			for {
				select {
				case ev := <-validator.eventChan:
					if _, ok := validator.blocks[ev.Block.NumberU64()]; ok {
						continue
					}

					// before hook
					err = runHook(test.getBeforeHook(index), test, ev.Block, validator, index)
					if err != nil {
						return err
					}

					validator.blocks[ev.Block.NumberU64()] = block{ev.Block.Hash(), len(ev.Block.Transactions())}
					validator.lastBlock = ev.Block.NumberU64()

					logger.Error("last mined block", "validator", index,
						"num", validator.lastBlock, "hash", validator.blocks[ev.Block.NumberU64()].hash,
						"txCount", validator.blocks[ev.Block.NumberU64()].txs)

					if atomic.LoadUint32(testCanBeStopped) == 1 {
						if atomic.LoadInt64(test.validatorsCanBeStopped) == int64(len(validators)) {
							break wgLoop
						}
						if atomic.LoadInt64(test.validatorsCanBeStopped) > int64(len(validators)) {
							return fmt.Errorf("something is wrong. %d of %d validators are ready to be stopped", atomic.LoadInt64(test.validatorsCanBeStopped), uint32(len(validators)))
						}
						continue
					}

					// actual forming and sending transaction
					logger.Debug("peer", "address", crypto.PubkeyToAddress(validator.privateKey.PublicKey).String(), "block", ev.Block.Number().Uint64(), "isRunning", validator.isRunning)

					if validator.isRunning {
						txsMu.Lock()
						if _, ok := txs[validator.lastBlock]; !ok {
							txs[validator.lastBlock] = ev.Block.Transactions().Len()
						}
						txsMu.Unlock()

						for _, tx := range ev.Block.Transactions() {
							validator.transactionsMu.Lock()
							if _, ok := validator.transactions[tx.Hash()]; ok {
								validator.txsChainCount[ev.Block.NumberU64()]++
								delete(validator.transactions, tx.Hash())
							}
							validator.transactionsMu.Unlock()
						}

						if int(validator.lastBlock) <= test.numBlocks {
							for i := 0; i < txPerPeer; i++ {
								nextValidatorIndex := (index + i + 1) % len(validators)
								toAddr := crypto.PubkeyToAddress(validators[nextValidatorIndex].privateKey.PublicKey)
								var tx *types.Transaction
								var innerErr error
								if f, ok := test.sendTransactionHooks[nextValidatorIndex]; ok {
									if tx, innerErr = f(validator.service, validator.privateKey, fromAddr, toAddr); innerErr != nil {
										return innerErr
									}
								} else {
									if tx, innerErr = sendTx(validator.service, validator.privateKey, fromAddr, toAddr, generateRandomTx); innerErr != nil {
										return innerErr
									}

								}

								atomic.AddInt64(validator.txsSendCount, 1)

								validator.transactionsMu.Lock()
								validator.transactions[tx.Hash()] = struct{}{}
								validator.transactionsMu.Unlock()
							}
						}
					}

					// after hook
					err = runHook(test.getAfterHook(index), test, ev.Block, validator, index)
					if err != nil {
						return err
					}

					// check transactions status if all blocks are passed
					if int(validator.lastBlock) > test.numBlocks {
						//all transactions were included into the chain
						if errorOnTx {
							validator.transactionsMu.Lock()
							if len(validator.transactions) == 0 {
								if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
									atomic.AddInt64(test.validatorsCanBeStopped, 1)
								}
							}
							validator.transactionsMu.Unlock()
						} else {
							if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
								atomic.AddInt64(test.validatorsCanBeStopped, 1)
							}
						}
					}

					if validator.isRunning && int(validator.lastBlock) >= test.numBlocks+blocksToWait {
						if errorOnTx {
							pending, queued := validator.service.TxPool().Stats()
							if pending > 0 {
								return fmt.Errorf("after a new block it should be 0 pending transactions got %d. block %d", pending, ev.Block.Number().Uint64())
							}
							if queued > 0 {
								return fmt.Errorf("after a new block it should be 0 queued transactions got %d. block %d", queued, ev.Block.Number().Uint64())
							}

							validator.transactionsMu.Lock()
							pendingTransactions := len(validator.transactions)
							havePendingTransactions := pendingTransactions != 0
							validator.transactionsMu.Unlock()

							if havePendingTransactions {
								var txsChainCount int64
								for _, txsBlockCount := range validator.txsChainCount {
									txsChainCount += txsBlockCount
								}

								if validator.wasStopped {
									//fixme an error should be returned
									logger.Error("test error!!!", "err", fmt.Errorf("a validator %d still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
										index,
										pendingTransactions, ev.Block.Number().Uint64(),
										atomic.LoadInt64(validator.txsSendCount), txsChainCount))

									if atomic.CompareAndSwapUint32(testCanBeStopped, 0, 1) {
										atomic.AddInt64(test.validatorsCanBeStopped, 1)
									}
								} else {
									return fmt.Errorf("a validator %d still have transactions to be mined %d. block %d. Total sent %d, total mined %d",
										index,
										pendingTransactions, ev.Block.Number().Uint64(),
										atomic.LoadInt64(validator.txsSendCount), txsChainCount)
								}
							}
						}
					}
				case innerErr := <-validator.subscription.Err():
					if innerErr != nil {
						return fmt.Errorf("error in blockchain %q", innerErr)
					}

					time.Sleep(500 * time.Millisecond)

					// after hook
					err = runHook(test.getAfterHook(index), test, nil, validator, index)
					if err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				default:
				}
				// allow to exit goroutine when no quorum expected
				// check that there's no quorum within the given noQuorumTimeout
				if !hasQuorum(validators) && test.noQuorumAfterBlock > 0 {
					log.Error("No Quorum", "index", index, "last_block", validator.lastBlock)
					// wait for quorum to get restored
					time.Sleep(test.noQuorumTimeout)
					break wgLoop
				}
			}
			return nil
		})
	}
	err := wg.Wait()
	keys := make([]int, 0, len(txs))
	for key := range txs {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	fmt.Println("Transactions per block")
	for _, key := range keys {
		count := txs[uint64(key)]
		fmt.Printf("Block %d has %d transactions\n", key, count)
	}
	fmt.Println("\nPending transactions")
	for index, validator := range validators {
		validator.transactionsMu.Lock()
		fmt.Printf("Validator %d has %d transactions\n", index, len(validator.transactions))
		validator.transactionsMu.Unlock()
	}
	if err != nil {
		t.Fatal(err)
	}

	// no blocks can be mined with no quorum
	if test.noQuorumAfterBlock > 0 {
		for index, validator := range validators {
			if validator.lastBlock < test.noQuorumAfterBlock-1 {
				t.Fatalf("validator [%d] should have mined blocks. expected block number %d, but got %d",
					index, test.noQuorumAfterBlock-1, validator.lastBlock)
			}

			if validator.lastBlock > test.noQuorumAfterBlock {
				t.Fatalf("validator [%d] mined blocks without quorum. expected block number %d, but got %d",
					index, test.noQuorumAfterBlock, validator.lastBlock)
			}
		}
	}

	// check that all nodes reached the same minimum blockchain height
	minHeight := math.MaxInt64
	for index, validator := range validators {
		if len(test.maliciousPeers) != 0 {
			if _, ok := test.maliciousPeers[index]; ok {
				// don't check chain for malicious peers
				continue
			}
		}

		validatorBlock := validator.lastBlock
		if minHeight > int(validatorBlock) {
			minHeight = int(validatorBlock)
		}

		if test.noQuorumAfterBlock > 0 {
			continue
		}

		if validatorBlock < uint64(test.numBlocks) {
			t.Fatalf("a validator is behind the network index %d and block %v - expected %d",
				index, validatorBlock, test.numBlocks)
		}
	}

	// check that all nodes got the same blocks
	for i := 1; i <= minHeight; i++ {
		blockHash := validators[0].blocks[uint64(i)].hash

		for index, validator := range validators[1:] {
			if len(test.maliciousPeers) != 0 {
				if _, ok := test.maliciousPeers[index+1]; ok {
					// don't check chain for malicious peers
					continue
				}
			}
			if validator.blocks[uint64(i)].hash != blockHash {
				t.Fatalf("validators %d and %d have different blocks %d - %q vs %s",
					0, index+1, i, validator.blocks[uint64(i)].hash.String(), blockHash.String())
			}
		}
	}
	fmt.Println("\nTransactions OK")
}

func hasQuorum(validators []*testNode) bool {
	active := 0
	for _, val := range validators {
		if val.isRunning {
			active++
		}
	}
	return quorum(len(validators), active)
}

func quorum(valCount, activeVals int) bool {
	return float64(activeVals) >= math.Ceil(float64(2)/float64(3)*float64(valCount))
}

func runHook(validatorHook hook, test *testCase, block *types.Block, validator *testNode, index int) error {
	if validatorHook == nil {
		return nil
	}

	err := validatorHook(block, validator, test, time.Now())
	if err != nil {
		return fmt.Errorf("error while executing before hook for validator index %d and block %v, err %v",
			index, block.NumberU64(), err)
	}

	return nil
}

func hookStopNode(nodeIndex int, blockNum uint64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		if block.Number().Uint64() == blockNum {
			err := validator.stopNode()
			if err != nil {
				return err
			}

			tCase.setStopTime(nodeIndex, currentTime)
		}

		return nil
	}
}

func hookForceStopNode(nodeIndex int, blockNum uint64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		if block.Number().Uint64() == blockNum {
			err := validator.forceStopNode()
			if err != nil {
				return err
			}
			tCase.setStopTime(nodeIndex, currentTime)
		}
		return nil
	}
}

func hookStartNode(nodeIndex int, durationAfterStop float64) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		stopTime := tCase.getStopTime(nodeIndex)
		if block == nil && currentTime.Sub(stopTime).Seconds() >= durationAfterStop {
			if err := validator.startNode(); err != nil {
				return err
			}

			if err := validator.startService(); err != nil {
				return err
			}
		}

		return nil
	}
}
