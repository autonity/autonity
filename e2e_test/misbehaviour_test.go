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
	committeeSize := 4
	users, err := test.Users(committeeSize, "10e18,v,1,0.0.0.0:%s,%s", 10998)
	require.NoError(t, err)
	for i := 0; i < committeeSize; i++ {
		users[i].UserType = params.UserValidator
		users[i].Stake = 1
	}

	network, err := test.NewNetworkWithMaliciousUser(users, uint8(faultdetector.PN))
	require.NoError(t, err)
	defer network.Shutdown()

	_, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	time.Sleep(90 * time.Second)
	// todo: fetch on-chain proofs of misbehaviour and verify them off chain.
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
