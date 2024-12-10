package e2e

import (
	"encoding/json"
	"testing"

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
	var result  []*p2p.PeerInfo
	err = json.Unmarshal(jsonPayload, &result)
	require.NoError(t, err)
	for _, peer := range result {
		require.True(t, peer.Network.Trusted)
	}
}
