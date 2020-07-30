package test

import (
	"fmt"
	"github.com/clearmatics/autonity/common/keygenerator"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/crypto"
	"math/big"
	//"sync"
	"testing"
)

func TestNetworkPermission(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// prepare chain operator
	operatorKey, err := keygenerator.Next()
	if err != nil {
		t.Fatal(err)
	}
	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	// remove node from white list hook
	// add node to white list hook
	// stop node hook
	// start node hook
	// genesis hook
	// checker hook

	// genesis hook
	unTrustLastNode := func(g *core.Genesis) *core.Genesis {
		g.Config.AutonityContractConfig.Operator = operatorAddress
		g.Alloc[operatorAddress] = core.GenesisAccount{
			Balance: big.NewInt(100000000000000000),
		}
		// remove the last node in the whitelist, take it as a untrusted peer at bootstrap phase.
		g.Config.AutonityContractConfig.Users = g.Config.AutonityContractConfig.
		return g
	}

	testCases := []*testCase{
		{
			name: "Un-permitted Lower TD Peer Cannot Connect to Network",
			numValidators: 5,
			numBlocks: 30,
			txPerPeer: 1,
			genesisHook: unTrustLastNode,
			finalAssert: nil,
		},
		{
			name: "UnPermitted High-TD Peer Finally Get Dropped From Network",
			numValidators: 4,
			numBlocks: 30,
			txPerPeer: 1,
			genesisHook: nil,
			finalAssert: nil,
		},
		{
			name: "Permitted Node Get Synced Without Available Whitelist",
			numValidators: 4,
			numBlocks: 30,
			txPerPeer: 1,
			genesisHook: nil,
			finalAssert: nil,
		},
		{
			name: "Permitted New Node Drop Un-Permitted Peer Once it Get Synced",
			numValidators: 4,
			numBlocks: 30,
			txPerPeer: 1,
			genesisHook: nil,
			finalAssert: nil,
		},
		{
			name: "Permitted New Node Drop Malicious Peer Once it Get Synced",
			numValidators: 4,
			numBlocks: 30,
			txPerPeer: 1,
			genesisHook: nil,
			finalAssert: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}

}

// todo peer with lower TD cannot connect to host if it is not presented on the host's whitelist.
func TestUnPermittedLowTDPeerCannotConnectToNetwork(t *testing.T) {

}

// todo peer with higher TD can connect to host but finally dropped if it is not presented on the host's whitelist.
func TestUnPermittedHighTDPeerFinallyDroppedFromNetwork(t *testing.T) {

}

// todo host get synced without genesis members are presented on the latest network.
func TestPermittedNodeGetSyncedWithoutAvailableWhitelist(t *testing.T) {

}

// todo host get synced and it drop those peers which are not presented on the latest whitelist.
func TestPermittedNodeDropUnPermittedPeerOnceGetSynced(t *testing.T) {

}

// todo host get synced and malicious peer who announce extremely high TD get dropped.
func TestPermittedNodeDropMaliciousPeerOnceGetSynced(t *testing.T) {

}