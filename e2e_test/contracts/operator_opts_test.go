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
	newOperator := crypto.PubkeyToAddress(newOptKey.PublicKey)
	err = client.AwaitSetOperator(initialOptKey, newOperator, tm)
	require.NoError(t, err)
	op, err := client.Interactor.Call(nil).GetOperator()
	require.NoError(t, err)
	require.Equal(t, newOperator, op)

	newBaseFee := new(big.Int).SetUint64(12)
	err = client.AwaitSetMinBaseFee(newOptKey, newBaseFee, tm)
	require.NoError(t, err)
	fee, err := client.Interactor.Call(nil).GetMinBaseFee()
	require.NoError(t, err)
	require.Equal(t, 0, newBaseFee.Cmp(fee))

	newSize := new(big.Int).SetUint64(66)
	err = client.AwaitSetCommitteeSize(newOptKey, newSize, tm)
	require.NoError(t, err)
	size, err := client.Interactor.Call(nil).GetMaxCommitteeSize()
	require.NoError(t, err)
	require.Equal(t, 0, newSize.Cmp(size))

	newBondingPeriod := new(big.Int).SetUint64(60)
	err = client.AwaitSetUnbondingPeriod(newOptKey, newBondingPeriod, tm)
	require.NoError(t, err)
	bp, err := client.Interactor.Call(nil).GetUnbondingPeriod()
	require.NoError(t, err)
	require.Equal(t, 0, newBondingPeriod.Cmp(bp))

	newEpochPeriod := new(big.Int).SetUint64(60)
	err = client.AwaitSetEpochPeriod(newOptKey, newEpochPeriod, tm)
	require.NoError(t, err)
	ep, err := client.Interactor.Call(nil).GetEpochPeriod()
	require.NoError(t, err)
	require.Equal(t, 0, newEpochPeriod.Cmp(ep))

	newTreasury, err := crypto.GenerateKey()
	require.NoError(t, err)
	newTreasuryAddr := crypto.PubkeyToAddress(newTreasury.PublicKey)
	err = client.AwaitSetTreasuryAccount(newOptKey, newTreasuryAddr, tm)
	require.NoError(t, err)
	treasury, err := client.Interactor.Call(nil).GetTreasuryAccount()
	require.NoError(t, err)
	require.Equal(t, newTreasuryAddr, treasury)

	newTreasuryFee := new(big.Int).SetUint64(1000)
	err = client.AwaitSetTreasuryFee(newOptKey, newTreasuryFee, tm)
	require.NoError(t, err)
	tFee, err := client.Interactor.Call(nil).GetTreasuryFee()
	require.NoError(t, err)
	require.Equal(t, 0, newTreasuryFee.Cmp(tFee))
}
