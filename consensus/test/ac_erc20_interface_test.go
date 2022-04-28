package test

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

// test ERC-20 interfaces
func TestACERC20Interfaces(t *testing.T) {
	numOfValidators := 3
	operator, err := makeAccount()
	require.NoError(t, err)
	operatorAddr := crypto.PubkeyToAddress(operator.PublicKey)

	owner, err := makeAccount()
	require.NoError(t, err)
	ownerAddr := crypto.PubkeyToAddress(owner.PublicKey)

	spender, err := makeAccount()
	require.NoError(t, err)
	spenderAddr := crypto.PubkeyToAddress(spender.PublicKey)

	amountTX := new(big.Int).SetUint64(50)
	initBalance := mintAmount
	amountApproved := new(big.Int).SetUint64(20)
	amountTransferFrom := new(big.Int).SetUint64(10)

	testCase := &testCase{
		name:          "ERC20 interface test",
		numValidators: numOfValidators,
		numBlocks:     15,
		txPerPeer:     1,
		// set operator and pre-mine native tokens.
		genesisHook: func(g *core.Genesis) *core.Genesis {
			g.Config.AutonityContractConfig.Operator = operatorAddr
			g.Alloc[operatorAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			g.Alloc[ownerAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			g.Alloc[spenderAddr] = core.GenesisAccount{
				Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
			}
			return g
		},
		afterHooks: map[string]hook{
			// mint 100 Newtons for owner by operator.
			"V0": mintNewtonHook(
				map[uint64]struct{}{
					1: {},
				},
				operator,
				ownerAddr,
			),
			// owner approve 20 Newtons to spender, operator send 50 Newtons to V1.
			"V1": transferApproveHook(
				map[uint64]struct{}{
					4: {},
				},
				owner,
				operator,
				spenderAddr,
				amountApproved,
				amountTX,
			),
			// spender transferFrom owner with 10 Newtons since owner approved 20.
			"V2": transferFromHook(
				map[uint64]struct{}{
					8: {},
				},
				spender,
				ownerAddr,
				amountTransferFrom,
			),
		},
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			node := validators["V1"]
			client := interact(node.rpcPort)
			defer client.close()
			v1Balance, err := client.call(node.lastBlock).balanceOf(node.EthAddress())
			require.NoError(t, err)
			operatorBalance, err := client.call(node.lastBlock).balanceOf(operatorAddr)
			require.NoError(t, err)
			ownerBalance, err := client.call(node.lastBlock).balanceOf(ownerAddr)
			require.NoError(t, err)
			spenderBalance, err := client.call(node.lastBlock).balanceOf(spenderAddr)
			require.NoError(t, err)

			require.Equal(t, amountTX.Uint64(), v1Balance.Uint64())
			require.Equal(t, initBalance.Uint64()-amountTX.Uint64(), operatorBalance.Uint64())
			require.Equal(t, initBalance.Uint64(), ownerBalance.Uint64()+amountTransferFrom.Uint64())
			require.Equal(t, spenderBalance.Uint64(), amountTransferFrom.Uint64())

			allowance, err := client.call(node.lastBlock).allowance(ownerAddr, spenderAddr)
			require.NoError(t, err)
			require.Equal(t, amountApproved.Uint64()-amountTransferFrom.Uint64(), allowance.Uint64())
		},
	}

	runTest(t, testCase)
}

func transferFromHook(toDo map[uint64]struct{}, spender *ecdsa.PrivateKey, owner common.Address, amountTX *big.Int) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := toDo[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		spenderAddr := crypto.PubkeyToAddress(spender.PublicKey)
		if _, err := interaction.tx(spender).transferFrom(owner, spenderAddr, amountTX); err != nil {
			return err
		}
		return nil
	}
}

func transferApproveHook(toDo map[uint64]struct{}, owner, op *ecdsa.PrivateKey, spender common.Address, amountApr, amountTX *big.Int) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := toDo[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(owner).approve(spender, amountApr); err != nil {
			return err
		}
		if _, err := interaction.tx(op).transfer(validator.EthAddress(), amountTX); err != nil {
			return err
		}
		return nil
	}
}

func mintNewtonHook(toDo map[uint64]struct{}, op *ecdsa.PrivateKey, owner common.Address) hook {
	return func(block *types.Block, validator *testNode, tCase *testCase, currentTime time.Time) error {
		blockNum := block.Number().Uint64()
		if _, ok := toDo[blockNum]; !ok {
			return nil
		}
		interaction := interact(validator.rpcPort)
		defer interaction.close()
		if _, err := interaction.tx(op).mint(owner, mintAmount); err != nil {
			return err
		}
		// let operator mint Newton to himself.
		opAddr := crypto.PubkeyToAddress(op.PublicKey)
		if _, err := interaction.tx(op).mint(opAddr, mintAmount); err != nil {
			return err
		}
		return nil
	}
}
