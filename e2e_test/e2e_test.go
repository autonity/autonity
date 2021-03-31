package test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/test"
	"github.com/stretchr/testify/require"
)

// This test checks that we can process transactions that transfer value from
// one participant to another.
func TestSendingValue(t *testing.T) {
	network, err := test.NewNetwork(2, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	defer network.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = network[0].SendETracked(ctx, network[1].Address, 10)
	require.NoError(t, err)
}

// This test checks that when a transaction is processed the fees are divided
// between validators.
func TestFeeRedistributionOnlyValidators(t *testing.T) {
	network, err := test.NewNetwork(2, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	defer network.Shutdown()

	n := network[0]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := n.SendE(ctx, network[1].Address, 10)
	require.NoError(t, err)
	err = network.AwaitTransactions(ctx, tx)
	require.NoError(t, err)

	zero := big.NewInt(0)
	num := n.ProcessedTxBlock(tx).Number()

	fee, err := n.TxFee(ctx, tx)
	require.NoError(t, err)
	halfFee := fee.Div(fee, big.NewInt(2))

	// We will have paid the whole fee but then received half back
	startingBalance, err := n.WsClient.BalanceAt(ctx, n.Address, zero)
	require.NoError(t, err)
	endingBalance, err := n.WsClient.BalanceAt(ctx, n.Address, num)
	require.NoError(t, err)
	expected := big.NewInt(0).Sub(startingBalance, halfFee)
	expected.Sub(expected, tx.Value())

	require.Equal(t, expected, endingBalance)

	n = network[1]
	// We expect this node to have collected half the fees
	startingBalance, err = n.WsClient.BalanceAt(ctx, n.Address, zero)
	require.NoError(t, err)
	endingBalance, err = n.WsClient.BalanceAt(ctx, n.Address, num)
	require.NoError(t, err)
	expected = big.NewInt(0).Add(startingBalance, halfFee)
	expected.Add(expected, tx.Value())

	require.Equal(t, expected, endingBalance)
}

// This test checks that when a transaction is processed the fees are divided
// between validators and stakeholders.
func TestFeeRedistributionValidatorsAndStakeholders(t *testing.T) {
	// We want 3 users, 2 validators with 1 stake each and a stakeholder with 2
	// stake.
	users, err := test.Users(3, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	users[2].UserType = params.UserStakeHolder
	users[2].Stake = 2

	network, err := test.NewNetworkFromUsers(users)
	require.NoError(t, err)
	defer network.Shutdown()

	n := network[0]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := n.SendE(ctx, network[1].Address, 10)
	require.NoError(t, err)
	err = network.AwaitTransactions(ctx, tx)
	require.NoError(t, err)

	zero := big.NewInt(0)
	num := n.ProcessedTxBlock(tx).Number()
	fee, err := n.TxFee(ctx, tx)
	require.NoError(t, err)
	quaterFee := big.NewInt(0).Div(fee, big.NewInt(4))

	// The sending node will have paid the whole fee but then received 1/4 back.
	startingBalance, err := n.WsClient.BalanceAt(ctx, n.Address, zero)
	require.NoError(t, err)
	endingBalance, err := n.WsClient.BalanceAt(ctx, n.Address, num)
	require.NoError(t, err)

	expected := big.NewInt(0).Sub(startingBalance, fee)
	expected.Add(expected, quaterFee)
	expected.Sub(expected, tx.Value())

	require.Equal(t, expected, endingBalance)

	// This node will have received the value sent + 1/4 fee
	n = network[1]
	startingBalance, err = n.WsClient.BalanceAt(ctx, n.Address, zero)
	require.NoError(t, err)
	endingBalance, err = n.WsClient.BalanceAt(ctx, n.Address, num)
	require.NoError(t, err)

	expected = big.NewInt(0).Add(startingBalance, quaterFee)
	expected.Add(expected, tx.Value())
	require.Equal(t, expected, endingBalance)

	// The stakeholder has twice as much stake as the other nodes so will have
	// received 1/2 fee.
	n = network[2]
	startingBalance, err = n.WsClient.BalanceAt(ctx, n.Address, zero)
	require.NoError(t, err)
	endingBalance, err = n.WsClient.BalanceAt(ctx, n.Address, num)
	require.NoError(t, err)

	expected = big.NewInt(0).Add(startingBalance, quaterFee)
	expected.Add(expected, quaterFee)
	require.Equal(t, expected, endingBalance)
}

func TestStartingAndStoppingNodes(t *testing.T) {
	t.Skip("This test will be unreliable until https://github.com/clearmatics/autonity/issues/84 is fixed")
	network, err := test.NewNetwork(5, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	defer network.Shutdown()
	n := network[0]

	// Send a tx to see that the network is working
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendETracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[1].Close()
	require.NoError(t, err)

	// Send a tx to see that the network is working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendETracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[2].Close()
	require.NoError(t, err)

	// We have now stopped more than F nodes, so we expect tx processing to time out.
	// Well wait 5 times the avgTransactionDuration before we assume the tx is not being processed.
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = n.SendETracked(ctx, network[1].Address, 10)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
	}

	// We start a node again and expect the previously unprocessed transaction to be processed
	err = network[2].Start()
	require.NoError(t, err)

	// Ensure that the previously sent transaction is now processed
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.AwaitSentTransactions(ctx)
	require.NoError(t, err)

	// Send a tx to see that the network is still working
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = n.SendETracked(ctx, network[1].Address, 10)
	require.NoError(t, err)

	// Start the last stopped node
	err = network[1].Start()
	require.NoError(t, err)

	// Send a tx to see that the network is still working
	err = n.SendETracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)
}

func TestFaultDetectorMaliciousBehaviourPN(t *testing.T) {

}

func TestFaultDetectorMaliciousBehaviourPO(t *testing.T) {

}

func TestFaultDetectorMaliciousBehaviourPVN(t *testing.T) {

}

func TestFaultDetectorMaliciousBehaviourC(t *testing.T) {

}

func TestFaultDetectorMaliciousBehaviourInvalidProposal(t *testing.T) {

}

func TestFaultDetectorMaliciousBehaviourInvalidProposer(t *testing.T) {

}

func TestFaultDetectorMaliciousBehaviourEquivocation(t *testing.T) {

}
