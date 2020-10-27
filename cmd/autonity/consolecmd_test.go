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
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/clearmatics/autonity/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 eth:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 tendermint:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "eth:1.0 net:1.0 rpc:1.0 web3:1.0"
)

var genesis = `{
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
				"tendermint" : {"block-period" : 1}
			}
		}`

func tmpDataDirWithGenesisFile(t *testing.T) (dir string, genesisFile string) {
	dir = tmpdir(t)
	genesisFile = filepath.Join(dir, "genesis.json")
	if err := ioutil.WriteFile(genesisFile, []byte(genesis), 0600); err != nil {
		t.Fatalf("failed to write genesis file: %v", err)
	}
	return dir, genesisFile
}

// Tests that, by using console sub-command to launch a client without a genesis block and genesis configuration, node
// embedded within a console can not be started up properly.
func TestConsoleStartWithoutGenesis(t *testing.T) {
	autonity := runAutonity(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none", "console")
	// Verify the actual welcome message to the required template
	autonity.Expect(`
Fatal: Error starting protocol stack: DB has no genesis block and there is no genesis file set by user
`)
	autonity.ExpectExit()
}

// Tests that, with a genesis configuration, a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	dir, jsonFile := tmpDataDirWithGenesisFile(t)
	defer os.RemoveAll(dir)

	// Start a autonity client via console sub command.
	autonity := runAutonity(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--datadir", dir, "--genesis", jsonFile, "console")

	// Gather all the infos the welcome message needs to contain
	autonity.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	autonity.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	autonity.SetTemplateFunc("gover", runtime.Version)
	autonity.SetTemplateFunc("autonityver", func() string { return params.VersionWithCommit("", "") })
	autonity.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)") })
	autonity.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	autonity.Expect(`
The embedded Autonity Console is no longer supported, use it a your own risk.
Consider the Autonity Node.js Console as replacement.

instance: Autonity/v{{autonityver}}/{{goos}}-{{goarch}}/{{gover}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	autonity.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\autonity` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "autonity.ipc")
	}
	dir, jsonFile := tmpDataDirWithGenesisFile(t)
	defer os.RemoveAll(dir)
	autonity := runAutonity(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--ipcpath", ipc, "--datadir", dir, "--genesis", jsonFile)

	defer func() {
		autonity.Interrupt()
		autonity.ExpectExit()
	}()

	waitForEndpoint(t, ipc, 3*time.Second)
	testAttachWelcome(t, autonity, "ipc:"+ipc, ipcAPIs)

}

func TestHTTPAttachWelcome(t *testing.T) {
	dir, jsonFile := tmpDataDirWithGenesisFile(t)
	defer os.RemoveAll(dir)
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	autonity := runAutonity(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--rpc", "--rpcport", port, "--datadir", dir, "--genesis", jsonFile)
	defer func() {
		autonity.Interrupt()
		autonity.ExpectExit()
	}()

	endpoint := "http://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 3*time.Second)
	testAttachWelcome(t, autonity, endpoint, httpAPIs)
}

func TestWSAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	dir, jsonFile := tmpDataDirWithGenesisFile(t)
	defer os.RemoveAll(dir)
	autonity := runAutonity(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--ws", "--wsport", port, "--datadir", dir, "--genesis", jsonFile)
	defer func() {
		autonity.Interrupt()
		autonity.ExpectExit()
	}()

	endpoint := "ws://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 3*time.Second)
	testAttachWelcome(t, autonity, endpoint, httpAPIs)
}

func testAttachWelcome(t *testing.T, autonity *testautonity, endpoint, apis string) {
	// Attach to a running autonity note and terminate immediately
	attach := runAutonity(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("autonityver", func() string { return params.VersionWithCommit("", "") })
	attach.SetTemplateFunc("etherbase", func() string { return autonity.Etherbase })
	attach.SetTemplateFunc("niltime", func() string {
		return time.Unix(0, 0).Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")
	})
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return autonity.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
The embedded Autonity Console is no longer supported, use it a your own risk.
Consider the Autonity Node.js Console as replacement.

instance: Autonity/v{{autonityver}}/{{goos}}-{{goarch}}/{{gover}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
