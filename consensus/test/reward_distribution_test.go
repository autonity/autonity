package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"math/big"
	"testing"
	"time"
)

func TestRewardDistribution(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	hookGenerator := func() hook {
		type EconomicMetaData struct {
			Accounts        []common.Address `abi:"accounts"`
			Usertypes       []uint8          `abi:"usertypes"`
			Stakes          []*big.Int       `abi:"stakes"`
			Commissionrates []*big.Int       `abi:"commissionrates"`
			Mingasprice     *big.Int         `abi:"mingasprice"`
			Stakesupply     *big.Int         `abi:"stakesupply"`
		}

		// each new mined block will be checked by reward distribution checking.
		fRewardChecker := func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
			// reward distribution checking start from height-2.
			if block.NumberU64() < 2 {
				return nil
			}
			blockReward := new(big.Int)
			gp := new(core.GasPool).AddGas(block.GasLimit())
			parent := validator.service.BlockChain().GetBlock(block.ParentHash(), block.NumberU64()-1)

			addr := crypto.PubkeyToAddress(validator.privateKey.PublicKey)
			// get parent block view.
			parentState, stateErr := validator.service.BlockChain().StateAt(parent.Root())
			if stateErr != nil {
				return stateErr
			}
			balanceBase := parentState.GetBalance(addr)
			// calculate sent, gasUsed, receive amount and reward base on new block.
			sentAmount := new(big.Int)
			receivedAmount := new(big.Int)
			usedGas := new(big.Int)
			rewardFraction := new(big.Int)
			for i, tx := range block.Transactions() {
				// Apply TX to parent view again to estimate the gasUsed of the TX.
				parentState.Prepare(tx.Hash(), block.Hash(), i)
				receipt, receiptErr := core.ApplyTransaction(validator.service.BlockChain().Config(), validator.service.BlockChain(), nil, gp, parentState, block.Header(), tx, new(uint64), *validator.service.BlockChain().GetVMConfig())
				if receiptErr != nil {
					return receiptErr
				}

				// count block reward.
				txGasConsumed := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipt.GasUsed))
				blockReward.Add(blockReward, txGasConsumed)

				// count received amount on the block.
				if *tx.To() == addr {
					receivedAmount.Add(receivedAmount, tx.Value())
				}

				// count sent and gasUsed amount on the block.
				validator.transactionsMu.Lock()
				if _, ok := validator.transactions[tx.Hash()]; ok {
					sentAmount.Add(sentAmount, tx.Value())
					usedGas.Add(usedGas, txGasConsumed)
				}
				validator.transactionsMu.Unlock()
			}

			// calculate reward fractions.
			parentEcomonimcState := EconomicMetaData{make([]common.Address, 32), make([]uint8, 32), make([]*big.Int, 32),
				make([]*big.Int, 32), new(big.Int), new(big.Int)}
			err := validator.service.BlockChain().GetAutonityContract().AutonityContractCall(parentState, parent.Header(), "dumpEconomicsMetricData", &parentEcomonimcState)
			if err != nil {
				return err
			}

			for i := 0; i < len(parentEcomonimcState.Accounts); i++ {
				if addr == parentEcomonimcState.Accounts[i] {
					rewardFraction = new(big.Int).Mul(blockReward, parentEcomonimcState.Stakes[i])
					rewardFraction = rewardFraction.Div(rewardFraction, parentEcomonimcState.Stakesupply)
					break
				}
			}

			// current block view.
			currentState, stateErr := validator.service.BlockChain().StateAt(block.Root())
			if stateErr != nil {
				return stateErr
			}
			balanceActual := currentState.GetBalance(addr)
			// check if balance is expected: base + reward + received - sent - gasUsed == bActual.
			balanceWant := new(big.Int).Add(balanceBase, rewardFraction)
			balanceWant.Add(balanceWant, receivedAmount)
			balanceWant.Sub(balanceWant, sentAmount)
			balanceWant.Sub(balanceWant, usedGas)
			if balanceWant.Cmp(balanceActual) != 0 {
				return fmt.Errorf("rewardFraction distribution might in-correct")
			}
			return nil
		}

		return fRewardChecker
	}

	rewardChecker := hookGenerator()

	cases := []*testCase{
		{
			name:          "check reward distribution per block",
			numValidators: 6,
			numBlocks:     30,
			txPerPeer:     1,
			// All nodes in the network check its reward distribution.
			beforeHooks: map[string]hook{
				"VA": rewardChecker,
				"VB": rewardChecker,
				"VC": rewardChecker,
				"VD": rewardChecker,
				"VE": rewardChecker,
				"VF": rewardChecker,
			},
			stopTime: make(map[string]time.Time),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}
