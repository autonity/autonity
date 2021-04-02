package test

import (
	"context"
	"github.com/clearmatics/autonity/faultdetector"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFaultDetectorMaliciousBehaviourPN(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.PN, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.PN)
}

func TestFaultDetectorMaliciousBehaviourPO(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.PO, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.PO)
}

func TestFaultDetectorMaliciousBehaviourPVN(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.PVN, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.PVN)
}

func TestFaultDetectorMaliciousBehaviourC(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.C, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.C)
}

func TestFaultDetectorMaliciousBehaviourInvalidProposal(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.InvalidProposal, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.InvalidProposal)
}

func TestFaultDetectorMaliciousBehaviourInvalidProposer(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.InvalidProposer, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.InvalidProposer)
}

func TestFaultDetectorMaliciousBehaviourEquivocation(t *testing.T) {
	network, cancel := setupNetworkWithMaliciousBehaviour(t, 4, faultdetector.Equivocation, 90*time.Second)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, faultdetector.Equivocation)
}

func setupNetworkWithMaliciousBehaviour(t *testing.T, committeeSize int, ruleID faultdetector.Rule,
	timeout time.Duration) (test.Network, context.CancelFunc) {
	users, err := test.Users(committeeSize, "10e18,v,1,0.0.0.0:%s,%s", 11111)
	require.NoError(t, err)
	for i := 0; i < committeeSize; i++ {
		users[i].UserType = params.UserValidator
		users[i].Stake = 1
	}

	network, err := test.NewNetworkWithMaliciousUser(users, uint8(ruleID))
	require.NoError(t, err)

	_, cancel := context.WithTimeout(context.Background(), timeout)
	return network, cancel
}

func verifyOnChainProof(t *testing.T, network test.Network, ruleID faultdetector.Rule) {
	// todo: fetch and check if on-chain proofs of misbehaviour presented.
}
