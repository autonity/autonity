package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dynamicRpcs = []string{
		"aut_retrieveContract",
		"aut_getValidators",
		"aut_getVersion",
		"aut_retrieveState",
		"aut_getStakeholders",
		"aut_getCommittee",
		"aut_getWhitelist",
		"aut_getAccountStake",
		"aut_getStake",
		"aut_getRate",
		"aut_myUserType",
		"aut_getMaxCommitteeSize",
		"aut_getCurrentCommiteeSize",
		"aut_getMinimumGasPrice",
		"aut_checkMember",
		"aut_getProposer",
		"aut_totalSupply",
		"aut_dumpEconomicsMetricData",
		"aut_enodesWhitelist",
		"aut_deployer",
		"aut_operatorAccount",
		"aut_bondingPeriod",
		"aut_committeeSize",
		"aut_contractVersion",
	}
)

func TestDynamicRpcs(t *testing.T) {

	cases := []*testCase{
		{
			numValidators: 1,
			numBlocks:     1,
			finalAssert: func(t *testing.T, validators map[string]*testNode) {
				n := validators["VA"]
				ep := n.node.HTTPEndpoint()
				// payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"%s","params":%s, "id":1}`, "aut_getAccountStake", `["0x499ea9ccfb49d1c9c207b7370d5e55cfd828858c"]`)
				for _, method := range dynamicRpcs {
					body := &rpcCall{
						Method:  method,
						Jsonrpc: "2.0",
						Id:      1,
					}
					switch method {
					case "aut_getAccountStake":
						body.Params = []string{"0x499ea9ccfb49d1c9c207b7370d5e55cfd828858c"}
					case "aut_getProposer":
						body.Params = []int{1, 0}
					case "aut_getRate":
						body.Params = []string{"0x499ea9ccfb49d1c9c207b7370d5e55cfd828858c"}
					case "aut_checkMember":
						body.Params = []string{"0x499ea9ccfb49d1c9c207b7370d5e55cfd828858c"}
					}
					payload, err := json.Marshal(body)
					require.NoError(t, err)
					println("calling", string(payload))
					callRpc(t, ep, payload)
					// resp, err := http.Post(ep, "application/json", bytes.NewBuffer([]byte(payload)))
					// assert.NoError(t, err)
					// defer resp.Body.Close()
					// respBytes, err := ioutil.ReadAll(resp.Body)
					// assert.NoError(t, err)
					// fmt.Println(string(respBytes))
				}
			},
		}}

	for _, testCase := range cases {
		testCase := testCase
		t.Run("", func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}

type rpcCall struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Id      int         `json:"id,omitempty"`
}

func callRpc(t *testing.T, ep string, payload []byte) {
	resp, err := http.Post(ep, "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	result := &rpcCall{}
	err = json.Unmarshal(respBytes, result)
	assert.NoError(t, err)
	fmt.Println(string(respBytes))
}
