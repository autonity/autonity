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
		name, key, treasury, output string
		useFile                     bool
	}{
		{
			name:     "incorrect key in file",
			key:      "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91@",
			treasury: "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:   "Fatal: Failed to load the private key: invalid character '@' at end of key file\n",
			useFile:  true,
		},
		{
			name:     "incorrect key hex",
			key:      "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91@",
			treasury: "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:   "Fatal: Failed to parse the private key: invalid hex character '@' in private key\n",
			useFile:  false,
		},
		{
			name:     "incorrect treasury format",
			key:      "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91",
			treasury: "850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:   "Fatal: Failed to decode: hex string without 0x prefix\n",
			useFile:  true,
		},
		{
			name:     "success case",
			key:      "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91",
			treasury: "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
			output:   "Signature hex: 0x2de76b1226530c7b1fc84e3cedfce6ade854ba8e3a65c0e214695fe814f4411c3c2c884053648d9ad8dfa57855d67f5911fd18687440314cd7916bc6619eaf5c01\n",
			useFile:  true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			genEnodeProofWithExpect(t, test.key, test.treasury, test.output, test.useFile)
		})
	}
}

func genEnodeProofWithExpect(t *testing.T, key, treasury, expected string, useFile bool) {
	dir := tmpdir(t)
	var geth *testautonity
	if useFile {
		keyfile := filepath.Join(dir, "key.prv")
		if err := ioutil.WriteFile(keyfile, []byte(key), 0600); err != nil {
			t.Error(err)
		}
		geth = runAutonity(t, "genEnodeProof", "--nodekey", keyfile, treasury)
	} else {

		geth = runAutonity(t, "genEnodeProof", "--nodekeyhex", key, treasury)
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
