// Modifications Copyright 2018 The klaytn Authors
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
//
// This file is derived from cmd/geth/consolecmd_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 eth:1.0 governance:1.0 istanbul:1.0 klay:1.0 net:1.0 personal:1.0 rpc:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "eth:1.0 klay:1.0 net:1.0 rpc:1.0 web3:1.0"
)

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	// Start a klay console, make sure it's cleaned up and terminate the console
	klay := runKlay(t,
		"klay-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none",
		"console")

	// Gather all the infos the welcome message needs to contain
	klay.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	klay.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	klay.SetTemplateFunc("gover", runtime.Version)
	klay.SetTemplateFunc("klayver", func() string { return params.VersionWithCommit(gitCommit) })
	klay.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	klay.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	klay.Expect(`
Welcome to the Klaytn JavaScript console!

instance: Klaytn/{{klayver}}/{{goos}}-{{goarch}}/{{gover}}
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	klay.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	// Configure the instance for IPC attachement
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\klay` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "klay.ipc")
	}
	// Note: we need --shh because testAttachWelcome checks for default
	// list of ipc modules and shh is included there.
	klay := runKlay(t,
		"klay-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none", "--ipcpath", ipc)

	time.Sleep(5 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, klay, "ipc:"+ipc, ipcAPIs)

	klay.Interrupt()
	klay.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	klay := runKlay(t,
		"klay-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none", "--rpc", "--rpcport", port)

	time.Sleep(5 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, klay, "http://localhost:"+port, httpAPIs)

	klay.Interrupt()
	klay.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	klay := runKlay(t,
		"klay-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none", "--ws", "--wsport", port)

	time.Sleep(5 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, klay, "ws://localhost:"+port, httpAPIs)

	klay.Interrupt()
	klay.ExpectExit()
}

func testAttachWelcome(t *testing.T, klay *testklay, endpoint, apis string) {
	// Attach to a running Klaytn node and terminate immediately
	attach := runKlay(t, "klay-test", "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("klayver", func() string { return params.VersionWithCommit(gitCommit) })
	attach.SetTemplateFunc("rewardbase", func() string { return klay.Rewardbase })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return klay.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the Klaytn JavaScript console!

instance: Klaytn/{{klayver}}/{{goos}}-{{goarch}}/{{gover}}{{if ipc}}
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
