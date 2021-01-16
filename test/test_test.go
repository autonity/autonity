package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	// log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
	users, err := Users(5, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	g, err := Genesis(users)
	require.NoError(t, err)
	network := make([]*Node, len(users))
	for i, u := range users {
		n, cleanup, err := NewNode(u, g)
		defer cleanup()
		require.NoError(t, err)
		err = n.Start()
		require.NoError(t, err)
		network[i] = n
	}

	for _, n := range network {
		err := n.Eth.StartMining(1)
		require.NoError(t, err)
	}
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transacting are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Strart does not wait for snapshotBlock to be set. Therfore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(20 * time.Millisecond)

	for i := range network {
		for j := range network {
			sender := network[i]
			receiver := network[j].Address
			err := sender.SendE(context.Background(), receiver, 10)
			require.NoError(t, err)
		}
	}
	for i := range network {
		err := network[i].AwaitSentTransactions(context.Background())
		require.NoError(t, err)
	}
}

func TestStartingAndStoppingNodes(t *testing.T) {
	// log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
	users, err := Users(5, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	g, err := Genesis(users)
	require.NoError(t, err)
	network := make([]*Node, len(users))
	for i, u := range users {
		n, cleanup, err := NewNode(u, g)
		defer cleanup()
		require.NoError(t, err)
		err = n.Start()
		require.NoError(t, err)
		network[i] = n
	}

	for _, n := range network {
		err := n.Eth.StartMining(1)
		require.NoError(t, err)
	}
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transacting are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Strart does not wait for snapshotBlock to be set. Therfore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(10 * time.Millisecond)

	n := network[0]
	println("----------------------------------1")
	// Send a tx to see that the network is working
	err = n.SendETracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[1].Close()
	require.NoError(t, err)

	println("----------------------------------2")
	// Send a tx to see that the network is working
	err = n.SendETracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)

	// Stop a node
	err = network[2].Close()
	require.NoError(t, err)

	// time.Sleep(time.Second * 5)

	println("lenlen", len(n.ProcessedTxs))

	println("----------------------------------3")
	// We have now stopped more than F nodes, so we expect tx sending to time out.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err = n.SendETracked(ctx, network[1].Address, 10)
	println("lenlen", len(n.ProcessedTxs))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
	}
	println("tx not mined while > f nodes stopped")

	// We start a node again and expect the previously unprocessed transaction to be processed
	println("sent tx len", len(n.SentTxs))
	// This is not working because the transaction tracker cannot update the
	// state of test.Node, so when we call sendE down below test.Node is
	// waiting for the transaction tracked here and the transaction just sent.
	tr, err := TrackTransactions(n.WsClient)
	require.NoError(t, err)

	err = network[2].Start()
	require.NoError(t, err)
	err = network[2].Eth.StartMining(1)
	require.NoError(t, err)

	_, err = tr.AwaitTransactions(context.Background(), []common.Hash{n.SentTxs[len(n.SentTxs)-1]})
	require.NoError(t, err)
	println("previously unmined tx now mined")
	n.SentTxs = n.SentTxs[:len(n.SentTxs)-1]

	// Send a tx to see that the network is working
	err = n.SendETracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)
	println("fresh tx mined")

	// Start the last stopped node
	err = network[1].Start()
	require.NoError(t, err)
	err = network[1].Eth.StartMining(1)
	require.NoError(t, err)

	// Send a tx to see that the network is working
	err = n.SendETracked(context.Background(), network[1].Address, 10)
	require.NoError(t, err)
	time.Sleep(time.Second)

	println("final tx mined")
}

//func TestWsSubscribes(t *testing.T) {
//	// log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
//	users, err := Users(1, "10e18,v,1,0.0.0.0:%s,%s", 6780)
//	require.NoError(t, err)
//	g, err := Genesis(users)
//	require.NoError(t, err)
//	n, cleanup, err := NewNode(users[0], g)
//	defer cleanup()
//	require.NoError(t, err)
//	err = n.Start()
//	require.NoError(t, err)
//	var ethereum *eth.Ethereum
//	if err := n.Service(&ethereum); err != nil {
//		require.NoError(t, err)
//	}
//	err = ethereum.StartMining(1)
//	require.NoError(t, err)

//	// There is a race condition in miner.worker its field snapshotBlock is set
//	// only when new transacting are received or commitNewWork is called. But
//	// both of these happen in goroutines separate to the call to miner.Start
//	// and miner.Strart does not wait for snapshotBlock to be set. Therfore
//	// there is currently no way to know when it is safe to call estimate gas.
//	// What we do here is sleep a bit and cross our fingers.
//	time.Sleep(10 * time.Millisecond)

//	// bal, err := n.WsClient.BalanceAt(context.Background(), n.Address, nil)
//	// require.NoError(t, err)
//	// fmt.Printf("BalancAt: %s\n", bal.String())

//	err = n.SendETracked(context.Background(), common.Address{}, 10)
//	require.NoError(t, err)

//	// bal, err = n.WsClient.BalanceAt(context.Background(), n.Address, nil)
//	// require.NoError(t, err)
//	// fmt.Printf("BalancAt: %s\n", bal.String())
//	//time.Sleep(10 * time.Second)

//	println("----------------------------------")
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
//	defer cancel()
//	err = n.SendETracked(ctx, common.Address{}, 10)
//	require.NoError(t, err)
//	// if !errors.Is(err, context.DeadlineExceeded) {
//	// 	t.Fatalf("expecting %q, instead got: %v ", context.DeadlineExceeded.Error(), err)
//	// }
//	// bal, err = n.WsClient.BalanceAt(context.Background(), n.Address, nil)
//	// require.NoError(t, err)
//	// fmt.Printf("BalancAt: %s\n", bal.String())

//	// for _, x := range n.SentTxs {
//	// 	r, err := n.WsClient.TransactionReceipt(context.Background(), x)
//	// 	require.NoError(t, err)
//	// 	fmt.Printf("tx block: %s\n", r.BlockNumber.String())
//	// }

//}
