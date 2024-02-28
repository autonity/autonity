package contracts

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidatorManagement(t *testing.T) {
	network, err := e2e.NewNetwork(t, 4, "10e18,v,1000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(2, 10, false)

	accounts, err := makeAccounts(2)
	require.NoError(t, err)
	oracle := accounts[0]
	delegator := accounts[1]

	validator, consensusKey, err := crypto.GenAutonityKeys()
	treasuryHex := crypto.PubkeyToAddress(validator.PublicKey).Hex()
	require.NoError(t, err)

	pop, err := crypto.AutonityPOPProof(validator, oracle, treasuryHex, consensusKey)
	require.NoError(t, err)

	operatorNode := network[0]
	operatorKey := operatorNode.Key

	err = fundingAccounts(operatorNode, []*ecdsa.PrivateKey{validator, delegator})
	require.NoError(t, err)
	timeout := 5 * time.Second
	oracleAddr := crypto.PubkeyToAddress(oracle.PublicKey)
	enodeURL := enode.V4DNSUrl(validator.PublicKey, "127.0.0.1", 30303, 30303) + ":30303"
	err = operatorNode.AwaitRegisterValidator(validator, enodeURL, oracleAddr, consensusKey.PublicKey().Marshal(), pop, timeout)
	require.NoError(t, err)

	delegatee := crypto.PubkeyToAddress(validator.PublicKey)
	err = operatorNode.AwaitMintNTN(operatorKey, crypto.PubkeyToAddress(delegator.PublicKey), mintAmount, timeout)
	require.NoError(t, err)
	err = operatorNode.AwaitBondStake(delegator, delegatee, mintAmount, timeout)
	require.NoError(t, err)

	epochPeriod := operatorNode.EthConfig.Genesis.Config.AutonityContractConfig.EpochPeriod
	for {
		if operatorNode.Eth.BlockChain().CurrentHeader().Number.Uint64()%epochPeriod == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	// un bond after epoch end
	err = operatorNode.AwaitUnbondStake(delegator, delegatee, mintAmount, timeout)
	require.NoError(t, err)
}
