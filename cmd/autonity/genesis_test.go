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
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	ipcEndpoint := func(t *testing.T) (ipc string, ws string) {
		ws = tmpdir(t)
		ipc = filepath.Join(ws, "autonity.ipc")
		return ipc, ws
	}

	startNode := func(t *testing.T, genesis, ipc, ws, datadir string) *testautonity {
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(genesis), 0600); err != nil {
			t.Fatalf("failed to write genesis file: %v", err)
		}
		// Start node with genesis.
		autonity := runAutonity(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--ipcpath", ipc, "--datadir", datadir, "--genesis", json)
		return autonity
	}

	// Query the custom genesis block, check if the chain db have the wanted genesis block.
	assertMatchGenesis := func(t *testing.T, datadir string, query string, wanted string) {
		autonity := runAutonity(t, "--nousb",
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", query, "console")
		autonity.ExpectRegexp(wanted)
	}

	// start node with genesis file and check if chain db have the expected genesis.
	checkNormalStartNode := func(t *testing.T, genesis, ipc, ws, datadir, query, want string) {
		// Start node with a genesis file on the data-dir.
		autonity := startNode(t, genesis, ipc, ws, datadir)
		waitForEndpoint(t, ipc, 5*time.Second)
		// Stop node and check if genesis is generated and matched with wanted result.
		autonity.Interrupt()
		autonity.ExpectExit()
		assertMatchGenesis(t, datadir, query, want)
	}

	// start node with in-compatible genesis file, check node start should be interrupted and genesis block is expected.
	checkIncompatibleStartNode := func(t *testing.T, rawGenesis, newGenesis, ipc, ws, datadir, query, want string) {
		for _, tt := range customGenesisTests {
			// start a node with raw genesis and check genesis block is as expected.
			checkNormalStartNode(t, rawGenesis, ipc, ws, datadir, tt.query, tt.result)

			// start node on the same data dir with a brand new incompatible genesis.
			autonity := startNode(t, newGenesis, ipc, ws, datadir)

			// with a incompatible genesis, client should do nothing with initGenesis by exit process.
			autonity.ExpectExit()

			// the genesis block should be matched with the raw one.
			assertMatchGenesis(t, datadir, query, want)
		}
	}

	t.Run("Tests that starting Autonity with a custom genesis block and chain definitions works correctly", func(t *testing.T) {
		for _, tt := range customGenesisTests {
			ipc, ws := ipcEndpoint(t)
			datadir := tmpdir(t)
			defer os.RemoveAll(ws)
			defer os.RemoveAll(datadir)
			checkNormalStartNode(t, tt.genesis, ipc, ws, datadir, tt.query, tt.result)
		}
	})

	t.Run("Tests that starting Autonity with a same genesis file, node should start normally.", func(t *testing.T) {
		for _, tt := range customGenesisTests {
			ipc, ws := ipcEndpoint(t)
			datadir := tmpdir(t)
			defer os.RemoveAll(ws)
			defer os.RemoveAll(datadir)
			// Start node with a genesis file and then stop it by checking genesis is generated as expected.
			checkNormalStartNode(t, tt.genesis, ipc, ws, datadir, tt.query, tt.result)
			// Start node again with the same genesis file on the same datadir.
			checkNormalStartNode(t, tt.genesis, ipc, ws, datadir, tt.query, tt.result)
		}
	})

	t.Run("Tests that starting Autonity with a mis-matched genesis file, node should stop running and keep genesis un-touched.", func(t *testing.T) {
		for _, tt := range customGenesisTests {
			ipc, ws := ipcEndpoint(t)
			datadir := tmpdir(t)
			defer os.RemoveAll(ws)
			defer os.RemoveAll(datadir)

			checkIncompatibleStartNode(t, tt.genesis, tt.misMatchGenesis, ipc, ws, datadir, tt.query, tt.result)
		}
	})

	t.Run("Tests that starting Autonity with a incompatible genesis file, node should stop running and keep genesis un-touched.", func(t *testing.T) {
		for _, tt := range customGenesisTests {
			ipc, ws := ipcEndpoint(t)
			datadir := tmpdir(t)
			defer os.RemoveAll(ws)
			defer os.RemoveAll(datadir)

			checkIncompatibleStartNode(t, tt.genesis, tt.inCompatibleGenesis, ipc, ws, datadir, tt.query, tt.result)
		}
	})
}
