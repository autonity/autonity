package rpc_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/clearmatics/autonity/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCs(t *testing.T) {
	users, err := test.Users(1, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	g, err := test.Genesis(users)
	require.NoError(t, err)
	n, err := test.NewNode(users[0], g)
	defer n.Close()
	n.Start()
	n.Eth.StopMining() // No need to mine

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	zero := big.NewInt(0)
	balance, err := n.WsClient.BalanceAt(ctx, n.Address, zero)
	require.NoError(t, err)
	assert.Equal(t, uint64(10e18), balance.Uint64(), "BalanceAt")

	chainID, err := n.WsClient.ChainID(ctx)
	require.NoError(t, err)
	assert.Equal(t, g.Config.ChainID, chainID, "ChainID")

	networkID, err := n.WsClient.NetworkID(ctx)
	require.NoError(t, err)
	assert.Equal(t, n.EthConfig.NetworkId, networkID.Uint64(), "NetworkID")

	pendingBalance, err := n.WsClient.PendingBalanceAt(ctx, n.Address)
	require.NoError(t, err)
	assert.Equal(t, uint64(10e18), pendingBalance.Uint64(), "PendingBalanceAt")

	suggestedGasPrice, err := n.WsClient.SuggestGasPrice(ctx)
	require.NoError(t, err)
	assert.Equal(t, g.Config.AutonityContractConfig.MinGasPrice, suggestedGasPrice.Uint64(), "SuggestGasPrice")
}
