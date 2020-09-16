package test

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/node"
	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
	users, err := Users(4, "10e18,v,1,0.0.0.0:%s,%s", 6780)
	require.NoError(t, err)
	g, err := Genesis(users)
	require.NoError(t, err)
	var network []*node.Node
	for _, u := range users {
		n, cleanup, err := NewNode(u, g)
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

	n := network[0]
	c, err := ethclient.Dial("ws://" + n.WSEndpoint())
	require.NoError(t, err)

	var txs []*types.Transaction
	var hashes []common.Hash
	for i := range network {
		for j := range network {
			sender := network[i]
			receiver := network[j]
			tx, err := ValueTransferTransaction(
				c,
				sender.Server().PrivateKey,
				crypto.PubkeyToAddress(sender.Server().PrivateKey.PublicKey),
				crypto.PubkeyToAddress(receiver.Server().PrivateKey.PublicKey),
				big.NewInt(10),
			)
			require.NoError(t, err)
			txs = append(txs, tx)
			hashes = append(hashes, tx.Hash())
		}
	}
	tr, err := TrackTransactions(c)
	require.NoError(t, err)
	for i := range txs {
		err = c.SendTransaction(context.Background(), txs[i])
		require.NoError(t, err)
	}

	err = tr.AwaitTransactions(context.Background(), hashes)
}
