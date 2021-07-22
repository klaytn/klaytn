// Modifications Copyright 2019 The klaytn Authors
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
// This file is derived from cmd/geth/run_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/console"
	metricutils "github.com/klaytn/klaytn/metrics/utils"
	"github.com/klaytn/klaytn/node"
	"gopkg.in/urfave/cli.v1"
)

func tmpdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "klay-test")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

type testklay struct {
	*utils.TestCmd

	// template variables for expect
	Datadir    string
	Rewardbase string
}

var (
	// The app that holds all commands and flags.
	app = utils.NewApp(GetGitCommit(), "the Klaytn command line interface")

	// flags that configure the node
	nodeFlags = CommonNodeFlags

	rpcFlags = CommonRPCFlags
)

func contains(list []cli.Flag, item cli.Flag) bool {
	for _, flag := range list {
		if flag.GetName() == item.GetName() {
			return true
		}
	}
	return false
}

func union(list1, list2 []cli.Flag) []cli.Flag {
	for _, item := range list2 {
		if !contains(list1, item) {
			list1 = append(list1, item)
		}
	}
	return list1
}

func init() {
	// Initialize the CLI app and start Klay
	app.Action = RunKlaytnNode
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2019 The klaytn Authors"
	app.Commands = []cli.Command{
		// See chaincmd.go:
		InitCommand,

		// See accountcmd.go
		AccountCommand,

		// See consolecmd.go:
		GetConsoleCommand(nodeFlags, rpcFlags),
		AttachCommand,

		// See versioncmd.go:
		VersionCommand,

		// See dumpconfigcmd.go:
		GetDumpConfigCommand(nodeFlags, rpcFlags),
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, ConsoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = union(app.Flags, KCNFlags)
	app.Flags = union(app.Flags, KPNFlags)
	app.Flags = union(app.Flags, KENFlags)
	app.Flags = union(app.Flags, KSCNFlags)
	app.Flags = union(app.Flags, KSPNFlags)
	app.Flags = union(app.Flags, KSENFlags)

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		logDir := (&node.Config{DataDir: utils.MakeDataDir(ctx)}).ResolvePath("logs")
		debug.CreateLogDir(logDir)
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		metricutils.StartMetricCollectionAndExport(ctx)
		utils.SetupNetwork(ctx)
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}

	// Run the app if we've been exec'd as "klay-test" in runKlay.
	reexec.Register("klay-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
	reexec.Register("klay-test-flag", func() {
		app.Action = RunTestKlaytnNode
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func TestMain(m *testing.M) {
	// check if we have been reexec'd
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

// spawns klay with the given command line args. If the args don't set --datadir, the
// child g gets a temporary data directory.
func runKlay(t *testing.T, name string, args ...string) *testklay {
	tt := &testklay{}
	tt.TestCmd = utils.NewTestCmd(t, tt)
	for i, arg := range args {
		switch {
		case arg == "-datadir" || arg == "--datadir":
			if i < len(args)-1 {
				tt.Datadir = args[i+1]
			}
		case arg == "-rewardbase" || arg == "--rewardbase":
			if i < len(args)-1 {
				tt.Rewardbase = args[i+1]
			}
		}
	}
	if tt.Datadir == "" {
		tt.Datadir = tmpdir(t)
		tt.Cleanup = func() { os.RemoveAll(tt.Datadir) }
		args = append([]string{"-datadir", tt.Datadir}, args...)
		// Remove the temporary datadir if something fails below.
		defer func() {
			if t.Failed() {
				tt.Cleanup()
			}
		}()
	}

	// Boot "klay". This actually runs the test binary but the TestMain
	// function will prevent any tests from running.
	tt.Run(name, args...)

	return tt
}

func RunTestKlaytnNode(ctx *cli.Context) error {
	fullNode := MakeFullNode(ctx)
	fullNode.Wait()
	return nil
}
