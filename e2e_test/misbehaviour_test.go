package test

import (
	"context"
	"github.com/clearmatics/autonity/faultdetector"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
	"github.com/clearmatics/autonity/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFaultDetectorMaliciousBehaviourPN(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.PN
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func TestFaultDetectorMaliciousBehaviourPO(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.PO
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func TestFaultDetectorMaliciousBehaviourPVN(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.PVN
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func TestFaultDetectorMaliciousBehaviourC(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.C
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func TestFaultDetectorMaliciousBehaviourInvalidProposal(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.InvalidProposal
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func TestFaultDetectorMaliciousBehaviourInvalidProposer(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.InvalidProposer
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func TestFaultDetectorMaliciousBehaviourEquivocation(t *testing.T) {
	to := 90 * time.Second
	rule := faultdetector.Equivocation
	committeeSize := 4
	ctx, network, cancel := setupNetworkWithMaliciousBehaviour(t, committeeSize, rule, to)
	defer network.Shutdown()
	defer cancel()
	keepCheckingProofExistence(t, to, ctx, network, rule)
}

func setupNetworkWithMaliciousBehaviour(t *testing.T, committeeSize int, ruleID faultdetector.Rule,
	timeout time.Duration) (context.Context, test.Network, context.CancelFunc) {
	users, err := test.Users(committeeSize, "10e18,v,1,0.0.0.0:%s,%s", 11111)
	require.NoError(t, err)
	for i := 0; i < committeeSize; i++ {
		users[i].UserType = params.UserValidator
		users[i].Stake = 1
	}

	network, err := test.NewNetworkWithMaliciousUser(users, uint8(ruleID))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return ctx, network, cancel
}

// keepCheckingProof will keep trying to check if the wanted proof is presented on-chain until either
// we get a result from doSomething() or the timeout expires
func keepCheckingProofExistence(t *testing.T, to time.Duration, c context.Context, net test.Network,
	ruleID faultdetector.Rule) bool {
	timeout := time.After(to)
	tick := time.Tick(1 * time.Second)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-timeout:
			return false
		case <-tick:
			if proofExists(t, c, net, ruleID) {
				return true
			}
		}
	}
}

func proofExists(t *testing.T, c context.Context, network test.Network, ruleID faultdetector.Rule) bool { // nolint
	mis, err := network[0].WsClient.GetOnChainProofs(c, "tendermint_getMisbehaviours")
	require.NoError(t, err)
	expected := false
	for _, p := range mis {
		rP := new(faultdetector.RawProof)
		err := rlp.DecodeBytes(p.Rawproof, rP)
		require.NoError(t, err)
		if rP.Rule == ruleID {
			expected = true
		}
	}
	return expected
}