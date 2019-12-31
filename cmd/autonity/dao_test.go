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
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/params"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

// Genesis block for nodes which don't care about the DAO fork (i.e. not configured)
var daoOldGenesis = `{
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
		"autonityContract" : {
			            "users": [
                {
                    "enode": "enode://b0ad42df49340361a0d9989a124acf58ec05e20e7bc35b994e8626c8f548a992d5d10f4d0acbedbd25ff30796590a55bcb090c3bed2f00c2f59f5f47391cbe6e@127.0.0.1:5000",
                    "address": "0xaD8dD560065cc44DC8011F80B4c4F974Cd5E8350",
                    "type": "validator",
                    "stake": 2
                },
                {
                    "enode": "enode://0964c84500648ce4ef06366beef0fee73fa7b56aed5a14c2861a238079dd188acc8c6da9da743d20216000413e3ec5bcd8a743d96b58ebff467019515aeb04d6@127.0.0.1:5001",
                    "address": "0xbEA4FF331eC70ff0506e716F673F2553574e22e8",
                    "type": "validator",
                    "stake": 1
                },
                {
                    "enode": "enode://535304a4bade7aed3120c21267baded1bae027aa9064db99abab5b326700710f1b32e9281c3f92d5fd5a4a218e40aecaa8f52e8e2632447683ea7d7651a07947@127.0.0.1:5002",
                    "address": "0x4d1deA1EC4dB139c5D19880aC7791dd54Cb5Aa41",
                    "type": "validator",
                    "stake": 1
                },
                {
                    "enode": "enode://237392d2d24688ea1e1f65f7284cd53825f55cf2ec2962b410bb559ddc3a57d1f71517dbe19fd8baf1c187e7e8e5912960ba09115de028709d46ec82f5c49849@127.0.0.1:5003",
                    "address": "0x4f2c44DFb93465027Ed042F4d38b41c8c5Cf662C",
                    "type": "validator",
                    "stake": 1
                }
            ]},
		"tendermint" : {}
	}
}`

// Genesis block for nodes which actively oppose the DAO fork
var daoNoForkGenesis = `{
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
		"daoForkBlock"   : 314,
		"daoForkSupport" : false,
		"autonityContract" : {
			            "users": [
                {
                    "enode": "enode://b0ad42df49340361a0d9989a124acf58ec05e20e7bc35b994e8626c8f548a992d5d10f4d0acbedbd25ff30796590a55bcb090c3bed2f00c2f59f5f47391cbe6e@127.0.0.1:5000",
                    "address": "0xaD8dD560065cc44DC8011F80B4c4F974Cd5E8350",
                    "type": "validator",
                    "stake": 2
                },
                {
                    "enode": "enode://0964c84500648ce4ef06366beef0fee73fa7b56aed5a14c2861a238079dd188acc8c6da9da743d20216000413e3ec5bcd8a743d96b58ebff467019515aeb04d6@127.0.0.1:5001",
                    "address": "0xbEA4FF331eC70ff0506e716F673F2553574e22e8",
                    "type": "validator",
                    "stake": 1
                },
                {
                    "enode": "enode://535304a4bade7aed3120c21267baded1bae027aa9064db99abab5b326700710f1b32e9281c3f92d5fd5a4a218e40aecaa8f52e8e2632447683ea7d7651a07947@127.0.0.1:5002",
                    "address": "0x4d1deA1EC4dB139c5D19880aC7791dd54Cb5Aa41",
                    "type": "validator",
                    "stake": 1
                },
                {
                    "enode": "enode://237392d2d24688ea1e1f65f7284cd53825f55cf2ec2962b410bb559ddc3a57d1f71517dbe19fd8baf1c187e7e8e5912960ba09115de028709d46ec82f5c49849@127.0.0.1:5003",
                    "address": "0x4f2c44DFb93465027Ed042F4d38b41c8c5Cf662C",
                    "type": "validator",
                    "stake": 1
                }
            ]},
		"tendermint" : {}
	}
}`

// Genesis block for nodes which actively support the DAO fork
var daoProForkGenesis = `{
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
		"daoForkBlock"   : 314,
		"daoForkSupport" : true,
		"autonityContract" : {
			            "users": [
                {
                    "enode": "enode://b0ad42df49340361a0d9989a124acf58ec05e20e7bc35b994e8626c8f548a992d5d10f4d0acbedbd25ff30796590a55bcb090c3bed2f00c2f59f5f47391cbe6e@127.0.0.1:5000",
                    "address": "0xaD8dD560065cc44DC8011F80B4c4F974Cd5E8350",
                    "type": "validator",
                    "stake": 2
                },
                {
                    "enode": "enode://0964c84500648ce4ef06366beef0fee73fa7b56aed5a14c2861a238079dd188acc8c6da9da743d20216000413e3ec5bcd8a743d96b58ebff467019515aeb04d6@127.0.0.1:5001",
                    "address": "0xbEA4FF331eC70ff0506e716F673F2553574e22e8",
                    "type": "validator",
                    "stake": 1
                },
                {
                    "enode": "enode://535304a4bade7aed3120c21267baded1bae027aa9064db99abab5b326700710f1b32e9281c3f92d5fd5a4a218e40aecaa8f52e8e2632447683ea7d7651a07947@127.0.0.1:5002",
                    "address": "0x4d1deA1EC4dB139c5D19880aC7791dd54Cb5Aa41",
                    "type": "validator",
                    "stake": 1
                },
                {
                    "enode": "enode://237392d2d24688ea1e1f65f7284cd53825f55cf2ec2962b410bb559ddc3a57d1f71517dbe19fd8baf1c187e7e8e5912960ba09115de028709d46ec82f5c49849@127.0.0.1:5003",
                    "address": "0x4f2c44DFb93465027Ed042F4d38b41c8c5Cf662C",
                    "type": "validator",
                    "stake": 1
                }
            ]
		},
		"tendermint" : {}
	}
}`

var daoGenesisHash = common.HexToHash("c626d8be5d39286154d44074a2f0a31192792385e89d6f1f60dea77be6b99ce5")
var daoGenesisForkBlock = big.NewInt(314)

// TestDAOForkBlockNewChain tests that the DAO hard-fork number and the nodes support/opposition is correctly
// set in the database after various initialization procedures and invocations.
func TestDAOForkBlockNewChain(t *testing.T) {
	for i, arg := range []struct {
		genesis     string
		expectBlock *big.Int
		expectVote  bool
	}{
		// Test DAO Default Mainnet
		{"", params.MainnetChainConfig.DAOForkBlock, true},
		// test DAO Init Old Privnet
		{daoOldGenesis, nil, false},
		// test DAO Default No Fork Privnet
		{daoNoForkGenesis, daoGenesisForkBlock, false},
		// test DAO Default Pro Fork Privnet
		{daoProForkGenesis, daoGenesisForkBlock, true},
	} {
		testDAOForkBlockNewChain(t, i, arg.genesis, arg.expectBlock, arg.expectVote)
	}
}

func testDAOForkBlockNewChain(t *testing.T, test int, genesis string, expectBlock *big.Int, expectVote bool) {
	// Create a temporary data directory to use and inspect later
	datadir := tmpdir(t)
	defer os.RemoveAll(datadir)

	// Start a Autonity instance with the requested flags set and immediately terminate
	if genesis != "" {
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", test, err)
		}
		runAutonity(t, "--datadir", datadir, "init", json).WaitExit()
	} else {
		// Force chain initialization
		args := []string{"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none", "--ipcdisable", "--datadir", datadir}
		autonity := runAutonity(t, append(args, []string{"--exec", "2+2", "console"}...)...)
		autonity.WaitExit()
	}
	// Retrieve the DAO config flag from the database
	path := filepath.Join(datadir, "autonity", "chaindata")
	db, err := rawdb.NewLevelDBDatabase(path, 0, 0, "")
	if err != nil {
		t.Fatalf("test %d: failed to open test database: %v", test, err)
	}
	defer db.Close()

	genesisHash := common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	if genesis != "" {
		genesisHash = daoGenesisHash
	}
	config := rawdb.ReadChainConfig(db, genesisHash)
	if config == nil {
		t.Errorf("test %d: failed to retrieve chain config: %v", test, err)
		return // we want to return here, the other checks can't make it past this point (nil panic).
	}
	// Validate the DAO hard-fork block number against the expected value
	if config.DAOForkBlock == nil {
		if expectBlock != nil {
			t.Errorf("test %d: dao hard-fork block mismatch: have nil, want %v", test, expectBlock)
		}
	} else if expectBlock == nil {
		t.Errorf("test %d: dao hard-fork block mismatch: have %v, want nil", test, config.DAOForkBlock)
	} else if config.DAOForkBlock.Cmp(expectBlock) != 0 {
		t.Errorf("test %d: dao hard-fork block mismatch: have %v, want %v", test, config.DAOForkBlock, expectBlock)
	}
	if config.DAOForkSupport != expectVote {
		t.Errorf("test %d: dao hard-fork support mismatch: have %v, want %v", test, config.DAOForkSupport, expectVote)
	}
}
