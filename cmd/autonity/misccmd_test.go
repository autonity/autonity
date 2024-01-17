package main

import (
	"fmt"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestGenPOP(t *testing.T) {

	tests := []struct {
		name, nodeKey, oracleKey, treasury, output      string
		useNodeKeyFile, useOracleKeyFile, legacyNodeKey bool
	}{
		{
			name:             "use correct legacy node key file to generate POP",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
			legacyNodeKey:    true,
		},
		{
			name:             "incorrect node key file",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3@",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to load the node private key: invalid character '@' at end of key file\n",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
		},
		{
			name:             "incorrect node key hex",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3@",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to parse the node private key: invalid hex character '@' in node key\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: false,
		},
		{
			name:             "incorrect oracle key file",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301@",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to load the oracle private key: invalid character '@' at end of key file\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: true,
		},
		{
			name:             "incorrect oracle key hex",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301@",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to parse the oracle private key: invalid hex character '@' in private key\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: false,
		},
		{
			name:             "incorrect treasury format",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to decode: hex string without 0x prefix\n",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
		},
		{
			name:             "success with key files",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "0xbb20d3fc5401cdf47ef59db8b581552a02023f10abd9bed4386ccdceb8522e7646bb72efcd669b98301fb040f2772c5553da50e0ee541aa84a5e67ff35e5c22c0123dce5d268ac692308a5f046c5a605da3ac66aff78810615eff008318b2e16001fe8d4e8735180cc46fe80e84a4580d80b82d51161676342cd54c6854faa823200808edced8b116bf603bbd4abf7e273093bdfd44391c8dca3c75dbe7d012b28b4cd7cdfe3f7cdf2c94b454c04ff4ab7d90dd8430597afffb6e5ff490bfebc4038af6a837ac01619069eddd6e63d06d8ba47e93d3a6845c029c8ce7904a7f43147\n",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
		},
		{
			name:             "success with key hex",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a5793cca398a63081b656790184794b9997073620e5d862750fd61b6e7fec3399ce3",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "0xbb20d3fc5401cdf47ef59db8b581552a02023f10abd9bed4386ccdceb8522e7646bb72efcd669b98301fb040f2772c5553da50e0ee541aa84a5e67ff35e5c22c0123dce5d268ac692308a5f046c5a605da3ac66aff78810615eff008318b2e16001fe8d4e8735180cc46fe80e84a4580d80b82d51161676342cd54c6854faa823200808edced8b116bf603bbd4abf7e273093bdfd44391c8dca3c75dbe7d012b28b4cd7cdfe3f7cdf2c94b454c04ff4ab7d90dd8430597afffb6e5ff490bfebc4038af6a837ac01619069eddd6e63d06d8ba47e93d3a6845c029c8ce7904a7f43147\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			genPOPWithExpect(t, test.nodeKey, test.oracleKey, test.treasury, test.output, test.useNodeKeyFile, test.useOracleKeyFile, test.legacyNodeKey)
		})
	}
}

func genPOPWithExpect(t *testing.T, nodeKey, oracleKey, treasury, expected string, useNodeKeyFile, useOracleKeyFile, legacyNodeKey bool) {
	dir := tmpdir(t)
	var geth *testautonity
	if useNodeKeyFile && useOracleKeyFile {
		nodeKeyFile := filepath.Join(dir, "nodekey.prv")
		if err := ioutil.WriteFile(nodeKeyFile, []byte(nodeKey), 0600); err != nil {
			t.Error(err)
		}
		oracleKeyFile := filepath.Join(dir, "oraclekey.prv")
		if err := ioutil.WriteFile(oracleKeyFile, []byte(oracleKey), 0600); err != nil {
			t.Error(err)
		}
		geth = runAutonity(t, "genOwnershipProof", "--autonitykeys", nodeKeyFile, "--oraclekey", oracleKeyFile, treasury)
	} else if !useNodeKeyFile && !useOracleKeyFile {
		geth = runAutonity(t, "genOwnershipProof", "--autonitykeyshex", nodeKey, "--oraclekeyhex", oracleKey, treasury)
	} else if !useNodeKeyFile && useOracleKeyFile {
		oracleKeyFile := filepath.Join(dir, "oraclekey.prv")
		if err := ioutil.WriteFile(oracleKeyFile, []byte(oracleKey), 0600); err != nil {
			t.Error(err)
		}
		geth = runAutonity(t, "genOwnershipProof", "--autonitykeyshex", nodeKey, "--oraclekey", oracleKeyFile, treasury)
	} else if useNodeKeyFile && !useOracleKeyFile {
		nodeKeyfile := filepath.Join(dir, "nodekey.prv")
		if err := ioutil.WriteFile(nodeKeyfile, []byte(nodeKey), 0600); err != nil {
			t.Error(err)
		}
		geth = runAutonity(t, "genOwnershipProof", "--autonitykeys", nodeKeyfile, "--oraclekeyhex", oracleKey, treasury)
	}
	defer geth.ExpectExit()
	if !legacyNodeKey {
		geth.Expect(expected)
		return
	}

	// wait for the child process to flush the new generated key to FS.
	time.Sleep(time.Second * 2)
	// construct expected string from generated consensus key.
	nKey, cKey, err := crypto.LoadAutonityKeys(filepath.Join(dir, "nodekey.prv"))
	require.NoError(t, err)
	oKey, err := crypto.LoadECDSA(filepath.Join(dir, "oraclekey.prv"))
	require.NoError(t, err)
	pop, err := crypto.AutonityPOPProof(nKey, oKey, treasury, cKey)
	require.NoError(t, err)
	expectedOutPut := fmt.Sprintf("%s\n", hexutil.Encode(pop))
	geth.Expect(expectedOutPut)
}

func TestGenNodeKey(t *testing.T) {
	tests := []struct {
		name, outKeyFile, output string
		writeAddr                bool
	}{
		{
			name:       "invalid key file",
			outKeyFile: "",
			writeAddr:  false,
		},
		{
			name:       "success",
			outKeyFile: "test.key",
			writeAddr:  false,
		},
		{
			name:       "success with write public address",
			outKeyFile: "test.key",
			writeAddr:  true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			genNodeKeyWithExpect(t, test.outKeyFile, test.writeAddr)
		})
	}
}

func genNodeKeyWithExpect(t *testing.T, fileName string, writeAddr bool) {
	dir := tmpdir(t)
	var geth *testautonity
	var expected string
	keyfile := filepath.Join(dir, fileName)
	if writeAddr {
		geth = runAutonity(t, "genAutonityKeys", keyfile, "--writeaddress")
	} else {
		geth = runAutonity(t, "genAutonityKeys", keyfile)
	}

	output := string(geth.Output())
	if len(fileName) != 0 {
		privateKey, consensusKey, err := crypto.LoadAutonityKeys(keyfile)
		if err != nil {
			t.Errorf("Failed to load the private key: %v", err)
			return
		}
		if writeAddr {
			expected = fmt.Sprintf("Node address: %s\nNode public key: 0x%x\nConsensus public key: %v\n",
				crypto.PubkeyToAddress(privateKey.PublicKey).String(),
				crypto.FromECDSAPub(&privateKey.PublicKey)[1:], consensusKey.PublicKey().Hex())
		}
	} else {
		expected = "Fatal: could not save key open " + keyfile + ": is a directory\n"
	}

	if output != expected {
		t.Error("output dosen't match, Actual: ", output, "\nExpected: ", expected)
	}

	defer geth.ExpectExit()
}
