package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/node"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
	users, err := Users(4, "10e18", "v", "1", 6780)
	require.NoError(t, err)
	g, err := Genesis(users)
	require.NoError(t, err)
	var network []*node.Node
	for _, u := range users {
		n, cleanup, err := Node(u, g)
		defer cleanup()
		require.NoError(t, err)
		err = n.Start()

		var ethereum *eth.Ethereum
		if err := n.Service(&ethereum); err != nil {
			require.NoError(t, err)
		}
		err = ethereum.StartMining(1)
		require.NoError(t, err)
		network = append(network, n)
	}
	time.Sleep(time.Second * 4)

	c, err := ethclient.Dial("ws://" + network[0].WSEndpoint())
	require.NoError(t, err)

	ctx := context.WithCancel(context.Background())
	hashes, errors, err := MinedTransactions(ctx, client*ethclient.Client)
	require.NoError(t, err)
	// sendtranaction

	for _, n := range network {
		spew.Dump("peers", crypto.PubkeyToAddress(n.Config().NodeKey().PublicKey), n.Server().Peers())
	}
	time.Sleep(time.Second * 20)
}
