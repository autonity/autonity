package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dynamicRpcs = []string{
		"aut_operatorAccount",
		"aut_committeeSize",
		"aut_deployer",
		"aut_allowance",
		"aut_getNewContract",
		"aut_getState",
		"aut_getVersion",
		"aut_getCommittee",
		"aut_getValidators",
		"aut_getStakeholders",
		"aut_getWhitelist",
		"aut_balanceOf",
		"aut_totalSupply",
		"aut_getUser",
		"aut_getMaxCommitteeSize",
		"aut_getMinimumGasPrice",
		"aut_getProposer",
		"aut_dumpEconomicMetrics",
	}
)

// This test makes an rpc call for each element in dynamicRpcs. It only checks
// that a result is returned and no error is encountered, it is assumed that
// the functionality of these calls is tested in more detail elsewhere.
func TestDynamicRpcs(t *testing.T) {
	tc := &testCase{
		numValidators: 1,
		numBlocks:     1,
		finalAssert: func(t *testing.T, validators map[string]*testNode) {
			n := validators["VA"]
			validatorAddress := n.EthAddress().String()
			ep := n.node.HTTPEndpoint()
			for _, method := range dynamicRpcs {
				body := &rpcCall{
					Method:  method,
					Jsonrpc: "2.0",
					Id:      1,
				}
				switch method {
				case "aut_allowance":
					body.Params = []string{validatorAddress, validatorAddress}
				case "aut_balanceOf":
					body.Params = []string{validatorAddress}
				case "aut_getUser":
					body.Params = []string{validatorAddress}
				case "aut_getProposer":
					height := 1
					round := 0
					body.Params = []int{height, round}
				}
				payload, err := json.Marshal(body)
				require.NoError(t, err)
				respBytes := callRPC(t, ep, payload)
				responseMap := make(map[string]interface{})
				err = json.Unmarshal(respBytes, &responseMap)
				require.NoError(t, err)

				// Check that there was no error and that a result was returned.
				assert.NotNil(t, responseMap["result"])
				assert.Nil(t, responseMap["error"])
			}
		},
	}
	runTest(t, tc)
}

type rpcCall struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"` // nolint
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Id      int         `json:"id,omitempty"` // nolint
}

func callRPC(t *testing.T, ep string, payload []byte) []byte {
	resp, err := http.Post(ep, "application/json", bytes.NewBuffer(payload)) // nolint gosec complains about variable url
	assert.NoError(t, err)
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	result := &rpcCall{}
	err = json.Unmarshal(respBytes, result)
	assert.NoError(t, err)
	return respBytes
}
