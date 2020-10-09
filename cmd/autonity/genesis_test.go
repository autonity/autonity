// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var customGenesisTests = []struct {
	genesis             string
	misMatchGenesis     string
	inCompatibleGenesis string
	query               string
	result              string
}{
	// Genesis file with specific chain configurations
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x1",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000001339",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
				"homesteadBlock" : 0,
				"daoForkBlock"   : 0,
				"daoForkSupport" : true,
				"homesteadBlock": 0,
				"eip150Block": 0,
				"eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"eip155Block": 0,
				"eip158Block": 0,
				"byzantiumBlock": 0,
				"constantinopleBlock": 0,
				"petersburgBlock": 0,
				"autonityContract"    : {
					"users" : [ 
						{
							"enode" : "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost:3",
							"type" : "validator",
							"stake" : 1
						}
					]
				},
				"tendermint" : {}
			}
		}`,
		misMatchGenesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x1",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000001339",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
				"homesteadBlock" : 0,
				"daoForkBlock"   : 0,
				"daoForkSupport" : true,
				"homesteadBlock": 0,
				"eip150Block": 0,
				"eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"eip155Block": 0,
				"eip158Block": 0,
				"byzantiumBlock": 0,
				"constantinopleBlock": 0,
				"petersburgBlock": 0,
				"autonityContract"    : {
					"users" : [ 
						{
							"enode" : "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost:3",
							"type" : "validator",
							"stake" : 2
						}
					]
				},
				"tendermint" : {}
			}
		}`,
		inCompatibleGenesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x1",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000001339",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
				"homesteadBlock" : 0,
				"daoForkBlock"   : 0,
				"daoForkSupport" : false,
				"homesteadBlock": 0,
				"eip150Block": 0,
				"eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"eip155Block": 0,
				"eip158Block": 0,
				"byzantiumBlock": 0,
				"constantinopleBlock": 0,
				"petersburgBlock": 0,
				"autonityContract"    : {
					"users" : [ 
						{
							"enode" : "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost:3",
							"type" : "validator",
							"stake" : 2
						}
					]
				},
				"tendermint" : {}
			}
		}`,
		query:  "eth.getBlock(0).nonce",
		result: "0x0000000000001339",
	},
}

// Tests that initializing Autonity with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
		var ipc string
		if runtime.GOOS == "windows" {
			ipc = `\\.\pipe\autonity` + strconv.Itoa(trulyRandInt(100000, 999999))
		} else {
			ws := tmpdir(t)
			defer os.RemoveAll(ws)
			ipc = filepath.Join(ws, "autonity.ipc")
		}

		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		// First create chain data with genesis, and then stop the node.
		autonity := runAutonity(t,
			"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
			"--etherbase", coinbase, "--ipcpath", ipc, "--datadir", datadir, "--genesis", json)
		waitForEndpoint(t, ipc, 5*time.Second)
		autonity.Interrupt()

		// Query the custom genesis block, check if start-up do init the wanted genesis block.
		autonity = runAutonity(t, "--nousb",
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		autonity.ExpectRegexp(tt.result)
		autonity.ExpectExit()
	}
}

// Test that if there is chain data with genesis block presented, start node with same genesis file, the node should
// perform a normal start up.
func TestStartNodeWithSameGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
		var ipc string
		if runtime.GOOS == "windows" {
			ipc = `\\.\pipe\autonity` + strconv.Itoa(trulyRandInt(100000, 999999))
		} else {
			ws := tmpdir(t)
			defer os.RemoveAll(ws)
			ipc = filepath.Join(ws, "autonity.ipc")
		}

		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}

		// First create chain data with genesis, and then stop the node.
		autonity := runAutonity(t,
			"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
			"--etherbase", coinbase, "--ipcpath", ipc, "--datadir", datadir, "--genesis", json)
		waitForEndpoint(t, ipc, 5*time.Second)
		autonity.Interrupt()

		// Make sure genesis block is created.
		autonity = runAutonity(t, "--nousb",
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		autonity.ExpectRegexp(tt.result)
		autonity.ExpectExit()

		// Then start node again with same genesis, node should perform a normal start up.
		autonity = runAutonity(t,
			"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
			"--etherbase", coinbase, "--ipcpath", ipc, "--datadir", datadir, "--genesis", json)

		defer func() {
			autonity.Interrupt()
			autonity.ExpectExit()
		}()

		// Check RPC endpoint service is loaded and RPC service is a available now.
		waitForEndpoint(t, ipc, 5*time.Second)
	}
}

// Test start node with mismatch genesis block and incompatible chain configuration, the node start up should end up,
func TestStartNodeWithIncompatibleGenesis(t *testing.T) {
	testStartNodeWithIncompatibleGenesis := func(t *testing.T, defaultGenesis string, newGenesis string) {
		for _, tt := range customGenesisTests {
			coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
			var ipc string
			if runtime.GOOS == "windows" {
				ipc = `\\.\pipe\autonity` + strconv.Itoa(trulyRandInt(100000, 999999))
			} else {
				ws := tmpdir(t)
				defer os.RemoveAll(ws)
				ipc = filepath.Join(ws, "autonity.ipc")
			}

			// Create a temporary data directory to use and inspect later
			datadir := tmpdir(t)
			defer os.RemoveAll(datadir)

			// Initialize the data directory with the default genesis block
			json := filepath.Join(datadir, "genesis.json")
			if err := ioutil.WriteFile(json, []byte(defaultGenesis), 0600); err != nil {
				t.Fatalf("failed to write genesis file: %v", err)
			}

			// First create chain data with genesis, and then stop the node.
			autonity := runAutonity(t,
				"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
				"--etherbase", coinbase, "--ipcpath", ipc, "--datadir", datadir, "--genesis", json)
			waitForEndpoint(t, ipc, 5*time.Second)
			autonity.Interrupt()

			// Make sure genesis block is created.
			autonity = runAutonity(t, "--nousb",
				"--datadir", datadir, "--maxpeers", "0", "--port", "0",
				"--nodiscover", "--nat", "none", "--ipcdisable",
				"--exec", tt.query, "console")
			autonity.ExpectRegexp(tt.result)
			autonity.ExpectExit()

			// Then start node again with a incompatible genesis file, node should not init a new genesis, and the start up should be interrupted.
			jsonV2 := filepath.Join(datadir, "genesisV2.json")
			if err := ioutil.WriteFile(jsonV2, []byte(newGenesis), 0600); err != nil {
				t.Fatalf("failed to write genesis file: %v", err)
			}
			autonity = runAutonity(t,
				"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
				"--etherbase", coinbase, "--ipcpath", ipc, "--datadir", datadir, "--genesis", jsonV2)

			defer func() {
				autonity.Interrupt()
				autonity.ExpectExit()
			}()

			// With a incompatible genesis, client should do nothing with initGenesis by exit the start up.
			autonity.ExpectExit()

			// Make sure the former genesis block is not over-write.
			autonity = runAutonity(t, "--nousb",
				"--datadir", datadir, "--maxpeers", "0", "--port", "0",
				"--nodiscover", "--nat", "none", "--ipcdisable",
				"--exec", tt.query, "console")
			autonity.ExpectRegexp(tt.result)
			autonity.ExpectExit()
		}
	}

	for _, tt := range customGenesisTests {
		testStartNodeWithIncompatibleGenesis(t, tt.genesis, tt.misMatchGenesis)
		testStartNodeWithIncompatibleGenesis(t, tt.genesis, tt.inCompatibleGenesis)
	}
}
