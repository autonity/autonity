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
	"bytes"
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/params"
	"github.com/stretchr/testify/require"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 eth:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 tendermint:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "eth:1.0 net:1.0 rpc:1.0 web3:1.0"
)

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	testInitGenesis(t, datadir)

	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"

	// Start a autonity console, make sure it's cleaned up and terminate the console
	autonity := runAutonity(t,
		"--datadir", datadir, "--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "console", "--nodekey", "/home/pierspowlesland/projects/autonity/userkey", "--mine")

	userKey, err := crypto.LoadECDSA("/home/pierspowlesland/projects/autonity/userkey")
	require.NoError(t, err)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey).Hex()
	userAddr = strings.ToLower(userAddr)

	// Gather all the infos the welcome message needs to contain
	autonity.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	autonity.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	autonity.SetTemplateFunc("gover", runtime.Version)
	autonity.SetTemplateFunc("autonityver", func() string { return params.VersionWithCommit("", "") })
	autonity.SetTemplateFunc("etherbase", func() string { return userAddr })
	autonity.SetTemplateFunc("datadir", func() string { return datadir })
	autonity.SetTemplateFunc("apis", func() string { return ipcAPIs })

	templateSource := `
Welcome to the Autonity JavaScript console!

instance: Autonity/v{{autonityver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{etherbase}}
at block: .*
 datadir: {{datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`
	tpl := template.Must(template.New("").Funcs(autonity.Func).Parse(templateSource))

	wantbuf := new(bytes.Buffer)
	if err := tpl.Execute(wantbuf, autonity.Data); err != nil {
		panic(err)
	}
	// Verify the actual welcome message to the required template
	autonity.ExpectRegexp(wantbuf.String())
	autonity.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	testInitGenesis(t, datadir)
	// Configure the instance for IPC attachement
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\autonity` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "autonity.ipc")
	}

	autonity := runAutonity(t,
		"--datadir", datadir, "--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--ipcpath", ipc, "--nodekey", "/home/pierspowlesland/projects/autonity/userkey", "--mine")

	waitForEndpoint(t, ipc, 3*time.Second)
	userKey, err := crypto.LoadECDSA("/home/pierspowlesland/projects/autonity/userkey")
	require.NoError(t, err)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey).Hex()
	userAddr = strings.ToLower(userAddr)
	testAttachWelcome(t, autonity, "ipc:"+ipc, ipcAPIs, userAddr)

	autonity.Interrupt()
	autonity.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	testInitGenesis(t, datadir)
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	autonity := runAutonity(t,
		"--datadir", datadir, "--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--rpc", "--rpcport", port, "--nodekey", "/home/pierspowlesland/projects/autonity/userkey", "--mine")

	endpoint := "http://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 3*time.Second)

	userKey, err := crypto.LoadECDSA("/home/pierspowlesland/projects/autonity/userkey")
	require.NoError(t, err)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey).Hex()
	userAddr = strings.ToLower(userAddr)
	testAttachWelcome(t, autonity, endpoint, httpAPIs, userAddr)

	autonity.Interrupt()
	autonity.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	testInitGenesis(t, datadir)
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	autonity := runAutonity(t,
		"--datadir", datadir, "--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--ws", "--wsport", port, "--nodekey", "/home/pierspowlesland/projects/autonity/userkey", "--mine")

	endpoint := "ws://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 3*time.Second)

	userKey, err := crypto.LoadECDSA("/home/pierspowlesland/projects/autonity/userkey")
	require.NoError(t, err)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey).Hex()
	userAddr = strings.ToLower(userAddr)

	testAttachWelcome(t, autonity, endpoint, httpAPIs, userAddr)

	autonity.Interrupt()
	autonity.ExpectExit()
}

func testAttachWelcome(t *testing.T, autonity *testautonity, endpoint, apis, userAddr string) {
	// Attach to a running autonity note and terminate immediately
	attach := runAutonity(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("autonityver", func() string { return params.VersionWithCommit("", "") })
	attach.SetTemplateFunc("etherbase", func() string { return userAddr })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return autonity.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })
	templateSource := `
Welcome to the Autonity JavaScript console!

instance: Autonity/v{{autonityver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{etherbase}}
at block: .*{{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`

	tpl := template.Must(template.New("").Funcs(attach.Func).Parse(templateSource))

	wantbuf := new(bytes.Buffer)
	if err := tpl.Execute(wantbuf, attach.Data); err != nil {
		panic(err)
	}
	// Verify the actual welcome message to the required template
	attach.ExpectRegexp(wantbuf.String())
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
