package tenclient

import (
	"context"
	"testing"

	"github.com/clearmatics/autonity/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContractAddress(t *testing.T) {

	test := &testCase{
		name:          "no malicious",
		numValidators: 3,
		numBlocks:     5,
		txPerPeer:     1,
	}

	nodeNames := []string{"VA", "VB", "VC"}
	nodes := setupNodes(t, test, nodeNames)
	require.NotNil(t, nodes)

	addr := nodes[nodeNames[0]].listener[1].Addr().String()

	c, err := rpc.DialContext(context.Background(), "http://"+addr)
	require.NoError(t, err)
	tc := TenClient{c}

	contractAddr, err := tc.ContractAddresss(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, contractAddr)
}
