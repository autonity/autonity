package contracts

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestOperatorOpts(t *testing.T) {
	network, err := e2e.NewNetwork(t, 4, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(2, 10, false)

	initialOptKey := network[0].Key
	newOptKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	client := network[0]
	// funding for new operator.
	err = fundingAccounts(client, []*ecdsa.PrivateKey{newOptKey})

	tm := 5 * time.Second
	err = client.AwaitSetOperator(initialOptKey, crypto.PubkeyToAddress(newOptKey.PublicKey), tm)
	require.NoError(t, err)

	err = client.AwaitSetMinBaseFee(newOptKey, new(big.Int).SetUint64(12), tm)
	require.NoError(t, err)

	err = client.AwaitSetCommitteeSize(newOptKey, new(big.Int).SetUint64(66), tm)
	require.NoError(t, err)

	err = client.AwaitSetUnbondingPeriod(newOptKey, new(big.Int).SetUint64(60), tm)
	require.NoError(t, err)

	err = client.AwaitSetEpochPeriod(newOptKey, new(big.Int).SetUint64(60), tm)
	require.NoError(t, err)

	newTreasury, err := crypto.GenerateKey()
	require.NoError(t, err)

	err = client.AwaitSetTreasuryAccount(newOptKey, crypto.PubkeyToAddress(newTreasury.PublicKey), tm)
	require.NoError(t, err)

	err = client.AwaitSetTreasuryFee(newOptKey, new(big.Int).SetUint64(1000), tm)
	require.NoError(t, err)
}
