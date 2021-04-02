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
	to := 90 * time.Second
	rule := faultdetector.PN
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
}

func TestFaultDetectorMaliciousBehaviourPO(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.PO
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
}

func TestFaultDetectorMaliciousBehaviourPVN(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.PVN
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
}

func TestFaultDetectorMaliciousBehaviourC(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.C
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
}

func TestFaultDetectorMaliciousBehaviourInvalidProposal(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.InvalidProposal
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
}

func TestFaultDetectorMaliciousBehaviourInvalidProposer(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.InvalidProposer
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
}

func TestFaultDetectorMaliciousBehaviourEquivocation(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.Equivocation
	committeeSize := 4
	network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	verifyOnChainProof(t, network, rule)
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

func verifyOnChainProof(t *testing.T, network test.Network, ruleID faultdetector.Rule) { // nolint
	// todo: fetch and check if on-chain proofs of misbehaviour presented. // nolint
}
