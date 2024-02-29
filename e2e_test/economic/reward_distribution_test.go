package economic

import (
	"context"
	"fmt"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

// This test checks that when a transaction is processed the fees are divided
// between validators and stakeholders.
func TestFeeRedistributionValidatorsAndDelegators(t *testing.T) {
	t.Skip("Is broken with Penalty Absorbing Stake")
	//todo: should be rewrite once the ATN inflation is introduced.
	//todo: fix. Genesis validators are no longer issued Liquid Newton. Need to introduce 3rd party delegators.
	vals, err := e2e.Validators(t, 3, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	vals[2].Stake = 25000

	network, err := e2e.NewNetworkFromValidators(t, vals, true)
	require.NoError(t, err)
	defer network.Shutdown()

	n := network[0]

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// retrieve current balance
	// send liquid newton to some random address
	// check balance - shouldnt have increased
	// wait for epoch
	// check claimable fees
	// redeem fees

	// Setup Bindings
	autonityContract, _ := autonity.NewAutonity(params.AutonityContractAddress, n.WsClient)
	valAddrs, _ := autonityContract.GetValidators(nil)
	liquidContracts := make([]*autonity.Liquid, len(valAddrs))
	validators := make([]autonity.AutonityValidator, len(valAddrs))
	for i, valAddr := range valAddrs {
		validators[i], _ = autonityContract.GetValidator(nil, valAddr)
		liquidContracts[i], _ = autonity.NewLiquid(validators[i].LiquidContract, n.WsClient)
	}
	transactor, _ := bind.NewKeyedTransactorWithChainID(vals[0].TreasuryKey, big.NewInt(1234))
	tx, err := liquidContracts[0].Transfer(
		transactor,
		common.Address{66, 66}, big.NewInt(1337))

	require.NoError(t, err)
	_ = network.WaitToMineNBlocks(2, 20, false)
	tx2, err := n.SendAUT(ctx, network[1].Address, 10)
	require.NoError(t, err)
	err = network.AwaitTransactions(ctx, tx, tx2)
	require.NoError(t, err)
	// claimable fees should be 0 before epoch
	for i := range liquidContracts {
		unclaimed, _ := liquidContracts[i].UnclaimedRewards(&bind.CallOpts{}, validators[i].Treasury)
		require.Equal(t, big.NewInt(0).Bytes(), unclaimed.Bytes())
	}

	// wait for epoch

	// calculate reward pool
	r1, _ := n.WsClient.TransactionReceipt(context.Background(), tx.Hash())
	r2, _ := n.WsClient.TransactionReceipt(context.Background(), tx2.Hash())

	b1, _ := n.WsClient.BlockByNumber(context.Background(), r1.BlockNumber)
	b2, _ := n.WsClient.BlockByNumber(context.Background(), r2.BlockNumber)

	rewardT1 := new(big.Int).Mul(new(big.Int).SetUint64(r1.CumulativeGasUsed), b1.BaseFee())
	rewardT2 := new(big.Int).Mul(new(big.Int).SetUint64(r2.CumulativeGasUsed), b2.BaseFee())

	totalFees := new(big.Int).Add(rewardT1, rewardT2)
	treasuryRewards := new(big.Int).Div(new(big.Int).Mul(totalFees, new(big.Int).SetUint64(15)), big.NewInt(10000))
	totalRewards := new(big.Int).Sub(totalFees, treasuryRewards)

	balanceBeforeEpoch, _ := n.WsClient.BalanceAt(context.Background(), params.AutonityContractAddress, nil)
	require.Equal(t, totalFees, balanceBeforeEpoch)

	err = network.WaitToMineNBlocks(30, 90, false)
	require.NoError(t, err)

	fmt.Println("total rewards", totalRewards)
	balanceGlobalTreasury, _ := n.WsClient.BalanceAt(context.Background(), common.Address{120}, nil)
	cfg, _ := autonityContract.Config(nil)
	fmt.Println(cfg)
	require.Equal(t, treasuryRewards, balanceGlobalTreasury)

	stake := []int64{10000 - 1337, 10000, 25000}
	epochStake := []int64{10000, 10000, 25000}
	totalStake := int64(45000)
	for i := range liquidContracts {
		unclaimed, _ := liquidContracts[i].UnclaimedRewards(&bind.CallOpts{}, validators[i].Treasury)
		totalValRewards := new(big.Int).Div(new(big.Int).Mul(totalRewards, big.NewInt(epochStake[i])), big.NewInt(totalStake))
		valCommission := new(big.Int).Div(new(big.Int).Mul(totalValRewards, big.NewInt(12)), big.NewInt(100))
		stakerReward := new(big.Int).Sub(totalValRewards, valCommission)
		require.Equal(t, new(big.Int).Div(new(big.Int).Mul(stakerReward, big.NewInt(stake[i])), big.NewInt(epochStake[i])), unclaimed)
	}
}
