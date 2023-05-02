package main

import (
	"fmt"
	"github.com/autonity/autonity/crypto"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenEnodeProof(t *testing.T) {
	tests := []struct {
		name, nodeKey, oracleKey, treasury, output string
		useNodeKeyFile, useOracleKeyFile           bool
	}{
		{
			name:             "incorrect node key file",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579@",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to load the node private key: invalid character '@' at end of key file\n",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
		},
		{
			name:             "incorrect node key hex",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579@",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to parse the node private key: invalid hex character '@' in private key\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: false,
		},
		{
			name:             "incorrect oracle key file",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301@",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to load the oracle private key: invalid character '@' at end of key file\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: true,
		},
		{
			name:             "incorrect oracle key hex",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301@",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to parse the oracle private key: invalid hex character '@' in private key\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: false,
		},
		{
			name:             "incorrect treasury format",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Fatal: Failed to decode: hex string without 0x prefix\n",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
		},
		{
			name:             "success with key files",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Signature hex: 0xbb20d3fc5401cdf47ef59db8b581552a02023f10abd9bed4386ccdceb8522e7646bb72efcd669b98301fb040f2772c5553da50e0ee541aa84a5e67ff35e5c22c0123dce5d268ac692308a5f046c5a605da3ac66aff78810615eff008318b2e16001fe8d4e8735180cc46fe80e84a4580d80b82d51161676342cd54c6854faa823200\n",
			useNodeKeyFile:   true,
			useOracleKeyFile: true,
		},
		{
			name:             "success with key hex",
			nodeKey:          "f1ab65d8d07ab6a7a2ab8419fa5bbaf8938f45556387d43f3f15967bc599a579",
			oracleKey:        "198227888008a50b57bfb4d70ef5c4a3ef085538b148842fe3628b9005d66301",
			treasury:         "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:           "Signature hex: 0xbb20d3fc5401cdf47ef59db8b581552a02023f10abd9bed4386ccdceb8522e7646bb72efcd669b98301fb040f2772c5553da50e0ee541aa84a5e67ff35e5c22c0123dce5d268ac692308a5f046c5a605da3ac66aff78810615eff008318b2e16001fe8d4e8735180cc46fe80e84a4580d80b82d51161676342cd54c6854faa823200\n",
			useNodeKeyFile:   false,
			useOracleKeyFile: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			genOwnershipProofWithExpect(t, test.nodeKey, test.oracleKey, test.treasury, test.output, test.useNodeKeyFile, test.useOracleKeyFile)
		})
	}
}

func genOwnershipProofWithExpect(t *testing.T, nodeKey, oracleKey, treasury, expected string, useNodeKeyFile, useOracleKeyFile bool) {
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
		geth = runAutonity(t, "genOwnershipProof", "--nodekey", nodeKeyFile, "--oraclekey", oracleKeyFile, treasury)
	} else if !useNodeKeyFile && !useOracleKeyFile {
		geth = runAutonity(t, "genOwnershipProof", "--nodekeyhex", nodeKey, "--oraclekeyhex", oracleKey, treasury)
	} else if !useNodeKeyFile && useOracleKeyFile {
		oracleKeyFile := filepath.Join(dir, "oraclekey.prv")
		if err := ioutil.WriteFile(oracleKeyFile, []byte(oracleKey), 0600); err != nil {
			t.Error(err)
		}
		geth = runAutonity(t, "genOwnershipProof", "--nodekeyhex", nodeKey, "--oraclekey", oracleKeyFile, treasury)
	} else if useNodeKeyFile && !useOracleKeyFile {
		nodeKeyfile := filepath.Join(dir, "nodekey.prv")
		if err := ioutil.WriteFile(nodeKeyfile, []byte(nodeKey), 0600); err != nil {
			t.Error(err)
		}
		geth = runAutonity(t, "genOwnershipProof", "--nodekey", nodeKeyfile, "--oraclekeyhex", oracleKey, treasury)
	}
	defer geth.ExpectExit()
	geth.Expect(expected)
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
		geth = runAutonity(t, "genNodeKey", keyfile, "--writeaddress")
	} else {
		geth = runAutonity(t, "genNodeKey", keyfile)
	}

	output := string(geth.Output())
	if len(fileName) != 0 {
		privateKey, err := crypto.LoadECDSA(keyfile)
		if err != nil {
			t.Errorf("Failed to load the private key: %v", err)
			return
		}
		if writeAddr {
			expected = fmt.Sprintf("%x\n", crypto.FromECDSAPub(&privateKey.PublicKey)[1:])
		}
	} else {
		expected = "Fatal: could not save key open " + keyfile + ": is a directory\n"
	}

	if output != expected {
		t.Error("output dosen't match, Actual: ", output, "\nExpected: ", expected)
	}

	defer geth.ExpectExit()
}
