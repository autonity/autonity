package contracts

import (
	"bytes"
	"encoding/json"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	getterRPCs = []string{
		"aut_getOperator",
		"aut_getMaxCommitteeSize",
		"aut_deployer",
		"aut_allowance",
		"aut_name",
		"aut_symbol",
		"aut_getNewContract",
		"aut_getVersion",
		"aut_getCommittee",
		"aut_getValidators",
		"aut_balanceOf",
		"aut_totalSupply",
		"aut_getMaxCommitteeSize",
		"aut_getMinimumBaseFee",
		"aut_getNewContract",
	}
)

// This test checks that those getter functions can be accessed via the client's HTTP RPC calls.
func TestACGetterRPCs(t *testing.T) {
	network, err := e2e.NewNetwork(t, 1, "10e18,v,1,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)
	defer network.Shutdown()

	// wait for the consensus engine to work.
	network.WaitToMineNBlocks(2, 10, false)

	validator := network[0].Address.String()
	ep := network[0].HTTPEndpoint()

	for _, method := range getterRPCs {
		body := &rpcCall{
			Method:  method,
			Jsonrpc: "2.0",
			Id:      1,
		}
		switch method {
		case "aut_allowance":
			body.Params = []string{validator, validator}
		case "aut_balanceOf":
			body.Params = []string{validator}
		case "aut_getUser":
			body.Params = []string{validator}
		}
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		respBytes := callRPC(t, ep, payload)
		responseMap := make(map[string]interface{})
		err = json.Unmarshal(respBytes, &responseMap)
		require.NoError(t, err)
		// Check that there was no error and that a result was returned.
		require.NotNil(t, responseMap["result"])
		require.Nil(t, responseMap["error"])
	}
}

type rpcCall struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"` // nolint
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Id      int         `json:"id,omitempty"` // nolint
}

func callRPC(t *testing.T, ep string, payload []byte) []byte {
	resp, err := http.Post(ep, "application/json", bytes.NewBuffer(payload)) // nolint gosec complains about variable url
	require.NoError(t, err)
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	result := &rpcCall{}
	err = json.Unmarshal(respBytes, result)
	require.NoError(t, err)
	return respBytes
}
