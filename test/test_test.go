package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/log"
	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
	users, err := Users(5, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	g, err := Genesis(users)
	require.NoError(t, err)
	var network []*Node
	for _, u := range users {
		n, cleanup, err := NewNode(u, g)
		defer cleanup()
		require.NoError(t, err)
		err = n.Start()
		require.NoError(t, err)

		network = append(network, n)
	}

	for _, n := range network {
		var ethereum *eth.Ethereum
		if err := n.Service(&ethereum); err != nil {
			require.NoError(t, err)
		}
		err = ethereum.StartMining(1)
		require.NoError(t, err)
	}
	// There is a race condition in miner.worker its field snapshotBlock is set
	// only when new transacting are received or commitNewWork is called. But
	// both of these happen in goroutines separate to the call to miner.Start
	// and miner.Strart does not wait for snapshotBlock to be set. Therfore
	// there is currently no way to know when it is safe to call estimate gas.
	// What we do here is sleep a bit and cross our fingers.
	time.Sleep(10 * time.Millisecond)

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
