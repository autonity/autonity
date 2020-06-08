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
)

var customGenesisTests = []struct {
	genesis string
	query   string
	result  string
}{
	// Plain genesis file without anything extra
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x20000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000000042",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"	 : {
				"tendermint": {
				  "block-period": 1
				},
				"autonityContract": {
				  "deployer": "0xe6191d6b8649b8cfbaf05ba416a35a1c5c868c40",
				  "bytecode": "",
				  "abi": "",
				  "minGasPrice": 1,
				  "operator": "0xb6aa12260f602cdc814c0e960f7dc8cfabf0558d",
				  "users": [
					{
					  "enode": "enode://f423075988abe9a14f44725c3f4c0eff330d0c94708d0236d3088da01ec6dec89278784bcba0d5116c7dc5e28c8acae7f14be0930bc389b80a9c117c8a25ea38@127.0.0.1:6789",
					  "type": "validator",
					  "stake": 1
					}
				  ]
				}
			}
		}`,
		query:  "eth.getBlock(0).nonce",
		result: "0x0000000000000042",
	},
	// Genesis file with an empty chain configuration (ensure missing fields work)
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x20000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000000042",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
				"tendermint": {
				  "block-period": 1
				},
				"autonityContract": {
				  "deployer": "0xe6191d6b8649b8cfbaf05ba416a35a1c5c868c40",
				  "bytecode": "",
				  "abi": "",
				  "minGasPrice": 1,
				  "operator": "0xb6aa12260f602cdc814c0e960f7dc8cfabf0558d",
				  "users": [
					{
					  "enode": "enode://f423075988abe9a14f44725c3f4c0eff330d0c94708d0236d3088da01ec6dec89278784bcba0d5116c7dc5e28c8acae7f14be0930bc389b80a9c117c8a25ea38@127.0.0.1:6789",
					  "type": "validator",
					  "stake": 1
					}
				  ]
				}
			}
		}`,
		query:  "eth.getBlock(0).nonce",
		result: "0x0000000000000042",
	},
	// Genesis file with specific chain configurations
	{
		genesis: `{
			"alloc"      : {},
			"coinbase"   : "0x0000000000000000000000000000000000000000",
			"difficulty" : "0x20000",
			"extraData"  : "",
			"gasLimit"   : "0x2fefd8",
			"nonce"      : "0x0000000000000042",
			"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
				"homesteadBlock" : 314,
				"daoForkBlock"   : 141,
				"daoForkSupport" : true,
				"tendermint": {
				  "block-period": 1
				},
				"autonityContract": {
				  "deployer": "0xe6191d6b8649b8cfbaf05ba416a35a1c5c868c40",
				  "bytecode": "",
				  "abi": "",
				  "minGasPrice": 1,
				  "operator": "0xb6aa12260f602cdc814c0e960f7dc8cfabf0558d",
				  "users": [
					{
					  "enode": "enode://f423075988abe9a14f44725c3f4c0eff330d0c94708d0236d3088da01ec6dec89278784bcba0d5116c7dc5e28c8acae7f14be0930bc389b80a9c117c8a25ea38@127.0.0.1:6789",
					  "type": "validator",
					  "stake": 1
					}
				  ]
				}
			}
		}`,
		query:  "eth.getBlock(0).nonce",
		result: "0x0000000000000042",
	},
}

// Tests that initializing Autonity with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runAutonity(t, "--datadir", datadir, "init", json).WaitExit()

		// Query the custom genesis block
		autonity := runAutonity(t,
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		autonity.ExpectRegexp(tt.result)
		autonity.ExpectExit()
	}
}
