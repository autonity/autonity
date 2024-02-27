package contract_opt

import (
	"context"
	"crypto/ecdsa"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

var (
	mintAmount     = new(big.Int).SetUint64(100)
	approvedAmount = new(big.Int).SetUint64(15)
)

func TestACERC20Interfaces(t *testing.T) {
	network, err := e2e.NewNetwork(t, 4, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	operatorNode := network[0]
	operatorKey := operatorNode.Key
	accounts, err := makeAccounts(2)
	require.NoError(t, err)

	err = fundingAccounts(operatorNode, accounts)
	require.NoError(t, err)

	// mint NTN for accounts
	timeout := 5 * time.Second
	for _, account := range accounts {
		err = operatorNode.AwaitMintNTN(operatorKey, crypto.PubkeyToAddress(account.PublicKey), mintAmount, timeout)
		require.NoError(t, err)
	}

	Alice := accounts[0]
	Bob := accounts[1]
	// Alice grant NTN transferFrom approval to Bob with 15 NTN
	err = operatorNode.AwaitApproveNTN(Alice, crypto.PubkeyToAddress(Bob.PublicKey), approvedAmount, timeout)
	require.NoError(t, err)

	// Bob transfer from Alice's account with the approved amount of NTN to Bob's account
	err = operatorNode.AwaitTransferFromNTN(Bob, crypto.PubkeyToAddress(Alice.PublicKey), approvedAmount, timeout)
	require.NoError(t, err)

	// Bot transfer 15 NTNs back to Alice.
	err = operatorNode.AwaitTransferNTN(Bob, crypto.PubkeyToAddress(Alice.PublicKey), approvedAmount, timeout)
	require.NoError(t, err)

	for _, account := range accounts {
		balance, err := operatorNode.BalanceNTN(crypto.PubkeyToAddress(account.PublicKey))
		require.NoError(t, err)
		require.Equal(t, mintAmount.Uint64(), balance.Uint64())
	}
}

func makeAccounts(num int) ([]*ecdsa.PrivateKey, error) {
	var accounts []*ecdsa.PrivateKey
	for i := 0; i < num; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, key)
	}
	return accounts, nil
}

// fundingAccounts distribute 1 ATN from operator account to each account in the accounts list.
func fundingAccounts(operatorNode *e2e.Node, accounts []*ecdsa.PrivateKey) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var txs []*types.Transaction
	for _, account := range accounts {
		tx, err := operatorNode.SendAUT(ctx, crypto.PubkeyToAddress(account.PublicKey), params.Ether)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
	}

	return operatorNode.AwaitTransactions(ctx, txs...)
}
