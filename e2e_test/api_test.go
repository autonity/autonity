package e2e

import (
	"encoding/json"
	"testing"

	"github.com/autonity/autonity/rpc"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/eth"
	"github.com/autonity/autonity/p2p"
)

func TestAPI_AcnPeers(t *testing.T) {
	network, err := NewNetwork(t, 7, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)
	err = network.WaitToMineNBlocks(5, 30, false)
	require.NoError(t, err)
	node := network[0]
	apis := node.Eth.APIs()
	var autContractAPI *eth.AutonityContractAPI
	for _, api := range apis {
		if api.Namespace == "aut" {
			autContractAPI = api.Service.(*eth.AutonityContractAPI)
			break
		}
	}
	require.NotNil(t, autContractAPI)
	peers := autContractAPI.AcnPeers()
	require.Equal(t, len(network)-1, len(peers))

	// test json marshaling
	jsonPayload, err := json.Marshal(peers)
	require.NoError(t, err)

	// test json un-marshaling
	var result []*p2p.PeerInfo
	err = json.Unmarshal(jsonPayload, &result)
	require.NoError(t, err)
	for _, peer := range result {
		require.True(t, peer.Network.Trusted)
	}
}

func getAutApi(node *Node) *eth.AutonityContractAPI {
	var autContractAPI *eth.AutonityContractAPI

	apis := node.Eth.APIs()
	for _, api := range apis {
		if api.Namespace == "aut" {
			autContractAPI = api.Service.(*eth.AutonityContractAPI)
			break
		}
	}
	return autContractAPI
}

func TestAPI_Versions(t *testing.T) {
	network, err := NewNetwork(t, 2, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown(t)

	err = network.WaitToMineNBlocks(5, 30, false)
	require.NoError(t, err)

	autContractAPI := getAutApi(network[0])
	require.NotNil(t, autContractAPI)

	versions, err := autContractAPI.ProtocolContractsVersions(nil)
	require.NoError(t, err)
	for _, version := range versions {
		t.Log(version)
		// all base protocol contract versions should be defined
		require.NotEqual(t, "", version.Version)
	}

	pendingBlockNumber := rpc.PendingBlockNumber
	pendingVersions, err := autContractAPI.ProtocolContractsVersions(&pendingBlockNumber)
	require.NoError(t, err)
	require.Equal(t, versions, pendingVersions)

	latestBlockNumber := rpc.LatestBlockNumber
	latestVersions, err := autContractAPI.ProtocolContractsVersions(&latestBlockNumber)
	require.NoError(t, err)
	require.Equal(t, versions, latestVersions)

	earliestBlockNumber := rpc.EarliestBlockNumber
	earliestVersions, err := autContractAPI.ProtocolContractsVersions(&earliestBlockNumber)
	require.NoError(t, err)
	require.Equal(t, versions, earliestVersions)
}
